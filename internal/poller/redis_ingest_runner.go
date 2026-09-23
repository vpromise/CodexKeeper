package poller

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"cpa-usage-keeper/internal/timeutil"
	"github.com/sirupsen/logrus"
)

type RedisIngestRunnerConfig struct {
	BatchSize     int
	RetryInterval time.Duration
}

// RedisIngestRunner uses the current CodexProxy protocol: subscribe first, then
// drain the fixed usage queue. It never probes or consumes a legacy endpoint.
type RedisIngestRunner struct {
	subscribeSource UsageSubscriptionSource
	redisSource     UsagePullSource
	writer          RedisInboxWriter
	controlObserver RedisControlMessageObserver
	config          RedisIngestRunnerConfig
	sleep           func(context.Context, time.Duration) bool
	mu              sync.Mutex
	status          Status
}

func NewRedisIngestRunner(sub UsageSubscriptionSource, redis UsagePullSource, writer RedisInboxWriter, cfg RedisIngestRunnerConfig) *RedisIngestRunner {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10000
	}
	if cfg.RetryInterval <= 0 {
		cfg.RetryInterval = time.Second
	}
	return &RedisIngestRunner{subscribeSource: sub, redisSource: redis, writer: writer, config: cfg, sleep: sleepContext}
}

func (r *RedisIngestRunner) SetControlMessageObserver(observer RedisControlMessageObserver) {
	if r != nil {
		r.controlObserver = observer
	}
}

func (r *RedisIngestRunner) Run(ctx context.Context) error {
	if r == nil || r.subscribeSource == nil || r.redisSource == nil || r.writer == nil {
		return fmt.Errorf("usage subscription, backlog source, and inbox writer are required")
	}
	r.mu.Lock()
	if r.status.Running {
		r.mu.Unlock()
		return ErrSyncAlreadyRunning
	}
	r.status.Running = true
	r.mu.Unlock()
	defer func() {
		r.mu.Lock()
		r.status.Running = false
		r.status.SyncRunning = false
		r.mu.Unlock()
	}()
	backoff := NewRedisIngestBackoff(r.config.RetryInterval, 30*time.Second)
	for ctx.Err() == nil {
		r.record("connecting", nil, false)
		sub, err := r.subscribeSource.Subscribe(ctx)
		if err == nil && sub == nil {
			err = fmt.Errorf("usage subscription is nil")
		}
		if err == nil {
			backoff.Reset()
			if r.controlObserver != nil {
				r.controlObserver.NotifyIngestConnected()
			}
			err = r.consume(ctx, sub)
			_ = sub.Close()
		}
		if ctx.Err() != nil {
			return nil
		}
		r.record("reconnecting", err, false)
		if r.controlObserver != nil {
			r.controlObserver.MarkRefreshPollingRequired("usage subscription disconnected")
		}
		if !r.sleep(ctx, backoff.NextDelay()) {
			return nil
		}
	}
	return nil
}

func (r *RedisIngestRunner) consume(ctx context.Context, sub UsageSubscription) error {
	r.record("subscribe_backfill", nil, true)
	// Once subscribed, CodexProxy publishes new events to the subscription and
	// keeps only pre-subscription events in the usage backlog.
	for ctx.Err() == nil {
		messages, err := r.redisSource.Pull(ctx)
		if err != nil {
			return fmt.Errorf("read usage backlog: %w", err)
		}
		if err := r.persist(ctx, RedisIngestSourceRedisPull, messages); err != nil {
			return err
		}
		if len(messages) < r.config.BatchSize {
			break
		}
	}
	for ctx.Err() == nil {
		r.record("subscribing", nil, true)
		first, err := sub.Receive(ctx)
		if err != nil {
			return err
		}
		batch := []string{first}
		window, cancel := context.WithTimeout(ctx, time.Second)
		var receiveErr error
		for len(batch) < r.config.BatchSize {
			message, err := sub.Receive(window)
			if err != nil {
				receiveErr = err
				break
			}
			batch = append(batch, message)
		}
		cancel()
		if err := r.persist(ctx, RedisIngestSourceSubscribe, batch); err != nil {
			return err
		}
		if receiveErr != nil && !errors.Is(receiveErr, context.DeadlineExceeded) {
			return receiveErr
		}
	}
	return ctx.Err()
}

// persist retains an already consumed batch until SQLite accepts it. A failed
// write must not cause another LPOP or silently discard subscription events.
func (r *RedisIngestRunner) persist(ctx context.Context, source string, messages []string) error {
	if len(messages) == 0 {
		return nil
	}
	receivedAt := timeutil.NormalizeStorageTime(time.Now())
	backoff := NewRedisIngestBackoff(r.config.RetryInterval, 30*time.Second)
	for {
		writeCtx := ctx
		cancel := func() {}
		if ctx.Err() != nil {
			// Finish the current batch before App closes the database during shutdown.
			writeCtx, cancel = context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		}
		_, err := r.writer.Insert(writeCtx, source, messages, receivedAt)
		cancel()
		if err == nil {
			return nil
		}
		r.record("inbox_write_failed", err, false)
		if ctx.Err() != nil {
			return fmt.Errorf("flush usage inbox: %w", err)
		}
		if !r.sleep(ctx, backoff.NextDelay()) {
			continue
		}
	}
}

func (r *RedisIngestRunner) record(state string, err error, available bool) {
	r.mu.Lock()
	previous := r.status.LastStatus
	r.status.LastStatus = state
	r.status.LastRunAt = time.Now()
	r.status.SyncRunning = available
	r.status.LastError = ""
	if err != nil {
		r.status.LastError = err.Error()
	}
	r.mu.Unlock()
	if err != nil && previous != state {
		logrus.WithError(err).WithField("state", state).Warn("usage ingestion interrupted")
	}
}

func (r *RedisIngestRunner) Status() Status {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status
}

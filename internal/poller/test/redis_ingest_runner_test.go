package poller_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"cpa-usage-keeper/internal/poller"

	"github.com/sirupsen/logrus"
)

type fakeSubscribeSource struct {
	sub poller.UsageSubscription
	err error
}

func (s fakeSubscribeSource) Subscribe(context.Context) (poller.UsageSubscription, error) {
	return s.sub, s.err
}

type blockingSubscription struct {
	messages chan string
}

func (s *blockingSubscription) Receive(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case message := <-s.messages:
		return message, nil
	}
}

func (s *blockingSubscription) Close() error { return nil }

type failingSubscription struct {
	err error
}

func (s failingSubscription) Receive(context.Context) (string, error) { return "", s.err }

func (s failingSubscription) Close() error { return nil }

type fakePullSource struct {
	mu      sync.Mutex
	batches [][]string
	errs    []error
	err     error
	calls   int
}

func (s *fakePullSource) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func (s *fakePullSource) Pull(context.Context) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if len(s.errs) > 0 {
		err := s.errs[0]
		s.errs = s.errs[1:]
		if err != nil {
			return nil, err
		}
	}
	if s.err != nil {
		return nil, s.err
	}
	if len(s.batches) == 0 {
		return nil, nil
	}
	batch := s.batches[0]
	s.batches = s.batches[1:]
	return batch, nil
}

type fakeNamedPullSource struct {
	*fakePullSource
	sourceName string
}

func (s *fakeNamedPullSource) SourceName() string { return s.sourceName }

type fakeInboxInsert struct {
	source   string
	messages []string
}

type fakeInboxWriter struct {
	attempts chan fakeInboxInsert
	ch       chan fakeInboxInsert
	err      error
}

func newFakeInboxWriter() *fakeInboxWriter {
	return &fakeInboxWriter{attempts: make(chan fakeInboxInsert, 10), ch: make(chan fakeInboxInsert, 10)}
}

func (w *fakeInboxWriter) Insert(_ context.Context, source string, messages []string, _ time.Time) (int, error) {
	entry := fakeInboxInsert{source: source, messages: append([]string(nil), messages...)}
	w.attempts <- entry
	if w.err != nil {
		return 0, w.err
	}
	w.ch <- entry
	return len(messages), nil
}

func (w *fakeInboxWriter) waitForAttempt(t *testing.T) fakeInboxInsert {
	t.Helper()
	select {
	case entry := <-w.attempts:
		return entry
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for write attempt")
		return fakeInboxInsert{}
	}
}

func (w *fakeInboxWriter) waitForInsert(t *testing.T) fakeInboxInsert {
	t.Helper()
	select {
	case entry := <-w.ch:
		return entry
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for insert")
		return fakeInboxInsert{}
	}
}

func waitForStatus(t *testing.T, runner *poller.RedisIngestRunner, match func(poller.Status) bool) poller.Status {
	t.Helper()
	deadline := time.After(time.Second)
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for {
		status := runner.Status()
		if match(status) {
			return status
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for status, got: %+v", status)
			return status
		case <-tick.C:
		}
	}
}

type lockedLogBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedLogBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedLogBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func capturePollerLogs(t *testing.T, level logrus.Level) *lockedLogBuffer {
	t.Helper()
	logs := &lockedLogBuffer{}
	previousOutput := logrus.StandardLogger().Out
	previousFormatter := logrus.StandardLogger().Formatter
	previousLevel := logrus.GetLevel()
	logrus.SetOutput(logs)
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logrus.SetLevel(level)
	t.Cleanup(func() {
		logrus.SetOutput(previousOutput)
		logrus.SetFormatter(previousFormatter)
		logrus.SetLevel(previousLevel)
	})
	return logs
}

func runNativeIngest(t *testing.T, runner *poller.RedisIngestRunner) context.CancelFunc {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runner.Run(ctx) }()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(2 * time.Second):
			t.Error("runner did not stop")
		}
	})
	return cancel
}

func TestNativeSubscriberDrainsFullBacklogAndReceivesLiveEvents(t *testing.T) {
	sub := &blockingSubscription{messages: make(chan string, 2)}
	pull := &fakePullSource{batches: [][]string{{"first", "second"}, {"third"}}}
	writer := newFakeInboxWriter()
	runner := poller.NewRedisIngestRunner(fakeSubscribeSource{sub: sub}, pull, writer, poller.RedisIngestRunnerConfig{BatchSize: 2, RetryInterval: time.Millisecond})
	runNativeIngest(t, runner)
	if got := writer.waitForInsert(t); got.source != poller.RedisIngestSourceRedisPull || len(got.messages) != 2 {
		t.Fatalf("first backlog: %+v", got)
	}
	if got := writer.waitForInsert(t); len(got.messages) != 1 || got.messages[0] != "third" {
		t.Fatalf("remaining backlog: %+v", got)
	}
	sub.messages <- "live-one"
	sub.messages <- "live-two"
	if got := writer.waitForInsert(t); got.source != poller.RedisIngestSourceSubscribe || len(got.messages) != 2 {
		t.Fatalf("live batch: %+v", got)
	}
	waitForStatus(t, runner, func(s poller.Status) bool { return s.SyncRunning })
}

type retryInboxWriter struct {
	mu       sync.Mutex
	calls    int
	messages [][]string
	success  chan struct{}
}

func (w *retryInboxWriter) Insert(_ context.Context, _ string, messages []string, _ time.Time) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.calls++
	w.messages = append(w.messages, append([]string(nil), messages...))
	if w.calls == 1 {
		return 0, errors.New("database is locked")
	}
	select {
	case w.success <- struct{}{}:
	default:
	}
	return len(messages), nil
}

func TestNativeSubscriberRetriesConsumedBacklogWithoutAnotherPop(t *testing.T) {
	pull := &fakePullSource{batches: [][]string{{"retained"}}}
	writer := &retryInboxWriter{success: make(chan struct{}, 1)}
	runner := poller.NewRedisIngestRunner(fakeSubscribeSource{sub: &blockingSubscription{messages: make(chan string)}}, pull, writer, poller.RedisIngestRunnerConfig{BatchSize: 2, RetryInterval: time.Millisecond})
	runNativeIngest(t, runner)
	select {
	case <-writer.success:
	case <-time.After(time.Second):
		t.Fatal("batch was not retried")
	}
	writer.mu.Lock()
	defer writer.mu.Unlock()
	if writer.calls != 2 || writer.messages[0][0] != "retained" || writer.messages[1][0] != "retained" {
		t.Fatalf("retry batches: %+v", writer.messages)
	}
	if pull.callCount() != 1 {
		t.Fatal("a failed inbox write consumed another backlog batch")
	}
}

func TestNativeSubscriberFlushesPartialBatchOnShutdown(t *testing.T) {
	sub := &blockingSubscription{messages: make(chan string, 1)}
	writer := newFakeInboxWriter()
	runner := poller.NewRedisIngestRunner(fakeSubscribeSource{sub: sub}, &fakePullSource{}, writer, poller.RedisIngestRunnerConfig{BatchSize: 10})
	cancel := runNativeIngest(t, runner)
	waitForStatus(t, runner, func(s poller.Status) bool { return s.LastStatus == "subscribing" })
	sub.messages <- "last-event"
	// Wait until Receive consumes the first event, then cancel the batch window.
	deadline := time.Now().Add(time.Second)
	for len(sub.messages) > 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()
	if got := writer.waitForInsert(t); len(got.messages) != 1 || got.messages[0] != "last-event" {
		t.Fatalf("shutdown batch: %+v", got)
	}
}

func TestNativeSubscriberDoesNotPopUntilSubscriptionConnects(t *testing.T) {
	pull := &fakePullSource{}
	runner := poller.NewRedisIngestRunner(fakeSubscribeSource{err: errors.New("offline")}, pull, newFakeInboxWriter(), poller.RedisIngestRunnerConfig{RetryInterval: time.Millisecond})
	runNativeIngest(t, runner)
	waitForStatus(t, runner, func(s poller.Status) bool { return s.LastError != "" })
	if pull.callCount() != 0 {
		t.Fatal("subscription failure triggered destructive fallback polling")
	}
}

type reconnectSubscribeSource struct {
	calls int
	next  poller.UsageSubscription
}

func (s *reconnectSubscribeSource) Subscribe(context.Context) (poller.UsageSubscription, error) {
	s.calls++
	if s.calls == 1 {
		return failingSubscription{err: errors.New("connection lost")}, nil
	}
	return s.next, nil
}

func TestNativeSubscriberDrainsBacklogAgainAfterReconnect(t *testing.T) {
	sub := &blockingSubscription{messages: make(chan string)}
	source := &reconnectSubscribeSource{next: sub}
	pull := &fakePullSource{batches: [][]string{nil, {"during-disconnect"}}}
	writer := newFakeInboxWriter()
	runner := poller.NewRedisIngestRunner(source, pull, writer, poller.RedisIngestRunnerConfig{BatchSize: 2, RetryInterval: time.Millisecond})
	runNativeIngest(t, runner)
	got := writer.waitForInsert(t)
	if got.source != poller.RedisIngestSourceRedisPull || len(got.messages) != 1 || got.messages[0] != "during-disconnect" {
		t.Fatalf("reconnect backlog: %+v", got)
	}
	if pull.callCount() != 2 {
		t.Fatalf("expected a backlog drain on each connection")
	}
}

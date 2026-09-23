package poller

import (
	"context"
	"fmt"

	"cpa-usage-keeper/internal/cpa"
)

// RedisPullSource drains only the current CodexProxy usage backlog.
type RedisPullSource struct{ client *cpa.RedisQueueClient }

func NewRedisPullSource(opts cpa.RedisQueueOptions) *RedisPullSource {
	opts.QueueKey = cpa.ManagementUsageQueueKey
	return &RedisPullSource{client: cpa.NewRedisQueueClientWithOptions(opts)}
}

func (s *RedisPullSource) Pull(ctx context.Context) ([]string, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("usage backlog client is nil")
	}
	return s.client.PopUsage(ctx)
}

func (*RedisPullSource) SourceName() string { return RedisIngestSourceRedisPull }

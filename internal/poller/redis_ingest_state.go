package poller

import "cpa-usage-keeper/internal/cpa"

const (
	RedisIngestSourceSubscribe = "redis_subscribe:" + cpa.ManagementUsageSubscribeChannel
	RedisIngestSourceRedisPull = "redis_pull:" + cpa.ManagementUsageQueueKey
)

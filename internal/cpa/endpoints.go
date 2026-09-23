package cpa

const (
	cpaManagementAuthFilesEndpoint       = "/v0/management/auth-files"
	cpaManagementAuthFilesStatusEndpoint = "/v0/management/auth-files/status"
	cpaManagementAPIKeysEndpoint         = "/v0/management/api-keys"
	cpaManagementCodexAPIKeyEndpoint     = "/v0/management/codex-api-key"
	cpaManagementClaudeAPIKeyEndpoint    = "/v0/management/claude-api-key"
	cpaManagementAPICallEndpoint         = "/v0/management/api-call"
	cpaManagementResetQuotaEndpoint      = "/v0/management/reset-quota"
	cpaManagementRequestLogByIDEndpoint  = "/v0/management/request-log-by-id"
	cpaModelsEndpoint                    = "/v1/models"

	cpaManagementRedisNetwork       = "tcp"
	ManagementRedisDefaultPort      = "8317"
	ManagementRedisAuthCommand      = "AUTH"
	ManagementRedisPopCommand       = "LPOP"
	ManagementRedisSubscribeCommand = "SUBSCRIBE"
	ManagementUsageQueueKey         = "usage"
	ManagementUsageSubscribeChannel = "usage"
	// ManagementErrorsSubscribeChannel 是 CPA 只广播、不提供 LPOP 补偿的凭证错误 channel。
	ManagementErrorsSubscribeChannel = "errors"
	ManagementUsageQueueMaxBatchSize = 10000
)

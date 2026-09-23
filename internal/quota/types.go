package quota

import (
	"context"

	"cpa-usage-keeper/internal/cpa/dto/apicall"
	"cpa-usage-keeper/internal/entities"
)

type ManagementAPICaller interface {
	CallManagementAPI(context.Context, apicall.Request) (*apicall.Response, error)
}

type ManagementClient interface {
	ManagementAPICaller
	ResetQuota(context.Context, string) error
}

type ProviderInput struct {
	Identity entities.UsageIdentity
}

type ProviderOutput struct {
	Provider string
	Result   any
}

type ProviderResetOutput struct {
	Code           string `json:"code,omitempty"`
	WindowsReset   int    `json:"windowsReset,omitempty"`
	RecoveryFailed bool   `json:"recoveryFailed,omitempty"`
}

type ProviderResetter interface {
	Reset(context.Context, ProviderInput) (ProviderResetOutput, error)
}

type ProviderResetCreditsOutput struct {
	AvailableCount *int                        `json:"availableCount"`
	Credits        []CodexRateLimitResetCredit `json:"credits"`
}

type ProviderResetCreditLister interface {
	ListResetCredits(context.Context, ProviderInput) (ProviderResetCreditsOutput, error)
}

type QuotaWindow struct {
	Duration *float64 `json:"duration,omitempty"`
	Unit     string   `json:"unit,omitempty"`
	Seconds  *int64   `json:"seconds,omitempty"`
}

type SubscriptionInfo struct {
	Provider string `json:"provider"`
	Plan     string `json:"plan"`
	TierID   string `json:"tierId,omitempty"`
	TierName string `json:"tierName,omitempty"`
}

type QuotaRow struct {
	Key               string       `json:"key"`
	Label             string       `json:"label,omitempty"`
	Scope             string       `json:"scope,omitempty"`
	Metric            string       `json:"metric,omitempty"`
	GroupKey          string       `json:"groupKey,omitempty"`
	GroupLabel        string       `json:"groupLabel,omitempty"`
	GroupDescription  string       `json:"groupDescription,omitempty"`
	Used              *float64     `json:"used,omitempty"`
	Limit             *float64     `json:"limit,omitempty"`
	Remaining         *float64     `json:"remaining,omitempty"`
	UsedPercent       *float64     `json:"usedPercent,omitempty"`
	RemainingFraction *float64     `json:"remainingFraction,omitempty"`
	Allowed           *bool        `json:"allowed,omitempty"`
	LimitReached      *bool        `json:"limitReached,omitempty"`
	Window            *QuotaWindow `json:"window,omitempty"`
	ResetAt           string       `json:"resetAt,omitempty"`
	ResetAfterSeconds *int64       `json:"resetAfterSeconds,omitempty"`
	WindowUsageTokens *int64       `json:"window_usage_tokens,omitempty"`
	WindowUsageCost   *float64     `json:"window_usage_cost,omitempty"`
}

type CodexUsageWindow struct {
	// UsedPercent 是上游返回的已用小数百分比；零值只有 HasUsedPercent=true 时才是明确事实。
	UsedPercent float64 `json:"usedPercent,omitempty"`
	// LimitWindowSeconds 是上游原始窗口秒数；历史周期身份不把它限制为当前已知枚举。
	LimitWindowSeconds int64 `json:"limitWindowSeconds,omitempty"`
	// ResetAfterSeconds 是相对观察时间的剩余秒数；明确零值表示立即重置。
	ResetAfterSeconds int64 `json:"resetAfterSeconds,omitempty"`
	// ResetAt 是上游绝对 Unix 秒；存在时优先于相对倒计时。
	ResetAt int64 `json:"resetAt,omitempty"`
	// WindowUsageTokens 是现有 quota cache 展示字段，历史表永远不复制或持久化它。
	WindowUsageTokens *int64 `json:"window_usage_tokens,omitempty"`
	// WindowUsageCost 是现有 quota cache 展示字段，历史表永远不复制或持久化它。
	WindowUsageCost *float64 `json:"window_usage_cost,omitempty"`
	// HasUsedPercent 区分上游缺失字段和明确 0% 已用，仅供内部历史提取使用。
	HasUsedPercent bool `json:"-"`
	// HasLimitWindowSeconds 区分上游缺失字段和无效零秒窗口，仅供内部历史提取使用。
	HasLimitWindowSeconds bool `json:"-"`
	// HasResetAfterSeconds 区分相对重置字段缺失和明确零秒，仅供内部历史提取使用。
	HasResetAfterSeconds bool `json:"-"`
	// HasResetAt 区分绝对重置字段缺失和无效零值，仅供内部历史提取使用。
	HasResetAt bool `json:"-"`
}

type CodexRateLimitInfo struct {
	Allowed         *bool             `json:"allowed,omitempty"`
	LimitReached    *bool             `json:"limitReached,omitempty"`
	PrimaryWindow   *CodexUsageWindow `json:"primaryWindow,omitempty"`
	SecondaryWindow *CodexUsageWindow `json:"secondaryWindow,omitempty"`
}

type CodexAdditionalRateLimit struct {
	LimitName      string              `json:"limitName,omitempty"`
	MeteredFeature string              `json:"meteredFeature,omitempty"`
	RateLimit      *CodexRateLimitInfo `json:"rateLimit,omitempty"`
}

type CodexRateLimitResetCredits struct {
	AvailableCount *int `json:"availableCount,omitempty"`
}

type CodexRateLimitResetCredit struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	GrantedAt string `json:"grantedAt,omitempty"`
	ExpiresAt string `json:"expiresAt"`
}

type CodexUsagePayload struct {
	PlanType              string                      `json:"planType,omitempty"`
	RateLimit             *CodexRateLimitInfo         `json:"rateLimit,omitempty"`
	CodeReviewRateLimit   *CodexRateLimitInfo         `json:"codeReviewRateLimit,omitempty"`
	AdditionalRateLimits  []CodexAdditionalRateLimit  `json:"additionalRateLimits,omitempty"`
	RateLimitResetCredits *CodexRateLimitResetCredits `json:"rateLimitResetCredits,omitempty"`
}

type ClaudeUsageWindow struct {
	Utilization float64 `json:"utilization,omitempty"`
	ResetsAt    string  `json:"resetsAt,omitempty"`
}

type ClaudeExtraUsage struct {
	IsEnabled    bool     `json:"isEnabled,omitempty"`
	MonthlyLimit float64  `json:"monthlyLimit,omitempty"`
	UsedCredits  float64  `json:"usedCredits,omitempty"`
	Utilization  *float64 `json:"utilization,omitempty"`
}

type ClaudeUsagePayload struct {
	FiveHour          *ClaudeUsageWindow `json:"fiveHour,omitempty"`
	SevenDay          *ClaudeUsageWindow `json:"sevenDay,omitempty"`
	SevenDayOAuthApps *ClaudeUsageWindow `json:"sevenDayOauthApps,omitempty"`
	SevenDayOpus      *ClaudeUsageWindow `json:"sevenDayOpus,omitempty"`
	SevenDaySonnet    *ClaudeUsageWindow `json:"sevenDaySonnet,omitempty"`
	SevenDayCowork    *ClaudeUsageWindow `json:"sevenDayCowork,omitempty"`
	IguanaNecktie     *ClaudeUsageWindow `json:"iguanaNecktie,omitempty"`
	ExtraUsage        *ClaudeExtraUsage  `json:"extraUsage,omitempty"`
}

type ClaudeProfileAccount struct {
	UUID         string `json:"uuid,omitempty"`
	FullName     string `json:"fullName,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Email        string `json:"email,omitempty"`
	HasClaudeMax *bool  `json:"hasClaudeMax,omitempty"`
	HasClaudePro *bool  `json:"hasClaudePro,omitempty"`
}

type ClaudeProfileOrganization struct {
	UUID                 string `json:"uuid,omitempty"`
	Name                 string `json:"name,omitempty"`
	OrganizationType     string `json:"organizationType,omitempty"`
	BillingType          string `json:"billingType,omitempty"`
	RateLimitTier        string `json:"rateLimitTier,omitempty"`
	HasExtraUsageEnabled bool   `json:"hasExtraUsageEnabled,omitempty"`
	SubscriptionStatus   string `json:"subscriptionStatus,omitempty"`
}

type ClaudeProfileResponse struct {
	Account      *ClaudeProfileAccount      `json:"account,omitempty"`
	Organization *ClaudeProfileOrganization `json:"organization,omitempty"`
}

type CodexResult struct {
	Usage *CodexUsagePayload `json:"usage"`
}

type ClaudeResult struct {
	Usage   *ClaudeUsagePayload    `json:"usage"`
	Profile *ClaudeProfileResponse `json:"profile"`
}

type ProviderHandler interface {
	Check(context.Context, ProviderInput) (ProviderOutput, error)
}

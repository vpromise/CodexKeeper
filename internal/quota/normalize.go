package quota

import (
	"strings"
	"time"

	"cpa-usage-keeper/internal/timeutil"
)

const (
	quotaWindowFiveHourSeconds     int64 = 5 * 60 * 60
	quotaWindowSevenDaySeconds     int64 = 7 * 24 * 60 * 60
	quotaWindowThirtyDaySeconds    int64 = 30 * 24 * 60 * 60
	quotaWindowAverageMonthSeconds int64 = 365 * 24 * 60 * 60 / 12
)

func NormalizeQuotaRows(output ProviderOutput) []QuotaRow {
	// 不在 provider 层强行统一原始结构，只在出口处转换为前端展示需要的 quota rows。
	switch result := output.Result.(type) {
	case CodexResult:
		return normalizeCodexQuotaRows(result)
	case *CodexResult:
		if result == nil {
			return nil
		}
		return normalizeCodexQuotaRows(*result)
	case ClaudeResult:
		return normalizeClaudeQuotaRows(result)
	case *ClaudeResult:
		if result == nil {
			return nil
		}
		return normalizeClaudeQuotaRows(*result)
	default:
		return nil
	}
}

func normalizeClaudeQuotaRows(result ClaudeResult) []QuotaRow {
	if result.Usage == nil {
		return nil
	}
	rows := make([]QuotaRow, 0, 8)
	rows = appendClaudeWindowQuotaRow(rows, "five_hour", "5h", "window", result.Usage.FiveHour)
	rows = appendClaudeWindowQuotaRow(rows, "seven_day", "Weekly", "window", result.Usage.SevenDay)
	rows = appendClaudeWindowQuotaRow(rows, "seven_day_oauth_apps", "7d OAuth Apps", "window", result.Usage.SevenDayOAuthApps)
	rows = appendClaudeWindowQuotaRow(rows, "seven_day_opus", "7d Opus", "model", result.Usage.SevenDayOpus)
	rows = appendClaudeWindowQuotaRow(rows, "seven_day_sonnet", "7d Sonnet", "model", result.Usage.SevenDaySonnet)
	rows = appendClaudeWindowQuotaRow(rows, "seven_day_cowork", "7d Cowork", "window", result.Usage.SevenDayCowork)
	rows = appendClaudeWindowQuotaRow(rows, "iguana_necktie", "Iguana Necktie", "window", result.Usage.IguanaNecktie)
	if result.Usage.ExtraUsage != nil {
		rows = append(rows, QuotaRow{
			Key:         "extra_usage",
			Label:       "Extra Usage",
			Scope:       "extra_usage",
			Used:        floatPtr(result.Usage.ExtraUsage.UsedCredits),
			Limit:       floatPtr(result.Usage.ExtraUsage.MonthlyLimit),
			UsedPercent: result.Usage.ExtraUsage.Utilization,
			Allowed:     boolPtr(result.Usage.ExtraUsage.IsEnabled),
		})
	}
	return rows
}

func appendClaudeWindowQuotaRow(rows []QuotaRow, key string, label string, scope string, window *ClaudeUsageWindow) []QuotaRow {
	if window == nil {
		return rows
	}
	row := QuotaRow{
		Key:         key,
		Label:       label,
		Scope:       scope,
		UsedPercent: floatPtr(window.Utilization),
		ResetAt:     window.ResetsAt,
	}
	// Claude 只给官方语义明确的 5h 会话窗口和 seven_day 系列补 seconds，其它未知 key 不猜测。
	if key == "five_hour" {
		row.Window = &QuotaWindow{Seconds: intPtr(quotaWindowFiveHourSeconds)}
	} else if strings.HasPrefix(key, "seven_day") {
		row.Window = &QuotaWindow{Seconds: intPtr(quotaWindowSevenDaySeconds)}
	}
	return append(rows, row)
}

func normalizeCodexQuotaRows(result CodexResult) []QuotaRow {
	// Codex 根据 limit_window_seconds 明确区分 5h/Weekly/Monthly；未知窗口回退到 primary/secondary 角色。
	if result.Usage == nil {
		return nil
	}
	rows := make([]QuotaRow, 0, 4+len(result.Usage.AdditionalRateLimits)*2)
	rows = appendCodexWindowQuotaRows(rows, "rate_limit", "5h", "Weekly", "window", "", result.Usage.RateLimit)
	rows = appendCodexWindowQuotaRows(rows, "code_review_rate_limit", "Code Review 5h", "Code Review Weekly", "code_review", "", result.Usage.CodeReviewRateLimit)
	for _, additional := range result.Usage.AdditionalRateLimits {
		// 主限额之外的 code review / spark 等窗口也保留为 extra quota，避免丢失上游数据。
		metric := additional.MeteredFeature
		if metric == "" {
			metric = additional.LimitName
		}
		primaryLabel := additional.LimitName + " 5h"
		secondaryLabel := additional.LimitName + " Weekly"
		rows = appendCodexWindowQuotaRows(rows, "additional_rate_limits."+additional.LimitName, primaryLabel, secondaryLabel, "additional", metric, additional.RateLimit)
	}
	return rows
}

func appendCodexWindowQuotaRows(rows []QuotaRow, keyPrefix string, primaryLabel string, secondaryLabel string, scope string, metric string, info *CodexRateLimitInfo) []QuotaRow {
	if info == nil {
		return rows
	}
	rows = appendCodexWindowQuotaRow(rows, keyPrefix+".primary_window", primaryLabel, scope, metric, info, info.PrimaryWindow)
	rows = appendCodexWindowQuotaRow(rows, keyPrefix+".secondary_window", secondaryLabel, scope, metric, info, info.SecondaryWindow)
	return rows
}

func codexWindowLabel(key string, label string, seconds int64) string {
	switch seconds {
	case quotaWindowFiveHourSeconds:
		return codexKnownWindowLabel(label, "5h")
	case quotaWindowSevenDaySeconds:
		return codexKnownWindowLabel(label, "Weekly")
	case quotaWindowThirtyDaySeconds, quotaWindowAverageMonthSeconds:
		return codexKnownWindowLabel(label, "Monthly")
	}
	return codexRoleWindowLabel(key, label)
}

func codexKnownWindowLabel(label string, replacement string) string {
	if label == "5h" || label == "Weekly" || label == "Monthly" || label == "Window" {
		return replacement
	}
	for _, candidate := range []string{"5h", "Weekly", "Monthly", "Window"} {
		if strings.Contains(label, candidate) {
			return strings.Replace(label, candidate, replacement, 1)
		}
	}
	return label
}

func codexRoleWindowLabel(key string, label string) string {
	role := codexWindowRole(key)
	if role == "" {
		return label
	}
	fallback := codexPrefixedWindowRole(key, role)
	if label == "" || label == "5h" || label == "Weekly" || label == "Monthly" || label == "Window" {
		return fallback
	}
	for _, candidate := range []string{"5h", "Weekly", "Monthly", "Window"} {
		if strings.Contains(label, candidate) {
			return strings.Replace(label, candidate, role, 1)
		}
	}
	return fallback
}

func codexWindowRole(key string) string {
	if strings.HasSuffix(key, ".primary_window") {
		return "Primary"
	}
	if strings.HasSuffix(key, ".secondary_window") {
		return "Secondary"
	}
	return ""
}

func codexPrefixedWindowRole(key string, role string) string {
	if strings.HasPrefix(key, "code_review_rate_limit.") {
		return "Code Review " + role
	}
	if strings.HasPrefix(key, "additional_rate_limits.") {
		name := strings.TrimPrefix(key, "additional_rate_limits.")
		name = strings.TrimSuffix(name, ".primary_window")
		name = strings.TrimSuffix(name, ".secondary_window")
		if name != "" {
			return name + " " + role
		}
	}
	return role
}

func appendCodexWindowQuotaRow(rows []QuotaRow, key string, label string, scope string, metric string, info *CodexRateLimitInfo, window *CodexUsageWindow) []QuotaRow {
	// 每个窗口单独展开为一条 quota row，前端再按窗口秒数决定主/次展示位置。
	if window == nil {
		return rows
	}
	label = codexWindowLabel(key, label, window.LimitWindowSeconds)
	row := QuotaRow{
		Key:               key,
		Label:             label,
		Scope:             scope,
		Metric:            metric,
		UsedPercent:       floatPtr(window.UsedPercent),
		Allowed:           info.Allowed,
		LimitReached:      info.LimitReached,
		WindowUsageTokens: window.WindowUsageTokens,
		WindowUsageCost:   window.WindowUsageCost,
	}
	if window.LimitWindowSeconds != 0 {
		row.Window = &QuotaWindow{Seconds: intPtr(window.LimitWindowSeconds)}
	}
	if window.ResetAfterSeconds != 0 {
		row.ResetAfterSeconds = intPtr(window.ResetAfterSeconds)
	}
	if window.ResetAt != 0 {
		row.ResetAt = timeutil.FormatStorageTime(time.Unix(window.ResetAt, 0))
	}
	return append(rows, row)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func floatPtr(value float64) *float64 {
	return &value
}

func intPtr(value int64) *int64 {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

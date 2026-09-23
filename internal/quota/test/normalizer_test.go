package test

import (
	"testing"
	"time"

	"cpa-usage-keeper/internal/quota"
	"cpa-usage-keeper/internal/timeutil"
)

func TestNormalizeClaudeQuotaRows(t *testing.T) {
	utilization := 25.0
	rows := quota.NormalizeQuotaRows(quota.ProviderOutput{Provider: "claude", Result: quota.ClaudeResult{
		Usage: &quota.ClaudeUsagePayload{
			FiveHour:       &quota.ClaudeUsageWindow{Utilization: 36, ResetsAt: "2026-05-09T12:00:00Z"},
			SevenDay:       &quota.ClaudeUsageWindow{Utilization: 72, ResetsAt: "2026-05-10T12:00:00Z"},
			SevenDaySonnet: &quota.ClaudeUsageWindow{Utilization: 18, ResetsAt: "2026-05-11T12:00:00Z"},
			ExtraUsage:     &quota.ClaudeExtraUsage{IsEnabled: true, MonthlyLimit: 1000, UsedCredits: 250, Utilization: &utilization},
		},
		Profile: &quota.ClaudeProfileResponse{Account: &quota.ClaudeProfileAccount{Email: "user@example.com"}},
	}})

	if len(rows) != 4 {
		t.Fatalf("expected 4 quota rows, got %#v", rows)
	}
	fiveHour := findQuotaRow(t, rows, "five_hour")
	assertQuotaText(t, fiveHour, "5h", "window", "")
	assertFloatField(t, fiveHour.UsedPercent, 36, "five_hour usedPercent")
	if fiveHour.ResetAt != "2026-05-09T12:00:00Z" {
		t.Fatalf("unexpected five_hour resetAt: %#v", fiveHour)
	}
	assertIntField(t, fiveHour.Window.Seconds, 18000, "five_hour window seconds")
	weekly := findQuotaRow(t, rows, "seven_day")
	assertQuotaText(t, weekly, "Weekly", "window", "")
	assertFloatField(t, weekly.UsedPercent, 72, "seven_day usedPercent")
	assertIntField(t, weekly.Window.Seconds, 604800, "seven_day window seconds")
	sonnet := findQuotaRow(t, rows, "seven_day_sonnet")
	assertQuotaText(t, sonnet, "7d Sonnet", "model", "")
	assertIntField(t, sonnet.Window.Seconds, 604800, "seven_day_sonnet window seconds")
	extra := findQuotaRow(t, rows, "extra_usage")
	assertQuotaText(t, extra, "Extra Usage", "extra_usage", "")
	assertFloatField(t, extra.Used, 250, "extra_usage used")
	assertFloatField(t, extra.Limit, 1000, "extra_usage limit")
	assertFloatField(t, extra.UsedPercent, 25, "extra_usage usedPercent")
	assertBoolField(t, extra.Allowed, true, "extra_usage allowed")
}

func TestNormalizeCodexQuotaRows(t *testing.T) {
	previousLocal := time.Local
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	t.Cleanup(func() { time.Local = previousLocal })
	time.Local = location

	allowed := true
	limitReached := false
	resetAt := int64(1760000000)
	rows := quota.NormalizeQuotaRows(quota.ProviderOutput{Provider: "codex", Result: quota.CodexResult{Usage: &quota.CodexUsagePayload{
		PlanType: "plus",
		RateLimit: &quota.CodexRateLimitInfo{
			Allowed:      &allowed,
			LimitReached: &limitReached,
			PrimaryWindow: &quota.CodexUsageWindow{
				UsedPercent:        25,
				LimitWindowSeconds: 18000,
				ResetAfterSeconds:  1200,
				ResetAt:            resetAt,
			},
			SecondaryWindow: &quota.CodexUsageWindow{UsedPercent: 65, LimitWindowSeconds: 604800, ResetAfterSeconds: 7200},
		},
		CodeReviewRateLimit: &quota.CodexRateLimitInfo{
			Allowed:       &allowed,
			PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 40, LimitWindowSeconds: 18000, ResetAfterSeconds: 600},
		},
		AdditionalRateLimits: []quota.CodexAdditionalRateLimit{{
			LimitName:      "codex-spark",
			MeteredFeature: "spark",
			RateLimit: &quota.CodexRateLimitInfo{
				Allowed:       &allowed,
				PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 12, LimitWindowSeconds: 18000, ResetAfterSeconds: 900},
			},
		}},
	}}})

	if len(rows) != 4 {
		t.Fatalf("expected 4 quota rows, got %#v", rows)
	}
	primary := findQuotaRow(t, rows, "rate_limit.primary_window")
	assertQuotaText(t, primary, "5h", "window", "")
	assertFloatField(t, primary.UsedPercent, 25, "primary usedPercent")
	assertIntField(t, primary.Window.Seconds, 18000, "primary window seconds")
	assertIntField(t, primary.ResetAfterSeconds, 1200, "primary resetAfterSeconds")
	if primary.ResetAt != timeutil.FormatStorageTime(time.Unix(resetAt, 0)) {
		t.Fatalf("unexpected primary resetAt: %#v", primary)
	}
	assertBoolField(t, primary.Allowed, true, "primary allowed")
	assertBoolField(t, primary.LimitReached, false, "primary limitReached")

	secondary := findQuotaRow(t, rows, "rate_limit.secondary_window")
	assertQuotaText(t, secondary, "Weekly", "window", "")
	assertFloatField(t, secondary.UsedPercent, 65, "secondary usedPercent")
	codeReview := findQuotaRow(t, rows, "code_review_rate_limit.primary_window")
	assertQuotaText(t, codeReview, "Code Review 5h", "code_review", "")
	additional := findQuotaRow(t, rows, "additional_rate_limits.codex-spark.primary_window")
	assertQuotaText(t, additional, "codex-spark 5h", "additional", "spark")
	assertFloatField(t, additional.UsedPercent, 12, "additional usedPercent")
}

func TestNormalizeCodexPrimaryWindowUsesWindowSecondsForWeeklyLabel(t *testing.T) {
	rows := quota.NormalizeQuotaRows(quota.ProviderOutput{Provider: "codex", Result: quota.CodexResult{Usage: &quota.CodexUsagePayload{
		RateLimit: &quota.CodexRateLimitInfo{
			PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 10, LimitWindowSeconds: 604800},
		},
	}}})

	primary := findQuotaRow(t, rows, "rate_limit.primary_window")
	assertQuotaText(t, primary, "Weekly", "window", "")
	assertIntField(t, primary.Window.Seconds, 604800, "primary weekly window seconds")
	if primary.ResetAfterSeconds != nil {
		t.Fatalf("expected missing Codex reset_after_seconds to stay nil, got %#v", primary.ResetAfterSeconds)
	}
}

func TestNormalizeCodexPrimaryWindowUsesWindowSecondsForMonthlyLabel(t *testing.T) {
	rows := quota.NormalizeQuotaRows(quota.ProviderOutput{Provider: "codex", Result: quota.CodexResult{Usage: &quota.CodexUsagePayload{
		RateLimit: &quota.CodexRateLimitInfo{
			PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 10, LimitWindowSeconds: 2628000},
		},
		CodeReviewRateLimit: &quota.CodexRateLimitInfo{
			PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 25, LimitWindowSeconds: 2592000},
		},
		AdditionalRateLimits: []quota.CodexAdditionalRateLimit{{
			LimitName: "codex-spark",
			RateLimit: &quota.CodexRateLimitInfo{
				PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 40, LimitWindowSeconds: 2628000},
			},
		}},
	}}})

	primary := findQuotaRow(t, rows, "rate_limit.primary_window")
	assertQuotaText(t, primary, "Monthly", "window", "")
	assertIntField(t, primary.Window.Seconds, 2628000, "primary average monthly window seconds")
	codeReview := findQuotaRow(t, rows, "code_review_rate_limit.primary_window")
	assertQuotaText(t, codeReview, "Code Review Monthly", "code_review", "")
	additional := findQuotaRow(t, rows, "additional_rate_limits.codex-spark.primary_window")
	assertQuotaText(t, additional, "codex-spark Monthly", "additional", "codex-spark")
}

func TestNormalizeCodexUnknownWindowKeepsRoleLabelWithoutGenericWindow(t *testing.T) {
	rows := quota.NormalizeQuotaRows(quota.ProviderOutput{Provider: "codex", Result: quota.CodexResult{Usage: &quota.CodexUsagePayload{
		RateLimit: &quota.CodexRateLimitInfo{
			PrimaryWindow:   &quota.CodexUsageWindow{UsedPercent: 10, LimitWindowSeconds: 3600},
			SecondaryWindow: &quota.CodexUsageWindow{UsedPercent: 20, LimitWindowSeconds: 7200},
		},
		CodeReviewRateLimit: &quota.CodexRateLimitInfo{
			PrimaryWindow: &quota.CodexUsageWindow{UsedPercent: 30, LimitWindowSeconds: 3600},
		},
		AdditionalRateLimits: []quota.CodexAdditionalRateLimit{{
			LimitName: "GPT-5.3-Codex-Spark",
			RateLimit: &quota.CodexRateLimitInfo{
				SecondaryWindow: &quota.CodexUsageWindow{UsedPercent: 40, LimitWindowSeconds: 7200},
			},
		}},
	}}})

	primary := findQuotaRow(t, rows, "rate_limit.primary_window")
	assertQuotaText(t, primary, "Primary", "window", "")
	assertIntField(t, primary.Window.Seconds, 3600, "primary unknown window seconds")
	secondary := findQuotaRow(t, rows, "rate_limit.secondary_window")
	assertQuotaText(t, secondary, "Secondary", "window", "")
	assertIntField(t, secondary.Window.Seconds, 7200, "secondary unknown window seconds")
	codeReview := findQuotaRow(t, rows, "code_review_rate_limit.primary_window")
	assertQuotaText(t, codeReview, "Code Review Primary", "code_review", "")
	additional := findQuotaRow(t, rows, "additional_rate_limits.GPT-5.3-Codex-Spark.secondary_window")
	assertQuotaText(t, additional, "GPT-5.3-Codex-Spark Secondary", "additional", "GPT-5.3-Codex-Spark")
}

func findQuotaRow(t *testing.T, rows []quota.QuotaRow, key string) quota.QuotaRow {
	t.Helper()
	for _, row := range rows {
		if row.Key == key {
			return row
		}
	}
	t.Fatalf("missing quota row %q in %#v", key, rows)
	return quota.QuotaRow{}
}

func assertQuotaText(t *testing.T, row quota.QuotaRow, label string, scope string, metric string) {
	t.Helper()
	if row.Label != label || row.Scope != scope || row.Metric != metric {
		t.Fatalf("unexpected quota row text: got %#v want label=%q scope=%q metric=%q", row, label, scope, metric)
	}
}

func assertFloatField(t *testing.T, value *float64, expected float64, label string) {
	t.Helper()
	if value == nil || *value != expected {
		t.Fatalf("unexpected %s: got %#v want %v", label, value, expected)
	}
}

func assertIntField(t *testing.T, value *int64, expected int64, label string) {
	t.Helper()
	if value == nil || *value != expected {
		t.Fatalf("unexpected %s: got %#v want %v", label, value, expected)
	}
}

func assertBoolField(t *testing.T, value *bool, expected bool, label string) {
	t.Helper()
	if value == nil || *value != expected {
		t.Fatalf("unexpected %s: got %#v want %v", label, value, expected)
	}
}

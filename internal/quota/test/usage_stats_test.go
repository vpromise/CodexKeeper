package test

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/pricing"
	. "cpa-usage-keeper/internal/quota"
	"cpa-usage-keeper/internal/repository"
	"cpa-usage-keeper/internal/repository/dto"
	"cpa-usage-keeper/internal/timeutil"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func TestQuotaRowUsageWindowParsesStorageTimeResetAt(t *testing.T) {
	windowSeconds := int64(5 * 60 * 60)
	resetAt := "2026-05-26 03:00:00"
	now := time.Date(2026, 5, 26, 2, 0, 0, 0, time.UTC)

	windowStart, windowEnd, ok := quotaRowUsageWindow(QuotaRow{
		ResetAt: resetAt,
		Window:  &QuotaWindow{Seconds: &windowSeconds},
	}, now)

	if !ok {
		t.Fatal("expected storage-time resetAt to produce usage window")
	}
	if windowStart.IsZero() || windowEnd.IsZero() {
		t.Fatalf("expected non-zero window, got [%s, %s)", windowStart, windowEnd)
	}
}

func TestQuotaRowUsageWindowUsesResetAfterSecondsWhenResetAtMissing(t *testing.T) {
	windowSeconds := int64(5 * 60 * 60)
	resetAfterSeconds := int64(60 * 60)
	now := time.Date(2026, 5, 26, 2, 0, 0, 0, time.UTC)

	windowStart, windowEnd, ok := quotaRowUsageWindow(QuotaRow{
		Window:            &QuotaWindow{Seconds: &windowSeconds},
		ResetAfterSeconds: &resetAfterSeconds,
	}, now)

	if !ok {
		t.Fatal("expected reset-after row to produce usage window")
	}
	wantEnd := now
	wantStart := now.Add(time.Duration(resetAfterSeconds-windowSeconds) * time.Second)
	if !windowEnd.Equal(wantEnd) || !windowStart.Equal(wantStart) {
		t.Fatalf("expected window [%s, %s), got [%s, %s)", wantStart, wantEnd, windowStart, windowEnd)
	}
}

func TestQuotaRowUsageWindowRejectsNegativeResetAfterSeconds(t *testing.T) {
	windowSeconds := int64(5 * 60 * 60)
	resetAfterSeconds := int64(-60)
	now := time.Date(2026, 5, 26, 2, 0, 0, 0, time.UTC)

	_, _, ok := quotaRowUsageWindow(QuotaRow{
		Window:            &QuotaWindow{Seconds: &windowSeconds},
		ResetAfterSeconds: &resetAfterSeconds,
	}, now)

	if ok {
		t.Fatal("expected negative reset_after_seconds to be ignored")
	}
}

func TestAttachWindowUsageStatsOnlyBackfillsMissingKnownWindowScopeRows(t *testing.T) {
	db := openQuotaTestDB(t)
	service := NewServiceWithRegistry(db, NewProviderRegistry(nil), emptyPricingCatalogForTest())
	defer service.StopRefreshTasks()
	windowSeconds := int64(5 * 60 * 60)
	weeklySeconds := int64(7 * 24 * 60 * 60)
	monthlySeconds := int64(30 * 24 * 60 * 60)
	averageMonthlySeconds := quotaWindowAverageMonthSeconds
	unknownSeconds := int64(60 * 60)
	resetAt := time.Date(2026, 6, 2, 5, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 2, 3, 0, 0, 0, time.UTC)

	if err := db.Create(&entities.UsageEvent{
		AuthIndex:   "auth-pro",
		Model:       "gpt-codex",
		Timestamp:   now.Add(-time.Hour),
		TotalTokens: 222,
	}).Error; err != nil {
		t.Fatalf("seed usage event: %v", err)
	}

	response := attachWindowUsageStats(service, context.Background(), "auth-pro", CheckResponse{ID: "auth-pro", Quota: []QuotaRow{
		{
			Key:               "rate_limit.primary_window",
			Label:             "5h",
			Scope:             "window",
			Window:            &QuotaWindow{Seconds: &windowSeconds},
			ResetAt:           timeutil.FormatStorageTime(resetAt),
			WindowUsageTokens: intPtr(11),
			WindowUsageCost:   floatPtr(0.42),
		},
		{
			Key:     "additional_rate_limits.GPT-5.3-Codex-Spark.primary_window",
			Label:   "GPT-5.3-Codex-Spark 5h",
			Scope:   "additional",
			Metric:  "codex_bengalfox",
			Window:  &QuotaWindow{Seconds: &windowSeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
		{
			Key:     "rate_limit.secondary_window",
			Label:   "Weekly",
			Scope:   "window",
			Window:  &QuotaWindow{Seconds: &weeklySeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
		{
			Key:     "rate_limit.monthly_window",
			Label:   "Monthly",
			Scope:   "window",
			Window:  &QuotaWindow{Seconds: &monthlySeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
		{
			Key:     "rate_limit.average_monthly_window",
			Label:   "Monthly",
			Scope:   "window",
			Window:  &QuotaWindow{Seconds: &averageMonthlySeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
		{
			Key:     "code_review_rate_limit.primary_window",
			Label:   "Code Review 5h",
			Scope:   "code_review",
			Window:  &QuotaWindow{Seconds: &windowSeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
		{
			Key:     "rate_limit.unknown_window",
			Label:   "Unknown",
			Scope:   "window",
			Window:  &QuotaWindow{Seconds: &unknownSeconds},
			ResetAt: timeutil.FormatStorageTime(resetAt),
		},
	}}, now)

	providerWindow := findQuotaRow(t, response.Quota, "rate_limit.primary_window")
	if providerWindow.WindowUsageTokens == nil || *providerWindow.WindowUsageTokens != 11 {
		t.Fatalf("expected provider window_usage_tokens to be preserved, got %#v", providerWindow.WindowUsageTokens)
	}
	if providerWindow.WindowUsageCost == nil || !(math.Abs(*providerWindow.WindowUsageCost-0.42) <= 0.000000001) {
		t.Fatalf("expected provider window_usage_cost to be preserved, got %#v", providerWindow.WindowUsageCost)
	}

	for _, key := range []string{"rate_limit.secondary_window", "rate_limit.monthly_window", "rate_limit.average_monthly_window"} {
		row := findQuotaRow(t, response.Quota, key)
		assertWindowUsage(t, row, 222, 0)
	}
	for _, key := range []string{"additional_rate_limits.GPT-5.3-Codex-Spark.primary_window", "code_review_rate_limit.primary_window", "rate_limit.unknown_window"} {
		row := findQuotaRow(t, response.Quota, key)
		if row.WindowUsageTokens != nil || row.WindowUsageCost != nil {
			t.Fatalf("%s unexpectedly backfilled: %+v", key, row)
		}
	}
}

func TestAttachWindowUsageStatsBackfillsBothFieldsWhenProviderWindowUsageIncomplete(t *testing.T) {
	db := openQuotaTestDB(t)
	service := NewServiceWithRegistry(db, NewProviderRegistry(nil), emptyPricingCatalogForTest())
	defer service.StopRefreshTasks()
	windowSeconds := int64(5 * 60 * 60)
	resetAt := time.Date(2026, 6, 2, 5, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 2, 3, 0, 0, 0, time.UTC)

	if err := db.Create(&entities.UsageEvent{
		AuthIndex:   "auth-partial",
		Model:       "gpt-codex",
		Timestamp:   now.Add(-time.Hour),
		TotalTokens: 333,
	}).Error; err != nil {
		t.Fatalf("seed usage event: %v", err)
	}

	response := attachWindowUsageStats(service, context.Background(), "auth-partial", CheckResponse{ID: "auth-partial", Quota: []QuotaRow{
		{
			Key:               "rate_limit.primary_window",
			Label:             "5h",
			Scope:             "window",
			Window:            &QuotaWindow{Seconds: &windowSeconds},
			ResetAt:           timeutil.FormatStorageTime(resetAt),
			WindowUsageTokens: intPtr(11),
		},
		{
			Key:             "rate_limit.secondary_window",
			Label:           "5h",
			Scope:           "window",
			Window:          &QuotaWindow{Seconds: &windowSeconds},
			ResetAt:         timeutil.FormatStorageTime(resetAt),
			WindowUsageCost: floatPtr(0.42),
		},
	}}, now)

	for _, key := range []string{"rate_limit.primary_window", "rate_limit.secondary_window"} {
		assertWindowUsage(t, findQuotaRow(t, response.Quota, key), 333, 0)
	}
}

func TestAttachWindowUsageStatsDropsIncompleteProviderWindowUsageWhenFallbackUnavailable(t *testing.T) {
	db := openQuotaTestDB(t)
	service := NewServiceWithRegistry(db, NewProviderRegistry(nil), emptyPricingCatalogForTest())
	defer service.StopRefreshTasks()
	windowSeconds := int64(5 * 60 * 60)

	response := attachWindowUsageStats(service, context.Background(), "auth-partial", CheckResponse{ID: "auth-partial", Quota: []QuotaRow{{
		Key:               "rate_limit.primary_window",
		Label:             "5h",
		Scope:             "window",
		Window:            &QuotaWindow{Seconds: &windowSeconds},
		WindowUsageTokens: intPtr(11),
	}}}, time.Date(2026, 6, 2, 3, 0, 0, 0, time.UTC))

	row := findQuotaRow(t, response.Quota, "rate_limit.primary_window")
	if row.WindowUsageTokens != nil || row.WindowUsageCost != nil {
		t.Fatalf("expected incomplete provider usage to be dropped when local fallback is unavailable, got tokens=%#v cost=%#v", row.WindowUsageTokens, row.WindowUsageCost)
	}
}

func TestAttachWindowUsageStatsDoesNotBackfillPartialAdditionalOrCodeReviewRows(t *testing.T) {
	db := openQuotaTestDB(t)
	service := NewServiceWithRegistry(db, NewProviderRegistry(nil), emptyPricingCatalogForTest())
	defer service.StopRefreshTasks()
	windowSeconds := int64(5 * 60 * 60)
	resetAt := time.Date(2026, 6, 2, 5, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 2, 3, 0, 0, 0, time.UTC)

	if err := db.Create(&entities.UsageEvent{
		AuthIndex:   "auth-special",
		Model:       "gpt-codex",
		Timestamp:   now.Add(-time.Hour),
		TotalTokens: 444,
	}).Error; err != nil {
		t.Fatalf("seed usage event: %v", err)
	}

	response := attachWindowUsageStats(service, context.Background(), "auth-special", CheckResponse{ID: "auth-special", Quota: []QuotaRow{
		{
			Key:               "additional_rate_limits.GPT-5.3-Codex-Spark.primary_window",
			Label:             "GPT-5.3-Codex-Spark 5h",
			Scope:             "additional",
			Metric:            "codex_bengalfox",
			Window:            &QuotaWindow{Seconds: &windowSeconds},
			ResetAt:           timeutil.FormatStorageTime(resetAt),
			WindowUsageTokens: intPtr(11),
		},
		{
			Key:             "code_review_rate_limit.primary_window",
			Label:           "Code Review 5h",
			Scope:           "code_review",
			Window:          &QuotaWindow{Seconds: &windowSeconds},
			ResetAt:         timeutil.FormatStorageTime(resetAt),
			WindowUsageCost: floatPtr(0.42),
		},
	}}, now)

	additional := findQuotaRow(t, response.Quota, "additional_rate_limits.GPT-5.3-Codex-Spark.primary_window")
	if additional.WindowUsageTokens == nil || *additional.WindowUsageTokens != 11 || additional.WindowUsageCost != nil {
		t.Fatalf("expected partial additional provider usage to be preserved without fallback, got tokens=%#v cost=%#v", additional.WindowUsageTokens, additional.WindowUsageCost)
	}
	codeReview := findQuotaRow(t, response.Quota, "code_review_rate_limit.primary_window")
	if codeReview.WindowUsageTokens != nil || codeReview.WindowUsageCost == nil || !(math.Abs(*codeReview.WindowUsageCost-0.42) <= 0.000000001) {
		t.Fatalf("expected partial code review provider usage to be preserved without fallback, got tokens=%#v cost=%#v", codeReview.WindowUsageTokens, codeReview.WindowUsageCost)
	}
}

func TestAttachWindowUsageStatsPreservesProviderZeroWindowUsage(t *testing.T) {
	db := openQuotaTestDB(t)
	service := NewServiceWithRegistry(db, NewProviderRegistry(nil), emptyPricingCatalogForTest())
	defer service.StopRefreshTasks()
	windowSeconds := int64(5 * 60 * 60)
	resetAt := time.Date(2026, 6, 2, 5, 0, 0, 0, time.UTC)
	now := time.Date(2026, 6, 2, 3, 0, 0, 0, time.UTC)
	if _, err := repository.UpsertModelPriceSetting(db, dto.ModelPriceSettingInput{Model: "priced", PromptPricePer1M: 10, CompletionPricePer1M: 20, CacheReadPricePer1M: 1}); err != nil {
		t.Fatalf("UpsertModelPriceSetting returned error: %v", err)
	}
	if err := db.Create(&entities.UsageEvent{
		AuthIndex:    "auth-zero",
		Model:        "priced",
		Timestamp:    now.Add(-time.Hour),
		InputTokens:  1_000_000,
		OutputTokens: 500_000,
		TotalTokens:  1_500_000,
	}).Error; err != nil {
		t.Fatalf("seed usage event: %v", err)
	}

	response := attachWindowUsageStats(service, context.Background(), "auth-zero", CheckResponse{ID: "auth-zero", Quota: []QuotaRow{{
		Key:               "rate_limit.primary_window",
		Label:             "5h",
		Scope:             "window",
		Window:            &QuotaWindow{Seconds: &windowSeconds},
		ResetAt:           timeutil.FormatStorageTime(resetAt),
		WindowUsageTokens: intPtr(0),
		WindowUsageCost:   floatPtr(0),
	}}}, now)

	row := findQuotaRow(t, response.Quota, "rate_limit.primary_window")
	if row.WindowUsageTokens == nil || *row.WindowUsageTokens != 0 {
		t.Fatalf("expected provider zero tokens to be preserved, got %#v", row.WindowUsageTokens)
	}
	if row.WindowUsageCost == nil || *row.WindowUsageCost != 0 {
		t.Fatalf("expected provider zero cost to be preserved, got %#v", row.WindowUsageCost)
	}
}

type failOnceGroupedUsageStatsProvider struct {
	groupedCalls int
}

func (provider *failOnceGroupedUsageStatsProvider) SumByAuthIndex(context.Context, string, time.Time, *time.Time) (repository.UsageWindowStats, error) {
	return repository.UsageWindowStats{}, nil
}

func (provider *failOnceGroupedUsageStatsProvider) SumGroupsByAuthIndex(context.Context, string, time.Time, *time.Time, repository.UsageWindowStatsGrouper) (repository.UsageWindowGroupedStats, error) {
	provider.groupedCalls++
	if provider.groupedCalls == 1 {
		return repository.UsageWindowGroupedStats{}, errors.New("transient grouped query failure")
	}
	return repository.UsageWindowGroupedStats{
		Complete: true,
		Groups: map[string]repository.UsageWindowStats{
			"antigravity-claude-and-gpt-models": {Tokens: 10, Cost: 1, CostAvailable: true},
		},
	}, nil
}

type quotaUsageQueryCounter struct {
	mu    sync.Mutex
	count int
}

func (counter *quotaUsageQueryCounter) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return counter
}
func (counter *quotaUsageQueryCounter) Info(context.Context, string, ...interface{})  {}
func (counter *quotaUsageQueryCounter) Warn(context.Context, string, ...interface{})  {}
func (counter *quotaUsageQueryCounter) Error(context.Context, string, ...interface{}) {}
func (counter *quotaUsageQueryCounter) Trace(_ context.Context, _ time.Time, sqlAndRows func() (string, int64), _ error) {
	sql, _ := sqlAndRows()
	if !strings.Contains(sql, "FROM `usage_events`") || !strings.Contains(sql, "GROUP BY") {
		return
	}
	counter.mu.Lock()
	counter.count++
	counter.mu.Unlock()
}

func (counter *quotaUsageQueryCounter) UsageWindowQueries() int {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	return counter.count
}

func quotaUsagePricingCatalog(t *testing.T, db *gorm.DB) *pricing.Catalog {
	t.Helper()
	snapshot, err := repository.LoadPricingSnapshot(context.Background(), db)
	if err != nil {
		t.Fatalf("LoadPricingSnapshot: %v", err)
	}
	return pricing.NewCatalog(snapshot)
}

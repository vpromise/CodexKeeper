package test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"cpa-usage-keeper/internal/cpa/dto/authfiles"
	"cpa-usage-keeper/internal/cpa/dto/cpaapikeys"
	"cpa-usage-keeper/internal/cpa/dto/providerconfig"
	"cpa-usage-keeper/internal/cpa/dto/response"
	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"

	"gorm.io/gorm"
)

func TestSyncMetadataPreservesAuthFilesAndManagementAPIKeySemantics(t *testing.T) {
	db := openMetadataTestDatabase(t, "auth-and-management.db")
	now := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	activeStart := now.Add(-24 * time.Hour)
	activeUntil := now.Add(6 * 24 * time.Hour)
	accountID := "acct-codex"
	planType := "team"
	fetcher := newMetadataTestFetcher()
	// Auth Files 覆盖 email/label/name/auth-index fallback 和专属扩展字段。
	fetcher.setAuthFiles([]authfiles.AuthFile{
		{AuthIndex: "auth-email", Name: "claude.json", Path: "/data/auths/claude.json", Email: "user@example.com", Label: "ignored-label", Type: " CLAUDE ", Provider: "Claude", Prefix: "oauth-prefix", Priority: metadataIntPointer(6), Disabled: metadataBoolPointer(false), Note: metadataStringPointer("auth note")},
		{AuthIndex: "auth-label", Name: "claude-label.json", Label: "Claude Label", Type: "claude", Provider: "Claude"},
		{AuthIndex: "auth-name", Name: "Codex Name", Type: "codex", Provider: "Codex", IDToken: &authfiles.AuthFileIDToken{AccountID: &accountID, ActiveStart: &activeStart, ActiveUntil: &activeUntil, PlanType: &planType}},
		{AuthIndex: "auth-index-fallback", Type: "codex", Provider: "Codex"},
	})
	fetcher.managementAPIKeysResult = &response.ManagementAPIKeysResult{StatusCode: 200, Payload: cpaapikeys.ManagementAPIKeysResponse{APIKeys: []string{"sk-alpha123456", "sk-beta654321"}}}
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	if err := syncer.SyncMetadata(context.Background()); err != nil {
		t.Fatalf("SyncMetadata returned error: %v", err)
	}
	if fetcher.callCount("auth-files") != 1 {
		t.Fatalf("auth-files calls = %d", fetcher.callCount("auth-files"))
	}
	if fetcher.callCount("management-api-keys") != 1 {
		t.Fatalf("management-api-keys calls = %d", fetcher.callCount("management-api-keys"))
	}
	identities := loadMetadataIdentityMap(t, db)
	emailRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAuthFile, "auth-email")]
	if emailRow.Name != "user@example.com" || emailRow.AuthTypeName != "oauth" || emailRow.Type != "claude" || emailRow.Provider != "Claude" || emailRow.Prefix != "oauth-prefix" || emailRow.IsDeleted {
		t.Fatalf("email auth identity = %+v", emailRow)
	}
	if emailRow.FileName == nil || *emailRow.FileName != "claude.json" || emailRow.FilePath == nil || *emailRow.FilePath != "/data/auths/claude.json" {
		t.Fatalf("email auth file metadata = %+v", emailRow)
	}
	if emailRow.Priority == nil || *emailRow.Priority != 6 || emailRow.Disabled == nil || *emailRow.Disabled || emailRow.Note == nil || *emailRow.Note != "auth note" {
		t.Fatalf("email auth optional metadata = %+v", emailRow)
	}
	labelRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAuthFile, "auth-label")]
	if labelRow.Name != "Claude Label" {
		t.Fatalf("label auth identity = %+v", labelRow)
	}
	codexRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAuthFile, "auth-name")]
	if codexRow.Name != "Codex Name" || codexRow.AccountID == nil || *codexRow.AccountID != accountID || codexRow.PlanType == nil || *codexRow.PlanType != planType || codexRow.ActiveStart == nil || !codexRow.ActiveStart.Equal(activeStart) || codexRow.ActiveUntil == nil || !codexRow.ActiveUntil.Equal(activeUntil) {
		t.Fatalf("codex auth identity = %+v", codexRow)
	}
	fallbackRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAuthFile, "auth-index-fallback")]
	if fallbackRow.Name != "auth-index-fallback" {
		t.Fatalf("fallback auth identity = %+v", fallbackRow)
	}
	// 所有新建身份必须采用本轮 now，不能使用 GORM 系统时钟或零时间。
	if !emailRow.CreatedAt.Equal(now) || !emailRow.UpdatedAt.Equal(now) {
		t.Fatalf("auth identity times = created:%s updated:%s", emailRow.CreatedAt, emailRow.UpdatedAt)
	}
	apiKeys, err := repository.ListActiveCPAAPIKeys(db)
	if err != nil {
		t.Fatalf("ListActiveCPAAPIKeys returned error: %v", err)
	}
	if len(apiKeys) != 2 || apiKeys[0].DisplayKey != "sk-*********123456" || apiKeys[0].APIKey != "sk-alpha123456" {
		t.Fatalf("management API keys = %+v", apiKeys)
	}
}

func TestSyncMetadataFetchFailuresPreserveAuthFilesAndManagementKeys(t *testing.T) {
	db := openMetadataTestDatabase(t, "preserve-on-fetch-failure.db")
	firstNow := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	secondNow := firstNow.Add(time.Hour)
	fetcher := newMetadataTestFetcher()
	fetcher.setAuthFiles([]authfiles.AuthFile{{AuthIndex: "preserved-auth", Email: "preserved@example.com", Type: "claude", Provider: "Claude"}})
	fetcher.managementAPIKeysResult = &response.ManagementAPIKeysResult{StatusCode: 200, Payload: cpaapikeys.ManagementAPIKeysResponse{APIKeys: []string{"sk-preserved123456"}}}
	currentNow := firstNow
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return currentNow })
	if err := syncer.SyncMetadata(context.Background()); err != nil {
		t.Fatalf("initial SyncMetadata returned error: %v", err)
	}
	currentNow = secondNow
	fetcher.authFilesResult = nil
	fetcher.authFilesErr = errors.New("auth unavailable")
	fetcher.managementAPIKeysResult = nil
	fetcher.managementAPIKeysErr = errors.New("management unavailable")
	err := syncer.SyncMetadata(context.Background())
	if err == nil || err.Error() != "fetch auth files: auth unavailable; fetch management api keys: management unavailable" {
		t.Fatalf("failure error = %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	preservedAuth := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAuthFile, "preserved-auth")]
	// fetch failure 不能把本地身份 stale 或刷新 updated_at。
	if preservedAuth.IsDeleted || preservedAuth.DeletedAt != nil || !preservedAuth.UpdatedAt.Equal(firstNow) {
		t.Fatalf("preserved auth identity = %+v", preservedAuth)
	}
	apiKeys, listErr := repository.ListActiveCPAAPIKeys(db)
	if listErr != nil {
		t.Fatalf("ListActiveCPAAPIKeys returned error: %v", listErr)
	}
	if len(apiKeys) != 1 || apiKeys[0].APIKey != "sk-preserved123456" {
		t.Fatalf("preserved management API keys = %+v", apiKeys)
	}
}

func TestSyncMetadataProviderWarningStillRunsHistoricalStatsCatchUp(t *testing.T) {
	db := openMetadataTestDatabase(t, "provider-warning-catch-up.db")
	eventTime := time.Date(2026, 7, 15, 8, 0, 0, 0, time.UTC)
	now := eventTime.Add(time.Hour)
	_, _, err := repository.InsertUsageEvents(db, []entities.UsageEvent{{EventKey: "provider-history", AuthType: "apikey", AuthIndex: "codex-history", Model: "gpt-5", Timestamp: eventTime, InputTokens: 11, OutputTokens: 13, TotalTokens: 24}})
	if err != nil {
		t.Fatalf("seed provider usage event: %v", err)
	}
	fetcher := newMetadataTestFetcher()
	fetcher.standardResults["codex"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "codex-secret", Name: "Codex History", AuthIndex: "codex-history"}}}
	fetcher.standardResults["claude"] = nil
	fetcher.standardErrors["claude"] = errors.New("claude unavailable")
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	err = syncer.SyncMetadata(context.Background())
	if err == nil || !strings.Contains(err.Error(), "fetch claude api keys: claude unavailable") {
		t.Fatalf("provider warning = %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	codexRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "codex-history")]
	if codexRow.Identity != "codex-history" || codexRow.LookupKey != "codex-secret" || codexRow.TotalRequests != 1 || codexRow.InputTokens != 11 || codexRow.OutputTokens != 13 || codexRow.TotalTokens != 24 || codexRow.LastAggregatedUsageEventID == 0 {
		t.Fatalf("provider catch-up identity = %+v", codexRow)
	}
	if codexRow.FirstUsedAt == nil || !codexRow.FirstUsedAt.Equal(eventTime) || codexRow.LastUsedAt == nil || !codexRow.LastUsedAt.Equal(eventTime) || codexRow.StatsUpdatedAt == nil || !codexRow.StatsUpdatedAt.Equal(now) {
		t.Fatalf("provider catch-up times = %+v", codexRow)
	}
}

func TestSyncMetadataProviderPersistenceErrorSuppressesFetchWarning(t *testing.T) {
	db := openMetadataTestDatabase(t, "provider-persistence-error.db")
	now := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	fetcher := newMetadataTestFetcher()
	fetcher.managementAPIKeysResult = &response.ManagementAPIKeysResult{StatusCode: 200, Payload: cpaapikeys.ManagementAPIKeysResponse{APIKeys: []string{"sk-independent123456"}}}
	fetcher.standardResults["codex"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "codex-secret", AuthIndex: "codex-persist-fail", Name: "Codex"}}}
	fetcher.standardResults["claude"] = nil
	fetcher.standardErrors["claude"] = errors.New("claude unavailable")
	if err := db.Callback().Create().Before("gorm:create").Register("test:block_provider_identity_create", func(tx *gorm.DB) {
		if tx.Statement.Table == "usage_identities" {
			tx.AddError(errors.New("provider persistence blocked"))
		}
	}); err != nil {
		t.Fatalf("register provider persistence callback: %v", err)
	}
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	err := syncer.SyncMetadata(context.Background())
	if err == nil || !strings.Contains(err.Error(), "sync provider usage identities: create usage identities: provider persistence blocked") {
		t.Fatalf("provider persistence error = %v", err)
	}
	if strings.Contains(err.Error(), "claude unavailable") {
		t.Fatalf("provider fetch warning was not suppressed: %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	if _, ok := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "codex-persist-fail")]; ok {
		t.Fatalf("provider identity survived rollback: %+v", identities)
	}
	apiKeys, listErr := repository.ListActiveCPAAPIKeys(db)
	if listErr != nil {
		t.Fatalf("ListActiveCPAAPIKeys returned error: %v", listErr)
	}
	if len(apiKeys) != 1 || apiKeys[0].APIKey != "sk-independent123456" {
		t.Fatalf("independent management keys = %+v", apiKeys)
	}
}

func TestSyncMetadataManagementPersistenceErrorKeepsProviderWriteAndStopsCatchUp(t *testing.T) {
	db := openMetadataTestDatabase(t, "management-persistence-error.db")
	eventTime := time.Date(2026, 7, 15, 8, 0, 0, 0, time.UTC)
	now := eventTime.Add(time.Hour)
	_, _, err := repository.InsertUsageEvents(db, []entities.UsageEvent{{EventKey: "blocked-catch-up", AuthType: "apikey", AuthIndex: "provider-survives", Model: "gpt-5", Timestamp: eventTime, TotalTokens: 10}})
	if err != nil {
		t.Fatalf("seed blocked catch-up event: %v", err)
	}
	if err := db.Migrator().DropTable(&entities.CPAAPIKey{}); err != nil {
		t.Fatalf("drop cpa_api_keys table: %v", err)
	}
	fetcher := newMetadataTestFetcher()
	fetcher.managementAPIKeysResult = &response.ManagementAPIKeysResult{StatusCode: 200, Payload: cpaapikeys.ManagementAPIKeysResponse{APIKeys: []string{"sk-fails123456"}}}
	fetcher.standardResults["codex"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "provider-secret", AuthIndex: "provider-survives", Name: "Provider Survives"}}}
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	err = syncer.SyncMetadata(context.Background())
	if err == nil || !strings.Contains(err.Error(), "sync management api keys") {
		t.Fatalf("management persistence error = %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	providerRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "provider-survives")]
	if providerRow.Identity != "provider-survives" || providerRow.LookupKey != "provider-secret" || providerRow.IsDeleted || providerRow.TotalRequests != 0 || providerRow.LastAggregatedUsageEventID != 0 || providerRow.StatsUpdatedAt != nil {
		t.Fatalf("provider after management persistence error = %+v", providerRow)
	}
}

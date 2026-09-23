package test

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"cpa-usage-keeper/internal/cpa/dto/providerconfig"
	"cpa-usage-keeper/internal/cpa/dto/response"
	"cpa-usage-keeper/internal/entities"

	"gorm.io/gorm"
)

func TestProviderMetadataSyncPreservesSourceFields(t *testing.T) {
	db := openMetadataTestDatabase(t, "existing-provider-fields.db")
	now := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	priority := 7
	disabled := true
	note := "provider note"
	fetcher := newMetadataTestFetcher()
	fetcher.standardResults["codex"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "secret-codex", AuthIndex: "auth-codex", Prefix: "prefix-codex", BaseURL: "https://codex.example/v1", Name: "Codex Team", Priority: &priority, Disabled: &disabled, Note: &note}}}
	fetcher.standardResults["claude"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "secret-claude", AuthIndex: "auth-claude", Prefix: "prefix-claude", BaseURL: "https://claude.example/v1", Name: "Claude Team"}}}
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	if err := syncer.SyncMetadata(context.Background()); err != nil {
		t.Fatalf("SyncMetadata returned error: %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	codexRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "auth-codex")]
	// Identity and lookup key come from distinct upstream fields.
	if codexRow.Name != "Codex Team" || codexRow.Provider != "Codex Team" || codexRow.Identity != "auth-codex" || codexRow.LookupKey != "secret-codex" || codexRow.Prefix != "prefix-codex" || codexRow.BaseURL != "https://codex.example/v1" || codexRow.Type != "codex" || codexRow.AuthTypeName != "apikey" || codexRow.IsDeleted {
		t.Fatalf("codex provider identity = %+v", codexRow)
	}
	if codexRow.Priority == nil || *codexRow.Priority != priority || codexRow.Disabled == nil || *codexRow.Disabled != disabled || codexRow.Note == nil || *codexRow.Note != note {
		t.Fatalf("codex provider optional metadata = %+v", codexRow)
	}
	claudeRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "auth-claude")]
	if claudeRow.Type != "claude" || claudeRow.Name != "Claude Team" || claudeRow.LookupKey != "secret-claude" {
		t.Fatalf("claude provider identity = %+v", claudeRow)
	}
	if _, ok := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "prefix-codex")]; ok {
		t.Fatalf("provider prefix created an identity: %+v", identities)
	}
	for _, source := range []string{"codex", "claude"} {
		if fetcher.callCount(source) != 1 {
			t.Fatalf("%s calls = %d", source, fetcher.callCount(source))
		}
	}
}

func TestProviderMetadataSyncKeepsFailedSourcesAndStalesOnlySuccessfulTypes(t *testing.T) {
	db := openMetadataTestDatabase(t, "provider-stale-boundaries.db")
	oldTime := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	now := oldTime.Add(24 * time.Hour)
	seed := []entities.UsageIdentity{
		{Name: "Old Claude", AuthType: entities.UsageIdentityAuthTypeAIProvider, AuthTypeName: "apikey", Identity: "old-claude", Type: "claude", Provider: "Claude", CreatedAt: oldTime, UpdatedAt: oldTime},
		{Name: "Old Codex", AuthType: entities.UsageIdentityAuthTypeAIProvider, AuthTypeName: "apikey", Identity: "old-codex", Type: "codex", Provider: "Codex", CreatedAt: oldTime, UpdatedAt: oldTime},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed provider identities: %v", err)
	}
	fetcher := newMetadataTestFetcher()
	fetcher.standardResults["claude"] = &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{}}
	fetcher.standardResults["codex"] = nil
	syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
	err := syncer.SyncMetadata(context.Background())
	if err == nil || !strings.Contains(err.Error(), "codex api keys response is nil") {
		t.Fatalf("provider boundary warning = %v", err)
	}
	identities := loadMetadataIdentityMap(t, db)
	codexRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "old-codex")]
	if codexRow.IsDeleted || codexRow.DeletedAt != nil || !codexRow.UpdatedAt.Equal(oldTime) {
		t.Fatalf("nil Codex identity = %+v", codexRow)
	}
	claudeRow := identities[metadataIdentityKey(entities.UsageIdentityAuthTypeAIProvider, "old-claude")]
	if !claudeRow.IsDeleted || claudeRow.DeletedAt == nil || !claudeRow.DeletedAt.Equal(now) || !claudeRow.UpdatedAt.Equal(now) {
		t.Fatalf("empty Claude identity = %+v", claudeRow)
	}
}

// Compare business fields independently of database-generated IDs.
type providerPersistenceProjection struct {
	Name         string
	ProviderType string
	Provider     string
	LookupKey    string
	Prefix       string
	BaseURL      string
	IsDeleted    bool
}

func TestProviderMetadataSyncCompletionOrderDoesNotChangeDatabase(t *testing.T) {
	registryOrder := []string{"codex", "claude"}
	orders := []struct {
		name            string
		completionOrder []string
	}{
		{name: "forward", completionOrder: []string{"codex", "claude"}},
		{name: "reverse", completionOrder: []string{"claude", "codex"}},
		{name: "mixed", completionOrder: []string{"codex", "claude"}},
	}
	want := make(map[string]providerPersistenceProjection, len(registryOrder))
	for _, source := range registryOrder {
		want["auth-"+source] = providerPersistenceProjection{
			Name: "name-" + source, ProviderType: source, Provider: "name-" + source,
			LookupKey: "secret-" + source, Prefix: "prefix-" + source, BaseURL: "https://" + source + ".example/v1",
		}
	}
	for _, order := range orders {
		t.Run(order.name, func(t *testing.T) {
			db := openMetadataTestDatabase(t, "completion-"+order.name+".db")
			fetcher := newMetadataTestFetcher()
			entered := make(chan string, len(registryOrder))
			done := make(chan string, len(registryOrder))
			gates := make(map[string]chan struct{}, len(registryOrder))
			gateOnce := make(map[string]*sync.Once, len(registryOrder))
			for _, source := range registryOrder {
				gates[source] = make(chan struct{})
				gateOnce[source] = &sync.Once{}
			}
			release := func(source string) {
				gateOnce[source].Do(func() {
					close(gates[source])
				})
			}
			// Release all endpoints on failure to avoid leaking goroutines.
			t.Cleanup(func() {
				for _, source := range registryOrder {
					release(source)
				}
			})
			for _, source := range registryOrder {
				fetcher.standardHooks[source] = func(ctx context.Context) (*response.ProviderKeyConfigResult, error) {
					entered <- source
					select {
					case <-gates[source]:
					case <-ctx.Done():
						return nil, ctx.Err()
					}
					done <- source
					return &response.ProviderKeyConfigResult{StatusCode: 200, Payload: []providerconfig.ProviderKeyConfig{{APIKey: "secret-" + source, AuthIndex: "auth-" + source, Name: "name-" + source, Prefix: "prefix-" + source, BaseURL: "https://" + source + ".example/v1"}}}, nil
				}
			}
			now := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
			syncer := newMetadataTestSyncer(db, fetcher, func() time.Time { return now })
			resultCh := make(chan error, 1)
			go func() {
				resultCh <- syncer.SyncMetadata(context.Background())
			}()
			waitForMetadataSourceSet(t, entered, registryOrder)
			for _, source := range order.completionOrder {
				release(source)
				waitForMetadataSourceSet(t, done, []string{source})
			}
			select {
			case err := <-resultCh:
				if err != nil {
					t.Fatalf("%s SyncMetadata returned error: %v", order.name, err)
				}
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for completion-order SyncMetadata")
			}
			projection := loadProviderPersistenceProjection(t, db)
			if !reflect.DeepEqual(projection, want) {
				t.Fatalf("%s projection = %#v, want %#v", order.name, projection, want)
			}
		})
	}
}

func waitForMetadataSourceSet(t *testing.T, events <-chan string, want []string) {
	t.Helper()
	seen := make(map[string]struct{}, len(want))
	// Bound the wait so concurrency regressions fail instead of hanging.
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for len(seen) < len(want) {
		select {
		case source := <-events:
			seen[source] = struct{}{}
		case <-timer.C:
			t.Fatalf("timed out waiting for metadata sources: got=%v want=%v", seen, want)
		}
	}
	for _, source := range want {
		if _, ok := seen[source]; !ok {
			t.Fatalf("metadata source %q did not run: %v", source, seen)
		}
	}
}

func loadProviderPersistenceProjection(t *testing.T, db *gorm.DB) map[string]providerPersistenceProjection {
	t.Helper()
	var rows []entities.UsageIdentity
	if err := db.Where("auth_type = ?", entities.UsageIdentityAuthTypeAIProvider).Find(&rows).Error; err != nil {
		t.Fatalf("load provider persistence projection: %v", err)
	}
	projection := make(map[string]providerPersistenceProjection, len(rows))
	for _, row := range rows {
		projection[row.Identity] = providerPersistenceProjection{Name: row.Name, ProviderType: row.Type, Provider: row.Provider, LookupKey: row.LookupKey, Prefix: row.Prefix, BaseURL: row.BaseURL, IsDeleted: row.IsDeleted}
	}
	return projection
}

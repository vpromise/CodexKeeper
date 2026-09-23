package providermetadata_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"cpa-usage-keeper/internal/cpa/dto/providerconfig"
	"cpa-usage-keeper/internal/cpa/dto/response"
	"cpa-usage-keeper/internal/service/providermetadata"
)

type providerFetcherStub struct {
	codexResult, claudeResult *response.ProviderKeyConfigResult
	codexErr, claudeErr       error
}

func (s *providerFetcherStub) FetchCodexAPIKeys(context.Context) (*response.ProviderKeyConfigResult, error) {
	return s.codexResult, s.codexErr
}
func (s *providerFetcherStub) FetchClaudeAPIKeys(context.Context) (*response.ProviderKeyConfigResult, error) {
	return s.claudeResult, s.claudeErr
}

func TestFetchNativeCredentialsPreservesNamesAndStableIndexes(t *testing.T) {
	priority, disabled, note := 7, false, "primary account"
	fetcher := &providerFetcherStub{
		codexResult: &response.ProviderKeyConfigResult{Payload: []providerconfig.ProviderKeyConfig{
			{APIKey: "codex-key", AuthIndex: "codex-index", Name: "Codex Team", Prefix: "pro01", BaseURL: "https://codex.example/v1", Priority: &priority, Disabled: &disabled, Note: &note},
			{APIKey: "duplicate", AuthIndex: "codex-index"},
			{APIKey: "missing-index"},
		}},
		claudeResult: &response.ProviderKeyConfigResult{Payload: []providerconfig.ProviderKeyConfig{{APIKey: "claude-key", AuthIndex: "claude-index", Name: "Claude Team"}}},
	}
	snapshot, err := providermetadata.Fetch(context.Background(), fetcher)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(snapshot.FetchedProviderTypes, []string{"codex", "claude"}) {
		t.Fatalf("unexpected provider scope: %v", snapshot.FetchedProviderTypes)
	}
	if len(snapshot.Credentials) != 2 {
		t.Fatalf("credentials: %+v", snapshot.Credentials)
	}
	got := snapshot.Credentials[0]
	if got.DisplayName != "Codex Team" || got.AuthIndex != "codex-index" || got.Prefix != "pro01" || got.LookupKey != "codex-key" || got.Note == nil || *got.Note != note || got.Priority == nil || *got.Priority != priority || got.Disabled == nil || *got.Disabled {
		t.Fatalf("native metadata changed: %+v", got)
	}
	if snapshot.Credentials[1].ProviderType != "claude" {
		t.Fatal("Claude source missing")
	}
}

func TestFetchEmptySuccessAndFailureHaveDifferentStaleScopes(t *testing.T) {
	fetcher := &providerFetcherStub{codexResult: &response.ProviderKeyConfigResult{}, claudeErr: errors.New("unavailable")}
	snapshot, err := providermetadata.Fetch(context.Background(), fetcher)
	if err == nil || !reflect.DeepEqual(snapshot.FetchedProviderTypes, []string{"codex"}) {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
	fetcher.claudeErr = nil
	snapshot, err = providermetadata.Fetch(context.Background(), fetcher)
	if err == nil || !reflect.DeepEqual(snapshot.FetchedProviderTypes, []string{"codex"}) {
		t.Fatalf("nil response must not stale Claude: %+v %v", snapshot, err)
	}
}

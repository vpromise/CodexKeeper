package cpa_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cpa-usage-keeper/internal/cpa"
	"cpa-usage-keeper/internal/cpa/dto/providerconfig"
)

type providerEndpointResult struct {
	statusCode int
	body       []byte
	keys       []providerconfig.ProviderKeyConfig
}

type providerEndpointCase struct {
	name        string
	path        string
	directBody  string
	wrappedBody string
}

func TestProviderAPIKeyClientsUseDedicatedEndpointsAndDecodePayloads(t *testing.T) {
	cases := []providerEndpointCase{
		{name: "codex", path: "/v0/management/codex-api-key", directBody: `[{"api-key":"codex-key","prefix":"codex-prefix","base-url":"https://codex.example/v1","name":"Codex","auth-index":"codex-auth"}]`, wrappedBody: `{"codex-api-key":[{"api-key":"codex-key","prefix":"codex-prefix","base-url":"https://codex.example/v1","name":"Codex","auth-index":"codex-auth"}]}`},
		{name: "claude", path: "/v0/management/claude-api-key", directBody: `[{"api-key":"claude-key","prefix":"claude-prefix","base-url":"https://claude.example/v1","name":"Claude","auth-index":"claude-auth"}]`, wrappedBody: `{"claude-api-key":[{"api-key":"claude-key","prefix":"claude-prefix","base-url":"https://claude.example/v1","name":"Claude","auth-index":"claude-auth"}]}`},
	}
	for _, tc := range cases {
		for _, variant := range []struct {
			name string
			body string
		}{
			{name: "direct", body: tc.directBody},
			{name: "wrapped", body: tc.wrappedBody},
		} {
			t.Run(tc.name+"/"+variant.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("method = %q, want GET", r.Method)
					}
					if r.URL.Path != tc.path {
						t.Errorf("path = %q, want %q", r.URL.Path, tc.path)
					}
					if got := r.Header.Get("Authorization"); got != "Bearer management-secret" {
						t.Errorf("Authorization = %q", got)
					}
					_, _ = w.Write([]byte(variant.body))
				}))
				defer server.Close()
				client := cpa.NewClient(server.URL, "management-secret", 2*time.Second, false)
				result, err := fetchProviderEndpoint(context.Background(), client, tc.name)
				if err != nil {
					t.Fatalf("fetch endpoint: %v", err)
				}
				if result.statusCode != http.StatusOK || len(result.body) == 0 {
					t.Fatalf("result metadata = status:%d body:%q", result.statusCode, string(result.body))
				}
				if len(result.keys) != 1 || result.keys[0].APIKey == "" || result.keys[0].Prefix == "" || result.keys[0].BaseURL == "" || result.keys[0].Name == "" || result.keys[0].AuthIndex == "" {
					t.Fatalf("provider payload = %#v", result.keys)
				}
			})
		}
	}
}

func TestProviderAPIKeyClientClassifiesEmptyAndBlankBodies(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		wantError bool
	}{
		{name: "direct-array", body: `[]`},
		{name: "direct-null", body: `null`},
		{name: "wrapped-array", body: `{"codex-api-key":[]}`},
		{name: "blank-body", body: "   \n", wantError: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			client := cpa.NewClient(server.URL, "management-secret", 2*time.Second, false)
			result, err := client.FetchCodexAPIKeys(context.Background())
			if tc.wantError {
				if err == nil || result == nil {
					t.Fatalf("blank body result=%#v err=%v", result, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("empty payload: %v", err)
			}
			if result.StatusCode != http.StatusOK || len(result.Payload) != 0 {
				t.Fatalf("empty result = %#v", result)
			}
		})
	}
}

func fetchProviderEndpoint(ctx context.Context, client *cpa.Client, source string) (providerEndpointResult, error) {
	switch source {
	case "codex":
		result, err := client.FetchCodexAPIKeys(ctx)
		return providerEndpointResult{statusCode: result.StatusCode, body: result.Body, keys: result.Payload}, err
	case "claude":
		result, err := client.FetchClaudeAPIKeys(ctx)
		return providerEndpointResult{statusCode: result.StatusCode, body: result.Body, keys: result.Payload}, err
	default:
		return providerEndpointResult{}, fmt.Errorf("unknown provider source: %s", source)
	}
}

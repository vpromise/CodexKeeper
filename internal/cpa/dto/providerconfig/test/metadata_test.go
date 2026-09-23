package providerconfig_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"cpa-usage-keeper/internal/cpa/dto/providerconfig"
)

func TestProviderKeyConfigDecodesSupportedFieldAliases(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		// kebab-case 是当前 CPA management response 的主格式。
		{name: "kebab-case", body: `{"api-key":"provider-key","prefix":"team","name":"Provider","base-url":"https://provider.example/v1","auth-index":"provider-auth","priority":8,"disabled":false,"note":"primary"}`},
		// snake-case 保留旧响应兼容。
		{name: "snake-case", body: `{"apiKey":"provider-key","prefix":"team","name":"Provider","base_url":"https://provider.example/v1","auth_index":"provider-auth","priority":8,"disabled":false,"note":"primary"}`},
		// camel-case 保留 client 既有别名兼容。
		{name: "camel-case", body: `{"key":"provider-key","prefix":"team","name":"Provider","baseURL":"https://provider.example/v1","authIndex":"provider-auth","priority":8,"disabled":false,"note":"primary"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cfg providerconfig.ProviderKeyConfig
			if err := json.Unmarshal([]byte(tc.body), &cfg); err != nil {
				t.Fatalf("unmarshal provider key config: %v", err)
			}
			if cfg.APIKey != "provider-key" || cfg.Prefix != "team" || cfg.Name != "Provider" || cfg.BaseURL != "https://provider.example/v1" || cfg.AuthIndex != "provider-auth" {
				t.Fatalf("provider key fields = %+v", cfg)
			}
			if cfg.Priority == nil || *cfg.Priority != 8 || cfg.Disabled == nil || *cfg.Disabled || cfg.Note == nil || *cfg.Note != "primary" {
				t.Fatalf("provider sync fields = %+v", cfg)
			}
		})
	}
}

func TestProviderKeyConfigInfersDisabledFromExcludedModels(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		wantDisabled *bool
	}{
		{name: "kebab-case wildcard", body: `{"excluded-models":["gemini-1.5-pro","*"]}`, wantDisabled: boolPtr(true)},
		{name: "snake-case wildcard", body: `{"excluded_models":[" * "]}`, wantDisabled: boolPtr(true)},
		{name: "camel-case wildcard", body: `{"excludedModels":["*"]}`, wantDisabled: boolPtr(true)},
		{name: "ordinary exclusions", body: `{"excluded-models":["gpt-*","claude-3"]}`},
		{name: "explicit disabled wins", body: `{"disabled":false,"excluded-models":["*"]}`, wantDisabled: boolPtr(false)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cfg providerconfig.ProviderKeyConfig
			if err := json.Unmarshal([]byte(tc.body), &cfg); err != nil {
				t.Fatalf("unmarshal provider key config: %v", err)
			}
			if !reflect.DeepEqual(cfg.Disabled, tc.wantDisabled) {
				t.Fatalf("disabled = %v, want %v", cfg.Disabled, tc.wantDisabled)
			}
		})
	}
}

func boolPtr(value bool) *bool {
	return &value
}

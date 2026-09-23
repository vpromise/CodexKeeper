package test

import (
	"cpa-usage-keeper/internal/quota"
	"testing"
)

func TestDefaultProviderConfigsContainsOnlyNativeAPICallTemplates(t *testing.T) {
	configs := quota.DefaultProviderConfigs()
	templates := configs.APICallTemplates()
	if len(templates) != 3 {
		t.Fatalf("expected three native templates, got %d", len(templates))
	}
	expected := []string{"https://chatgpt.com/backend-api/wham/usage", "https://api.anthropic.com/api/oauth/usage", "https://api.anthropic.com/api/oauth/profile"}
	for i, template := range templates {
		if template.Method != "GET" || template.URL != expected[i] || template.Headers["Authorization"] != "Bearer $TOKEN$" {
			t.Fatalf("unexpected template: %+v", template)
		}
	}
}

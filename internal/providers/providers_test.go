package providers

import "testing"

func TestUsageProvider(t *testing.T) {
	for _, test := range []struct {
		provider, executor, want string
		supported                bool
	}{
		{" Codex ", "CodexExecutor", "codex", true},
		{"claude", "ClaudeExecutor", "claude", true},
		{"", "CodexWebsocketsExecutor", "codex", true},
		{"", "ClaudeExecutor", "claude", true},
		{"antigravity", "ClaudeExecutor", "antigravity", false},
		{"openai", "CodexExecutor", "openai", false},
		{"", "OpenAICompatExecutor", "", false},
		{"", "", "", false},
	} {
		t.Run(test.provider+"/"+test.executor, func(t *testing.T) {
			got, supported := UsageProvider(test.provider, test.executor)
			if got != test.want || supported != test.supported {
				t.Fatalf("UsageProvider = (%q, %v), want (%q, %v)", got, supported, test.want, test.supported)
			}
		})
	}
}

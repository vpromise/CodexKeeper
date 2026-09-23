// Package providers defines the native providers supported by CodexKeeper.
package providers

import "strings"

func Normalize(provider string) string {
	return strings.ToLower(strings.TrimSpace(provider))
}

func Supported(provider string) bool {
	switch Normalize(provider) {
	case "codex", "claude":
		return true
	default:
		return false
	}
}

// UsageProvider accepts an explicit native provider, or a known native executor
// when the provider is absent. Model names never establish provider identity.
func UsageProvider(provider, executor string) (string, bool) {
	if provider = Normalize(provider); provider != "" {
		return provider, Supported(provider)
	}
	switch Normalize(executor) {
	case "codexexecutor", "codexwebsocketsexecutor", "codexautoexecutor":
		return "codex", true
	case "claudeexecutor":
		return "claude", true
	default:
		return "", false
	}
}

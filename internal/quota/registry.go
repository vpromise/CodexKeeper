package quota

import "strings"

type ProviderRegistry struct {
	handlers map[string]ProviderHandler
}

func NewDefaultProviderRegistry(caller ManagementClient, configs ProviderConfigs) ProviderRegistry {
	return NewProviderRegistry(map[string]ProviderHandler{
		"codex":  NewCodexProvider(caller, configs.Codex),
		"claude": NewClaudeProvider(caller, configs.ClaudeUsage, configs.ClaudeProfile),
	})
}

func NewProviderRegistry(handlers map[string]ProviderHandler) ProviderRegistry {
	registry := ProviderRegistry{handlers: make(map[string]ProviderHandler, len(handlers))}
	for identityType, handler := range handlers {
		identityType = normalizeIdentityType(identityType)
		if identityType == "" || handler == nil {
			continue
		}
		registry.handlers[identityType] = handler
	}
	return registry
}

func (r ProviderRegistry) Provider(identityType string) (ProviderHandler, bool) {
	handler, ok := r.handlers[normalizeIdentityType(identityType)]
	return handler, ok
}

func normalizeIdentityType(identityType string) string {
	return strings.ToLower(strings.TrimSpace(identityType))
}

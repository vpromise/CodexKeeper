package providerconfig

import (
	"encoding/json"
	"fmt"
)

// ProviderKeyConfig 是标准 API key provider 配置的兼容归一化视图，支持 CPA 返回的多种 key 命名。
type ProviderKeyConfig struct {
	APIKey         string
	Prefix         string
	Name           string
	BaseURL        string
	AuthIndex      string
	Priority       *int
	Disabled       *bool
	ExcludedModels []string
	Note           *string
}

func (p *ProviderKeyConfig) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode provider key config: %w", err)
	}
	p.APIKey = firstString(raw, "apiKey", "api-key", "key")
	p.Prefix = firstString(raw, "prefix")
	p.Name = firstString(raw, "name")
	p.BaseURL = firstString(raw, "base-url", "base_url", "baseURL")
	p.AuthIndex = firstString(raw, "auth-index", "auth_index", "authIndex")
	p.Priority = firstInt(raw, "priority")
	p.Disabled = firstBool(raw, "disabled")
	p.ExcludedModels = firstStringList(raw, "excluded-models", "excluded_models", "excludedModels")
	// CPA 用 excluded-models 中精确的 * 表示普通 Provider 整条停用；显式 disabled 始终优先。
	if p.Disabled == nil && stringListContains(raw, "*", "excluded-models", "excluded_models", "excludedModels") {
		disabled := true
		p.Disabled = &disabled
	}
	p.Note = firstStringPtr(raw, "note")
	return nil
}

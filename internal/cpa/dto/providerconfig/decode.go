package providerconfig

import (
	"strings"
)

func firstString(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		if text, ok := value.(string); ok {
			return text
		}
	}
	return ""
}

func firstStringPtr(raw map[string]any, keys ...string) *string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		if text, ok := value.(string); ok {
			return &text
		}
	}
	return nil
}

func firstInt(raw map[string]any, keys ...string) *int {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			converted := int(typed)
			return &converted
		case int:
			return &typed
		}
	}
	return nil
}

func firstBool(raw map[string]any, keys ...string) *bool {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		if typed, ok := value.(bool); ok {
			return &typed
		}
	}
	return nil
}

func stringListContains(raw map[string]any, target string, keys ...string) bool {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		items, ok := value.([]any)
		if !ok {
			continue
		}
		for _, item := range items {
			text, ok := item.(string)
			if ok && strings.TrimSpace(text) == target {
				return true
			}
		}
	}
	return false
}

func firstStringList(raw map[string]any, keys ...string) []string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		items, ok := value.([]any)
		if !ok {
			continue
		}
		out := make([]string, 0, len(items))
		for _, item := range items {
			if text, ok := item.(string); ok {
				out = append(out, text)
			}
		}
		return out
	}
	return nil
}

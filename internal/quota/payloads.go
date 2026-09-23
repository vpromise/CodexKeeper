package quota

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"cpa-usage-keeper/internal/cpa/dto/apicall"
)

func parseCodexUsagePayload(response *apicall.Response) (*CodexUsagePayload, error) {
	object, err := parseResponseObject(response)
	if err != nil {
		return nil, err
	}
	payload := &CodexUsagePayload{
		PlanType:              stringField(object, "plan_type", "planType"),
		RateLimit:             parseCodexRateLimitInfo(objectField(object, "rate_limit", "rateLimit")),
		CodeReviewRateLimit:   parseCodexRateLimitInfo(objectField(object, "code_review_rate_limit", "codeReviewRateLimit")),
		RateLimitResetCredits: parseCodexRateLimitResetCredits(objectField(object, "rate_limit_reset_credits", "rateLimitResetCredits")),
	}
	for _, raw := range arrayField(object, "additional_rate_limits", "additionalRateLimits") {
		limitObject := rawObject(raw)
		if limitObject == nil {
			continue
		}
		payload.AdditionalRateLimits = append(payload.AdditionalRateLimits, CodexAdditionalRateLimit{
			LimitName:      stringField(limitObject, "limit_name", "limitName"),
			MeteredFeature: stringField(limitObject, "metered_feature", "meteredFeature"),
			RateLimit:      parseCodexRateLimitInfo(objectField(limitObject, "rate_limit", "rateLimit")),
		})
	}
	return payload, nil
}

func parseCodexRateLimitResetCredits(object map[string]json.RawMessage) *CodexRateLimitResetCredits {
	if object == nil {
		return nil
	}
	// 官方刷新结果可能省略 reset credit 字段；只有明确返回 available_count 时才暴露给前端。
	availableCount := intPtrField(object, "available_count", "availableCount")
	if availableCount == nil {
		return nil
	}
	count := int(*availableCount)
	return &CodexRateLimitResetCredits{AvailableCount: &count}
}

func parseCodexRateLimitInfo(object map[string]json.RawMessage) *CodexRateLimitInfo {
	if object == nil {
		return nil
	}
	return &CodexRateLimitInfo{
		Allowed:         boolPtrField(object, "allowed"),
		LimitReached:    boolPtrField(object, "limit_reached", "limitReached"),
		PrimaryWindow:   parseCodexUsageWindow(objectField(object, "primary_window", "primaryWindow")),
		SecondaryWindow: parseCodexUsageWindow(objectField(object, "secondary_window", "secondaryWindow")),
	}
}

func parseCodexUsageWindow(object map[string]json.RawMessage) *CodexUsageWindow {
	if object == nil {
		return nil
	}
	// 逐字段保留 presence，历史采集必须区分缺失字段和上游明确返回的零值。
	usedPercent, hasUsedPercent := floatValue(object, "used_percent", "usedPercent")
	limitWindowSeconds, hasLimitWindowSeconds := floatValue(object, "limit_window_seconds", "limitWindowSeconds")
	resetAfterSeconds, hasResetAfterSeconds := floatValue(object, "reset_after_seconds", "resetAfterSeconds")
	resetAt, hasResetAt := floatValue(object, "reset_at", "resetAt")
	return &CodexUsageWindow{
		// 数值字段保持既有零值/截断语义，presence 字段只为历史提取提供额外事实。
		UsedPercent:           usedPercent,
		LimitWindowSeconds:    int64(limitWindowSeconds),
		ResetAfterSeconds:     int64(resetAfterSeconds),
		ResetAt:               int64(resetAt),
		WindowUsageTokens:     intPtrField(object, "window_usage_tokens", "windowUsageTokens"),
		WindowUsageCost:       floatPtrField(object, "window_usage_cost", "windowUsageCost"),
		HasUsedPercent:        hasUsedPercent,
		HasLimitWindowSeconds: hasLimitWindowSeconds,
		HasResetAfterSeconds:  hasResetAfterSeconds,
		HasResetAt:            hasResetAt,
	}
}

func parseClaudeUsagePayload(response *apicall.Response) (*ClaudeUsagePayload, error) {
	object, err := parseResponseObject(response)
	if err != nil {
		return nil, err
	}
	return &ClaudeUsagePayload{
		FiveHour:          parseClaudeUsageWindow(objectField(object, "five_hour", "fiveHour")),
		SevenDay:          parseClaudeUsageWindow(objectField(object, "seven_day", "sevenDay")),
		SevenDayOAuthApps: parseClaudeUsageWindow(objectField(object, "seven_day_oauth_apps", "sevenDayOauthApps")),
		SevenDayOpus:      parseClaudeUsageWindow(objectField(object, "seven_day_opus", "sevenDayOpus")),
		SevenDaySonnet:    parseClaudeUsageWindow(objectField(object, "seven_day_sonnet", "sevenDaySonnet")),
		SevenDayCowork:    parseClaudeUsageWindow(objectField(object, "seven_day_cowork", "sevenDayCowork")),
		IguanaNecktie:     parseClaudeUsageWindow(objectField(object, "iguana_necktie", "iguanaNecktie")),
		ExtraUsage:        parseClaudeExtraUsage(objectField(object, "extra_usage", "extraUsage")),
	}, nil
}

func parseClaudeUsageWindow(object map[string]json.RawMessage) *ClaudeUsageWindow {
	if object == nil {
		return nil
	}
	return &ClaudeUsageWindow{
		Utilization: floatField(object, "utilization"),
		ResetsAt:    stringField(object, "resets_at", "resetsAt"),
	}
}

func parseClaudeExtraUsage(object map[string]json.RawMessage) *ClaudeExtraUsage {
	if object == nil {
		return nil
	}
	return &ClaudeExtraUsage{
		IsEnabled:    boolField(object, "is_enabled", "isEnabled"),
		MonthlyLimit: floatField(object, "monthly_limit", "monthlyLimit"),
		UsedCredits:  floatField(object, "used_credits", "usedCredits"),
		Utilization:  floatPtrField(object, "utilization"),
	}
}

func parseClaudeProfilePayload(response *apicall.Response) (*ClaudeProfileResponse, error) {
	object, err := parseResponseObject(response)
	if err != nil {
		return nil, err
	}
	return &ClaudeProfileResponse{
		Account:      parseClaudeProfileAccount(objectField(object, "account")),
		Organization: parseClaudeProfileOrganization(objectField(object, "organization")),
	}, nil
}

func parseClaudeProfileAccount(object map[string]json.RawMessage) *ClaudeProfileAccount {
	if object == nil {
		return nil
	}
	return &ClaudeProfileAccount{
		UUID:         stringField(object, "uuid"),
		FullName:     stringField(object, "full_name", "fullName"),
		DisplayName:  stringField(object, "display_name", "displayName"),
		Email:        stringField(object, "email"),
		HasClaudeMax: boolPtrField(object, "has_claude_max", "hasClaudeMax"),
		HasClaudePro: boolPtrField(object, "has_claude_pro", "hasClaudePro"),
	}
}

func parseClaudeProfileOrganization(object map[string]json.RawMessage) *ClaudeProfileOrganization {
	if object == nil {
		return nil
	}
	return &ClaudeProfileOrganization{
		UUID:                 stringField(object, "uuid"),
		Name:                 stringField(object, "name"),
		OrganizationType:     stringField(object, "organization_type", "organizationType"),
		BillingType:          stringField(object, "billing_type", "billingType"),
		RateLimitTier:        stringField(object, "rate_limit_tier", "rateLimitTier"),
		HasExtraUsageEnabled: boolField(object, "has_extra_usage_enabled", "hasExtraUsageEnabled"),
		SubscriptionStatus:   stringField(object, "subscription_status", "subscriptionStatus"),
	}
}

func parseCodexResetCreditResponse(response *apicall.Response) (ProviderResetOutput, error) {
	object, err := parseResponseObject(response)
	if err != nil {
		return ProviderResetOutput{}, err
	}
	// consume 成功必须返回 code=reset 和 windows_reset，避免把未知成功体当成已重置。
	code := stringField(object, "code")
	if code != "reset" {
		return ProviderResetOutput{}, fmt.Errorf("unexpected reset credit response code: %s", code)
	}
	windowsReset := intPtrField(object, "windows_reset", "windowsReset")
	if windowsReset == nil {
		return ProviderResetOutput{}, fmt.Errorf("invalid reset credit response")
	}
	return ProviderResetOutput{
		Code:         code,
		WindowsReset: int(*windowsReset),
	}, nil
}

func parseCodexResetCreditsResponse(response *apicall.Response) (ProviderResetCreditsOutput, error) {
	object, err := parseResponseObject(response)
	if err != nil {
		return ProviderResetCreditsOutput{}, err
	}
	if _, hasCredits := object["credits"]; !hasCredits {
		if _, hasCount := object["available_count"]; !hasCount {
			if _, hasCamelCount := object["availableCount"]; !hasCamelCount {
				return ProviderResetCreditsOutput{}, fmt.Errorf("invalid reset credits response")
			}
		}
	}
	credits := make([]CodexRateLimitResetCredit, 0)
	for _, raw := range arrayField(object, "credits") {
		creditObject := rawObject(raw)
		if creditObject == nil || stringField(creditObject, "reset_type", "resetType") != "codex_rate_limits" || stringField(creditObject, "status") != "available" {
			continue
		}
		expiresAt := stringField(creditObject, "expires_at", "expiresAt")
		if expiresAt == "" {
			continue
		}
		credits = append(credits, CodexRateLimitResetCredit{
			ID:        stringField(creditObject, "id"),
			Status:    "available",
			GrantedAt: stringField(creditObject, "granted_at", "grantedAt"),
			ExpiresAt: expiresAt,
		})
	}
	var availableCount *int
	if count := intPtrField(object, "available_count", "availableCount"); count != nil && *count >= 0 {
		parsedCount := int(*count)
		availableCount = &parsedCount
	}
	return ProviderResetCreditsOutput{AvailableCount: availableCount, Credits: credits}, nil
}

func parseResponseObject(response *apicall.Response) (map[string]json.RawMessage, error) {
	if response == nil {
		return nil, fmt.Errorf("missing quota response")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, targetHTTPError(response)
	}
	if object := rawObject(response.Body); object != nil {
		return object, nil
	}
	trimmed := strings.TrimSpace(response.BodyText)
	if trimmed == "" {
		return nil, fmt.Errorf("empty quota response body")
	}
	if object := rawObject([]byte(trimmed)); object != nil {
		return object, nil
	}
	return nil, fmt.Errorf("parse quota response body")
}

func rawObject(data []byte) map[string]json.RawMessage {
	if len(data) == 0 {
		return nil
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err == nil && object != nil {
		return object
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		return rawObject([]byte(strings.TrimSpace(text)))
	}
	return nil
}

func objectField(object map[string]json.RawMessage, keys ...string) map[string]json.RawMessage {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			return rawObject(raw)
		}
	}
	return nil
}

func arrayField(object map[string]json.RawMessage, keys ...string) []json.RawMessage {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			var values []json.RawMessage
			if err := json.Unmarshal(raw, &values); err == nil {
				return values
			}
		}
	}
	return nil
}

func stringField(object map[string]json.RawMessage, keys ...string) string {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			var text string
			if err := json.Unmarshal(raw, &text); err == nil {
				return strings.TrimSpace(text)
			}
			var number json.Number
			if err := json.Unmarshal(raw, &number); err == nil {
				return number.String()
			}
		}
	}
	return ""
}

func floatField(object map[string]json.RawMessage, keys ...string) float64 {
	value, _ := floatValue(object, keys...)
	return value
}

func floatPtrField(object map[string]json.RawMessage, keys ...string) *float64 {
	value, ok := floatValue(object, keys...)
	if !ok {
		return nil
	}
	return &value
}

func quotaFractionPtrField(object map[string]json.RawMessage, keys ...string) *float64 {
	for _, key := range keys {
		raw, ok := object[key]
		if !ok || rawJSONNull(raw) {
			continue
		}
		var number float64
		if err := json.Unmarshal(raw, &number); err == nil {
			if !math.IsNaN(number) && !math.IsInf(number, 0) {
				return &number
			}
			continue
		}
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		scale := 1.0
		if strings.HasSuffix(text, "%") {
			text = strings.TrimSpace(strings.TrimSuffix(text, "%"))
			scale = 0.01
		}
		parsed, err := strconv.ParseFloat(text, 64)
		if err != nil {
			continue
		}
		parsed *= scale
		if !math.IsNaN(parsed) && !math.IsInf(parsed, 0) {
			return &parsed
		}
	}
	return nil
}

func floatValue(object map[string]json.RawMessage, keys ...string) (float64, bool) {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			if rawJSONNull(raw) {
				continue
			}
			var number float64
			if err := json.Unmarshal(raw, &number); err == nil {
				return number, true
			}
			var text string
			if err := json.Unmarshal(raw, &text); err == nil {
				parsed, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
				if err == nil {
					return parsed, true
				}
			}
		}
	}
	return 0, false
}

func rawJSONNull(raw json.RawMessage) bool {
	return len(raw) == 4 && raw[0] == 'n' && raw[1] == 'u' && raw[2] == 'l' && raw[3] == 'l'
}

func intField(object map[string]json.RawMessage, keys ...string) int64 {
	value, ok := floatValue(object, keys...)
	if !ok {
		return 0
	}
	return int64(value)
}

func intPtrField(object map[string]json.RawMessage, keys ...string) *int64 {
	value, ok := floatValue(object, keys...)
	if !ok {
		return nil
	}
	parsed := int64(value)
	return &parsed
}

func boolField(object map[string]json.RawMessage, keys ...string) bool {
	value := boolPtrField(object, keys...)
	return value != nil && *value
}

func boolPtrField(object map[string]json.RawMessage, keys ...string) *bool {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			if rawJSONNull(raw) {
				continue
			}
			var value bool
			if err := json.Unmarshal(raw, &value); err == nil {
				return &value
			}
			var text string
			if err := json.Unmarshal(raw, &text); err == nil {
				switch strings.ToLower(strings.TrimSpace(text)) {
				case "true", "1", "yes", "y", "on":
					parsed := true
					return &parsed
				case "false", "0", "no", "n", "off":
					parsed := false
					return &parsed
				}
			}
		}
	}
	return nil
}

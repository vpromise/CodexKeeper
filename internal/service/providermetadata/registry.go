package providermetadata

import "fmt"

// providerSources is the complete native provider metadata surface.
func providerSources() []source {
	return []source{codexSource(), claudeSource()}
}

// validateSources 拒绝空标识和重复 source/type，防止 registry 维护错误改变 stale 范围。
func validateSources(sources []source) error {
	// seenIDs 记录 registry 内部 source ID。
	seenIDs := make(map[string]struct{}, len(sources))
	// seenTypes 记录最终写入的 provider type。
	seenTypes := make(map[string]struct{}, len(sources))
	// 按声明顺序检查每个固定来源。
	for _, item := range sources {
		// 空 source ID 无法生成稳定 warning。
		if item.id == "" {
			return fmt.Errorf("provider metadata source id is required")
		}
		// 空 provider type 无法形成正确 stale scope。
		if item.providerType == "" {
			return fmt.Errorf("provider metadata source %q type is required", item.id)
		}
		// 空默认名会让缺 name 的正常 CPA entry 被错误过滤。
		if item.defaultDisplayName == "" {
			return fmt.Errorf("provider metadata source %q default display name is required", item.id)
		}
		// 空 warning 名无法保持来源错误文本稳定。
		if item.warningName == "" {
			return fmt.Errorf("provider metadata source %q warning name is required", item.id)
		}
		// 每个 registry source 必须绑定唯一 endpoint 执行函数。
		if item.fetch == nil {
			return fmt.Errorf("provider metadata source %q fetch is required", item.id)
		}
		// 重复 source ID 表示 registry 声明冲突。
		if _, ok := seenIDs[item.id]; ok {
			return fmt.Errorf("provider metadata source id %q is duplicated", item.id)
		}
		// 记录已通过校验的 source ID。
		seenIDs[item.id] = struct{}{}
		// 重复 provider type 会让两个 endpoint 互相 stale，必须拒绝。
		if _, ok := seenTypes[item.providerType]; ok {
			return fmt.Errorf("provider metadata type %q is duplicated", item.providerType)
		}
		// 记录已通过校验的 provider type。
		seenTypes[item.providerType] = struct{}{}
	}
	// 全部固定来源通过校验。
	return nil
}

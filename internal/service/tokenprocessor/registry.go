package tokenprocessor

import (
	"fmt"
	"strings"
)

// executorDefinition 把一个 CPA Go executor alias 显式绑定到已有协议 handler。
type executorDefinition struct {
	alias     string
	handlerID HandlerID
}

// identityAliasDefinition 只定义旧事件 fallback hint，绝不代表 CPA parser contract。
type identityAliasDefinition struct {
	alias     string
	handlerID HandlerID
}

// registryIndexes 是只读静态索引；所有 resolver 调用共享同一份校验后的 registry。
type registryIndexes struct {
	executors  map[string]executorDefinition
	identities map[string]identityAliasDefinition
}

// handlerDefinitions 明确列出唯一协议 handler，registry 校验用它阻止 executor 指向悬空实现。
var handlerDefinitions = []HandlerID{HandlerClaude, HandlerResponsesInclusive, HandlerStrictPassThrough}

// tokenHandlers 是协议实现的唯一 registry；executor 与 identity registry 只能引用这里的 HandlerID。
var tokenHandlers = map[HandlerID]tokenHandler{HandlerClaude: claudeHandler{}, HandlerResponsesInclusive: responsesInclusiveHandler{}, HandlerStrictPassThrough: strictPassThroughHandler{}}

// executorDefinitions 是中央静态引用表；新增 executor 必须同时拥有独立文件、此处一条引用和黑盒测试。
var executorDefinitions = []executorDefinition{claudeExecutorDefinition, codexExecutorDefinition, codexWebsocketsExecutorDefinition, codexAutoExecutorDefinition}

// identityAliasDefinitions 原样保留 Keeper 已验证的 identity 路由，但每一项都只有 identity_hint 强度。
var identityAliasDefinitions = []identityAliasDefinition{{alias: "claude", handlerID: HandlerClaude}, {alias: "codex", handlerID: HandlerResponsesInclusive}}

// defaultRegistry 在包加载时一次性构建并校验静态表；这里没有 init 自注册或运行时插件顺序问题。
var defaultRegistry, defaultRegistryErr = buildRegistryIndexes()

func buildRegistryIndexes() (registryIndexes, error) {
	// 先校验 handler ID 唯一，避免 alias 指向含糊的协议实现。
	handlers := make(map[HandlerID]struct{}, len(handlerDefinitions))
	for _, handlerID := range handlerDefinitions {
		if _, exists := handlers[handlerID]; exists {
			return registryIndexes{}, fmt.Errorf("duplicate token handler id %q", handlerID)
		}
		handlers[handlerID] = struct{}{}
	}
	// 每个声明的 handler 都必须有且只有一个纯计算实现，防止 resolver 通过但 Process 才发现悬空实现。
	if len(tokenHandlers) != len(handlers) {
		return registryIndexes{}, fmt.Errorf("token handler registry size %d does not match definitions %d", len(tokenHandlers), len(handlers))
	}
	for handlerID, handler := range tokenHandlers {
		if _, exists := handlers[handlerID]; !exists {
			return registryIndexes{}, fmt.Errorf("token handler implementation %q is not declared", handlerID)
		}
		if handler == nil || handler.ID() != handlerID {
			return registryIndexes{}, fmt.Errorf("token handler implementation %q reports mismatched id", handlerID)
		}
	}

	// executor alias 经过 trim/lower 后必须唯一，保证 exact match 的结果稳定。
	executors := make(map[string]executorDefinition, len(executorDefinitions))
	for _, definition := range executorDefinitions {
		alias := normalizeRegistryAlias(definition.alias)
		if alias == "" {
			return registryIndexes{}, fmt.Errorf("empty token executor alias")
		}
		if _, exists := handlers[definition.handlerID]; !exists {
			return registryIndexes{}, fmt.Errorf("executor alias %q references unknown handler %q", definition.alias, definition.handlerID)
		}
		if _, exists := executors[alias]; exists {
			return registryIndexes{}, fmt.Errorf("duplicate token executor alias %q", definition.alias)
		}
		executors[alias] = definition
	}

	// identity alias 也必须唯一，但它们只生成 identity_hint，不与 executor registry 合并权限。
	identities := make(map[string]identityAliasDefinition, len(identityAliasDefinitions))
	for _, definition := range identityAliasDefinitions {
		alias := normalizeRegistryAlias(definition.alias)
		if alias == "" {
			return registryIndexes{}, fmt.Errorf("empty token identity alias")
		}
		if _, exists := handlers[definition.handlerID]; !exists {
			return registryIndexes{}, fmt.Errorf("identity alias %q references unknown handler %q", definition.alias, definition.handlerID)
		}
		if _, exists := identities[alias]; exists {
			return registryIndexes{}, fmt.Errorf("duplicate token identity alias %q", definition.alias)
		}
		identities[alias] = definition
	}

	return registryIndexes{executors: executors, identities: identities}, nil
}

func normalizeRegistryAlias(value string) string {
	// alias 只做大小写和首尾空白归一化，不改写内部字符，防止相似名称被误识别成已知协议。
	return strings.ToLower(strings.TrimSpace(value))
}

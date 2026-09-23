package tokenprocessor

type strictPassThroughHandler struct{}

func (strictPassThroughHandler) ID() HandlerID { return HandlerStrictPassThrough }

func (strictPassThroughHandler) Normalize(context normalizationContext) handlerResult {
	// Unknown contracts retain cache aliases and zero-only total reconciliation.
	tokens, actions, violations := normalizeNonClaudeCacheRead(context)
	// cache alias 自己的损坏不能关闭只依赖 Input/Output 的既有零 Total 回填。
	return handlerResult{tokens: tokens, authority: totalAuthorityZeroOnly, actions: actions, violations: violations}
}

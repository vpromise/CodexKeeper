import { describe, expect, it } from 'vitest'
import { buildCredentialProviderFilterOptions, normalizeCredentialProviderFilterKey } from '../credentialProviderFilters'
describe('native credential filters', () => {
  it.each(['auth-files', 'ai-provider'] as const)('counts only native credentials in %s', scope => {
    const options = buildCredentialProviderFilterOptions(scope, [{type: 'codex', count: 3}, {type: 'claude', count: 2}, {type: 'gemini', count: 8}])
    expect(options[0].count).toBe(5)
    expect(options.map(option => option.key).sort()).toEqual(['all', 'claude', 'codex'])
    expect(normalizeCredentialProviderFilterKey(scope, 'gemini')).toBe('all')
  })
  it('ignores invalid counts', () => {
    expect(buildCredentialProviderFilterOptions('ai-provider', [{type: 'codex', count: NaN}, {type: 'claude', count: -1}])).toEqual([])
  })
})

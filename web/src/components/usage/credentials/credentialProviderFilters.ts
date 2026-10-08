import type { UsageIdentityTypeCount } from '@/lib/types'
import type { ProviderBrandIconKey } from '@/components/ProviderBrandIcon'

export type CredentialProviderFilterScope = 'auth-files' | 'ai-provider'
export type KnownCredentialProviderFilterKey = ProviderBrandIconKey
export type CredentialProviderFilterKey = 'all' | KnownCredentialProviderFilterKey

export interface CredentialProviderFilterOption {
  key: CredentialProviderFilterKey
  count: number
  labelKey: string
  knownKey?: KnownCredentialProviderFilterKey
}

interface KnownCredentialProviderFilter {
  key: KnownCredentialProviderFilterKey
  labelKey: string
  types: string[]
}

const AUTH_FILE_PROVIDER_FILTERS: KnownCredentialProviderFilter[] = [
  { key: 'claude', labelKey: 'usage_stats.credentials_filter_claude', types: ['claude'] },
  { key: 'codex', labelKey: 'usage_stats.credentials_filter_codex', types: ['codex'] },
]

const AI_PROVIDER_FILTERS: KnownCredentialProviderFilter[] = [
  { key: 'codex', labelKey: 'usage_stats.credentials_filter_codex', types: ['codex'] },
  { key: 'claude', labelKey: 'usage_stats.credentials_filter_claude', types: ['claude'] },
]

const FILTERS_BY_SCOPE: Record<CredentialProviderFilterScope, KnownCredentialProviderFilter[]> = {
  'auth-files': AUTH_FILE_PROVIDER_FILTERS,
  'ai-provider': AI_PROVIDER_FILTERS,
}

function credentialProviderFiltersForScope(scope: CredentialProviderFilterScope): KnownCredentialProviderFilter[] {
  return FILTERS_BY_SCOPE[scope]
}

export function credentialProviderFilterTypes(scope: CredentialProviderFilterScope, filter: CredentialProviderFilterKey): string[] {
  if (filter === 'all') {
    return []
  }
  return credentialProviderFiltersForScope(scope).find((item) => item.key === filter)?.types ?? []
}

export function normalizeCredentialProviderFilterKey(scope: CredentialProviderFilterScope, value: unknown): CredentialProviderFilterKey {
  if (typeof value !== 'string' || value === 'all') {
    return 'all'
  }
  return credentialProviderFiltersForScope(scope).some((item) => item.key === value)
    ? value as CredentialProviderFilterKey
    : 'all'
}

// null 表示没有有效的链接筛选；显式 all 必须区别于回退到本地偏好。
export function resolveCredentialProviderFilterFromSearch(scope: CredentialProviderFilterScope, search: string): CredentialProviderFilterKey | null {
  const value = new URLSearchParams(search).get('provider')
  return value !== null && normalizeCredentialProviderFilterKey(scope, value) === value
    ? value as CredentialProviderFilterKey
    : null
}

export function buildCredentialProviderSearch(search: string, filter: CredentialProviderFilterKey): string {
  const params = new URLSearchParams(search)
  params.set('provider', filter)
  return `?${params.toString()}`
}

export function buildCredentialProviderFilterOptions(scope: CredentialProviderFilterScope, typeCounts: UsageIdentityTypeCount[], selected: CredentialProviderFilterKey = 'all'): CredentialProviderFilterOption[] {
  const countsByType = new Map<string, number>()
  let allCount = 0

  for (const item of typeCounts) {
    if (item.type !== 'codex' && item.type !== 'claude') continue
    const count = finiteCount(item.count)
    if (count <= 0) {
      continue
    }
    allCount += count
    countsByType.set(item.type, (countsByType.get(item.type) ?? 0) + count)
  }

  if (allCount <= 0 && selected === 'all') {
    return []
  }

  const options: CredentialProviderFilterOption[] = [{ key: 'all', labelKey: 'usage_stats.credentials_filter_all', count: allCount }]

  // Only native credential types contribute to counts and filter buttons.
  for (const filter of credentialProviderFiltersForScope(scope)) {
    const count = filter.types.reduce((sum, type) => sum + (countsByType.get(type) ?? 0), 0)
    // 链接指定的供应商即使暂时没有凭据，也保留选中项和返回 All 的入口。
    if (count <= 0 && filter.key !== selected) {
      continue
    }
    options.push({ key: filter.key, labelKey: filter.labelKey, count, knownKey: filter.key })
  }

  return options
}

function finiteCount(value: number): number {
  return Number.isFinite(value) && value > 0 ? value : 0
}

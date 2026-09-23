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

export function buildCredentialProviderFilterOptions(scope: CredentialProviderFilterScope, typeCounts: UsageIdentityTypeCount[]): CredentialProviderFilterOption[] {
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

  if (allCount <= 0) {
    return []
  }

  const options: CredentialProviderFilterOption[] = [{ key: 'all', labelKey: 'usage_stats.credentials_filter_all', count: allCount }]

  // Only native credential types contribute to counts and filter buttons.
  for (const filter of credentialProviderFiltersForScope(scope)) {
    const count = filter.types.reduce((sum, type) => sum + (countsByType.get(type) ?? 0), 0)
    if (count <= 0) {
      continue
    }
    options.push({ key: filter.key, labelKey: filter.labelKey, count, knownKey: filter.key })
  }

  return options
}

function finiteCount(value: number): number {
  return Number.isFinite(value) && value > 0 ? value : 0
}

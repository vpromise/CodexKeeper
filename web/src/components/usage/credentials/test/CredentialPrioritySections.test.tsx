import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it, vi } from 'vitest'
import type { UsageIdentity } from '@/lib/types'
import { AiProviderCredentialsSection } from '../AiProviderCredentialsSection'
import { AuthFileCredentialsSection } from '../AuthFileCredentialsSection'
import { buildAiProviderCredentialRows, buildAuthFileCredentialRows } from '../credentialViewModels'
import { createAiProviderSectionProps, createAuthFileSectionProps } from './credentialSectionFixtures'

vi.mock('react-i18next', () => ({
  initReactI18next: { type: '3rdParty', init: () => undefined },
  useTranslation: () => ({ t: (key: string) => key }),
}))

const identity = (authType: 1 | 2): UsageIdentity => ({
  id: String(authType), identity: `index-${authType}`, auth_type: authType,
  name: 'Credential', type: 'codex', provider: 'provider', is_deleted: false,
} as UsageIdentity)

describe('credential priority editing is limited to auth files', () => {
  it('renders an actionable P0 only for an auth file with no priority', () => {
    const onSavePriority = vi.fn(async () => undefined)
    const authRows = buildAuthFileCredentialRows([identity(1)])
    const aiRows = buildAiProviderCredentialRows([identity(2)])
    const authHTML = renderToStaticMarkup(<AuthFileCredentialsSection {...createAuthFileSectionProps({ rows: authRows, total: 1, onSavePriority })} />)
    const aiHTML = renderToStaticMarkup(<AiProviderCredentialsSection {...createAiProviderSectionProps({ rows: aiRows, total: 1 })} />)
    expect(authHTML).toContain('>P0</button>')
    expect(authHTML).toContain('usage_stats.credentials_priority_edit')
    expect(aiHTML).not.toContain('usage_stats.credentials_priority_edit')
  })

  it.each([
    { type: 'codex', is_deleted: true },
    { type: 'antigravity', is_deleted: false },
  ])('does not offer priority editing for unavailable auth files: $type / $is_deleted', (overrides) => {
    const rows = buildAuthFileCredentialRows([{ ...identity(1), ...overrides }])
    const html = renderToStaticMarkup(<AuthFileCredentialsSection {...createAuthFileSectionProps({ rows, total: 1, onSavePriority: async () => undefined })} />)
    expect(html).not.toContain('usage_stats.credentials_priority_edit')
  })
})

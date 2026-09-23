import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it, vi } from 'vitest'
import { CredentialProviderFilterBar } from '../CredentialProviderFilterBar'
vi.mock('react-i18next', () => ({initReactI18next: {type: '3rdParty', init: () => undefined}, useTranslation: () => ({t: (key: string) => key})}))
describe('native provider filter bar', () => {
  it('keeps Codex and Claude branding and excludes removed providers', () => {
    const html = renderToStaticMarkup(<CredentialProviderFilterBar scope="auth-files" typeCounts={[{type: 'codex', count: 2}, {type: 'claude', count: 3}, {type: 'gemini', count: 10}]} value="all" onChange={() => undefined} />)
    expect(html).toContain('data-provider-brand-icon="codex"')
    expect(html).toContain('data-provider-brand-icon="claude"')
    expect(html).not.toContain('gemini')
    expect(html).toContain('>5</span>')
  })
  it('hides an empty filter bar', () => {
    expect(renderToStaticMarkup(<CredentialProviderFilterBar scope="auth-files" typeCounts={[]} value="all" onChange={() => undefined} />)).toBe('')
  })
})

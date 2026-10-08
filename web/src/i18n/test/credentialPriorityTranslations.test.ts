import { describe, expect, it } from 'vitest'
import i18n, { SUPPORTED_LANGUAGES } from '../index'

describe('credential priority translations', () => {
  it.each(SUPPORTED_LANGUAGES)('includes the editor and errors in %s', (language) => {
    const labels = i18n.getResourceBundle(language, 'translation').usage_stats
    for (const key of [
      'credentials_priority_edit', 'credentials_priority_input', 'credentials_priority_save',
      'credentials_priority_cancel', 'credentials_priority_saving', 'credentials_priority_invalid',
      'credentials_priority_order_hint',
      'credentials_priority_save_success', 'credentials_priority_save_failed',
      'credentials_priority_stale_target', 'credentials_priority_conflict_auth_file',
    ]) {
      expect(labels[key]).toEqual(expect.any(String))
      expect(labels[key].length).toBeGreaterThan(0)
    }
  })
})

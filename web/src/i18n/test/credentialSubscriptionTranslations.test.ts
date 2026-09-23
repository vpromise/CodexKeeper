import { describe, expect, it } from 'vitest'
import i18n, { SUPPORTED_LANGUAGES } from '../index'

const OFFICIAL_PLAN_LABELS = {
  'usage_stats.credentials_subscription_claude_free': 'Free',
  'usage_stats.credentials_subscription_claude_pro': 'Pro',
  'usage_stats.credentials_subscription_claude_max': 'Max',
  'usage_stats.credentials_subscription_claude_team': 'Team',
} as const

describe('credential subscription translations', () => {
  it('keeps official plan names unchanged across languages', () => {
    for (const language of SUPPORTED_LANGUAGES) {
      for (const [key, label] of Object.entries(OFFICIAL_PLAN_LABELS)) {
        expect(i18n.getResource(language, 'translation', key), `${language}:${key}`).toBe(label)
      }
    }
  })
})

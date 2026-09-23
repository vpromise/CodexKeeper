import { renderToStaticMarkup } from 'react-dom/server';
import { expect, it, vi } from 'vitest';
import { LoginPage } from '../LoginPage';

vi.mock('react-i18next', async (importOriginal) => ({
  ...await importOriginal<typeof import('react-i18next')>(),
  useTranslation: () => ({ t: (key: string) => key }),
}));
vi.mock('@/components/ui/LanguageSwitcher', () => ({ LanguageSwitcher: () => null }));

it('shows the administrator password error', () => {
  const html = renderToStaticMarkup(<LoginPage onPasswordSubmit={vi.fn()} adminError="bad password" />);
  expect(html).toContain('bad password');
  expect(html).not.toContain('api_key');
});

it('exposes all theme options in a labelled control', () => {
  const html = renderToStaticMarkup(<LoginPage onPasswordSubmit={vi.fn()} />);
  expect(html).toContain('role="tablist" aria-label="usage_stats.theme_switch"');
  for (const theme of ['light', 'dark', 'auto']) {
    expect(html).toContain(`usage_stats.theme_${theme}`);
  }
});

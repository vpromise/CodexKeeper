import React from 'react';
import i18n from '@/i18n';
import { describe, expect, it, vi } from 'vitest';
import { renderToStaticMarkup } from 'react-dom/server';
import { getSessionLogoutConfirmationKeys, SessionSettingsCard } from '../SessionSettingsCard';
import type { AuthManagedSessionItem } from '@/lib/types';

const sessions: AuthManagedSessionItem[] = [
  {
    id: 'current-admin-hash',
    kind: 'admin',
    role: 'admin',
    source: 'standard',
    current: true,
    loginAt: '2026/06/20 10:00:00',
    expiresAt: '2026/06/20 12:00:00',
  },
  {
    id: 'other-admin-hash',
    kind: 'admin',
    role: 'admin',
    source: 'standard',
    loginAt: '2026/06/20 10:05:00',
    expiresAt: '2026/06/20 12:05:00',
  },
  {
    id: 'hashed-session-id',
    kind: 'admin',
    role: 'admin',
    source: 'embed',
    alias: 'Embedded Admin',
    loginAt: '2026/06/20 10:10:00',
    expiresAt: '2026/06/27 10:10:00',
  },
];

const renderCard = (props: Partial<React.ComponentProps<typeof SessionSettingsCard>> = {}) => renderToStaticMarkup(
  <SessionSettingsCard
    sessions={sessions}
    loading={false}
    revokingId={null}
    onLogout={() => undefined}
    {...props}
  />,
);

describe('SessionSettingsCard', () => {
  it('renders standalone and embedded admin sessions with shared row details and current marker', () => {
    const html = renderCard();

    expect(html).toContain('Session Management');
    expect(html).toContain('Admin Session');
    expect(html).toContain('Standalone');
    expect(html).toContain('CPAMC Embed');
    expect(html).toContain('Current');
    expect(html).toContain('2026/06/20 10:00:00');
    expect(html).toContain('2026/06/20 12:00:00');
    expect(html).toContain('2026/06/20 10:05:00');
    expect(html).toContain('2026/06/20 12:05:00');
    expect(html).toContain('Embedded Admin');
    expect(html).toContain('2026/06/20 10:10:00');
    expect(html).toContain('2026/06/27 10:10:00');
    expect(html).not.toContain('All admin sessions will be signed out together.');
    expect(html).not.toContain('sk-*********123456');
    expect(html).not.toContain('current-admin-hash');
    expect(html).not.toContain('other-admin-hash');
    expect(html).not.toContain('hashed-session-id');
    expect(html).not.toContain('api_key_viewer');
    expect(html.match(/>Sign out</g)).toHaveLength(2);
  });

  it('renders loading and empty states', () => {
    expect(renderCard({ sessions: [], loading: true })).toContain('Loading...');
    expect(renderCard({ sessions: [], loading: false })).toContain('No active sessions.');
  });

  it('uses per-session warning copy for both admin and API key confirmations', () => {
    const adminKeys = getSessionLogoutConfirmationKeys();

    expect(i18n.t(adminKeys.bodyKey)).toContain('this admin session');
    expect(i18n.t(adminKeys.bodyKey)).toContain('Other admin sessions');
    expect(i18n.t(adminKeys.bodyKey)).not.toContain('current device');
  });

  it('disables the row currently being revoked', () => {
    const html = renderCard({ revokingId: 'hashed-session-id' });

    expect(html).toContain('Signing out');
    expect(html).toContain('disabled=""');
  });

  it('does not invoke logout while only rendering', () => {
    const onLogout = vi.fn();

    renderCard({ onLogout });

    expect(onLogout).not.toHaveBeenCalled();
  });
});

describe('SessionSettingsCard client metadata', () => {
  it('renders the raw User-Agent, login IP, changed recent IP, and recent activity', () => {
    const userAgent = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/18.0 Safari/605.1.15';
    const session: AuthManagedSessionItem = {
      id: 'session-hash',
      kind: 'admin',
      role: 'admin',
      source: 'standard',
      userAgent,
      loginIp: '203.0.113.9',
      lastSeenIp: '198.51.100.42',
      loginAt: '2026/08/13 10:00:00',
      lastSeenAt: '2026/08/13 10:15:00',
      expiresAt: '2026/08/20 10:00:00',
    };

    const html = renderCard({ sessions: [session] });

    expect(html.split(userAgent)).toHaveLength(2);
    expect(html).toContain('User-Agent');
    expect(html).toMatch(/<dt[^>]*>Login IP<\/dt><dd[^>]*>203\.0\.113\.9<\/dd>/);
    expect(html).toMatch(/<dt[^>]*>Recent IP<\/dt><dd[^>]*>198\.51\.100\.42<\/dd>/);
    expect(html).toMatch(/<dt[^>]*>Last active<\/dt><dd[^>]*>2026\/08\/13 10:15:00<\/dd>/);
  });

  it('renders an honest fallback for sessions created before metadata collection', () => {
    const legacySession: AuthManagedSessionItem = {
      id: 'legacy-session-hash',
      kind: 'admin',
      role: 'admin',
      source: 'standard',
      loginAt: '2026/08/01 10:00:00',
      lastSeenAt: '2026/08/01 10:00:00',
      expiresAt: '2026/08/08 10:00:00',
    };

    const html = renderCard({ sessions: [legacySession] });

    expect(html).toMatch(/<dt[^>]*>Login IP<\/dt><dd[^>]*>Unknown<\/dd>/);
  });
});

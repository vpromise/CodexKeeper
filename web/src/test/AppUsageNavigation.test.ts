import { describe, expect, it } from 'vitest';
import { getAdminTargetPath } from '../App';
describe('administrator navigation', () => {
  it.each(['/overview', '/analysis', '/realtime', '/auth-files', '/ai-provider', '/request-events', '/settings'])('keeps %s', path => {
    expect(getAdminTargetPath(path)).toBe(path);
  });
  it.each(['/', '/ranking', '/key-overview', '/key-analysis', '/key-ranking', '/auth-files/', '//example.com/auth-files'])('normalizes unavailable path %s', path => {
    expect(getAdminTargetPath(path)).toBe('/');
  });
});

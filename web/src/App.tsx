import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import './index.css';
import './App.css';
import './embed/cpamcEmbed.css';
import { ApiError, appPath, clearEmbedSessionToken, getSession, login } from './lib/api';
import { AppFooter } from './components/AppFooter';
import { LoginPage } from './pages/LoginPage';
import { UsagePage } from './pages/UsagePage';
import { cpamcEmbedSearch, isCPAMCEmbed, notifyCPAMCEmbedReady } from './embed/cpamcEmbed';
import { getUsageTabPath, resolveUsageTabFromPath, stripAppBasePath } from './lib/usageNavigation';
import { useUsageStatsStore } from './stores/useUsageStatsStore';

type AuthState = 'checking' | 'authenticated' | 'unauthenticated';

export const getAdminTargetPath = (currentPath: string): string => {
  const tab = resolveUsageTabFromPath(currentPath);
  return tab ? getUsageTabPath(tab) : '/';
};

function App() {
  const { t } = useTranslation();
  const [authState, setAuthState] = useState<AuthState>('checking');
  const [loginError, setLoginError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const clearUsageStats = useUsageStatsStore((state) => state.clearUsageStats);
  const isEmbeddedInCPAMC = isCPAMCEmbed();

  const clearSession = useCallback(() => {
    clearEmbedSessionToken();
    clearUsageStats();
    setAuthState('unauthenticated');
  }, [clearUsageStats]);

  const loadSession = useCallback(async () => {
    const session = await getSession();
    if (!session.authenticated || (session.role && session.role !== 'admin')) {
      clearSession();
      return false;
    }
    setAuthState('authenticated');
    const currentPath = stripAppBasePath(window.location.pathname, window.__APP_BASE_PATH__);
    const targetPath = getAdminTargetPath(currentPath ?? '/');
    if (currentPath !== targetPath) {
      window.history.replaceState(null, '', appPath(targetPath) + cpamcEmbedSearch());
    }
    return true;
  }, [clearSession]);

  useEffect(() => {
    void loadSession().catch(clearSession);
  }, [clearSession, loadSession]);

  useEffect(() => {
    notifyCPAMCEmbedReady();
  }, []);

  const handlePasswordLogin = useCallback(async (password: string) => {
    setSubmitting(true);
    setLoginError('');
    try {
      await login(password);
      if (!await loadSession()) setLoginError(t('auth.login_failed'));
    } catch (error) {
      setLoginError(t(error instanceof ApiError && error.status === 401
        ? 'auth.invalid_password'
        : error instanceof ApiError && error.status === 429
          ? 'auth.login_rate_limited'
          : 'auth.login_failed'));
      clearSession();
    } finally {
      setSubmitting(false);
    }
  }, [clearSession, loadSession, t]);

  let page: ReactNode;
  if (authState === 'checking') {
    page = <div className="app-checking" aria-busy="true" />;
  } else if (authState === 'unauthenticated') {
    page = <LoginPage loading={submitting} adminError={loginError} onPasswordSubmit={handlePasswordLogin} />;
  } else {
    page = <UsagePage onAuthRequired={clearSession} />;
  }

  return (
    <div className="app-frame" data-embed={isEmbeddedInCPAMC ? 'cpamc' : undefined}>
      <main className="app-main">{page}</main>
      <AppFooter loadVersion={authState === 'authenticated'} />
    </div>
  );
}

export default App;

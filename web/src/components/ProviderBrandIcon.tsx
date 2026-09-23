import claudeIcon from '@/assets/icons/claude.svg'
import codexIcon from '@/assets/icons/codex.svg'
import styles from './ProviderBrandIcon.module.scss'

export const PROVIDER_BRAND_ICON_KEYS = [
  'claude',
  'codex',
] as const

export type ProviderBrandIconKey = typeof PROVIDER_BRAND_ICON_KEYS[number]

export interface ProviderBrandIconProps {
  providerType: string | null | undefined
  size: number | string
  ariaLabel?: string
  className?: string
}

// Only native Codex and Claude identities have provider branding.
const providerBrandIconKeyByType: Readonly<Record<string, ProviderBrandIconKey>> = {
  claude: 'claude',
  codex: 'codex',
}

// 品牌资源取自 Lobe Icons；统一复用其 Avatar 圆形容器、背景色与缩放比例。
const providerBrandIconUrlByKey: Readonly<Record<ProviderBrandIconKey, string>> = {
  claude: claudeIcon,
  codex: codexIcon,
}

export function providerBrandIconKey(providerType: string | null | undefined): ProviderBrandIconKey | undefined {
  const normalized = providerType?.trim().toLowerCase()
  return normalized ? providerBrandIconKeyByType[normalized] : undefined
}

export function ProviderBrandIcon({ providerType, size, ariaLabel, className }: ProviderBrandIconProps) {
  const iconKey = providerBrandIconKey(providerType)
  if (!iconKey) {
    return null
  }

  const rootClassName = `${styles.providerBrandIcon} ${styles.providerBrandIconAvatar} ${className ?? ''}`.trim()
  const accessibleLabel = ariaLabel?.trim() || undefined

  return (
    <span
      className={rootClassName}
      data-provider-brand-icon={iconKey}
      data-provider-brand-icon-tone="avatar"
      style={{ width: size, height: size }}
      role={accessibleLabel ? 'img' : undefined}
      aria-label={accessibleLabel}
      aria-hidden={accessibleLabel ? undefined : true}
    >
      <img className={styles.providerBrandIconAsset} src={providerBrandIconUrlByKey[iconKey]} alt="" />
    </span>
  )
}

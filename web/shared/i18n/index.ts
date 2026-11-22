import { createI18n } from 'vue-i18n'
import zhCN from './messages/zh-CN'
import en from './messages/en'

export const SUPPORTED_LOCALES = ['zh-CN', 'en'] as const
export type AppLocale = (typeof SUPPORTED_LOCALES)[number]
export const LOCALE_LABELS: Record<AppLocale, string> = {
  'zh-CN': '简体中文',
  en: 'English',
}
export const LOCALE_STORAGE_KEY = 'frp-locale'

export function getInitialLocale(): AppLocale {
  if (typeof localStorage !== 'undefined') {
    const saved = localStorage.getItem(LOCALE_STORAGE_KEY)
    if (saved && (SUPPORTED_LOCALES as readonly string[]).includes(saved)) {
      return saved as AppLocale
    }
  }
  const nav = typeof navigator !== 'undefined' ? navigator.language : ''
  if (nav && nav.toLowerCase().startsWith('zh')) return 'zh-CN'
  return 'zh-CN'
}

export function createI18nForApp(_app: 'frps' | 'frpc') {
  return createI18n({
    legacy: false,
    globalInjection: true,
    locale: getInitialLocale(),
    fallbackLocale: 'zh-CN',
    messages: {
      'zh-CN': zhCN,
      en,
    },
  })
}

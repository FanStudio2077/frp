import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import en from './en'

const saved = typeof localStorage !== 'undefined' ? localStorage.getItem('console-locale') : null

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: saved === 'en' ? 'en' : 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: { 'zh-CN': zhCN, en },
})

export default i18n

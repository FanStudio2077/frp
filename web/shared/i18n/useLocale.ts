import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useStorage } from '@vueuse/core'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import {
  LOCALE_LABELS,
  LOCALE_STORAGE_KEY,
  SUPPORTED_LOCALES,
  type AppLocale,
} from './index'

export function useLocale() {
  const { locale } = useI18n()
  const stored = useStorage<AppLocale>(LOCALE_STORAGE_KEY, locale.value as AppLocale)

  const applyToDocument = (value: string) => {
    if (typeof document !== 'undefined') document.documentElement.lang = value
  }

  watch(
    stored,
    (value) => {
      locale.value = value
      applyToDocument(value)
    },
    { immediate: true },
  )

  const setLocale = (value: AppLocale) => {
    stored.value = value
    locale.value = value
    applyToDocument(value)
  }

  const elLocale = computed(() => (locale.value === 'en' ? en : zhCn))
  const availableLocales = computed(() =>
    SUPPORTED_LOCALES.map((value) => ({ value, label: LOCALE_LABELS[value] })),
  )

  return { locale, elLocale, availableLocales, setLocale }
}

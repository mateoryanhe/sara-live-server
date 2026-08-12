import {createI18n} from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import es from 'element-plus/es/locale/lang/es'
import pt from 'element-plus/es/locale/lang/pt'
import type {Language} from 'element-plus/es/locale'
import zhCN from './locales/zh-CN'
import enUS from './locales/en'
import esES from './locales/es'
import ptBR from './locales/pt'
import hiIN from './locales/hi'
import idID from './locales/id'
import {CMS_LOCALE_STORAGE_KEY, SUPPORTED_LOCALES, type CmsLocale} from './locales/types'

export {CMS_LOCALE_STORAGE_KEY, SUPPORTED_LOCALES, LOCALE_LABELS} from './locales/types'
export type {CmsLocale} from './locales/types'

const elementPlusLocales: Record<CmsLocale, Language> = {
  'zh-CN': zhCn,
  en,
  es,
  pt,
  hi: en,
  id: en,
}

export function getElementPlusLocale(locale: CmsLocale): Language {
  return elementPlusLocales[locale] || en
}

function normalizeLocale(raw: string): CmsLocale | null {
  const low = raw.toLowerCase().trim()
  if (low === 'zh-cn' || low === 'zh' || low.startsWith('zh-')) {
    return 'zh-CN'
  }
  if (low === 'en' || low.startsWith('en-')) {
    return 'en'
  }
  if (low === 'es' || low.startsWith('es-')) {
    return 'es'
  }
  if (low === 'pt' || low.startsWith('pt-')) {
    return 'pt'
  }
  if (low === 'hi' || low.startsWith('hi-')) {
    return 'hi'
  }
  if (low === 'id' || low.startsWith('id-')) {
    return 'id'
  }
  return null
}

export function detectBrowserLocale(): CmsLocale {
  const stored = localStorage.getItem(CMS_LOCALE_STORAGE_KEY)
  if (stored) {
    const normalized = normalizeLocale(stored)
    if (normalized) {
      return normalized
    }
  }

  const candidates = navigator.languages?.length
    ? navigator.languages
    : [navigator.language]

  for (const lang of candidates) {
    const normalized = normalizeLocale(lang)
    if (normalized) {
      return normalized
    }
  }

  return 'zh-CN'
}

export function persistLocale(locale: CmsLocale) {
  localStorage.setItem(CMS_LOCALE_STORAGE_KEY, locale)
}

export const i18n = createI18n({
  legacy: false,
  locale: detectBrowserLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    en: enUS,
    es: esES,
    pt: ptBR,
    hi: hiIN,
    id: idID,
  },
})

export function setAppLocale(locale: CmsLocale) {
  i18n.global.locale.value = locale
  persistLocale(locale)
  document.documentElement.lang = locale
}

document.documentElement.lang = i18n.global.locale.value

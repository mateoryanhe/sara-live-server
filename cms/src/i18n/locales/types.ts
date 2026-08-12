export const CMS_LOCALE_STORAGE_KEY = 'cms_locale'

export const SUPPORTED_LOCALES = ['zh-CN', 'en', 'es', 'pt', 'hi', 'id'] as const
export type CmsLocale = (typeof SUPPORTED_LOCALES)[number]

export const LOCALE_LABELS: Record<CmsLocale, string> = {
  'zh-CN': '中文',
  en: 'EN',
  es: 'ES',
  pt: 'PT',
  hi: 'HI',
  id: 'ID',
}

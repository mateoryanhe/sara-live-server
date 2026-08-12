export type PageLocale = 'zh-CN' | 'en' | 'es' | 'pt' | 'hi' | 'id'

export function definePageMessages<T extends Record<string, string>>(
  zh: T,
  en: Record<keyof T, string>,
  es: Record<keyof T, string>,
  pt: Record<keyof T, string>,
  hi: Record<keyof T, string>,
  id: Record<keyof T, string>,
) {
  return {
    'zh-CN': zh,
    en,
    es,
    pt,
    hi,
    id,
  } as const
}

export function definePageMessagesFromEn<T extends Record<string, string>>(
  zh: T,
  en: Record<keyof T, string>,
) {
  return definePageMessages(zh, en, en as Record<keyof T, string>, en as Record<keyof T, string>, en as Record<keyof T, string>, en as Record<keyof T, string>)
}

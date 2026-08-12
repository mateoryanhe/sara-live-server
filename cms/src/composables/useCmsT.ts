import {useI18n} from 'vue-i18n'

/** vue-i18n translate with named params */
export function useCmsT() {
  const {t, te, locale} = useI18n()
  return {t, te, locale}
}

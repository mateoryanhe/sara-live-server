import {computed} from 'vue'
import {useRoute} from 'vue-router'
import {
  coinMerchantPaymentCountryCfgApi,
  paymentCountryCfgApi,
  userWithdrawPayoutCountryCfgApi,
} from '@/api/modules/payment-country-cfg'

export type PaymentCountryCfgScope = 'normal' | 'coinMerchant' | 'userWithdrawPayout'

export function usePaymentCountryCfgScope() {
  const route = useRoute()
  const cfgScope = computed<PaymentCountryCfgScope>(() => {
    const raw = String(route.meta.paymentCountryCfgScope || 'normal')
    if (raw === 'coinMerchant' || raw === 'userWithdrawPayout') {
      return raw
    }
    return 'normal'
  })
  const isCoinMerchant = computed(() => cfgScope.value === 'coinMerchant')
  const isUserWithdrawPayout = computed(() => cfgScope.value === 'userWithdrawPayout')
  const usesCombinedPaymentFields = computed(() => isCoinMerchant.value || isUserWithdrawPayout.value)

  const configApi = computed(() => {
    if (isUserWithdrawPayout.value) {
      return userWithdrawPayoutCountryCfgApi
    }
    if (isCoinMerchant.value) {
      return coinMerchantPaymentCountryCfgApi
    }
    return paymentCountryCfgApi
  })

  const rootRouteName = computed(() => {
    if (isUserWithdrawPayout.value) {
      return 'UserWithdrawPayoutCountryCfgManagement'
    }
    if (isCoinMerchant.value) {
      return 'CoinMerchantPaymentCountryCfgManagement'
    }
    return 'PaymentCountryCfgManagement'
  })

  const continentRouteName = computed(() => {
    if (isUserWithdrawPayout.value) {
      return 'UserWithdrawPayoutCountryCfgContinent'
    }
    if (isCoinMerchant.value) {
      return 'CoinMerchantPaymentCountryCfgContinent'
    }
    return 'PaymentCountryCfgContinent'
  })

  const editRouteName = computed(() => {
    if (isUserWithdrawPayout.value) {
      return 'UserWithdrawPayoutCountryCfgEdit'
    }
    if (isCoinMerchant.value) {
      return 'CoinMerchantPaymentCountryCfgEdit'
    }
    return 'PaymentCountryCfgEdit'
  })

  const permissionPage = computed(() => {
    if (isUserWithdrawPayout.value) {
      return 'UserWithdrawPayoutCountryCfgManagement'
    }
    if (isCoinMerchant.value) {
      return 'CoinMerchantPaymentCountryCfgManagement'
    }
    return 'PaymentCountryCfgManagement'
  })

  return {
    cfgScope,
    isCoinMerchant,
    isUserWithdrawPayout,
    usesCombinedPaymentFields,
    configApi,
    rootRouteName,
    continentRouteName,
    editRouteName,
    permissionPage,
  }
}

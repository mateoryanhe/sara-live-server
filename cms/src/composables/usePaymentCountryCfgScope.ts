import {computed} from 'vue'
import {useRoute} from 'vue-router'
import {
  coinMerchantPaymentCountryCfgApi,
  paymentCountryCfgApi,
  userWithdrawPayoutCountryCfgApi,
} from '@/api/modules/payment-country-cfg'

export type PaymentCountryCfgScope = 'normal' | 'coinMerchant' | 'userWithdrawPayout'

/** 用户提现代付：部分国家固定支付类型（与 App/HaiPay 约定一致） */
export const USER_WITHDRAW_PAYOUT_FIXED_PAY_TYPE_BY_COUNTRY: Readonly<Record<string, string>> = {
  ID: 'EWALLET',
  MY: 'EWALLET',
}

export function userWithdrawPayoutFixedPayType(countryCode: string): string | null {
  const code = String(countryCode || '').trim().toUpperCase()
  return USER_WITHDRAW_PAYOUT_FIXED_PAY_TYPE_BY_COUNTRY[code] ?? null
}

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
  const isUserWithdrawPayout = computed(() => {
    if (cfgScope.value === 'userWithdrawPayout') {
      return true
    }
    const name = String(route.name ?? '')
    return name === 'UserWithdrawPayoutCountryCfgManagement'
        || name === 'UserWithdrawPayoutCountryCfgContinent'
        || name === 'UserWithdrawPayoutCountryCfgEdit'
        || route.path.includes('/user-withdraw-payout-country-cfg')
  })
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

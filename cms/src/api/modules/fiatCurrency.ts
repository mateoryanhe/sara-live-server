import {request} from '../request'

export interface HaiPayRegionItem {
    currencyCode: string
    fiatCurrencyCode: string
    fiatCurrencyCodes: string[]
    name: string
    nameEn: string
    nameZh: string
    symbol: string
    icon: string
    currencyType: number
    sort: number
    paymentMethods: Array<{
        payType: string
        inBankCode: string
        minAmount: string
        maxAmount: string
        description: string
    }>
}

export const fiatCurrencyApi = {
    haiPayRegionList: () => {
        return request.post<{ list: HaiPayRegionItem[] }>('/fiatCurrency/haiPayRegionList', {})
    },
}

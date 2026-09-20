import {request} from '../request'

export interface PaymentCountryMethodSelection {
    currencyCode: string
    payType: string
    inBankCode: string
}

export interface PaymentCountryCfgItem {
    id: string
    countryCode: string
    countryNameEn: string
    countryNameZh: string
    continent: string
    supportedCurrencies: string[]
    currencyCode: string
    enabled: boolean
    paymentMethods: PaymentCountryPaymentMethod[]
    selectedPaymentMethods: PaymentCountryMethodSelection[]
}

export interface PaymentCountryPaymentMethod {
    currencyCode: string
    payType: string
    inBankCode: string
    minAmount: string
    maxAmount: string
    description: string
    available: boolean
}

export interface PaymentCountryCfgGroup {
    continent: string
    countries: PaymentCountryCfgItem[]
}

export interface PaymentCountryCfgResponse {
    continents: PaymentCountryCfgGroup[]
    paymentTypes: string[]
}

/** 未持久化的目录数据没有配置 ID，App 可见性必须按关闭处理。 */
export function normalizePaymentCountryCfgGroups(groups?: PaymentCountryCfgGroup[]): PaymentCountryCfgGroup[] {
    return (groups || []).map((group) => ({
        ...group,
        countries: (group.countries || []).map((row) => ({
            ...row,
            paymentMethods: row.paymentMethods || [],
            selectedPaymentMethods: row.selectedPaymentMethods || [],
            enabled: Boolean(row.id && row.id !== '0') && row.enabled,
        })),
    }))
}

const prefix = '/paymentCountryCfg'

export const paymentCountryCfgApi = {
    getConfig: () => request.post<PaymentCountryCfgResponse>(`${prefix}/getCollectionCountryCfg`, {}),

    saveBasic: (data: Pick<PaymentCountryCfgItem, 'countryCode' | 'currencyCode' | 'enabled'>) => {
        return request.post<{ success: boolean; id: string }>(`${prefix}/saveCollectionCountryCfg`, {
            ...data,
            saveSection: 'basic',
        })
    },
}

export const coinMerchantPaymentCountryCfgApi = {
    getConfig: () => request.post<PaymentCountryCfgResponse>(
        '/coinMerchantPaymentCountryCfg/getCollectionCountryCfg', {},
    ),

    saveConfig: (data: Pick<PaymentCountryCfgItem, 'countryCode' | 'currencyCode' | 'enabled'> & {
        payType: string
        inBankCode: string
    }) => request.post<{ success: boolean; id: string }>(
        '/coinMerchantPaymentCountryCfg/saveCollectionCountryCfg', data,
    ),
}

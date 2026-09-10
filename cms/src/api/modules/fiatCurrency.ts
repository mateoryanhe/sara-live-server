import {request} from '../request'
import type {
    FiatCurrency,
    FiatCurrencyQuery,
    PageResponse,
} from '@/types/api'

export const fiatCurrencyApi = {
    getList: (params: FiatCurrencyQuery) => {
        return request.post<PageResponse<FiatCurrency>>('/fiatCurrency/fiatCurrencyList', params)
    },

    create: (data: {
        currencyCode: string
        name: string
        symbol: string
        icon?: string
        currencyType: number
        sort?: number
        status?: number
    }) => {
        return request.post<{ id: string }>('/fiatCurrency/createFiatCurrency', data)
    },

    update: (data: {
        id: string | number
        currencyCode: string
        name: string
        symbol: string
        icon?: string
        currencyType: number
        sort?: number
        status?: number
    }) => {
        return request.post<{ success: boolean }>('/fiatCurrency/updateFiatCurrency', data)
    },

    remove: (id: string | number) => {
        return request.post<{ success: boolean }>('/fiatCurrency/deleteFiatCurrency', {id})
    },

    reloadCfgCache: () => {
        return request.post<{ success: boolean }>('/fiatCurrency/reloadFiatCurrencyCache', {})
    },
}

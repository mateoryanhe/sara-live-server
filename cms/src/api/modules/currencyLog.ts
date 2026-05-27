import {request} from '../request'
import type {CurrencyLogQuery, PageResponse, CurrencyLogItem} from '@/types/api'

export const currencyLogApi = {
    getCurrencyLogList: (params: CurrencyLogQuery) => {
        return request.post<PageResponse<CurrencyLogItem>>('/currencyLog/cmsCurrencyLogList', params)
    },
}

export default currencyLogApi

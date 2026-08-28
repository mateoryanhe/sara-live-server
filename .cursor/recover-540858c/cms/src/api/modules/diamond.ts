import {request} from '../request'
import type {CMSAdjustCurrencyPayload} from './gold'

export const diamondApi = {
    add: (data: CMSAdjustCurrencyPayload) => {
        return request.post<{ diamond: number }>('/diamond/add', data)
    },

    sub: (data: CMSAdjustCurrencyPayload) => {
        return request.post<{ diamond: number }>('/diamond/sub', data)
    },
}

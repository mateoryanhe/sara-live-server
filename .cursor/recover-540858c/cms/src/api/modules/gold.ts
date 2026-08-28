import {request} from '../request'

export interface CMSAdjustCurrencyPayload {
    userId: string | number
    amount: number
    reason: number
}

export const goldApi = {
    add: (data: CMSAdjustCurrencyPayload) => {
        return request.post<{ gold: number }>('/gold/add', data)
    },

    sub: (data: CMSAdjustCurrencyPayload) => {
        return request.post<{ gold: number }>('/gold/sub', data)
    },
}

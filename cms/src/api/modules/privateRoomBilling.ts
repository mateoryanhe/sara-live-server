import {request} from '../request'
import type {PageResponse, PrivateRoomBilling, PrivateRoomBillingQuery} from '@/types/api'

export const privateRoomBillingApi = {
    getBillingList: (params: PrivateRoomBillingQuery) => {
        return request.post<PageResponse<PrivateRoomBilling>>('/privateRoomBilling/billingList', params)
    },

    createBilling: (data: {
        pricePerMinute: number
        sort: number
    }) => {
        return request.post<{ id: string }>('/privateRoomBilling/createBilling', data)
    },

    updateBilling: (data: {
        id: string | number
        pricePerMinute: number
        sort: number
    }) => {
        return request.post<boolean>('/privateRoomBilling/updateBilling', data)
    },

    deleteBilling: (id: string | number) => {
        return request.post<boolean>('/privateRoomBilling/deleteBilling', {id})
    },

    onShelfBilling: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/privateRoomBilling/onShelfBilling', {id})
    },

    offShelfBilling: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/privateRoomBilling/offShelfBilling', {id})
    },
}

import {request} from '../request'
import type {PageResponse, RechargeOrder, RechargeOrderQuery} from '@/types/api'

export const rechargeOrderApi = {
    getRechargeOrderList: (params: RechargeOrderQuery) => {
        return request.post<PageResponse<RechargeOrder>>('/rechargeOrder/rechargeOrderList', params)
    },

    manualRecharge: (orderId: string) => {
        return request.post<{ orderId: string; gold: number; after: number; success: boolean }>(
            '/rechargeOrder/manualRecharge',
            {orderId},
        )
    },

    manualCreateOrder: (params: { userId: string; amount?: number; cfgId?: string }) => {
        return request.post<{
            orderId: string
            price: number
            gold: number
            currency: string
            status: number
            success: boolean
        }>('/rechargeOrder/manualCreateOrder', params)
    },

    createChannelRechargeOrderTest: (params: {
        userId: string
        cfgId: number
        currencyCode: string
        packageName?: string
    }) => {
        return request.post<{
            orderId: string
            payUrl: string
            price: number
            payAmount: number
            currency: string
            status: number
        }>('/rechargeOrder/createChannelRechargeOrderTest', params)
    },
}

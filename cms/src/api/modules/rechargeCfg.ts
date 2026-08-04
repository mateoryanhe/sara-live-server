import {request} from '../request'
import type {PageResponse, RechargeCfg, RechargeCfgQuery} from '@/types/api'

export const rechargeCfgApi = {
    getRechargeCfgList: (params: RechargeCfgQuery) => {
        return request.post<PageResponse<RechargeCfg>>('/rechargeCfg/rechargeCfgList', params)
    },

    createRechargeCfg: (data: {
        name: string
        packageName: string
        cfgType: number
        icon: string
        gold: number
        extraGold: number
        price: number
        productId: string
        sort: number
        description: string
    }) => {
        return request.post<{ id: string }>('/rechargeCfg/createRechargeCfg', data)
    },

    updateRechargeCfg: (data: {
        id: string | number
        name: string
        packageName: string
        cfgType: number
        icon: string
        gold: number
        extraGold: number
        price: number
        productId: string
        sort: number
        description: string
    }) => {
        return request.post<boolean>('/rechargeCfg/updateRechargeCfg', data)
    },

    deleteRechargeCfg: (id: string | number) => {
        return request.post<boolean>('/rechargeCfg/deleteRechargeCfg', {id})
    },

    onShelfRechargeCfg: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/rechargeCfg/onShelfRechargeCfg', {id})
    },

    offShelfRechargeCfg: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/rechargeCfg/offShelfRechargeCfg', {id})
    },
}

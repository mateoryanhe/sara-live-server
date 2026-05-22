import {request} from '../request'
import type {PageResponse, VipCfg, VipCfgQuery} from '@/types/api'

export const vipCfgApi = {
    getVipCfgList: (params: VipCfgQuery) => {
        return request.post<PageResponse<VipCfg>>('/vipCfg/vipCfgList', params)
    },

    createVipCfg: (data: {
        level: number
        levelName: string
        status: number
        upgradeRechargeLimit: number
        minWithdrawAmount: number
        maxWithdrawAmount: number
        fee: number
    }) => {
        return request.post<{ id: string }>('/vipCfg/createVipCfg', data)
    },

    updateVipCfg: (data: {
        id: string | number
        level: number
        levelName: string
        status: number
        upgradeRechargeLimit: number
        minWithdrawAmount: number
        maxWithdrawAmount: number
        fee: number
    }) => {
        return request.post<boolean>('/vipCfg/updateVipCfg', data)
    },

    deleteVipCfg: (id: string | number) => {
        return request.post<boolean>('/vipCfg/deleteVipCfg', {id})
    },
}

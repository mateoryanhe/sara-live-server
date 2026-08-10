import {request} from '../request'
import type {PageResponse, VipCfg, VipCfgQuery} from '@/types/api'

type VipCfgSavePayload = {
    level: number
    levelName: string
    levelIcon?: string
    withdrawSwitch: number
    animationSwitch: number
    commentEffectSwitch: number
    customerServiceSwitch: number
    upgradeRechargeLimit: number
    minWithdrawAmount: number
    maxWithdrawAmount: number
    fee: number
    animation?: string
    animationIcon?: string
    animationDescEn?: string
    animationDescEs?: string
    animationDescPt?: string
    animationDescHi?: string
    commentEffect?: string
    commentEffectIcon?: string
    commentEffectDescEn?: string
    commentEffectDescEs?: string
    commentEffectDescPt?: string
    commentEffectDescHi?: string
    withdrawIcon?: string
    withdrawNoticeEn?: string
    withdrawNoticeEs?: string
    withdrawNoticePt?: string
    withdrawNoticeHi?: string
    customerServiceIcon?: string
    customerServiceDescEn?: string
    customerServiceDescEs?: string
    customerServiceDescPt?: string
    customerServiceDescHi?: string
}

export const vipCfgApi = {
    getVipCfgList: (params: VipCfgQuery) => {
        return request.post<PageResponse<VipCfg>>('/vipCfg/vipCfgList', params)
    },

    createVipCfg: (data: VipCfgSavePayload) => {
        return request.post<{ id: string }>('/vipCfg/createVipCfg', data)
    },

    updateVipCfg: (data: VipCfgSavePayload & { id: string | number }) => {
        return request.post<boolean>('/vipCfg/updateVipCfg', data)
    },

    deleteVipCfg: (id: string | number) => {
        return request.post<boolean>('/vipCfg/deleteVipCfg', {id})
    },
}

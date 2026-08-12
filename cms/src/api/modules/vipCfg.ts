import {request} from '../request'
import type {PageResponse, VipCfg, VipCfgQuery} from '@/types/api'

type VipCfgSavePayload = {
    level: number
    levelName: string
    levelIcon?: string
    animationSwitch: number
    commentEffectSwitch: number
    customerServiceSwitch: number
    upgradeRechargeLimit: number
    animation?: string
    animationIcon?: string
    animationTitleEn?: string
    animationTitleEs?: string
    animationTitlePt?: string
    animationTitleHi?: string
    animationTitleId?: string
    animationDescEn?: string
    animationDescEs?: string
    animationDescPt?: string
    animationDescHi?: string
    animationDescId?: string
    commentEffect?: string
    commentEffectIcon?: string
    commentEffectTitleEn?: string
    commentEffectTitleEs?: string
    commentEffectTitlePt?: string
    commentEffectTitleHi?: string
    commentEffectTitleId?: string
    commentEffectDescEn?: string
    commentEffectDescEs?: string
    commentEffectDescPt?: string
    commentEffectDescHi?: string
    commentEffectDescId?: string
    customerServiceIcon?: string
    customerServiceTitleEn?: string
    customerServiceTitleEs?: string
    customerServiceTitlePt?: string
    customerServiceTitleHi?: string
    customerServiceTitleId?: string
    customerServiceDescEn?: string
    customerServiceDescEs?: string
    customerServiceDescPt?: string
    customerServiceDescHi?: string
    customerServiceDescId?: string
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

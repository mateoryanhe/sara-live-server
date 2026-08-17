import {request} from '../request'
import type {
    AnchorIncomeSettlementLogItem,
    Guild,
    GuildQuery,
    ImportGuildAnchorsReq,
    ImportGuildAnchorsRes,
    MyGuildAnchorIncomeSettlementLogQuery,
    MyGuildProfileListRes,
    PageResponse,
    UpdateMyGuildProfileReq,
} from '@/types/api'

// 直播工会管理API
export const guildApi = {
    // 获取工会列表
    getGuildList: (params: GuildQuery) => {
        return request.post<PageResponse<Guild>>('/guild/guildList', params)
    },

    // 创建工会
    createGuild: (data: {
        name: string
        leaderId: number
        description: string
    }) => {
        return request.post<{ id: string }>('/guild/createGuild', data)
    },

    // 更新工会
    updateGuild: (data: {
        id: string
        name: string
        leaderId: number
        description: string
    }) => {
        return request.post<boolean>('/guild/updateGuild', data)
    },

    // 下架工会
    deleteGuild: (id: string) => {
        return request.post<boolean>('/guild/deleteGuild', {id})
    },

    // 工会垃圾库：已下架列表
    getOffShelfGuildList: (params: GuildQuery) => {
        return request.post<PageResponse<Guild>>('/guild/offShelfGuildList', params)
    },

    // 上架工会
    onShelfGuild: (id: string) => {
        return request.post<{ success: boolean }>('/guild/onShelfGuild', {id})
    },

    getMyGuildProfile: () => {
        return request.post<MyGuildProfileListRes>('/guild/getMyGuildProfile', {})
    },

    getMyGuildAnchorIncomeSettlementLogList: (params: MyGuildAnchorIncomeSettlementLogQuery) => {
        return request.post<PageResponse<AnchorIncomeSettlementLogItem>>('/guild/cmsMyGuildAnchorIncomeSettlementLogList', params)
    },

    updateMyGuildProfile: (data: UpdateMyGuildProfileReq) => {
        return request.post<{ success: boolean }>('/guild/updateMyGuildProfile', data)
    },

    importGuildAnchors: (data: ImportGuildAnchorsReq) => {
        return request.post<ImportGuildAnchorsRes>('/guild/importGuildAnchors', data)
    },
}

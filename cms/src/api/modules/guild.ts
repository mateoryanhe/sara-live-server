import {request} from '../request'
import type {
    Guild,
    GuildQuery,
    ImportGuildAnchorsReq,
    ImportGuildAnchorsRes,
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
        timezone: number
    }) => {
        return request.post<{ id: string }>('/guild/createGuild', data)
    },

    // 更新工会
    updateGuild: (data: {
        id: string
        name: string
        leaderId: number
        description: string
        timezone: number
    }) => {
        return request.post<boolean>('/guild/updateGuild', data)
    },

    // 删除工会
    deleteGuild: (id: string) => {
        return request.post<boolean>('/guild/deleteGuild', {id})
    },

    getMyGuildProfile: () => {
        return request.post<MyGuildProfileListRes>('/guild/getMyGuildProfile', {})
    },

    updateMyGuildProfile: (data: UpdateMyGuildProfileReq) => {
        return request.post<{ success: boolean }>('/guild/updateMyGuildProfile', data)
    },

    importGuildAnchors: (data: ImportGuildAnchorsReq) => {
        return request.post<ImportGuildAnchorsRes>('/guild/importGuildAnchors', data)
    },

    // 批量更新工会时区
    batchUpdateGuildTimezone: (data: { guildIds: string[], timezone: number }) => {
        return request.post<{ success: boolean }>('/guild/batchUpdateGuildTimezone', data)
    },
}

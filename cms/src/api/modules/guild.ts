import {request} from '../request'
import type {Guild, GuildQuery, PageResponse} from '@/types/api'

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
        contact: string
        description: string
        status: number
    }) => {
        return request.post<{ id: string }>('/guild/createGuild', data)
    },

    // 更新工会
    updateGuild: (data: {
        id: string
        name: string
        leaderId: number
        contact: string
        description: string
        status: number
    }) => {
        return request.post<boolean>('/guild/updateGuild', data)
    },

    // 删除工会
    deleteGuild: (id: string) => {
        return request.post<boolean>('/guild/deleteGuild', {id})
    },
}

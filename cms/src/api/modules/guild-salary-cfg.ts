import {request} from '../request'
import type {GuildSalaryCfg, GuildSalaryCfgQuery, PageResponse} from '@/types/api'

export const guildSalaryCfgApi = {
  getList: (params: GuildSalaryCfgQuery) => {
    return request.post<PageResponse<GuildSalaryCfg>>('/guildSalaryCfg/guildSalaryCfgList', params)
  },

  create: (data: {
    weeklyWorkDays: number
    dailyLiveDurationMinutes: number
    salaryAmount: number
    sort: number
  }) => {
    return request.post<{ id: string }>('/guildSalaryCfg/createGuildSalaryCfg', data)
  },

  update: (data: {
    id: string | number
    weeklyWorkDays: number
    dailyLiveDurationMinutes: number
    salaryAmount: number
    sort: number
  }) => {
    return request.post<{ success: boolean }>('/guildSalaryCfg/updateGuildSalaryCfg', data)
  },

  remove: (id: string | number) => {
    return request.post<{ success: boolean }>('/guildSalaryCfg/deleteGuildSalaryCfg', {id})
  },
}

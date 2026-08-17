import {request} from '../request'
import type {AnchorSalaryCfg, AnchorSalaryCfgQuery, PageResponse} from '@/types/api'

export const anchorSalaryCfgApi = {
  getList: (params: AnchorSalaryCfgQuery) => {
    return request.post<PageResponse<AnchorSalaryCfg>>('/anchorSalaryCfg/anchorSalaryCfgList', params)
  },

  create: (data: {
    weeklyWorkDays: number
    dailyLiveDurationMinutes: number
    salaryAmount: number
    sort: number
  }) => {
    return request.post<{ id: string }>('/anchorSalaryCfg/createAnchorSalaryCfg', data)
  },

  update: (data: {
    id: string | number
    weeklyWorkDays: number
    dailyLiveDurationMinutes: number
    salaryAmount: number
    sort: number
  }) => {
    return request.post<{ success: boolean }>('/anchorSalaryCfg/updateAnchorSalaryCfg', data)
  },

  remove: (id: string | number) => {
    return request.post<{ success: boolean }>('/anchorSalaryCfg/deleteAnchorSalaryCfg', {id})
  },
}

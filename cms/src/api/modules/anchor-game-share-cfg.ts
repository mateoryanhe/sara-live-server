import {request} from '../request'
import type {AnchorGameShareCfg, PageResponse} from '@/types/api'

export const anchorGameShareCfgApi = {
  getList: (data: {salaryType: 1 | 2; pageIndex: number; pageSize: number}) => {
    return request.post<PageResponse<AnchorGameShareCfg>>(
      '/anchorGameShareCfg/anchorGameShareCfgList',
      data,
    )
  },

  create: (data: {
    salaryType: 1 | 2
    level: number
    gameTotalGoldRevenue: number
    anchorGameSharePercent: number
    guildGameSharePercent: number
  }) => {
    return request.post<{id: string}>(
      '/anchorGameShareCfg/createAnchorGameShareCfg',
      data,
    )
  },

  update: (data: {
    id: string | number
    level: number
    gameTotalGoldRevenue: number
    anchorGameSharePercent: number
    guildGameSharePercent: number
  }) => {
    return request.post<{success: boolean}>(
      '/anchorGameShareCfg/updateAnchorGameShareCfg',
      data,
    )
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>(
      '/anchorGameShareCfg/deleteAnchorGameShareCfg',
      {id},
    )
  },
}

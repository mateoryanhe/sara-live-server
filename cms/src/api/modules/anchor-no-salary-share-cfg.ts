import {request} from '../request'
import type {AnchorNoSalaryShareCfg, PageResponse} from '@/types/api'

export const anchorNoSalaryShareCfgApi = {
  getList: (data: {pageIndex: number; pageSize: number}) => {
    return request.post<PageResponse<AnchorNoSalaryShareCfg>>(
      '/anchorNoSalaryShareCfg/anchorNoSalaryShareCfgList',
      data,
    )
  },

  create: (data: {
    level: number
    socialTotalDiamondRevenue: number
    anchorSocialSharePercent: number
    guildSocialSharePercent: number
  }) => {
    return request.post<{id: string}>(
      '/anchorNoSalaryShareCfg/createAnchorNoSalaryShareCfg',
      data,
    )
  },

  update: (data: {
    id: string | number
    level: number
    socialTotalDiamondRevenue: number
    anchorSocialSharePercent: number
    guildSocialSharePercent: number
  }) => {
    return request.post<{success: boolean}>(
      '/anchorNoSalaryShareCfg/updateAnchorNoSalaryShareCfg',
      data,
    )
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>(
      '/anchorNoSalaryShareCfg/deleteAnchorNoSalaryShareCfg',
      {id},
    )
  },
}
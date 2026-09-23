import {request} from '../request'
import type {AnchorSalarySocialShareCfg, PageResponse} from '@/types/api'

export const anchorSalarySocialShareCfgApi = {
  getList: (data: {pageIndex: number; pageSize: number}) => {
    return request.post<PageResponse<AnchorSalarySocialShareCfg>>(
      '/anchorSalarySocialShareCfg/anchorSalarySocialShareCfgList',
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
      '/anchorSalarySocialShareCfg/createAnchorSalarySocialShareCfg',
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
      '/anchorSalarySocialShareCfg/updateAnchorSalarySocialShareCfg',
      data,
    )
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>(
      '/anchorSalarySocialShareCfg/deleteAnchorSalarySocialShareCfg',
      {id},
    )
  },
}

import {request} from '../request'
import type {PageResponse, StaticCacheRule, StaticCacheRuleQuery} from '@/types/api'

export const staticCacheCfgApi = {
  getList: (params: StaticCacheRuleQuery) => {
    return request.post<PageResponse<StaticCacheRule>>('/staticCacheCfg/staticCacheRuleList', params)
  },

  create: (data: {fileName: string; remark: string}) => {
    return request.post<{id: string}>('/staticCacheCfg/createStaticCacheRule', data)
  },

  update: (data: {id: string | number; fileName: string; remark: string}) => {
    return request.post<{success: boolean}>('/staticCacheCfg/updateStaticCacheRule', data)
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/staticCacheCfg/deleteStaticCacheRule', {id})
  },
}

export default staticCacheCfgApi

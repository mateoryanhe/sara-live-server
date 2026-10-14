import {request} from '../request'
import type {PromotionChannelUrlListResponse, PromotionChannelUrlQuery} from '@/types/api'

export const promotionChannelUrlApi = {
  getList: (params: PromotionChannelUrlQuery) => {
    return request.post<PromotionChannelUrlListResponse>('/promotionChannelUrl/promotionChannelUrlList', params)
  },

  create: (data: {channelName: string; channelCode: number}) => {
    return request.post<{id: string}>('/promotionChannelUrl/createPromotionChannelUrl', data)
  },

  update: (data: {id: string | number; channelName: string; channelCode: number}) => {
    return request.post<{success: boolean}>('/promotionChannelUrl/updatePromotionChannelUrl', data)
  },

  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/promotionChannelUrl/deletePromotionChannelUrl', {id})
  },
}

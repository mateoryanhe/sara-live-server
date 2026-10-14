import {request} from '../request'
import type {PromotionChannelDailyStatListResponse, PromotionChannelDailyStatQuery} from '@/types/api'

export const promotionChannelDailyStatApi = {
  getList: (params: PromotionChannelDailyStatQuery) => {
    return request.post<PromotionChannelDailyStatListResponse>('/promotionChannelUrl/promotionChannelDailyStatList', params)
  },
}

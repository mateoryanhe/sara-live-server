import {request} from '../request'
import type {LuckyGiftPoolPendingBalanceItem, PageResponse} from '@/types/api'

export interface ListLuckyGiftPoolPendingBalanceParams {
  status?: number
  pageIndex: number
  pageSize: number
}

export const luckyGiftPoolPendingBalanceApi = {
  listPendingBalance: (params: ListLuckyGiftPoolPendingBalanceParams) =>
    request.post<PageResponse<LuckyGiftPoolPendingBalanceItem>>(
      '/luckyGiftPoolCfg/listPendingBalance',
      params,
    ),
}

export default luckyGiftPoolPendingBalanceApi

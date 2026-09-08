import {request} from '../request'
import type {PageResponse} from '@/types/api'

export interface CoinMerchantItem {
  id: string
  username: string
  gold: number
  cancel: boolean
  channel: number
  createdAt?: string
}

export interface CoinMerchantQuery {
  pageIndex: number
  pageSize: number
  key?: string
  cancel?: number | null
}

export const coinMerchantApi = {
  list: (params: CoinMerchantQuery) => {
    return request.post<PageResponse<CoinMerchantItem>>('/coinMerchant/coinMerchantList', params)
  },
  create: (data: {username: string; password: string}) => {
    return request.post<{id: string; username: string}>('/coinMerchant/createCoinMerchant', data)
  },
  resetPassword: (data: {accountId: string | number; password: string}) => {
    return request.post<{success: boolean}>('/coinMerchant/resetCoinMerchantPassword', data)
  },
  cancel: (data: {accountId: string | number}) => {
    return request.post<{success: boolean}>('/coinMerchant/cancelCoinMerchant', data)
  },
}

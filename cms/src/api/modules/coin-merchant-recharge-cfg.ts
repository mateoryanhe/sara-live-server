import {request} from '../request'
import type {PageResponse} from '@/types/api'

export interface CoinMerchantRechargeCfg {
  id: string
  name: string
  price: number
  gold: number
  status: number
  createdAt?: string
  updatedAt?: string
}

export interface CoinMerchantRechargeCfgQuery {
  pageIndex: number
  pageSize: number
  name?: string
  statusFilter?: number
}

export const coinMerchantRechargeCfgApi = {
  list: (params: CoinMerchantRechargeCfgQuery) => {
    return request.post<PageResponse<CoinMerchantRechargeCfg>>(
        '/coinMerchantRechargeCfg/coinMerchantRechargeCfgList',
        params,
    )
  },
  create: (data: {name: string; price: number; gold: number}) => {
    return request.post<{id: string}>('/coinMerchantRechargeCfg/createCoinMerchantRechargeCfg', data)
  },
  update: (data: {id: string | number; name: string; price: number; gold: number}) => {
    return request.post<{success: boolean}>('/coinMerchantRechargeCfg/updateCoinMerchantRechargeCfg', data)
  },
  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/coinMerchantRechargeCfg/deleteCoinMerchantRechargeCfg', {id})
  },
  onShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/coinMerchantRechargeCfg/onShelfCoinMerchantRechargeCfg', {id})
  },
  offShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/coinMerchantRechargeCfg/offShelfCoinMerchantRechargeCfg', {id})
  },
}

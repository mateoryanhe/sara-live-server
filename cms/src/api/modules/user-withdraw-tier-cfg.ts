import {request} from '../request'
import type {PageResponse} from '@/types/api'

export interface UserWithdrawTierCfg {
  id: string
  name: string
  withdrawUsd: number
  minGameConsumeGold: number
  status: number
  createdAt?: string
  updatedAt?: string
}

export interface UserWithdrawTierCfgQuery {
  pageIndex: number
  pageSize: number
  name?: string
  statusFilter?: number
}

export const userWithdrawTierCfgApi = {
  list: (params: UserWithdrawTierCfgQuery) => {
    return request.post<PageResponse<UserWithdrawTierCfg>>('/userWithdrawTierCfg/userWithdrawTierCfgList', params)
  },
  create: (data: {name: string; withdrawUsd: number; minGameConsumeGold?: number}) => {
    return request.post<{id: string}>('/userWithdrawTierCfg/createUserWithdrawTierCfg', data)
  },
  update: (data: {id: string | number; name: string; withdrawUsd: number; minGameConsumeGold?: number}) => {
    return request.post<{success: boolean}>('/userWithdrawTierCfg/updateUserWithdrawTierCfg', data)
  },
  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/userWithdrawTierCfg/deleteUserWithdrawTierCfg', {id})
  },
  onShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/userWithdrawTierCfg/onShelfUserWithdrawTierCfg', {id})
  },
  offShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/userWithdrawTierCfg/offShelfUserWithdrawTierCfg', {id})
  },
}

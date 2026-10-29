import {request} from '../request'
import type {PageResponse} from '@/types/api'

export interface VoiceChatMicSeatTier {
  id: string
  seatCount: number
  sort: number
  status: number
  createdAt?: string
  updatedAt?: string
}

export interface VoiceChatMicSeatTierQuery {
  pageIndex: number
  pageSize: number
  seatCount?: number
  statusFilter?: number
}

export const voiceChatMicSeatTierApi = {
  list: (params: VoiceChatMicSeatTierQuery) => {
    return request.post<PageResponse<VoiceChatMicSeatTier>>('/voiceChatMicSeatTier/voiceChatMicSeatTierList', params)
  },
  create: (data: {seatCount: number; sort?: number}) => {
    return request.post<{id: string}>('/voiceChatMicSeatTier/createVoiceChatMicSeatTier', data)
  },
  update: (data: {id: string | number; seatCount: number; sort?: number}) => {
    return request.post<{success: boolean}>('/voiceChatMicSeatTier/updateVoiceChatMicSeatTier', data)
  },
  remove: (id: string | number) => {
    return request.post<{success: boolean}>('/voiceChatMicSeatTier/deleteVoiceChatMicSeatTier', {id})
  },
  onShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/voiceChatMicSeatTier/onShelfVoiceChatMicSeatTier', {id})
  },
  offShelf: (id: string | number) => {
    return request.post<{success: boolean; status: number}>('/voiceChatMicSeatTier/offShelfVoiceChatMicSeatTier', {id})
  },
}

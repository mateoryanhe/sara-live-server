import {request} from '../request'
import type {AnchorVideoPriceTier, AnchorVideoPriceTierQuery, PageResponse} from '@/types/api'

export const anchorVideoPriceTierApi = {
  getList: (params: AnchorVideoPriceTierQuery) => {
    return request.post<PageResponse<AnchorVideoPriceTier>>('/anchorVideoPriceTier/anchorVideoPriceTierList', params)
  },

  create: (data: { price: number }) => {
    return request.post<{ id: string }>('/anchorVideoPriceTier/createAnchorVideoPriceTier', data)
  },

  update: (data: { id: string | number; price: number }) => {
    return request.post<boolean>('/anchorVideoPriceTier/updateAnchorVideoPriceTier', data)
  },

  delete: (id: string | number) => {
    return request.post<boolean>('/anchorVideoPriceTier/deleteAnchorVideoPriceTier', {id})
  },

  onShelf: (id: string | number) => {
    return request.post<{ success: boolean; status: number }>('/anchorVideoPriceTier/onShelfAnchorVideoPriceTier', {id})
  },

  offShelf: (id: string | number) => {
    return request.post<{ success: boolean; status: number }>('/anchorVideoPriceTier/offShelfAnchorVideoPriceTier', {id})
  },
}

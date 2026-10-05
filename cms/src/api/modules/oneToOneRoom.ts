import {request} from '../request'
import type {OneToOneRoomItem, PageQuery, PageResponse} from '@/types/api'

export interface OneToOneRoomQuery extends PageQuery {
    key?: string
    status?: number
}

export const oneToOneRoomApi = {
    list: (params: OneToOneRoomQuery) => {
        return request.post<PageResponse<OneToOneRoomItem>>('/oneToOneRoom/list', params)
    },
    create: (data: { userId: string; billing: number }) => {
        return request.post<{ userId: string }>('/oneToOneRoom/create', data)
    },
    update: (data: { userId: string; billing: number }) => {
        return request.post<boolean>('/oneToOneRoom/update', data)
    },
    setStatus: (userId: string, status: number) => {
        return request.post<boolean>('/oneToOneRoom/setStatus', {userId, status})
    },
}

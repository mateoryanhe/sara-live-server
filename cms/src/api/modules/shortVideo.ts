import {request} from '../request'
import type {PageResponse, ShortVideo, ShortVideoQuery} from '@/types/api'

export const shortVideoApi = {
    getShortVideoList: (params: ShortVideoQuery) => {
        return request.post<PageResponse<ShortVideo>>('/shortVideo/shortVideoList', params)
    },

    createShortVideo: (data: {
        title: string
        video: string
        cover: string
        sort: number
        description: string
    }) => {
        return request.post<{ id: string }>('/shortVideo/createShortVideo', data)
    },

    updateShortVideo: (data: {
        id: string | number
        title: string
        video: string
        cover: string
        sort: number
        description: string
    }) => {
        return request.post<boolean>('/shortVideo/updateShortVideo', data)
    },

    deleteShortVideo: (id: string | number) => {
        return request.post<boolean>('/shortVideo/deleteShortVideo', {id})
    },

    onShelfShortVideo: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/shortVideo/onShelfShortVideo', {id})
    },

    offShelfShortVideo: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/shortVideo/offShelfShortVideo', {id})
    },
}

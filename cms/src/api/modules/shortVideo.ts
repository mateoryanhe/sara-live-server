import {request} from '../request'
import type {
    GetShortVideoCfgRes,
    PageResponse,
    SaveShortVideoCfgReq,
    SaveShortVideoCfgRes,
    ShortVideo,
    ShortVideoCategory,
    ShortVideoCategoryQuery,
    ShortVideoQuery,
    ShortVideoWatchQuery,
    ShortVideoWatchRecord,
} from '@/types/api'

export const shortVideoApi = {
    getShortVideoList: (params: ShortVideoQuery) => {
        return request.post<PageResponse<ShortVideo>>('/shortVideo/shortVideoList', params)
    },

    getShortVideoWatchList: (params: ShortVideoWatchQuery) => {
        return request.post<PageResponse<ShortVideoWatchRecord>>('/shortVideo/shortVideoWatchList', params)
    },

    getShortVideoCfg: () => {
        return request.post<GetShortVideoCfgRes>('/shortVideo/getShortVideoCfg', {})
    },

    saveShortVideoCfg: (data: SaveShortVideoCfgReq) => {
        return request.post<SaveShortVideoCfgRes>('/shortVideo/saveShortVideoCfg', data)
    },

    updateShortVideo: (data: {
        id: string | number
        title: string
        cover: string
        sort: number
        isPaid: number
        diamondPerMinute: number
        categoryId: number
        source: number
        authorId: string | number
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

    getShortVideoCategoryList: (params: ShortVideoCategoryQuery) => {
        return request.post<PageResponse<ShortVideoCategory>>('/shortVideo/shortVideoCategoryList', params)
    },

    createShortVideoCategory: (data: { name: string; sort: number }) => {
        return request.post<{ id: string }>('/shortVideo/createShortVideoCategory', data)
    },

    updateShortVideoCategory: (data: { id: string | number; name: string; sort: number }) => {
        return request.post<boolean>('/shortVideo/updateShortVideoCategory', data)
    },

    deleteShortVideoCategory: (id: string | number) => {
        return request.post<boolean>('/shortVideo/deleteShortVideoCategory', {id})
    },
}

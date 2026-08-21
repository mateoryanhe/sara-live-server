import {request} from '../request'
import type {
    CreateShortVideoReq,
    CreateShortVideoRes,
    GetShortVideoCfgRes,
    PageResponse,
    SaveShortVideoCfgReq,
    SaveShortVideoCfgRes,
    ShortVideo,
    ShortVideoCategory,
    ShortVideoCategoryQuery,
    ShortVideoPriceTier,
    ShortVideoPriceTierQuery,
    ShortVideoQuery,
    ShortVideoStorageStat,
    ShortVideoWatchQuery,
    ShortVideoWatchRecord,
} from '@/types/api'

const SHORT_VIDEO_LIST_TIMEOUT_MS = 30 * 60 * 1000
const SHORT_VIDEO_STORAGE_STAT_TIMEOUT_MS = 10 * 60 * 1000

export const shortVideoApi = {
    getShortVideoList: (params: ShortVideoQuery) => {
        return request.post<PageResponse<ShortVideo>>('/shortVideo/shortVideoList', params, {
            timeout: SHORT_VIDEO_LIST_TIMEOUT_MS,
        })
    },

    createShortVideo: (data: CreateShortVideoReq) => {
        return request.post<CreateShortVideoRes>('/shortVideo/createShortVideo', data)
    },

    getShortVideoWatchList: (params: ShortVideoWatchQuery) => {
        return request.post<PageResponse<ShortVideoWatchRecord>>('/shortVideo/shortVideoWatchList', params)
    },

    getShortVideoCfg: () => {
        return request.post<GetShortVideoCfgRes>('/shortVideo/getShortVideoCfg', {})
    },

    getShortVideoStorageStat: () => {
        return request.post<ShortVideoStorageStat>('/shortVideo/shortVideoStorageStat', {}, {
            timeout: SHORT_VIDEO_STORAGE_STAT_TIMEOUT_MS,
        })
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
        payDiamond: number
        freeWatchSeconds: number
        categoryId: number
        source: number
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

    getShortVideoPriceTierList: (params: ShortVideoPriceTierQuery) => {
        return request.post<PageResponse<ShortVideoPriceTier>>('/shortVideo/shortVideoPriceTierList', params)
    },

    createShortVideoPriceTier: (data: { price: number }) => {
        return request.post<{ id: string }>('/shortVideo/createShortVideoPriceTier', data)
    },

    updateShortVideoPriceTier: (data: { id: string | number; price: number }) => {
        return request.post<boolean>('/shortVideo/updateShortVideoPriceTier', data)
    },

    deleteShortVideoPriceTier: (id: string | number) => {
        return request.post<boolean>('/shortVideo/deleteShortVideoPriceTier', {id})
    },

    onShelfShortVideoPriceTier: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/shortVideo/onShelfShortVideoPriceTier', {id})
    },

    offShelfShortVideoPriceTier: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/shortVideo/offShelfShortVideoPriceTier', {id})
    },
}

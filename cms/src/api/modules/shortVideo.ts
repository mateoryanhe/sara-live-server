import {request} from '../request'
import type {
    GetShortVideoCfgRes,
    PageResponse,
    SaveShortVideoCfgReq,
    SaveShortVideoCfgRes,
    ShortVideo,
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

    uploadShortVideo: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<{ fileName: string; url: string }>('/shortVideo/uploadShortVideo', formData, {
            headers: {'Content-Type': 'multipart/form-data'},
        })
    },

    createShortVideo: (data: {
        title: string
        video: string
        cover: string
        sort: number
        isPaid: number
        diamondPerSecond: number
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
        isPaid: number
        diamondPerSecond: number
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

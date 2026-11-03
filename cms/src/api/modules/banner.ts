import {request} from '../request'
import type {Banner, BannerQuery, PageResponse} from '@/types/api'

export type BannerSavePayload = {
    title: string
    image: string
    link: string
    scene: number
    direction: number
    sort: number
    excludePackageNames?: string
}

export const bannerApi = {
    getBannerList: (params: BannerQuery) => {
        return request.post<PageResponse<Banner>>('/banner/bannerList', params)
    },

    createBanner: (data: BannerSavePayload) => {
        return request.post<{ id: string }>('/banner/createBanner', data)
    },

    updateBanner: (data: BannerSavePayload & { id: string | number }) => {
        return request.post<boolean>('/banner/updateBanner', data)
    },

    deleteBanner: (id: string | number) => {
        return request.post<boolean>('/banner/deleteBanner', {id})
    },

    onShelfBanner: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/banner/onShelfBanner', {id})
    },

    offShelfBanner: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/banner/offShelfBanner', {id})
    },
}

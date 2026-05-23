import {request} from '../request'
import type {Gift, GiftQuery, PageResponse} from '@/types/api'

// 礼物管理API
export const giftApi = {
    // 获取礼物列表
    getGiftList: (params: GiftQuery) => {
        return request.post<PageResponse<Gift>>('/gift/giftList', params)
    },

    // 创建礼物
    createGift: (data: {
        name: string
        icon: string
        animation: string
        price: number
        category: string
        sort: number
        description: string
    }) => {
        return request.post<{ id: string }>('/gift/createGift', data)
    },

    // 更新礼物
    updateGift: (data: {
        id: string | number
        name: string
        icon: string
        animation: string
        price: number
        category: string
        sort: number
        description: string
    }) => {
        return request.post<boolean>('/gift/updateGift', data)
    },

    // 删除礼物
    deleteGift: (id: string | number) => {
        return request.post<boolean>('/gift/deleteGift', {id})
    },

    // 上架礼物
    onShelfGift: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/gift/onShelfGift', {id})
    },

    // 下架礼物
    offShelfGift: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/gift/offShelfGift', {id})
    },
}

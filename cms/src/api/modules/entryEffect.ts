import {request} from '../request'
import type {EntryEffect, EntryEffectQuery, PageResponse} from '@/types/api'

export const entryEffectApi = {
    getEntryEffectList: (params: EntryEffectQuery) => {
        return request.post<PageResponse<EntryEffect>>('/entryEffect/entryEffectList', params)
    },

    createEntryEffect: (data: {
        name: string
        levelStart: number
        levelEnd: number
        animation: string
    }) => {
        return request.post<{ id: string }>('/entryEffect/createEntryEffect', data)
    },

    updateEntryEffect: (data: {
        id: string | number
        name: string
        levelStart: number
        levelEnd: number
        animation: string
    }) => {
        return request.post<boolean>('/entryEffect/updateEntryEffect', data)
    },

    deleteEntryEffect: (id: string | number) => {
        return request.post<boolean>('/entryEffect/deleteEntryEffect', {id})
    },

    onShelfEntryEffect: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/entryEffect/onShelfEntryEffect', {id})
    },

    offShelfEntryEffect: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/entryEffect/offShelfEntryEffect', {id})
    },
}

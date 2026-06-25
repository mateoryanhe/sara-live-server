import {request} from '../request'
import type {LiveRoomTag, LiveRoomTagQuery, PageResponse} from '@/types/api'

export const liveRoomTagApi = {
    getLiveRoomTagList: (params: LiveRoomTagQuery) => {
        return request.post<PageResponse<LiveRoomTag>>('/liveRoomTag/liveRoomTagList', params)
    },

    createLiveRoomTag: (data: { name: string; sort: number }) => {
        return request.post<{ id: string }>('/liveRoomTag/createLiveRoomTag', data)
    },

    updateLiveRoomTag: (data: { id: string | number; name: string; sort: number }) => {
        return request.post<boolean>('/liveRoomTag/updateLiveRoomTag', data)
    },

    deleteLiveRoomTag: (id: string | number) => {
        return request.post<boolean>('/liveRoomTag/deleteLiveRoomTag', {id})
    },
}

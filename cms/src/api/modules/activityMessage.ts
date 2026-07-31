import {request} from '../request'
import type {ActivityMessage, ActivityMessageQuery, PageResponse} from '@/types/api'

export type ActivityMessageForm = {
    id: string
    iconEn: string
    iconEs: string
    iconPt: string
    iconHi: string
    bgEn: string
    bgEs: string
    bgPt: string
    bgHi: string
    titleEn: string
    titleEs: string
    titlePt: string
    titleHi: string
    contentEn: string
    contentEs: string
    contentPt: string
    contentHi: string
    urlEn: string
    urlEs: string
    urlPt: string
    urlHi: string
}

export const activityMessageApi = {
    getActivityMessageList: (params: ActivityMessageQuery) => {
        return request.post<PageResponse<ActivityMessage>>('/activityMessage/activityMessageList', params)
    },

    createActivityMessage: (data: Omit<ActivityMessageForm, 'id'>) => {
        return request.post<{ id: string }>('/activityMessage/createActivityMessage', data)
    },

    updateActivityMessage: (data: ActivityMessageForm & { id: string | number }) => {
        return request.post<{ success: boolean }>('/activityMessage/updateActivityMessage', data)
    },

    deleteActivityMessage: (id: string | number) => {
        return request.post<{ success: boolean }>('/activityMessage/deleteActivityMessage', {id})
    },

    publishActivityMessage: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/activityMessage/publishActivityMessage', {id})
    },

    unpublishActivityMessage: (id: string | number) => {
        return request.post<{ success: boolean; status: number }>('/activityMessage/unpublishActivityMessage', {id})
    },
}

export default activityMessageApi

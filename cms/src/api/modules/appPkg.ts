import {request} from '../request'
import type {AppPkg, AppPkgQuery, PageResponse} from '@/types/api'

export type AppPkgSavePayload = {
    packageName: string
    remark?: string
    attributionEnabled?: boolean
    attributionProvider?: string
    appsFlyerDevKey?: string
    appsFlyerAppId?: string
}

export const appPkgApi = {
    getAppPkgList: (params: AppPkgQuery) => {
        return request.post<PageResponse<AppPkg>>('/appPkg/appPkgList', params)
    },

    createAppPkg: (data: AppPkgSavePayload) => {
        return request.post<{ id: string }>('/appPkg/createAppPkg', data)
    },

    updateAppPkg: (data: AppPkgSavePayload & { id: string | number }) => {
        return request.post<boolean>('/appPkg/updateAppPkg', data)
    },

    deleteAppPkg: (id: string | number) => {
        return request.post<boolean>('/appPkg/deleteAppPkg', {id})
    },
}

export default appPkgApi

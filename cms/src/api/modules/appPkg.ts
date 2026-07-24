import {request} from '../request'
import type {AppPkg, AppPkgQuery, PageResponse} from '@/types/api'

export const appPkgApi = {
    getAppPkgList: (params: AppPkgQuery) => {
        return request.post<PageResponse<AppPkg>>('/appPkg/appPkgList', params)
    },

    createAppPkg: (data: {
        packageName: string
        secretKey: string
        privacyPolicyUrl?: string
        termsOfServiceUrl?: string
        remark?: string
    }) => {
        return request.post<{ id: string }>('/appPkg/createAppPkg', data)
    },

    updateAppPkg: (data: {
        id: string | number
        packageName: string
        secretKey: string
        privacyPolicyUrl?: string
        termsOfServiceUrl?: string
        remark?: string
    }) => {
        return request.post<boolean>('/appPkg/updateAppPkg', data)
    },

    deleteAppPkg: (id: string | number) => {
        return request.post<boolean>('/appPkg/deleteAppPkg', {id})
    },
}

export default appPkgApi

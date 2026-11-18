import {request} from '../request'

export interface AppsFlyerCfgItem {
    id: string
    androidDevKey?: string
    iosDevKey?: string
    iosAppId?: string
    createdAt?: string
    updatedAt?: string
}

export interface GetAppsFlyerCfgRes {
    cfg?: AppsFlyerCfgItem | null
}

export interface SaveAppsFlyerCfgReq {
    id: number
    androidDevKey?: string
    iosDevKey?: string
    iosAppId?: string
}

export interface SaveAppsFlyerCfgRes {
    success: boolean
    id?: string
}

export const appsFlyerCfgApi = {
    getAppsFlyerCfg: () => request.post<GetAppsFlyerCfgRes>('/appsFlyer/getAppsFlyerCfg', {}),
    saveAppsFlyerCfg: (data: SaveAppsFlyerCfgReq) =>
        request.post<SaveAppsFlyerCfgRes>('/appsFlyer/saveAppsFlyerCfg', data),
}

export default appsFlyerCfgApi
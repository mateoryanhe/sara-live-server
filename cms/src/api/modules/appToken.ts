import {request} from '../request'
import type {AppToken, GetAppTokenReq, PageResponse, SaveAppTokenReq} from '@/types/api'

const appTokenApi = {
    getAppToken: (data?: GetAppTokenReq) => {
        return request.post<PageResponse<AppToken>>('/appToken/getAppToken', data)
    },
    saveAppToken: (data: SaveAppTokenReq) => {
        return request.post<boolean>('/appToken/saveAppToken', data)
    },
}

export default appTokenApi

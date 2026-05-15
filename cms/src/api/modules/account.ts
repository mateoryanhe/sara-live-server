import {request} from '../request'
import type {BanReq, CancelReq, PageResponse, QueryUserInfoReq, UnBanReq, UnCancelReq, UserInfo} from '@/types/api'

const accountApi = {
    // 封号
    ban: (data: BanReq) => {
        return request.post<boolean>('/account/ban', data)
    },

    // 解封
    unBan: (data: UnBanReq) => {
        return request.post<boolean>('/account/unBan', data)
    },

    // 注销
    cancel: (data: CancelReq) => {
        return request.post<boolean>('/account/cancel', data)
    },

    // 取消注销
    unCancel: (data: UnCancelReq) => {
        return request.post<boolean>('/account/unCancel', data)
    },

    // 获取用户信息
    getUserInfo: (data: QueryUserInfoReq) => {
        return request.post<PageResponse<UserInfo>>('/account/getUserInfo', data)
    }
}

export default accountApi

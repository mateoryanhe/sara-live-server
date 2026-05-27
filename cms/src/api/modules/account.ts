import {request} from '../request'
import type {
    AnchorListItem,
    BanAnchorReq,
    BanAnchorReq,
    BanReq,
    CancelReq,
    PageResponse,
    QueryAnchorListReq,
    QueryUserInfoReq,
    SetAnchorReq,
    UnBanAnchorReq,
    UnBanReq,
    UnCancelReq,
    UserInfo
} from '@/types/api'

const accountApi = {
    // 封号
    ban: (data: BanReq) => {
        return request.post<boolean>('/account/ban', data)
    },

    // 封禁主播(含App推送)
    banAnchor: (data: BanAnchorReq) => {
        return request.post<boolean>('/account/banAnchor', data)
    },

    // 解封主播直播间
    unBanAnchor: (data: UnBanAnchorReq) => {
        return request.post<boolean>('/account/unBanAnchor', data)
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
    },

    getAnchorList: (data: QueryAnchorListReq) => {
        return request.post<PageResponse<AnchorListItem>>('/account/getAnchorList', data)
    },

    setAnchor: (data: SetAnchorReq) => {
        return request.post<boolean>('/account/setAnchor', data)
    }
}

export default accountApi

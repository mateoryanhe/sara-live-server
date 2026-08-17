import {request} from '../request'
import type {
    AnchorDetail,
    AnchorDailyEffectiveLiveItem,
    AnchorDailyEffectiveLiveQuery,
    AnchorListItem,
    BanAnchorReq,
    BatchSetAnchorReq,
    BatchSetAnchorRes,
    BatchSetSeniorAnchorReq,
    CancelReq,
    ExitGuildReq,
    ExitGuildRes,
    PageResponse,
    QueryAnchorListReq,
    QueryOffShelfLiveRoomListReq,
    OffShelfLiveRoomItem,
    QueryUserInfoReq,
    SetAnchorReq,
    SetCanRankReq,
    SetLiveRoomStatusReq,
    SetLiveRoomStatusRes,
    SetRechargeWhitelistReq,
    SetSeniorAnchorReq,
    SetUserTypeReq,
    UnBanAnchorReq,
    UnBanReq,
    UnCancelReq,
    UserDetail,
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

    getAnchorDetail: (anchorId: string | number) => {
        return request.post<AnchorDetail>('/account/getAnchorDetail', {anchorId})
    },

    getAnchorDailyEffectiveLiveList: (data: AnchorDailyEffectiveLiveQuery) => {
        return request.post<PageResponse<AnchorDailyEffectiveLiveItem>>('/account/getAnchorDailyEffectiveLiveList', data)
    },

    getUserDetail: (userId: string | number) => {
        return request.post<UserDetail>('/account/getUserDetail', {userId})
    },

    getOffShelfLiveRoomList: (data: QueryOffShelfLiveRoomListReq) => {
        return request.post<PageResponse<OffShelfLiveRoomItem>>('/account/getOffShelfLiveRoomList', data)
    },

    setAnchor: (data: SetAnchorReq) => {
        return request.post<boolean>('/account/setAnchor', data)
    },

    setSeniorAnchor: (data: SetSeniorAnchorReq) => {
        return request.post<boolean>('/account/setSeniorAnchor', data)
    },

    batchSetAnchor: (data: BatchSetAnchorReq) => {
        return request.post<BatchSetAnchorRes>('/account/batchSetAnchor', data)
    },

    batchSetSeniorAnchor: (data: BatchSetSeniorAnchorReq) => {
        return request.post<BatchSetAnchorRes>('/account/batchSetSeniorAnchor', data)
    },

    setUserType: (data: SetUserTypeReq) => {
        return request.post<boolean>('/account/setUserType', data)
    },

    setCanRank: (data: SetCanRankReq) => {
        return request.post<boolean>('/account/setCanRank', data)
    },

    setRechargeWhitelist: (data: SetRechargeWhitelistReq) => {
        return request.post<boolean>('/account/setRechargeWhitelist', data)
    },

    // 上架/下架主播直播间
    setLiveRoomStatus: (data: SetLiveRoomStatusReq) => {
        return request.post<SetLiveRoomStatusRes>('/account/setLiveRoomStatus', data)
    },

    // 退出工会(将工会ID置为0)
    exitGuild: (data: ExitGuildReq) => {
        return request.post<ExitGuildRes>('/account/exitGuild', data)
    }
}

export default accountApi

import {request} from '../request'
import type {
    AnchorDailyEffectiveLiveItem,
    AnchorIncomeSettlementLogItem,
    Guild,
    GuildDailyEffectiveLiveItem,
    GuildDailyEffectiveLiveQuery,
    GuildAnchorDailyEffectiveLiveItem,
    GuildAnchorDailyEffectiveLiveQuery,
    GuildDetailIncome,
    GuildIncomeArchivesRes,
    GuildQuery,
    GuildTransferInfo,
    ImportGuildAnchorsReq,
    ImportGuildAnchorsRes,
    JoinGuildAnchorReq,
    JoinGuildAnchorRes,
    SaveGuildTransferInfoReq,
    SetGuildAnchorTypeReq,
    SetGuildAnchorTypeRes,
    MyGuildAnchorIncomeSettlementLogQuery,
    MyGuildAnchorListItem,
    MyGuildAnchorListQuery,
    MyGuildAnchorDailyEffectiveLiveQuery,
    MyOwnedGuildAnchorDailyEffectiveLiveQuery,
    MyOwnedGuildAnchorListQuery,
    MyGuildProfileListRes,
    PageResponse,
    UpdateMyGuildProfileReq,
} from '@/types/api'

// 直播工会管理API
export const guildApi = {
    // 获取工会列表
    getGuildList: (params: GuildQuery) => {
        return request.post<PageResponse<Guild>>('/guild/guildList', params)
    },

    // 创建工会
    createGuild: (data: {
        name: string
        leaderId: number
        description: string
    }) => {
        return request.post<{ id: string }>('/guild/createGuild', data)
    },

    // 更新工会
    updateGuild: (data: {
        id: string
        name: string
        leaderId: number
        description: string
    }) => {
        return request.post<boolean>('/guild/updateGuild', data)
    },

    // 下架工会
    deleteGuild: (id: string) => {
        return request.post<boolean>('/guild/deleteGuild', {id})
    },

    // 工会垃圾库：已下架列表
    getOffShelfGuildList: (params: GuildQuery) => {
        return request.post<PageResponse<Guild>>('/guild/offShelfGuildList', params)
    },

    // 上架工会
    onShelfGuild: (id: string) => {
        return request.post<{ success: boolean }>('/guild/onShelfGuild', {id})
    },

    getMyGuildProfile: () => {
        return request.post<MyGuildProfileListRes>('/guild/getMyGuildProfile', {})
    },

    getMyGuildAnchorList: (params: MyGuildAnchorListQuery) => {
        return request.post<PageResponse<MyGuildAnchorListItem>>('/guild/getMyGuildAnchorList', params)
    },

    getMyGuildAnchorDailyEffectiveLiveList: (params: MyGuildAnchorDailyEffectiveLiveQuery) => {
        return request.post<PageResponse<AnchorDailyEffectiveLiveItem>>('/guild/getMyGuildAnchorDailyEffectiveLiveList', params)
    },

    getGuildAnchorIncomeSettlementLogList: (params: MyGuildAnchorIncomeSettlementLogQuery) => {
        return request.post<PageResponse<AnchorIncomeSettlementLogItem>>('/guild/cmsGuildAnchorIncomeSettlementLogList', params)
    },

    updateMyGuildProfile: (data: UpdateMyGuildProfileReq) => {
        return request.post<{ success: boolean }>('/guild/updateMyGuildProfile', data)
    },

    importGuildAnchors: (data: ImportGuildAnchorsReq) => {
        return request.post<ImportGuildAnchorsRes>('/guild/importGuildAnchors', data)
    },

    joinGuildAnchor: (data: JoinGuildAnchorReq) => {
        return request.post<JoinGuildAnchorRes>('/guild/joinGuildAnchor', data)
    },

    setGuildAnchorType: (data: SetGuildAnchorTypeReq) => {
        return request.post<SetGuildAnchorTypeRes>('/guild/setGuildAnchorType', data)
    },

    getGuildTransferInfo: (guildId: string | number) => {
        return request.post<{info: GuildTransferInfo | null}>('/guild/getGuildTransferInfo', {guildId})
    },

    saveGuildTransferInfo: (data: SaveGuildTransferInfoReq) => {
        return request.post<{success: boolean}>('/guild/saveGuildTransferInfo', data)
    },

    getGuildDetail: (guildId: string) => {
        return request.post<GuildDetailIncome>('/guild/getGuildDetail', {guildId})
    },

    getGuildIncomeArchives: (guildId: string) => {
        return request.post<GuildIncomeArchivesRes>('/guild/getGuildIncomeArchives', {guildId})
    },

    getGuildDailyEffectiveLiveList: (data: GuildDailyEffectiveLiveQuery) => {
        return request.post<PageResponse<GuildDailyEffectiveLiveItem>>('/guild/getGuildDailyEffectiveLiveList', data)
    },

    getGuildAnchorDailyEffectiveLiveList: (data: GuildAnchorDailyEffectiveLiveQuery) => {
        return request.post<PageResponse<GuildAnchorDailyEffectiveLiveItem>>('/guild/cmsGuildAnchorDailyEffectiveLiveList', data)
    },

    getMyOwnedGuildAnchorDailyEffectiveLiveList: (params: MyOwnedGuildAnchorDailyEffectiveLiveQuery) => {
        return request.post<PageResponse<GuildAnchorDailyEffectiveLiveItem>>('/guild/cmsMyGuildAnchorDailyEffectiveLiveList', params)
    },

    getMyOwnedGuildAnchorList: (params: MyOwnedGuildAnchorListQuery) => {
        return request.post<PageResponse<MyGuildAnchorListItem>>('/guild/getMyOwnedGuildAnchorList', params)
    },
}

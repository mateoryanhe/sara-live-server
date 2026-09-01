import request from '../request'
import type {
    GetDataSyncCfgRes,
    SaveDataSyncCfgReq,
    SaveDataSyncCfgRes,
    SyncActivityMessageReq,
    SyncActivityMessageRes,
    SyncBatchRes,
    SyncFirstRechargeActivityCfgRes,
    SyncIdsReq,
    SyncVipCfgReq,
    SyncVipCfgRes,
} from '@/types/api'

export const dataSyncApi = {
    getDataSyncCfg: () => {
        return request.post<GetDataSyncCfgRes>('/dataSync/getDataSyncCfg', {})
    },

    saveDataSyncCfg: (data: SaveDataSyncCfgReq) => {
        return request.post<SaveDataSyncCfgRes>('/dataSync/saveDataSyncCfg', data)
    },

    syncVipCfg: (data: SyncVipCfgReq) => {
        return request.post<SyncVipCfgRes>('/dataSync/syncVipCfg', data)
    },

    syncVipCfgAssets: (data: SyncVipCfgReq) => {
        return request.post<SyncVipCfgRes>('/dataSync/syncVipCfgAssets', data)
    },

    syncActivityMessage: (data: SyncActivityMessageReq) => {
        return request.post<SyncActivityMessageRes>('/dataSync/syncActivityMessage', data)
    },

    syncBanner: (data: SyncIdsReq) => {
        return request.post<SyncBatchRes>('/dataSync/syncBanner', data)
    },

    syncGift: (data: SyncIdsReq) => {
        return request.post<SyncBatchRes>('/dataSync/syncGift', data)
    },

    syncGiftAssets: (data: SyncIdsReq) => {
        return request.post<SyncBatchRes>('/dataSync/syncGiftAssets', data)
    },

    syncRechargeCfg: (data: SyncIdsReq) => {
        return request.post<SyncBatchRes>('/dataSync/syncRechargeCfg', data)
    },

    syncFirstRechargeActivityCfg: () => {
        return request.post<SyncFirstRechargeActivityCfgRes>('/dataSync/syncFirstRechargeActivityCfg', {})
    },
}

export default dataSyncApi

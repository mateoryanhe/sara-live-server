import request from '../request'
import type {GetDataSyncCfgRes, SaveDataSyncCfgReq, SaveDataSyncCfgRes, SyncVipCfgReq, SyncVipCfgRes} from '@/types/api'

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
}

export default dataSyncApi

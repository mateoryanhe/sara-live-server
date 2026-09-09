import request from '../request'
import type {
    GetSyncLocalStorageToS3StatusRes,
    GetUploadResourceCfgRes,
    SaveUploadResourceCfgReq,
    SaveUploadResourceCfgRes,
    SyncLocalStorageToS3Res,
} from '@/types/api'

export const uploadResourceApi = {
    getUploadResourceCfg: () => {
        return request.post<GetUploadResourceCfgRes>('/upload/getUploadResourceCfg', {})
    },

    saveUploadResourceCfg: (data: SaveUploadResourceCfgReq) => {
        return request.post<SaveUploadResourceCfgRes>('/upload/saveUploadResourceCfg', data)
    },

    syncLocalStorageToS3: () => {
        return request.post<SyncLocalStorageToS3Res>('/upload/syncLocalStorageToS3', {})
    },

    getSyncLocalStorageToS3Status: () => {
        return request.post<GetSyncLocalStorageToS3StatusRes>('/upload/getSyncLocalStorageToS3Status', {})
    },
}

export default uploadResourceApi

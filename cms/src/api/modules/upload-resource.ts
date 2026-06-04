import request from '../request'
import type {
    GetUploadResourceCfgRes,
    SaveUploadResourceCfgReq,
    SaveUploadResourceCfgRes,
} from '@/types/api'

export const uploadResourceApi = {
    getUploadResourceCfg: () => {
        return request.post<GetUploadResourceCfgRes>('/upload/getUploadResourceCfg', {})
    },

    saveUploadResourceCfg: (data: SaveUploadResourceCfgReq) => {
        return request.post<SaveUploadResourceCfgRes>('/upload/saveUploadResourceCfg', data)
    },
}

export default uploadResourceApi

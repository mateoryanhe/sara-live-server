import request from '../request'
import type {
    GetCfEmailCfgRes,
    SaveCfEmailCfgReq,
    SaveCfEmailCfgRes,
    SendCfEmailTestReq,
    SendCfEmailTestRes,
} from '@/types/api'

export const cfEmailApi = {
    getCfEmailCfg: () => {
        return request.post<GetCfEmailCfgRes>('/cfEmail/getCfEmailCfg', {})
    },

    saveCfEmailCfg: (data: SaveCfEmailCfgReq) => {
        return request.post<SaveCfEmailCfgRes>('/cfEmail/saveCfEmailCfg', data)
    },

    sendCfEmailTest: (data: SendCfEmailTestReq) => {
        return request.post<SendCfEmailTestRes>('/cfEmail/sendCfEmailTest', data)
    },
}

export default cfEmailApi

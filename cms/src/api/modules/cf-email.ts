import request from '../request'
import type {GetCfEmailCfgRes, SaveCfEmailCfgReq, SaveCfEmailCfgRes} from '@/types/api'

export const cfEmailApi = {
    getCfEmailCfg: () => {
        return request.post<GetCfEmailCfgRes>('/cfEmail/getCfEmailCfg', {})
    },

    saveCfEmailCfg: (data: SaveCfEmailCfgReq) => {
        return request.post<SaveCfEmailCfgRes>('/cfEmail/saveCfEmailCfg', data)
    },
}

export default cfEmailApi

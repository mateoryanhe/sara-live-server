import {request} from '../request'
import type {GetLiveCfgRes, SaveLiveCfgReq, SaveLiveCfgRes} from '@/types/api'

export const liveCfgApi = {
    getLiveCfg: () => {
        return request.post<GetLiveCfgRes>('/liveCfg/getLiveCfg', {})
    },

    saveLiveCfg: (data: SaveLiveCfgReq) => {
        return request.post<SaveLiveCfgRes>('/liveCfg/saveLiveCfg', data)
    },
}

export default liveCfgApi

import request from '../request'
import type {GetPreloadCfgRes, SavePreloadCfgReq, SavePreloadCfgRes} from '@/types/api'

export const preloadCfgApi = {
    getPreloadCfg: () => {
        return request.post<GetPreloadCfgRes>('/preloadCfg/getPreloadCfg', {})
    },

    savePreloadCfg: (data: SavePreloadCfgReq) => {
        return request.post<SavePreloadCfgRes>('/preloadCfg/savePreloadCfg', data)
    },
}

export default preloadCfgApi

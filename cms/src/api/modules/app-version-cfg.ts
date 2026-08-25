import request from '../request'
import type {GetAppVersionCfgRes, SaveAppVersionCfgReq, SaveAppVersionCfgRes} from '@/types/api'

export const appVersionCfgApi = {
    getAppVersionCfg: () => {
        return request.post<GetAppVersionCfgRes>('/appVersionCfg/getAppVersionCfg', {})
    },

    saveAppVersionCfg: (data: SaveAppVersionCfgReq) => {
        return request.post<SaveAppVersionCfgRes>('/appVersionCfg/saveAppVersionCfg', data)
    },
}

export default appVersionCfgApi

import request from '../request'
import type {GetGooglePlayCfgRes, SaveGooglePlayCfgReq, SaveGooglePlayCfgRes} from '@/types/api'

export const googlePlayApi = {
    getGooglePlayCfg: () => {
        return request.post<GetGooglePlayCfgRes>('/googlePlay/getGooglePlayCfg', {})
    },

    saveGooglePlayCfg: (data: SaveGooglePlayCfgReq) => {
        return request.post<SaveGooglePlayCfgRes>('/googlePlay/saveGooglePlayCfg', data)
    },
}

export default googlePlayApi

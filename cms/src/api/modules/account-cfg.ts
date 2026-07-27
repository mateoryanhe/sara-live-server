import request from '../request'
import type {GetAccountCfgRes, SaveAccountCfgReq, SaveAccountCfgRes} from '@/types/api'

export const accountCfgApi = {
    getAccountCfg: () => {
        return request.post<GetAccountCfgRes>('/accountCfg/getAccountCfg', {})
    },

    saveAccountCfg: (data: SaveAccountCfgReq) => {
        return request.post<SaveAccountCfgRes>('/accountCfg/saveAccountCfg', data)
    },
}

export default accountCfgApi

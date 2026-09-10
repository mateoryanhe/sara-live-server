import request from '../request'
import type {GetHaiPayCfgRes, SaveHaiPayCfgReq, SaveHaiPayCfgRes} from '@/types/api'

export const haipayApi = {
    getHaiPayCfg: () => {
        return request.post<GetHaiPayCfgRes>('/haipay/getHaiPayCfg', {})
    },

    saveHaiPayCfg: (data: SaveHaiPayCfgReq) => {
        return request.post<SaveHaiPayCfgRes>('/haipay/saveHaiPayCfg', data)
    },
}

export default haipayApi

import request from '../request'
import type {GetYhPayCfgRes, SaveYhPayCfgReq, SaveYhPayCfgRes} from '@/types/api'

export const yhpayApi = {
    getYhPayCfg: () => {
        return request.post<GetYhPayCfgRes>('/yhpay/getYhPayCfg', {})
    },

    saveYhPayCfg: (data: SaveYhPayCfgReq) => {
        return request.post<SaveYhPayCfgRes>('/yhpay/saveYhPayCfg', data)
    },
}

export default yhpayApi

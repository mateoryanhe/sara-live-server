import request from '../request'
import type {GetCustomerServiceCfgRes, SaveCustomerServiceCfgReq, SaveCustomerServiceCfgRes} from '@/types/api'

export const customerServiceApi = {
    getCustomerServiceCfg: () => {
        return request.post<GetCustomerServiceCfgRes>('/customerService/getCustomerServiceCfg', {})
    },

    saveCustomerServiceCfg: (data: SaveCustomerServiceCfgReq) => {
        return request.post<SaveCustomerServiceCfgRes>('/customerService/saveCustomerServiceCfg', data)
    },
}

export default customerServiceApi

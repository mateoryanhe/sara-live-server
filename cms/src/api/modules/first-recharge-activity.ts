import request from '../request'
import type {
    FirstRechargeActivityCfg,
    GetFirstRechargeActivityCfgRes,
    SaveFirstRechargeActivityCfgReq,
    SaveFirstRechargeActivityCfgRes,
} from '@/types/api'

export const firstRechargeActivityApi = {
    getFirstRechargeActivityCfg: () => {
        return request.post<GetFirstRechargeActivityCfgRes>('/firstRechargeActivity/getFirstRechargeActivityCfg', {})
    },

    saveFirstRechargeActivityCfg: (data: SaveFirstRechargeActivityCfgReq) => {
        return request.post<SaveFirstRechargeActivityCfgRes>('/firstRechargeActivity/saveFirstRechargeActivityCfg', data)
    },
}

export default firstRechargeActivityApi

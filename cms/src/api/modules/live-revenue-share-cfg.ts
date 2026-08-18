import request from '../request'
import type {
  GetLiveRevenueShareCfgRes,
  SaveLiveRevenueShareCfgReq,
  SaveLiveRevenueShareCfgRes,
} from '@/types/api'

export const liveRevenueShareCfgApi = {
  getCfg: () => {
    return request.post<GetLiveRevenueShareCfgRes>('/liveRevenueShareCfg/getLiveRevenueShareCfg', {})
  },

  saveCfg: (data: SaveLiveRevenueShareCfgReq) => {
    return request.post<SaveLiveRevenueShareCfgRes>('/liveRevenueShareCfg/saveLiveRevenueShareCfg', data)
  },
}

export default liveRevenueShareCfgApi

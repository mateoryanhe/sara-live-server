import {request} from '../request'
import type {GetGameActivityCfgRes, SaveGameActivityCfgReq, SaveGameActivityCfgRes} from '@/types/api'

export const gameActivityApi = {
  getGameActivityCfg: () => {
    return request.post<GetGameActivityCfgRes>('/gameActivity/getGameActivityCfg', {})
  },

  saveGameActivityCfg: (data: SaveGameActivityCfgReq) => {
    return request.post<SaveGameActivityCfgRes>('/gameActivity/saveGameActivityCfg', data)
  },
}

export default gameActivityApi

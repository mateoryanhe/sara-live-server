import {request} from '../request'
import type {GetFirebasePushCfgRes, SaveFirebasePushCfgReq, SaveFirebasePushCfgRes} from '@/types/api'

export const firebasePushApi = {
  getFirebasePushCfg: () => {
    return request.post<GetFirebasePushCfgRes>('/firebasePush/getFirebasePushCfg', {})
  },

  saveFirebasePushCfg: (data: SaveFirebasePushCfgReq) => {
    return request.post<SaveFirebasePushCfgRes>('/firebasePush/saveFirebasePushCfg', data)
  },
}

export default firebasePushApi

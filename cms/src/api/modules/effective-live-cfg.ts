import {request} from '../request'
import type {
  GetEffectiveLiveCfgRes,
  SaveEffectiveLiveCfgReq,
  SaveEffectiveLiveCfgRes,
} from '@/types/api'

export const effectiveLiveCfgApi = {
  getEffectiveLiveCfg: () => request.post<GetEffectiveLiveCfgRes>('/effectiveLiveCfg/getEffectiveLiveCfg', {}),
  saveEffectiveLiveCfg: (data: SaveEffectiveLiveCfgReq) => request.post<SaveEffectiveLiveCfgRes>('/effectiveLiveCfg/saveEffectiveLiveCfg', data),
}

export default effectiveLiveCfgApi

import {request} from '../request'
import type {
  GetPlatformAnchorSettlementCfgRes,
  SavePlatformAnchorSettlementCfgReq,
  SavePlatformAnchorSettlementCfgRes,
} from '@/types/api'

export const platformAnchorSettlementCfgApi = {
  getPlatformAnchorSettlementCfg: () => request.post<GetPlatformAnchorSettlementCfgRes>('/platformAnchorSettlementCfg/getPlatformAnchorSettlementCfg', {}),
  savePlatformAnchorSettlementCfg: (data: SavePlatformAnchorSettlementCfgReq) => request.post<SavePlatformAnchorSettlementCfgRes>('/platformAnchorSettlementCfg/savePlatformAnchorSettlementCfg', data),
}

export default platformAnchorSettlementCfgApi

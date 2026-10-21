import {request} from '../request'
import type {
  GetLuckyGiftSettlementCfgRes,
  SaveLuckyGiftSettlementCfgReq,
  SaveLuckyGiftSettlementCfgRes,
} from '@/types/api'

export const luckyGiftSettlementCfgApi = {
  getLuckyGiftSettlementCfg: () =>
      request.post<GetLuckyGiftSettlementCfgRes>('/luckyGiftSettlementCfg/getLuckyGiftSettlementCfg', {}),
  saveLuckyGiftSettlementCfg: (data: SaveLuckyGiftSettlementCfgReq) =>
      request.post<SaveLuckyGiftSettlementCfgRes>('/luckyGiftSettlementCfg/saveLuckyGiftSettlementCfg', data),
}

export default luckyGiftSettlementCfgApi

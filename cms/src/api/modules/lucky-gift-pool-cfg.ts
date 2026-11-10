import {request} from '../request'
import type {
  GetLuckyGiftPoolCfgRes,
  SaveLuckyGiftPoolCfgReq,
  SaveLuckyGiftPoolCfgRes,
} from '@/types/api'

export const luckyGiftPoolCfgApi = {
  getLuckyGiftPoolCfg: () =>
    request.post<GetLuckyGiftPoolCfgRes>('/luckyGiftPoolCfg/getLuckyGiftPoolCfg', {}),
  saveLuckyGiftPoolCfg: (data: SaveLuckyGiftPoolCfgReq) =>
    request.post<SaveLuckyGiftPoolCfgRes>('/luckyGiftPoolCfg/saveLuckyGiftPoolCfg', data),
}

export default luckyGiftPoolCfgApi

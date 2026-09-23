import {request} from '../request'
import type {
  GetAnchorNoSalaryShareCfgRes,
  SaveAnchorNoSalaryShareCfgReq,
  SaveAnchorNoSalaryShareCfgRes,
} from '@/types/api'

export const anchorNoSalaryShareCfgApi = {
  getCfg: () => {
    return request.post<GetAnchorNoSalaryShareCfgRes>(
      '/anchorNoSalaryShareCfg/getAnchorNoSalaryShareCfg',
      {},
    )
  },

  saveCfg: (data: SaveAnchorNoSalaryShareCfgReq) => {
    return request.post<SaveAnchorNoSalaryShareCfgRes>(
      '/anchorNoSalaryShareCfg/saveAnchorNoSalaryShareCfg',
      data,
    )
  },
}

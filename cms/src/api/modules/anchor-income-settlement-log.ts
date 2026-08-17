import request from '../request'
import type {AnchorIncomeSettlementLogItem, AnchorIncomeSettlementLogQuery, PageResponse} from '@/types/api'

export const anchorIncomeSettlementLogApi = {
  getList: (params: AnchorIncomeSettlementLogQuery) => {
    return request.post<PageResponse<AnchorIncomeSettlementLogItem>>('/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList', params)
  },
}

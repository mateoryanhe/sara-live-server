import {request} from '../request'
import type {AnchorIncomeSettlementLogItem, AnchorIncomeSettlementLogQuery, PageResponse} from '@/types/api'

export const anchorIncomeSettlementLogApi = {
  getList: (params: AnchorIncomeSettlementLogQuery) => {
    return request.post<PageResponse<AnchorIncomeSettlementLogItem>>('/anchorIncomeSettlementLog/cmsAnchorIncomeSettlementLogList', params)
  },
  batchApprove: (data: {ids: string[]}) => {
    return request.post<{successCount: number; failCount: number}>(
        '/anchorIncomeSettlementLog/cmsBatchApproveAnchorSettlement',
        data,
    )
  },
  batchTransfer: (data: {ids: string[]}) => {
    return request.post<{successCount: number; failCount: number; message: string}>(
        '/anchorIncomeSettlementLog/cmsBatchTransferAnchorSettlement',
        data,
    )
  },
}

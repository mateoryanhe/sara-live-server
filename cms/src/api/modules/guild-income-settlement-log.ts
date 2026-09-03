import request from '../request'
import type {
  GuildIncomeSettlementLogItem,
  GuildIncomeSettlementLogQuery,
  PageResponse,
} from '@/types/api'

export const guildIncomeSettlementLogApi = {
  getList: (params: GuildIncomeSettlementLogQuery) => {
    return request.post<PageResponse<GuildIncomeSettlementLogItem>>(
        '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList',
        params,
    )
  },
  batchApprove: (data: { ids: string[] }) => {
    return request.post<{ successCount: number; failCount: number }>(
        '/guildIncomeSettlementLog/cmsBatchApproveGuildSettlement',
        data,
    )
  },
  batchTransfer: (data: { ids: string[] }) => {
    return request.post<{ reserved: boolean; message: string }>(
        '/guildIncomeSettlementLog/cmsBatchTransferGuildSettlement',
        data,
    )
  },
}

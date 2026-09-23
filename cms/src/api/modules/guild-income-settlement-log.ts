import {request} from '../request'
import type {
  GuildIncomeSettlementLogItem,
  GuildIncomeSettlementLogDetailRes,
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
  getDetail: (data: { id: string }) => {
    return request.post<GuildIncomeSettlementLogDetailRes>(
        '/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogDetail',
        data,
    )
  },
  reopenApproval: (data: { id: string }) => {
    return request.post<{ success: boolean }>(
        '/guildIncomeSettlementLog/cmsReopenGuildSettlementApproval',
        data,
    )
  },
  copyPayout: (data: { id: string }) => {
    return request.post<{ id: string }>(
        '/guildIncomeSettlementLog/cmsCopyGuildSettlementPayout',
        data,
    )
  },
  updateReceivableUsd: (data: { id: string; settlementReceivableUsd: number }) => {
    return request.post<{ success: boolean }>(
        '/guildIncomeSettlementLog/cmsUpdateGuildSettlementReceivableUsd',
        data,
    )
  },
  batchTransfer: (data: { ids: string[] }) => {
    return request.post<{ successCount: number; failCount: number; message: string }>(
        '/guildIncomeSettlementLog/cmsBatchTransferGuildSettlement',
        data,
    )
  },
}

import request from '../request'
import type {GuildIncomeSettlementLogItem, GuildIncomeSettlementLogQuery, PageResponse} from '@/types/api'

export const guildIncomeSettlementLogApi = {
  getList: (params: GuildIncomeSettlementLogQuery) => {
    return request.post<PageResponse<GuildIncomeSettlementLogItem>>('/guildIncomeSettlementLog/cmsGuildIncomeSettlementLogList', params)
  },
}

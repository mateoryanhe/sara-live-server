import {request} from '../request'
import type {
    AccessLogItem,
    AccessLogQuery,
    AccessLogStats,
    AccessTrendData,
    AccessTrendQuery,
    DetailLogItem,
    DetailLogQuery,
    ErrorLogItem,
    ErrorLogQuery,
    LogPathsConfig,
    PageResponse,
    TraceLogDetail,
} from '@/types/api'

export const logQueryApi = {
    getLogPaths: () => {
        return request.post<LogPathsConfig>('/logQuery/getLogPaths', {})
    },

    queryDetailLogs: (params: DetailLogQuery) => {
        return request.post<PageResponse<DetailLogItem>>('/logQuery/queryDetailLogs', params)
    },

    queryAccessLogs: (params: AccessLogQuery) => {
        return request.post<PageResponse<AccessLogItem>>('/logQuery/queryAccessLogs', params)
    },

    queryErrorLogs: (params: ErrorLogQuery) => {
        return request.post<PageResponse<ErrorLogItem>>('/logQuery/queryErrorLogs', params)
    },

    getTraceLogs: (traceId: string, date: string) => {
        return request.post<TraceLogDetail>('/logQuery/getTraceLogs', {traceId, date})
    },

    getAccessStats: (params: { startDate: string; endDate: string; topN?: number }) => {
        return request.post<AccessLogStats>('/logQuery/getAccessStats', params)
    },

    getAccessTrend: (params: AccessTrendQuery) => {
        return request.post<AccessTrendData>('/logQuery/getAccessTrend', params)
    },
}

export default logQueryApi

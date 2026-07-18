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
    LogQueryJobResult,
    LogQueryJobSubmitResult,
    PageResponse,
    TraceLogDetail,
} from '@/types/api'

const LOG_QUERY_POLL_INTERVAL_MS = 2000
const LOG_QUERY_POLL_TIMEOUT_MS = 90000

const sleep = (ms: number) => new Promise<void>((resolve) => {
    setTimeout(resolve, ms)
})

export type LogQueryStatusHandler = (job: LogQueryJobResult) => void

export async function runAsyncLogQuery<T>(
    queryType: string,
    payload: object,
    onStatus?: LogQueryStatusHandler,
): Promise<T> {
    const submit = await request.post<LogQueryJobSubmitResult>('/logQuery/submitJob', {
        queryType,
        payload,
    })

    const deadline = Date.now() + LOG_QUERY_POLL_TIMEOUT_MS
    while (Date.now() < deadline) {
        await sleep(LOG_QUERY_POLL_INTERVAL_MS)
        const job = await request.post<LogQueryJobResult<T>>('/logQuery/getJob', {jobId: submit.jobId})
        onStatus?.(job)
        if (job.status === 'done') {
            return job.result as T
        }
        if (job.status === 'failed') {
            throw new Error(job.errorMessage || '日志查询失败')
        }
    }
    throw new Error('日志查询超时，请稍后重试')
}

export const logQueryApi = {
    getLogPaths: () => {
        return request.post<LogPathsConfig>('/logQuery/getLogPaths', {})
    },

    queryDetailLogs: (params: DetailLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<PageResponse<DetailLogItem>>('detail', params, onStatus)
    },

    queryAccessLogs: (params: AccessLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<PageResponse<AccessLogItem>>('access', params, onStatus)
    },

    queryErrorLogs: (params: ErrorLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<PageResponse<ErrorLogItem>>('error', params, onStatus)
    },

    getTraceLogs: (traceId: string, startDate: string, endDate: string, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<TraceLogDetail>('trace', {traceId, startDate, endDate}, onStatus)
    },

    getAccessStats: (params: { startDate: string; endDate: string; topN?: number }, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<AccessLogStats>('accessStats', params, onStatus)
    },

    getAccessTrend: (params: AccessTrendQuery, onStatus?: LogQueryStatusHandler) => {
        return runAsyncLogQuery<AccessTrendData>('accessTrend', params, onStatus)
    },
}

export default logQueryApi

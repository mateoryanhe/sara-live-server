import {request} from '../request'
import {
    buildPageResponse,
    estimateLogQueryPageTotal,
    parseAccessLogLine,
    parseDetailLogLine,
    parseErrorLogLines,
    parseLogExportPage,
    parseStatsExport,
    parseTraceExport,
    parseTrendExport,
} from '@/utils/logParsers'
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
    LogQueryExportResult,
    LogQueryJobResult,
    LogQueryJobSubmitResult,
    PageResponse,
    TraceLogDetail,
} from '@/types/api'

const LOG_QUERY_POLL_INTERVAL_MS = 2000
const LOG_QUERY_POLL_TIMEOUT_MS = 30 * 60 * 1000

const sleep = (ms: number) => new Promise<void>((resolve) => {
    setTimeout(resolve, ms)
})

export type LogQueryStatusHandler = (job: LogQueryJobResult) => void

export async function runAsyncLogQuery<T>(
    queryType: string,
    payload: object,
    onStatus?: LogQueryStatusHandler,
    pollTimeoutMs: number = LOG_QUERY_POLL_TIMEOUT_MS,
): Promise<T> {
    const submit = await request.post<LogQueryJobSubmitResult>('/logQuery/submitJob', {
        queryType,
        payload,
    })

    const deadline = Date.now() + pollTimeoutMs
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

const resolveExportUrl = (fileUrl: string) => {
    if (/^https?:\/\//i.test(fileUrl)) {
        return fileUrl
    }
    const origin = window.location.origin.replace(/\/$/, '')
    return `${origin}${fileUrl.startsWith('/') ? fileUrl : `/${fileUrl}`}`
}

const downloadExportFile = async (exportRes: LogQueryExportResult) => {
    const response = await fetch(resolveExportUrl(exportRes.fileUrl))
    if (!response.ok) {
        throw new Error('下载日志导出文件失败')
    }
    return response.text()
}

const deleteExportFile = async (exportId: string) => {
    await request.post('/logQuery/deleteExport', {exportId})
}

async function runExportQuery<T>(
    queryType: string,
    payload: object,
    parse: (text: string, exportRes: LogQueryExportResult) => T,
    onStatus?: LogQueryStatusHandler,
    pollTimeoutMs: number = LOG_QUERY_POLL_TIMEOUT_MS,
): Promise<T> {
    const exportRes = await runAsyncLogQuery<LogQueryExportResult>(queryType, payload, onStatus, pollTimeoutMs)
    try {
        const text = await downloadExportFile(exportRes)
        return parse(text, exportRes)
    } finally {
        try {
            await deleteExportFile(exportRes.exportId)
        } catch (error) {
            console.warn('删除日志导出文件失败:', error)
        }
    }
}

export const logQueryApi = {
    getLogPaths: () => {
        return request.post<LogPathsConfig>('/logQuery/getLogPaths', {})
    },

    queryDetailLogs: (params: DetailLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<PageResponse<DetailLogItem>>('detail', params, (text, exportRes) => {
            const items = parseLogExportPage(text, parseDetailLogLine)
            return buildPageResponse(items, estimateLogQueryPageTotal(exportRes.pageIndex, exportRes.pageSize, items.length))
        }, onStatus)
    },

    queryAccessLogs: (params: AccessLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<PageResponse<AccessLogItem>>('access', params, (text, exportRes) => {
            const items = parseLogExportPage(text, parseAccessLogLine)
            return buildPageResponse(items, estimateLogQueryPageTotal(exportRes.pageIndex, exportRes.pageSize, items.length))
        }, onStatus)
    },

    queryErrorLogs: (params: ErrorLogQuery, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<PageResponse<ErrorLogItem>>('error', params, (text, exportRes) => {
            const items = parseErrorLogLines(text)
            return buildPageResponse(items, estimateLogQueryPageTotal(exportRes.pageIndex, exportRes.pageSize, items.length))
        }, onStatus)
    },

    getTraceLogs: (traceId: string, startDate: string, endDate: string, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<TraceLogDetail>('trace', {traceId, startDate, endDate}, (text) =>
            parseTraceExport(text, traceId, startDate, endDate), onStatus)
    },

    getAccessStats: (params: { startDate: string; endDate: string; topN?: number }, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<AccessLogStats>('accessStats', params, (text) => parseStatsExport(text), onStatus)
    },

    getAccessTrend: (params: AccessTrendQuery, onStatus?: LogQueryStatusHandler) => {
        return runExportQuery<AccessTrendData>('accessTrend', params, (text) => parseTrendExport(text), onStatus)
    },

    deleteExport: deleteExportFile,
}

export default logQueryApi

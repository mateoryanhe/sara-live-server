import {cmsExportApi} from '@/api/modules/cmsExport'
import type {CMSExportJobResult, CMSExportResult} from '@/types/api'

const CMS_EXPORT_POLL_INTERVAL_MS = 2000
const CMS_EXPORT_POLL_TIMEOUT_MS = 30 * 60 * 1000

const sleep = (ms: number) => new Promise<void>((resolve) => {
    setTimeout(resolve, ms)
})

export type CMSExportStatusHandler = (job: CMSExportJobResult) => void

type TranslateFn = (key: string, params?: Record<string, unknown>) => string

export async function runAsyncCmsExport(
    exportType: string,
    payload: object,
    onStatus?: CMSExportStatusHandler,
    pollTimeoutMs: number = CMS_EXPORT_POLL_TIMEOUT_MS,
): Promise<CMSExportResult> {
    const submit = await cmsExportApi.submitJob({exportType, payload})

    const deadline = Date.now() + pollTimeoutMs
    while (Date.now() < deadline) {
        await sleep(CMS_EXPORT_POLL_INTERVAL_MS)
        const job = await cmsExportApi.getJob({jobId: submit.jobId})
        onStatus?.(job)
        if (job.status === 'done') {
            if (!job.result) {
                throw new Error('Export result is empty')
            }
            return job.result
        }
        if (job.status === 'failed') {
            throw new Error(job.errorMessage || 'Export failed')
        }
    }
    throw new Error('Export timeout')
}

const resolveExportUrl = (fileUrl: string) => {
    if (/^https?:\/\//i.test(fileUrl)) {
        return fileUrl
    }
    const origin = window.location.origin.replace(/\/$/, '')
    return `${origin}${fileUrl.startsWith('/') ? fileUrl : `/${fileUrl}`}`
}

export function downloadExportFile(fileUrl: string, fileName: string): void {
    const link = document.createElement('a')
    link.href = resolveExportUrl(fileUrl)
    link.download = fileName.endsWith('.csv') ? fileName : `${fileName}.csv`
    link.click()
}

export async function exportAndDownloadFile(
    exportType: string,
    payload: object,
    fileName: string,
    onStatus?: CMSExportStatusHandler,
): Promise<CMSExportResult> {
    const result = await runAsyncCmsExport(exportType, payload, onStatus)
    try {
        downloadExportFile(result.fileUrl, fileName || result.fileName)
    } finally {
        try {
            await cmsExportApi.deleteExport({exportId: result.exportId})
        } catch (error) {
            console.warn('delete export file failed:', error)
        }
    }
    return result
}

export function buildExportStatusTip(t: TranslateFn, job: CMSExportJobResult): string {
    if (job.status === 'pending') {
        return job.queuePosition > 1
            ? t('common.exportQueuePending', {count: job.queuePosition - 1})
            : t('common.exportQueueStarting')
    }
    if (job.status === 'running') {
        const progress = job.progress
        if (progress?.totalRows) {
            return t('common.exportProgress', {
                exported: progress.exportedRows,
                total: progress.totalRows,
            })
        }
        return t('common.exportRunning')
    }
    return ''
}

export const CMS_EXPORT_TYPE_LIVE_RECORD = 'liveRecord'
export const CMS_EXPORT_TYPE_LIVE_REVENUE_LOG = 'liveRevenueLog'
export const CMS_EXPORT_TYPE_VIDEO_CALL_LOG = 'videoCallLog'
export const CMS_EXPORT_TYPE_ANCHOR_INCOME_SETTLEMENT_LOG = 'anchorIncomeSettlementLog'
export const CMS_EXPORT_TYPE_GUILD_INCOME_SETTLEMENT_LOG = 'guildIncomeSettlementLog'
export const CMS_EXPORT_TYPE_GUILD_ANCHOR_INCOME_SETTLEMENT_LOG = 'guildAnchorIncomeSettlementLog'
export const CMS_EXPORT_TYPE_MY_GUILD_ANCHOR_INCOME_SETTLEMENT_LOG = 'myGuildAnchorIncomeSettlementLog'
export const CMS_EXPORT_TYPE_ANCHOR_DAILY_EFFECTIVE_LIVE = 'anchorDailyEffectiveLive'
export const CMS_EXPORT_TYPE_GUILD_DAILY_EFFECTIVE_LIVE = 'guildDailyEffectiveLive'
export const CMS_EXPORT_TYPE_GUILD_ANCHOR_DAILY_EFFECTIVE_LIVE = 'guildAnchorDailyEffectiveLive'

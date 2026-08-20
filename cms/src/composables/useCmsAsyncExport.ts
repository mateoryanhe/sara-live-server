import {ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {buildExportStatusTip, exportAndDownloadFile} from '@/utils/cms-async-export'

type TranslateFn = (key: string, params?: Record<string, unknown>) => string

export function handleCmsExportError(error: unknown, t: TranslateFn) {
    console.error('cms export failed:', error)
    const message = error instanceof Error ? error.message : ''
    if (message.includes('没有可导出的数据') || message.toLowerCase().includes('no data')) {
        ElMessage.warning(t('common.exportEmpty'))
        return
    }
    if (message.includes('超时') || message.toLowerCase().includes('timeout')) {
        ElMessage.error(t('common.exportTimeout'))
        return
    }
    ElMessage.error(message || t('common.exportFailed'))
}

export function buildCsvHeaders(columns: Array<{header: string}>): string[] {
    return columns.map(column => column.header)
}

export function buildDailyEffectiveLiveExportLabels(t: TranslateFn, ns = 'pages.anchorList') {
    return {
        settledYesText: t(`${ns}.dailySettledYes`),
        settledNoText: t(`${ns}.dailySettledNo`),
    }
}

export function useCmsAsyncExport() {
    const exporting = ref(false)
    const exportStatusTip = ref('')
    const {t} = useI18n()

    const runExport = async (exportType: string, payload: object, fileName: string) => {
        exporting.value = true
        exportStatusTip.value = t('common.exportRunning')
        try {
            await exportAndDownloadFile(exportType, payload, fileName, job => {
                exportStatusTip.value = buildExportStatusTip(t, job)
            })
            ElMessage.success(t('common.exportSuccess'))
        } catch (error) {
            handleCmsExportError(error, t)
        } finally {
            exporting.value = false
            exportStatusTip.value = ''
        }
    }

    return {exporting, exportStatusTip, runExport}
}

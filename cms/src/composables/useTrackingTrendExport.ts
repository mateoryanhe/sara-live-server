import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {usePagePermission} from '@/composables/usePagePermission'
import type {TrackingEventTrendPoint} from '@/types/api'
import {
  buildTrackingExportFilename,
  exportTrackingCountTrendCsv,
  type TrackingTrendExportSummaryLabels,
} from '@/utils/tracking-trend-export'

export type TrackingTrendExportContext = {
  fileSlug: string
  eventSlug: string
  period: 'daily' | 'weekly' | 'monthly'
  todayCount: number
  weekCount: number
  monthCount: number
  points: TrackingEventTrendPoint[]
  countColumnHeader: string
  summaryLabels: TrackingTrendExportSummaryLabels
}

export function useTrackingTrendExport(
  pageName: string,
  getContext: () => TrackingTrendExportContext,
) {
  const {t} = useI18n()
  const {can} = usePagePermission(pageName)
  const canExport = computed(() => can('export'))

  const handleExport = () => {
    const ctx = getContext()
    const ok = exportTrackingCountTrendCsv(
        buildTrackingExportFilename(ctx.fileSlug, ctx.eventSlug, ctx.period),
        {
          todayCount: ctx.todayCount,
          weekCount: ctx.weekCount,
          monthCount: ctx.monthCount,
          points: ctx.points,
          countColumnHeader: ctx.countColumnHeader,
          summaryLabels: ctx.summaryLabels,
          timeHeader: t('pages.dashboard.exportCsvTime'),
        },
    )
    if (ok) {
      ElMessage.success(t('common.exportSuccess'))
    } else {
      ElMessage.warning(t('common.exportEmpty'))
    }
  }

  return {canExport, handleExport}
}

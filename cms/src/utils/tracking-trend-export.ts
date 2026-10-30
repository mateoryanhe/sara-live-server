import type {TrackingEventTrendPoint} from '@/types/api'
import {buildCsvContent, type CsvColumn} from './csv-export'

export type TrackingTrendExportSummaryLabels = {
  today: string
  week: string
  month: string
}

export type TrackingTrendExportInput = {
  todayCount: number
  weekCount: number
  monthCount: number
  points: TrackingEventTrendPoint[]
  countColumnHeader: string
  summaryLabels: TrackingTrendExportSummaryLabels
  timeHeader: string
}

type TrackingTrendExportRow = {
  time: string
  count: number
}

export function buildTrackingExportFilename(
  fileSlug: string,
  eventSlug: string,
  period: string,
  stamp?: string,
): string {
  const ts = stamp || new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')
  return `tracking-${fileSlug}-${eventSlug}-${period}-${ts}`
}

export function pickTrackingTrendPeriodPoints(
  trend: {
    daily?: TrackingEventTrendPoint[]
    weekly?: TrackingEventTrendPoint[]
    monthly?: TrackingEventTrendPoint[]
  },
  period: 'daily' | 'weekly' | 'monthly',
): TrackingEventTrendPoint[] {
  if (period === 'weekly') {
    return trend.weekly || []
  }
  if (period === 'monthly') {
    return trend.monthly || []
  }
  return trend.daily || []
}

export function buildTrackingCountTrendCsvContent(input: TrackingTrendExportInput): string | null {
  const points = Array.isArray(input.points) ? input.points : []
  const hasSummary =
      input.todayCount !== 0 || input.weekCount !== 0 || input.monthCount !== 0
  if (points.length === 0 && !hasSummary) {
    return null
  }

  const rows: TrackingTrendExportRow[] = [
    {time: input.summaryLabels.today, count: input.todayCount},
    {time: input.summaryLabels.week, count: input.weekCount},
    {time: input.summaryLabels.month, count: input.monthCount},
    ...points.map((point) => ({
      time: point.time,
      count: Number(point.count) || 0,
    })),
  ]

  const columns: CsvColumn<TrackingTrendExportRow>[] = [
    {header: input.timeHeader, value: (row) => row.time},
    {header: input.countColumnHeader, value: (row) => row.count},
  ]
  return buildCsvContent(columns, rows)
}

export function exportTrackingCountTrendCsv(filename: string, input: TrackingTrendExportInput): boolean {
  const content = buildTrackingCountTrendCsvContent(input)
  if (!content) {
    return false
  }
  const blob = new Blob([content], {type: 'text/csv;charset=utf-8;'})
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename.endsWith('.csv') ? filename : `${filename}.csv`
  link.click()
  URL.revokeObjectURL(url)
  return true
}

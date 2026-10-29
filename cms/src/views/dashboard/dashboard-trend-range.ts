import {formatServerDateStringFromDate} from '@/utils/server-datetime'

/** 资源指标(在线/游戏统计)服务端细采样保留天数 */
export const DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS = 3

/** 仪表盘折线：统一快捷粒度与选手动范围上限（自然日） */
export const DASHBOARD_TREND_PRESET_DAYS = [1, 7, 30] as const
export const DASHBOARD_TREND_MAX_CALENDAR_DAYS = 30

/** @deprecated 使用 DASHBOARD_TREND_MAX_CALENDAR_DAYS */
export const DASHBOARD_STAT_TREND_MAX_CALENDAR_DAYS = DASHBOARD_TREND_MAX_CALENDAR_DAYS

export const DASHBOARD_ONLINE_TREND_MAX_CALENDAR_DAYS = DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS

/** 与后端 resourcemonitor maxTrendPoints 一致 */
export const DASHBOARD_ONLINE_TREND_POINT_LIMIT = 26_000

/** 服务端单次拉取上限；前端用图表滑块浏览全量点 */
export function diffDashboardCalendarDays(start: string, end: string) {
  const startMs = new Date(`${start}T00:00:00`).getTime()
  const endMs = new Date(`${end}T00:00:00`).getTime()
  if (Number.isNaN(startMs) || Number.isNaN(endMs)) {
    return 0
  }
  return Math.floor(Math.abs(endMs - startMs) / 86_400_000) + 1
}

export function buildDashboardTrendPresetRange(days: number): [string, string] {
  const capped = Math.min(Math.max(days, 1), DASHBOARD_TREND_MAX_CALENDAR_DAYS)
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - (capped - 1))
  return [
    formatServerDateStringFromDate(start),
    formatServerDateStringFromDate(end),
  ]
}

export function buildDashboardDefaultTrendRange(): [string, string] {
  return buildDashboardTrendPresetRange(1)
}

/** 判断当前日期范围是否等于某快捷天数（用于按钮高亮） */
export function matchDashboardTrendPresetDays(range: [string, string] | null): number | null {
  if (!range?.[0] || !range?.[1]) {
    return null
  }
  for (const days of DASHBOARD_TREND_PRESET_DAYS) {
    const preset = buildDashboardTrendPresetRange(days)
    if (preset[0] === range[0] && preset[1] === range[1]) {
      return days
    }
  }
  return null
}

/**
 * 最近 N 天资源指标趋势（细采样 10 秒，与后端保留策略一致）。
 * 不传 start/end：走服务端默认 FineRetention + 10 分钟分桶缓存。
 * 切勿每次轮询带「当前秒」的滑动窗口——后端曾按完整时间串永久缓存，会导致数 GB 泄漏。
 */
export function buildDashboardResourceMetricTrendQuery(
    _days: number = DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS,
    bypassDashboardCache = false,
) {
  return {
    startTime: '',
    endTime: '',
    limit: DASHBOARD_ONLINE_TREND_POINT_LIMIT,
    bypassDashboardCache,
  }
}

/** 将日期范围截断为最多 maxDays 个自然日（保留结束日及向前 maxDays-1 天） */
export function clampDashboardTrendDateRange(
    range: [string, string] | null,
    maxDays: number,
): {range: [string, string] | null; clamped: boolean} {
  if (!range?.[0] || !range?.[1] || maxDays <= 0) {
    return {range, clamped: false}
  }
  const spanDays = diffDashboardCalendarDays(range[0], range[1])
  if (spanDays <= maxDays) {
    return {range, clamped: false}
  }
  const end = new Date(`${range[1]}T00:00:00`)
  const start = new Date(end)
  start.setDate(end.getDate() - (maxDays - 1))
  return {
    range: [
      formatServerDateStringFromDate(start),
      formatServerDateStringFromDate(end),
    ],
    clamped: true,
  }
}

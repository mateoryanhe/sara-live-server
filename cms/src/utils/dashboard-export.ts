import type {
  GameBetUserTrendPoint,
  MinuteLiveActivityTrendPoint,
  ResourceMetricPoint,
  UserStatTrendPoint,
} from '@/types/api'
import {buildCsvContent, type CsvColumn} from './csv-export'

export type DashboardExportLabels = {
  section: string
  metric: string
  value: string
  time: string
  audienceOnline: string
  liveAnchorOnline: string
  gameOnline: string
  gameBetUser: string
  activeUsers: string
  registerUsers: string
  rechargeUser: string
  appRecharge: string
  h5Recharge: string
  appLogin: string
  h5Login: string
  join: string
  leave: string
  gift: string
  chat: string
  bet: string
  call1v1: string
  rechargeOrder: string
}

type BasicExportRow = {section: string; metric: string; value: string | number}

const downloadCsvFile = (filename: string, content: string): void => {
  const blob = new Blob([content], {type: 'text/csv;charset=utf-8;'})
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename.endsWith('.csv') ? filename : `${filename}.csv`
  link.click()
  URL.revokeObjectURL(url)
}

const sleep = (ms: number) => new Promise<void>((resolve) => {
  setTimeout(resolve, ms)
})

export function buildDashboardExportStamp(): string {
  return new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')
}

export function buildDashboardExportFilename(tabKey: string, suffix?: string, stamp?: string): string {
  const ts = stamp || buildDashboardExportStamp()
  const part = suffix ? `${tabKey}-${suffix}` : tabKey
  return `dashboard-${part}-${ts}`
}

export function buildDashboardBasicCsvContent(
  sections: {title: string; cards: {label: string; value: number | string}[]}[],
  labels: DashboardExportLabels,
): string | null {
  const rows: BasicExportRow[] = []
  for (const section of sections) {
    for (const card of section.cards) {
      rows.push({
        section: section.title,
        metric: card.label,
        value: card.value,
      })
    }
  }
  if (rows.length === 0) {
    return null
  }
  const columns: CsvColumn<BasicExportRow>[] = [
    {header: labels.section, value: (r) => r.section},
    {header: labels.metric, value: (r) => r.metric},
    {header: labels.value, value: (r) => r.value},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardBasicCsv(
  filename: string,
  sections: {title: string; cards: {label: string; value: number | string}[]}[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardBasicCsvContent(sections, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

const resolveResourceMetricAudience = (point: ResourceMetricPoint): number => {
  if (point.audienceCount !== undefined && point.audienceCount !== '') {
    return Number(point.audienceCount)
  }
  const total = Number(point.onlineCount || 0)
  const anchors = Number(point.liveAnchorCount || 0)
  return total > anchors ? total - anchors : 0
}

type OnlineTrendRow = {
  time: string
  audience: number
  liveAnchor: number
}

export function buildDashboardOnlineCsvContent(
  realtime: {audienceOnlineCount: number; liveAnchorCount: number},
  points: ResourceMetricPoint[],
  labels: DashboardExportLabels,
): string | null {
  const trendRows: OnlineTrendRow[] = points.map((p) => ({
    time: p.time,
    audience: resolveResourceMetricAudience(p),
    liveAnchor: Number(p.liveAnchorCount || 0),
  }))
  if (trendRows.length === 0 && realtime.audienceOnlineCount === 0 && realtime.liveAnchorCount === 0) {
    return null
  }

  type Row = {time: string; audience: number | string; liveAnchor: number | string}
  const rows: Row[] = [
    {time: '', audience: realtime.audienceOnlineCount, liveAnchor: realtime.liveAnchorCount},
    ...trendRows.map((r) => ({
      time: r.time,
      audience: r.audience,
      liveAnchor: r.liveAnchor,
    })),
  ]

  const columns: CsvColumn<Row>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.audienceOnline, value: (r) => r.audience},
    {header: labels.liveAnchorOnline, value: (r) => r.liveAnchor},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardOnlineCsv(
  filename: string,
  realtime: {audienceOnlineCount: number; liveAnchorCount: number},
  points: ResourceMetricPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardOnlineCsvContent(realtime, points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardGameStatsCsvContent(
  realtime: {gameOnlineCount: number; gameBetUserCount: number},
  points: GameBetUserTrendPoint[],
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0 && realtime.gameOnlineCount === 0 && realtime.gameBetUserCount === 0) {
    return null
  }

  type Row = {time: string; gameOnline: number | string; gameBetUser: number | string}
  const rows: Row[] = [
    {time: 'realtime', gameOnline: realtime.gameOnlineCount, gameBetUser: realtime.gameBetUserCount},
    ...points.map((p) => ({
      time: p.time,
      gameOnline: Number(p.onlineCount ?? 0),
      gameBetUser: Number(p.count ?? 0),
    })),
  ]

  const columns: CsvColumn<Row>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.gameOnline, value: (r) => r.gameOnline},
    {header: labels.gameBetUser, value: (r) => r.gameBetUser},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardGameStatsCsv(
  filename: string,
  realtime: {gameOnlineCount: number; gameBetUserCount: number},
  points: GameBetUserTrendPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardGameStatsCsvContent(realtime, points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardUserEventsCsvContent(
  points: MinuteLiveActivityTrendPoint[],
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0) {
    return null
  }
  const columns: CsvColumn<MinuteLiveActivityTrendPoint>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.join, value: (r) => r.joinCount},
    {header: labels.leave, value: (r) => r.leaveCount},
    {header: labels.gift, value: (r) => r.giftCount},
    {header: labels.chat, value: (r) => r.chatCount},
    {header: labels.bet, value: (r) => r.betCount},
    {header: labels.call1v1, value: (r) => r.call1v1Count},
    {header: labels.rechargeOrder, value: (r) => r.rechargeOrderCount},
  ]
  return buildCsvContent(columns, points)
}

export function exportDashboardUserEventsCsv(
  filename: string,
  points: MinuteLiveActivityTrendPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardUserEventsCsvContent(points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardUserTrendCsvContent(
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0) {
    return null
  }
  const columns: CsvColumn<UserStatTrendPoint>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.activeUsers, value: (r) => r.activeUserCount},
    {header: labels.registerUsers, value: (r) => r.registerUserCount},
  ]
  return buildCsvContent(columns, points)
}

export function exportDashboardUserTrendCsv(
  filename: string,
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardUserTrendCsvContent(points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardRechargeUserBarCsvContent(
  points: UserStatTrendPoint[],
  metricKey: string,
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0) {
    return null
  }
  type Row = {time: string; value: number}
  const rows: Row[] = points.map((p) => ({
    time: p.time,
    value: Number(p.barMetrics?.[metricKey] ?? 0),
  }))
  const columns: CsvColumn<Row>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.rechargeUser, value: (r) => r.value},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardRechargeUserBarCsv(
  filename: string,
  points: UserStatTrendPoint[],
  metricKey: string,
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardRechargeUserBarCsvContent(points, metricKey, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardClientRechargeCsvContent(
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0) {
    return null
  }
  type Row = {time: string; app: number; h5: number}
  const rows: Row[] = points.map((p) => ({
    time: p.time,
    app: Number(p.appDeviceRechargeAmount || 0),
    h5: Number(p.h5DeviceRechargeAmount || 0),
  }))
  const columns: CsvColumn<Row>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.appRecharge, value: (r) => r.app},
    {header: labels.h5Recharge, value: (r) => r.h5},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardClientRechargeCsv(
  filename: string,
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardClientRechargeCsvContent(points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export function buildDashboardClientLoginCsvContent(
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): string | null {
  if (points.length === 0) {
    return null
  }
  type Row = {time: string; app: number; h5: number}
  const rows: Row[] = points.map((p) => ({
    time: p.time,
    app: Number(p.appDeviceLoginCount || 0),
    h5: Number(p.h5DeviceLoginCount || 0),
  }))
  const columns: CsvColumn<Row>[] = [
    {header: labels.time, value: (r) => r.time},
    {header: labels.appLogin, value: (r) => r.app},
    {header: labels.h5Login, value: (r) => r.h5},
  ]
  return buildCsvContent(columns, rows)
}

export function exportDashboardClientLoginCsv(
  filename: string,
  points: UserStatTrendPoint[],
  labels: DashboardExportLabels,
): boolean {
  const content = buildDashboardClientLoginCsvContent(points, labels)
  if (!content) {
    return false
  }
  downloadCsvFile(filename, content)
  return true
}

export type DashboardCsvBundleFile = {
  filename: string
  content: string | null
}

/** 依次触发浏览器下载多个 CSV（间隔避免被拦截） */
export async function downloadDashboardCsvBundle(files: DashboardCsvBundleFile[]): Promise<number> {
  const toDownload = files.filter((f): f is DashboardCsvBundleFile & {content: string} => Boolean(f.content))
  for (let i = 0; i < toDownload.length; i++) {
    downloadCsvFile(toDownload[i].filename, toDownload[i].content)
    if (i < toDownload.length - 1) {
      await sleep(180)
    }
  }
  return toDownload.length
}

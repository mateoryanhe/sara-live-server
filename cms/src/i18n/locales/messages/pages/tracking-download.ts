import {definePageMessagesFromEn} from './_define'

const zh = {
  tabH5: 'H5下载量',
  tabOfficial: '官网下载量',
  h5Desc:
      '统计 H5 落地页下载点击：POST /trackingEvent/reportH5Download（免登录，暂不要求 body），每次成功 +1（event h5_download）。',
  officialDesc:
      '统计官网下载点击：POST /trackingEvent/reportOfficialSiteDownload（免登录，暂不要求 body），每次成功 +1（event official_site_download）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  h5SeriesName: 'H5下载次数',
  officialSeriesName: '官网下载次数',
  h5ChartDaily: 'H5下载 · 按日（最近 30 天）',
  h5ChartWeekly: 'H5下载 · 按周（最近 12 周）',
  h5ChartMonthly: 'H5下载 · 按月（最近 12 月）',
  officialChartDaily: '官网下载 · 按日（最近 30 天）',
  officialChartWeekly: '官网下载 · 按周（最近 12 周）',
  officialChartMonthly: '官网下载 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingDownloadMessages = definePageMessagesFromEn(zh, {
  tabH5: 'H5 downloads',
  tabOfficial: 'Official site downloads',
  h5Desc:
      'H5 landing download clicks: POST /trackingEvent/reportH5Download (no login, empty body ok), +1 each (event h5_download).',
  officialDesc:
      'Official site download clicks: POST /trackingEvent/reportOfficialSiteDownload (no login, empty body ok), +1 each (event official_site_download).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  h5SeriesName: 'H5 download count',
  officialSeriesName: 'Official download count',
  h5ChartDaily: 'H5 download · daily (last 30 days)',
  h5ChartWeekly: 'H5 download · weekly (last 12 weeks)',
  h5ChartMonthly: 'H5 download · monthly (last 12 months)',
  officialChartDaily: 'Official download · daily (last 30 days)',
  officialChartWeekly: 'Official download · weekly (last 12 weeks)',
  officialChartMonthly: 'Official download · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

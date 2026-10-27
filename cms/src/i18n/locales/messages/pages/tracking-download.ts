import {definePageMessagesFromEn} from './_define'

const zh = {
  tabH5: 'H5下载量',
  tabOfficial: '官网下载量',
  tabShareShort: '分享短链点击',
  h5Desc:
      '统计 H5 落地页下载点击：POST /trackingEvent/reportH5Download（免登录，暂不要求 body），每次成功 +1（event h5_download）。',
  officialDesc:
      '统计官网下载点击：POST /trackingEvent/reportOfficialSiteDownload（免登录，暂不要求 body），每次成功 +1（event official_site_download）。',
  shareShortDesc:
      '统计分享文案第三行短链点击：GET /dl 成功 302 到官网最大 APK 时服务端 +1（event share_short_link_click）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  h5SeriesName: 'H5下载次数',
  officialSeriesName: '官网下载次数',
  shareShortSeriesName: '短链点击次数',
  h5ChartDaily: 'H5下载 · 按日（最近 30 天）',
  h5ChartWeekly: 'H5下载 · 按周（最近 12 周）',
  h5ChartMonthly: 'H5下载 · 按月（最近 12 月）',
  officialChartDaily: '官网下载 · 按日（最近 30 天）',
  officialChartWeekly: '官网下载 · 按周（最近 12 周）',
  officialChartMonthly: '官网下载 · 按月（最近 12 月）',
  shareShortChartDaily: '短链点击 · 按日（最近 30 天）',
  shareShortChartWeekly: '短链点击 · 按周（最近 12 周）',
  shareShortChartMonthly: '短链点击 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingDownloadMessages = definePageMessagesFromEn(zh, {
  tabH5: 'H5 downloads',
  tabOfficial: 'Official site downloads',
  tabShareShort: 'Share short-link clicks',
  h5Desc:
      'H5 landing download clicks: POST /trackingEvent/reportH5Download (no login, empty body ok), +1 each (event h5_download).',
  officialDesc:
      'Official site download clicks: POST /trackingEvent/reportOfficialSiteDownload (no login, empty body ok), +1 each (event official_site_download).',
  shareShortDesc:
      'Share copy line-3 short link: GET /dl successful 302 to largest official APK, server +1 (event share_short_link_click).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  h5SeriesName: 'H5 download count',
  officialSeriesName: 'Official download count',
  shareShortSeriesName: 'Short-link click count',
  h5ChartDaily: 'H5 download · daily (last 30 days)',
  h5ChartWeekly: 'H5 download · weekly (last 12 weeks)',
  h5ChartMonthly: 'H5 download · monthly (last 12 months)',
  officialChartDaily: 'Official download · daily (last 30 days)',
  officialChartWeekly: 'Official download · weekly (last 12 weeks)',
  officialChartMonthly: 'Official download · monthly (last 12 months)',
  shareShortChartDaily: 'Short-link click · daily (last 30 days)',
  shareShortChartWeekly: 'Short-link click · weekly (last 12 weeks)',
  shareShortChartMonthly: 'Short-link click · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

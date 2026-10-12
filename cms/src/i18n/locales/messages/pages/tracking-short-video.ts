import {definePageMessagesFromEn} from './_define'

const zh = {
  tabPlayStart: '视频起播',
  tabPaywallClick: '付费解锁点击',
  tabPaySuccess: '付费成功',
  tabWatchEnd: '播放结束',
  playStartDesc:
      '统计短视频起播次数：App POST /shortVideo/watchShortVideoStart（Sofie POST /sofie/clip/beginWatch）成功且视频已上架 +1（event short_video_play_start）。',
  paywallClickDesc:
      '统计确认付费解锁：App POST /shortVideo/payShortVideo（Sofie POST /sofie/clip/unlockMedia），未购买过该视频时每次请求 +1（event short_video_paywall_click）；扣款失败仍计点击，已购直接返回不计。',
  paySuccessDesc:
      '统计付费观看扣款成功：同上 payShortVideo 首次扣款成功 +1（event short_video_pay_success）。',
  watchEndDesc:
      '统计单次播放结束/划走：App POST /shortVideo/watchShortVideoEnd（Sofie POST /sofie/clip/finishWatch）成功 +1（event short_video_watch_end）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  playStartSeriesName: '起播次数',
  paywallClickSeriesName: '点击次数',
  paySuccessSeriesName: '付费成功次数',
  watchEndSeriesName: '结束次数',
  playStartChartDaily: '起播 · 按日（最近 30 天）',
  playStartChartWeekly: '起播 · 按周（最近 12 周）',
  playStartChartMonthly: '起播 · 按月（最近 12 月）',
  paywallClickChartDaily: '付费解锁点击 · 按日（最近 30 天）',
  paywallClickChartWeekly: '付费解锁点击 · 按周（最近 12 周）',
  paywallClickChartMonthly: '付费解锁点击 · 按月（最近 12 月）',
  paySuccessChartDaily: '付费成功 · 按日（最近 30 天）',
  paySuccessChartWeekly: '付费成功 · 按周（最近 12 周）',
  paySuccessChartMonthly: '付费成功 · 按月（最近 12 月）',
  watchEndChartDaily: '播放结束 · 按日（最近 30 天）',
  watchEndChartWeekly: '播放结束 · 按周（最近 12 周）',
  watchEndChartMonthly: '播放结束 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingShortVideoMessages = definePageMessagesFromEn(zh, {
  tabPlayStart: 'Play start',
  tabPaywallClick: 'Paywall click',
  tabPaySuccess: 'Pay success',
  tabWatchEnd: 'Watch end',
  playStartDesc:
      'Short video play starts: POST /shortVideo/watchShortVideoStart (Sofie /sofie/clip/beginWatch) when video is on shelf (+1, short_video_play_start).',
  paywallClickDesc:
      'Confirm unlock (paywall): POST /shortVideo/payShortVideo (Sofie /sofie/clip/unlockMedia), +1 per attempt when not yet purchased (short_video_paywall_click); failed deduct still counts; already-paid skip.',
  paySuccessDesc:
      'Paid watch success: same payShortVideo on first successful diamond deduct (+1, short_video_pay_success).',
  watchEndDesc:
      'Watch session end: POST /shortVideo/watchShortVideoEnd (Sofie /sofie/clip/finishWatch) (+1, short_video_watch_end).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  playStartSeriesName: 'Play start count',
  paywallClickSeriesName: 'Paywall click count',
  paySuccessSeriesName: 'Pay success count',
  watchEndSeriesName: 'Watch end count',
  playStartChartDaily: 'Play start · daily (last 30 days)',
  playStartChartWeekly: 'Play start · weekly (last 12 weeks)',
  playStartChartMonthly: 'Play start · monthly (last 12 months)',
  paywallClickChartDaily: 'Paywall click · daily (last 30 days)',
  paywallClickChartWeekly: 'Paywall click · weekly (last 12 weeks)',
  paywallClickChartMonthly: 'Paywall click · monthly (last 12 months)',
  paySuccessChartDaily: 'Pay success · daily (last 30 days)',
  paySuccessChartWeekly: 'Pay success · weekly (last 12 weeks)',
  paySuccessChartMonthly: 'Pay success · monthly (last 12 months)',
  watchEndChartDaily: 'Watch end · daily (last 30 days)',
  watchEndChartWeekly: 'Watch end · weekly (last 12 weeks)',
  watchEndChartMonthly: 'Watch end · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

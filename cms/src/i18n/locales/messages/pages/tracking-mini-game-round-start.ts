import {definePageMessagesFromEn} from './_define'

const zh = {
  tabRoundStart: '点击开始游戏',
  tabRoundResult: '单局游戏结算',
  tabGameExposure: '半屏游戏曝光',
  tabGameWebViewLoad: 'WebView 加载成功',
  desc:
      '统计用户点击「开始游戏」并成功获取启动链接的次数：App 调用 POST /game/appGameStart（Sofie POST /sofie/play/openTitle）且第三方返回有效 link 后 +1（对齐埋点 mini_game_round_start）。',
  roundResultDesc:
      '统计单局游戏结算次数：第三方 vendor transfer 回调扣款/派彩成功且 transaction_id 首次处理完成后 +1（对齐埋点 mini_game_round_result）。',
  gameExposureDesc:
      '统计半屏游戏窗口曝光次数：App 调用 POST /liveRoom/reportMiniGameExposure（Sofie POST /sofie/studio/notifyMiniGameExposure），登录即可、无需 body 参数，每次成功 +1（对齐埋点 mini_game_exposure）。',
  gameWebViewLoadDesc:
      '统计游戏 WebView 成功加载次数：App 在 WebView 加载游戏成功后调用 POST /liveRoom/reportMiniGameWebViewLoadSuccess（Sofie POST /sofie/studio/notifyMiniGameWebViewLoadSuccess），登录即可、无需 body，每次成功 +1（对齐埋点 mini_game_webview_load_success）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  seriesName: '开始游戏次数',
  roundResultSeriesName: '结算次数',
  gameExposureSeriesName: '曝光次数',
  gameWebViewLoadSeriesName: '加载成功次数',
  chartDaily: '按日趋势（最近 30 天）',
  chartWeekly: '按周趋势（最近 12 周）',
  chartMonthly: '按月趋势（最近 12 月）',
  roundResultChartDaily: '单局结算 · 按日（最近 30 天）',
  roundResultChartWeekly: '单局结算 · 按周（最近 12 周）',
  roundResultChartMonthly: '单局结算 · 按月（最近 12 月）',
  gameExposureChartDaily: '半屏曝光 · 按日（最近 30 天）',
  gameExposureChartWeekly: '半屏曝光 · 按周（最近 12 周）',
  gameExposureChartMonthly: '半屏曝光 · 按月（最近 12 月）',
  gameWebViewLoadChartDaily: 'WebView 加载 · 按日（最近 30 天）',
  gameWebViewLoadChartWeekly: 'WebView 加载 · 按周（最近 12 周）',
  gameWebViewLoadChartMonthly: 'WebView 加载 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingMiniGameRoundStartMessages = definePageMessagesFromEn(zh, {
  tabRoundStart: 'Start game',
  tabRoundResult: 'Round settlement',
  tabGameExposure: 'Half-screen exposure',
  tabGameWebViewLoad: 'WebView load success',
  desc:
      'Counts when the user taps start game and the app successfully gets a launch URL: POST /game/appGameStart (Sofie /sofie/play/openTitle) returns a non-empty link (+1, event mini_game_round_start).',
  roundResultDesc:
      'Counts successful single-round game settlements: vendor transfer callback completes wallet update for a new transaction_id (+1, event mini_game_round_result).',
  gameExposureDesc:
      'Half-screen mini-game exposures: POST /liveRoom/reportMiniGameExposure (Sofie /sofie/studio/notifyMiniGameExposure), login only, empty body, +1 per success (event mini_game_exposure).',
  gameWebViewLoadDesc:
      'Game WebView load success: POST /liveRoom/reportMiniGameWebViewLoadSuccess (Sofie /sofie/studio/notifyMiniGameWebViewLoadSuccess) after the game loads in WebView; login only, empty body, +1 (event mini_game_webview_load_success).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  seriesName: 'Start game count',
  roundResultSeriesName: 'Settlement count',
  gameExposureSeriesName: 'Exposure count',
  gameWebViewLoadSeriesName: 'Load success count',
  chartDaily: 'Daily trend (last 30 days)',
  chartWeekly: 'Weekly trend (last 12 weeks)',
  chartMonthly: 'Monthly trend (last 12 months)',
  roundResultChartDaily: 'Settlement · daily (last 30 days)',
  roundResultChartWeekly: 'Settlement · weekly (last 12 weeks)',
  roundResultChartMonthly: 'Settlement · monthly (last 12 months)',
  gameExposureChartDaily: 'Exposure · daily (last 30 days)',
  gameExposureChartWeekly: 'Exposure · weekly (last 12 weeks)',
  gameExposureChartMonthly: 'Exposure · monthly (last 12 months)',
  gameWebViewLoadChartDaily: 'WebView load · daily (last 30 days)',
  gameWebViewLoadChartWeekly: 'WebView load · weekly (last 12 weeks)',
  gameWebViewLoadChartMonthly: 'WebView load · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

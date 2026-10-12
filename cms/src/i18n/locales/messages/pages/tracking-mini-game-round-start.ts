import {definePageMessagesFromEn} from './_define'

const zh = {
  tabGameJoin: '进入游戏直播间',
  tabFirstFrame: '直播首帧画面渲染完成',
  tabRoundStart: '点击开始游戏',
  tabRoundResult: '单局游戏结算',
  tabGameExposure: '半屏游戏曝光',
  tabGameWebViewLoad: 'WebView 加载成功',
  gameJoinDesc:
      '统计观众进入游戏直播类直播间（category=2）的次数；每次 joinRoom 成功 +1，不含主播本人（含重复进入，event game_live_room_join）。',
  firstFrameDesc:
      '统计游戏直播类房间（category=2）首帧画面渲染完成：App POST /liveRoom/reportLiveFirstFrameRendered（Sofie notifyLiveFirstFrameReady）传 roomId，每次成功 +1，不含主播（event game_live_watch_stream）。',
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
  gameJoinSeriesName: '进入次数',
  firstFrameSeriesName: '首帧完成次数',
  seriesName: '开始游戏次数',
  roundResultSeriesName: '结算次数',
  gameExposureSeriesName: '曝光次数',
  gameWebViewLoadSeriesName: '加载成功次数',
  gameJoinChartDaily: '进房 · 按日（最近 30 天）',
  gameJoinChartWeekly: '进房 · 按周（最近 12 周）',
  gameJoinChartMonthly: '进房 · 按月（最近 12 月）',
  firstFrameChartDaily: '首帧完成 · 按日（最近 30 天）',
  firstFrameChartWeekly: '首帧完成 · 按周（最近 12 周）',
  firstFrameChartMonthly: '首帧完成 · 按月（最近 12 月）',
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
  tabGameJoin: 'Enter game live room',
  tabFirstFrame: 'First video frame rendered',
  tabRoundStart: 'Start game',
  tabRoundResult: 'Round settlement',
  tabGameExposure: 'Half-screen exposure',
  tabGameWebViewLoad: 'WebView load success',
  gameJoinDesc:
      'Join count for game live rooms (category=2); each successful joinRoom +1, excludes anchor (includes re-join, event game_live_room_join).',
  firstFrameDesc:
      'First frame rendered in game live rooms (category=2): POST /liveRoom/reportLiveFirstFrameRendered with roomId (+1, excludes anchor, event game_live_watch_stream).',
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
  gameJoinSeriesName: 'Join count',
  gameJoinChartDaily: 'Join · daily (last 30 days)',
  gameJoinChartWeekly: 'Join · weekly (last 12 weeks)',
  gameJoinChartMonthly: 'Join · monthly (last 12 months)',
  firstFrameSeriesName: 'First frame count',
  firstFrameChartDaily: 'First frame · daily (last 30 days)',
  firstFrameChartWeekly: 'First frame · weekly (last 12 weeks)',
  firstFrameChartMonthly: 'First frame · monthly (last 12 months)',
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

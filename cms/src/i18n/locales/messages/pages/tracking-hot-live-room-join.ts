import {definePageMessagesFromEn} from './_define'

const zh = {
  tabShowcaseJoin: '进入秀场直播间',
  tabFirstFrame: '直播首帧画面渲染完成',
  tabShowcaseLeave: '退出直播间',
  tabGameJoin: '进入游戏直播类房间',
  desc: '统计观众进入秀场类型直播间的次数（category=1，每次 joinRoom 成功 +1，不含主播本人；含重复进入）。',
  leaveDesc:
      '统计观众退出秀场直播间（category=1）的次数；每次 leaveRoom 成功 +1，不含主播本人（仅主动退房接口，不含踢出/心跳清理等路径）。',
  gameJoinDesc:
      '统计观众进入游戏直播类直播间（category=2）的次数；每次 joinRoom 成功 +1，不含主播本人（含重复进入）。',
  firstFrameDesc:
      '统计观众在秀场直播间（category=1）首帧画面渲染完成的次数；App 调用 POST /liveRoom/reportLiveFirstFrameRendered 上报，每次成功 +1（不含主播本人；同一会话可多次上报，由客户端控制去重）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  seriesName: '进入次数',
  firstFrameSeriesName: '首帧完成次数',
  leaveSeriesName: '退出次数',
  gameJoinSeriesName: '进入次数',
  chartDaily: '按日趋势（最近 30 天）',
  chartWeekly: '按周趋势（最近 12 周）',
  chartMonthly: '按月趋势（最近 12 月）',
  firstFrameChartDaily: '首帧完成 · 按日（最近 30 天）',
  firstFrameChartWeekly: '首帧完成 · 按周（最近 12 周）',
  firstFrameChartMonthly: '首帧完成 · 按月（最近 12 月）',
  leaveChartDaily: '退出 · 按日（最近 30 天）',
  leaveChartWeekly: '退出 · 按周（最近 12 周）',
  leaveChartMonthly: '退出 · 按月（最近 12 月）',
  gameJoinChartDaily: '游戏直播进房 · 按日（最近 30 天）',
  gameJoinChartWeekly: '游戏直播进房 · 按周（最近 12 周）',
  gameJoinChartMonthly: '游戏直播进房 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingHotLiveRoomJoinMessages = definePageMessagesFromEn(zh, {
  tabShowcaseJoin: 'Showcase room joins',
  tabFirstFrame: 'First video frame rendered',
  tabShowcaseLeave: 'Leave live room',
  tabGameJoin: 'Game live room joins',
  desc: 'Join count for showcase live rooms (category=1; each successful joinRoom +1, excludes anchor; includes re-join).',
  leaveDesc:
      'Leave count for showcase rooms (category=1); each successful leaveRoom +1, excludes anchor (API leave only, not kick/heartbeat cleanup).',
  gameJoinDesc:
      'Join count for game live rooms (category=2); each successful joinRoom +1, excludes anchor (includes re-join).',
  firstFrameDesc:
      'Counts when the first live video frame is rendered in a showcase room (category=1). App reports via POST /liveRoom/reportLiveFirstFrameRendered; each success +1 (excludes anchor; client may dedupe per session).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  seriesName: 'Join count',
  firstFrameSeriesName: 'First frame count',
  leaveSeriesName: 'Leave count',
  gameJoinSeriesName: 'Join count',
  chartDaily: 'Daily trend (last 30 days)',
  firstFrameChartDaily: 'First frame · daily (last 30 days)',
  firstFrameChartWeekly: 'First frame · weekly (last 12 weeks)',
  firstFrameChartMonthly: 'First frame · monthly (last 12 months)',
  leaveChartDaily: 'Leave · daily (last 30 days)',
  leaveChartWeekly: 'Leave · weekly (last 12 weeks)',
  leaveChartMonthly: 'Leave · monthly (last 12 months)',
  gameJoinChartDaily: 'Game live joins · daily (last 30 days)',
  gameJoinChartWeekly: 'Game live joins · weekly (last 12 weeks)',
  gameJoinChartMonthly: 'Game live joins · monthly (last 12 months)',
  chartWeekly: 'Weekly trend (last 12 weeks)',
  chartMonthly: 'Monthly trend (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

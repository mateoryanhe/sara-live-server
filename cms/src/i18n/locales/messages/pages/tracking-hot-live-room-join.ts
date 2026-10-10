import {definePageMessagesFromEn} from './_define'

const zh = {
  tabShowcaseJoin: '进入秀场直播间',
  tabFirstFrame: '直播首帧画面渲染完成',
  tabShowcaseLeave: '退出直播间',
  desc: '统计观众进入秀场类型直播间的次数（category=1，每次 joinRoom 成功 +1，不含主播本人；含重复进入）。',
  leaveDesc:
      '统计观众退出秀场直播间（category=1）的次数；每次 leaveRoom 成功 +1，不含主播本人（仅主动退房接口，不含踢出/心跳清理等路径）。',
  firstFrameDesc:
      '统计观众在秀场直播间（category=1）首帧画面渲染完成的次数；App 调用 POST /liveRoom/reportLiveFirstFrameRendered 上报，每次成功 +1（不含主播本人；同一会话可多次上报，由客户端控制去重）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  seriesName: '进入次数',
  firstFrameSeriesName: '首帧完成次数',
  leaveSeriesName: '退出次数',
  chartDaily: '按日趋势（最近 30 天）',
  chartWeekly: '按周趋势（最近 12 周）',
  chartMonthly: '按月趋势（最近 12 月）',
  firstFrameChartDaily: '首帧完成 · 按日（最近 30 天）',
  firstFrameChartWeekly: '首帧完成 · 按周（最近 12 周）',
  firstFrameChartMonthly: '首帧完成 · 按月（最近 12 月）',
  leaveChartDaily: '退出 · 按日（最近 30 天）',
  leaveChartWeekly: '退出 · 按周（最近 12 周）',
  leaveChartMonthly: '退出 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingHotLiveRoomJoinMessages = definePageMessagesFromEn(zh, {
  tabShowcaseJoin: 'Showcase room joins',
  tabFirstFrame: 'First video frame rendered',
  tabShowcaseLeave: 'Leave live room',
  desc: 'Join count for showcase live rooms (category=1; each successful joinRoom +1, excludes anchor; includes re-join).',
  leaveDesc:
      'Leave count for showcase rooms (category=1); each successful leaveRoom +1, excludes anchor (API leave only, not kick/heartbeat cleanup).',
  firstFrameDesc:
      'Counts when the first live video frame is rendered in a showcase room (category=1). App reports via POST /liveRoom/reportLiveFirstFrameRendered; each success +1 (excludes anchor; client may dedupe per session).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  seriesName: 'Join count',
  firstFrameSeriesName: 'First frame count',
  leaveSeriesName: 'Leave count',
  chartDaily: 'Daily trend (last 30 days)',
  firstFrameChartDaily: 'First frame · daily (last 30 days)',
  firstFrameChartWeekly: 'First frame · weekly (last 12 weeks)',
  firstFrameChartMonthly: 'First frame · monthly (last 12 months)',
  leaveChartDaily: 'Leave · daily (last 30 days)',
  leaveChartWeekly: 'Leave · weekly (last 12 weeks)',
  leaveChartMonthly: 'Leave · monthly (last 12 months)',
  chartWeekly: 'Weekly trend (last 12 weeks)',
  chartMonthly: 'Monthly trend (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  desc:
      '统计秀场直播间（category=1）内发起的 1v1 视频通话次数：call_order.source=1（直播间来源），观众调用 liveRoomCall 成功创建通话单并推送呼叫后 +1（对齐埋点 call_1v1_initiate）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  seriesName: '发起次数',
  chartDaily: '按日趋势（最近 30 天）',
  chartWeekly: '按周趋势（最近 12 周）',
  chartMonthly: '按月趋势（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingCall1v1InitiateMessages = definePageMessagesFromEn(zh, {
  desc:
      '1v1 video call initiations in showcase rooms (category=1): call_order.source=1 (live room); +1 when liveRoomCall succeeds (event call_1v1_initiate).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  seriesName: 'Initiation count',
  chartDaily: 'Daily trend (last 30 days)',
  chartWeekly: 'Weekly trend (last 12 weeks)',
  chartMonthly: 'Monthly trend (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

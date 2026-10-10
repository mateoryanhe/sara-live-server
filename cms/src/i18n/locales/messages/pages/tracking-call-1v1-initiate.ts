import {definePageMessagesFromEn} from './_define'

const zh = {
  tabLiveRoom: '秀场直播间发起',
  tabOneToOneRoom: '1v1房间视频通话',
  tabConnectSuccess: '双方成功接通',
  desc:
      '统计秀场直播间（category=1）内发起的 1v1 视频通话次数：call_order.source=1（直播间来源），观众调用 liveRoomCall 成功创建通话单并推送呼叫后 +1（对齐埋点 call_1v1_initiate）。',
  oneToOneRoomDesc:
      '统计 1v1 房间内发起的视频通话次数：call_order.source=3（1v1 房间来源），oneToOneRoomCall 成功创建通话单并推送呼叫后 +1（对齐埋点 call_1v1_room）。',
  connectSuccessDesc:
      '统计双方成功接通的 1v1 视频通话次数：主被叫各 confirmCall 一次且首次扣费成功、订单进入通话中后 +1。含秀场直播间 source=1（category=1）与 1v1 房间 source=3（对齐埋点 call_1v1_connect_success）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  seriesName: '发起次数',
  oneToOneRoomSeriesName: '通话次数',
  connectSuccessSeriesName: '接通次数',
  chartDaily: '按日趋势（最近 30 天）',
  chartWeekly: '按周趋势（最近 12 周）',
  chartMonthly: '按月趋势（最近 12 月）',
  oneToOneRoomChartDaily: '1v1房间 · 按日（最近 30 天）',
  oneToOneRoomChartWeekly: '1v1房间 · 按周（最近 12 周）',
  oneToOneRoomChartMonthly: '1v1房间 · 按月（最近 12 月）',
  connectSuccessChartDaily: '成功接通 · 按日（最近 30 天）',
  connectSuccessChartWeekly: '成功接通 · 按周（最近 12 周）',
  connectSuccessChartMonthly: '成功接通 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingCall1v1InitiateMessages = definePageMessagesFromEn(zh, {
  tabLiveRoom: 'Showcase live room',
  tabOneToOneRoom: '1v1 room video calls',
  tabConnectSuccess: 'Both parties connected',
  desc:
      '1v1 video call initiations in showcase rooms (category=1): call_order.source=1 (live room); +1 when liveRoomCall succeeds (event call_1v1_initiate).',
  oneToOneRoomDesc:
      'Video call count in 1v1 rooms: call_order.source=3; +1 when oneToOneRoomCall succeeds (event call_1v1_room).',
  connectSuccessDesc:
      'Successful 1v1 video connections: +1 when both sides confirmCall and first charge succeeds (in-call). Includes showcase live room source=1 (category=1) and 1v1 room source=3 (event call_1v1_connect_success).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  seriesName: 'Initiation count',
  oneToOneRoomSeriesName: 'Call count',
  connectSuccessSeriesName: 'Connected count',
  chartDaily: 'Daily trend (last 30 days)',
  chartWeekly: 'Weekly trend (last 12 weeks)',
  chartMonthly: 'Monthly trend (last 12 months)',
  oneToOneRoomChartDaily: '1v1 room · daily (last 30 days)',
  oneToOneRoomChartWeekly: '1v1 room · weekly (last 12 weeks)',
  oneToOneRoomChartMonthly: '1v1 room · monthly (last 12 months)',
  connectSuccessChartDaily: 'Connected · daily (last 30 days)',
  connectSuccessChartWeekly: 'Connected · weekly (last 12 weeks)',
  connectSuccessChartMonthly: 'Connected · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

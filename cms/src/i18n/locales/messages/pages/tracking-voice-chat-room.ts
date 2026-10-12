import {definePageMessagesFromEn} from './_define'

const zh = {
  tabEnter: '进入语聊房',
  tabMicRequest: '申请上麦',
  tabGiftClick: '礼物面板点击',
  enterDesc:
      '统计进入语聊房：观众 joinRoom 成功且 live_room_cfgs.category=5(语聊)，不含主播；每次 +1（event voice_room_enter）。',
  micRequestDesc:
      '统计申请上麦：POST /liveRoom/applyVoiceChatMic（Sofie /sofie/studio/requestVoiceChatMic）成功 +1（event voice_room_mic_request）。',
  giftClickDesc:
      '统计礼物面板点击：POST /liveRoom/reportVoiceRoomGiftClick（Sofie /sofie/studio/notifyVoiceRoomGiftClick），登录即可、暂不要求 body，+1（event voice_room_gift_click）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  enterSeriesName: '进房次数',
  micRequestSeriesName: '申请上麦次数',
  giftClickSeriesName: '礼物点击次数',
  enterChartDaily: '进房 · 按日（最近 30 天）',
  enterChartWeekly: '进房 · 按周（最近 12 周）',
  enterChartMonthly: '进房 · 按月（最近 12 月）',
  micRequestChartDaily: '申请上麦 · 按日（最近 30 天）',
  micRequestChartWeekly: '申请上麦 · 按周（最近 12 周）',
  micRequestChartMonthly: '申请上麦 · 按月（最近 12 月）',
  giftClickChartDaily: '礼物点击 · 按日（最近 30 天）',
  giftClickChartWeekly: '礼物点击 · 按周（最近 12 周）',
  giftClickChartMonthly: '礼物点击 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingVoiceChatRoomMessages = definePageMessagesFromEn(zh, {
  tabEnter: 'Enter voice room',
  tabMicRequest: 'Mic request',
  tabGiftClick: 'Gift panel click',
  enterDesc:
      'Voice room enters: joinRoom success with category=5 (voice chat), excluding host (+1, voice_room_enter).',
  micRequestDesc:
      'Mic apply: POST /liveRoom/applyVoiceChatMic (Sofie requestVoiceChatMic) on success (+1, voice_room_mic_request).',
  giftClickDesc:
      'Gift panel click: POST /liveRoom/reportVoiceRoomGiftClick (Sofie notifyVoiceRoomGiftClick), login only (+1, voice_room_gift_click).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  enterSeriesName: 'Enter count',
  micRequestSeriesName: 'Mic request count',
  giftClickSeriesName: 'Gift click count',
  enterChartDaily: 'Enter · daily (last 30 days)',
  enterChartWeekly: 'Enter · weekly (last 12 weeks)',
  enterChartMonthly: 'Enter · monthly (last 12 months)',
  micRequestChartDaily: 'Mic request · daily (last 30 days)',
  micRequestChartWeekly: 'Mic request · weekly (last 12 weeks)',
  micRequestChartMonthly: 'Mic request · monthly (last 12 months)',
  giftClickChartDaily: 'Gift click · daily (last 30 days)',
  giftClickChartWeekly: 'Gift click · weekly (last 12 weeks)',
  giftClickChartMonthly: 'Gift click · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: '按 UTC 分钟聚合：进房(秀场/游戏/语聊)、离房(已有 Leave 事件)、送礼成功、公屏文字+付费弹幕、游戏下注、1v1 发起、充值下单。syndb 入库，保留约 3 天，每日 0 点清理过期桶。',
  chartTitle: '用户事件',
  chartTitleLast3Days: '用户事件（最近3天）',
  chartHint: '按 UTC 分钟聚合，服务端仅保留最近 3 天；点数按比例铺满 X 轴（点多会压缩，点少会拉开）。',
  rangeTruncated: '时间范围内采样点超过 {max}，仅展示部分分钟数据。',
  fetchFailed: '加载趋势失败',
  join: '进房',
  leave: '离房',
  gift: '送礼',
  chat: '聊天',
  bet: '游戏下注',
  call1v1: '1v1发起',
  rechargeOrder: '充值下单',
}

export const trackingMinuteLiveActivityMessages = definePageMessagesFromEn(zh, {
  pageDesc: 'UTC minute buckets: join (showcase/game/voice), leave (existing leave events), paid gifts, chat+paid danmaku, game bets, 1v1 initiate, recharge orders. syndb, ~3d retention, daily cleanup.',
  chartTitle: 'User events',
  chartTitleLast3Days: 'User events (last 3 days)',
  chartHint: 'UTC minute buckets; server keeps the last 3 days. Points fill the X-axis (compress when dense).',
  rangeTruncated: 'More than {max} points in range; showing a partial sample.',
  fetchFailed: 'Failed to load trend',
  join: 'Join',
  leave: 'Leave',
  gift: 'Gift',
  chat: 'Chat',
  bet: 'Game bet',
  call1v1: '1v1 initiate',
  rechargeOrder: 'Recharge order',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: '按 user_id：user_ext.email 从空首次写入时当日 val +1（bindEmail 或邮箱新号注册）；同用户重复登录不计；注销后同邮箱新号再计一次。',
  todayTotal: '今日人数',
  chartDaily: '按日趋势（最近 30 天）',
  seriesName: '人数',
  fetchFailed: '获取趋势失败',
} as const

export const trackingEmailAccountMessages = definePageMessagesFromEn(zh, {
  pageDesc:
    'Per user_id: first non-empty write to user_ext.email adds +1 daily val (bindEmail or new email signup). Repeat logins do not count; re-register after cancel (new user_id) counts again.',
  todayTotal: 'Today',
  chartDaily: 'Daily trend (last 30 days)',
  seriesName: 'Count',
  fetchFailed: 'Failed to load trend',
})
import {definePageMessagesFromEn} from './_define'

const zh = {
  button: '批量立即结算',
  selectRequired: '请先勾选需要立即结算的工会',
  confirmTitle: '确认立即结算',
  confirmMessage: '确定立即结算已选择的 {count} 个工会吗？系统会结算当前未结算流水，并生成新的审核中结算单。',
  result: '立即结算完成：生成 {settled} 条结算单，无待结算数据 {noData} 个，失败 {fail} 个',
  failGuildIds: '结算失败的工会ID：{ids}',
  requestFailed: '批量立即结算失败',
} as const

const en = {
  button: 'Settle Selected Now',
  selectRequired: 'Select at least one guild to settle',
  confirmTitle: 'Confirm Immediate Settlement',
  confirmMessage: 'Settle the selected {count} guilds now? Current unsettled income will be settled and new pending settlement orders will be created.',
  result: 'Settlement completed: {settled} orders created, {noData} guilds had no data, {fail} failed',
  failGuildIds: 'Failed guild IDs: {ids}',
  requestFailed: 'Batch immediate settlement failed',
} as const

export const guildImmediateSettlementMessages = definePageMessagesFromEn(zh, en)

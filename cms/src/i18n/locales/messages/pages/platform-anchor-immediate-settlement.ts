import {definePageMessagesFromEn} from './_define'

const zh = {
  button: '批量立即结算',
  selectRequired: '请先勾选需要立即结算的平台主播',
  confirmTitle: '确认立即结算',
  confirmMessage: '确定立即结算已选择的 {count} 个平台主播吗？系统会结算当前未结算流水，并生成新的审核中主播代付单。',
  result: '立即结算完成：生成 {settled} 条主播代付单，无待结算数据 {noData} 个，失败 {fail} 个',
  failAnchorIds: '结算失败的平台主播ID：{ids}',
  requestFailed: '平台主播批量立即结算失败',
} as const

const en = {
  button: 'Settle Selected Now',
  selectRequired: 'Select at least one platform anchor to settle',
  confirmTitle: 'Confirm Immediate Settlement',
  confirmMessage: 'Settle the selected {count} platform anchors now? Current unsettled income will be settled and new pending anchor payout orders will be created.',
  result: 'Settlement completed: {settled} payout orders created, {noData} anchors had no data, {fail} failed',
  failAnchorIds: 'Failed platform anchor IDs: {ids}',
  requestFailed: 'Batch immediate platform anchor settlement failed',
} as const

export const platformAnchorImmediateSettlementMessages = definePageMessagesFromEn(zh, en)

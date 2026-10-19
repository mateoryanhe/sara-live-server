import {definePageMessagesFromEn} from './_define'

const zh = {
  button: '批量重置未结算收益',
  selectRequired: '请先勾选需要重置未结算收益的平台主播',
  confirmTitle: '确认重置未结算收益',
  confirmMessage: '确定将已选择的 {count} 个平台主播的未结算收益直接清零吗？不会生成结算单、归档或任何日志，且不可自动恢复。',
  result: '重置完成：清零 {reset} 个，已为 0 的 {noData} 个，失败 {fail} 个',
  failAnchorIds: '失败的平台主播ID：{ids}',
  requestFailed: '批量重置未结算收益失败',
} as const

const en = {
  button: 'Reset Unsettled Income',
  selectRequired: 'Select at least one platform anchor',
  confirmTitle: 'Confirm Reset Unsettled Income',
  confirmMessage: 'Clear unsettled income for the selected {count} platform anchors? No payout orders, archives, or logs will be created. This cannot be undone automatically.',
  result: 'Done: {reset} cleared, {noData} already zero, {fail} failed',
  failAnchorIds: 'Failed platform anchor IDs: {ids}',
  requestFailed: 'Batch reset unsettled income failed',
} as const

export const platformAnchorResetUnsettledIncomeMessages = definePageMessagesFromEn(zh, en)

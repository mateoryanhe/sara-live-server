import {definePageMessagesFromEn} from './_define'

const zh = {
  button: '批量重置未结算',
  selectRequired: '请先勾选需要重置未结算的工会',
  confirmTitle: '确认重置未结算',
  confirmMessage:
    '确定将已选择的 {count} 个工会及其名下主播的未结算收益直接清零吗？不会生成结算单、归档或任何日志，且不可自动恢复。',
  result: '重置完成：清零 {reset} 个，已为 0 的 {noData} 个，失败 {fail} 个',
  failGuildIds: '失败的工会ID：{ids}',
  requestFailed: '批量重置未结算失败',
} as const

const en = {
  button: 'Reset Unsettled',
  selectRequired: 'Select at least one guild',
  confirmTitle: 'Confirm Reset Unsettled Income',
  confirmMessage:
    'Clear unsettled income for the selected {count} guild(s) and their anchors? No payout orders, archives, or logs will be created. This cannot be undone automatically.',
  result: 'Done: {reset} cleared, {noData} already zero, {fail} failed',
  failGuildIds: 'Failed guild IDs: {ids}',
  requestFailed: 'Batch reset unsettled income failed',
} as const

export const guildResetUnsettledIncomeMessages = definePageMessagesFromEn(zh, en)

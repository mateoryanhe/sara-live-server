import {definePageMessages} from './_define'

const zh = {
  auditGroup: '审核与代付',
  userId: '用户ID',
  userIdPlaceholder: '可选',
  status: '状态',
  id: '提现单ID',
  nickname: '昵称',
  diamondAmount: '钻石',
  goldAmountLegacy: '金币(旧)',
  withdrawUsd: 'USD',
  transferCurrency: '代付币种',
  transferLocalAmount: '代付金额',
  failMsg: '失败原因',
  createdAt: '申请时间',
  batchApprove: '批量通过',
  batchReject: '批量拒绝',
  batchSubmitPayout: '提交HaiPay代付',
  rejectTitle: '拒绝原因',
  rejectReasonPlaceholder: '可选，将展示给用户',
  submitPayoutConfirm: '确认对所选记录提交 HaiPay 代付？',
  batchResult: '成功 {ok} 条，失败 {fail} 条',
  fetchFailed: '加载列表失败',
  statusPending: '待审核',
  statusApproved: '已通过',
  statusTransferred: '已到账',
  statusTransferring: '代付中',
  statusRejected: '已拒绝',
}

const en = {
  auditGroup: 'Review & payout',
  userId: 'User ID',
  userIdPlaceholder: 'Optional',
  status: 'Status',
  id: 'Withdrawal ID',
  nickname: 'Nickname',
  diamondAmount: 'Diamonds',
  goldAmountLegacy: 'Gold (legacy)',
  withdrawUsd: 'USD',
  transferCurrency: 'Payout currency',
  transferLocalAmount: 'Payout amount',
  failMsg: 'Failure reason',
  createdAt: 'Applied at',
  batchApprove: 'Batch approve',
  batchReject: 'Batch reject',
  batchSubmitPayout: 'Submit HaiPay payout',
  rejectTitle: 'Reject reason',
  rejectReasonPlaceholder: 'Optional, visible to user',
  submitPayoutConfirm: 'Submit HaiPay payout for selected records?',
  batchResult: 'Success {ok}, failed {fail}',
  fetchFailed: 'Failed to load list',
  statusPending: 'Pending',
  statusApproved: 'Approved',
  statusTransferred: 'Paid',
  statusTransferring: 'Paying out',
  statusRejected: 'Rejected',
}

const es = {...en}
const pt = {...en}
const hi = {...en}
const id = {...en}

export const userWithdrawLogListMessages = definePageMessages(zh, en, es, pt, hi, id)

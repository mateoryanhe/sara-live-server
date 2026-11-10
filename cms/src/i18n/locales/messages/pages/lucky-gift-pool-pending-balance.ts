import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: '查看幸运礼物奖池临时入账记录（冷启/排空时不直加奖池，延时合并后标记已处理）。',
  amount: '待入池金额',
  status: '处理状态',
  statusAll: '全部',
  statusPending: '未处理',
  statusProcessed: '已处理',
  statusProcessing: '处理中',
  createdAt: '创建时间',
  updatedAt: '更新时间',
  fetchFailed: '加载临时入账列表失败',
} as const

export const luckyGiftPoolPendingBalanceMessages = definePageMessagesFromEn(zh, {
  pageDesc: 'Lucky-gift pool pending credits (written during cold-start/drain; merged into pool later).',
  amount: 'Pending amount',
  status: 'Status',
  statusAll: 'All',
  statusPending: 'Pending',
  statusProcessed: 'Processed',
  statusProcessing: 'Processing',
  createdAt: 'Created at',
  updatedAt: 'Updated at',
  fetchFailed: 'Failed to load pending balance list',
})

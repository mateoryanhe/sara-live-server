import {definePageMessagesFromEn} from './_define'

const zh = {
  dateRange: '统计日期',
  startDate: '开始日期',
  endDate: '结束日期',
  dateRangeSeparator: '至',
  statDate: '日期',
  channelName: '渠道名称',
  channelCode: '渠道编码',
  channelCodePlaceholder: '全部留空',
  registerCount: '注册人数',
  orderCount: '下单笔数',
  payUserCount: '支付人次',
  orderTotalAmount: '下单总金额(USD)',
  rechargeSuccessAmount: '成功充值总金额(USD)',
  dateRangeRequired: '请选择统计日期范围',
  channelCodeInvalid: '渠道编码须为正整数或留空',
  fetchFailed: '获取推广渠道日统计失败',
} as const

export const promotionChannelDailyStatListMessages = definePageMessagesFromEn(zh, {
  dateRange: 'Date range',
  startDate: 'Start',
  endDate: 'End',
  dateRangeSeparator: 'to',
  statDate: 'Date',
  channelName: 'Channel name',
  channelCode: 'Channel code',
  channelCodePlaceholder: 'All if empty',
  registerCount: 'Registrations',
  orderCount: 'Orders',
  payUserCount: 'Paying users',
  orderTotalAmount: 'Order total (USD)',
  rechargeSuccessAmount: 'Paid total (USD)',
  dateRangeRequired: 'Select a date range',
  channelCodeInvalid: 'Channel code must be a positive integer or empty',
  fetchFailed: 'Failed to load promotion daily stats',
})

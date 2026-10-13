import {definePageMessages} from './_define'

const zh = {
  detailTitle: '经销商工会结算详情',
  backToList: '返回经销商代付列表',
} as const

const en: Record<keyof typeof zh, string> = {
  detailTitle: 'Dealer Guild Settlement Details',
  backToList: 'Back to Dealer Payout List',
}

export const coinMerchantGuildTransferListMessages = definePageMessages(zh, en, en, en, en, en)

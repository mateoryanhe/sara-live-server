import {definePageMessages} from './_define'

const zh = {
  detailTitle: '币商工会结算详情',
  backToList: '返回币商代付列表',
} as const

const en: Record<keyof typeof zh, string> = {
  detailTitle: 'Coin Merchant Guild Settlement Details',
  backToList: 'Back to Coin Merchant Payout List',
}

export const coinMerchantGuildTransferListMessages = definePageMessages(zh, en, en, en, en, en)

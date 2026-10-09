import {definePageMessagesFromEn} from './_define'

const zh = {
  addGuild: '新增币商工会',
  editGuild: '编辑币商工会',
  sharePercentRequired: '请输入分佣比例',
  giftIncomeCumulative: '礼物流水累计',
  detailTitle: '币商工会详情',
  detailTitleWithName: '币商工会详情 - {name}',
  detailTitleWithId: '币商工会详情 - {id}',
} as const

export const coinMerchantGuildListMessages = definePageMessagesFromEn(zh, {
  addGuild: 'Add Coin Merchant Guild',
  editGuild: 'Edit Coin Merchant Guild',
  sharePercentRequired: 'Please enter share percent',
  giftIncomeCumulative: 'Cumulative gift flow',
  detailTitle: 'Coin Merchant Guild Detail',
  detailTitleWithName: 'Coin Merchant Guild Detail - {name}',
  detailTitleWithId: 'Coin Merchant Guild Detail - {id}',
})

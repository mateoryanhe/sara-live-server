import {definePageMessagesFromEn} from './_define'

const zh = {
  addGuild: '新增经销商工会',
  editGuild: '编辑经销商工会',
  sharePercentRequired: '请输入分佣比例',
  giftIncomeCumulative: '礼物流水累计',
  detailTitle: '经销商工会详情',
  detailTitleWithName: '经销商工会详情 - {name}',
  detailTitleWithId: '经销商工会详情 - {id}',
} as const

export const coinMerchantGuildListMessages = definePageMessagesFromEn(zh, {
  addGuild: 'Add Dealer Guild',
  editGuild: 'Edit Dealer Guild',
  sharePercentRequired: 'Please enter share percent',
  giftIncomeCumulative: 'Cumulative gift flow',
  detailTitle: 'Dealer Guild Detail',
  detailTitleWithName: 'Dealer Guild Detail - {name}',
  detailTitleWithId: 'Dealer Guild Detail - {id}',
})

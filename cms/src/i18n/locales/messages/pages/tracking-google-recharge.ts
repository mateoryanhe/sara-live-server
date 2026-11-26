import {definePageMessagesFromEn} from './_define'

const zh = {
  tabList: '面板点击',
  tabOrder: '档位点击',
  tabResult: '充值到账',
  listDesc:
      'Android / X-Platform=android：POST /rechargeCfg/rechargeCfgListForApp（Sofie browsePackages）成功返回 +1（event google_recharge_list）。',
  orderDesc:
      'Google Play 档位建单：createRechargeOrder 且 payChannel=Google(2) 成功 +1（event google_recharge_order）。',
  resultDesc:
      'Google Play 订单首次 completeOrder 发放金币成功 +1（event google_recharge_result）。',
  todayTotal: '今日',
  listSeriesName: '面板点击次数',
  orderSeriesName: '档位点击次数',
  resultSeriesName: '到账次数',
  chartDaily: '按日（最近 30 天）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingGoogleRechargeMessages = definePageMessagesFromEn(zh, {
  tabList: 'Panel click',
  tabOrder: 'Tier click',
  tabResult: 'Recharge credited',
  listDesc:
      'Android X-Platform: rechargeCfgListForApp / Sofie browsePackages success +1 (google_recharge_list).',
  orderDesc:
      'createRechargeOrder with Google pay channel success +1 (google_recharge_order).',
  resultDesc:
      'First completeOrder for Google Play order +1 (google_recharge_result).',
  todayTotal: 'Today',
  listSeriesName: 'List fetches',
  orderSeriesName: 'Tier clicks',
  resultSeriesName: 'Credits',
  chartDaily: 'Daily (last 30 days)',
  fetchFailed: 'Failed to load trend',
})
import {definePageMessagesFromEn} from './_define'

const zh = {
  tabPanelEnter: '充值面板曝光',
  tabChargeClick: '充值档位点击',
  tabChargeResult: '充值到账成功',
  panelEnterDesc:
      '统计充值/礼物面板曝光：POST /rechargeOrder/reportChargePanelEnter（Sofie /sofie/topup/notifyChargePanelEnter），登录即可、暂不要求 body，+1（event charge_panel_enter）。',
  chargeClickDesc:
      '统计确认充值档位/自定义金额建单：App createRechargeOrder / createChannelRechargeOrder / 经销商渠道建单成功 +1（event charge_click）。',
  chargeResultDesc:
      '统计支付回调到账成功：订单首次 completeOrder 发放金币成功 +1（event charge_result_server）。',
  todayTotal: '今日',
  weekTotal: '本周',
  monthTotal: '本月',
  panelEnterSeriesName: '曝光次数',
  chargeClickSeriesName: '点击建单次数',
  chargeResultSeriesName: '到账成功次数',
  panelEnterChartDaily: '面板曝光 · 按日（最近 30 天）',
  panelEnterChartWeekly: '面板曝光 · 按周（最近 12 周）',
  panelEnterChartMonthly: '面板曝光 · 按月（最近 12 月）',
  chargeClickChartDaily: '档位点击 · 按日（最近 30 天）',
  chargeClickChartWeekly: '档位点击 · 按周（最近 12 周）',
  chargeClickChartMonthly: '档位点击 · 按月（最近 12 月）',
  chargeResultChartDaily: '到账成功 · 按日（最近 30 天）',
  chargeResultChartWeekly: '到账成功 · 按周（最近 12 周）',
  chargeResultChartMonthly: '到账成功 · 按月（最近 12 月）',
  fetchFailed: '加载趋势数据失败',
}

export const trackingRechargeMessages = definePageMessagesFromEn(zh, {
  tabPanelEnter: 'Recharge panel view',
  tabChargeClick: 'Recharge tier click',
  tabChargeResult: 'Recharge success',
  panelEnterDesc:
      'Recharge/gift panel views: POST /rechargeOrder/reportChargePanelEnter (Sofie notifyChargePanelEnter), login only (+1, charge_panel_enter).',
  chargeClickDesc:
      'Confirm tier/custom amount order created: createRechargeOrder / channel / dealer order success (+1, charge_click).',
  chargeResultDesc:
      'Payment callback first success: completeOrder grants gold (+1, charge_result_server).',
  todayTotal: 'Today',
  weekTotal: 'This week',
  monthTotal: 'This month',
  panelEnterSeriesName: 'Panel view count',
  chargeClickSeriesName: 'Order click count',
  chargeResultSeriesName: 'Success count',
  panelEnterChartDaily: 'Panel view · daily (last 30 days)',
  panelEnterChartWeekly: 'Panel view · weekly (last 12 weeks)',
  panelEnterChartMonthly: 'Panel view · monthly (last 12 months)',
  chargeClickChartDaily: 'Tier click · daily (last 30 days)',
  chargeClickChartWeekly: 'Tier click · weekly (last 12 weeks)',
  chargeClickChartMonthly: 'Tier click · monthly (last 12 months)',
  chargeResultChartDaily: 'Success · daily (last 30 days)',
  chargeResultChartWeekly: 'Success · weekly (last 12 weeks)',
  chargeResultChartMonthly: 'Success · monthly (last 12 months)',
  fetchFailed: 'Failed to load trend data',
})

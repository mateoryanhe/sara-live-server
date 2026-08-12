/** User stat bar chart metric config (matches backend statdto.UserStatBarMetric*) */
export interface UserStatBarSeriesItem {
  key: string
  labelKey: string
  color: string
  enabled: boolean
}

export const USER_STAT_BAR_SERIES: UserStatBarSeriesItem[] = [
  {key: 'rechargeUser', labelKey: 'barMetricRechargeUser', color: '#409EFF', enabled: true},
  {key: 'vipUser', labelKey: 'barMetricVipUser', color: '#E6A23C', enabled: false},
  {key: 'goldConsumeUser', labelKey: 'barMetricGoldConsumeUser', color: '#67C23A', enabled: true},
  {key: 'diamondConsumeUser', labelKey: 'barMetricDiamondConsumeUser', color: '#9f7aea', enabled: true},
]

export const getEnabledUserStatBarSeries = () =>
    USER_STAT_BAR_SERIES.filter((item) => item.enabled)

export const getUserStatBarMetricTabs = getEnabledUserStatBarSeries

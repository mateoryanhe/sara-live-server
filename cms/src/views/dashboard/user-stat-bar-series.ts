/** 用户数据柱形图指标配置(与后端 statdto.UserStatBarMetric* 一致) */
export interface UserStatBarSeriesItem {
    /** 指标 key,对应接口 barMetrics 字段名 */
    key: string
    label: string
    color: string
    /** 是否展示(预留指标默认关闭,接入数据后改为 true) */
    enabled: boolean
}

/**
 * 柱形图系列配置: 新增指标时在此追加, 并在后端 BuildUserStatBarMetrics 中赋值
 */
export const USER_STAT_BAR_SERIES: UserStatBarSeriesItem[] = [
    {key: 'rechargeUser', label: '充值人数', color: '#409EFF', enabled: true},
    {key: 'vipUser', label: 'VIP用户', color: '#E6A23C', enabled: false},
    {key: 'goldConsumeUser', label: '金币消费用户', color: '#67C23A', enabled: true},
    {key: 'diamondConsumeUser', label: '钻石消费用户', color: '#9f7aea', enabled: true},
]

export const getEnabledUserStatBarSeries = () =>
    USER_STAT_BAR_SERIES.filter((item) => item.enabled)

/** 柱形图指标 Tab 列表(仅已启用指标) */
export const getUserStatBarMetricTabs = getEnabledUserStatBarSeries

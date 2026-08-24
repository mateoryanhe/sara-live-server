type TranslateFn = (key: string) => string

export const LIVE_REVENUE_TYPE_OPTIONS = [
  {value: 1, labelKey: 'pages.revenueLogList.revenueGift'},
  {value: 2, labelKey: 'pages.revenueLogList.revenuePaidDanmaku'},
  {value: 3, labelKey: 'pages.revenueLogList.revenueGameBet'},
  {value: 4, labelKey: 'pages.revenueLogList.revenuePrivateRoom'},
  {value: 5, labelKey: 'pages.revenueLogList.revenueTicket'},
  {value: 6, labelKey: 'pages.revenueLogList.revenueLiveRoomVideoCallTicket'},
  {value: 7, labelKey: 'pages.revenueLogList.revenueLiveRoomVideoCallBilling'},
] as const

export function createLiveRevenueTypeFormatter(t: TranslateFn) {
  const map = Object.fromEntries(
    LIVE_REVENUE_TYPE_OPTIONS.map(option => [option.value, t(option.labelKey)]),
  ) as Record<number, string>

  return (type: number) => map[type] || t('pages.revenueLogList.unknown')
}

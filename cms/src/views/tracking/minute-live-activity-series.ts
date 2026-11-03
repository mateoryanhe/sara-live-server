export const MINUTE_LIVE_ACTIVITY_SERIES_KEYS = [
  'joinCount',
  'leaveCount',
  'giftCount',
  'chatCount',
  'betCount',
  'call1v1Count',
  'call1v1ConnectSuccessCount',
  'call1v1ConsumeCount',
  'rechargeOrderCount',
] as const

export type MinuteLiveActivitySeriesKey = (typeof MINUTE_LIVE_ACTIVITY_SERIES_KEYS)[number]

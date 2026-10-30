export const MINUTE_LIVE_ACTIVITY_SERIES_KEYS = [
  'joinCount',
  'leaveCount',
  'giftCount',
  'chatCount',
  'betCount',
  'call1v1Count',
  'rechargeOrderCount',
] as const

export type MinuteLiveActivitySeriesKey = (typeof MINUTE_LIVE_ACTIVITY_SERIES_KEYS)[number]

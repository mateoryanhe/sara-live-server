/** 后端直播时长单位为秒，CMS 统一换算为分钟 */

export function liveDurationSecondsToMinutes(seconds: number | null | undefined): number | null {
  if (seconds === null || seconds === undefined || seconds <= 0) {
    return null
  }
  return Math.round(seconds / 60 * 10) / 10
}

export function formatLiveDurationMinutes(
  seconds: number | null | undefined,
  t: (key: string, params?: Record<string, unknown>) => string,
  i18nKey = 'pages.anchorList.liveDurationMinutes',
): string {
  const minutes = liveDurationSecondsToMinutes(seconds)
  if (minutes == null) {
    return '-'
  }
  return t(i18nKey, {minutes})
}

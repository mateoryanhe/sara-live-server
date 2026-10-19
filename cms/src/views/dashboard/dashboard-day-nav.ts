import {formatServerDateOnly, formatServerDateStringFromDate} from '@/utils/server-datetime'

export function getDashboardTodayDateString(): string {
  return formatServerDateStringFromDate(new Date())
}

export function buildDashboardSingleDayRange(day: string): [string, string] {
  return [day, day]
}

export function shiftDashboardFocusDay(day: string, deltaDays: number): string {
  const ms = Date.parse(`${day}T00:00:00.000Z`)
  if (Number.isNaN(ms)) {
    return getDashboardTodayDateString()
  }
  return formatServerDateStringFromDate(new Date(ms + deltaDays * 86_400_000))
}

export function isDashboardFocusDayToday(day: string): boolean {
  return day === getDashboardTodayDateString()
}

export function canGoDashboardNextDay(day: string): boolean {
  return day < getDashboardTodayDateString()
}

export function formatDashboardFocusDayLabel(day: string): string {
  return formatServerDateOnly(day)
}

export function clampDashboardFocusDay(day: string): string {
  const trimmed = day.trim()
  const today = getDashboardTodayDateString()
  if (!trimmed || trimmed > today) {
    return today
  }
  return trimmed
}

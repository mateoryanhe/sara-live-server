/** CMS 统一按服务器时区 UTC+0 展示与筛选时间 */

function pad2(n: number): string {
    return String(n).padStart(2, '0')
}

export function parseServerDateInput(value: string | number | Date): Date | null {
    if (value instanceof Date) {
        return Number.isNaN(value.getTime()) ? null : value
    }
    if (typeof value === 'number') {
        const ms = value < 1e12 ? value * 1000 : value
        const date = new Date(ms)
        return Number.isNaN(date.getTime()) ? null : date
    }
    const date = new Date(value)
    return Number.isNaN(date.getTime()) ? null : date
}

function formatUtcDateTimeParts(date: Date): string {
    return `${date.getUTCFullYear()}-${pad2(date.getUTCMonth() + 1)}-${pad2(date.getUTCDate())} ${pad2(date.getUTCHours())}:${pad2(date.getUTCMinutes())}:${pad2(date.getUTCSeconds())}`
}

function formatUtcDateOnlyParts(date: Date): string {
    return `${date.getUTCFullYear()}-${pad2(date.getUTCMonth() + 1)}-${pad2(date.getUTCDate())}`
}

function formatUtcDateOnlyFromMs(ms: number): string {
    const date = new Date(ms)
    return formatUtcDateOnlyParts(date)
}

/** 展示服务器下发的时间（UTC+0） */
export function formatServerDateTime(value: string | number | Date | null | undefined): string {
    if (value == null || value === '') {
        return '-'
    }
    const date = parseServerDateInput(value)
    if (!date) {
        return '-'
    }
    return formatUtcDateTimeParts(date)
}

/** CSV 等导出场景：UTC 时间，不带后缀 */
export function formatServerDateTimeForExport(value: string | number | Date | null | undefined): string {
    if (value == null || value === '') {
        return ''
    }
    const date = parseServerDateInput(value)
    if (!date) {
        return ''
    }
    return formatUtcDateTimeParts(date)
}

/** 仅日期部分（UTC+0），如 2025-08-27 */
export function formatServerDateOnly(value: string | number | Date | null | undefined): string {
    if (value == null || value === '') {
        return '-'
    }
    const date = parseServerDateInput(value)
    if (!date) {
        return '-'
    }
    return formatUtcDateOnlyParts(date)
}

/** 日期筛选：UTC 当天 00:00:00 的 Unix 秒 */
export function toServerDayStartUnix(dateStr: string): number {
    return Math.floor(Date.parse(`${dateStr}T00:00:00.000Z`) / 1000)
}

/** 日期筛选：UTC 当天 23:59:59 的 Unix 秒 */
export function toServerDayEndUnix(dateStr: string): number {
    return Math.floor(Date.parse(`${dateStr}T23:59:59.000Z`) / 1000)
}

/** UTC 本周一至周日（YYYY-MM-DD），与后端 WeekDateRange 在 UTC 服务器上一致 */
export function getServerWeekDateRange(ref: Date = new Date()): { start: string; end: string } {
    const day = ref.getUTCDay()
    const diff = day === 0 ? -6 : 1 - day
    const startMs = Date.UTC(ref.getUTCFullYear(), ref.getUTCMonth(), ref.getUTCDate() + diff)
    const endMs = startMs + 6 * 86400000
    return {
        start: formatUtcDateOnlyFromMs(startMs),
        end: formatUtcDateOnlyFromMs(endMs),
    }
}

/** UTC 上周一至上周日（YYYY-MM-DD） */
export function getServerLastWeekDateRange(ref: Date = new Date()): { start: string; end: string } {
    const thisWeek = getServerWeekDateRange(ref)
    const thisStartMs = Date.parse(`${thisWeek.start}T00:00:00.000Z`)
    const lastStartMs = thisStartMs - 7 * 86400000
    const lastEndMs = lastStartMs + 6 * 86400000
    return {
        start: formatUtcDateOnlyFromMs(lastStartMs),
        end: formatUtcDateOnlyFromMs(lastEndMs),
    }
}

/** 当前 UTC 时间 + N 天，用于默认表单值 */
export function formatServerNowPlusDays(days: number): string {
    return formatUtcDateTimeParts(new Date(Date.now() + days * 86400000))
}

/** Date 对象格式化为 UTC 日期 YYYY-MM-DD */
export function formatServerDateStringFromDate(date: Date): string {
    return formatUtcDateOnlyParts(date)
}

/** Date 对象格式化为 UTC 日期时间（无后缀，供表单默认值） */
export function formatServerDateTimeFromDate(date: Date): string {
    return formatUtcDateTimeParts(date)
}

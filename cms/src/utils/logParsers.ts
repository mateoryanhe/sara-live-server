import type {
    AccessLogItem,
    AccessLogStats,
    AccessTrendData,
    DetailLogItem,
    ErrorLogItem,
    PageResponse,
    TopStatItem,
    TraceLogDetail,
} from '@/types/api'

const LOG_QUERY_PATH = '/logQuery/'

const detailLogRe = /^(\S+)\s+\[(\w+)\]\s+\{([a-f0-9]+)\}\s+(.+)$/
const accessLogRe = /^(\S+)\s+\{([a-f0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP\/[\d.]+"\s+([\d.]+),\s+([^,]+),/
const errorLogHeaderRe = /^(\S+)\s+\[(\w+)\]\s+\{([a-f0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP\/[\d.]+"\s+([\d.]+),\s+([^,]+),/
const errorMetaRe = /,\s*(-?\d+),\s*"([^"]*)"(?:,\s*(.*))?$/
const reqIdRe = /reqId=(\d+)/
const authIdRe = /authId=(\d+)/
const userIdRe = /userid=(\d+)/
const playerRe = /玩家=(\d+)/
const urlRe = /url=([^,\s]+)/
const headersRe = /headers=(\{.*\})$/
const elapsedMsRes = [
    /totalMs=(-?\d+)ms/,
    /handlerMs=(-?\d+)ms/,
    /writeMs=(-?\d+)ms/,
    /bodyMs=(-?\d+)ms/,
    /authMs=(-?\d+)ms/,
    /afterOutputMs=(-?\d+)ms/,
    /从队列进入到中间件时间间隔Ms=(-?\d+)ms/,
]

const splitLines = (text: string) =>
    text.replace(/\r\n/g, '\n').split('\n').map((line) => line.trimEnd()).filter((line) => line.length > 0)

const extractUserAgent = (line: string) => {
    const lastQuote = line.lastIndexOf('"')
    if (lastQuote <= 0) {
        return ''
    }
    const prevQuote = line.lastIndexOf('"', lastQuote - 1)
    if (prevQuote < 0 || prevQuote + 1 >= lastQuote) {
        return ''
    }
    return line.slice(prevQuote + 1, lastQuote)
}

const extractAuthIdFromMessage = (message: string) => {
    const authMatch = message.match(authIdRe)
    if (authMatch?.[1]) {
        return authMatch[1]
    }
    const userMatch = message.match(userIdRe)
    if (userMatch?.[1]) {
        return userMatch[1]
    }
    const playerMatch = message.match(playerRe)
    if (playerMatch?.[1]) {
        return playerMatch[1]
    }
    return ''
}

const extractIdsFromLogHeaders = (message: string) => {
    const match = message.match(headersRe)
    if (!match?.[1]) {
        return {reqId: '', authId: ''}
    }
    try {
        const headers = JSON.parse(match[1]) as Record<string, string[]>
        let reqId = ''
        let authId = ''
        for (const [key, values] of Object.entries(headers)) {
            if (!values?.length || !values[0]) {
                continue
            }
            const lowerKey = key.toLowerCase()
            if (lowerKey === 'reqid' && !reqId) {
                reqId = values[0]
            }
            if (lowerKey === 'authorization' && !authId) {
                authId = values[0].split('.', 1)[0]
            }
            if (lowerKey === 'authid' && !authId) {
                authId = values[0]
            }
        }
        return {reqId, authId}
    } catch {
        return {reqId: '', authId: ''}
    }
}

const extractElapsedMsFromMessage = (message: string) => {
    for (const re of elapsedMsRes) {
        const match = message.match(re)
        if (match?.[1]) {
            const value = Number(match[1])
            if (!Number.isNaN(value)) {
                return value
            }
        }
    }
    return undefined
}

export const parseDetailLogLine = (line: string): DetailLogItem | null => {
    const trimmed = line.trim()
    if (!trimmed) {
        return null
    }
    const match = trimmed.match(detailLogRe)
    if (!match) {
        return null
    }
    const message = match[4].trim()
    const entry: DetailLogItem = {
        time: match[1],
        level: match[2],
        traceId: match[3],
        reqId: '',
        authId: '',
        url: '',
        message,
        raw: trimmed,
    }
    const reqMatch = message.match(reqIdRe)
    if (reqMatch?.[1]) {
        entry.reqId = reqMatch[1]
    }
    entry.authId = extractAuthIdFromMessage(message)
    if ((entry.reqId === '' || entry.authId === '') && message.includes('headers=')) {
        const headerIds = extractIdsFromLogHeaders(message)
        if (!entry.reqId) {
            entry.reqId = headerIds.reqId
        }
        if (!entry.authId) {
            entry.authId = headerIds.authId
        }
    }
    const urlMatch = message.match(urlRe)
    if (urlMatch?.[1]) {
        entry.url = urlMatch[1]
    }
    const elapsedMs = extractElapsedMsFromMessage(message)
    if (elapsedMs !== undefined) {
        entry.elapsedMs = elapsedMs
    }
    return entry
}

export const parseAccessLogLine = (line: string): AccessLogItem | null => {
    const trimmed = line.trim()
    if (!trimmed) {
        return null
    }
    const match = trimmed.match(accessLogRe)
    if (!match) {
        return null
    }
    return {
        time: match[1],
        traceId: match[2],
        statusCode: Number(match[3]) || 0,
        method: match[4],
        url: match[5],
        handlerMs: Number(match[6]) * 1000,
        ip: match[7].trim(),
        userAgent: extractUserAgent(trimmed),
        raw: trimmed,
    }
}

const fillErrorLogMeta = (entry: ErrorLogItem, line: string) => {
    const metaMatch = line.match(errorMetaRe)
    if (!metaMatch) {
        return
    }
    entry.errorCode = Number(metaMatch[1]) || 0
    entry.errorMessage = metaMatch[2] || ''
    entry.detail = (metaMatch[3] || '').trim()
}

const finalizeErrorLogEntry = (entry: ErrorLogItem, body: string) => {
    fillErrorLogMeta(entry, body)
    if (!entry.authId) {
        const authMatch = body.match(authIdRe)
        if (authMatch?.[1]) {
            entry.authId = authMatch[1]
        }
    }
    entry.stack = body
    entry.raw = body
}

const isLogQueryRelatedError = (entry: ErrorLogItem) => {
    for (const field of [entry.url, entry.raw, entry.detail, entry.stack, entry.errorMessage]) {
        if (field.includes(LOG_QUERY_PATH)) {
            return true
        }
    }
    return false
}

const parseErrorLogHeader = (line: string): ErrorLogItem | null => {
    const trimmed = line.trim()
    if (!trimmed) {
        return null
    }
    const match = trimmed.match(errorLogHeaderRe)
    if (!match) {
        return null
    }
    const entry: ErrorLogItem = {
        time: match[1],
        level: match[2],
        traceId: match[3],
        statusCode: Number(match[4]) || 0,
        method: match[5],
        url: match[6],
        handlerMs: Number(match[7]) * 1000,
        ip: match[8].trim(),
        errorCode: 0,
        errorMessage: '',
        detail: '',
        stack: '',
        raw: trimmed,
    }
    fillErrorLogMeta(entry, trimmed)
    return entry
}

const parseErrorLogBlock = (lines: string[]): ErrorLogItem | null => {
    if (!lines.length) {
        return null
    }
    const header = parseErrorLogHeader(lines[0])
    if (header) {
        finalizeErrorLogEntry(header, lines.join('\n'))
        return header
    }
    const detail = parseDetailLogLine(lines[0])
    if (!detail || !detail.message.includes('ErrorLog')) {
        return null
    }
    const body = lines.join('\n')
    const entry: ErrorLogItem = {
        time: detail.time,
        level: detail.level,
        traceId: detail.traceId,
        statusCode: 0,
        method: '',
        url: detail.url,
        handlerMs: detail.elapsedMs || 0,
        ip: '',
        authId: detail.authId,
        errorCode: 0,
        errorMessage: detail.message,
        detail: body,
        stack: body,
        raw: body,
    }
    return entry
}

export const parseErrorLogLines = (text: string): ErrorLogItem[] => {
    const entries: ErrorLogItem[] = []
    let block: string[] = []
    const flush = () => {
        if (!block.length) {
            return
        }
        const entry = parseErrorLogBlock(block)
        block = []
        if (entry && !isLogQueryRelatedError(entry)) {
            entries.push(entry)
        }
    }
    for (const line of splitLines(text)) {
        if (errorLogHeaderRe.test(line) || (line.includes('ErrorLog') && detailLogRe.test(line))) {
            flush()
        }
        block.push(line)
    }
    flush()
    return entries
}

export const parseLogExportPage = <T>(
    text: string,
    parser: (line: string) => T | null,
): T[] => {
    const items: T[] = []
    for (const line of splitLines(text)) {
        const item = parser(line)
        if (item) {
            items.push(item)
        }
    }
    return items
}

export const parseTraceExport = (text: string, traceId: string, startDate: string, endDate: string): TraceLogDetail => {
    const sections: Record<string, string[]> = {detail: [], access: [], error: []}
    let current = ''
    for (const line of text.replace(/\r\n/g, '\n').split('\n')) {
        if (line.startsWith('@detail')) {
            current = 'detail'
            continue
        }
        if (line.startsWith('@access')) {
            current = 'access'
            continue
        }
        if (line.startsWith('@error')) {
            current = 'error'
            continue
        }
        if (current && line.trim()) {
            sections[current].push(line)
        }
    }
    return {
        traceId,
        startDate,
        endDate,
        detailLogs: sections.detail.map((line) => parseDetailLogLine(line)).filter(Boolean) as DetailLogItem[],
        accessLogs: sections.access.map((line) => parseAccessLogLine(line)).filter(Boolean) as AccessLogItem[],
        errorLogs: parseErrorLogLines(sections.error.join('\n')),
    }
}

const parseTopSection = (lines: string[]): TopStatItem[] => {
    const items: TopStatItem[] = []
    for (const line of lines) {
        const match = line.trim().match(/^(\d+)\s+(.+)$/)
        if (!match) {
            continue
        }
        items.push({count: Number(match[1]), key: match[2].trim()})
    }
    return items
}

export const parseStatsExport = (text: string): AccessLogStats => {
    const sections: Record<string, string[]> = {urlTop: [], ipTop: []}
    let current = ''
    for (const line of text.replace(/\r\n/g, '\n').split('\n')) {
        if (line === 'urlTop' || line === 'ipTop') {
            current = line
            continue
        }
        if (current && line.trim()) {
            sections[current].push(line)
        }
    }
    return {
        urlTop: parseTopSection(sections.urlTop),
        ipTop: parseTopSection(sections.ipTop),
    }
}

export const parseTrendExport = (text: string): AccessTrendData => {
    const result: AccessTrendData = {
        intervalMinutes: 15,
        points: [],
        totalCount: 0,
        peakTime: '',
        peakCount: 0,
    }
    let section = ''
    for (const line of text.replace(/\r\n/g, '\n').split('\n')) {
        if (line === 'points' || line === 'meta') {
            section = line
            continue
        }
        if (!line.trim()) {
            continue
        }
        if (section === 'points') {
            const [time, countText] = line.split('\t')
            if (time && countText) {
                result.points.push({time, count: Number(countText) || 0})
            }
            continue
        }
        if (section === 'meta') {
            const [key, value] = line.split('\t')
            if (key === 'intervalMinutes') {
                result.intervalMinutes = Number(value) || result.intervalMinutes
            } else if (key === 'totalCount') {
                result.totalCount = Number(value) || 0
            } else if (key === 'peakTime') {
                result.peakTime = value || ''
            } else if (key === 'peakCount') {
                result.peakCount = Number(value) || 0
            }
        }
    }
    result.points.sort((a, b) => a.time.localeCompare(b.time))
    return result
}

export const buildPageResponse = <T>(items: T[], total: number): PageResponse<T> => ({
    total,
    data: items,
})

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

const detailLogRe = /^(\S+)\s+\[(\w+)\]\s+(?:\{([a-fA-F0-9]+)\}\s+)?(.+)$/
const accessLogRe = /^(\S+)\s+\{([a-fA-F0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP\/[\d.]+"\s+([\d.]+),\s+([^,]+),/
const accessLogRawRe = /^(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP\/[\d.]+"\s+([\d.]+),\s+([^,]+),/
const errorLogHeaderRe = /^(\S+)\s+\[(\w+)\]\s+\{([a-f0-9]+)\}\s+(\d+)\s+"(\w+)\s+\S+\s+\S+\s+(\S+)\s+HTTP\/[\d.]+"\s+([\d.]+),\s+([^,]+),/
const errorMetaRe = /,\s*(-?\d+),\s*"([^"]*)"(?:,\s*(.*))?$/
const errorLogBodyRe = /^ErrorLog\s+/
const errorLogStackLineRe = /^ErrorLog\s+stack\|/
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

const stripLogHeader = (line: string) => {
    const normalized = line.replace(/\r\n/g, '\n')
    const newlineIdx = normalized.indexOf('\n')
    const firstLine = newlineIdx >= 0 ? normalized.slice(0, newlineIdx) : normalized
    const rest = newlineIdx >= 0 ? normalized.slice(newlineIdx + 1) : ''
    const match = firstLine.match(detailLogRe)
    if (match) {
        const message = match[4].trim()
        return rest ? `${message}\n${rest}` : message
    }
    return normalized.trim()
}

const parseLegacyInlineError = (message: string) => {
    const stackIdx = message.indexOf(' stack=')
    const errIdx = message.indexOf('err=')
    if (errIdx < 0) {
        return {summary: message.trim(), stack: ''}
    }
    const summary = stackIdx < 0
        ? message.slice(errIdx + 4).trim()
        : message.slice(errIdx + 4, stackIdx).trim()
    if (stackIdx < 0) {
        return {summary, stack: ''}
    }
    let stack = message.slice(stackIdx + 7).trim()
    if (summary) {
        const escaped = summary.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
        stack = stack.replace(new RegExp(`^\\d+[.)\\]]?\\)?\\s*${escaped}`), '').trim()
    }
    return {
        summary,
        stack: normalizeErrorStack(stack),
    }
}

const isStackSectionMarker = (line: string) => {
    const trimmed = line.trim()
    return trimmed === 'Stack:' || trimmed === 'Stack'
}

const isRawStackFrameLine = (line: string) => {
    const trimmed = line.trim()
    if (!trimmed) {
        return false
    }
    if (isStackSectionMarker(trimmed)) {
        return true
    }
    if (/^\d+[.)]\)?\s*\S/.test(trimmed)) {
        return true
    }
    if (/^[A-Za-z]:[/\\]/.test(trimmed) || /^\/[\w./-]+/.test(trimmed)) {
        return true
    }
    if (/^\s+\S/.test(line) && !trimmed.includes('=')) {
        return true
    }
    return false
}

const parseErrorLogMessage = (message: string) => {
    const body = message.trim()
    if (!errorLogBodyRe.test(body)) {
        return {summary: body, stack: ''}
    }
    if (errorLogStackLineRe.test(body)) {
        return {summary: '', stack: normalizeErrorStack(body.replace(errorLogStackLineRe, ''))}
    }
    if (body.includes(' stack=')) {
        return parseLegacyInlineError(body)
    }
    const errMatch = body.match(/err=(.+)$/)
    if (errMatch?.[1]) {
        return {summary: errMatch[1].trim(), stack: ''}
    }
    return {summary: body.replace(/^ErrorLog\s+/, '').trim(), stack: ''}
}

export const normalizeErrorStack = (stack: string) => {
    const lines = stack
        .replace(/\r\n/g, '\n')
        .split('\n')
        .map((line) => line.trim())
        .filter((line) => line.length > 0)
    if (!lines.length) {
        return ''
    }
    const normalized: string[] = []
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i]
        if (/^\d+[.)]\)?\s/.test(line)) {
            const frame = line.replace(/^\d+[.)]\)?\s*/, '')
            const location = lines[i + 1]?.trim()
            if (location && !/^\d+[.)]\)?\s/.test(location) && !isStackSectionMarker(location)) {
                normalized.push(`${frame}\n    ${location}`)
                i++
            } else {
                normalized.push(frame)
            }
            continue
        }
        if (isStackSectionMarker(line)) {
            if (normalized.length > 0) {
                normalized.push('')
            }
            continue
        }
        if (line.startsWith('interface conversion:')) {
            continue
        }
        normalized.push(line)
    }
    return normalized.join('\n')
}

export const formatErrorSummary = (item: Pick<ErrorLogItem, 'errorMessage' | 'detail' | 'raw'>) => {
    if (item.errorMessage && !item.errorMessage.includes('ErrorLog')) {
        return item.errorMessage
    }
    const parsed = parseErrorLogMessage(stripLogHeader(item.raw || item.detail || item.errorMessage || ''))
    return parsed.summary || item.errorMessage || item.detail || item.raw || ''
}

const parseErrorStackFromText = (text: string): string => {
    if (!text) {
        return ''
    }
    const entry = parseErrorLogBlock(splitLines(text))
    if (entry?.stack) {
        return entry.stack
    }
    if (text.includes(' stack=')) {
        return parseLegacyInlineError(stripLogHeader(text)).stack
    }
    const parsed = parseErrorLogMessage(stripLogHeader(text))
    return parsed.stack
}

export const formatErrorStack = (item: Pick<ErrorLogItem, 'stack' | 'detail' | 'raw' | 'errorMessage'>) => {
    if (item.stack && !item.stack.includes('ErrorLog source=')) {
        const normalized = normalizeErrorStack(item.stack)
        if (normalized) {
            return normalized
        }
    }
    return parseErrorStackFromText(item.raw || item.detail || '')
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
        traceId: match[3] || '',
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
    if (trimmed.includes(LOG_QUERY_PATH)) {
        return null
    }
    const match = trimmed.match(accessLogRe)
    if (match) {
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
    const rawMatch = trimmed.match(accessLogRawRe)
    if (!rawMatch) {
        return null
    }
    return {
        time: '',
        traceId: '',
        statusCode: Number(rawMatch[1]) || 0,
        method: rawMatch[2],
        url: rawMatch[3],
        handlerMs: Number(rawMatch[4]) * 1000,
        ip: rawMatch[5].trim(),
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
    entry.detail = body
    entry.raw = body
    entry.errorMessage = formatErrorSummary(entry)
    if (!entry.stack) {
        entry.stack = parseErrorStackFromText(body)
    }
}

const isErrorLogBlockStart = (line: string) => {
    if (errorLogHeaderRe.test(line)) {
        return true
    }
    const detail = parseDetailLogLine(line)
    return !!detail && errorLogBodyRe.test(detail.message) && !errorLogStackLineRe.test(detail.message)
}

const isErrorLogStackContinuation = (line: string) => {
    const detail = parseDetailLogLine(line)
    if (detail && errorLogStackLineRe.test(detail.message)) {
        return true
    }
    if (!detail && !isErrorLogBlockStart(line) && isRawStackFrameLine(line)) {
        return true
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
    if (!detail || !errorLogBodyRe.test(detail.message)) {
        return null
    }
    const stackParts: string[] = []
    let summary = ''
    for (const line of lines) {
        const body = stripLogHeader(line)
        const parsed = parseErrorLogMessage(body)
        if (parsed.summary && errorLogBodyRe.test(body)) {
            summary = parsed.summary
        }
        if (parsed.stack) {
            stackParts.push(parsed.stack)
        } else if (isRawStackFrameLine(body)) {
            stackParts.push(isStackSectionMarker(body) ? 'Stack:' : body.trim())
        }
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
        errorMessage: summary,
        detail: body,
        stack: normalizeErrorStack(stackParts.join('\n')),
        raw: body,
    }
    if (!entry.errorMessage) {
        entry.errorMessage = formatErrorSummary(entry)
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
        if (isErrorLogBlockStart(line)) {
            flush()
        } else if (block.length && !isErrorLogStackContinuation(line)) {
            flush()
        }
        block.push(line)
    }
    flush()
    return entries
}

const isLogQueryRelatedError = (entry: ErrorLogItem) => {
    for (const field of [entry.url, entry.raw, entry.detail, entry.stack, entry.errorMessage]) {
        if (field.includes(LOG_QUERY_PATH)) {
            return true
        }
    }
    return false
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

export const LOG_JSON_FIELD_OPTIONS = [
    {key: 'headers', label: 'Header'},
    {key: 'respContent', label: '响应'},
    {key: 'bodyContent', label: 'Body'},
] as const

export type LogJsonFieldKey = typeof LOG_JSON_FIELD_OPTIONS[number]['key']

const logJsonFieldRe = (field: LogJsonFieldKey) => new RegExp(`(?:^|[,\\s])${field}=([\\s\\S]*)`)

const tryParseJsonValue = (raw: string): unknown | null => {
    const trimmed = raw.trim()
    if (!trimmed) {
        return null
    }
    try {
        return JSON.parse(trimmed)
    } catch {
        const starts = [
            trimmed.indexOf('{'),
            trimmed.indexOf('['),
        ].filter((index) => index >= 0)
        if (!starts.length) {
            return null
        }
        const start = Math.min(...starts)
        const candidate = trimmed.slice(start)
        for (let end = candidate.length; end > 0; end--) {
            const slice = candidate.slice(0, end).trim()
            if (!slice) {
                continue
            }
            try {
                return JSON.parse(slice)
            } catch {
                continue
            }
        }
        return null
    }
}

export const hasLogJsonField = (text: string, field: LogJsonFieldKey) => logJsonFieldRe(field).test(text)

export const listLogJsonFields = (text?: string) => {
    if (!text) {
        return []
    }
    return LOG_JSON_FIELD_OPTIONS.filter((item) => hasLogJsonField(text, item.key))
}

export const extractAndFormatLogJsonField = (text: string, field: LogJsonFieldKey) => {
    const match = text.match(logJsonFieldRe(field))
    if (!match?.[1]) {
        return null
    }
    const raw = match[1].trim()
    const parsed = tryParseJsonValue(raw)
    if (parsed === null) {
        return {
            raw,
            formatted: raw,
            parsed: false,
        }
    }
    return {
        raw,
        formatted: JSON.stringify(parsed, null, 2),
        parsed: true,
    }
}

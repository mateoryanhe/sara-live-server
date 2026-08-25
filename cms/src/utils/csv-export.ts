export interface CsvColumn<T> {
  header: string
  value: (row: T) => string | number | boolean | null | undefined
}

/** Excel 打开 CSV 时超长数字会丢精度; 超 2^53-1 的 ID 加 \\t 前缀强制按文本显示 */
export function formatCsvSnowflakeId(value: string | number | null | undefined): string {
  if (value === null || value === undefined || value === '') {
    return ''
  }
  const text = String(value).trim()
  if (!/^\d+$/.test(text)) {
    return text
  }
  if (text.length > 15 || BigInt(text) > BigInt(Number.MAX_SAFE_INTEGER)) {
    return `\t${text}`
  }
  return text
}

function escapeCsvCell(value: unknown): string {
  if (value === null || value === undefined) {
    return ''
  }
  const text = typeof value === 'string' || typeof value === 'number'
    ? formatCsvSnowflakeId(value)
    : String(value)
  if (/[",\r\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

export function buildCsvContent<T>(columns: CsvColumn<T>[], rows: T[]): string {
  const headerLine = columns.map(column => escapeCsvCell(column.header)).join(',')
  const dataLines = rows.map(row =>
    columns.map(column => escapeCsvCell(column.value(row))).join(','),
  )
  return `\uFEFF${[headerLine, ...dataLines].join('\r\n')}`
}

export function downloadCsv<T>(filename: string, columns: CsvColumn<T>[], rows: T[]): void {
  const content = buildCsvContent(columns, rows)
  const blob = new Blob([content], {type: 'text/csv;charset=utf-8;'})
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename.endsWith('.csv') ? filename : `${filename}.csv`
  link.click()
  URL.revokeObjectURL(url)
}

export async function fetchAllPagedRows<T>(
  fetchPage: (pageIndex: number, pageSize: number) => Promise<{data?: T[]; total?: number}>,
  pageSize = 100,
): Promise<T[]> {
  const all: T[] = []
  let pageIndex = 1
  let knownTotal: number | undefined

  while (true) {
    const response = await fetchPage(pageIndex, pageSize)
    const rows = Array.isArray(response.data) ? response.data : []
    if (rows.length === 0) {
      break
    }

    if (knownTotal === undefined && typeof response.total === 'number') {
      knownTotal = response.total
    }

    all.push(...rows)

    if (knownTotal !== undefined && all.length >= knownTotal) {
      break
    }

    pageIndex += 1
  }

  return all
}

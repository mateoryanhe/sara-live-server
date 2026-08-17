export interface CsvColumn<T> {
  header: string
  value: (row: T) => string | number | boolean | null | undefined
}

function escapeCsvCell(value: unknown): string {
  if (value === null || value === undefined) {
    return ''
  }
  const text = String(value)
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
  pageSize = 500,
): Promise<T[]> {
  const all: T[] = []
  let pageIndex = 1
  let total = 0

  while (true) {
    const response = await fetchPage(pageIndex, pageSize)
    const rows = response.data ?? []
    total = response.total ?? rows.length
    all.push(...rows)
    if (rows.length === 0 || all.length >= total) {
      break
    }
    pageIndex += 1
  }

  return all
}

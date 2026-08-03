/** CMS 数值展示/输入小数位数，全局改这里即可(如改为 4) */
export const NUMBER_DISPLAY_DECIMALS = 2

/** 表单数值输入框小数位数 */
export const NUMBER_INPUT_DECIMALS = NUMBER_DISPLAY_DECIMALS

function normalizeNumberString(value: number | string): string | null {
  const raw = typeof value === 'number' ? String(value) : String(value).trim()
  if (!raw) {
    return null
  }
  if (/^-?\d+(\.\d+)?$/.test(raw)) {
    return raw
  }
  const num = Number(value)
  if (!Number.isFinite(num)) {
    return null
  }
  return String(num)
}

/** 截断小数位(不四舍五入)，用于只读展示 */
export function formatNumberDisplay(
  value: number | string | null | undefined,
  emptyText = '-',
  decimals: number = NUMBER_DISPLAY_DECIMALS,
): string {
  if (value === null || value === undefined || value === '') {
    return emptyText
  }
  const source = normalizeNumberString(value)
  if (source === null) {
    return emptyText
  }
  const negative = source.startsWith('-')
  const unsigned = negative ? source.slice(1) : source
  const [intPart, fracPart = ''] = unsigned.split('.')
  const truncatedFrac = fracPart.slice(0, decimals).padEnd(decimals, '0')
  const formatted = decimals > 0 ? `${intPart}.${truncatedFrac}` : intPart
  return negative ? `-${formatted}` : formatted
}

/** 截断小数位(不四舍五入)，用于编辑表单赋值 */
export function truncateNumber(
  value: number | string | null | undefined,
  decimals: number = NUMBER_DISPLAY_DECIMALS,
): number {
  const text = formatNumberDisplay(value, '0', decimals)
  return Number(text)
}

export const formatAmount = formatNumberDisplay
export const formatPrice = formatNumberDisplay

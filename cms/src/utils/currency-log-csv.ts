import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string) => string

export type CurrencyLogCsvRow = {
  id: string
  action: number
  amount: number
  before: number
  after: number
  reasonText?: string
  createdAt?: string | null
}

export function formatCurrencyLogCsvDate(dateString: string | null | undefined): string {
  if (!dateString) {
    return ''
  }
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return ''
  }
}

export function buildCurrencyLogCsvColumns(t: TranslateFn, currencyType: 1 | 2): CsvColumn<CurrencyLogCsvRow>[] {
  const amountHeader = currencyType === 2
      ? t('pages.currencyLogList.diamondChange')
      : t('pages.currencyLogList.goldChange')
  return [
    {header: t('pages.currencyLogList.logId'), value: row => row.id},
    {
      header: t('pages.currencyLogList.changeType'),
      value: row => {
        if (row.action === 1) return t('pages.currencyLogList.actionIncrease')
        if (row.action === 2) return t('pages.currencyLogList.actionDecrease')
        return ''
      },
    },
    {header: amountHeader, value: row => row.amount},
    {header: t('pages.currencyLogList.beforeChange'), value: row => row.before},
    {header: t('pages.currencyLogList.afterChange'), value: row => row.after},
    {header: t('pages.currencyLogList.reason'), value: row => row.reasonText || ''},
    {header: t('pages.currencyLogList.time'), value: row => formatCurrencyLogCsvDate(row.createdAt)},
  ]
}

/** @deprecated use buildCurrencyLogCsvColumns */
export function buildGoldCurrencyLogCsvColumns(t: TranslateFn): CsvColumn<CurrencyLogCsvRow>[] {
  return buildCurrencyLogCsvColumns(t, 1)
}

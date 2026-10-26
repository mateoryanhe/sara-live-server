import type {CsvColumn} from './csv-export'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'

type TranslateFn = (key: string) => string

export type CurrencyLogCsvRow = {
  id: string
  currencyType?: number
  action: number
  amount: number
  before: number
  after: number
  reasonText?: string
  createdAt?: string | null
}

export function buildCurrencyLogCsvColumns(
    t: TranslateFn,
    currencyType: 0 | 1 | 2 = 0,
): CsvColumn<CurrencyLogCsvRow>[] {
  const amountHeader = currencyType === 2
      ? t('pages.currencyLogList.diamondChange')
      : currencyType === 1
          ? t('pages.currencyLogList.goldChange')
          : t('pages.currencyLogList.amountChange')
  return [
    {header: t('pages.currencyLogList.logId'), value: row => row.id},
    {
      header: t('pages.currencyLogList.currencyType'),
      value: row => {
        if (row.currencyType === 2) return t('pages.currencyLogList.currencyTypeDiamond')
        if (row.currencyType === 1) return t('pages.currencyLogList.currencyTypeGold')
        return ''
      },
    },
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
    {header: t('pages.currencyLogList.time'), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

/** @deprecated use buildCurrencyLogCsvColumns */
export function buildGoldCurrencyLogCsvColumns(t: TranslateFn): CsvColumn<CurrencyLogCsvRow>[] {
  return buildCurrencyLogCsvColumns(t, 1)
}

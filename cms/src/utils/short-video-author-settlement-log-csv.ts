import type {CsvColumn} from './csv-export'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'

type TranslateFn = (key: string) => string

export type ShortVideoAuthorSettlementLogCsvRow = {
  id: string
  unsettledIncome: number
  settlementDiamond: number
  anchorSharePercent: number
  createdAt?: string | null
}

export function buildShortVideoAuthorSettlementLogCsvColumns(
    t: TranslateFn,
): CsvColumn<ShortVideoAuthorSettlementLogCsvRow>[] {
  const ns = 'pages.shortVideoAuthorSettlementLogList'
  return [
    {header: t(`${ns}.logId`), value: row => row.id},
    {header: t(`${ns}.unsettledIncome`), value: row => row.unsettledIncome},
    {header: t(`${ns}.settlementDiamond`), value: row => row.settlementDiamond},
    {header: t(`${ns}.anchorSharePercent`), value: row => row.anchorSharePercent},
    {header: t(`${ns}.time`), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

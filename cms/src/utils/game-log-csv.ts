import type {GameBetLogItem, GameWinLogItem} from '@/types/api'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'
import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string) => string

export function buildGameBetLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.gameBetLogList',
): CsvColumn<GameBetLogItem>[] {
  return [
    {header: t(`${ns}.recordId`), value: row => row.id},
    {header: t('common.userId'), value: row => row.userId ?? ''},
    {header: t(`${ns}.gameCode`), value: row => row.gameCode ?? ''},
    {header: t(`${ns}.nameEn`), value: row => row.nameEn ?? ''},
    {header: t(`${ns}.betAmount`), value: row => row.amount ?? ''},
    {header: t(`${ns}.platform`), value: row => row.platformType ?? ''},
    {header: t(`${ns}.liveRoomId`), value: row => row.liveRoomId ?? ''},
    {header: t(`${ns}.liveRecordId`), value: row => row.liveRecordId ?? ''},
    {header: t(`${ns}.orderId`), value: row => row.orderId ?? ''},
    {header: t(`${ns}.time`), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

export function buildGameWinLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.gameWinLogList',
): CsvColumn<GameWinLogItem>[] {
  return [
    {header: t(`${ns}.recordId`), value: row => row.id},
    {header: t('common.userId'), value: row => row.userId ?? ''},
    {header: t(`${ns}.gameCode`), value: row => row.gameCode ?? ''},
    {header: t(`${ns}.nameEn`), value: row => row.nameEn ?? ''},
    {header: t(`${ns}.winAmount`), value: row => row.amount ?? ''},
    {header: t(`${ns}.platform`), value: row => row.platformType ?? ''},
    {header: t(`${ns}.orderId`), value: row => row.orderId ?? ''},
    {header: t(`${ns}.time`), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

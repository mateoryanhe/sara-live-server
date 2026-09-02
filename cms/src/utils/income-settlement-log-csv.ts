import type {IncomeSettlementLogAmounts} from '@/types/api'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'
import type {CsvColumn} from './csv-export'
import {liveDurationSecondsToMinutes} from './live-duration-format'

export type SettlementLogCsvRow = IncomeSettlementLogAmounts & {
  id: string
  roomId?: string
  roomNickname?: string
  guildId?: string
  guildName?: string
  createdAt?: string | null
}

type TranslateFn = (key: string) => string

function incomeAmountColumns(t: TranslateFn, ns: string): CsvColumn<SettlementLogCsvRow>[] {
  return [
    {header: t(`${ns}.totalIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.totalGiftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.totalPaidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.totalVideoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.totalVideoCallTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.totalVideoCallBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.totalShortVideoIncome`), value: row => row.totalShortVideoIncome},
    {header: t(`${ns}.totalGameIncome`), value: row => row.totalGameIncome},
    {header: t(`${ns}.totalLiveDuration`), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
  ]
}

function amountColumns(t: TranslateFn, ns: string): CsvColumn<SettlementLogCsvRow>[] {
  return [
    ...incomeAmountColumns(t, ns),
    {header: t(`${ns}.settlementSalary`), value: row => row.settlementSalary},
    {header: t(`${ns}.settlementFlowCommission`), value: row => row.settlementShareAmount ?? ''},
    {header: t(`${ns}.settlementShareAmountUsd`), value: row => row.settlementShareAmountUsd ?? ''},
    {header: t(`${ns}.anchorSharePercent`), value: row => row.anchorSharePercent ?? ''},
  ]
}

export function buildAnchorSettlementLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.anchorIncomeSettlementLogList',
): CsvColumn<SettlementLogCsvRow>[] {
  return [
    {header: t(`${ns}.logId`), value: row => row.id},
    {header: t(`${ns}.roomId`), value: row => row.roomId ?? ''},
    {header: t(`${ns}.roomNickname`), value: row => row.roomNickname ?? ''},
    ...amountColumns(t, ns),
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

function guildAmountColumns(t: TranslateFn, ns: string): CsvColumn<SettlementLogCsvRow>[] {
  return [
    ...incomeAmountColumns(t, ns),
    {header: t(`${ns}.guildSharePercent`), value: row => row.guildSharePercent ?? ''},
    {header: t(`${ns}.settlementReceivableUsd`), value: row => row.settlementReceivableUsd ?? ''},
  ]
}

export function buildGuildSettlementLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.guildIncomeSettlementLogList',
): CsvColumn<SettlementLogCsvRow>[] {
  return [
    {header: t(`${ns}.logId`), value: row => row.id},
    {header: t(`${ns}.guildId`), value: row => row.guildId ?? ''},
    {header: t(`${ns}.guildName`), value: row => row.guildName ?? ''},
    {header: t(`${ns}.settlementSalary`), value: row => row.settlementSalary ?? ''},
    {header: t(`${ns}.settlementShareAmount`), value: row => row.settlementShareAmount ?? ''},
    ...guildAmountColumns(t, ns),
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

export function buildGuildAnchorSettlementLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.guildAnchorIncomeSettlementLogList',
): CsvColumn<SettlementLogCsvRow>[] {
  return [
    {header: t(`${ns}.logId`), value: row => row.id},
    {header: t(`${ns}.guildName`), value: row => row.guildName ?? ''},
    {header: t(`${ns}.roomId`), value: row => row.roomId ?? ''},
    {header: t(`${ns}.roomNickname`), value: row => row.roomNickname ?? ''},
    ...amountColumns(t, ns),
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

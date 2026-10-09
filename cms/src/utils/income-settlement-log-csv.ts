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
  status?: number
  transferAt?: string | null
  transferOrderId?: string
  transferPlatformNo?: string
  transferLocalAmount?: number
  transferCurrency?: string
  transferFailMsg?: string
  createdAt?: string | null
}

type TranslateFn = (key: string) => string

function incomeAmountColumns(t: TranslateFn, ns: string): CsvColumn<SettlementLogCsvRow>[] {
  return [
    {header: t(`${ns}.totalIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.totalSocialIncome`), value: row => row.totalSocialIncome},
    {header: t(`${ns}.totalGiftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.totalPaidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.totalVideoCallIncome`), value: row => row.totalVideoCallIncome},
	{header: t(`${ns}.totalVideoCallTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.totalVideoCallBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.totalShortVideoIncome`), value: row => row.totalShortVideoIncome},
    {header: t(`${ns}.totalGameIncome`), value: row => row.totalGameIncome},
    {header: t(`${ns}.totalLiveDuration`), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
    {header: t(`${ns}.effectiveLiveDays`), value: row => row.effectiveLiveDays ?? ''},
  ]
}

function amountColumns(t: TranslateFn, ns: string): CsvColumn<SettlementLogCsvRow>[] {
  return [
    ...incomeAmountColumns(t, ns),
    {header: t(`${ns}.settlementSalary`), value: row => row.settlementSalary},
    {header: t(`${ns}.settlementFlowCommission`), value: row => row.settlementShareAmount ?? ''},
    {header: t(`${ns}.settlementShareAmountUsd`), value: row => row.settlementShareAmountUsd ?? ''},
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
    {header: t(`${ns}.hasSalary`), value: row => row.hasSalary ? t('common.yes') : t('common.no')},
    {header: t(`${ns}.anchorSocialSharePercent`), value: row => row.anchorSocialSharePercent ?? ''},
    {header: t(`${ns}.anchorSocialShareAmount`), value: row => row.anchorSocialShareAmount ?? ''},
    {header: t(`${ns}.guildSocialSharePercent`), value: row => row.guildSocialSharePercent ?? ''},
    {header: t(`${ns}.guildSocialShareAmount`), value: row => row.guildSocialShareAmount ?? ''},
    {header: t(`${ns}.anchorGameSharePercent`), value: row => row.anchorGameSharePercent ?? ''},
    {header: t(`${ns}.anchorGameShareAmountGold`), value: row => row.anchorGameShareAmountGold ?? ''},
    {header: t(`${ns}.guildGameSharePercent`), value: row => row.guildGameSharePercent ?? ''},
    {header: t(`${ns}.guildGameShareAmountGold`), value: row => row.guildGameShareAmountGold ?? ''},
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
  ]
}

export function buildPlatformAnchorPayoutCsvColumns(
  t: TranslateFn,
): CsvColumn<SettlementLogCsvRow>[] {
  const settlementColumns = buildAnchorSettlementLogCsvColumns(t)
  const settlementDetailColumns = settlementColumns.slice(3, settlementColumns.length - 1)
  return [
    {header: t('pages.anchorIncomeSettlementLogList.logId'), value: row => row.id},
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
    {header: t('pages.anchorIncomeSettlementLogList.roomId'), value: row => row.roomId ?? ''},
    {header: t('pages.anchorIncomeSettlementLogList.roomNickname'), value: row => row.roomNickname ?? ''},
    {header: t('pages.guildTransferList.settlementReceivableUsd'), value: row => row.settlementReceivableUsd ?? ''},
    {header: t('pages.guildTransferList.status'), value: row => row.status ?? ''},
    {header: t('pages.guildTransferList.transferCurrency'), value: row => row.transferCurrency ?? ''},
    {header: t('pages.guildTransferList.transferLocalAmount'), value: row => row.transferLocalAmount ?? ''},
    {header: t('pages.guildTransferList.transferOrderId'), value: row => row.transferOrderId ?? ''},
    {header: t('pages.guildTransferList.thirdPartyOrderId'), value: row => row.transferPlatformNo ?? ''},
    {header: t('pages.guildTransferList.transferAt'), value: row => formatServerDateTimeForExport(row.transferAt)},
    {header: t('pages.guildTransferList.transferFailMsg'), value: row => row.transferFailMsg ?? ''},
    ...settlementDetailColumns,
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
    ...guildAmountColumns(t, ns),
    {header: t(`${ns}.settlementRuleType`), value: row => row.settlementRuleType === 1
      ? t(`${ns}.settlementRuleTiered`)
      : t(`${ns}.settlementRuleLegacy`)},
    {header: t(`${ns}.anchorSocialShareAmount`), value: row => row.anchorSocialShareAmount ?? ''},
    {header: t(`${ns}.guildSocialShareAmount`), value: row => row.guildSocialShareAmount ?? ''},
    {header: t(`${ns}.anchorGameShareAmountGold`), value: row => row.anchorGameShareAmountGold ?? ''},
    {header: t(`${ns}.guildGameShareAmountGold`), value: row => row.guildGameShareAmountGold ?? ''},
    {header: t(`${ns}.goldToDiamondRate`), value: row => row.goldToDiamondRate ?? ''},
    {header: t(`${ns}.usdToGoldRate`), value: row => row.usdToGoldRate ?? ''},
    {header: t(`${ns}.gameShareAmountDiamond`), value: row => row.gameShareAmountDiamond ?? ''},
    {header: t(`${ns}.totalSettlementDiamond`), value: row => row.totalSettlementDiamond ?? ''},
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

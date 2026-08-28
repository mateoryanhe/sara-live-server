import type {AnchorDailyEffectiveLiveItem, GuildAnchorDailyEffectiveLiveItem, GuildDailyEffectiveLiveItem} from '@/types/api'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'
import type {CsvColumn} from './csv-export'
import {liveDurationSecondsToMinutes} from './live-duration-format'

type TranslateFn = (key: string) => string

export function buildAnchorDailyEffectiveLiveCsvColumns(
  t: TranslateFn,
  ns = 'pages.anchorList',
): CsvColumn<AnchorDailyEffectiveLiveItem>[] {
  return [
    {header: t(`${ns}.dailyRecordId`), value: row => row.id},
    {header: t(`${ns}.dailyLiveDate`), value: row => row.liveDate ?? ''},
    {header: t(`${ns}.dailyLiveDuration`), value: row => liveDurationSecondsToMinutes(row.liveDuration) ?? ''},
    {header: t(`${ns}.dailyReportedLiveDuration`), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
    {header: t(`${ns}.liveIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.giftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.paidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.shortVideoIncome`), value: row => row.totalShortVideoIncome},
    {header: t(`${ns}.gameIncome`), value: row => row.totalGameIncome},
    {
      header: t(`${ns}.dailySettled`),
      value: row => (row.settled ? t(`${ns}.dailySettledYes`) : t(`${ns}.dailySettledNo`)),
    },
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
    {header: t(`${ns}.roomUpdatedAt`), value: row => formatServerDateTimeForExport(row.updatedAt)},
  ]
}

export function buildGuildDailyEffectiveLiveCsvColumns(
  t: TranslateFn,
  ns = 'pages.anchorList',
): CsvColumn<GuildDailyEffectiveLiveItem>[] {
  return [
    {header: t(`${ns}.dailyRecordId`), value: row => row.id},
    {header: t(`${ns}.dailyLiveDate`), value: row => row.liveDate ?? ''},
    {header: t(`${ns}.dailyLiveDuration`), value: row => liveDurationSecondsToMinutes(row.liveDuration) ?? ''},
    {header: t(`${ns}.dailyReportedLiveDuration`), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
    {header: t(`${ns}.liveIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.giftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.paidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.shortVideoIncome`), value: row => row.totalShortVideoIncome},
    {header: t(`${ns}.gameIncome`), value: row => row.totalGameIncome},
    {
      header: t(`${ns}.dailySettled`),
      value: row => (row.settled ? t(`${ns}.dailySettledYes`) : t(`${ns}.dailySettledNo`)),
    },
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
    {header: t(`${ns}.roomUpdatedAt`), value: row => formatServerDateTimeForExport(row.updatedAt)},
  ]
}

export function buildGuildAnchorDailyEffectiveLiveCsvColumns(
  t: TranslateFn,
  ns = 'pages.anchorList',
): CsvColumn<GuildAnchorDailyEffectiveLiveItem>[] {
  return [
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomId'), value: row => row.roomId ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomNickname'), value: row => row.roomNickname ?? ''},
    {header: t(`${ns}.dailyLiveDate`), value: row => row.liveDate ?? ''},
    {header: t('pages.liveDailyEffectiveLiveList.unsettledTotalIncome'), value: row => row.unsettledTotalIncome},
    {header: t('pages.liveDailyEffectiveLiveList.dailyLiveIncome'), value: row => row.totalIncome},
    {header: t(`${ns}.dailyLiveDuration`), value: row => liveDurationSecondsToMinutes(row.liveDuration) ?? ''},
  ]
}

export function buildLiveDailyEffectiveLiveListCsvColumns(
  t: TranslateFn,
): CsvColumn<GuildAnchorDailyEffectiveLiveItem>[] {
  const ns = 'pages.liveDailyEffectiveLiveList'
  return [
    {header: t('pages.anchorList.dailyLiveDate'), value: row => row.liveDate ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomId'), value: row => row.roomId ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomNickname'), value: row => row.roomNickname ?? ''},
    {header: t('pages.anchorList.dailyLiveDuration'), value: row => liveDurationSecondsToMinutes(row.liveDuration) ?? ''},
    {header: t('pages.anchorList.dailyReportedLiveDuration'), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
    {header: t(`${ns}.dailyLiveIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.dailyGiftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.dailyPaidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.dailyVideoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.dailyVideoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.dailyVideoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.dailyShortVideoIncome`), value: row => row.totalShortVideoIncome},
    {header: t(`${ns}.dailyGameIncome`), value: row => row.totalGameIncome},
    {
      header: t('pages.anchorList.dailySettled'),
      value: row => (row.settled ? t('pages.anchorList.dailySettledYes') : t('pages.anchorList.dailySettledNo')),
    },
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
    {header: t('pages.anchorList.roomUpdatedAt'), value: row => formatServerDateTimeForExport(row.updatedAt)},
  ]
}

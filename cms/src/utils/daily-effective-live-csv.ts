import type {AnchorDailyEffectiveLiveItem, GuildAnchorDailyEffectiveLiveItem, GuildDailyEffectiveLiveItem} from '@/types/api'
import type {CsvColumn} from './csv-export'
import {liveDurationSecondsToMinutes} from './live-duration-format'

type TranslateFn = (key: string) => string

export function formatDailyEffectiveLiveCsvDate(dateString: string | null | undefined): string {
  if (!dateString) {
    return ''
  }
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return ''
  }
}

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
    {header: t(`${ns}.privateRoomTicketIncome`), value: row => row.totalPrivateRoomTicketIncome},
    {header: t(`${ns}.privateRoomWatchIncome`), value: row => row.totalPrivateRoomWatchIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {
      header: t(`${ns}.dailySettled`),
      value: row => (row.settled ? t(`${ns}.dailySettledYes`) : t(`${ns}.dailySettledNo`)),
    },
    {header: t('common.createdAt'), value: row => formatDailyEffectiveLiveCsvDate(row.createdAt)},
    {header: t(`${ns}.roomUpdatedAt`), value: row => formatDailyEffectiveLiveCsvDate(row.updatedAt)},
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
    {header: t(`${ns}.privateRoomTicketIncome`), value: row => row.totalPrivateRoomTicketIncome},
    {header: t(`${ns}.privateRoomWatchIncome`), value: row => row.totalPrivateRoomWatchIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {
      header: t(`${ns}.dailySettled`),
      value: row => (row.settled ? t(`${ns}.dailySettledYes`) : t(`${ns}.dailySettledNo`)),
    },
    {header: t('common.createdAt'), value: row => formatDailyEffectiveLiveCsvDate(row.createdAt)},
    {header: t(`${ns}.roomUpdatedAt`), value: row => formatDailyEffectiveLiveCsvDate(row.updatedAt)},
  ]
}

export function buildGuildAnchorDailyEffectiveLiveCsvColumns(
  t: TranslateFn,
  ns = 'pages.anchorList',
): CsvColumn<GuildAnchorDailyEffectiveLiveItem>[] {
  return [
    {header: t(`${ns}.dailyRecordId`), value: row => row.id},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomId'), value: row => row.roomId ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomNickname'), value: row => row.roomNickname ?? ''},
    {header: t(`${ns}.dailyLiveDate`), value: row => row.liveDate ?? ''},
    {header: t(`${ns}.dailyLiveDuration`), value: row => liveDurationSecondsToMinutes(row.liveDuration) ?? ''},
    {header: t(`${ns}.dailyReportedLiveDuration`), value: row => liveDurationSecondsToMinutes(row.totalLiveDuration) ?? ''},
    {header: t(`${ns}.liveIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.giftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.paidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.privateRoomTicketIncome`), value: row => row.totalPrivateRoomTicketIncome},
    {header: t(`${ns}.privateRoomWatchIncome`), value: row => row.totalPrivateRoomWatchIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {
      header: t(`${ns}.dailySettled`),
      value: row => (row.settled ? t(`${ns}.dailySettledYes`) : t(`${ns}.dailySettledNo`)),
    },
    {header: t('common.createdAt'), value: row => formatDailyEffectiveLiveCsvDate(row.createdAt)},
    {header: t(`${ns}.roomUpdatedAt`), value: row => formatDailyEffectiveLiveCsvDate(row.updatedAt)},
  ]
}

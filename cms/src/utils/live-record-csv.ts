import type {LiveRecordItem} from '@/types/api'
import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string) => string

export function formatLiveRecordCsvDate(dateString: string | null | undefined): string {
  if (!dateString) {
    return ''
  }
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return ''
  }
}

export function buildLiveRecordCsvColumns(
  t: TranslateFn,
  ns = 'pages.liveRecordList',
): CsvColumn<LiveRecordItem>[] {
  return [
    {header: t(`${ns}.recordId`), value: row => row.id},
    {header: t(`${ns}.anchorId`), value: row => row.anchorId},
    {header: t(`${ns}.anchorNickname`), value: row => row.nickname ?? ''},
    {header: t('common.startTime'), value: row => formatLiveRecordCsvDate(row.startTime)},
    {header: t('common.endTime'), value: row => formatLiveRecordCsvDate(row.endTime)},
    {header: t(`${ns}.totalAudience`), value: row => row.totalAudience},
    {header: t(`${ns}.liveDuration`), value: row => row.totalLiveDuration},
    {header: t(`${ns}.totalIncome`), value: row => row.totalIncome},
    {header: t(`${ns}.giftIncome`), value: row => row.totalGiftIncome},
    {header: t(`${ns}.paidDanmakuIncome`), value: row => row.totalPaidDanmakuIncome},
    {header: t(`${ns}.videoTicketIncome`), value: row => row.totalVideoCallTicketIncome},
    {header: t(`${ns}.videoBillingIncome`), value: row => row.totalVideoCallBillingIncome},
    {header: t(`${ns}.videoCallIncome`), value: row => row.totalVideoCallIncome},
    {header: t(`${ns}.giftSenderCount`), value: row => row.totalGiftSender},
    {header: t(`${ns}.newFollowers`), value: row => row.totalNewFollower},
    {header: t(`${ns}.totalGameBet`), value: row => row.totalGameBet},
    {header: t('common.createdAt'), value: row => formatLiveRecordCsvDate(row.createdAt)},
  ]
}

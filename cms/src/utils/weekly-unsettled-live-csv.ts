import type {WeeklyUnsettledLiveItem} from '@/types/api'
import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string) => string

export function buildLiveWeeklyUnsettledLiveListCsvColumns(
  t: TranslateFn,
): CsvColumn<WeeklyUnsettledLiveItem>[] {
  const ns = 'pages.liveWeeklyUnsettledLiveList'
  return [
    {header: t('pages.anchorList.dailyLiveDate'), value: row => row.liveDate ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomId'), value: row => row.roomId ?? ''},
    {header: t('pages.guildAnchorIncomeSettlementLogList.roomNickname'), value: row => row.roomNickname ?? ''},
    {header: t(`${ns}.weeklyUnsettledTotalIncome`), value: row => row.weeklyUnsettledTotalIncome},
    {header: t('pages.anchorList.dailyLiveDuration'), value: row => row.liveDuration},
    {header: t('pages.liveDailyEffectiveLiveList.dailyLiveIncome'), value: row => row.totalIncome},
  ]
}

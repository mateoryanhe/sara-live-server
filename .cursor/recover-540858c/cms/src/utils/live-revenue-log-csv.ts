import type {LiveRevenueLogItem} from '@/types/api'
import type {CsvColumn} from './csv-export'
import {formatLiveRecordCsvDate} from './live-record-csv'

type TranslateFn = (key: string) => string

export function buildLiveRevenueLogCsvColumns(
  t: TranslateFn,
  formatRevenueType: (type: number) => string,
  ns = 'pages.revenueLogList',
): CsvColumn<LiveRevenueLogItem>[] {
  return [
    {header: t(`${ns}.liveRecordId`), value: row => row.liveRecordId},
    {header: t(`${ns}.revenueType`), value: row => row.revenueTypeText || formatRevenueType(row.revenueType)},
    {header: t(`${ns}.liveRoomId`), value: row => row.roomId},
    {header: t(`${ns}.anchorNickname`), value: row => row.receiverNickname ?? ''},
    {header: t(`${ns}.payerUserId`), value: row => row.senderId},
    {header: t(`${ns}.payerNickname`), value: row => row.senderNickname ?? ''},
    {header: t(`${ns}.bizName`), value: row => row.bizName ?? ''},
    {header: t(`${ns}.count`), value: row => row.count},
    {header: t(`${ns}.unitPriceDiamond`), value: row => row.unitPrice},
    {header: t(`${ns}.totalAmountDiamond`), value: row => row.totalAmount},
    {
      header: t('common.status'),
      value: row => row.statusText || (row.status === 1 ? t(`${ns}.refunded`) : t('common.normal')),
    },
    {header: t('common.createdAt'), value: row => formatLiveRecordCsvDate(row.createdAt)},
  ]
}

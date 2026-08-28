import type {VideoCallLogItem} from '@/types/api'
import {formatServerDateTimeForExport} from '@/utils/server-datetime'
import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string) => string

export function buildVideoCallLogCsvColumns(
  t: TranslateFn,
  ns = 'pages.videoCallLogList',
): CsvColumn<VideoCallLogItem>[] {
  return [
    {header: t('common.createdAt'), value: row => formatServerDateTimeForExport(row.createdAt)},
    {header: t('pages.rechargeOrderList.orderId'), value: row => row.id},
    {header: t(`${ns}.callDuration`), value: row => row.callDuration},
    {header: t(`${ns}.totalCostDiamond`), value: row => row.totalCost},
    {header: t('common.status'), value: row => row.statusText ?? ''},
    {header: t(`${ns}.callerId`), value: row => row.callerId},
    {header: t(`${ns}.callerNickname`), value: row => row.callerNickname ?? ''},
    {header: t(`${ns}.receiverId`), value: row => row.receiverId},
    {header: t(`${ns}.receiverNickname`), value: row => row.receiverNickname ?? ''},
    {header: t(`${ns}.callTime`), value: row => formatServerDateTimeForExport(row.callStartTime)},
    {header: t(`${ns}.answerTime`), value: row => formatLiveRecordCsvDate(row.answerTime)},
    {header: t(`${ns}.callerLastHeart`), value: row => formatLiveRecordCsvDate(row.callerHeartTime)},
    {header: t(`${ns}.receiverLastHeart`), value: row => formatLiveRecordCsvDate(row.receiverHeartTime)},
    {header: t(`${ns}.endTime`), value: row => formatLiveRecordCsvDate(row.orderEndTime)},
    {header: t(`${ns}.callDuration`), value: row => row.callDuration},
    {header: t(`${ns}.ticketDiamond`), value: row => row.ticketPrice},
    {header: t(`${ns}.pricePerMinuteDiamond`), value: row => row.pricePerMinute},
    {header: t(`${ns}.billingDurationMinutes`), value: row => row.billingDuration},
    {header: t(`${ns}.totalCostDiamond`), value: row => row.totalCost},
    {header: t(`${ns}.lastChargeTime`), value: row => formatLiveRecordCsvDate(row.chargeTime)},
    {header: t('common.createdAt'), value: row => formatLiveRecordCsvDate(row.createdAt)},
  ]
}

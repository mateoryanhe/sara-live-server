import type {RechargeOrder} from '@/types/api'
import type {CsvColumn} from './csv-export'

type TranslateFn = (key: string, params?: Record<string, unknown>) => string

export function buildRechargeOrderCsvColumns(
  t: TranslateFn,
  ns = 'pages.rechargeOrderList',
): CsvColumn<RechargeOrder>[] {
  return [
    {header: t(`${ns}.orderId`), value: row => row.id},
    {header: t(`${ns}.userCountry`), value: row => row.userCountry ?? ''},
    {header: t('common.userId'), value: row => row.userId ?? ''},
    {header: t(`${ns}.userNickname`), value: row => row.nickname ?? ''},
    {header: t(`${ns}.cfgId`), value: row => row.cfgId ?? ''},
    {header: t(`${ns}.price`), value: row => row.price ?? ''},
    {header: t(`${ns}.payAmount`), value: row => row.payAmount ?? ''},
    {header: t(`${ns}.payCurrency`), value: row => row.currency ?? ''},
    {header: t(`${ns}.gold`), value: row => row.gold ?? ''},
    {header: t('common.status'), value: row => String(row.status ?? '')},
    {header: t(`${ns}.source`), value: row => String(row.source ?? '')},
    {header: t(`${ns}.thirdOrderId`), value: row => row.thirdOrderId ?? ''},
    {header: t('common.remark'), value: row => row.remark ?? ''},
    {header: t('common.createdAt'), value: row => row.createdAt ?? ''},
    {header: t(`${ns}.paidAt`), value: row => row.paidAt ?? ''},
  ]
}

export function buildRechargeOrderStatusLabels(t: TranslateFn, ns = 'pages.rechargeOrderList') {
  return {
    '0': t(`${ns}.statusPending`),
    '1': t(`${ns}.statusCompleted`),
    '2': t(`${ns}.statusCancelled`),
    '3': t(`${ns}.statusFailed`),
  }
}

export function buildRechargeOrderSourceLabels(t: TranslateFn, ns = 'pages.rechargeOrderList') {
  return {
    '1': t(`${ns}.sourceAndroid`),
    '2': t(`${ns}.sourceManual`),
    '3': t(`${ns}.sourceH5`),
    '5': t(`${ns}.sourceIOS`),
    '6': t(`${ns}.sourceCoinMerchant`),
  }
}

import {definePageMessagesFromEn} from './_define'

const zh = {
  tabTitle: '平台主播结算配置',
  sectionTitle: '平台主播最低结算金额',
  tipTitle: '最低结算规则',
  tipLine1: '平台主播折算后的可收金额严格大于该金额时，才生成结算单；未达到时保留未结算流水，参与下一次结算。',
  tipLine2: '数据库未配置或配置值无效时，服务端默认使用 5.00 USD。保存后立即刷新缓存，无需重启服务。',
  minimumSettlementUsd: '最低结算金额',
  usd: 'USD',
  lastUpdated: '最近更新',
  fetchFailed: '获取平台主播结算配置失败',
  saveFailed: '保存平台主播结算配置失败',
  saveSuccess: '平台主播结算配置已保存并生效',
  rangeWarning: '最低结算金额必须大于 0 且不能超过 1000000 USD',
} as const

const en = {
  tabTitle: 'Platform anchor settlement',
  sectionTitle: 'Platform anchor minimum settlement amount',
  tipTitle: 'Minimum settlement rule',
  tipLine1: 'A settlement order is created only when the converted platform-anchor receivable is strictly greater than this amount. Otherwise, unsettled revenue is retained for the next settlement.',
  tipLine2: 'When no valid database value exists, the server uses 5.00 USD. Saving refreshes the cache immediately without a restart.',
  minimumSettlementUsd: 'Minimum settlement amount',
  usd: 'USD',
  lastUpdated: 'Last updated',
  fetchFailed: 'Failed to load platform anchor settlement config',
  saveFailed: 'Failed to save platform anchor settlement config',
  saveSuccess: 'Platform anchor settlement config saved and applied',
  rangeWarning: 'The minimum settlement amount must be greater than 0 and no more than 1000000 USD',
} as const

export const platformAnchorSettlementCfgMessages = definePageMessagesFromEn(zh, en)

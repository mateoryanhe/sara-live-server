import {definePageMessagesFromEn} from './_define'

const zh = {
  tabTitle: '幸运礼物',
  tipTitle: '幸运礼物结算比例',
  tipLine1: '幸运礼物相关收益/流水按该比例参与结算（百分比）。',
  tipLine2: '保存后写入 wallet_exchange_cfgs 并刷新缓存，默认 10%。',
  settlementPercent: '结算比例',
  percentUnit: '%',
  lastUpdated: '最近更新',
  fetchFailed: '获取幸运礼物结算配置失败',
  saveFailed: '保存幸运礼物结算配置失败',
  saveSuccess: '幸运礼物结算配置已保存并生效',
  rangeWarning: '结算比例须在 0–100 之间',
} as const

const en = {
  tabTitle: 'Lucky gift',
  tipTitle: 'Lucky gift settlement ratio',
  tipLine1: 'Lucky-gift related revenue uses this percentage for settlement.',
  tipLine2: 'Saved to wallet_exchange_cfgs with cache refresh; default 10%.',
  settlementPercent: 'Settlement ratio',
  percentUnit: '%',
  lastUpdated: 'Last updated',
  fetchFailed: 'Failed to load lucky gift settlement config',
  saveFailed: 'Failed to save lucky gift settlement config',
  saveSuccess: 'Lucky gift settlement config saved and applied',
  rangeWarning: 'Settlement ratio must be between 0 and 100',
} as const

export const luckyGiftSettlementCfgMessages = definePageMessagesFromEn(zh, en)

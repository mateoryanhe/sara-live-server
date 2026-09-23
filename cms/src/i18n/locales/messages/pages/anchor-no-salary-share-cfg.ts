import {definePageMessagesFromEn} from './_define'

const zh = {
  tip: '配置无底薪主播的社交钻石流水提成比例；游戏金币流水使用对应的动态档位配置。',
  anchorSocialSharePercent: '主播社交提成比(%)',
  guildSocialSharePercent: '工会社交提成比(%)',
  lastUpdated: '最近更新',
  fetchFailed: '获取无底薪主播提成配置失败',
  saveFailed: '保存无底薪主播提成配置失败',
  percentRangeInvalid: '提成比例需在 0～100 之间',
} as const

const en: Record<keyof typeof zh, string> = {
  tip: 'Configure social diamond revenue commissions for non-salaried anchors. Game gold revenue uses its tier configuration.',
  anchorSocialSharePercent: 'Anchor Social Commission (%)',
  guildSocialSharePercent: 'Guild Social Commission (%)',
  lastUpdated: 'Last Updated',
  fetchFailed: 'Failed to load the non-salaried anchor commission configuration',
  saveFailed: 'Failed to save the non-salaried anchor commission configuration',
  percentRangeInvalid: 'The commission percentage must be between 0 and 100',
}

export const anchorNoSalaryShareCfgMessages = definePageMessagesFromEn(zh, en)

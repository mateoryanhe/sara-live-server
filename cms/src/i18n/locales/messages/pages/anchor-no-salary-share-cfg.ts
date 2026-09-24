import {definePageMessagesFromEn} from './_define'

const zh = {
  dynamicTierHint: '等级和档位均可动态配置；无底薪主播的社交流水按钻石统计，达到多个档位时取最高门槛。',
  level: '等级',
  socialTotalDiamondRevenue: '社交总钻石流水',
  anchorSocialSharePercent: '主播社交提成比(%)',
  guildSocialSharePercent: '工会社交提成比(%)',
  editTier: '编辑分佣配置',
  addTier: '新增档位',
  deleteConfirm: '确定删除等级 {level} 的分佣配置吗？',
  fetchFailed: '获取无底薪社交流水分佣配置失败',
  socialTotalRevenueRequired: '请输入社交总流水',
  levelRequired: '请输入大于 0 的等级',
  anchorSocialSharePercentRequired: '请输入主播社交提成比',
  guildSocialSharePercentRequired: '请输入工会社交提成比',
  sharePercentInvalid: '提成比例需在 0～100 之间',
} as const

const en: Record<keyof typeof zh, string> = {
  dynamicTierHint: 'Levels and tiers are configurable. Non-salaried social revenue is measured in diamonds, and the highest matched threshold applies.',
  level: 'Level',
  socialTotalDiamondRevenue: 'Total Social Revenue (Diamonds)',
  anchorSocialSharePercent: 'Anchor Social Commission (%)',
  guildSocialSharePercent: 'Guild Social Commission (%)',
  editTier: 'Edit Commission Tier',
  addTier: 'Add Tier',
  deleteConfirm: 'Delete the commission configuration for level {level}?',
  fetchFailed: 'Failed to load non-salaried social commission tiers',
  socialTotalRevenueRequired: 'Enter total social revenue',
  levelRequired: 'Enter a level greater than 0',
  anchorSocialSharePercentRequired: 'Enter the anchor social commission percentage',
  guildSocialSharePercentRequired: 'Enter the guild social commission percentage',
  sharePercentInvalid: 'The commission percentage must be between 0 and 100',
}

export const anchorNoSalaryShareCfgMessages = definePageMessagesFromEn(zh, en)
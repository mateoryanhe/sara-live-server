import {definePageMessagesFromEn} from './_define'

const zh = {
  dynamicTierHint: '等级和档位均可动态配置；每一档的社交总钻石流水是升级到下一等级的边界，达到边界立即使用下一等级比例，最高等级封顶。',
  level: '等级',
  socialTotalDiamondRevenue: '社交总钻石流水',
  anchorSocialSharePercent: '主播社交提成比(%)',
  guildSocialSharePercent: '工会社交提成比(%)',
  editTier: '编辑分佣配置',
  addTier: '新增档位',
  deleteConfirm: '确定删除等级 {level} 的分佣配置吗？',
  fetchFailed: '获取有底薪社交流水分佣配置失败',
  socialTotalRevenueRequired: '请输入社交总流水',
  socialTotalRevenueInteger: '社交总钻石流水必须是非负整数',
  levelRequired: '请输入大于 0 的等级',
  anchorSocialSharePercentRequired: '请输入主播社交提成比',
  guildSocialSharePercentRequired: '请输入工会社交提成比',
  sharePercentInvalid: '提成比例需在 0～100 之间',
} as const

const en: Record<keyof typeof zh, string> = {
  dynamicTierHint: 'Levels and tiers are configurable. Each diamond-revenue value is the boundary for promotion to the next level; reaching it immediately applies the next level, capped at the highest level.',
  level: 'Level',
  socialTotalDiamondRevenue: 'Total Social Revenue (Diamonds)',
  anchorSocialSharePercent: 'Anchor Social Commission (%)',
  guildSocialSharePercent: 'Guild Social Commission (%)',
  editTier: 'Edit Commission Tier',
  addTier: 'Add Tier',
  deleteConfirm: 'Delete the commission configuration for level {level}?',
  fetchFailed: 'Failed to load salaried social commission tiers',
  socialTotalRevenueRequired: 'Enter total social revenue',
  socialTotalRevenueInteger: 'Total social diamond revenue must be a non-negative integer',
  levelRequired: 'Enter a level greater than 0',
  anchorSocialSharePercentRequired: 'Enter the anchor social commission percentage',
  guildSocialSharePercentRequired: 'Enter the guild social commission percentage',
  sharePercentInvalid: 'The commission percentage must be between 0 and 100',
}

export const anchorSalarySocialShareCfgListMessages = definePageMessagesFromEn(zh, en)

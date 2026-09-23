import {definePageMessagesFromEn} from './_define'

const zh = {
  withSalaryTitle: '有底薪游戏流水档位分佣配置',
  noSalaryTitle: '无底薪游戏流水档位分佣配置',
  dynamicTierHint: '等级和档位数量均可动态配置；游戏流水按金币统计，达到多个档位时取最高门槛。',
  level: '等级',
  gameTotalGoldRevenue: '游戏总金币流水',
  anchorGameSharePercent: '主播游戏提成比(%)',
  guildGameSharePercent: '工会游戏提成比(%)',
  addTier: '新增档位',
  editTier: '编辑游戏分佣配置',
  levelRequired: '请输入大于 0 的等级',
  fetchFailed: '获取游戏流水档位分佣配置失败',
  gameTotalGoldRevenueRequired: '请输入游戏总金币流水',
  anchorGameSharePercentRequired: '请输入主播游戏提成比',
  guildGameSharePercentRequired: '请输入工会游戏提成比',
  sharePercentInvalid: '提成比例需在 0～100 之间',
  deleteConfirm: '确认删除等级 {level} 的游戏分佣配置吗？',
} as const

const en: Record<keyof typeof zh, string> = {
  withSalaryTitle: 'Salaried Game Revenue Tier Commission',
  noSalaryTitle: 'Non-salaried Game Revenue Tier Commission',
  dynamicTierHint: 'Levels and tier count are configurable. Game revenue is measured in gold, and the highest matched threshold applies.',
  level: 'Level',
  gameTotalGoldRevenue: 'Total Game Revenue (Gold)',
  anchorGameSharePercent: 'Anchor Game Commission (%)',
  guildGameSharePercent: 'Guild Game Commission (%)',
  addTier: 'Add Tier',
  editTier: 'Edit Game Commission Tier',
  levelRequired: 'Enter a level greater than 0',
  fetchFailed: 'Failed to load game revenue commission tiers',
  gameTotalGoldRevenueRequired: 'Enter total game gold revenue',
  anchorGameSharePercentRequired: 'Enter the anchor game commission percentage',
  guildGameSharePercentRequired: 'Enter the guild game commission percentage',
  sharePercentInvalid: 'The commission percentage must be between 0 and 100',
  deleteConfirm: 'Delete the game commission configuration for level {level}?',
}

export const anchorGameShareCfgListMessages = definePageMessagesFromEn(zh, en)

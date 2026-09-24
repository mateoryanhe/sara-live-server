import {definePageMessagesFromEn} from './_define'

const zh = {
  walletSectionTitle: '金币兑换配置',
  sectionTitle: '有效直播时长配置',
  tipTitle: '有效直播时长规则',
  tipLine1: '单场直播时长严格大于该门槛时，整场时长才计入主播和工会的每日有效直播时长。',
  tipLine2: '数据库未配置或配置值无效时，服务端默认使用 30 分钟。保存后立即生效，无需重启服务。',
  minSessionMinutes: '单场有效直播门槛',
  minutes: '分钟',
  lastUpdated: '最近更新',
  fetchFailed: '获取有效直播时长配置失败',
  saveFailed: '保存有效直播时长配置失败',
  saveSuccess: '有效直播时长配置已保存并生效',
  rangeWarning: '门槛必须在 1 到 1440 分钟之间',
} as const

const en = {
  walletSectionTitle: 'Gold exchange configuration',
  sectionTitle: 'Effective live duration configuration',
  tipTitle: 'Effective live duration rule',
  tipLine1: 'A full session is counted toward the anchor and guild daily effective live duration only when its duration is strictly greater than this threshold.',
  tipLine2: 'When no valid database value exists, the server uses 30 minutes. Changes take effect immediately without a restart.',
  minSessionMinutes: 'Session threshold',
  minutes: 'minutes',
  lastUpdated: 'Last updated',
  fetchFailed: 'Failed to load effective live duration config',
  saveFailed: 'Failed to save effective live duration config',
  saveSuccess: 'Effective live duration config saved and applied',
  rangeWarning: 'The threshold must be between 1 and 1440 minutes',
} as const

export const effectiveLiveCfgMessages = definePageMessagesFromEn(zh, en)

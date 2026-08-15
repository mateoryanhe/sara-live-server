import {definePageMessagesFromEn} from './_define'

const zh = {
  guildName: '工会名称',
  guildNamePlaceholder: '请输入工会名称',
  leader: '会长',
  description: '描述',
  fetchFailed: '获取工会垃圾库失败',
  onShelfConfirm: '确定要上架工会 "{name}" 吗？该工会下所有主播间将同步上架。',
  onShelfSuccess: '上架成功',
  onShelfFailed: '上架失败',
} as const

const en: Record<keyof typeof zh, string> = {
  guildName: 'Guild Name',
  guildNamePlaceholder: 'Enter guild name',
  leader: 'Leader',
  description: 'Description',
  fetchFailed: 'Failed to load guild recycle bin',
  onShelfConfirm: 'Publish guild "{name}"? All live rooms under this guild will also be published.',
  onShelfSuccess: 'Published successfully',
  onShelfFailed: 'Failed to publish',
}

export const guildRecycleBinMessages = definePageMessagesFromEn(zh, en)

import {definePageMessagesFromEn} from './_define'

const zh = {
  keywordPlaceholder: '用户ID/昵称/手机号/分享码',
  guildId: '工会ID',
  roomTitle: '直播间标题',
  roomType: '房间类型',
  categoryHot: 'Hot',
  categoryGame: 'Game',
  categoryPrivate: '私密',
  fetchFailed: '获取回收站列表失败',
  onShelfConfirm: '确定要上架主播间 {id} 吗？上架后将重新进入App列表。',
  onShelfSuccess: '上架成功',
  onShelfFailed: '上架失败',
} as const

const en: Record<keyof typeof zh, string> = {
  keywordPlaceholder: 'User ID / nickname / phone / share code',
  guildId: 'Guild ID',
  roomTitle: 'Room Title',
  roomType: 'Room Type',
  categoryHot: 'Hot',
  categoryGame: 'Game',
  categoryPrivate: 'Private',
  fetchFailed: 'Failed to load recycle bin',
  onShelfConfirm: 'Publish live room {id}? It will appear in the App list again.',
  onShelfSuccess: 'Published successfully',
  onShelfFailed: 'Failed to publish',
}

export const liveRoomRecycleBinMessages = definePageMessagesFromEn(zh, en)

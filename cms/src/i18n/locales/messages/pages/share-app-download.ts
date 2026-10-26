import {definePageMessagesFromEn} from './_define'

const zh = {
  downloadUrl: '当前 APK 下载链接',
  downloadUrlEmpty: '（暂无：请先在「官网部署」上传 APK）',
  downloadUrlTip:
      '只读。取自官网部署目录 assets 下体积最大的 .apk（如 web.saralive.net/assets/…）。保存域名/部署 zip/上传 APK 后会自动刷新内存。',
  shareText: '分享文案',
  shareTextPlaceholder: '邀请好友一起来玩…',
  shareTextTip:
      'App getUserInfo.shareCode 会返回本字段并带上邀请码。可用占位符 {inviteCode}、{downloadUrl}；未写 {inviteCode} 时会在文案后换行追加邀请码。',
  shareTextMaxLength: '分享文案最多 4000 字符',
  lastUpdated: '最近更新',
  fetchCfgFailed: '加载配置失败',
  saveCfgFailed: '保存配置失败',
  saveSuccess: '保存成功',
}

export const shareAppDownloadMessages = definePageMessagesFromEn(zh, {
  downloadUrl: 'Current APK download URL',
  downloadUrlEmpty: '(None yet — upload an APK on Official Site Deploy)',
  downloadUrlTip:
      'Read-only. Largest .apk under official-site assets (e.g. web.saralive.net/assets/…). Refreshed after domain save / zip deploy / APK upload.',
  shareText: 'Share text',
  shareTextPlaceholder: 'Invite friends to play…',
  shareTextTip:
      'Returned as getUserInfo.shareCode with invite code. Placeholders: {inviteCode}, {downloadUrl}. If {inviteCode} is omitted, invite code is appended on a new line.',
  shareTextMaxLength: 'Share text max 4000 characters',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveCfgFailed: 'Failed to save config',
  saveSuccess: 'Saved successfully',
})

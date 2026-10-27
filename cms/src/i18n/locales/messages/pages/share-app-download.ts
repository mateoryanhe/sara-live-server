import {definePageMessagesFromEn} from './_define'

const zh = {
  downloadUrl: '当前 APK 直链',
  downloadUrlEmpty: '（暂无：请先在「官网部署」上传 APK）',
  downloadUrlTip:
      '只读。取自官网部署目录 assets 下体积最大的 .apk。短链点击后会 302 到这里。',
  shortUrl: '分享短链',
  shortUrlEmpty: '（暂无：需先有 APK 且配置官网域名）',
  shortUrlTip: '只读。固定路径 /dl，302 到上方直链；App 分享内容第 3 行使用此地址。',
  shareText: '分享文案',
  shareTextPlaceholder: '邀请好友一起来玩…',
  shareTextTip:
      'App getUserInfo.shareCode 固定三行：①本字段 ②Invite Code: 邀请码 ③上方短链。无需再写占位符。',
  shareTextMaxLength: '分享文案最多 4000 字符',
  previewShareCode: '下发给 App 的最终文案',
  previewShareCodeEmpty: '（保存文案后可预览）',
  previewShareCodeTip:
      '与 App 实际下发规则一致；邀请码位置用 YOUR_INVITE_CODE 占位。可复制后自行替换为真实邀请码。',
  copyPreview: '复制',
  copyPreviewSuccess: '已复制到剪贴板',
  copyPreviewFailed: '复制失败',
  lastUpdated: '最近更新',
  fetchCfgFailed: '加载配置失败',
  saveCfgFailed: '保存配置失败',
  saveSuccess: '保存成功',
}

export const shareAppDownloadMessages = definePageMessagesFromEn(zh, {
  downloadUrl: 'Direct APK URL',
  downloadUrlEmpty: '(None yet — upload an APK on Official Site Deploy)',
  downloadUrlTip:
      'Read-only. Largest .apk under official-site assets. The short link 302-redirects here.',
  shortUrl: 'Share short link',
  shortUrlEmpty: '(None — need an APK and official-site domain)',
  shortUrlTip: 'Read-only. Fixed path /dl → 302 to the direct URL; used as line 3 of App share text.',
  shareText: 'Share text',
  shareTextPlaceholder: 'Invite friends to play…',
  shareTextTip:
      'getUserInfo.shareCode is always: (1) this text (2) Invite Code: {code} (3) short link. No placeholders needed.',
  shareTextMaxLength: 'Share text max 4000 characters',
  previewShareCode: 'Final text sent to App',
  previewShareCodeEmpty: '(Save share text to preview)',
  previewShareCodeTip:
      'Same rules as App. Invite code uses YOUR_INVITE_CODE placeholder. Copy and replace with a real code if needed.',
  copyPreview: 'Copy',
  copyPreviewSuccess: 'Copied to clipboard',
  copyPreviewFailed: 'Copy failed',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveCfgFailed: 'Failed to save config',
  saveSuccess: 'Saved successfully',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: '说明',
  noticeLine1: '本页仅配置游戏活动入口：按钮图标与游戏列表地址。',
  noticeLine2: '保存后立即入库；App 下发与跳转逻辑后续再接，当前仅做配置管理。',
  buttonIcon: '游戏按钮 Icon',
  buttonIconPlaceholder: 'CDN 图片 URL，可上传或粘贴',
  buttonIconTip: '用于前端展示的游戏活动入口图标',
  clickUploadIcon: '上传图标',
  removeIcon: '清除图标',
  gameListUrl: '游戏列表 URL',
  gameListUrlPlaceholder: 'https://example.com/games',
  gameListUrlTip: '点击入口后打开的游戏列表地址',
  lastUpdated: '最近更新',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功，游戏活动配置已生效',
  saveFailed: '保存失败',
  uploadFailed: '图标上传失败',
}

const en = {
  noticeTitle: 'Notice',
  noticeLine1: 'Configure the game activity entry: button icon and game list URL.',
  noticeLine2: 'Saving writes to the database immediately. App delivery and navigation will be wired later; this page is config-only for now.',
  buttonIcon: 'Game button icon',
  buttonIconPlaceholder: 'CDN image URL; upload or paste',
  buttonIconTip: 'Icon shown for the game activity entry',
  clickUploadIcon: 'Upload icon',
  removeIcon: 'Clear icon',
  gameListUrl: 'Game list URL',
  gameListUrlPlaceholder: 'https://example.com/games',
  gameListUrlTip: 'URL opened when the entry is tapped',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved. Game activity config is active',
  saveFailed: 'Save failed',
  uploadFailed: 'Icon upload failed',
}

export const gameActivityMessages = definePageMessagesFromEn(zh, en)

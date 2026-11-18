import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: '说明',
  noticeLine1: '此处配置全局 AppsFlyer 密钥，与包名配置无关。',
  noticeLine2: 'App 通过 POST /appsFlyer/getCfgForApp 读取内存缓存，请求头 X-Platform 须为 android 或 ios。',
  noticeLine3: 'iOS 返回 appsFlyerDevKey 与 appsFlyerAppId；Android 仅返回 appsFlyerDevKey。',
  tabAndroid: 'Android',
  tabIos: 'iOS',
  androidDevKey: 'Dev Key',
  iosDevKey: 'Dev Key',
  iosAppId: 'iOS App ID',
  androidDevKeyPlaceholder: 'Android AppsFlyer Dev Key',
  iosDevKeyPlaceholder: 'iOS AppsFlyer Dev Key',
  iosAppIdPlaceholder: 'App Store 数字 ID',
  lastUpdated: '最近更新',
  saveSuccess: 'AppsFlyer 配置已保存',
  needOneField: '请至少填写一个平台的 Dev Key',
} as const

export const appsFlyerCfgMessages = definePageMessagesFromEn(zh, {
  noticeTitle: 'Notes',
  noticeLine1: 'Global AppsFlyer keys here; separate from per-package settings.',
  noticeLine2: 'App calls POST /appsFlyer/getCfgForApp (cached). Set header X-Platform to android or ios.',
  noticeLine3: 'iOS returns appsFlyerDevKey and appsFlyerAppId; Android returns appsFlyerDevKey only.',
  tabAndroid: 'Android',
  tabIos: 'iOS',
  androidDevKey: 'Dev Key',
  iosDevKey: 'Dev Key',
  iosAppId: 'iOS App ID',
  androidDevKeyPlaceholder: 'Android AppsFlyer Dev Key',
  iosDevKeyPlaceholder: 'iOS AppsFlyer Dev Key',
  iosAppIdPlaceholder: 'App Store numeric ID',
  lastUpdated: 'Last updated',
  saveSuccess: 'AppsFlyer settings saved',
  needOneField: 'Enter at least one platform Dev Key',
})
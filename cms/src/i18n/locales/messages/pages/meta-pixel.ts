import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: 'Meta Pixel 服务端上报(Conversions API)',
  noticeLine1: '客户端无需接入 Pixel JS；Go 在「用户注册成功」「普通用户充值到账」后向 Meta 上报。',
  noticeLine2: '注册事件：CompleteRegistration；充值事件：Purchase(USD 为订单配置美金 price)。',
  noticeLine3: '币商充值、充值白名单虚拟美金不上报。',
  enabled: '启用上报',
  pixelId: 'Pixel ID',
  pixelIdPlaceholder: 'Events Manager 中的 Pixel 编号',
  accessToken: 'Access Token',
  accessTokenPlaceholder: 'Conversions API 用访问令牌',
  accessTokenRequired: '启用时 Access Token 不能为空',
  pixelIdRequired: '启用时 Pixel ID 不能为空',
  testEventCode: '测试事件代码',
  testEventCodePlaceholder: '可选，Events Manager 测试用',
  lastUpdated: '最近更新',
  saveSuccess: '保存成功，Meta Pixel 服务端上报配置已生效',
} as const

export const metaPixelMessages = definePageMessagesFromEn(zh, {
  noticeTitle: 'Meta Pixel server-side (Conversions API)',
  noticeLine1: 'No client Pixel JS; Go sends events after registration and normal-user recharge success.',
  noticeLine2: 'CompleteRegistration on sign-up; Purchase (USD order price) on recharge.',
  noticeLine3: 'Coin-merchant and whitelist virtual USD recharges are excluded.',
  enabled: 'Enable reporting',
  pixelId: 'Pixel ID',
  pixelIdPlaceholder: 'Pixel ID from Events Manager',
  accessToken: 'Access Token',
  accessTokenPlaceholder: 'Conversions API access token',
  accessTokenRequired: 'Access Token is required when enabled',
  pixelIdRequired: 'Pixel ID is required when enabled',
  testEventCode: 'Test event code',
  testEventCodePlaceholder: 'Optional, for Events Manager test mode',
  lastUpdated: 'Last updated',
  saveSuccess: 'Saved. Meta Pixel server-side config is active',
})

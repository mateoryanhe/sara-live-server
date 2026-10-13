import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: 'Meta Pixel 服务端上报(Conversions API)',
  noticeLine1: '客户端无需接入 Pixel JS；Go 在「用户注册成功」「普通用户充值到账」后向 Meta 上报。',
  noticeLine2: '注册 CompleteRegistration；充值 Purchase 在**用户下单成功**时上报(非入账)，value=档位 USD price；user_data 同前(充值无 fbp/fbc)。',
  noticeLine3: '经销商充值、充值白名单虚拟美金不上报。',
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
  noticeLine2: 'CompleteRegistration on sign-up; Purchase on order created (not gold credited), value=USD price; same user_data rules (no fbp/fbc on purchase).',
  noticeLine3: 'Dealer and whitelist virtual USD recharges are excluded.',
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

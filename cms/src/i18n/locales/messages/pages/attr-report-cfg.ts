import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: '归因数据上报配置',
  noticeLine1: 'Base API 只填主机根地址；服务端自动拼接 /tj/api/accept.php。密钥对应请求体 api_token。',
  noticeLine2: '示例主机：https://cdntk.showsme.click → 实际上报 https://cdntk.showsme.click/tj/api/accept.php',
  baseApi: 'Base API',
  baseApiPlaceholder: '如 https://cdntk.showsme.click（勿填 /tj/api/accept.php）',
  baseApiRequired: '请填写 Base API',
  secretKey: '密钥',
  secretKeyPlaceholder: '对应 accept 接口 api_token',
  secretKeyRequired: '请填写密钥',
  lastUpdated: '最近更新',
  saveSuccess: '保存成功，归因数据上报配置已生效',
} as const

export const attrReportCfgMessages = definePageMessagesFromEn(zh, {
  noticeTitle: 'Attribution report config',
  noticeLine1: 'Base API is host root only; server appends /tj/api/accept.php. Secret maps to body api_token.',
  noticeLine2: 'Example host: https://cdntk.showsme.click → posts to https://cdntk.showsme.click/tj/api/accept.php',
  baseApi: 'Base API',
  baseApiPlaceholder: 'e.g. https://cdntk.showsme.click (do not include /tj/api/accept.php)',
  baseApiRequired: 'Enter Base API',
  secretKey: 'Secret Key',
  secretKeyPlaceholder: 'Maps to accept api_token',
  secretKeyRequired: 'Enter secret key',
  lastUpdated: 'Last updated',
  saveSuccess: 'Saved. Attribution report config is active',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: 'CMS 域名映射',
  noticeBody: '保存后写入数据库并立即刷新当前服务的域名注册表，无需重启。访问配置域名时直接返回 CMS 静态页面。',
  domain: 'CMS 域名',
  domainPlaceholder: '例如 lzg.saralive.net',
  domainTip: '只填写一个 Host，不含 http/https、端口、路径或逗号。',
  domainRequired: '请输入 CMS 域名',
  domainInvalid: '域名格式无效，只允许填写一个 Host，不能使用逗号',
  urlPrefix: '内部路径前缀',
  urlPrefixTip: '固定为 /cms，用于标识该站点映射，不需要手动修改。',
  root: 'CMS 部署目录',
  rootPlaceholder: '/home/ec2-user/cdn/cms',
  rootTip: '填写服务器上的绝对目录；目录不存在时服务端会尝试创建。',
  rootRequired: '请输入 CMS 部署目录',
  rootAbsolute: '部署目录必须是绝对路径',
  updatedAt: '最近更新时间',
  save: '保存并刷新映射',
  reload: '重新加载',
  fetchFailed: '获取 CMS 域名映射失败',
  saveSuccess: 'CMS 域名映射已保存并生效',
  saveFailed: '保存 CMS 域名映射失败',
}

const en = {
  noticeTitle: 'CMS Domain Mapping',
  noticeBody: 'Saving writes the mapping to the database and refreshes the current server registry immediately. No restart is required.',
  domain: 'CMS domain',
  domainPlaceholder: 'e.g. lzg.saralive.net',
  domainTip: 'Enter one host without a scheme, port, path, or comma.',
  domainRequired: 'Enter the CMS domain',
  domainInvalid: 'Invalid domain. Only one host is allowed; commas are not allowed.',
  urlPrefix: 'Internal URL prefix',
  urlPrefixTip: 'Fixed at /cms to identify this site mapping.',
  root: 'CMS deployment directory',
  rootPlaceholder: '/home/ec2-user/cdn/cms',
  rootTip: 'Enter an absolute server path. The server will try to create it when it does not exist.',
  rootRequired: 'Enter the CMS deployment directory',
  rootAbsolute: 'The deployment directory must be an absolute path',
  updatedAt: 'Last updated',
  save: 'Save and refresh',
  reload: 'Reload',
  fetchFailed: 'Failed to load the CMS domain mapping',
  saveSuccess: 'CMS domain mapping saved and applied',
  saveFailed: 'Failed to save the CMS domain mapping',
}

export const cmsDomainMappingMessages = definePageMessagesFromEn(zh, en)

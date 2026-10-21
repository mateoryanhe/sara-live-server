import {definePageMessagesFromEn} from './_define'

const zh = {
  noticeTitle: 'CMS 备用域名',
  noticeBody: '在此追加更多 Host，与「CMS域名映射」共用同一静态目录；主域名请在 CMS域名映射 页配置。',
  primaryDomain: '主 CMS 域名',
  deployRoot: 'CMS 部署目录',
  domainPlaceholder: '例如 backup.example.com',
  domainRequired: '请输入域名',
  domainInvalid: '域名格式无效，只允许一个 Host',
  addDomain: '添加域名',
  columnDomain: '域名',
  columnUpdatedAt: '更新时间',
  empty: '暂无备用域名',
  fetchFailed: '加载备用域名失败',
  addSuccess: '备用域名已添加并生效',
  addFailed: '添加备用域名失败',
  deleteConfirm: '确定删除备用域名 {domain}？',
  deleteSuccess: '已删除',
  deleteFailed: '删除失败',
}

const en = {
  noticeTitle: 'CMS extra domains',
  noticeBody: 'Add more hosts that serve the same CMS static files as the primary mapping on the CMS Domain Mapping page.',
  primaryDomain: 'Primary CMS domain',
  deployRoot: 'CMS deploy directory',
  domainPlaceholder: 'e.g. backup.example.com',
  domainRequired: 'Enter a domain',
  domainInvalid: 'Invalid domain; use a single host only',
  addDomain: 'Add domain',
  columnDomain: 'Domain',
  columnUpdatedAt: 'Updated at',
  empty: 'No extra domains',
  fetchFailed: 'Failed to load extra domains',
  addSuccess: 'Extra domain added',
  addFailed: 'Failed to add domain',
  deleteConfirm: 'Remove extra domain {domain}?',
  deleteSuccess: 'Removed',
  deleteFailed: 'Failed to remove',
}

export const cmsDomainAliasMessages = definePageMessagesFromEn(zh, en)

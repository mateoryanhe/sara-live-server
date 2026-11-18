import {definePageMessages} from './_define'

const zh = {
  title: '确认上传并部署',
  confirmBtn: '确认部署',
  cancelBtn: '取消',
  intro: '即将上传以下 ZIP 包，服务端会解压并写入部署目录（同名文件会被覆盖）：',
  site: '站点：{site}',
  file: 'ZIP 文件：{file}',
  size: '文件大小：{size}',
  domain: '访问域名：{domain}',
  deployPath: '部署目录：{path}',
  urlPrefix: '静态访问前缀：{prefix}',
  chunked: '上传方式：8 MB 分片上传，全部分片校验通过后合并并解压',
  direct: '上传方式：整包上传后服务端解压',
}

const en = {
  title: 'Confirm upload and deploy',
  confirmBtn: 'Deploy',
  cancelBtn: 'Cancel',
  intro: 'The following ZIP will be uploaded, extracted, and written to the deploy directory (existing files with the same names will be overwritten):',
  site: 'Site: {site}',
  file: 'ZIP file: {file}',
  size: 'File size: {size}',
  domain: 'Domain: {domain}',
  deployPath: 'Deploy directory: {path}',
  urlPrefix: 'Static URL prefix: {prefix}',
  chunked: 'Upload: 8 MB chunks; merged and extracted on the server after all chunks are verified',
  direct: 'Upload: single package upload, then server-side extraction',
}

const es = {...en}
const pt = {...en}
const hi = {...en}
const id = {...en}

export const deployZipConfirmMessages = definePageMessages(zh, en, es, pt, hi, id)
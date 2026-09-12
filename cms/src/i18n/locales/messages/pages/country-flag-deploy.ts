import {definePageMessagesFromEn} from './_define'

const zh = {
  version: '当前版本',
  urlPrefix: '静态访问前缀',
  deployPath: '资源根目录',
  lastUpdated: '最近更新',
  uploadTitle: '上传国旗 ZIP',
  uploadTip: '将 pub-tool/country-flags/flags 打成 zip 上传。服务端生成新 version 目录、写入数据库，并删除旧 version 目录。',
  dragTip: '将 zip 文件拖到此处，或点击选择',
  fileTip: '仅支持 .zip；包内需含两位简码 png（如 id.png），可放在任意子目录',
  deployBtn: '上传并发布',
  deploySuccess: '发布成功',
  deployFailed: '发布失败',
  fetchInfoFailed: '获取部署信息失败',
  zipOnly: '仅支持上传 zip 文件',
  singleFileOnly: '一次只能上传一个 zip 文件',
  resultVersion: '新版本：{version}',
  resultPath: '版本目录：{path}',
  resultFiles: '写入 PNG：{count} 个',
  resultRemoved: '清理旧版本目录：{count} 个',
  resultUrl: '访问前缀：{url}',
  emptyVersion: '尚未发布',
}

const en = {
  version: 'Current version',
  urlPrefix: 'Static URL prefix',
  deployPath: 'Resource root',
  lastUpdated: 'Last updated',
  uploadTitle: 'Upload flag ZIP',
  uploadTip: 'Zip pub-tool/country-flags/flags and upload. Server creates a new version dir, saves it to DB, and deletes old version dirs.',
  dragTip: 'Drop zip here or click to select',
  fileTip: 'Only .zip. Include 2-letter png files (e.g. id.png); nested folders are OK.',
  deployBtn: 'Upload & publish',
  deploySuccess: 'Published successfully',
  deployFailed: 'Publish failed',
  fetchInfoFailed: 'Failed to load deploy info',
  zipOnly: 'Only zip files are allowed',
  singleFileOnly: 'Only one zip file at a time',
  resultVersion: 'New version: {version}',
  resultPath: 'Version path: {path}',
  resultFiles: 'PNG files written: {count}',
  resultRemoved: 'Old version dirs removed: {count}',
  resultUrl: 'URL prefix: {url}',
  emptyVersion: 'Not published yet',
}

export const countryFlagDeployMessages = definePageMessagesFromEn(zh, en)

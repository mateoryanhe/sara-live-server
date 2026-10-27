import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: '探活 gw 连续检测不到直播 Go 进程并触发冷启后写入。表 process_down_logs。',
  host: '主机',
  hostPlaceholder: '主机名（可选）',
  dateRange: '判定时间',
  dateRangeSeparator: '至',
  appName: '进程名',
  missingSince: '开始缺失',
  detectedAt: '判定挂掉',
  coldStartOK: '冷启结果',
  coldStartSuccess: '成功',
  coldStartFailed: '失败',
  newPid: '新 PID',
  detail: '说明',
  createdAt: '入库时间',
  fetchFailed: '加载失败',
}

export const processDownLogMessages = definePageMessagesFromEn(zh, {
  pageDesc: 'Written by live-hot-restart-gw after continuous process miss triggers cold start. Table process_down_logs.',
  host: 'Host',
  hostPlaceholder: 'Hostname (optional)',
  dateRange: 'Detected at',
  dateRangeSeparator: 'to',
  appName: 'App name',
  missingSince: 'Missing since',
  detectedAt: 'Detected down',
  coldStartOK: 'Cold start',
  coldStartSuccess: 'OK',
  coldStartFailed: 'Failed',
  newPid: 'New PID',
  detail: 'Detail',
  createdAt: 'Created at',
  fetchFailed: 'Failed to load',
})

import {definePageMessagesFromEn} from './_define'

const zh = {
  tipTitle: '说明',
  tipLine1: '观众 App 登录成功（getUserInfo）时，按 X-Language / Accept-Language 选语种群发私信；同一观众每天最多一次。',
  tipLine2: '从当前 WebSocket 在线的普通/高级主播中，按 1v1 累计通话次数从高到低选取最多 N 人；排队依次私信，每条（含第一条）发送前在 [最小, 最大] 秒内随机等待（默认 30–30）。',
  tipLine3: '招呼文案固定为服务端 Go 内置（default_greet_lines.go，各语种各 30 条），发送时随机一条；本页仅配置开关、主播人数与排队间隔。',
  enabled: '开启',
  anchorCount: '打招呼主播人数',
  anchorCountHint: '1–20，默认 5',
  queueTitle: '排队招呼间隔（秒）',
  intervalMinSec: '间隔下限',
  intervalMaxSec: '间隔上限',
  intervalHint: '0–3600；相等则为固定间隔',
  intervalRangeInvalid: '间隔上限不能小于下限',
  fetchFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
} as const

const en: Record<keyof typeof zh, string> = {
  tipTitle: 'Notes',
  tipLine1: 'On App login (getUserInfo), locale from X-Language / Accept-Language; once per viewer per day.',
  tipLine2: 'Pick up to N online anchors by descending 1v1 call count; each PM (including the first) waits a random delay in [min, max] seconds before send (default 30–30).',
  tipLine3: 'Greeting copy is fixed to built-in Go lines (default_greet_lines.go, 30 per locale), one picked at random per send. This page only toggles enable, anchor count, and queue delay.',
  enabled: 'Enabled',
  anchorCount: 'Number of anchors to greet',
  anchorCountHint: '1–20, default 5',
  queueTitle: 'Queue delay (seconds)',
  intervalMinSec: 'Min interval',
  intervalMaxSec: 'Max interval',
  intervalHint: '0–3600; equal min/max = fixed delay',
  intervalRangeInvalid: 'Max interval cannot be less than min',
  fetchFailed: 'Failed to load config',
  saveSuccess: 'Saved',
  saveFailed: 'Failed to save',
}

export const audienceOnlineAnchorGreetMessages = definePageMessagesFromEn(zh, en)
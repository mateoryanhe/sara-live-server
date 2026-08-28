import {definePageMessages} from './_define'

const zh = {
  keywordPlaceholder: '流水ID / 主播ID / 昵称',
  fetchFailed: '获取本周流水失败',
  weekRangeHint: '当前查询范围：本周 {start} 至 {end}，仅展示未结算数据',
  weeklyUnsettledTotalIncome: '本周未结算总收益',
}

const en = {
  keywordPlaceholder: 'Flow ID / anchor ID / nickname',
  fetchFailed: 'Failed to load weekly unsettled live flow',
  weekRangeHint: 'Current range: {start} to {end} (this week, unsettled only)',
  weeklyUnsettledTotalIncome: 'Weekly Unsettled Total Income',
}

const es = {
  keywordPlaceholder: 'ID de flujo / ID de anfitrión / apodo',
  fetchFailed: 'Error al cargar flujo semanal no liquidado',
  weekRangeHint: 'Rango actual: {start} a {end} (esta semana, solo no liquidado)',
  weeklyUnsettledTotalIncome: 'Ingresos totales no liquidados de la semana',
}

const pt = {
  keywordPlaceholder: 'ID do fluxo / ID do anfitrião / apelido',
  fetchFailed: 'Falha ao carregar fluxo semanal não liquidado',
  weekRangeHint: 'Intervalo atual: {start} a {end} (esta semana, apenas não liquidado)',
  weeklyUnsettledTotalIncome: 'Receita total não liquidada da semana',
}

const hi = {
  keywordPlaceholder: 'फ्लो ID / एंकर ID / उपनाम',
  fetchFailed: 'साप्ताहिक अनिपटारा फ्लो लोड विफल',
  weekRangeHint: 'वर्तमान सीमा: {start} से {end} (इस सप्ताह, केवल अनिपटारा)',
  weeklyUnsettledTotalIncome: 'साप्ताहिक अनिपटारा कुल आय',
}

const id = {
  keywordPlaceholder: 'ID aliran / ID anchor / nickname',
  fetchFailed: 'Gagal memuat aliran mingguan belum diselesaikan',
  weekRangeHint: 'Rentang saat ini: {start} hingga {end} (minggu ini, hanya belum diselesaikan)',
  weeklyUnsettledTotalIncome: 'Total pendapatan belum diselesaikan minggu ini',
}

export const liveWeeklyUnsettledLiveListMessages = definePageMessages(zh, en, es, pt, hi, id)

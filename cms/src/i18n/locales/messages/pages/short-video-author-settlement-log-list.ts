import {definePageMessages} from './_define'

const zh = {
  logId: '日志ID',
  unsettledIncome: '结算流水',
  settlementDiamond: '到账钻石',
  anchorSharePercent: '主播分佣比例(%)',
  time: '结算时间',
  dateRangeSeparator: '至',
  fetchFailed: '获取短视频作者结算日志失败',
} as const

const en: Record<keyof typeof zh, string> = {
  logId: 'Log ID',
  unsettledIncome: 'Settled Income',
  settlementDiamond: 'Settlement Diamond',
  anchorSharePercent: 'Anchor Share (%)',
  time: 'Settlement Time',
  dateRangeSeparator: 'to',
  fetchFailed: 'Failed to load short video author settlement logs',
}

const es: Record<keyof typeof zh, string> = {
  logId: 'ID de registro',
  unsettledIncome: 'Ingresos liquidados',
  settlementDiamond: 'Diamantes liquidados',
  anchorSharePercent: 'Reparto ancla (%)',
  time: 'Hora de liquidación',
  dateRangeSeparator: 'a',
  fetchFailed: 'Error al cargar registros de liquidación de autores de video corto',
}

const pt: Record<keyof typeof zh, string> = {
  logId: 'ID do registro',
  unsettledIncome: 'Receita liquidada',
  settlementDiamond: 'Diamantes liquidados',
  anchorSharePercent: 'Reparto âncora (%)',
  time: 'Hora da liquidação',
  dateRangeSeparator: 'até',
  fetchFailed: 'Falha ao carregar registros de liquidação de autores de vídeo curto',
}

const hi: Record<keyof typeof zh, string> = {
  logId: 'लॉग ID',
  unsettledIncome: 'सेटल की गई आय',
  settlementDiamond: 'सेटलमेंट डायमंड',
  anchorSharePercent: 'एंकर हिस्सा (%)',
  time: 'सेटलमेंट समय',
  dateRangeSeparator: 'से',
  fetchFailed: 'शॉर्ट वीडियो लेखक सेटलमेंट लॉग लोड करने में विफल',
}

const id: Record<keyof typeof zh, string> = {
  logId: 'ID Log',
  unsettledIncome: 'Pendapatan diselesaikan',
  settlementDiamond: 'Diamond settlement',
  anchorSharePercent: 'Bagi hasil anchor (%)',
  time: 'Waktu settlement',
  dateRangeSeparator: 'sampai',
  fetchFailed: 'Gagal memuat log settlement penulis video pendek',
}

export const shortVideoAuthorSettlementLogListMessages = definePageMessages(zh, en, es, pt, hi, id)

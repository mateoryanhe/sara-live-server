import {definePageMessages} from './_define'

const zh = {
  logId: '流水ID',
  userNickname: '用户昵称',
  changeType: '变动类型',
  goldChange: '金币变动',
  diamondChange: '钻石变动',
  beforeChange: '变动前',
  afterChange: '变动后',
  reason: '原因',
  time: '时间',
  actionIncrease: '增加',
  actionDecrease: '减少',
  userIdPlaceholder: '请输入用户ID',
  fetchFailed: '获取流水失败',
} as const

const en: Record<keyof typeof zh, string> = {
  logId: 'Log ID',
  userNickname: 'Nickname',
  changeType: 'Change Type',
  goldChange: 'Gold Change',
  diamondChange: 'Diamond Change',
  beforeChange: 'Before',
  afterChange: 'After',
  reason: 'Reason',
  time: 'Time',
  actionIncrease: 'Increase',
  actionDecrease: 'Decrease',
  userIdPlaceholder: 'Enter user ID',
  fetchFailed: 'Failed to load currency logs',
}

const es: Record<keyof typeof zh, string> = {
  logId: 'ID de registro',
  userNickname: 'Apodo',
  changeType: 'Tipo de cambio',
  goldChange: 'Cambio de oro',
  diamondChange: 'Cambio de diamantes',
  beforeChange: 'Antes',
  afterChange: 'Después',
  reason: 'Motivo',
  time: 'Hora',
  actionIncrease: 'Aumento',
  actionDecrease: 'Disminución',
  userIdPlaceholder: 'Introduzca ID de usuario',
  fetchFailed: 'Error al cargar registros',
}

const pt: Record<keyof typeof zh, string> = {
  logId: 'ID do registro',
  userNickname: 'Apelido',
  changeType: 'Tipo de alteração',
  goldChange: 'Alteração de ouro',
  diamondChange: 'Alteração de diamantes',
  beforeChange: 'Antes',
  afterChange: 'Depois',
  reason: 'Motivo',
  time: 'Hora',
  actionIncrease: 'Aumento',
  actionDecrease: 'Redução',
  userIdPlaceholder: 'Digite o ID do usuário',
  fetchFailed: 'Falha ao carregar registros',
}

const hi: Record<keyof typeof zh, string> = {
  logId: 'लॉग ID',
  userNickname: 'उपनाम',
  changeType: 'परिवर्तन प्रकार',
  goldChange: 'गोल्ड परिवर्तन',
  diamondChange: 'डायमंड परिवर्तन',
  beforeChange: 'पहले',
  afterChange: 'बाद',
  reason: 'कारण',
  time: 'समय',
  actionIncrease: 'बढ़ोतरी',
  actionDecrease: 'कमी',
  userIdPlaceholder: 'उपयोगकर्ता ID दर्ज करें',
  fetchFailed: 'लॉग लोड विफल',
}

const id: Record<keyof typeof zh, string> = {
  logId: 'ID log',
  userNickname: 'Nama panggilan',
  changeType: 'Tipe perubahan',
  goldChange: 'Perubahan emas',
  diamondChange: 'Perubahan berlian',
  beforeChange: 'Sebelum',
  afterChange: 'Sesudah',
  reason: 'Alasan',
  time: 'Waktu',
  actionIncrease: 'Tambah',
  actionDecrease: 'Kurang',
  userIdPlaceholder: 'Masukkan ID pengguna',
  fetchFailed: 'Gagal memuat log',
}

export const currencyLogListMessages = definePageMessages(zh, en, es, pt, hi, id)

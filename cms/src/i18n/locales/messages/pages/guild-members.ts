import {definePageMessages} from './_define'

const zh = {
  title: '工会成员',
  titleWithName: '工会成员 - {name}',
  back: '返回',
  noMembers: '暂无成员',
  fetchFailed: '获取工会成员失败',
} as const

export const guildMembersMessages = definePageMessages(
  zh,
  {
    title: 'Guild Members',
    titleWithName: 'Guild Members - {name}',
    back: 'Back',
    noMembers: 'No members',
    fetchFailed: 'Failed to load guild members',
  },
  {
    title: 'Miembros del gremio',
    titleWithName: 'Miembros del gremio - {name}',
    back: 'Volver',
    noMembers: 'Sin miembros',
    fetchFailed: 'Error al cargar miembros',
  },
  {
    title: 'Membros da guilda',
    titleWithName: 'Membros da guilda - {name}',
    back: 'Voltar',
    noMembers: 'Sem membros',
    fetchFailed: 'Falha ao carregar membros',
  },
  {
    title: 'गिल्ड सदस्य',
    titleWithName: 'गिल्ड सदस्य - {name}',
    back: 'वापस',
    noMembers: 'कोई सदस्य नहीं',
    fetchFailed: 'सदस्य लोड विफल',
  },
  {
    title: 'Anggota Guild',
    titleWithName: 'Anggota Guild - {name}',
    back: 'Kembali',
    noMembers: 'Tidak ada anggota',
    fetchFailed: 'Gagal memuat anggota',
  },
)

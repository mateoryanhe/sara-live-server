import {definePageMessages} from './_define'

const zh = {
  back: '返回',
  savePermissions: '保存权限配置',
  expandAll: '全部展开',
  permissionHint: '勾选页面可授予整页权限；展开后可单独勾选各按钮',
  saveSuccess: '已保存 {count} 项权限',
  saveFailed: '保存权限配置失败',
} as const

const en: Record<keyof typeof zh, string> = {
  back: 'Back',
  savePermissions: 'Save Permissions',
  expandAll: 'Expand All',
  permissionHint: 'Check a page for full access; expand to select individual buttons',
  saveSuccess: 'Saved {count} permissions',
  saveFailed: 'Failed to save permissions',
}

const es: Record<keyof typeof zh, string> = {
  back: 'Volver',
  savePermissions: 'Guardar permisos',
  expandAll: 'Expandir todo',
  permissionHint: 'Marque una página para acceso completo; expanda para botones individuales',
  saveSuccess: 'Guardados {count} permisos',
  saveFailed: 'Error al guardar permisos',
}

const pt: Record<keyof typeof zh, string> = {
  back: 'Voltar',
  savePermissions: 'Salvar permissões',
  expandAll: 'Expandir tudo',
  permissionHint: 'Marque uma página para acesso total; expanda para botões individuais',
  saveSuccess: 'Salvas {count} permissões',
  saveFailed: 'Falha ao salvar permissões',
}

const hi: Record<keyof typeof zh, string> = {
  back: 'वापस',
  savePermissions: 'अनुमतियाँ सहेजें',
  expandAll: 'सभी विस्तृत करें',
  permissionHint: 'पूर्ण पहुंच के लिए पृष्ठ चुनें; बटन के लिए विस्तार करें',
  saveSuccess: '{count} अनुमतियाँ सहेजी गईं',
  saveFailed: 'अनुमतियाँ सहेजने में विफल',
}

const id: Record<keyof typeof zh, string> = {
  back: 'Kembali',
  savePermissions: 'Simpan izin',
  expandAll: 'Buka semua',
  permissionHint: 'Centang halaman untuk akses penuh; perluas untuk tombol individual',
  saveSuccess: 'Tersimpan {count} izin',
  saveFailed: 'Gagal menyimpan izin',
}

export const moduleListMessages = definePageMessages(zh, en, es, pt, hi, id)

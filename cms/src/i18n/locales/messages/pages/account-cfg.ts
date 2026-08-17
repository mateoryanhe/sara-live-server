import {definePageMessages} from './_define'

const zh = {
  cancelAccountByCodeEnabled: '注销码销户(官网)',
  cancelAccountByCodeTip: '控制官网公开接口 POST /userInfo/cancelAccountByCode 是否可用，默认关闭',
  simulatorLoginEnabled: '模拟器登录',
  simulatorLoginTip: '关闭后，将根据 App 上报的 cpuModel 识别模拟器并拒绝登录/注册，默认开启',
  lastUpdated: '最近更新',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
}

const en = {
  cancelAccountByCodeEnabled: 'Cancel account by code (official site)',
  cancelAccountByCodeTip: 'Controls whether public API POST /userInfo/cancelAccountByCode is enabled; disabled by default',
  simulatorLoginEnabled: 'Simulator login',
  simulatorLoginTip: 'When off, devices whose reported cpuModel looks like an emulator/simulator cannot login/register; enabled by default',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved successfully',
  saveFailed: 'Save failed',
}

const es = {
  cancelAccountByCodeEnabled: 'Cancelar cuenta por código (sitio oficial)',
  cancelAccountByCodeTip: 'Controla si la API pública POST /userInfo/cancelAccountByCode está habilitada; desactivada por defecto',
  simulatorLoginEnabled: 'Inicio de sesión en simulador',
  simulatorLoginTip: 'Si está desactivado, se rechazan dispositivos cuyo cpuModel indique emulador/simulador; activado por defecto',
  lastUpdated: 'Última actualización',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccess: 'Guardado correctamente',
  saveFailed: 'Error al guardar',
}

const pt = {
  cancelAccountByCodeEnabled: 'Cancelar conta por código (site oficial)',
  cancelAccountByCodeTip: 'Controla se a API pública POST /userInfo/cancelAccountByCode está habilitada; desativada por padrão',
  simulatorLoginEnabled: 'Login no simulador',
  simulatorLoginTip: 'Se desativado, dispositivos cujo cpuModel indicar emulador/simulador serão bloqueados; ativado por padrão',
  lastUpdated: 'Última atualização',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccess: 'Salvo com sucesso',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  cancelAccountByCodeEnabled: 'कोड से खाता रद्द (आधिकारिक साइट)',
  cancelAccountByCodeTip: 'नियंत्रित करता है कि सार्वजनिक API POST /userInfo/cancelAccountByCode उपलब्ध है या नहीं; डिफ़ॉल्ट बंद',
  simulatorLoginEnabled: 'सिम्युलेटर लॉगिन',
  simulatorLoginTip: 'बंद होने पर cpuModel से पहचाने गए सिम्युलेटर/एमुलेटर लॉगिन/रजिस्टर नहीं कर सकते; डिफ़ॉल्ट खुला',
  lastUpdated: 'अंतिम अपडेट',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccess: 'सफलतापूर्वक सहेजा',
  saveFailed: 'सहेजना विफल',
}

const id = {
  cancelAccountByCodeEnabled: 'Batalkan akun via kode (situs resmi)',
  cancelAccountByCodeTip: 'Mengontrol apakah API publik POST /userInfo/cancelAccountByCode tersedia; default nonaktif',
  simulatorLoginEnabled: 'Login simulator',
  simulatorLoginTip: 'Jika dimatikan, perangkat dengan cpuModel mirip emulator/simulator akan ditolak login/daftar; default aktif',
  lastUpdated: 'Terakhir diperbarui',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const accountCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

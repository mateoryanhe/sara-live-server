import {definePageMessages} from './_define'

const zh = {
  envType: '环境类型',
  envTypeProd: '正式服',
  envTypeReview: '提审服',
  envTypeTest: '测试服',
  envTypeTip: '写入账号配置，App 通过系统配置接口读取：0正式服 / 1提审服 / 2测试服',
  cancelAccountByCodeEnabled: '注销码销户(官网)',
  cancelAccountByCodeTip: '控制官网公开接口 POST /userInfo/cancelAccountByCode 是否可用，默认关闭',
  blockSimulatorLogin: '拦截模拟器登录',
  blockSimulatorLoginTip: '开启后，将根据 App 上报的 cpuModel 识别模拟器并拒绝登录/注册，默认关闭(不拦截)',
  lastUpdated: '最近更新',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
}

const en = {
  envType: 'Environment type',
  envTypeProd: 'Production',
  envTypeReview: 'Review',
  envTypeTest: 'Test',
  envTypeTip: 'Stored in account config; App reads via system cfg: 0 prod / 1 review / 2 test',
  cancelAccountByCodeEnabled: 'Cancel account by code (official site)',
  cancelAccountByCodeTip: 'Controls whether public API POST /userInfo/cancelAccountByCode is enabled; disabled by default',
  blockSimulatorLogin: 'Block simulator login',
  blockSimulatorLoginTip: 'When on, devices whose reported cpuModel looks like an emulator/simulator cannot login/register; off by default (not blocked)',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved successfully',
  saveFailed: 'Save failed',
}

const es = {
  envType: 'Tipo de entorno',
  envTypeProd: 'Producción',
  envTypeReview: 'Revisión',
  envTypeTest: 'Prueba',
  envTypeTip: 'Se guarda en config de cuenta; App lo lee vía cfg del sistema: 0 prod / 1 revisión / 2 prueba',
  cancelAccountByCodeEnabled: 'Cancelar cuenta por código (sitio oficial)',
  cancelAccountByCodeTip: 'Controla si la API pública POST /userInfo/cancelAccountByCode está habilitada; desactivada por defecto',
  blockSimulatorLogin: 'Bloquear login en simulador',
  blockSimulatorLoginTip: 'Si está activado, se rechazan dispositivos cuyo cpuModel indique emulador/simulador; desactivado por defecto (sin bloqueo)',
  lastUpdated: 'Última actualización',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccess: 'Guardado correctamente',
  saveFailed: 'Error al guardar',
}

const pt = {
  envType: 'Tipo de ambiente',
  envTypeProd: 'Produção',
  envTypeReview: 'Revisão',
  envTypeTest: 'Teste',
  envTypeTip: 'Salvo na config de conta; App lê via cfg do sistema: 0 prod / 1 revisão / 2 teste',
  cancelAccountByCodeEnabled: 'Cancelar conta por código (site oficial)',
  cancelAccountByCodeTip: 'Controla se a API pública POST /userInfo/cancelAccountByCode está habilitada; desativada por padrão',
  blockSimulatorLogin: 'Bloquear login no simulador',
  blockSimulatorLoginTip: 'Se ativado, dispositivos cujo cpuModel indicar emulador/simulador serão bloqueados; desativado por padrão (sem bloqueio)',
  lastUpdated: 'Última atualização',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccess: 'Salvo com sucesso',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  envType: 'परिवेश प्रकार',
  envTypeProd: 'प्रोडक्शन',
  envTypeReview: 'रिव्यू',
  envTypeTest: 'टेस्ट',
  envTypeTip: 'अकाउंट कॉन्फ़िग में सेव; App सिस्टम cfg से पढ़े: 0 प्रोड / 1 रिव्यू / 2 टेस्ट',
  cancelAccountByCodeEnabled: 'कोड से खाता रद्द (आधिकारिक साइट)',
  cancelAccountByCodeTip: 'नियंत्रित करता है कि सार्वजनिक API POST /userInfo/cancelAccountByCode उपलब्ध है या नहीं; डिफ़ॉल्ट बंद',
  blockSimulatorLogin: 'सिम्युलेटर लॉगिन ब्लॉक',
  blockSimulatorLoginTip: 'चालू होने पर cpuModel से पहचाने गए सिम्युलेटर/एमुलेटर लॉगिन/रजिस्टर नहीं कर सकते; डिफ़ॉल्ट बंद (कोई ब्लॉक नहीं)',
  lastUpdated: 'अंतिम अपडेट',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccess: 'सफलतापूर्वक सहेजा',
  saveFailed: 'सहेजना विफल',
}

const id = {
  envType: 'Jenis lingkungan',
  envTypeProd: 'Produksi',
  envTypeReview: 'Review',
  envTypeTest: 'Tes',
  envTypeTip: 'Disimpan di config akun; App membaca via cfg sistem: 0 prod / 1 review / 2 tes',
  cancelAccountByCodeEnabled: 'Batalkan akun via kode (situs resmi)',
  cancelAccountByCodeTip: 'Mengontrol apakah API publik POST /userInfo/cancelAccountByCode tersedia; default nonaktif',
  blockSimulatorLogin: 'Blokir login simulator',
  blockSimulatorLoginTip: 'Jika aktif, perangkat dengan cpuModel mirip emulator/simulator akan ditolak login/daftar; default nonaktif (tidak diblokir)',
  lastUpdated: 'Terakhir diperbarui',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const accountCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

import {definePageMessages} from './_define'

const zh = {
  recentLoginLimit: '最近登录用户预热数量',
  recentLoginLimitTip: '服务启动时按 user_infos.last_login_time 倒序预热最近 N 个用户的缓存，默认 100；修改后需重启服务生效',
  lastUpdated: '最近更新',
  preloadCountRequired: '请输入预热数量',
  preloadCountRange: '预热数量需在 1-10000 之间',
  fetchCfgFailed: '获取配置失败',
  saveSuccessRestart: '保存成功，重启服务后生效',
  saveFailed: '保存失败',
}

const en = {
  recentLoginLimit: 'Recent login user preload count',
  recentLoginLimitTip: 'On startup, preload cache for the N most recently logged-in users (by user_infos.last_login_time), default 100. Restart required after changes.',
  lastUpdated: 'Last updated',
  preloadCountRequired: 'Please enter preload count',
  preloadCountRange: 'Preload count must be between 1 and 10000',
  fetchCfgFailed: 'Failed to load config',
  saveSuccessRestart: 'Saved successfully; restart service to apply',
  saveFailed: 'Save failed',
}

const es = {
  recentLoginLimit: 'Cantidad de precarga de usuarios recientes',
  recentLoginLimitTip: 'Al iniciar, precarga caché de los N usuarios con último login más reciente (user_infos.last_login_time), predeterminado 100. Requiere reinicio.',
  lastUpdated: 'Última actualización',
  preloadCountRequired: 'Introduce la cantidad de precarga',
  preloadCountRange: 'La cantidad debe estar entre 1 y 10000',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccessRestart: 'Guardado; reinicie el servicio para aplicar',
  saveFailed: 'Error al guardar',
}

const pt = {
  recentLoginLimit: 'Quantidade de pré-carga de usuários recentes',
  recentLoginLimitTip: 'Na inicialização, pré-carrega cache dos N usuários com login mais recente (user_infos.last_login_time), padrão 100. Reinício necessário.',
  lastUpdated: 'Última atualização',
  preloadCountRequired: 'Digite a quantidade de pré-carga',
  preloadCountRange: 'A quantidade deve estar entre 1 e 10000',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccessRestart: 'Salvo; reinicie o serviço para aplicar',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  recentLoginLimit: 'हाल के लॉगिन उपयोगकर्ता प्रीलोड संख्या',
  recentLoginLimitTip: 'स्टार्टअप पर user_infos.last_login_time के अनुसार हाल के N उपयोगकर्ताओं का कैश प्रीलोड, डिफ़ॉल्ट 100। बदलाव के बाद रीस्टार्ट आवश्यक।',
  lastUpdated: 'अंतिम अपडेट',
  preloadCountRequired: 'प्रीलोड संख्या दर्ज करें',
  preloadCountRange: 'प्रीलोड संख्या 1-10000 के बीच होनी चाहिए',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccessRestart: 'सहेजा गया; लागू करने के लिए सेवा रीस्टार्ट करें',
  saveFailed: 'सहेजना विफल',
}

const id = {
  recentLoginLimit: 'Jumlah preload pengguna login terbaru',
  recentLoginLimitTip: 'Saat startup, preload cache N pengguna login terbaru (user_infos.last_login_time), default 100. Perlu restart setelah perubahan.',
  lastUpdated: 'Terakhir diperbarui',
  preloadCountRequired: 'Masukkan jumlah preload',
  preloadCountRange: 'Jumlah preload harus antara 1 dan 10000',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccessRestart: 'Disimpan; restart layanan untuk menerapkan',
  saveFailed: 'Gagal menyimpan',
}

export const preloadCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

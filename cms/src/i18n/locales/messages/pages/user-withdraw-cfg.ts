import {definePageMessages} from './_define'

const zh = {
  enabled: '用户提现开关',
  enabledTip: '关闭后 App 无法申请提现、无法维护收款信息；已提交的提现单仍可在 CMS 处理。',
  dailyMaxCount: '每天提现次数',
  dailyMaxCountTip: '单用户自然日最多可申请次数；0 表示不限制。按服务器本地日 0 点起算。',
  minUsd: '最低提现金额(USD)',
  minUsdTip: '按钱包「1 USD = N 金币」折算，申请金币数对应 USD 须不低于此值。',
  lastUpdated: '最近更新',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
}

const en = {
  enabled: 'User withdrawal',
  enabledTip: 'When off, the app cannot apply for withdrawal or edit payout info. Pending CMS records can still be processed.',
  dailyMaxCount: 'Daily withdrawal limit',
  dailyMaxCountTip: 'Max applications per user per calendar day (server local time). 0 = unlimited.',
  minUsd: 'Minimum withdrawal (USD)',
  minUsdTip: 'Uses wallet USD-to-gold rate; requested gold must convert to at least this USD amount.',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved successfully',
  saveFailed: 'Save failed',
}

const es = {
  enabled: 'Retiro de usuario',
  enabledTip: 'Desactivado: la app no puede solicitar retiro ni editar datos de cobro. Los registros pendientes en CMS siguen procesables.',
  dailyMaxCount: 'Límite diario de retiros',
  dailyMaxCountTip: 'Máx. solicitudes por usuario y día (hora local del servidor). 0 = sin límite.',
  minUsd: 'Retiro mínimo (USD)',
  minUsdTip: 'Usa la tasa USD-oro de la billetera; el oro solicitado debe equivaler al menos a este USD.',
  lastUpdated: 'Última actualización',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccess: 'Guardado correctamente',
  saveFailed: 'Error al guardar',
}

const pt = {
  enabled: 'Saque do usuário',
  enabledTip: 'Desligado: o app não pode solicitar saque nem editar dados de recebimento. Registros pendentes no CMS ainda podem ser processados.',
  dailyMaxCount: 'Limite diário de saques',
  dailyMaxCountTip: 'Máx. solicitações por usuário por dia (fuso local do servidor). 0 = ilimitado.',
  minUsd: 'Saque mínimo (USD)',
  minUsdTip: 'Usa a taxa USD-ouro da carteira; o ouro solicitado deve equivaler a pelo menos este USD.',
  lastUpdated: 'Última atualização',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccess: 'Salvo com sucesso',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  enabled: 'उपयोगकर्ता निकासी',
  enabledTip: 'बंद: ऐप निकासी या भुगतान जानकारी संपादित नहीं कर सकता। CMS में लंबित रिकॉर्ड अभी भी संसाधित हो सकते हैं।',
  dailyMaxCount: 'दैनिक निकासी सीमा',
  dailyMaxCountTip: 'प्रति उपयोगकर्ता प्रति दिन अधिकतम आवेदन (सर्वर स्थानीय समय)। 0 = असीमित।',
  minUsd: 'न्यूनतम निकासी (USD)',
  minUsdTip: 'वॉलेट USD-सोने की दर; अनुरोधित सोना कम से कम इस USD के बराबर होना चाहिए।',
  lastUpdated: 'अंतिम अपडेट',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccess: 'सफलतापूर्वक सहेजा गया',
  saveFailed: 'सहेजना विफल',
}

const id = {
  enabled: 'Penarikan pengguna',
  enabledTip: 'Mati: app tidak bisa ajukan penarikan atau edit info penerima. Record pending di CMS tetap bisa diproses.',
  dailyMaxCount: 'Batas penarikan harian',
  dailyMaxCountTip: 'Maks. pengajuan per pengguna per hari (waktu lokal server). 0 = tanpa batas.',
  minUsd: 'Penarikan minimum (USD)',
  minUsdTip: 'Pakai kurs USD-emas dompet; emas yang diminta harus setara minimal USD ini.',
  lastUpdated: 'Terakhir diperbarui',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const userWithdrawCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

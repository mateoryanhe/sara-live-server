import {definePageMessages} from './_define'

const zh = {
  enabled: '用户提现开关',
  enabledTip: '关闭后 App 无法申请提现、无法维护收款信息；已提交的提现单仍可在 CMS 处理。',
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
  minUsd: 'Penarikan minimum (USD)',
  minUsdTip: 'Pakai kurs USD-emas dompet; emas yang diminta harus setara minimal USD ini.',
  lastUpdated: 'Terakhir diperbarui',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const userWithdrawCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

import {definePageMessages} from './_define'

const zh = {
  enabled: '设备码注册风控',
  enabledTip: '开启后：同设备账号数超限禁止新注册；同设备每日注销次数超限禁止再注销。关闭后不做限制。',
  maxCount: '最大注册数',
  maxCountTip: '同一设备码(含已注销账号)允许保留的账号上限，超出后无法再注册新账号',
  dailyCancelLimit: '每天最多注销次',
  dailyCancelLimitTip: '同一设备码每天允许成功注销的次数上限',
  lastUpdated: '最近更新',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
}

const en = {
  enabled: 'Device register risk control',
  enabledTip: 'When on: block new registration if device account count exceeds max; block cancel if daily cancel limit is reached. Off = no limits.',
  maxCount: 'Max accounts per device',
  maxCountTip: 'Max accounts for the same device code (including cancelled). Over limit cannot register again.',
  dailyCancelLimit: 'Max cancels per day',
  dailyCancelLimitTip: 'Max successful account cancels per device code per day',
  lastUpdated: 'Last updated',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved successfully',
  saveFailed: 'Save failed',
}

const es = {
  enabled: 'Control de registro por dispositivo',
  enabledTip: 'Activo: bloquea registro si se supera el máximo de cuentas; bloquea cancelación si se supera el límite diario. Inactivo = sin límites.',
  maxCount: 'Máx. cuentas por dispositivo',
  maxCountTip: 'Máximo de cuentas del mismo código de dispositivo (incluye canceladas)',
  dailyCancelLimit: 'Máx. cancelaciones/día',
  dailyCancelLimitTip: 'Máximo de cancelaciones exitosas por dispositivo al día',
  lastUpdated: 'Última actualización',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccess: 'Guardado correctamente',
  saveFailed: 'Error al guardar',
}

const pt = {
  enabled: 'Controle de registro por dispositivo',
  enabledTip: 'Ativo: bloqueia registro se exceder o máximo de contas; bloqueia cancelamento se exceder o limite diário. Inativo = sem limites.',
  maxCount: 'Máx. contas por dispositivo',
  maxCountTip: 'Máximo de contas do mesmo código de dispositivo (inclui canceladas)',
  dailyCancelLimit: 'Máx. cancelamentos/dia',
  dailyCancelLimitTip: 'Máximo de cancelamentos bem-sucedidos por dispositivo por dia',
  lastUpdated: 'Última atualização',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccess: 'Salvo com sucesso',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  enabled: 'डिवाइस पंजीकरण जोखिम नियंत्रण',
  enabledTip: 'चालू: अधिकतम खातों से अधिक पर नया रजिस्टर रोकें; दैनिक रद्द सीमा पर रद्द रोकें। बंद = कोई सीमा नहीं।',
  maxCount: 'प्रति डिवाइस अधिकतम खाते',
  maxCountTip: 'एक ही डिवाइस कोड के खातों की अधिकतम संख्या (रद्द शामिल)',
  dailyCancelLimit: 'प्रतिदिन अधिकतम रद्द',
  dailyCancelLimitTip: 'प्रति डिवाइस प्रति दिन सफल रद्द की अधिकतम संख्या',
  lastUpdated: 'अंतिम अपडेट',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccess: 'सफलतापूर्वक सहेजा',
  saveFailed: 'सहेजना विफल',
}

const id = {
  enabled: 'Kontrol risiko registrasi perangkat',
  enabledTip: 'Aktif: blokir daftar baru jika melebihi maks akun; blokir batal jika melebihi batas harian. Nonaktif = tanpa batas.',
  maxCount: 'Maks. akun per perangkat',
  maxCountTip: 'Batas akun untuk kode perangkat yang sama (termasuk yang dibatalkan)',
  dailyCancelLimit: 'Maks. batalkan/hari',
  dailyCancelLimitTip: 'Batas pembatalan sukses per kode perangkat per hari',
  lastUpdated: 'Terakhir diperbarui',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const deviceRegisterRiskCfgMessages = definePageMessages(zh, en, es, pt, hi, id)

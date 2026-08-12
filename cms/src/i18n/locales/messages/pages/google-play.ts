import {definePageMessages} from './_define'

const zh = {
  noticeTitle: '说明',
  noticeLine1: '配置 Google Play RTDN 回调与验单参数。Pub/Sub Push 地址：POST /webhook/googlePlay/rtdn',
  noticeLine2: 'App 创建订单后，发起 Google 支付时需将返回的 obfuscatedAccountId 传入 BillingFlowParams.setObfuscatedAccountId。',
  noticeLine3: '充值档位请在「运营 → 充值配置」中维护 Google 类型商品 SKU（productId）。发货成功后会自动 consume 消耗型商品。',
  serviceAccountJson: '服务账号 JSON',
  serviceAccountJsonPlaceholder: '粘贴 Google Cloud 服务账号 JSON 全文',
  serviceAccountJsonTip: '需具备 Google Play Android Developer API 权限，并在 Play Console 关联该服务账号',
  lastUpdated: '最近更新',
  serviceAccountJsonRequired: '请填写服务账号 JSON',
  serviceAccountJsonInvalid: '服务账号 JSON 格式无效',
  fetchCfgFailed: '获取配置失败',
  saveSuccess: '保存成功',
  saveFailed: '保存失败',
}

const en = {
  noticeTitle: 'Notice',
  noticeLine1: 'Configure Google Play RTDN callback and purchase verification. Pub/Sub push URL: POST /webhook/googlePlay/rtdn',
  noticeLine2: 'After creating an order, pass the returned obfuscatedAccountId to BillingFlowParams.setObfuscatedAccountId when starting Google payment.',
  noticeLine3: 'Maintain Google product SKUs (productId) under Operations → Recharge config. Consumable products are auto-consumed after successful delivery.',
  serviceAccountJson: 'Service account JSON',
  serviceAccountJsonPlaceholder: 'Paste full Google Cloud service account JSON',
  serviceAccountJsonTip: 'Requires Google Play Android Developer API access and linking in Play Console',
  lastUpdated: 'Last updated',
  serviceAccountJsonRequired: 'Please enter service account JSON',
  serviceAccountJsonInvalid: 'Invalid service account JSON format',
  fetchCfgFailed: 'Failed to load config',
  saveSuccess: 'Saved successfully',
  saveFailed: 'Save failed',
}

const es = {
  noticeTitle: 'Aviso',
  noticeLine1: 'Configura RTDN de Google Play y verificación. URL Pub/Sub: POST /webhook/googlePlay/rtdn',
  noticeLine2: 'Tras crear pedido, pasa obfuscatedAccountId a BillingFlowParams.setObfuscatedAccountId al pagar.',
  noticeLine3: 'Mantén SKU Google (productId) en Operaciones → Config de recarga. Consumibles se consumen tras entrega.',
  serviceAccountJson: 'JSON cuenta de servicio',
  serviceAccountJsonPlaceholder: 'Pega el JSON completo de la cuenta de servicio',
  serviceAccountJsonTip: 'Requiere API Google Play Android Developer y vinculación en Play Console',
  lastUpdated: 'Última actualización',
  serviceAccountJsonRequired: 'Introduce JSON de cuenta de servicio',
  serviceAccountJsonInvalid: 'Formato JSON de cuenta de servicio inválido',
  fetchCfgFailed: 'Error al cargar la configuración',
  saveSuccess: 'Guardado correctamente',
  saveFailed: 'Error al guardar',
}

const pt = {
  noticeTitle: 'Aviso',
  noticeLine1: 'Configure RTDN Google Play e verificação. URL Pub/Sub: POST /webhook/googlePlay/rtdn',
  noticeLine2: 'Após criar pedido, passe obfuscatedAccountId para BillingFlowParams.setObfuscatedAccountId no pagamento.',
  noticeLine3: 'Mantenha SKUs Google (productId) em Operações → Config de recarga. Consumíveis são consumidos após entrega.',
  serviceAccountJson: 'JSON conta de serviço',
  serviceAccountJsonPlaceholder: 'Cole o JSON completo da conta de serviço Google Cloud',
  serviceAccountJsonTip: 'Requer API Google Play Android Developer e vínculo no Play Console',
  lastUpdated: 'Última atualização',
  serviceAccountJsonRequired: 'Digite JSON da conta de serviço',
  serviceAccountJsonInvalid: 'Formato JSON da conta de serviço inválido',
  fetchCfgFailed: 'Falha ao carregar configuração',
  saveSuccess: 'Salvo com sucesso',
  saveFailed: 'Falha ao salvar',
}

const hi = {
  noticeTitle: 'सूचना',
  noticeLine1: 'Google Play RTDN कॉलबैक और सत्यापन कॉन्फ़िग करें। Pub/Sub push: POST /webhook/googlePlay/rtdn',
  noticeLine2: 'ऑर्डर बनाने के बाद Google भुगतान में obfuscatedAccountId BillingFlowParams.setObfuscatedAccountId में पास करें।',
  noticeLine3: 'ऑपरेशन → रिचार्ज कॉन्फ़िग में Google SKU (productId) रखें। डिलीवरी के बाद consumable auto-consume।',
  serviceAccountJson: 'सेवा खाता JSON',
  serviceAccountJsonPlaceholder: 'Google Cloud सेवा खाता JSON पूरा चिपकाएँ',
  serviceAccountJsonTip: 'Google Play Android Developer API और Play Console लिंक आवश्यक',
  lastUpdated: 'अंतिम अपडेट',
  serviceAccountJsonRequired: 'सेवा खाता JSON दर्ज करें',
  serviceAccountJsonInvalid: 'सेवा खाता JSON अमान्य',
  fetchCfgFailed: 'कॉन्फ़िग लोड विफल',
  saveSuccess: 'सफलतापूर्वक सहेजा',
  saveFailed: 'सहेजना विफल',
}

const id = {
  noticeTitle: 'Catatan',
  noticeLine1: 'Konfigurasi callback RTDN Google Play dan verifikasi. URL Pub/Sub: POST /webhook/googlePlay/rtdn',
  noticeLine2: 'Setelah buat pesanan, kirim obfuscatedAccountId ke BillingFlowParams.setObfuscatedAccountId saat bayar.',
  noticeLine3: 'Kelola SKU Google (productId) di Operasi → Konfigurasi isi ulang. Produk consumable auto-consume setelah kirim.',
  serviceAccountJson: 'JSON akun layanan',
  serviceAccountJsonPlaceholder: 'Tempel JSON lengkap akun layanan Google Cloud',
  serviceAccountJsonTip: 'Butuh API Google Play Android Developer dan tautan di Play Console',
  lastUpdated: 'Terakhir diperbarui',
  serviceAccountJsonRequired: 'Masukkan JSON akun layanan',
  serviceAccountJsonInvalid: 'Format JSON akun layanan tidak valid',
  fetchCfgFailed: 'Gagal memuat konfigurasi',
  saveSuccess: 'Berhasil disimpan',
  saveFailed: 'Gagal menyimpan',
}

export const googlePlayMessages = definePageMessages(zh, en, es, pt, hi, id)

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.HaiPayCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.haipay.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.haipay.noticeLine1') }}</p>
        <p>{{ t('pages.haipay.noticeLine2') }}</p>
        <p>{{ t('pages.haipay.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.haipay.enabled')" prop="enabled">
          <el-switch v-model="formData.enabled"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.appId')" prop="appId">
          <el-input v-model="formData.appId" clearable :placeholder="t('pages.haipay.appIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.apiHost')" prop="apiHost">
          <el-input v-model="formData.apiHost" clearable :placeholder="t('pages.haipay.apiHostPlaceholder')"/>
          <span class="form-tip">{{ t('pages.haipay.apiHostTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.merchantSecretKey')" prop="merchantSecretKey">
          <el-input v-model="formData.merchantSecretKey" clearable show-password type="password"
                    :placeholder="t('pages.haipay.merchantSecretKeyPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.merchantPrivateKey')" prop="merchantPrivateKey">
          <el-input v-model="formData.merchantPrivateKey" :rows="5" type="textarea"
                    :placeholder="t('pages.haipay.merchantPrivateKeyPlaceholder')"/>
          <span class="form-tip">{{ t('pages.haipay.keyTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.haiPayPublicKey')" prop="haiPayPublicKey">
          <el-input v-model="formData.haiPayPublicKey" :rows="5" type="textarea"
                    :placeholder="t('pages.haipay.haiPayPublicKeyPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.callbackBaseUrl')" prop="callbackBaseUrl">
          <el-input v-model="formData.callbackBaseUrl" clearable
                    :placeholder="t('pages.haipay.callbackBaseUrlPlaceholder')"/>
          <span class="form-tip">{{ t('pages.haipay.callbackBaseUrlTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.returnUrl')" prop="returnUrl">
          <el-input v-model="formData.returnUrl" clearable :placeholder="t('pages.haipay.returnUrlPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.failReturnUrl')" prop="failReturnUrl">
          <el-input v-model="formData.failReturnUrl" clearable
                    :placeholder="t('pages.haipay.failReturnUrlPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.cancelUrl')" prop="cancelUrl">
          <el-input v-model="formData.cancelUrl" clearable :placeholder="t('pages.haipay.cancelUrlPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.paymentMethods')" prop="paymentMethods">
          <el-input v-model="formData.paymentMethods" clearable
                    :placeholder="t('pages.haipay.paymentMethodsPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.subject')" prop="subject">
          <el-input v-model="formData.subject" clearable :placeholder="t('pages.haipay.subjectPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.haipay.defaultRegion')" prop="defaultRegion">
          <el-input v-model="formData.defaultRegion" clearable :placeholder="t('pages.haipay.defaultRegionPlaceholder')"/>
          <span class="form-tip">{{ t('pages.haipay.defaultRegionTip') }}</span>
        </el-form-item>
        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.haipay.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {haipayApi} from '@/api/modules/haipay'
import type {HaiPayCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  enabled: false,
  appId: '',
  apiHost: '',
  merchantSecretKey: '',
  merchantPrivateKey: '',
  haiPayPublicKey: '',
  callbackBaseUrl: '',
  returnUrl: '',
  failReturnUrl: '',
  cancelUrl: '',
  paymentMethods: '',
  subject: 'Recharge',
  defaultRegion: 'ID',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  appId: [{required: true, message: t('pages.haipay.appIdRequired'), trigger: 'blur'}],
  apiHost: [{required: true, message: t('pages.haipay.apiHostRequired'), trigger: 'blur'}],
  merchantSecretKey: [{required: true, message: t('pages.haipay.merchantSecretKeyRequired'), trigger: 'blur'}],
  merchantPrivateKey: [{required: true, message: t('pages.haipay.merchantPrivateKeyRequired'), trigger: 'blur'}],
  haiPayPublicKey: [{required: true, message: t('pages.haipay.haiPayPublicKeyRequired'), trigger: 'blur'}],
}))

const applyCfg = (cfg: HaiPayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.appId = ''
    formData.apiHost = ''
    formData.merchantSecretKey = ''
    formData.merchantPrivateKey = ''
    formData.haiPayPublicKey = ''
    formData.callbackBaseUrl = ''
    formData.returnUrl = ''
    formData.failReturnUrl = ''
    formData.cancelUrl = ''
    formData.paymentMethods = ''
    formData.subject = 'Recharge'
    formData.defaultRegion = 'ID'
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.appId = cfg.appId ? String(cfg.appId) : ''
  formData.apiHost = cfg.apiHost || ''
  formData.merchantSecretKey = cfg.merchantSecretKey || ''
  formData.merchantPrivateKey = cfg.merchantPrivateKey || ''
  formData.haiPayPublicKey = cfg.haiPayPublicKey || ''
  formData.callbackBaseUrl = cfg.callbackBaseUrl || ''
  formData.returnUrl = cfg.returnUrl || ''
  formData.failReturnUrl = cfg.failReturnUrl || ''
  formData.cancelUrl = cfg.cancelUrl || ''
  formData.paymentMethods = cfg.paymentMethods || ''
  formData.subject = cfg.subject || 'Recharge'
  formData.defaultRegion = cfg.defaultRegion || 'ID'
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await haipayApi.getHaiPayCfg()
    applyCfg(response?.cfg)
  } catch (error) {
    console.error('fetch haipay cfg failed:', error)
    ElMessage.error(t('pages.haipay.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    const appIdNum = Number(formData.appId)
    const response = await haipayApi.saveHaiPayCfg({
      id: formData.id && formData.id !== '0' ? Number(formData.id) : undefined,
      enabled: formData.enabled,
      appId: appIdNum,
      apiHost: formData.apiHost.trim(),
      merchantSecretKey: formData.merchantSecretKey.trim(),
      merchantPrivateKey: formData.merchantPrivateKey.trim(),
      haiPayPublicKey: formData.haiPayPublicKey.trim(),
      callbackBaseUrl: formData.callbackBaseUrl.trim(),
      returnUrl: formData.returnUrl.trim(),
      failReturnUrl: formData.failReturnUrl.trim(),
      cancelUrl: formData.cancelUrl.trim(),
      paymentMethods: formData.paymentMethods.trim(),
      subject: formData.subject.trim() || 'Recharge',
      defaultRegion: (formData.defaultRegion || 'ID').trim().toUpperCase(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.haipay.saveSuccess'))
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.haipay.saveFailed'))
    }
  } catch (error) {
    console.error('save haipay cfg failed:', error)
    ElMessage.error(t('pages.haipay.saveFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(fetchCfg)
</script>

<style scoped>
.page-container {
  padding: 16px;
}

.card-header {
  font-weight: 600;
}

.tip-alert {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 860px;
}

.form-tip {
  display: block;
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}
</style>

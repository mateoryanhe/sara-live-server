<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.YhPayCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.yhpay.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.yhpay.noticeLine1') }}</p>
        <p>{{ t('pages.yhpay.noticeLine2') }}</p>
        <p>{{ t('pages.yhpay.noticeLine3') }}</p>
        <p>{{ t('pages.yhpay.noticeLine4') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.yhpay.merchantCode')" prop="merchantCode">
          <el-input v-model="formData.merchantCode" clearable :placeholder="t('pages.yhpay.merchantCodePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.apiKey')" prop="apiKey">
          <el-input v-model="formData.apiKey" clearable show-password type="password" :placeholder="t('pages.yhpay.apiKeyPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.apiHost')" prop="apiHost">
          <el-input v-model="formData.apiHost" clearable :placeholder="t('pages.yhpay.apiHostPlaceholder')"/>
          <span class="form-tip">{{ t('pages.yhpay.apiHostTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.cryptoApiHost')" prop="cryptoApiHost">
          <el-input v-model="formData.cryptoApiHost" clearable :placeholder="t('pages.yhpay.cryptoApiHostPlaceholder')"/>
          <span class="form-tip">{{ t('pages.yhpay.cryptoApiHostTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.callbackBaseUrl')" prop="callbackBaseUrl">
          <el-input v-model="formData.callbackBaseUrl" clearable :placeholder="t('pages.yhpay.callbackBaseUrlPlaceholder')"/>
          <span class="form-tip">{{ t('pages.yhpay.callbackBaseUrlTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.returnUrl')" prop="returnUrl">
          <el-input v-model="formData.returnUrl" clearable :placeholder="t('pages.yhpay.returnUrlPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.failedReturnUrl')" prop="failedReturnUrl">
          <el-input v-model="formData.failedReturnUrl" clearable :placeholder="t('pages.yhpay.failedReturnUrlPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.yhpay.cryptoNetwork')">
          <span>{{ t('pages.yhpay.cryptoNetworkFixed') }}</span>
          <span class="form-tip">{{ t('pages.yhpay.cryptoNetworkTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.yhpay.lastUpdated')">
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
import {yhpayApi} from '@/api/modules/yhpay'
import type {YhPayCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  merchantCode: '',
  apiKey: '',
  apiHost: '',
  cryptoApiHost: '',
  callbackBaseUrl: '',
  returnUrl: '',
  failedReturnUrl: '',
})

const FIXED_CRYPTO_NETWORK = 'TRC20'

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  merchantCode: [{required: true, message: t('pages.yhpay.merchantCodeRequired'), trigger: 'blur'}],
  apiKey: [{required: true, message: t('pages.yhpay.apiKeyRequired'), trigger: 'blur'}],
  apiHost: [{required: true, message: t('pages.yhpay.apiHostRequired'), trigger: 'blur'}],
  cryptoApiHost: [{required: true, message: t('pages.yhpay.cryptoApiHostRequired'), trigger: 'blur'}],
}))

const applyCfg = (cfg: YhPayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.merchantCode = ''
    formData.apiKey = ''
    formData.apiHost = ''
    formData.cryptoApiHost = ''
    formData.callbackBaseUrl = ''
    formData.returnUrl = ''
    formData.failedReturnUrl = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.merchantCode = cfg.merchantCode || ''
  formData.apiKey = cfg.apiKey || ''
  formData.apiHost = cfg.apiHost || ''
  formData.cryptoApiHost = cfg.cryptoApiHost || ''
  formData.callbackBaseUrl = cfg.callbackBaseUrl || ''
  formData.returnUrl = cfg.returnUrl || ''
  formData.failedReturnUrl = cfg.failedReturnUrl || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await yhpayApi.getYhPayCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch yhpay cfg failed:', error)
    ElMessage.error(t('pages.yhpay.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await yhpayApi.saveYhPayCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: true,
      merchantCode: formData.merchantCode.trim(),
      apiKey: formData.apiKey.trim(),
      apiHost: formData.apiHost.trim(),
      cryptoApiHost: formData.cryptoApiHost.trim(),
      callbackBaseUrl: formData.callbackBaseUrl.trim(),
      returnUrl: formData.returnUrl.trim(),
      failedReturnUrl: formData.failedReturnUrl.trim(),
      cryptoNetwork: FIXED_CRYPTO_NETWORK,
    })
    if (response?.success) {
      ElMessage.success(t('pages.yhpay.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.yhpay.saveFailed'))
    }
  } catch (error) {
    console.error('save yhpay cfg failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCfg()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.cfg-form {
  max-width: 760px;
}

.form-tip {
  display: block;
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}
</style>

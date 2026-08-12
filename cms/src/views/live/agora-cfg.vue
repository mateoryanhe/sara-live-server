<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AgoraCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="AppId" prop="appId">
          <el-input v-model="formData.appId" clearable :placeholder="t('pages.agoraCfg.appIdPlaceholder')"/>
        </el-form-item>

        <el-form-item label="AppCertificate" prop="appCertificate">
          <el-input
              v-model="formData.appCertificate"
              clearable
              :placeholder="t('pages.agoraCfg.appCertificatePlaceholder')"
              show-password
              type="password"
          />
        </el-form-item>

        <el-form-item label="REST CustomerId" prop="restCustomerId">
          <el-input v-model="formData.restCustomerId" clearable :placeholder="t('pages.agoraCfg.restCustomerIdPlaceholder')"/>
          <span class="form-tip">{{ t('pages.agoraCfg.restCustomerTip') }}</span>
        </el-form-item>

        <el-form-item label="REST CustomerSecret" prop="restCustomerSecret">
          <el-input
              v-model="formData.restCustomerSecret"
              clearable
              :placeholder="t('pages.agoraCfg.restCustomerSecretPlaceholder')"
              show-password
              type="password"
          />
          <span class="form-tip">{{ t('pages.agoraCfg.restCustomerTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.agoraCfg.cloudPlayerRegion')" prop="cloudPlayerRegion">
          <el-select v-model="formData.cloudPlayerRegion" :placeholder="t('pages.agoraCfg.selectRegion')" style="width: 220px">
            <el-option :label="t('pages.agoraCfg.regionCn')" value="cn"/>
            <el-option :label="t('pages.agoraCfg.regionAp')" value="ap"/>
            <el-option :label="t('pages.agoraCfg.regionEu')" value="eu"/>
            <el-option :label="t('pages.agoraCfg.regionNa')" value="na"/>
          </el-select>
        </el-form-item>

        <el-form-item :label="t('pages.agoraCfg.tokenExpireHours')" prop="tokenExpireHours">
          <el-input-number
              v-model="formData.tokenExpireHours"
              :max="TOKEN_EXPIRE_MAX_HOURS"
              :min="TOKEN_EXPIRE_MIN_HOURS"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ tokenExpireTip }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.agoraCfg.lastUpdated')">
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
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {agoraApi} from '@/api/modules/agora'
import type {AgoraCfg} from '@/types/api'

const {t} = useI18n()
const TOKEN_EXPIRE_MIN_HOURS = 4
const TOKEN_EXPIRE_MAX_HOURS = 24
const TOKEN_EXPIRE_DEFAULT_HOURS = 24
const TOKEN_REFRESH_AHEAD_GAP_HOURS = 2
const TOKEN_REFRESH_MIN_HOURS = 2
const SECONDS_PER_HOUR = 3600

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  appId: '',
  appCertificate: '',
  restCustomerId: '',
  restCustomerSecret: '',
  cloudPlayerRegion: 'cn',
  tokenExpireHours: TOKEN_EXPIRE_DEFAULT_HOURS,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const tokenExpireTip = computed(() =>
    t('pages.agoraCfg.tokenExpireTip', {
      min: TOKEN_EXPIRE_MIN_HOURS,
      max: TOKEN_EXPIRE_MAX_HOURS,
      default: TOKEN_EXPIRE_DEFAULT_HOURS,
    })
)

const formRules = computed(() => ({
  appId: [
    {required: true, message: t('pages.agoraCfg.appIdRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.agoraCfg.appIdLength'), trigger: 'blur'},
  ],
  appCertificate: [
    {required: true, message: t('pages.agoraCfg.appCertificateRequired'), trigger: 'blur'},
    {min: 1, max: 128, message: t('pages.agoraCfg.appCertificateLength'), trigger: 'blur'},
  ],
  restCustomerId: [
    {max: 64, message: t('pages.agoraCfg.restCustomerIdMax'), trigger: 'blur'},
  ],
  restCustomerSecret: [
    {max: 128, message: t('pages.agoraCfg.restCustomerSecretMax'), trigger: 'blur'},
  ],
  tokenExpireHours: [
    {required: true, message: t('pages.agoraCfg.tokenExpireRequired'), trigger: 'blur'},
    {
      type: 'number',
      min: TOKEN_EXPIRE_MIN_HOURS,
      max: TOKEN_EXPIRE_MAX_HOURS,
      message: t('pages.agoraCfg.tokenExpireRange', {min: TOKEN_EXPIRE_MIN_HOURS, max: TOKEN_EXPIRE_MAX_HOURS}),
      trigger: 'blur',
    },
  ],
}))

const clampTokenExpireHours = (hours: number) => {
  if (!Number.isFinite(hours) || hours <= 0) {
    return TOKEN_EXPIRE_DEFAULT_HOURS
  }
  return Math.min(TOKEN_EXPIRE_MAX_HOURS, Math.max(TOKEN_EXPIRE_MIN_HOURS, Math.round(hours)))
}

const secondsToHours = (seconds: number, fallback: number) => {
  if (!Number.isFinite(seconds) || seconds < 0) {
    return fallback
  }
  return Math.round(seconds / SECONDS_PER_HOUR)
}

const hoursToSeconds = (hours: number) => Math.round(hours) * SECONDS_PER_HOUR

const buildLegacyTokenRefreshSeconds = (expireHours: number) => {
  const refreshHours = Math.max(TOKEN_REFRESH_MIN_HOURS, expireHours - TOKEN_REFRESH_AHEAD_GAP_HOURS)
  return hoursToSeconds(refreshHours)
}

const applyCfg = (cfg: AgoraCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.appId = ''
    formData.appCertificate = ''
    formData.restCustomerId = ''
    formData.restCustomerSecret = ''
    formData.cloudPlayerRegion = 'cn'
    formData.tokenExpireHours = TOKEN_EXPIRE_DEFAULT_HOURS
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.appId = cfg.appId || ''
  formData.appCertificate = cfg.appCertificate || ''
  formData.restCustomerId = cfg.restCustomerId || ''
  formData.restCustomerSecret = cfg.restCustomerSecret || ''
  formData.cloudPlayerRegion = cfg.cloudPlayerRegion || 'cn'
  formData.tokenExpireHours = clampTokenExpireHours(secondsToHours(cfg.tokenExpireSeconds, TOKEN_EXPIRE_DEFAULT_HOURS))
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await agoraApi.getAgoraCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch agora config failed:', error)
    ElMessage.error(t('pages.agoraCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await agoraApi.saveAgoraCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      appId: formData.appId.trim(),
      appCertificate: formData.appCertificate.trim(),
      restCustomerId: formData.restCustomerId.trim(),
      restCustomerSecret: formData.restCustomerSecret.trim(),
      cloudPlayerRegion: formData.cloudPlayerRegion || 'cn',
      tokenExpireSeconds: hoursToSeconds(formData.tokenExpireHours),
      tokenRefreshSeconds: buildLegacyTokenRefreshSeconds(formData.tokenExpireHours),
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.agoraCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save agora config failed:', error)
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

.cfg-form {
  max-width: 720px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>

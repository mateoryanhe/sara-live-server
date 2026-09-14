<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.DeviceRegisterRiskCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="200px">
        <el-form-item :label="t('pages.deviceRegisterRiskCfg.enabled')">
          <el-switch
              v-model="formData.deviceRegisterRiskEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">
            {{ t('pages.deviceRegisterRiskCfg.enabledTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.deviceRegisterRiskCfg.maxCount')">
          <el-input-number
              v-model="formData.deviceAccountMaxCount"
              :min="1"
              :max="100"
              :disabled="!formData.deviceRegisterRiskEnabled"
          />
          <div class="form-tip">
            {{ t('pages.deviceRegisterRiskCfg.maxCountTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.deviceRegisterRiskCfg.dailyCancelLimit')">
          <el-input-number
              v-model="formData.deviceCancelDailyLimit"
              :min="1"
              :max="100"
              :disabled="!formData.deviceRegisterRiskEnabled"
          />
          <div class="form-tip">
            {{ t('pages.deviceRegisterRiskCfg.dailyCancelLimitTip') }}
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.deviceRegisterRiskCfg.lastUpdated')">
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import accountCfgApi from '@/api/modules/account-cfg'
import type {AccountCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  cancelAccountByCodeEnabled: false,
  blockSimulatorLogin: false,
  envType: 0,
  deviceRegisterRiskEnabled: true,
  deviceAccountMaxCount: 3,
  deviceCancelDailyLimit: 1,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const normalizeEnvType = (v: unknown) => {
  const n = Number(v)
  return n === 1 || n === 2 ? n : 0
}

const normalizePositiveInt = (v: unknown, fallback: number) => {
  const n = Number(v)
  return Number.isFinite(n) && n >= 1 ? Math.floor(n) : fallback
}

const applyCfg = (cfg: AccountCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.cancelAccountByCodeEnabled = !!cfg?.cancelAccountByCodeEnabled
  formData.blockSimulatorLogin = !!cfg?.blockSimulatorLogin
  formData.envType = normalizeEnvType(cfg?.envType)
  formData.deviceRegisterRiskEnabled = cfg?.deviceRegisterRiskEnabled !== false
  formData.deviceAccountMaxCount = normalizePositiveInt(cfg?.deviceAccountMaxCount, 3)
  formData.deviceCancelDailyLimit = normalizePositiveInt(cfg?.deviceCancelDailyLimit, 1)
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await accountCfgApi.getAccountCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch device register risk cfg failed:', error)
    ElMessage.error(t('pages.deviceRegisterRiskCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    const response = await accountCfgApi.saveAccountCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      cancelAccountByCodeEnabled: formData.cancelAccountByCodeEnabled,
      blockSimulatorLogin: formData.blockSimulatorLogin,
      envType: normalizeEnvType(formData.envType),
      deviceRegisterRiskEnabled: formData.deviceRegisterRiskEnabled,
      deviceAccountMaxCount: normalizePositiveInt(formData.deviceAccountMaxCount, 3),
      deviceCancelDailyLimit: normalizePositiveInt(formData.deviceCancelDailyLimit, 1),
    })
    if (response?.success) {
      ElMessage.success(t('pages.deviceRegisterRiskCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.deviceRegisterRiskCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save device register risk cfg failed:', error)
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
}

.cfg-form {
  max-width: 720px;
}

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}
</style>

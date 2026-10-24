<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.UserWithdrawCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="200px">
        <el-form-item :label="t('pages.userWithdrawCfg.enabled')">
          <el-switch
              v-model="formData.userWithdrawEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">{{ t('pages.userWithdrawCfg.enabledTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.userWithdrawCfg.androidVisible')">
          <el-switch
              v-model="formData.userWithdrawAndroidVisible"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">{{ t('pages.userWithdrawCfg.androidVisibleTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.userWithdrawCfg.dailyMaxCount')">
          <el-input-number
              v-model="formData.userWithdrawDailyMaxCount"
              :min="0"
              :max="999"
              :precision="0"
              :step="1"
              :disabled="!formData.userWithdrawEnabled"
          />
          <div class="form-tip">{{ t('pages.userWithdrawCfg.dailyMaxCountTip') }}</div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.userWithdrawCfg.lastUpdated')">
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

const formData = reactive({
  id: '0',
  cancelAccountByCodeEnabled: false,
  blockSimulatorLogin: false,
  envType: 0,
  deviceRegisterRiskEnabled: true,
  deviceAccountMaxCount: 3,
  deviceCancelDailyLimit: 1,
  userWithdrawEnabled: false,
  userWithdrawAndroidVisible: false,
  userWithdrawDailyMaxCount: 1,
})

const metaInfo = reactive({createdAt: '', updatedAt: ''})
/** 账号配置整表保存，保留服务端原最低提现金额，避免误覆盖 */
const preservedUserWithdrawMinUsd = ref(5)

const normalizePositiveInt = (v: unknown, fallback: number) => {
  const n = Number(v)
  return Number.isFinite(n) && n >= 1 ? Math.floor(n) : fallback
}

const applyCfg = (cfg: AccountCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.cancelAccountByCodeEnabled = !!cfg?.cancelAccountByCodeEnabled
  formData.blockSimulatorLogin = !!cfg?.blockSimulatorLogin
  formData.envType = Number(cfg?.envType) === 1 || Number(cfg?.envType) === 2 ? Number(cfg?.envType) : 0
  formData.deviceRegisterRiskEnabled = cfg?.deviceRegisterRiskEnabled !== false
  formData.deviceAccountMaxCount = normalizePositiveInt(cfg?.deviceAccountMaxCount, 3)
  formData.deviceCancelDailyLimit = normalizePositiveInt(cfg?.deviceCancelDailyLimit, 1)
  formData.userWithdrawEnabled = !!cfg?.userWithdrawEnabled
  formData.userWithdrawAndroidVisible = !!cfg?.userWithdrawAndroidVisible
  const dailyMax = Number(cfg?.userWithdrawDailyMaxCount)
  formData.userWithdrawDailyMaxCount = Number.isFinite(dailyMax) && dailyMax >= 0 ? Math.floor(dailyMax) : 1
  const minUsd = Number(cfg?.userWithdrawMinUsd)
  preservedUserWithdrawMinUsd.value = Number.isFinite(minUsd) && minUsd > 0 ? minUsd : 5
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await accountCfgApi.getAccountCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch user withdraw cfg failed:', error)
    ElMessage.error(t('pages.userWithdrawCfg.fetchCfgFailed'))
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
      envType: formData.envType,
      deviceRegisterRiskEnabled: formData.deviceRegisterRiskEnabled,
      deviceAccountMaxCount: formData.deviceAccountMaxCount,
      deviceCancelDailyLimit: formData.deviceCancelDailyLimit,
      userWithdrawEnabled: formData.userWithdrawEnabled,
      userWithdrawAndroidVisible: formData.userWithdrawAndroidVisible,
      userWithdrawMinUsd: preservedUserWithdrawMinUsd.value,
      userWithdrawDailyMaxCount: formData.userWithdrawDailyMaxCount,
    })
    if (response?.success) {
      ElMessage.success(t('pages.userWithdrawCfg.saveSuccess'))
      if (response.id) formData.id = response.id
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.userWithdrawCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save user withdraw cfg failed:', error)
  }
}

onMounted(fetchCfg)
</script>

<style scoped>
.page-container {
  padding: 16px;
}
.form-tip {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>

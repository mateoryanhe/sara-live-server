<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.InviteRechargeRewardManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.inviteRechargeRewardCfg.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.inviteRechargeRewardCfg.tipLine1') }}</p>
        <p>{{ t('pages.inviteRechargeRewardCfg.tipLine2') }}</p>
      </el-alert>

      <el-form :model="formData" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.inviteRechargeRewardCfg.enabled')">
          <el-switch
              v-model="formData.enabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
        </el-form-item>

        <el-form-item :label="t('pages.inviteRechargeRewardCfg.rewardPercent')">
          <el-input-number
              v-model="formData.rewardPercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
          />
          <span class="field-tip">%</span>
          <div class="form-tip">{{ t('pages.inviteRechargeRewardCfg.rewardPercentTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.inviteRechargeRewardCfg.validDays')">
          <el-input-number
              v-model="formData.validDays"
              :max="3650"
              :min="0"
              :precision="0"
              :step="1"
              controls-position="right"
          />
          <span class="field-tip">{{ t('pages.inviteRechargeRewardCfg.daysUnit') }}</span>
          <div class="form-tip">{{ t('pages.inviteRechargeRewardCfg.validDaysTip') }}</div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.inviteRechargeRewardCfg.lastUpdated')">
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
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {inviteRechargeRewardApi} from '@/api/modules/invite-recharge-reward'
import type {InviteRechargeRewardCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)

const formData = reactive({
  id: '0',
  enabled: false,
  rewardPercent: 5,
  validDays: 30,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: InviteRechargeRewardCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.rewardPercent = 5
    formData.validDays = 30
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.rewardPercent = cfg.rewardPercent ?? 5
  formData.validDays = cfg.validDays ?? 30
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await inviteRechargeRewardApi.getCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch invite recharge reward cfg failed:', error)
    ElMessage.error(t('pages.inviteRechargeRewardCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (formData.rewardPercent < 0 || formData.rewardPercent > 100) {
    ElMessage.warning(t('pages.inviteRechargeRewardCfg.percentRangeInvalid'))
    return
  }
  if (formData.validDays < 0) {
    ElMessage.warning(t('pages.inviteRechargeRewardCfg.daysInvalid'))
    return
  }

  loading.value = true
  try {
    const response = await inviteRechargeRewardApi.saveCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: formData.enabled,
      rewardPercent: formData.rewardPercent,
      validDays: formData.validDays,
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.inviteRechargeRewardCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save invite recharge reward cfg failed:', error)
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

.field-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.form-tip {
  width: 100%;
  margin-top: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}
</style>

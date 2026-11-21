<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AudienceOnlineAnchorGreetManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.audienceOnlineAnchorGreet.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.audienceOnlineAnchorGreet.tipLine1') }}</p>
        <p>{{ t('pages.audienceOnlineAnchorGreet.tipLine2') }}</p>
        <p>{{ t('pages.audienceOnlineAnchorGreet.tipLine3') }}</p>
      </el-alert>

      <el-form class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.audienceOnlineAnchorGreet.enabled')">
          <el-switch v-model="formData.enabled"/>
        </el-form-item>

        <el-form-item :label="t('pages.audienceOnlineAnchorGreet.anchorCount')" required>
          <el-input-number v-model="formData.anchorCount" :max="20" :min="1" style="max-width: 200px"/>
          <span class="field-hint">{{ t('pages.audienceOnlineAnchorGreet.anchorCountHint') }}</span>
        </el-form-item>

        <el-divider>{{ t('pages.audienceOnlineAnchorGreet.queueTitle') }}</el-divider>

        <el-form-item :label="t('pages.audienceOnlineAnchorGreet.intervalMinSec')" required>
          <el-input-number v-model="formData.intervalMinSec" :max="3600" :min="0" style="max-width: 200px"/>
          <span class="field-hint">{{ t('pages.audienceOnlineAnchorGreet.intervalHint') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.audienceOnlineAnchorGreet.intervalMaxSec')" required>
          <el-input-number v-model="formData.intervalMaxSec" :max="3600" :min="0" style="max-width: 200px"/>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.customerServiceCfg.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="canSave" type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
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
import {audienceOnlineAnchorGreetApi} from '@/api/modules/audience-online-anchor-greet'
import type {AudienceOnlineAnchorGreetCfg} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const loading = ref(false)
const {can} = usePagePermission('AudienceOnlineAnchorGreetManagement')
const canSave = computed(() => can('save'))

const formData = reactive({
  id: '0',
  enabled: false,
  anchorCount: 5,
  intervalMinSec: 30,
  intervalMaxSec: 30,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: AudienceOnlineAnchorGreetCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.anchorCount = 5
    formData.intervalMinSec = 30
    formData.intervalMaxSec = 30
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.anchorCount = cfg.anchorCount && cfg.anchorCount > 0 ? cfg.anchorCount : 5
  formData.intervalMinSec = cfg.intervalMinSec != null && cfg.intervalMinSec >= 0 ? cfg.intervalMinSec : 30
  formData.intervalMaxSec = cfg.intervalMaxSec != null && cfg.intervalMaxSec >= 0 ? cfg.intervalMaxSec : 30
  if (formData.intervalMaxSec < formData.intervalMinSec) {
    formData.intervalMaxSec = formData.intervalMinSec
  }
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await audienceOnlineAnchorGreetApi.getCfg()
    applyCfg(res.cfg)
  } catch {
    ElMessage.error(t('pages.audienceOnlineAnchorGreet.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (formData.intervalMaxSec < formData.intervalMinSec) {
    ElMessage.warning(t('pages.audienceOnlineAnchorGreet.intervalRangeInvalid'))
    return
  }
  loading.value = true
  try {
    const payload = {
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: formData.enabled,
      anchorCount: formData.anchorCount,
      intervalMinSec: formData.intervalMinSec,
      intervalMaxSec: formData.intervalMaxSec,
    }
    const res = await audienceOnlineAnchorGreetApi.saveCfg(payload)
    if (res.success) {
      ElMessage.success(t('pages.audienceOnlineAnchorGreet.saveSuccess'))
      await fetchCfg()
    }
  } catch {
    ElMessage.error(t('pages.audienceOnlineAnchorGreet.saveFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCfg()
})
</script>

<style scoped>
.tip-alert {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 780px;
}

.field-hint {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
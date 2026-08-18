<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRevenueShareCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.liveRevenueShareCfg.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.liveRevenueShareCfg.tipLine1') }}</p>
        <p>{{ t('pages.liveRevenueShareCfg.tipLine2') }}</p>
      </el-alert>

      <el-form :model="formData" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.liveRevenueShareCfg.anchorSharePercent')">
          <el-input-number
              v-model="formData.anchorSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
          />
          <span class="field-tip">%</span>
        </el-form-item>

        <el-form-item :label="t('pages.liveRevenueShareCfg.guildSharePercent')">
          <el-input-number
              v-model="formData.guildSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
          />
          <span class="field-tip">%</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.liveRevenueShareCfg.lastUpdated')">
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
import {liveRevenueShareCfgApi} from '@/api/modules/live-revenue-share-cfg'
import type {LiveRevenueShareCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)

const formData = reactive({
  id: '0',
  anchorSharePercent: 30,
  guildSharePercent: 10,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: LiveRevenueShareCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.anchorSharePercent = 30
    formData.guildSharePercent = 10
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.anchorSharePercent = cfg.anchorSharePercent ?? 30
  formData.guildSharePercent = cfg.guildSharePercent ?? 10
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await liveRevenueShareCfgApi.getCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch live revenue share config failed:', error)
    ElMessage.error(t('pages.liveRevenueShareCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (formData.anchorSharePercent < 0 || formData.anchorSharePercent > 100) {
    ElMessage.warning(t('pages.liveRevenueShareCfg.percentRangeInvalid'))
    return
  }
  if (formData.guildSharePercent < 0 || formData.guildSharePercent > 100) {
    ElMessage.warning(t('pages.liveRevenueShareCfg.percentRangeInvalid'))
    return
  }

  loading.value = true
  try {
    const response = await liveRevenueShareCfgApi.saveCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      anchorSharePercent: formData.anchorSharePercent,
      guildSharePercent: formData.guildSharePercent,
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.liveRevenueShareCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save live revenue share config failed:', error)
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
</style>

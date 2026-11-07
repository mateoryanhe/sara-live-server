<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GooglePlayCfgManagement') }}</span>
          <el-button
              v-if="can('sync')"
              :loading="syncingToTarget"
              type="warning"
              @click="handleSyncToTarget"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.googlePlay.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.googlePlay.noticeLine1') }}</p>
        <p>{{ t('pages.googlePlay.noticeLine2') }}</p>
        <p>{{ t('pages.googlePlay.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.googlePlay.packageName')" prop="packageName">
          <el-input
              v-model="formData.packageName"
              clearable
              :placeholder="t('pages.googlePlay.packageNamePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.googlePlay.packageNameTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.googlePlay.serviceAccountJson')" prop="serviceAccountJson">
          <el-input
              v-model="formData.serviceAccountJson"
              :rows="8"
              :placeholder="t('pages.googlePlay.serviceAccountJsonPlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.googlePlay.serviceAccountJsonTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.googlePlay.lastUpdated')">
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
import {dataSyncApi} from '@/api/modules/data-sync'
import {googlePlayApi} from '@/api/modules/google-play'
import {confirmDataSync} from '@/utils/confirm-data-sync'
import {usePagePermission} from '@/composables/usePagePermission'
import type {GooglePlayCfg} from '@/types/api'

const {t} = useI18n()
const {can} = usePagePermission('GooglePlayCfgManagement')
const loading = ref(false)
const syncingToTarget = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  packageName: '',
  serviceAccountJson: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const validateServiceAccountJson = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const json = value?.trim()
  if (!json) {
    callback(new Error(t('pages.googlePlay.serviceAccountJsonRequired')))
    return
  }
  try {
    JSON.parse(json)
  } catch {
    callback(new Error(t('pages.googlePlay.serviceAccountJsonInvalid')))
    return
  }
  callback()
}

const formRules = computed(() => ({
  packageName: [
    {required: true, message: t('pages.googlePlay.packageNameRequired'), trigger: 'blur'},
    {max: 128, message: t('pages.googlePlay.packageNameMaxLength'), trigger: 'blur'},
  ],
  serviceAccountJson: [{validator: validateServiceAccountJson, trigger: 'blur'}],
}))

const applyCfg = (cfg: GooglePlayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.packageName = ''
    formData.serviceAccountJson = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.packageName = cfg.packageName || ''
  formData.serviceAccountJson = cfg.serviceAccountJson || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await googlePlayApi.getGooglePlayCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch google play cfg failed:', error)
    ElMessage.error(t('pages.googlePlay.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await googlePlayApi.saveGooglePlayCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: true,
      packageName: formData.packageName.trim(),
      serviceAccountJson: formData.serviceAccountJson.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.googlePlay.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.googlePlay.saveFailed'))
    }
  } catch (error) {
    console.error('save google play cfg failed:', error)
  } finally {
    loading.value = false
  }
}

const handleSyncToTarget = async () => {
  if (!formData.id || formData.id === '0') {
    ElMessage.warning(t('pages.googlePlay.syncNeedSave'))
    return
  }
  try {
    await confirmDataSync({
      detail: t('pages.googlePlay.syncToTargetConfirm'),
      title: t('common.syncData'),
    })
  } catch {
    return
  }
  syncingToTarget.value = true
  try {
    const response = await dataSyncApi.syncGooglePlayCfg()
    if (response?.success) {
      ElMessage.success(response.message || t('pages.googlePlay.syncToTargetSuccess'))
    } else {
      ElMessage.error(t('pages.googlePlay.syncToTargetFailed'))
    }
  } catch (error) {
    console.error('sync google play cfg to target failed:', error)
    ElMessage.error(t('pages.googlePlay.syncToTargetFailedCheckConfig'))
  } finally {
    syncingToTarget.value = false
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
  max-width: 860px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}

code {
  padding: 0 4px;
  background: #f4f4f5;
  border-radius: 4px;
}
</style>

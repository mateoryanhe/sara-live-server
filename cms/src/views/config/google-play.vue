<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GooglePlayCfgManagement') }}</span>
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
import {googlePlayApi} from '@/api/modules/google-play'
import type {GooglePlayCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
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
  serviceAccountJson: [{validator: validateServiceAccountJson, trigger: 'blur'}],
}))

const applyCfg = (cfg: GooglePlayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.serviceAccountJson = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
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

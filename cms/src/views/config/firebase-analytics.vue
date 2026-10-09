<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.FirebaseAnalyticsCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.firebaseAnalytics.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.firebaseAnalytics.noticeLine1') }}</p>
        <p>{{ t('pages.firebaseAnalytics.noticeLine2') }}</p>
        <p>{{ t('pages.firebaseAnalytics.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.firebaseAnalytics.enabled')" prop="enabled">
          <el-switch v-model="formData.enabled" :active-value="1" :inactive-value="0"/>
        </el-form-item>
        <el-form-item :label="t('pages.firebaseAnalytics.projectId')" prop="projectId">
          <el-input
              v-model="formData.projectId"
              clearable
              :placeholder="t('pages.firebaseAnalytics.projectIdPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.firebaseAnalytics.clientConfigJson')" prop="clientConfigJson">
          <el-input
              v-model="formData.clientConfigJson"
              :rows="10"
              :placeholder="t('pages.firebaseAnalytics.clientConfigJsonPlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.firebaseAnalytics.clientConfigJsonTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.firebaseAnalytics.serviceAccountJson')" prop="serviceAccountJson">
          <el-input
              v-model="formData.serviceAccountJson"
              :rows="12"
              :placeholder="t('pages.firebaseAnalytics.serviceAccountJsonPlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.firebaseAnalytics.serviceAccountJsonTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.firebaseAnalytics.measurementApiSecret')" prop="measurementApiSecret">
          <el-input
              v-model="formData.measurementApiSecret"
              clearable
              show-password
              :placeholder="t('pages.firebaseAnalytics.measurementApiSecretPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.firebaseAnalytics.measurementApiSecretTip') }}</span>
        </el-form-item>
        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.firebaseAnalytics.lastUpdated')">
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
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {firebaseAnalyticsApi} from '@/api/modules/firebase-analytics'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  enabled: 0 as 0 | 1,
  projectId: '',
  clientConfigJson: '',
  serviceAccountJson: '',
  measurementApiSecret: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  projectId: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.firebaseAnalytics.projectIdRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
  clientConfigJson: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.firebaseAnalytics.clientConfigJsonRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
  serviceAccountJson: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.firebaseAnalytics.serviceAccountJsonRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
  measurementApiSecret: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.firebaseAnalytics.measurementApiSecretRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
}))

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await firebaseAnalyticsApi.getFirebaseAnalyticsCfg()
    const cfg = res.cfg
    if (!cfg) {
      formData.id = '0'
      formData.enabled = 0
      formData.projectId = ''
      formData.clientConfigJson = ''
      formData.serviceAccountJson = ''
      formData.measurementApiSecret = ''
      metaInfo.updatedAt = ''
      return
    }
    formData.id = cfg.id
    formData.enabled = cfg.enabled === 1 ? 1 : 0
    formData.projectId = cfg.projectId || ''
    formData.clientConfigJson = cfg.clientConfigJson || ''
    formData.serviceAccountJson = cfg.serviceAccountJson || ''
    formData.measurementApiSecret = cfg.measurementApiSecret || ''
    metaInfo.updatedAt = cfg.updatedAt || ''
  } catch (error) {
    console.error('load firebase analytics cfg failed:', error)
    ElMessage.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }
  loading.value = true
  try {
    await firebaseAnalyticsApi.saveFirebaseAnalyticsCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: formData.enabled,
      projectId: formData.projectId.trim(),
      clientConfigJson: formData.clientConfigJson.trim(),
      serviceAccountJson: formData.serviceAccountJson.trim(),
      measurementApiSecret: formData.measurementApiSecret.trim(),
    })
    ElMessage.success(t('pages.firebaseAnalytics.saveSuccess'))
    await fetchCfg()
  } catch (error) {
    console.error('save firebase analytics cfg failed:', error)
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
  font-size: 16px;
  font-weight: bold;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.cfg-form {
  max-width: 920px;
}

.form-tip {
  display: block;
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>

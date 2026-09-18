<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.FirebaseCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.firebase.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.firebase.noticeLine1') }}</p>
        <p>{{ t('pages.firebase.noticeLine2') }}</p>
        <p>{{ t('pages.firebase.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.firebase.projectId')" prop="projectId">
          <el-input
              v-model="formData.projectId"
              clearable
              :placeholder="t('pages.firebase.projectIdPlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('pages.firebase.serviceAccountJson')" prop="serviceAccountJson">
          <el-input
              v-model="formData.serviceAccountJson"
              :rows="12"
              :placeholder="t('pages.firebase.serviceAccountJsonPlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.firebase.serviceAccountJsonTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.firebase.clientConfigJson')" prop="clientConfigJson">
          <el-input
              v-model="formData.clientConfigJson"
              :rows="10"
              :placeholder="t('pages.firebase.clientConfigJsonPlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.firebase.clientConfigJsonTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.firebase.lastUpdated')">
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
import {firebaseApi} from '@/api/modules/firebase'
import type {FirebaseCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  projectId: '',
  serviceAccountJson: '',
  clientConfigJson: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const validateServiceAccountJson = (_: unknown, value: string, callback: (error?: Error) => void) => {
  const raw = value?.trim()
  if (!raw) {
    callback(new Error(t('pages.firebase.serviceAccountJsonRequired')))
    return
  }
  try {
    const parsed = JSON.parse(raw) as {project_id?: string}
    if (parsed.project_id !== formData.projectId.trim()) {
      callback(new Error(t('pages.firebase.serviceAccountProjectMismatch')))
      return
    }
  } catch {
    callback(new Error(t('pages.firebase.serviceAccountJsonInvalid')))
    return
  }
  callback()
}

const parseClientConfig = (value: string) => {
  const raw = value?.trim()
  if (!raw) throw new Error('empty client config')
  try {
    return JSON.parse(raw) as {apiKey?: string; projectId?: string; appId?: string}
  } catch {
    const firstBrace = raw.indexOf('{')
    const lastBrace = raw.lastIndexOf('}')
    if (firstBrace < 0 || lastBrace <= firstBrace) throw new Error('invalid client config')
    const objectLiteral = raw
        .slice(firstBrace, lastBrace + 1)
        .replace(/([{,]\s*)([A-Za-z_$][\w$]*)(\s*:)/g, '$1"$2"$3')
        .replace(/,\s*}/g, '}')
    return JSON.parse(objectLiteral) as {apiKey?: string; projectId?: string; appId?: string}
  }
}

const validateClientConfigJson = (_: unknown, value: string, callback: (error?: Error) => void) => {
  try {
    const parsed = parseClientConfig(value)
    if (!parsed.apiKey || !parsed.projectId || !parsed.appId) {
      callback(new Error(t('pages.firebase.clientConfigIncomplete')))
      return
    }
    if (parsed.projectId !== formData.projectId.trim()) {
      callback(new Error(t('pages.firebase.clientConfigProjectMismatch')))
      return
    }
  } catch {
    callback(new Error(t('pages.firebase.clientConfigJsonInvalid')))
    return
  }
  callback()
}

const formRules = computed(() => ({
  projectId: [{required: true, message: t('pages.firebase.projectIdRequired'), trigger: 'blur'}],
  serviceAccountJson: [{validator: validateServiceAccountJson, trigger: 'blur'}],
  clientConfigJson: [{validator: validateClientConfigJson, trigger: 'blur'}],
}))

const applyCfg = (cfg: FirebaseCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.projectId = ''
    formData.serviceAccountJson = ''
    formData.clientConfigJson = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.projectId = cfg.projectId || ''
  formData.serviceAccountJson = cfg.serviceAccountJson || ''
  formData.clientConfigJson = cfg.clientConfigJson || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await firebaseApi.getFirebaseCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch firebase cfg failed:', error)
    ElMessage.error(t('pages.firebase.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const clientConfigJson = JSON.stringify(parseClientConfig(formData.clientConfigJson), null, 2)
    formData.clientConfigJson = clientConfigJson
    const response = await firebaseApi.saveFirebaseCfg({
      id: Number(formData.id) || 0,
      projectId: formData.projectId.trim(),
      serviceAccountJson: formData.serviceAccountJson.trim(),
      clientConfigJson,
    })
    if (response?.success) {
      ElMessage.success(t('pages.firebase.saveSuccess'))
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.firebase.saveFailed'))
    }
  } catch (error) {
    console.error('save firebase cfg failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchCfg)
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
  line-height: 1.5;
}

.cfg-form {
  max-width: 860px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style>

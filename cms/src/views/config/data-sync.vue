<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.DataSyncCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.dataSync.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.dataSync.noticeBody') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.dataSync.targetApiBase')" prop="targetApiBase">
          <el-input
              v-model="formData.targetApiBase"
              clearable
              :placeholder="t('pages.dataSync.targetApiBasePlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('pages.dataSync.syncToken')" prop="token">
          <div class="token-input">
            <el-input
                v-model="formData.token"
                clearable
                :placeholder="t('pages.dataSync.syncTokenPlaceholder')"
                show-password
                type="password"
            />
            <el-button @click="generateSyncToken">{{ t('pages.dataSync.generateToken') }}</el-button>
            <el-button :disabled="!formData.token" @click="copySyncToken">{{ t('common.copy') }}</el-button>
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.dataSync.lastUpdated')">
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
import type {DataSyncCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  targetApiBase: '',
  token: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  targetApiBase: [{required: true, message: t('pages.dataSync.targetApiBaseRequired'), trigger: 'blur'}],
  token: [{required: true, message: t('pages.dataSync.syncTokenRequired'), trigger: 'blur'}],
}))

const generateSyncToken = () => {
  formData.token = crypto.randomUUID().replace(/-/g, '')
}

const copySyncToken = async () => {
  const value = formData.token.trim()
  if (!value) {
    ElMessage.warning(t('pages.dataSync.tokenEmptyCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(t('pages.dataSync.copiedToClipboard'))
  } catch (error) {
    console.error('copy token failed:', error)
    ElMessage.error(t('pages.dataSync.copyFailed'))
  }
}

const applyCfg = (cfg: DataSyncCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.targetApiBase = ''
    formData.token = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.targetApiBase = cfg.targetApiBase || ''
  formData.token = cfg.token || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await dataSyncApi.getDataSyncCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch data sync cfg failed:', error)
    ElMessage.error(t('pages.dataSync.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await dataSyncApi.saveDataSyncCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      targetApiBase: formData.targetApiBase.trim(),
      token: formData.token.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.dataSync.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.dataSync.saveFailed'))
    }
  } catch (error) {
    console.error('save data sync cfg failed:', error)
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

.token-input {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.token-input .el-input {
  flex: 1;
}
</style>

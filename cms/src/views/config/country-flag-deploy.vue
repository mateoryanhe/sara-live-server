<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.CountryFlagDeployManagement') }}</span>
      </template>

      <el-descriptions v-if="deployInfo" :column="1" border class="deploy-info">
        <el-descriptions-item :label="t('pages.countryFlagDeploy.version')">
          {{ deployInfo.version || t('pages.countryFlagDeploy.emptyVersion') }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.countryFlagDeploy.urlPrefix')">
          {{ deployInfo.urlPrefix }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.countryFlagDeploy.deployPath')">
          {{ deployInfo.deployPath }}
        </el-descriptions-item>
        <el-descriptions-item v-if="deployInfo.updatedAt" :label="t('pages.countryFlagDeploy.lastUpdated')">
          {{ deployInfo.updatedAt }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="upload-section">
        <div class="section-title">{{ t('pages.countryFlagDeploy.uploadTitle') }}</div>
        <div class="form-tip">{{ t('pages.countryFlagDeploy.uploadTip') }}</div>

        <el-upload
            v-if="can('deploy')"
            :auto-upload="false"
            :disabled="uploading"
            :limit="1"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
            :on-remove="handleRemove"
            accept=".zip,application/zip"
            action="#"
            class="zip-uploader"
            drag
        >
          <el-icon class="upload-icon">
            <UploadFilled/>
          </el-icon>
          <div class="el-upload__text">
            {{ t('pages.countryFlagDeploy.dragTip') }}
          </div>
          <template #tip>
            <div class="form-tip">{{ t('pages.countryFlagDeploy.fileTip') }}</div>
          </template>
        </el-upload>

        <el-progress
            v-if="uploading || uploadPercent > 0"
            :percentage="uploadPercent"
            :status="uploadStatus"
            class="upload-progress"
        />

        <div v-if="can('deploy')" class="action-row">
          <el-button
              :disabled="!selectedFile"
              :loading="uploading"
              type="primary"
              @click="handleDeploy"
          >
            {{ t('pages.countryFlagDeploy.deployBtn') }}
          </el-button>
        </div>

        <el-alert
            v-if="deployResult"
            :closable="false"
            :title="deployResultTitle"
            class="deploy-result"
            show-icon
            type="success"
        >
          <template #default>
            <div>{{ t('pages.countryFlagDeploy.resultVersion', {version: deployResult.version}) }}</div>
            <div>{{ t('pages.countryFlagDeploy.resultUrl', {url: deployResult.urlPrefix}) }}</div>
            <div>{{ t('pages.countryFlagDeploy.resultPath', {path: deployResult.deployPath}) }}</div>
            <div>{{ t('pages.countryFlagDeploy.resultFiles', {count: deployResult.fileCount}) }}</div>
            <div>{{ t('pages.countryFlagDeploy.resultRemoved', {count: deployResult.removed}) }}</div>
          </template>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref} from 'vue'
import {UploadFilled} from '@element-plus/icons-vue'
import type {UploadFile, UploadFiles} from 'element-plus'
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {countryFlagDeployApi} from '@/api/modules/country-flag-deploy'
import type {CountryFlagDeployInfo, DeployCountryFlagZipRes} from '@/api/modules/country-flag-deploy'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('CountryFlagDeployManagement')

const loading = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')
const deployInfo = ref<CountryFlagDeployInfo | null>(null)
const selectedFile = ref<File | null>(null)
const deployResult = ref<DeployCountryFlagZipRes | null>(null)

const deployResultTitle = computed(() => t('pages.countryFlagDeploy.deploySuccess'))

const fetchDeployInfo = async () => {
  loading.value = true
  try {
    const response = await countryFlagDeployApi.getCountryFlagDeployInfo()
    deployInfo.value = response.info ?? null
  } catch {
    ElMessage.error(t('pages.countryFlagDeploy.fetchInfoFailed'))
  } finally {
    loading.value = false
  }
}

const handleFileChange = (uploadFile: UploadFile, _uploadFiles: UploadFiles) => {
  deployResult.value = null
  uploadPercent.value = 0
  uploadStatus.value = ''
  selectedFile.value = uploadFile.raw ?? null
}

const handleRemove = () => {
  selectedFile.value = null
  uploadPercent.value = 0
  uploadStatus.value = ''
}

const handleExceed = () => {
  ElMessage.warning(t('pages.countryFlagDeploy.singleFileOnly'))
}

const handleDeploy = async () => {
  if (!selectedFile.value) {
    return
  }
  const fileName = selectedFile.value.name.toLowerCase()
  if (!fileName.endsWith('.zip')) {
    ElMessage.error(t('pages.countryFlagDeploy.zipOnly'))
    return
  }

  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  deployResult.value = null
  try {
    const response = await countryFlagDeployApi.deployZip(selectedFile.value, (percent) => {
      uploadPercent.value = percent
    })
    deployResult.value = response
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.countryFlagDeploy.deploySuccess'))
    await fetchDeployInfo()
  } catch {
    uploadStatus.value = 'exception'
    ElMessage.error(t('pages.countryFlagDeploy.deployFailed'))
  } finally {
    uploading.value = false
  }
}

onMounted(() => {
  fetchDeployInfo()
})
</script>

<style scoped>
.deploy-info {
  margin-bottom: 24px;
}

.upload-section {
  margin-top: 8px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.form-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 12px;
}

.zip-uploader {
  width: 100%;
}

.upload-icon {
  color: var(--el-color-primary);
  font-size: 48px;
  margin-bottom: 8px;
}

.upload-progress {
  margin-top: 16px;
}

.action-row {
  margin-top: 16px;
}

.deploy-result {
  margin-top: 16px;
}
</style>

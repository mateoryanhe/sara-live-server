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
        <el-descriptions-item :label="t('pages.countryFlagDeploy.fileCount')">
          {{ deployInfo.fileCount ?? 0 }}
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

      <div class="preview-section">
        <div class="section-title">{{ t('pages.countryFlagDeploy.previewTitle') }}</div>
        <div class="preview-toolbar">
          <el-input
              v-model="flagKeyword"
              :placeholder="t('pages.countryFlagDeploy.previewSearch')"
              clearable
              class="preview-search"
          />
          <span class="preview-count">
            {{ t('pages.countryFlagDeploy.previewCount', {total: flagTotal, shown: filteredFlags.length}) }}
          </span>
        </div>
        <el-empty
            v-if="flagTotal === 0"
            :description="t('pages.countryFlagDeploy.previewEmpty')"
        />
        <el-empty
            v-else-if="filteredFlags.length === 0"
            :description="t('pages.countryFlagDeploy.previewFilteredEmpty')"
        />
        <div v-else class="flag-grid">
          <div v-for="item in filteredFlags" :key="item.code" class="flag-card">
            <el-image
                :alt="item.code"
                :preview-src-list="[item.icon]"
                :src="item.icon"
                class="flag-img"
                fit="contain"
                preview-teleported
            >
              <template #error>
                <div class="flag-img-error">{{ item.code }}</div>
              </template>
            </el-image>
            <div class="flag-code">{{ item.code }}</div>
            <div class="flag-name" :title="flagDisplayName(item)">{{ flagDisplayName(item) }}</div>
            <div class="flag-file">{{ item.file }}</div>
          </div>
        </div>
      </div>

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
import type {
  CountryFlagDeployInfo,
  CountryFlagPreviewItem,
  DeployCountryFlagZipRes,
} from '@/api/modules/country-flag-deploy'
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
const flagKeyword = ref('')

const deployResultTitle = computed(() => t('pages.countryFlagDeploy.deploySuccess'))

const flagTotal = computed(() => deployInfo.value?.flags?.length ?? 0)

const filteredFlags = computed(() => {
  const list = deployInfo.value?.flags ?? []
  const kw = flagKeyword.value.trim().toLowerCase()
  if (!kw) {
    return list
  }
  return list.filter((item) => {
    return (
        item.code.toLowerCase().includes(kw) ||
        item.file.toLowerCase().includes(kw) ||
        item.nameEn.toLowerCase().includes(kw) ||
        item.nameZh.toLowerCase().includes(kw)
    )
  })
})

const flagDisplayName = (item: CountryFlagPreviewItem) => {
  if (item.nameZh && item.nameEn) {
    return `${item.nameZh} / ${item.nameEn}`
  }
  return item.nameZh || item.nameEn || item.code
}

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

.preview-section {
  margin-bottom: 28px;
}

.preview-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.preview-search {
  max-width: 360px;
}

.preview-count {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.flag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.flag-card {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 10px 8px;
  text-align: center;
  background: var(--el-bg-color);
}

.flag-img {
  width: 64px;
  height: 42px;
  margin: 0 auto 8px;
}

.flag-img-error {
  width: 64px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.flag-code {
  font-weight: 600;
  font-size: 14px;
}

.flag-name {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-regular);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.flag-file {
  margin-top: 2px;
  font-size: 11px;
  color: var(--el-text-color-secondary);
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

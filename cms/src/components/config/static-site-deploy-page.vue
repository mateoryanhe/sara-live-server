<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t(titleKey) }}</span>
      </template>

      <el-descriptions v-if="deployInfo" :column="1" border class="deploy-info">
        <el-descriptions-item :label="t('pages.officialSiteDeploy.domain')">
          <el-link
              v-if="deployInfo.domain"
              :href="`https://${deployInfo.domain}`"
              target="_blank"
              type="primary"
          >
            {{ deployInfo.domain }}
          </el-link>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.officialSiteDeploy.urlPrefix')">
          {{ deployInfo.urlPrefix }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.officialSiteDeploy.deployPath')">
          {{ deployInfo.deployPath }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.officialSiteDeploy.lastUploadAt')">
          {{ deployInfo.lastUploadAt || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form
          v-if="deployApi.saveMapping"
          ref="mappingFormRef"
          :model="mappingForm"
          :rules="mappingRules"
          class="mapping-form"
          label-width="120px"
      >
        <div class="section-title">{{ t('pages.officialSiteDeploy.siteMappingTitle') }}</div>
        <div class="form-tip">{{ t('pages.officialSiteDeploy.siteMappingTip') }}</div>

        <el-form-item :label="t('pages.officialSiteDeploy.domain')" prop="domain">
          <el-input
              v-model="mappingForm.domain"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.officialSiteDeploy.domainPlaceholder', {domain: domainExample ?? '-'})"
          />
        </el-form-item>

        <el-form-item :label="t('pages.officialSiteDeploy.deployPath')" prop="deployPath">
          <el-input
              v-model="mappingForm.deployPath"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.officialSiteDeploy.deployPathPlaceholder', {path: defaultDeployPath})"
          />
        </el-form-item>

        <el-form-item v-if="can('save')">
          <el-button :loading="savingMapping" type="primary" @click="handleSaveMapping">
            {{ t('pages.officialSiteDeploy.saveAndRefresh') }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="upload-section">
        <div class="section-title">{{ t('pages.officialSiteDeploy.uploadTitle') }}</div>
        <div class="form-tip">
          {{ t(enableChunkedZipUpload ? 'pages.officialSiteDeploy.chunkedZipUploadTip' : 'pages.officialSiteDeploy.uploadTip', {site: t(titleKey)}) }}
        </div>

        <el-upload
            v-if="can('deploy')"
            :auto-upload="false"
            :disabled="uploading || apkUploading"
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
            {{ t('pages.officialSiteDeploy.dragTip') }}
          </div>
          <template #tip>
            <div class="form-tip">{{ t('pages.officialSiteDeploy.fileTip') }}</div>
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
              :disabled="!selectedFile || apkUploading"
              :loading="uploading"
              type="primary"
              @click="handleDeploy"
          >
            {{ t('pages.officialSiteDeploy.deployBtn') }}
          </el-button>
        </div>

        <el-alert
            v-if="deployResult"
            :closable="false"
            :title="t('pages.officialSiteDeploy.deploySuccess')"
            class="deploy-result"
            show-icon
            type="success"
        >
          <template #default>
            <div>{{ t('pages.officialSiteDeploy.resultPath', {path: deployResult.deployPath}) }}</div>
            <div>{{ t('pages.officialSiteDeploy.resultFiles', {count: deployResult.fileCount}) }}</div>
            <div>{{ t('pages.officialSiteDeploy.resultDirs', {count: deployResult.dirCount}) }}</div>
          </template>
        </el-alert>

        <template v-if="enableApkUpload">
          <el-divider/>
          <div class="section-title">{{ t('pages.officialSiteDeploy.apkUploadTitle') }}</div>
          <div class="form-tip">
            {{ t('pages.officialSiteDeploy.apkUploadTip', {path: apkDeployPath}) }}
          </div>

          <el-upload
              v-if="can('deploy')"
              :auto-upload="false"
              :disabled="apkUploading || uploading"
              :limit="1"
              :on-change="handleApkFileChange"
              :on-exceed="handleApkExceed"
              :on-remove="handleApkRemove"
              accept=".apk,application/vnd.android.package-archive"
              action="#"
              class="apk-uploader"
              drag
          >
            <el-icon class="upload-icon">
              <UploadFilled/>
            </el-icon>
            <div class="el-upload__text">
              {{ t('pages.officialSiteDeploy.apkDragTip') }}
            </div>
            <template #tip>
              <div class="form-tip">{{ t('pages.officialSiteDeploy.apkFileTip') }}</div>
            </template>
          </el-upload>

          <el-progress
              v-if="apkUploading || apkUploadPercent > 0"
              :percentage="apkUploadPercent"
              :status="apkUploadStatus"
              class="upload-progress"
          />

          <div v-if="can('deploy')" class="action-row">
            <el-button
                :disabled="!selectedApkFile || uploading"
                :loading="apkUploading"
                type="primary"
                @click="handleApkUpload"
            >
              {{ t('pages.officialSiteDeploy.apkUploadBtn') }}
            </el-button>
          </div>

          <el-alert
              v-if="apkUploadResult"
              :closable="false"
              :title="t('pages.officialSiteDeploy.apkUploadSuccess')"
              class="deploy-result"
              show-icon
              type="success"
          >
            <template #default>
              <div>{{ t('pages.officialSiteDeploy.apkResultFile', {file: apkUploadResult.fileName}) }}</div>
              <div>{{ t('pages.officialSiteDeploy.apkResultSize', {size: formatFileSize(apkUploadResult.fileSize)}) }}</div>
              <div>{{ t('pages.officialSiteDeploy.resultPath', {path: apkUploadResult.deployPath}) }}</div>
              <div v-if="apkUploadResult.downloadUrl">
                {{ t('pages.officialSiteDeploy.apkDownloadUrl') }}：
                <el-link :href="apkUploadResult.downloadUrl" target="_blank" type="primary">
                  {{ apkUploadResult.downloadUrl }}
                </el-link>
              </div>
            </template>
          </el-alert>
        </template>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {UploadFilled} from '@element-plus/icons-vue'
import type {FormInstance, FormRules, UploadFile, UploadFiles} from 'element-plus'
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import type {
  CompleteStaticSiteFileUploadRes,
  DeployStaticSiteZipRes,
  StaticSiteDeployApi,
  StaticSiteDeployInfo,
} from '@/api/modules/official-site-deploy'
import {usePagePermission} from '@/composables/usePagePermission'

const props = defineProps<{
  pageName: string
  titleKey: string
  deployApi: StaticSiteDeployApi
  defaultDeployPath?: string
  domainExample?: string
  apkSubdirectory?: string
  enableApkUpload?: boolean
  enableChunkedZipUpload?: boolean
}>()

const {t} = useI18n()
const {can} = usePagePermission(props.pageName)

const loading = ref(false)
const savingMapping = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')
const deployInfo = ref<StaticSiteDeployInfo | null>(null)
const selectedFile = ref<File | null>(null)
const deployResult = ref<DeployStaticSiteZipRes | null>(null)
const apkUploading = ref(false)
const apkUploadPercent = ref(0)
const apkUploadStatus = ref<'success' | 'exception' | ''>('')
const selectedApkFile = ref<File | null>(null)
const apkUploadResult = ref<CompleteStaticSiteFileUploadRes | null>(null)
const mappingFormRef = ref<FormInstance>()
const mappingForm = reactive({
  domain: '',
  deployPath: props.defaultDeployPath ?? '',
})
const mappingRules = computed<FormRules>(() => ({
  domain: [{required: true, message: t('pages.officialSiteDeploy.domainRequired'), trigger: 'blur'}],
  deployPath: [{required: true, message: t('pages.officialSiteDeploy.deployPathRequired'), trigger: 'blur'}],
}))

const apkDeployPath = computed(() => {
  const deployPath = deployInfo.value?.deployPath?.replace(/[\\/]+$/, '') ?? ''
  const subdirectory = props.apkSubdirectory?.replace(/^[\\/]+|[\\/]+$/g, '') ?? ''
  if (!deployPath) {
    return subdirectory || '-'
  }
  if (!subdirectory) {
    return deployPath
  }
  const separator = deployPath.includes('\\') ? '\\' : '/'
  return `${deployPath}${separator}${subdirectory}`
})

const fetchDeployInfo = async () => {
  loading.value = true
  try {
    const response = await props.deployApi.getDeployInfo()
    deployInfo.value = response.info ?? null
    mappingForm.domain = response.info?.domain ?? ''
    mappingForm.deployPath = response.info?.deployPath ?? props.defaultDeployPath ?? ''
  } catch {
    ElMessage.error(t('pages.officialSiteDeploy.fetchInfoFailed'))
  } finally {
    loading.value = false
  }
}

const handleSaveMapping = async () => {
  const saveMapping = props.deployApi.saveMapping
  if (!saveMapping || !mappingFormRef.value) {
    return
  }
  try {
    await mappingFormRef.value.validate()
  } catch {
    return
  }
  savingMapping.value = true
  try {
    await saveMapping(mappingForm.domain.trim(), mappingForm.deployPath.trim())
    ElMessage.success(t('pages.officialSiteDeploy.saveSuccess'))
    await fetchDeployInfo()
  } catch {
    ElMessage.error(t('pages.officialSiteDeploy.saveFailed'))
  } finally {
    savingMapping.value = false
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
  ElMessage.warning(t('pages.officialSiteDeploy.singleFileOnly'))
}

const handleApkFileChange = (uploadFile: UploadFile, _uploadFiles: UploadFiles) => {
  apkUploadResult.value = null
  apkUploadPercent.value = 0
  apkUploadStatus.value = ''
  selectedApkFile.value = uploadFile.raw ?? null
}

const handleApkRemove = () => {
  selectedApkFile.value = null
  apkUploadPercent.value = 0
  apkUploadStatus.value = ''
}

const handleApkExceed = () => {
  ElMessage.warning(t('pages.officialSiteDeploy.singleApkOnly'))
}

const wait = (milliseconds: number) => new Promise<void>((resolve) => {
  window.setTimeout(resolve, milliseconds)
})

const uploadFileInChunks = async (
    file: File,
    onProgress: (percent: number) => void,
): Promise<CompleteStaticSiteFileUploadRes> => {
  const initFileUpload = props.deployApi.initFileUpload
  const uploadFileChunk = props.deployApi.uploadFileChunk
  const completeFileUpload = props.deployApi.completeFileUpload
  const abortFileUpload = props.deployApi.abortFileUpload
  if (!initFileUpload || !uploadFileChunk || !completeFileUpload || !abortFileUpload) {
    throw new Error('chunk upload api is unavailable')
  }

  const session = await initFileUpload(file.name, file.size)
  if (!session.uploadId || session.chunkSize <= 0 || session.totalChunks <= 0) {
    throw new Error('invalid chunk upload session')
  }

  try {
    for (let chunkIndex = 0; chunkIndex < session.totalChunks; chunkIndex += 1) {
      const start = chunkIndex * session.chunkSize
      const end = Math.min(start + session.chunkSize, file.size)
      const chunk = file.slice(start, end)
      let uploaded = false
      let lastError: unknown
      for (let attempt = 1; attempt <= 3; attempt += 1) {
        try {
          await uploadFileChunk(session.uploadId, chunkIndex, chunk, file.name, (chunkPercent) => {
            const loadedBytes = start + (chunk.size * chunkPercent) / 100
            onProgress(Math.min(99, Math.round((loadedBytes * 100) / file.size)))
          })
          uploaded = true
          break
        } catch (error) {
          lastError = error
          if (attempt < 3) {
            await wait(attempt * 500)
          }
        }
      }
      if (!uploaded) {
        throw lastError instanceof Error ? lastError : new Error('upload chunk failed')
      }
      onProgress(Math.min(99, Math.round((end * 100) / file.size)))
    }
    const result = await completeFileUpload(session.uploadId)
    onProgress(100)
    return result
  } catch (error) {
    try {
      await abortFileUpload(session.uploadId)
    } catch {
      // 过期分片由服务端定期清理，不覆盖原始上传错误。
    }
    throw error
  }
}

const handleDeploy = async () => {
  if (!selectedFile.value) {
    return
  }
  if (!selectedFile.value.name.toLowerCase().endsWith('.zip')) {
    ElMessage.error(t('pages.officialSiteDeploy.zipOnly'))
    return
  }

  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  deployResult.value = null
  try {
    if (props.enableChunkedZipUpload) {
      const response = await uploadFileInChunks(selectedFile.value, (percent) => {
        uploadPercent.value = percent
      })
      deployResult.value = {
        fileCount: response.fileCount,
        dirCount: response.dirCount,
        deployPath: response.deployPath,
        urlPrefix: response.urlPrefix,
      }
    } else {
      deployResult.value = await props.deployApi.deployZip(selectedFile.value, (percent) => {
        uploadPercent.value = percent
      })
    }
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.officialSiteDeploy.deploySuccess'))
    await fetchDeployInfo()
  } catch {
    uploadStatus.value = 'exception'
    ElMessage.error(t('pages.officialSiteDeploy.deployFailed'))
  } finally {
    uploading.value = false
  }
}

const handleApkUpload = async () => {
  const file = selectedApkFile.value
  if (!file) {
    return
  }
  if (!file.name.toLowerCase().endsWith('.apk')) {
    ElMessage.error(t('pages.officialSiteDeploy.apkOnly'))
    return
  }

  apkUploading.value = true
  apkUploadPercent.value = 0
  apkUploadStatus.value = ''
  apkUploadResult.value = null
  try {
    apkUploadResult.value = await uploadFileInChunks(file, (percent) => {
      apkUploadPercent.value = percent
    })
    apkUploadStatus.value = 'success'
    ElMessage.success(t('pages.officialSiteDeploy.apkUploadSuccess'))
  } catch {
    apkUploadStatus.value = 'exception'
    ElMessage.error(t('pages.officialSiteDeploy.apkUploadFailed'))
  } finally {
    apkUploading.value = false
  }
}

const formatFileSize = (fileSize: number) => {
  if (fileSize >= 1024 * 1024 * 1024) {
    return `${(fileSize / (1024 * 1024 * 1024)).toFixed(2)} GB`
  }
  if (fileSize >= 1024 * 1024) {
    return `${(fileSize / (1024 * 1024)).toFixed(2)} MB`
  }
  return `${(fileSize / 1024).toFixed(2)} KB`
}

onMounted(fetchDeployInfo)
</script>

<style scoped>
.deploy-info {
  margin-bottom: 24px;
}

.mapping-form {
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

.zip-uploader,
.apk-uploader {
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

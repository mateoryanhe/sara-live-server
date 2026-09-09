<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.CoinMerchantDeployManagement') }}</span>
      </template>

      <el-descriptions v-if="deployInfo" :column="1" border class="deploy-info">
        <el-descriptions-item :label="t('pages.coinMerchantDeploy.urlPrefix')">
          {{ deployInfo.urlPrefix }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.coinMerchantDeploy.deployPath')">
          {{ deployInfo.deployPath }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="120px">
        <div class="section-title">{{ t('pages.coinMerchantDeploy.secretTitle') }}</div>
        <div class="form-tip">{{ t('pages.coinMerchantDeploy.secretTip') }}</div>

        <el-form-item :label="t('pages.coinMerchantDeploy.deploySecret')" prop="deploySecret">
          <div class="secret-input">
            <el-input
                v-model="formData.deploySecret"
                clearable
                :disabled="!can('save')"
                :placeholder="t('pages.coinMerchantDeploy.deploySecretPlaceholder')"
                show-password
                type="password"
            />
            <el-button v-if="can('save')" @click="generateDeploySecret">
              {{ t('pages.coinMerchantDeploy.generateSecret') }}
            </el-button>
            <el-button :disabled="!formData.deploySecret" @click="copyDeploySecret">
              {{ t('common.copy') }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.coinMerchantDeploy.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item v-if="can('save')">
          <el-button :loading="saving" type="primary" @click="handleSaveSecret">
            {{ t('common.saveConfig') }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="upload-section">
        <div class="section-title">{{ t('pages.coinMerchantDeploy.uploadTitle') }}</div>
        <div class="form-tip">{{ t('pages.coinMerchantDeploy.uploadTip') }}</div>

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
            {{ t('pages.coinMerchantDeploy.dragTip') }}
          </div>
          <template #tip>
            <div class="form-tip">{{ t('pages.coinMerchantDeploy.fileTip') }}</div>
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
            {{ t('pages.coinMerchantDeploy.deployBtn') }}
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
            <div>{{ t('pages.coinMerchantDeploy.resultPath', {path: deployResult.deployPath}) }}</div>
            <div>{{ t('pages.coinMerchantDeploy.resultFiles', {count: deployResult.fileCount}) }}</div>
            <div>{{ t('pages.coinMerchantDeploy.resultDirs', {count: deployResult.dirCount}) }}</div>
          </template>
        </el-alert>
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
import {coinMerchantDeployApi} from '@/api/modules/coin-merchant-deploy'
import type {DeployCoinMerchantZipRes, CoinMerchantDeployInfo} from '@/api/modules/coin-merchant-deploy'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('CoinMerchantDeployManagement')

const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')
const deployInfo = ref<CoinMerchantDeployInfo | null>(null)
const selectedFile = ref<File | null>(null)
const deployResult = ref<DeployCoinMerchantZipRes | null>(null)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: 0,
  deploySecret: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  deploySecret: [{required: true, message: t('pages.coinMerchantDeploy.deploySecretRequired'), trigger: 'blur'}],
}))

const deployResultTitle = computed(() => t('pages.coinMerchantDeploy.deploySuccess'))

const applyDeployInfo = (info: CoinMerchantDeployInfo | null | undefined) => {
  deployInfo.value = info ?? null
  formData.id = info?.id ? Number(info.id) : 0
  formData.deploySecret = info?.deploySecret ?? ''
  metaInfo.updatedAt = info?.updatedAt ?? ''
}

const fetchDeployInfo = async () => {
  loading.value = true
  try {
    const response = await coinMerchantDeployApi.getCoinMerchantDeployInfo()
    applyDeployInfo(response.info)
  } catch {
    ElMessage.error(t('pages.coinMerchantDeploy.fetchInfoFailed'))
  } finally {
    loading.value = false
  }
}

const generateDeploySecret = () => {
  formData.deploySecret = crypto.randomUUID().replace(/-/g, '')
}

const copyDeploySecret = async () => {
  const value = formData.deploySecret.trim()
  if (!value) {
    ElMessage.warning(t('pages.coinMerchantDeploy.secretEmptyCopy'))
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(t('pages.coinMerchantDeploy.copiedToClipboard'))
  } catch (error) {
    console.error('copy deploy secret failed:', error)
    ElMessage.error(t('pages.coinMerchantDeploy.copyFailed'))
  }
}

const handleSaveSecret = async () => {
  if (!formRef.value) {
    return
  }
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    await coinMerchantDeployApi.saveCoinMerchantDeployCfg({
      id: formData.id,
      deploySecret: formData.deploySecret.trim(),
    })
    ElMessage.success(t('pages.coinMerchantDeploy.saveSuccess'))
    await fetchDeployInfo()
  } catch {
    ElMessage.error(t('pages.coinMerchantDeploy.saveFailed'))
  } finally {
    saving.value = false
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
  ElMessage.warning(t('pages.coinMerchantDeploy.singleFileOnly'))
}

const handleDeploy = async () => {
  if (!selectedFile.value) {
    return
  }
  const fileName = selectedFile.value.name.toLowerCase()
  if (!fileName.endsWith('.zip')) {
    ElMessage.error(t('pages.coinMerchantDeploy.zipOnly'))
    return
  }

  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  deployResult.value = null
  try {
    const response = await coinMerchantDeployApi.deployZip(selectedFile.value, (percent) => {
      uploadPercent.value = percent
    })
    deployResult.value = response
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.coinMerchantDeploy.deploySuccess'))
  } catch {
    uploadStatus.value = 'exception'
    ElMessage.error(t('pages.coinMerchantDeploy.deployFailed'))
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

.cfg-form {
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

.secret-input {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  width: 100%;
}

.secret-input :deep(.el-input) {
  flex: 1;
  min-width: 260px;
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

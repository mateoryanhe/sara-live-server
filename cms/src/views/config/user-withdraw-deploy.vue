<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.UserWithdrawDeployManagement') }}</span>
      </template>

      <el-descriptions v-if="deployInfo" :column="1" border class="deploy-info">
        <el-descriptions-item :label="t('pages.userWithdrawDeploy.urlPrefix')">
          {{ deployInfo.urlPrefix }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.userWithdrawDeploy.deployPath')">
          {{ deployInfo.deployPath }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.userWithdrawDeploy.lastUploadAt')">
          {{ deployInfo.lastUploadAt || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="120px">
        <div class="section-title">{{ t('pages.userWithdrawDeploy.siteMappingTitle') }}</div>
        <div class="form-tip">{{ t('pages.userWithdrawDeploy.siteMappingTip') }}</div>

        <el-form-item :label="t('pages.userWithdrawDeploy.domain')" prop="domain">
          <el-input
              v-model="formData.domain"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.userWithdrawDeploy.domainPlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('pages.userWithdrawDeploy.deployPath')" prop="deployPath">
          <el-input
              v-model="formData.deployPath"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.userWithdrawDeploy.deployPathPlaceholder')"
          />
        </el-form-item>

        <el-form-item v-if="can('save')">
          <el-button :loading="saving" type="primary" @click="handleSaveMapping">
            {{ t('pages.userWithdrawDeploy.saveAndRefresh') }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="upload-section">
        <div class="section-title">{{ t('pages.userWithdrawDeploy.uploadTitle') }}</div>
        <div class="form-tip">{{ t('pages.userWithdrawDeploy.uploadTip') }}</div>

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
            {{ t('pages.userWithdrawDeploy.dragTip') }}
          </div>
          <template #tip>
            <div class="form-tip">{{ t('pages.userWithdrawDeploy.fileTip') }}</div>
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
            {{ t('pages.userWithdrawDeploy.deployBtn') }}
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
            <div>{{ t('pages.userWithdrawDeploy.resultPath', {path: deployResult.deployPath}) }}</div>
            <div>{{ t('pages.userWithdrawDeploy.resultFiles', {count: deployResult.fileCount}) }}</div>
            <div>{{ t('pages.userWithdrawDeploy.resultDirs', {count: deployResult.dirCount}) }}</div>
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
import {userWithdrawDeployApi} from '@/api/modules/user-withdraw-deploy'
import type {DeployUserWithdrawZipRes, UserWithdrawDeployInfo} from '@/api/modules/user-withdraw-deploy'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('UserWithdrawDeployManagement')

const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')
const deployInfo = ref<UserWithdrawDeployInfo | null>(null)
const selectedFile = ref<File | null>(null)
const deployResult = ref<DeployUserWithdrawZipRes | null>(null)
const formRef = ref<FormInstance>()

const formData = reactive({
  domain: '',
  deployPath: '/home/ec2-user/cdn/user-withdraw',
})

const formRules = computed<FormRules>(() => ({
  domain: [{required: true, message: t('pages.userWithdrawDeploy.domainRequired'), trigger: 'blur'}],
  deployPath: [{required: true, message: t('pages.userWithdrawDeploy.deployPathRequired'), trigger: 'blur'}],
}))

const deployResultTitle = computed(() => t('pages.userWithdrawDeploy.deploySuccess'))

const applyDeployInfo = (info: UserWithdrawDeployInfo | null | undefined) => {
  deployInfo.value = info ?? null
  formData.domain = info?.domain ?? ''
  formData.deployPath = info?.deployPath ?? '/home/ec2-user/cdn/user-withdraw'
}

const fetchDeployInfo = async () => {
  loading.value = true
  try {
    const response = await userWithdrawDeployApi.getUserWithdrawDeployInfo()
    applyDeployInfo(response.info)
  } catch {
    ElMessage.error(t('pages.userWithdrawDeploy.fetchInfoFailed'))
  } finally {
    loading.value = false
  }
}

const handleSaveMapping = async () => {
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
    await userWithdrawDeployApi.saveUserWithdrawDeployCfg({
      domain: formData.domain.trim(),
      deployPath: formData.deployPath.trim(),
    })
    ElMessage.success(t('pages.userWithdrawDeploy.saveSuccess'))
    await fetchDeployInfo()
  } catch {
    ElMessage.error(t('pages.userWithdrawDeploy.saveFailed'))
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
  ElMessage.warning(t('pages.userWithdrawDeploy.singleFileOnly'))
}

const handleDeploy = async () => {
  if (!selectedFile.value) {
    return
  }
  const fileName = selectedFile.value.name.toLowerCase()
  if (!fileName.endsWith('.zip')) {
    ElMessage.error(t('pages.userWithdrawDeploy.zipOnly'))
    return
  }

  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  deployResult.value = null
  try {
    const response = await userWithdrawDeployApi.deployZip(selectedFile.value, (percent) => {
      uploadPercent.value = percent
    })
    deployResult.value = response
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.userWithdrawDeploy.deploySuccess'))
    await fetchDeployInfo()
  } catch {
    uploadStatus.value = 'exception'
    ElMessage.error(t('pages.userWithdrawDeploy.deployFailed'))
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

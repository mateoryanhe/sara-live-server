<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.ShareActivityDeployManagement') }}</span>
      </template>

      <el-descriptions v-if="deployInfo" :column="1" border class="deploy-info">
        <el-descriptions-item :label="t('pages.shareActivityDeploy.urlPrefix')">
          {{ deployInfo.urlPrefix }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.shareActivityDeploy.deployPath')">
          {{ deployInfo.deployPath }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.shareActivityDeploy.lastUploadAt')">
          {{ deployInfo.lastUploadAt || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="120px">
        <div class="section-title">{{ t('pages.shareActivityDeploy.siteMappingTitle') }}</div>
        <div class="form-tip">{{ t('pages.shareActivityDeploy.siteMappingTip') }}</div>

        <el-form-item :label="t('pages.shareActivityDeploy.domain')" prop="domain">
          <el-input
              v-model="formData.domain"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.shareActivityDeploy.domainPlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('pages.shareActivityDeploy.deployPath')" prop="deployPath">
          <el-input
              v-model="formData.deployPath"
              clearable
              :disabled="!can('save')"
              :placeholder="t('pages.shareActivityDeploy.deployPathPlaceholder')"
          />
        </el-form-item>

		<el-form-item v-if="metaInfo.updatedAt" :label="t('pages.shareActivityDeploy.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item v-if="can('save')">
          <el-button :loading="saving" type="primary" @click="handleSaveCfg">
            {{ t('pages.shareActivityDeploy.saveAndRefresh') }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="upload-section">
        <div class="section-title">{{ t('pages.shareActivityDeploy.uploadTitle') }}</div>
        <div class="form-tip">{{ t('pages.shareActivityDeploy.uploadTip') }}</div>

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
            {{ t('pages.shareActivityDeploy.dragTip') }}
          </div>
          <template #tip>
            <div class="form-tip">{{ t('pages.shareActivityDeploy.fileTip') }}</div>
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
            {{ t('pages.shareActivityDeploy.deployBtn') }}
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
            <div>{{ t('pages.shareActivityDeploy.resultPath', {path: deployResult.deployPath}) }}</div>
            <div>{{ t('pages.shareActivityDeploy.resultFiles', {count: deployResult.fileCount}) }}</div>
            <div>{{ t('pages.shareActivityDeploy.resultDirs', {count: deployResult.dirCount}) }}</div>
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
import {shareActivityDeployApi} from '@/api/modules/share-activity-deploy'
import type {DeployShareActivityZipRes, ShareActivityDeployInfo} from '@/api/modules/share-activity-deploy'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('ShareActivityDeployManagement')

const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')
const deployInfo = ref<ShareActivityDeployInfo | null>(null)
const selectedFile = ref<File | null>(null)
const deployResult = ref<DeployShareActivityZipRes | null>(null)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: 0,
  domain: '',
  deployPath: '/home/ec2-user/cdn/share-activity',
})

const metaInfo = reactive({
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  domain: [{required: true, message: t('pages.shareActivityDeploy.domainRequired'), trigger: 'blur'}],
  deployPath: [{required: true, message: t('pages.shareActivityDeploy.deployPathRequired'), trigger: 'blur'}],
}))

const deployResultTitle = computed(() => t('pages.shareActivityDeploy.deploySuccess'))

const applyDeployInfo = (info: ShareActivityDeployInfo | null | undefined) => {
  deployInfo.value = info ?? null
  formData.id = info?.id ? Number(info.id) : 0
  formData.domain = info?.domain ?? ''
  formData.deployPath = info?.deployPath ?? '/home/ec2-user/cdn/share-activity'
  metaInfo.updatedAt = info?.updatedAt ?? ''
}

const fetchDeployInfo = async () => {
  loading.value = true
  try {
    const response = await shareActivityDeployApi.getShareActivityDeployInfo()
    applyDeployInfo(response.info)
  } catch {
    ElMessage.error(t('pages.shareActivityDeploy.fetchInfoFailed'))
  } finally {
    loading.value = false
  }
}

const handleSaveCfg = async () => {
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
    await shareActivityDeployApi.saveShareActivityDeployCfg({
      id: formData.id,
      domain: formData.domain.trim(),
      deployPath: formData.deployPath.trim(),
    })
    ElMessage.success(t('pages.shareActivityDeploy.saveSuccess'))
    await fetchDeployInfo()
  } catch {
    ElMessage.error(t('pages.shareActivityDeploy.saveFailed'))
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
  ElMessage.warning(t('pages.shareActivityDeploy.singleFileOnly'))
}

const handleDeploy = async () => {
  if (!selectedFile.value) {
    return
  }
  const fileName = selectedFile.value.name.toLowerCase()
  if (!fileName.endsWith('.zip')) {
    ElMessage.error(t('pages.shareActivityDeploy.zipOnly'))
    return
  }

  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  deployResult.value = null
  try {
    const response = await shareActivityDeployApi.deployZip(selectedFile.value, (percent) => {
      uploadPercent.value = percent
    })
    deployResult.value = response
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.shareActivityDeploy.deploySuccess'))
    await fetchDeployInfo()
  } catch {
    uploadStatus.value = 'exception'
    ElMessage.error(t('pages.shareActivityDeploy.deployFailed'))
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

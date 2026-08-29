<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.UploadResourceCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-divider content-position="left">{{ t('pages.uploadResource.staticResources') }}</el-divider>

        <el-form-item :label="t('pages.uploadResource.resourceDomain')" prop="resourceDomain">
          <el-input
              v-model="formData.resourceDomain"
              clearable
              :placeholder="t('pages.uploadResource.resourceDomainPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.uploadResource.resourceDomainTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.uploadResource.storagePath')" prop="storagePath">
          <el-input
              v-model="formData.storagePath"
              clearable
              :placeholder="t('pages.uploadResource.storagePathPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.uploadResource.storagePathTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.uploadResource.cmsExportTtlMinutes')" prop="cmsExportTtlMinutes">
          <el-input-number v-model="formData.cmsExportTtlMinutes" :max="10080" :min="0" :step="1" controls-position="right"/>
          <span class="form-tip">{{ t('pages.uploadResource.cmsExportTtlMinutesTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.uploadResource.appImageMaxSizeMB')" prop="appImageMaxSizeMB">
          <el-input-number v-model="formData.appImageMaxSizeMB" :max="1024" :min="1" :step="1"/>
          <span class="form-tip">{{ t('pages.uploadResource.appImageMaxSizeTip') }}</span>
        </el-form-item>

        <el-divider content-position="left">{{ t('pages.uploadResource.imageModerationSection') }}</el-divider>

        <el-form-item :label="t('pages.uploadResource.enableImageModeration')">
          <el-switch
              v-model="formData.imageModerationEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <span class="form-tip">{{ t('pages.uploadResource.enableImageModerationTip') }}</span>
        </el-form-item>

        <template v-if="formData.imageModerationEnabled">
          <el-form-item label="AccessKey ID" prop="imageModerationAccessKeyId">
            <el-input v-model="formData.imageModerationAccessKeyId" clearable :placeholder="t('pages.uploadResource.accessKeyIdPlaceholder')"/>
          </el-form-item>

          <el-form-item label="AccessKey Secret" prop="imageModerationAccessKeySecret">
            <el-input
                v-model="formData.imageModerationAccessKeySecret"
                clearable
                :placeholder="t('pages.uploadResource.accessKeySecretPlaceholder')"
                show-password
                type="password"
            />
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.regionId')" prop="imageModerationRegionId">
            <el-input v-model="formData.imageModerationRegionId" clearable :placeholder="t('pages.uploadResource.imageModerationRegionPlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.endpoint')" prop="imageModerationEndpoint">
            <el-input
                v-model="formData.imageModerationEndpoint"
                clearable
                :placeholder="t('pages.uploadResource.imageModerationEndpointPlaceholder')"
            />
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.imageModerationService')" prop="imageModerationService">
            <el-input
                v-model="formData.imageModerationService"
                clearable
                :placeholder="t('pages.uploadResource.imageModerationServicePlaceholder')"
            />
            <span class="form-tip">{{ t('pages.uploadResource.imageModerationServiceTip') }}</span>
          </el-form-item>
        </template>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.uploadResource.lastUpdated')">
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {uploadResourceApi} from '@/api/modules/upload-resource'
import type {UploadResourceCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()
const imageSecretTouched = ref(false)

const DEFAULT_STORAGE_PATH = '/home/ec2-user/cdn/images'

/** Linux(/...) 或 Windows(D:\...) 绝对路径 */
function isAbsoluteStoragePath(path: string | undefined): boolean {
  if (!path) {
    return false
  }
  if (path.startsWith('/')) {
    return true
  }
  if (/^[A-Za-z]:[/\\]/.test(path)) {
    return true
  }
  return path.startsWith('\\\\') || path.startsWith('//')
}

const formData = reactive({
  id: '0',
  resourceDomain: 'http://127.0.0.1',
  storagePath: DEFAULT_STORAGE_PATH,
  cmsExportTtlMinutes: 30,
  appImageMaxSizeMB: 1,
  imageModerationEnabled: false,
  imageModerationAccessKeyId: '',
  imageModerationAccessKeySecret: '',
  imageModerationRegionId: 'cn-shanghai',
  imageModerationEndpoint: 'green-cip.cn-shanghai.aliyuncs.com',
  imageModerationService: 'profilePhotoCheck',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  resourceDomain: [{max: 256, message: t('pages.uploadResource.domainMaxLength'), trigger: 'blur'}],
  storagePath: [
    {required: true, message: t('pages.uploadResource.storagePathRequired'), trigger: 'blur'},
    {max: 512, message: t('pages.uploadResource.storagePathMaxLength'), trigger: 'blur'},
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!isAbsoluteStoragePath(value?.trim())) {
          callback(new Error(t('pages.uploadResource.storagePathAbsolute')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  cmsExportTtlMinutes: [
    {required: true, message: t('pages.uploadResource.cmsExportTtlRequired'), trigger: 'change'},
    {type: 'number', min: 0, max: 10080, message: t('pages.uploadResource.cmsExportTtlRange'), trigger: 'change'},
  ],
  appImageMaxSizeMB: [
    {required: true, message: t('pages.uploadResource.appImageMaxSizeRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.uploadResource.appImageMaxSizeMin'), trigger: 'blur'},
  ],
  imageModerationAccessKeyId: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.imageModerationEnabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.uploadResource.imageModerationAccessKeyRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
}))

watch(
    () => formData.imageModerationAccessKeySecret,
    () => {
      imageSecretTouched.value = true
    },
)

const applyCfg = (cfg: UploadResourceCfg | null | undefined) => {
  imageSecretTouched.value = false
  if (!cfg) {
    formData.id = '0'
    formData.resourceDomain = 'http://127.0.0.1'
    formData.storagePath = DEFAULT_STORAGE_PATH
    formData.cmsExportTtlMinutes = 30
    formData.appImageMaxSizeMB = 1
    formData.imageModerationEnabled = false
    formData.imageModerationAccessKeyId = ''
    formData.imageModerationAccessKeySecret = ''
    formData.imageModerationRegionId = 'cn-shanghai'
    formData.imageModerationEndpoint = 'green-cip.cn-shanghai.aliyuncs.com'
    formData.imageModerationService = 'profilePhotoCheck'
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.resourceDomain = cfg.resourceDomain || 'http://127.0.0.1'
  formData.storagePath = cfg.storagePath || DEFAULT_STORAGE_PATH
  formData.cmsExportTtlMinutes = cfg.cmsExportTtlMinutes ?? 30
  formData.appImageMaxSizeMB = cfg.appImageMaxSizeMB || 1
  formData.imageModerationEnabled = !!cfg.imageModerationEnabled
  formData.imageModerationAccessKeyId = cfg.imageModerationAccessKeyId || ''
  formData.imageModerationAccessKeySecret = ''
  formData.imageModerationRegionId = cfg.imageModerationRegionId || 'cn-shanghai'
  formData.imageModerationEndpoint = cfg.imageModerationEndpoint || 'green-cip.cn-shanghai.aliyuncs.com'
  formData.imageModerationService = cfg.imageModerationService || 'profilePhotoCheck'
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await uploadResourceApi.getUploadResourceCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch upload resource cfg failed:', error)
    ElMessage.error(t('pages.uploadResource.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await uploadResourceApi.saveUploadResourceCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      resourceDomain: formData.resourceDomain.trim(),
      storagePath: formData.storagePath.trim(),
      cmsExportTtlMinutes: formData.cmsExportTtlMinutes,
      appImageMaxSizeMB: formData.appImageMaxSizeMB,
      imageModerationEnabled: formData.imageModerationEnabled,
      imageModerationAccessKeyId: formData.imageModerationAccessKeyId.trim(),
      imageModerationAccessKeySecret: imageSecretTouched.value
          ? formData.imageModerationAccessKeySecret.trim()
          : '',
      imageModerationRegionId: formData.imageModerationRegionId.trim(),
      imageModerationEndpoint: formData.imageModerationEndpoint.trim(),
      imageModerationService: formData.imageModerationService.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.uploadResource.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.uploadResource.saveFailed'))
    }
  } catch (error) {
    console.error('save upload resource cfg failed:', error)
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

.cfg-form {
  max-width: 760px;
}

.form-tip {
  display: block;
  margin-top: 6px;
  margin-left: 0;
  color: #909399;
  font-size: 13px;
  line-height: 1.4;
}
</style>

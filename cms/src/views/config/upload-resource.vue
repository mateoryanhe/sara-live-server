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

        <el-divider content-position="left">{{ t('pages.uploadResource.s3Section') }}</el-divider>

        <el-form-item :label="t('pages.uploadResource.enableS3')">
          <div class="switch-block">
            <el-switch
                v-model="formData.s3Enabled"
                :active-value="true"
                :inactive-value="false"
                :active-text="t('common.open')"
                :inactive-text="t('common.close')"
                @change="onS3EnabledChange"
            />
            <div class="form-tip">{{ t('pages.uploadResource.enableS3Tip') }}</div>
          </div>
        </el-form-item>

        <template v-if="formData.s3Enabled === true">
          <el-form-item :label="t('pages.uploadResource.s3PublicDomain')" prop="s3PublicDomain">
            <el-input
                v-model="formData.s3PublicDomain"
                clearable
                :placeholder="t('pages.uploadResource.s3PublicDomainPlaceholder')"
            />
            <div class="form-tip">{{ t('pages.uploadResource.s3PublicDomainTip') }}</div>
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.s3Endpoint')" prop="s3Endpoint">
            <el-input
                v-model="formData.s3Endpoint"
                clearable
                :placeholder="t('pages.uploadResource.s3EndpointPlaceholder')"
            />
            <div class="form-tip">{{ t('pages.uploadResource.s3EndpointTip') }}</div>
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.s3Bucket')" prop="s3Bucket">
            <el-input v-model="formData.s3Bucket" clearable :placeholder="t('pages.uploadResource.s3BucketPlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.uploadResource.s3KeyPrefix')" prop="s3KeyPrefix">
            <el-input
                v-model="formData.s3KeyPrefix"
                clearable
                :placeholder="t('pages.uploadResource.s3KeyPrefixPlaceholder')"
            />
            <div class="form-tip">{{ t('pages.uploadResource.s3KeyPrefixTip') }}</div>
          </el-form-item>

          <el-form-item label="AccessKey ID" prop="s3AccessKeyId">
            <el-input v-model="formData.s3AccessKeyId" clearable :placeholder="t('pages.uploadResource.s3AccessKeyIdPlaceholder')"/>
          </el-form-item>

          <el-form-item label="Secret Access Key" prop="s3SecretAccessKey">
            <el-input
                v-model="formData.s3SecretAccessKey"
                clearable
                :placeholder="t('pages.uploadResource.accessKeySecretPlaceholder')"
                show-password
                type="password"
            />
          </el-form-item>
        </template>

        <el-divider content-position="left">{{ t('pages.uploadResource.imageModerationSection') }}</el-divider>

        <el-form-item :label="t('pages.uploadResource.enableImageModeration')">
          <div class="switch-block">
            <el-switch
                v-model="formData.imageModerationEnabled"
                :active-value="true"
                :inactive-value="false"
                :active-text="t('common.open')"
                :inactive-text="t('common.close')"
            />
            <div class="form-tip">{{ t('pages.uploadResource.enableImageModerationTip') }}</div>
          </div>
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
          <el-button
              :disabled="syncing || !formData.s3Enabled"
              :loading="syncing"
              type="warning"
              @click="handleSyncLocalToS3"
          >
            {{ t('pages.uploadResource.syncLocalToS3') }}
          </el-button>
          <div class="form-tip">{{ t('pages.uploadResource.syncLocalToS3Tip') }}</div>
        </el-form-item>

        <el-form-item v-if="syncStatus.startedAt || syncStatus.running" :label="t('pages.uploadResource.syncLocalStatus')">
          <div class="sync-status">
            <div>
              <el-tag :type="syncStatus.running ? 'warning' : 'success'" size="small">
                {{ syncStatus.running ? t('pages.uploadResource.syncLocalRunning') : t('pages.uploadResource.syncLocalIdle') }}
              </el-tag>
              <span class="sync-meta">
                {{ syncStatus.root || formData.storagePath }}
                <template v-if="syncStatus.keyPrefix"> · prefix={{ syncStatus.keyPrefix }}</template>
              </span>
            </div>
            <div class="sync-line">
              {{
                t('pages.uploadResource.syncLocalProgress', {
                  done: syncStatus.done,
                  total: syncStatus.total,
                  success: syncStatus.success,
                  failed: syncStatus.failed,
                  skipped: syncStatus.skipped,
                })
              }}
            </div>
            <div v-if="syncStatus.lastKey" class="sync-line muted">
              {{ t('pages.uploadResource.syncLocalLastKey') }}: {{ syncStatus.lastKey }}
            </div>
            <div v-if="syncStatus.lastError" class="sync-line error">
              {{ t('pages.uploadResource.syncLocalLastError') }}: {{ syncStatus.lastError }}
            </div>
            <div v-if="syncStatus.startedAt || syncStatus.finishedAt" class="sync-line muted">
              <template v-if="syncStatus.startedAt">start {{ syncStatus.startedAt }}</template>
              <template v-if="syncStatus.finishedAt"> · end {{ syncStatus.finishedAt }}</template>
            </div>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {uploadResourceApi} from '@/api/modules/upload-resource'
import type {SyncLocalToS3Status, UploadResourceCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const syncing = ref(false)
const formRef = ref()
const imageSecretTouched = ref(false)
const s3SecretTouched = ref(false)
const applyingCfg = ref(false)
let syncPollTimer: ReturnType<typeof setInterval> | null = null

const syncStatus = reactive<SyncLocalToS3Status>({
  running: false,
  root: '',
  keyPrefix: '',
  startedAt: '',
  finishedAt: '',
  total: 0,
  done: 0,
  success: 0,
  failed: 0,
  skipped: 0,
  lastKey: '',
  lastError: '',
})

const applySyncStatus = (s: Partial<SyncLocalToS3Status> | null | undefined) => {
  if (!s) {
    return
  }
  syncStatus.running = s.running === true
  syncStatus.root = s.root || ''
  syncStatus.keyPrefix = s.keyPrefix || ''
  syncStatus.startedAt = s.startedAt || ''
  syncStatus.finishedAt = s.finishedAt || ''
  syncStatus.total = s.total ?? 0
  syncStatus.done = s.done ?? 0
  syncStatus.success = s.success ?? 0
  syncStatus.failed = s.failed ?? 0
  syncStatus.skipped = s.skipped ?? 0
  syncStatus.lastKey = s.lastKey || ''
  syncStatus.lastError = s.lastError || ''
}

const stopSyncPoll = () => {
  if (syncPollTimer) {
    clearInterval(syncPollTimer)
    syncPollTimer = null
  }
}

const pollSyncStatusOnce = async () => {
  try {
    const status = await uploadResourceApi.getSyncLocalStorageToS3Status()
    applySyncStatus(status)
    if (!status?.running) {
      stopSyncPoll()
      syncing.value = false
    }
  } catch (error) {
    console.error('poll sync local to s3 status failed:', error)
  }
}

const startSyncPoll = () => {
  stopSyncPoll()
  syncPollTimer = setInterval(() => {
    void pollSyncStatusOnce()
  }, 2000)
  void pollSyncStatusOnce()
}

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
  s3Enabled: false,
  s3PublicDomain: '',
  s3Endpoint: '',
  s3Bucket: '',
  s3AccessKeyId: '',
  s3SecretAccessKey: '',
  s3KeyPrefix: '',
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
  s3PublicDomain: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.s3Enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.uploadResource.s3PublicDomainRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  s3Endpoint: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.s3Enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.uploadResource.s3EndpointRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  s3Bucket: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.s3Enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.uploadResource.s3BucketRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  s3AccessKeyId: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.s3Enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.uploadResource.s3AccessKeyRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
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
      if (!applyingCfg.value) {
        imageSecretTouched.value = true
      }
    },
)

watch(
    () => formData.s3SecretAccessKey,
    () => {
      if (!applyingCfg.value) {
        s3SecretTouched.value = true
      }
    },
)

const onS3EnabledChange = async (val: string | number | boolean) => {
  formData.s3Enabled = val === true
  await nextTick()
  formRef.value?.clearValidate?.(['s3PublicDomain', 's3Endpoint', 's3Bucket', 's3AccessKeyId', 's3SecretAccessKey'])
}

const resetFormDefaults = () => {
  formData.id = '0'
  formData.resourceDomain = 'http://127.0.0.1'
  formData.storagePath = DEFAULT_STORAGE_PATH
  formData.cmsExportTtlMinutes = 30
  formData.appImageMaxSizeMB = 1
  formData.s3Enabled = false
  formData.s3PublicDomain = ''
  formData.s3Endpoint = ''
  formData.s3Bucket = ''
  formData.s3AccessKeyId = ''
  formData.s3SecretAccessKey = ''
  formData.s3KeyPrefix = ''
  formData.imageModerationEnabled = false
  formData.imageModerationAccessKeyId = ''
  formData.imageModerationAccessKeySecret = ''
  formData.imageModerationRegionId = 'cn-shanghai'
  formData.imageModerationEndpoint = 'green-cip.cn-shanghai.aliyuncs.com'
  formData.imageModerationService = 'profilePhotoCheck'
  metaInfo.createdAt = ''
  metaInfo.updatedAt = ''
}

const applyCfg = (cfg: UploadResourceCfg | null | undefined) => {
  applyingCfg.value = true
  imageSecretTouched.value = false
  s3SecretTouched.value = false
  try {
    if (!cfg) {
      resetFormDefaults()
      return
    }
    formData.id = cfg.id || '0'
    formData.resourceDomain = cfg.resourceDomain || 'http://127.0.0.1'
    formData.storagePath = cfg.storagePath || DEFAULT_STORAGE_PATH
    formData.cmsExportTtlMinutes = cfg.cmsExportTtlMinutes ?? 30
    formData.appImageMaxSizeMB = cfg.appImageMaxSizeMB || 1
    formData.s3Enabled = cfg.s3Enabled === true
    formData.s3PublicDomain = cfg.s3PublicDomain || ''
    formData.s3Endpoint = cfg.s3Endpoint || ''
    formData.s3Bucket = cfg.s3Bucket || ''
    formData.s3AccessKeyId = cfg.s3AccessKeyId || ''
    formData.s3SecretAccessKey = ''
    formData.s3KeyPrefix = cfg.s3KeyPrefix || ''
    formData.imageModerationEnabled = cfg.imageModerationEnabled === true
    formData.imageModerationAccessKeyId = cfg.imageModerationAccessKeyId || ''
    formData.imageModerationAccessKeySecret = ''
    formData.imageModerationRegionId = cfg.imageModerationRegionId || 'cn-shanghai'
    formData.imageModerationEndpoint = cfg.imageModerationEndpoint || 'green-cip.cn-shanghai.aliyuncs.com'
    formData.imageModerationService = cfg.imageModerationService || 'profilePhotoCheck'
    metaInfo.createdAt = cfg.createdAt || ''
    metaInfo.updatedAt = cfg.updatedAt || ''
  } finally {
    nextTick(() => {
      applyingCfg.value = false
    })
  }
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
  const wantS3 = formData.s3Enabled === true
  try {
    await formRef.value.validate()
  } catch {
    ElMessage.warning(t('pages.uploadResource.formInvalid'))
    return
  }
  loading.value = true
  try {
    const response = await uploadResourceApi.saveUploadResourceCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      resourceDomain: formData.resourceDomain.trim(),
      storagePath: formData.storagePath.trim(),
      cmsExportTtlMinutes: formData.cmsExportTtlMinutes,
      appImageMaxSizeMB: formData.appImageMaxSizeMB,
      s3Enabled: wantS3,
      s3PublicDomain: formData.s3PublicDomain.trim(),
      s3Endpoint: formData.s3Endpoint.trim(),
      s3Bucket: formData.s3Bucket.trim(),
      s3AccessKeyId: formData.s3AccessKeyId.trim(),
      s3SecretAccessKey: s3SecretTouched.value ? formData.s3SecretAccessKey.trim() : '',
      s3KeyPrefix: formData.s3KeyPrefix.trim(),
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
      if (wantS3 && formData.s3Enabled !== true) {
        ElMessage.warning(t('pages.uploadResource.s3NotPersisted'))
      }
    } else {
      ElMessage.error(t('pages.uploadResource.saveFailed'))
    }
  } catch (error) {
    console.error('save upload resource cfg failed:', error)
  } finally {
    loading.value = false
  }
}

const handleSyncLocalToS3 = async () => {
  if (!formData.s3Enabled) {
    ElMessage.warning(t('pages.uploadResource.syncLocalNeedS3'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.uploadResource.syncLocalToS3Confirm', {path: formData.storagePath || DEFAULT_STORAGE_PATH}),
        t('pages.uploadResource.syncLocalToS3'),
        {type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel')},
    )
  } catch {
    return
  }
  syncing.value = true
  try {
    const res = await uploadResourceApi.syncLocalStorageToS3()
    applySyncStatus(res)
    if (res?.alreadyRunning) {
      ElMessage.warning(t('pages.uploadResource.syncLocalAlreadyRunning'))
    } else if (res?.started) {
      ElMessage.success(t('pages.uploadResource.syncLocalStarted'))
    } else {
      ElMessage.warning(res?.message || t('pages.uploadResource.syncLocalFailed'))
    }
    startSyncPoll()
  } catch (error) {
    console.error('sync local to s3 failed:', error)
    syncing.value = false
    ElMessage.error(t('pages.uploadResource.syncLocalFailed'))
  }
}

onMounted(() => {
  fetchCfg()
  void pollSyncStatusOnce().then(() => {
    if (syncStatus.running) {
      syncing.value = true
      startSyncPoll()
    }
  })
})

onBeforeUnmount(() => {
  stopSyncPoll()
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

.switch-block {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.sync-status {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-width: 640px;
}

.sync-meta {
  margin-left: 8px;
  color: #606266;
  font-size: 13px;
}

.sync-line {
  font-size: 13px;
  color: #303133;
  line-height: 1.4;
  word-break: break-all;
}

.sync-line.muted {
  color: #909399;
}

.sync-line.error {
  color: #f56c6c;
}
</style>

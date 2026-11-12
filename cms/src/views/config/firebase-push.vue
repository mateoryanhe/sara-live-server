<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.FirebasePushCfgManagement') }}</span>
          <span v-if="metaInfo.updatedAt" class="header-meta">
            {{ t('pages.firebasePush.lastUpdated') }}: {{ metaInfo.updatedAt }}
          </span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.firebasePush.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.firebasePush.noticeLine1') }}</p>
        <p>{{ t('pages.firebasePush.noticeLine2') }}</p>
        <p>{{ t('pages.firebasePush.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="160px">
        <el-tabs v-model="activeTab" class="config-tabs">
          <el-tab-pane :label="t('pages.firebasePush.tabBasic')" name="basic">
            <div class="tab-pane-body">
              <el-form-item :label="t('pages.firebasePush.enabled')">
                <el-switch v-model="formData.enabled" :active-value="1" :inactive-value="0"/>
                <div class="form-tip">{{ t('pages.firebasePush.enabledTip') }}</div>
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.firebasePush.tabOfflineCall')" name="offlineCall">
            <div class="tab-pane-body">
              <el-form-item :label="t('pages.firebasePush.offlineCallTitle')" prop="offlineCallTitle">
                <el-input
                    v-model="formData.offlineCallTitle"
                    clearable
                    maxlength="128"
                    show-word-limit
                    :placeholder="t('pages.firebasePush.offlineCallTitlePlaceholder')"
                />
              </el-form-item>
              <el-form-item :label="t('pages.firebasePush.offlineCallBody')" prop="offlineCallBody">
                <el-input
                    v-model="formData.offlineCallBody"
                    :rows="3"
                    maxlength="512"
                    show-word-limit
                    type="textarea"
                    :placeholder="t('pages.firebasePush.offlineCallBodyPlaceholder')"
                />
                <span class="form-tip">{{ t('pages.firebasePush.offlineCallBodyTip') }}</span>
              </el-form-item>
              <el-form-item :label="t('pages.firebasePush.offlineCallIcon')" prop="offlineCallIcon">
                <div class="icon-upload-wrap">
                  <el-upload
                      :before-upload="beforeIconUpload"
                      :disabled="iconUploading"
                      :http-request="doIconUpload"
                      :show-file-list="false"
                      accept="image/*"
                      action="#"
                      class="icon-uploader"
                  >
                    <img v-if="formData.offlineCallIcon" :src="formData.offlineCallIcon" alt="icon" class="icon-preview"/>
                    <div v-else class="icon-uploader-placeholder">
                      <el-icon class="icon-uploader-icon">
                        <Plus/>
                      </el-icon>
                      <span>{{ t('pages.firebasePush.clickUploadIcon') }}</span>
                    </div>
                  </el-upload>
                  <el-input
                      v-model="formData.offlineCallIcon"
                      clearable
                      class="icon-url-input"
                      :placeholder="t('pages.firebasePush.offlineCallIconPlaceholder')"
                  />
                  <el-button
                      v-if="formData.offlineCallIcon"
                      link
                      type="danger"
                      @click="formData.offlineCallIcon = ''"
                  >
                    {{ t('pages.firebasePush.removeIcon') }}
                  </el-button>
                </div>
                <span class="form-tip">{{ t('pages.firebasePush.offlineCallIconTip') }}</span>
              </el-form-item>
              <el-form-item>
                <el-button @click="applyDefaults">{{ t('pages.firebasePush.resetDefaults') }}</el-button>
              </el-form-item>
            </div>
          </el-tab-pane>
        </el-tabs>

        <div class="form-actions">
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {firebasePushApi} from '@/api/modules/firebase-push'
import {uploadApi} from '@/api/modules/upload'
import type {FirebasePushCfg} from '@/types/api'

const DEFAULT_OFFLINE_CALL_TITLE = 'Incoming video call'
const DEFAULT_OFFLINE_CALL_BODY = '{nickname} is calling you'

const {t} = useI18n()
const loading = ref(false)
const iconUploading = ref(false)
const formRef = ref<FormInstance>()
const activeTab = ref<'basic' | 'offlineCall'>('basic')

const formData = reactive({
  id: '0',
  enabled: 0 as 0 | 1,
  offlineCallTitle: DEFAULT_OFFLINE_CALL_TITLE,
  offlineCallBody: DEFAULT_OFFLINE_CALL_BODY,
  offlineCallIcon: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const applyDefaults = () => {
  formData.offlineCallTitle = DEFAULT_OFFLINE_CALL_TITLE
  formData.offlineCallBody = DEFAULT_OFFLINE_CALL_BODY
  formData.offlineCallIcon = ''
}

const applyCfg = (cfg: FirebasePushCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = 0
    applyDefaults()
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = cfg.enabled === 1 ? 1 : 0
  formData.offlineCallTitle = (cfg.offlineCallTitle || '').trim() || DEFAULT_OFFLINE_CALL_TITLE
  formData.offlineCallBody = (cfg.offlineCallBody || '').trim() || DEFAULT_OFFLINE_CALL_BODY
  formData.offlineCallIcon = (cfg.offlineCallIcon || '').trim()
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await firebasePushApi.getFirebasePushCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch firebase push cfg failed:', error)
    ElMessage.error(t('pages.firebasePush.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const beforeIconUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error(t('pages.firebasePush.uploadFailed'))
    return false
  }
  return true
}

const doIconUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  iconUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    const url = (res.fileUrl || '').trim()
    if (!url) {
      ElMessage.error(t('pages.firebasePush.uploadFailed'))
      return
    }
    formData.offlineCallIcon = url
  } catch (error) {
    console.error('upload push icon failed:', error)
    ElMessage.error(t('pages.firebasePush.uploadFailed'))
  } finally {
    iconUploading.value = false
  }
}

const handleSave = async () => {
  loading.value = true
  try {
    const response = await firebasePushApi.saveFirebasePushCfg({
      id: Number(formData.id) || 0,
      enabled: formData.enabled,
      offlineCallTitle: formData.offlineCallTitle.trim() || DEFAULT_OFFLINE_CALL_TITLE,
      offlineCallBody: formData.offlineCallBody.trim() || DEFAULT_OFFLINE_CALL_BODY,
      offlineCallIcon: formData.offlineCallIcon.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.firebasePush.saveSuccess'))
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.firebasePush.saveFailed'))
    }
  } catch (error) {
    console.error('save firebase push cfg failed:', error)
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

.header-meta {
  color: var(--el-text-color-secondary);
  font-size: 13px;
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

.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.tab-pane-body {
  padding-top: 4px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.form-actions {
  margin-top: 8px;
}

.icon-upload-wrap {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.icon-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  width: 96px;
  height: 96px;
}

.icon-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.icon-preview {
  width: 96px;
  height: 96px;
  object-fit: cover;
  display: block;
}

.icon-uploader-placeholder {
  width: 96px;
  height: 96px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  gap: 6px;
}

.icon-uploader-icon {
  font-size: 22px;
}

.icon-url-input {
  max-width: 520px;
}
</style>

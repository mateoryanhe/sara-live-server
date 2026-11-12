<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GameActivityCfgManagement') }}</span>
          <span v-if="metaInfo.updatedAt" class="header-meta">
            {{ t('pages.gameActivity.lastUpdated') }}: {{ metaInfo.updatedAt }}
          </span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.gameActivity.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.gameActivity.noticeLine1') }}</p>
        <p>{{ t('pages.gameActivity.noticeLine2') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.gameActivity.buttonIcon')" prop="buttonIcon">
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
              <img v-if="formData.buttonIcon" :src="formData.buttonIcon" alt="icon" class="icon-preview"/>
              <div v-else class="icon-uploader-placeholder">
                <el-icon class="icon-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>{{ t('pages.gameActivity.clickUploadIcon') }}</span>
              </div>
            </el-upload>
            <el-input
                v-model="formData.buttonIcon"
                clearable
                class="icon-url-input"
                :placeholder="t('pages.gameActivity.buttonIconPlaceholder')"
            />
            <el-button
                v-if="formData.buttonIcon"
                link
                type="danger"
                @click="formData.buttonIcon = ''"
            >
              {{ t('pages.gameActivity.removeIcon') }}
            </el-button>
          </div>
          <span class="form-tip">{{ t('pages.gameActivity.buttonIconTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.gameActivity.gameListUrl')" prop="gameListUrl">
          <el-input
              v-model="formData.gameListUrl"
              clearable
              maxlength="1024"
              show-word-limit
              :placeholder="t('pages.gameActivity.gameListUrlPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.gameActivity.gameListUrlTip') }}</span>
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
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {gameActivityApi} from '@/api/modules/game-activity'
import {uploadApi} from '@/api/modules/upload'
import type {GameActivityCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const iconUploading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  buttonIcon: '',
  gameListUrl: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const applyCfg = (cfg: GameActivityCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.buttonIcon = ''
    formData.gameListUrl = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.buttonIcon = (cfg.buttonIcon || '').trim()
  formData.gameListUrl = (cfg.gameListUrl || '').trim()
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await gameActivityApi.getGameActivityCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch game activity cfg failed:', error)
    ElMessage.error(t('pages.gameActivity.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const beforeIconUpload = (file: File) => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.gameActivity.uploadFailed'))
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
      ElMessage.error(t('pages.gameActivity.uploadFailed'))
      return
    }
    formData.buttonIcon = url
  } catch (error) {
    console.error('upload game activity icon failed:', error)
    ElMessage.error(t('pages.gameActivity.uploadFailed'))
  } finally {
    iconUploading.value = false
  }
}

const handleSave = async () => {
  loading.value = true
  try {
    const response = await gameActivityApi.saveGameActivityCfg({
      id: Number(formData.id) || 0,
      buttonIcon: formData.buttonIcon.trim(),
      gameListUrl: formData.gameListUrl.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.gameActivity.saveSuccess'))
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.gameActivity.saveFailed'))
    }
  } catch (error) {
    console.error('save game activity cfg failed:', error)
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

.form-tip {
  display: block;
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
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

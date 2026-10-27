<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionShareAppDownloadCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.shareAppDownload.downloadUrl')">
          <el-input
              :model-value="formData.downloadUrl || t('pages.shareAppDownload.downloadUrlEmpty')"
              readonly
          />
          <div class="form-tip">{{ t('pages.shareAppDownload.downloadUrlTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.shareAppDownload.shortUrl')">
          <el-input
              :model-value="formData.shortUrl || t('pages.shareAppDownload.shortUrlEmpty')"
              readonly
          />
          <div class="form-tip">{{ t('pages.shareAppDownload.shortUrlTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.shareAppDownload.shareText')" prop="shareText">
          <el-input
              v-model="formData.shareText"
              type="textarea"
              :rows="8"
              maxlength="4000"
              show-word-limit
              :placeholder="t('pages.shareAppDownload.shareTextPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.shareAppDownload.shareTextTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.shareAppDownload.previewShareCode')">
          <div class="preview-wrap">
            <el-input
                :model-value="previewText || t('pages.shareAppDownload.previewShareCodeEmpty')"
                type="textarea"
                :rows="6"
                readonly
            />
            <el-button
                class="copy-btn"
                type="primary"
                plain
                :disabled="!previewText"
                @click="handleCopyPreview"
            >
              {{ t('pages.shareAppDownload.copyPreview') }}
            </el-button>
          </div>
          <div class="form-tip">{{ t('pages.shareAppDownload.previewShareCodeTip') }}</div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.shareAppDownload.lastUpdated')">
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
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {shareAppDownloadApi} from '@/api/modules/share-app-download'
import type {ShareAppDownloadCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  downloadUrl: '',
  shortUrl: '',
  shareText: '',
  previewShareCode: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  shareText: [
    {max: 4000, message: t('pages.shareAppDownload.shareTextMaxLength'), trigger: 'blur'},
  ],
}))

/** 本地即时预览(与服务端规则一致);保存后仍以接口 previewShareCode 为准刷新 */
const localPreview = computed(() => {
  const text = formData.shareText.trim().replaceAll('{inviteCode}', '').replaceAll('{downloadUrl}', '').trim()
  const parts: string[] = []
  if (text) {
    parts.push(text.replace(/\r?\n+$/, ''))
  }
  parts.push('Invite Code: YOUR_INVITE_CODE')
  if (formData.shortUrl) {
    parts.push(formData.shortUrl)
  }
  return parts.join('\n')
})

const previewText = computed(() => formData.previewShareCode || localPreview.value)

const applyCfg = (cfg: ShareAppDownloadCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.downloadUrl = ''
    formData.shortUrl = ''
    formData.shareText = ''
    formData.previewShareCode = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.downloadUrl = cfg.downloadUrl || ''
  formData.shortUrl = cfg.shortUrl || ''
  formData.shareText = cfg.shareText || ''
  formData.previewShareCode = cfg.previewShareCode || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await shareAppDownloadApi.getShareAppDownloadCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch share app download cfg failed:', error)
    ElMessage.error(t('pages.shareAppDownload.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await shareAppDownloadApi.saveShareAppDownloadCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      shareText: formData.shareText.trim(),
    })
    if (response?.id) {
      formData.id = response.id
    }
    ElMessage.success(t('pages.shareAppDownload.saveSuccess'))
    await fetchCfg()
  } catch (error) {
    if (error !== false) {
      console.error('save share app download cfg failed:', error)
      ElMessage.error(t('pages.shareAppDownload.saveCfgFailed'))
    }
  } finally {
    loading.value = false
  }
}

const handleCopyPreview = async () => {
  const text = previewText.value
  if (!text) {
    return
  }
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('pages.shareAppDownload.copyPreviewSuccess'))
  } catch (error) {
    console.error('copy preview failed:', error)
    ElMessage.error(t('pages.shareAppDownload.copyPreviewFailed'))
  }
}

onMounted(() => {
  void fetchCfg()
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
  max-width: 780px;
}

.form-tip {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
}

.preview-wrap {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.copy-btn {
  align-self: flex-start;
}
</style>

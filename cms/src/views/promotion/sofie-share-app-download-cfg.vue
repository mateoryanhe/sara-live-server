<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionSofieShareAppDownloadCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.sofieShareAppDownload.packageName')" prop="packageName">
          <el-input
              v-model="formData.packageName"
              clearable
              :placeholder="t('pages.sofieShareAppDownload.packageNamePlaceholder')"
          />
          <div class="form-tip">{{ t('pages.sofieShareAppDownload.packageNameTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.sofieShareAppDownload.googlePlayUrl')" prop="googlePlayUrl">
          <el-input
              v-model="formData.googlePlayUrl"
              clearable
              :placeholder="t('pages.sofieShareAppDownload.googlePlayUrlPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.sofieShareAppDownload.googlePlayUrlTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.sofieShareAppDownload.shortUrl')">
          <el-input
              :model-value="formData.shortUrl || t('pages.sofieShareAppDownload.shortUrlEmpty')"
              readonly
          />
          <div class="form-tip">{{ t('pages.sofieShareAppDownload.shortUrlTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.sofieShareAppDownload.shareText')" prop="shareText">
          <el-input
              v-model="formData.shareText"
              type="textarea"
              :rows="8"
              maxlength="4000"
              show-word-limit
              :placeholder="t('pages.sofieShareAppDownload.shareTextPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.sofieShareAppDownload.shareTextTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.sofieShareAppDownload.previewShareCode')">
          <div class="preview-wrap">
            <el-input
                :model-value="previewText || t('pages.sofieShareAppDownload.previewShareCodeEmpty')"
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
              {{ t('pages.sofieShareAppDownload.copyPreview') }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.sofieShareAppDownload.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="can('save')" type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
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
import {sofieShareAppDownloadApi} from '@/api/modules/sofie-share-app-download'
import type {SofieShareAppDownloadCfgItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('PromotionSofieShareAppDownloadCfgManagement')
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  packageName: '',
  googlePlayUrl: '',
  shortUrl: '',
  shareText: '',
  previewShareCode: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  googlePlayUrl: [
    {required: true, message: t('pages.sofieShareAppDownload.googlePlayRequired'), trigger: 'blur'},
  ],
  shareText: [
    {max: 4000, message: t('pages.sofieShareAppDownload.shareTextMaxLength'), trigger: 'blur'},
  ],
}))

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

const applyCfg = (cfg: SofieShareAppDownloadCfgItem | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.packageName = ''
    formData.googlePlayUrl = ''
    formData.shortUrl = ''
    formData.shareText = ''
    formData.previewShareCode = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.packageName = cfg.packageName || ''
  formData.googlePlayUrl = cfg.googlePlayUrl || ''
  formData.shortUrl = cfg.shortUrl || ''
  formData.shareText = cfg.shareText || ''
  formData.previewShareCode = cfg.previewShareCode || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await sofieShareAppDownloadApi.getSofieShareAppDownloadCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch sofie share cfg failed:', error)
    ElMessage.error(t('pages.sofieShareAppDownload.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true
    const response = await sofieShareAppDownloadApi.saveSofieShareAppDownloadCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      packageName: formData.packageName.trim(),
      googlePlayUrl: formData.googlePlayUrl.trim(),
      shareText: formData.shareText.trim(),
    })
    if (response?.id) {
      formData.id = response.id
    }
    ElMessage.success(t('pages.sofieShareAppDownload.saveSuccess'))
    await fetchCfg()
  } catch (error) {
    if (error !== false) {
      console.error('save sofie share cfg failed:', error)
      ElMessage.error(t('pages.sofieShareAppDownload.saveFailed'))
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
    ElMessage.success(t('pages.sofieShareAppDownload.copyPreviewSuccess'))
  } catch (error) {
    ElMessage.error(t('pages.sofieShareAppDownload.copyPreviewFailed'))
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
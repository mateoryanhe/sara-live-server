<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionShareAppDownloadCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.shareAppDownload.downloadUrl')">
          <el-input
              :model-value="formData.downloadUrl || t('pages.shareAppDownload.downloadUrlEmpty')"
              readonly
          />
          <div class="form-tip">{{ t('pages.shareAppDownload.downloadUrlTip') }}</div>
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
  shareText: '',
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

const applyCfg = (cfg: ShareAppDownloadCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.downloadUrl = ''
    formData.shareText = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.downloadUrl = cfg.downloadUrl || ''
  formData.shareText = cfg.shareText || ''
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
  max-width: 720px;
}

.form-tip {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
}
</style>

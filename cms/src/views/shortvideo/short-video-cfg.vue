<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ShortVideoCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.shortVideoCfg.maxFileSize')" prop="maxFileSizeMB">
          <el-input-number
              v-model="formData.maxFileSizeMB"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.maxFileSizeTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.maxCoverFileSize')" prop="maxCoverFileSize">
          <el-input-number
              v-model="formData.maxCoverFileSize"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.maxCoverFileSizeTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.maxDuration')" prop="maxDuration">
          <el-input-number
              v-model="formData.maxDuration"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.seconds') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.freeWatchSeconds')" prop="freeWatchSeconds">
          <el-input-number
              v-model="formData.freeWatchSeconds"
              :min="0"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.freeWatchSecondsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.anchorDailyUploadLimit')" prop="anchorDailyUploadLimit">
          <el-input-number
              v-model="formData.anchorDailyUploadLimit"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.anchorDailyUploadLimitTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.normalUserDailyUploadLimit')" prop="normalUserDailyUploadLimit">
          <el-input-number
              v-model="formData.normalUserDailyUploadLimit"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoCfg.normalUserDailyUploadLimitTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.shortVideoCfg.entryEnabled')" prop="entryEnabled">
          <el-switch
              v-model="formData.entryEnabled"
              :active-value="1"
              :inactive-value="0"
              :active-text="t('pages.shortVideoCfg.switchOn')"
              :inactive-text="t('pages.shortVideoCfg.switchOff')"
          />
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.shortVideoCfg.lastUpdated')">
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
import {shortVideoApi} from '@/api/modules/shortVideo'
import type {ShortVideoCfg} from '@/types/api'

const {t} = useI18n()
const MB = 1024 * 1024

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  maxFileSizeMB: 100,
  maxCoverFileSize: 5,
  maxDuration: 60,
  freeWatchSeconds: 8,
  anchorDailyUploadLimit: 100,
  normalUserDailyUploadLimit: 1,
  entryEnabled: 1,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  maxFileSizeMB: [
    {required: true, message: t('pages.shortVideoCfg.maxFileSizeRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.shortVideoCfg.maxFileSizeMin'), trigger: 'blur'},
  ],
  maxCoverFileSize: [
    {required: true, message: t('pages.shortVideoCfg.maxCoverFileSizeRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.shortVideoCfg.maxCoverFileSizeMin'), trigger: 'blur'},
  ],
  maxDuration: [
    {required: true, message: t('pages.shortVideoCfg.maxDurationRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.shortVideoCfg.maxDurationMin'), trigger: 'blur'},
  ],
  freeWatchSeconds: [
    {type: 'number', min: 0, message: t('pages.shortVideoCfg.freeWatchSecondsMin'), trigger: 'blur'},
  ],
  anchorDailyUploadLimit: [
    {required: true, message: t('pages.shortVideoCfg.anchorDailyUploadLimitRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.shortVideoCfg.anchorDailyUploadLimitMin'), trigger: 'blur'},
  ],
  normalUserDailyUploadLimit: [
    {required: true, message: t('pages.shortVideoCfg.normalUserDailyUploadLimitRequired'), trigger: 'blur'},
    {type: 'number', min: 1, message: t('pages.shortVideoCfg.normalUserDailyUploadLimitMin'), trigger: 'blur'},
  ],
  entryEnabled: [
    {required: true, message: t('pages.shortVideoCfg.entryEnabledRequired'), trigger: 'change'},
  ],
}))

const applyCfg = (cfg: ShortVideoCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.maxFileSizeMB = 100
    formData.maxCoverFileSize = 5
    formData.maxDuration = 60
    formData.freeWatchSeconds = 8
    formData.anchorDailyUploadLimit = 100
    formData.normalUserDailyUploadLimit = 1
    formData.entryEnabled = 1
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.maxFileSizeMB = Math.max(1, Math.round(cfg.maxFileSize / MB))
  formData.maxCoverFileSize = cfg.maxCoverFileSize ?? 5
  formData.maxDuration = cfg.maxDuration || 60
  formData.freeWatchSeconds = cfg.freeWatchSeconds ?? 8
  formData.anchorDailyUploadLimit = cfg.anchorDailyUploadLimit ?? 100
  formData.normalUserDailyUploadLimit = cfg.normalUserDailyUploadLimit ?? 1
  formData.entryEnabled = cfg.entryEnabled ?? 1
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch short video cfg failed:', error)
    ElMessage.error(t('pages.shortVideoCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await shortVideoApi.saveShortVideoCfg({
      id: formData.id,
      maxFileSize: formData.maxFileSizeMB * MB,
      maxCoverFileSize: formData.maxCoverFileSize,
      maxDuration: formData.maxDuration,
      freeWatchSeconds: formData.freeWatchSeconds,
      anchorDailyUploadLimit: formData.anchorDailyUploadLimit,
      normalUserDailyUploadLimit: formData.normalUserDailyUploadLimit,
      entryEnabled: formData.entryEnabled,
    })
    if (response?.success) {
      ElMessage.success(t('pages.shortVideoCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.shortVideoCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save short video cfg failed:', error)
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
  max-width: 720px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>

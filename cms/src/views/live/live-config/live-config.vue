<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.liveConfig.paidDanmakuPrice')" prop="paidDanmakuPrice">
          <el-input-number
              v-model="formData.paidDanmakuPrice"
              :min="0"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.liveConfig.paidDanmakuTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.liveConfig.privateRoomFreeWatchSeconds')" prop="privateRoomFreeWatchSeconds">
          <el-input-number
              v-model="formData.privateRoomFreeWatchSeconds"
              :min="0"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.liveConfig.privateRoomFreeWatchTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.liveConfig.lastUpdated')">
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
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {liveCfgApi} from '@/api/modules/liveCfg'
import type {LiveCfg} from '@/types/api'
import {NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  paidDanmakuPrice: 0,
  privateRoomFreeWatchSeconds: 420,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  paidDanmakuPrice: [
    {required: true, message: t('pages.liveConfig.paidDanmakuRequired'), trigger: 'blur'},
    {type: 'number', min: 0, message: t('pages.liveConfig.priceMinZero'), trigger: 'blur'},
  ],
  privateRoomFreeWatchSeconds: [
    {required: true, message: t('pages.liveConfig.freeWatchRequired'), trigger: 'blur'},
    {type: 'number', min: 0, message: t('pages.liveConfig.freeWatchMinZero'), trigger: 'blur'},
  ],
}))

const applyCfg = (cfg: LiveCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.paidDanmakuPrice = 0
    formData.privateRoomFreeWatchSeconds = 420
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.paidDanmakuPrice = truncateNumber(cfg.paidDanmakuPrice ?? 0)
  formData.privateRoomFreeWatchSeconds = cfg.privateRoomFreeWatchSeconds ?? 420
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await liveCfgApi.getLiveCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch live config failed:', error)
    ElMessage.error(t('pages.liveConfig.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await liveCfgApi.saveLiveCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      paidDanmakuPrice: formData.paidDanmakuPrice,
      privateRoomFreeWatchSeconds: formData.privateRoomFreeWatchSeconds,
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.liveConfig.saveFailed'))
    }
  } catch (error) {
    console.error('save live config failed:', error)
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

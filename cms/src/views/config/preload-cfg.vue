<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PreloadCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="200px">
        <el-form-item :label="t('pages.preloadCfg.recentLoginLimit')" prop="recentLoginLimit">
          <el-input-number v-model="formData.recentLoginLimit" :min="1" :max="10000" controls-position="right"/>
          <div class="form-tip">
            {{ t('pages.preloadCfg.recentLoginLimitTip') }}
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.preloadCfg.lastUpdated')">
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
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import preloadCfgApi from '@/api/modules/preload-cfg'
import type {PreloadCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  recentLoginLimit: 100,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  recentLoginLimit: [
    {required: true, message: t('pages.preloadCfg.preloadCountRequired'), trigger: 'change'},
    {type: 'number', min: 1, max: 10000, message: t('pages.preloadCfg.preloadCountRange'), trigger: 'change'},
  ],
}))

const applyCfg = (cfg: PreloadCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.recentLoginLimit = Number(cfg?.recentLoginLimit) || 100
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await preloadCfgApi.getPreloadCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch preload cfg failed:', error)
    ElMessage.error(t('pages.preloadCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const response = await preloadCfgApi.savePreloadCfg({
        id: formData.id === '0' ? 0 : Number(formData.id),
        recentLoginLimit: formData.recentLoginLimit,
      })
      if (response?.success) {
        ElMessage.success(t('pages.preloadCfg.saveSuccessRestart'))
        if (response.id) {
          formData.id = response.id
        }
        await fetchCfg()
      } else {
        ElMessage.error(t('pages.preloadCfg.saveFailed'))
      }
    } catch (error) {
      console.error('save preload cfg failed:', error)
    }
  })
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
}

.cfg-form {
  max-width: 720px;
}

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}
</style>

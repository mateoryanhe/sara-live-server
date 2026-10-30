<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ServerRuntimeCfgManagement') }}</span>
          <span v-if="metaInfo.updatedAt" class="header-meta">
            {{ t('pages.serverRuntimeCfg.lastUpdated') }}: {{ metaInfo.updatedAt }}
          </span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="200px">
        <el-tabs v-model="activeTab" class="config-tabs">
          <el-tab-pane :label="t('pages.serverRuntimeCfg.sectionPreload')" name="preload">
            <div class="tab-pane-body">
              <el-form-item :label="t('pages.serverRuntimeCfg.recentLoginLimit')" prop="recentLoginLimit">
                <el-input-number v-model="formData.recentLoginLimit" :min="1" :max="10000" controls-position="right"/>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.recentLoginLimitTip') }}</div>
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.serverRuntimeCfg.sectionRuntime')" name="runtime">
            <div class="tab-pane-body">
              <el-form-item :label="t('pages.serverRuntimeCfg.resourceMetricCollectionEnabled')">
                <el-switch
                    v-model="formData.resourceMetricCollectionEnabled"
                    :active-text="t('pages.serverRuntimeCfg.enabled')"
                    :inactive-text="t('pages.serverRuntimeCfg.disabled')"
                />
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.resourceMetricCollectionEnabledTip') }}</div>
              </el-form-item>
              <el-form-item :label="t('pages.serverRuntimeCfg.memoryLimitM')" prop="memoryLimitM">
                <el-input-number v-model="formData.memoryLimitM" :min="64" :max="32768" controls-position="right"/>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.memoryLimitMTip') }}</div>
              </el-form-item>
              <el-form-item :label="t('pages.serverRuntimeCfg.ipGeoDbPath')" prop="ipGeoDbPath">
                <el-input v-model="formData.ipGeoDbPath" maxlength="512" placeholder="/path/to/GeoLite2-Country.mmdb"/>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.ipGeoDbPathTip') }}</div>
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.serverRuntimeCfg.sectionHotRestart')" name="hotRestart">
            <div class="tab-pane-body">
              <el-form-item :label="t('pages.serverRuntimeCfg.hotRestartAuth')" prop="hotRestartAuth">
                <el-input v-model="formData.hotRestartAuth" maxlength="128" show-password autocomplete="new-password"/>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.hotRestartAuthTip') }}</div>
              </el-form-item>
              <el-form-item :label="t('pages.serverRuntimeCfg.hotRestartPhase1Time')" prop="hotRestartFlushTimeout">
                <el-input-number v-model="formData.hotRestartFlushTimeout" :max="3600" :min="1" controls-position="right"/>
                <span class="field-unit">{{ t('pages.serverRuntimeCfg.secondsUnit') }}</span>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.hotRestartPhase1TimeTip') }}</div>
              </el-form-item>
              <el-form-item :label="t('pages.serverRuntimeCfg.hotRestartPhase2Time')" prop="hotRestartExitTimeout">
                <el-input-number v-model="formData.hotRestartExitTimeout" :max="3600" :min="1" controls-position="right"/>
                <span class="field-unit">{{ t('pages.serverRuntimeCfg.secondsUnit') }}</span>
                <div class="form-tip">{{ t('pages.serverRuntimeCfg.hotRestartPhase2TimeTip') }}</div>
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
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import preloadCfgApi from '@/api/modules/preload-cfg'
import type {PreloadCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()
const activeTab = ref<'preload' | 'runtime' | 'hotRestart'>('preload')

const formData = reactive({
  id: '0',
  recentLoginLimit: 100,
  hotRestartAuth: '',
  hotRestartFlushTimeout: 60,
  hotRestartExitTimeout: 60,
  memoryLimitM: 300,
  ipGeoDbPath: '',
  resourceMetricCollectionEnabled: true,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  recentLoginLimit: [
    {required: true, message: t('pages.serverRuntimeCfg.preloadCountRequired'), trigger: 'change'},
    {type: 'number', min: 1, max: 10000, message: t('pages.serverRuntimeCfg.preloadCountRange'), trigger: 'change'},
  ],
  hotRestartAuth: [
    {required: true, message: t('pages.serverRuntimeCfg.hotRestartAuthRequired'), trigger: 'blur'},
    {min: 8, max: 128, message: t('pages.serverRuntimeCfg.hotRestartAuthLength'), trigger: 'blur'},
  ],
  hotRestartFlushTimeout: [
    {required: true, message: t('pages.serverRuntimeCfg.hotRestartPhase1TimeRequired'), trigger: 'change'},
    {type: 'number', min: 1, max: 3600, message: t('pages.serverRuntimeCfg.hotRestartPhase1TimeRange'), trigger: 'change'},
  ],
  hotRestartExitTimeout: [
    {required: true, message: t('pages.serverRuntimeCfg.hotRestartPhase2TimeRequired'), trigger: 'change'},
    {type: 'number', min: 1, max: 3600, message: t('pages.serverRuntimeCfg.hotRestartPhase2TimeRange'), trigger: 'change'},
  ],
  memoryLimitM: [
    {required: true, message: t('pages.serverRuntimeCfg.memoryLimitMRequired'), trigger: 'change'},
    {type: 'number', min: 64, max: 32768, message: t('pages.serverRuntimeCfg.memoryLimitMRange'), trigger: 'change'},
  ],
  ipGeoDbPath: [
    {required: true, message: t('pages.serverRuntimeCfg.ipGeoDbPathRequired'), trigger: 'blur'},
    {max: 512, message: t('pages.serverRuntimeCfg.ipGeoDbPathLength'), trigger: 'blur'},
  ],
}))

const applyCfg = (cfg: PreloadCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.recentLoginLimit = Number(cfg?.recentLoginLimit) || 100
  formData.hotRestartAuth = cfg?.hotRestartAuth || ''
  formData.hotRestartFlushTimeout = Number(cfg?.hotRestartFlushTimeout) || 60
  formData.hotRestartExitTimeout = Number(cfg?.hotRestartExitTimeout) || 60
  formData.memoryLimitM = Number(cfg?.memoryLimitM) || 300
  formData.ipGeoDbPath = cfg?.ipGeoDbPath || ''
  formData.resourceMetricCollectionEnabled = cfg?.resourceMetricCollectionEnabled !== false
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await preloadCfgApi.getPreloadCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch server runtime cfg failed:', error)
    ElMessage.error(t('pages.serverRuntimeCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid, fields) => {
    if (!valid) {
      const hotRestartKeys = ['hotRestartAuth', 'hotRestartFlushTimeout', 'hotRestartExitTimeout'] as const
      const runtimeKeys = ['memoryLimitM', 'ipGeoDbPath'] as const
      if (fields && hotRestartKeys.some((k) => fields[k])) {
        activeTab.value = 'hotRestart'
      } else if (fields && runtimeKeys.some((k) => fields[k])) {
        activeTab.value = 'runtime'
      } else if (fields) {
        activeTab.value = 'preload'
      }
      return
    }
    try {
      const response = await preloadCfgApi.savePreloadCfg({
        id: formData.id === '0' ? 0 : Number(formData.id),
        recentLoginLimit: formData.recentLoginLimit,
        hotRestartAuth: formData.hotRestartAuth.trim(),
        hotRestartFlushTimeout: formData.hotRestartFlushTimeout,
        hotRestartExitTimeout: formData.hotRestartExitTimeout,
        memoryLimitM: formData.memoryLimitM,
        ipGeoDbPath: formData.ipGeoDbPath.trim(),
        resourceMetricCollectionEnabled: formData.resourceMetricCollectionEnabled,
      })
      if (response?.success) {
        ElMessage.success(t('pages.serverRuntimeCfg.saveSuccessRestart'))
        if (response.id) {
          formData.id = response.id
        }
        await fetchCfg()
      } else {
        ElMessage.error(t('pages.serverRuntimeCfg.saveFailed'))
      }
    } catch (error) {
      console.error('save server runtime cfg failed:', error)
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
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.header-meta {
  font-size: 13px;
  font-weight: normal;
  color: var(--el-text-color-secondary);
}

.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.config-tabs :deep(.el-tabs__content) {
  padding-top: 4px;
}

.tab-pane-body {
  max-width: 760px;
  padding: 16px 4px 8px;
  min-height: 280px;
}

.cfg-form {
  max-width: 100%;
}

.field-unit {
  margin-left: 8px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}

.form-actions {
  max-width: 760px;
  margin-top: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>

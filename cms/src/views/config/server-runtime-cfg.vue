<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ServerRuntimeCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="200px">
        <el-divider content-position="left">{{ t('pages.serverRuntimeCfg.sectionPreload') }}</el-divider>
        <el-form-item :label="t('pages.serverRuntimeCfg.recentLoginLimit')" prop="recentLoginLimit">
          <el-input-number v-model="formData.recentLoginLimit" :min="1" :max="10000" controls-position="right"/>
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.recentLoginLimitTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.serverRuntimeCfg.initGold')" prop="initGold">
          <el-input-number
              v-model="formData.initGold"
              :min="0"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="1"
              controls-position="right"
          />
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.initGoldTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.serverRuntimeCfg.initDiamond')" prop="initDiamond">
          <el-input-number
              v-model="formData.initDiamond"
              :min="0"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="1"
              controls-position="right"
          />
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.initDiamondTip') }}</div>
        </el-form-item>

        <el-divider content-position="left">{{ t('pages.serverRuntimeCfg.sectionRuntime') }}</el-divider>
        <el-form-item :label="t('pages.serverRuntimeCfg.hotRestartAuth')" prop="hotRestartAuth">
          <el-input v-model="formData.hotRestartAuth" maxlength="128" show-password autocomplete="new-password"/>
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.hotRestartAuthTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.serverRuntimeCfg.memoryLimitM')" prop="memoryLimitM">
          <el-input-number v-model="formData.memoryLimitM" :min="64" :max="32768" controls-position="right"/>
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.memoryLimitMTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.serverRuntimeCfg.ipGeoDbPath')" prop="ipGeoDbPath">
          <el-input v-model="formData.ipGeoDbPath" maxlength="512" placeholder="/path/to/GeoLite2-Country.mmdb"/>
          <div class="form-tip">{{ t('pages.serverRuntimeCfg.ipGeoDbPathTip') }}</div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.serverRuntimeCfg.lastUpdated')">
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
import {NUMBER_INPUT_DECIMALS} from '@/utils/number-format'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  recentLoginLimit: 100,
  initGold: 0,
  initDiamond: 0,
  hotRestartAuth: '',
  memoryLimitM: 300,
  ipGeoDbPath: '',
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
  initGold: [
    {required: true, message: t('pages.serverRuntimeCfg.initGoldRequired'), trigger: 'change'},
    {type: 'number', min: 0, message: t('pages.serverRuntimeCfg.initGoldMin'), trigger: 'change'},
  ],
  initDiamond: [
    {required: true, message: t('pages.serverRuntimeCfg.initDiamondRequired'), trigger: 'change'},
    {type: 'number', min: 0, message: t('pages.serverRuntimeCfg.initDiamondMin'), trigger: 'change'},
  ],
  hotRestartAuth: [
    {required: true, message: t('pages.serverRuntimeCfg.hotRestartAuthRequired'), trigger: 'blur'},
    {min: 8, max: 128, message: t('pages.serverRuntimeCfg.hotRestartAuthLength'), trigger: 'blur'},
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
  formData.initGold = Number(cfg?.initGold) || 0
  formData.initDiamond = Number(cfg?.initDiamond) || 0
  formData.hotRestartAuth = cfg?.hotRestartAuth || ''
  formData.memoryLimitM = Number(cfg?.memoryLimitM) || 300
  formData.ipGeoDbPath = cfg?.ipGeoDbPath || ''
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
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const response = await preloadCfgApi.savePreloadCfg({
        id: formData.id === '0' ? 0 : Number(formData.id),
        recentLoginLimit: formData.recentLoginLimit,
        initGold: formData.initGold,
        initDiamond: formData.initDiamond,
        hotRestartAuth: formData.hotRestartAuth.trim(),
        memoryLimitM: formData.memoryLimitM,
        ipGeoDbPath: formData.ipGeoDbPath.trim(),
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

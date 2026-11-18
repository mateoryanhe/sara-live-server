<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AppsFlyerCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.appsFlyerCfg.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.appsFlyerCfg.noticeLine1') }}</p>
        <p>{{ t('pages.appsFlyerCfg.noticeLine2') }}</p>
        <p>{{ t('pages.appsFlyerCfg.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="160px">
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="t('pages.appsFlyerCfg.tabAndroid')" name="android">
            <el-form-item :label="t('pages.appsFlyerCfg.androidDevKey')">
              <el-input
                  v-model="formData.androidDevKey"
                  clearable
                  :placeholder="t('pages.appsFlyerCfg.androidDevKeyPlaceholder')"
              />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane :label="t('pages.appsFlyerCfg.tabIos')" name="ios">
            <el-form-item :label="t('pages.appsFlyerCfg.iosDevKey')">
              <el-input
                  v-model="formData.iosDevKey"
                  clearable
                  :placeholder="t('pages.appsFlyerCfg.iosDevKeyPlaceholder')"
              />
            </el-form-item>
            <el-form-item :label="t('pages.appsFlyerCfg.iosAppId')">
              <el-input
                  v-model="formData.iosAppId"
                  clearable
                  :placeholder="t('pages.appsFlyerCfg.iosAppIdPlaceholder')"
              />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.appsFlyerCfg.lastUpdated')">
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
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance} from 'element-plus'
import {appsFlyerCfgApi} from '@/api/modules/apps-flyer-cfg'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()
const activeTab = ref('android')

const formData = reactive({
  id: 0,
  androidDevKey: '',
  iosDevKey: '',
  iosAppId: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await appsFlyerCfgApi.getAppsFlyerCfg()
    const cfg = res.cfg
    if (!cfg) {
      formData.id = 0
      formData.androidDevKey = ''
      formData.iosDevKey = ''
      formData.iosAppId = ''
      metaInfo.updatedAt = ''
      return
    }
    formData.id = Number(cfg.id) || 0
    formData.androidDevKey = cfg.androidDevKey ?? ''
    formData.iosDevKey = cfg.iosDevKey ?? ''
    formData.iosAppId = cfg.iosAppId ?? ''
    metaInfo.updatedAt = cfg.updatedAt ?? ''
  } catch (e) {
    console.error('fetch apps flyer cfg failed:', e)
    ElMessage.error(t('common.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  const android = formData.androidDevKey.trim()
  const iosKey = formData.iosDevKey.trim()
  const iosId = formData.iosAppId.trim()
  if (!android && !iosKey && !iosId) {
    ElMessage.warning(t('pages.appsFlyerCfg.needOneField'))
    return
  }
  loading.value = true
  try {
    const res = await appsFlyerCfgApi.saveAppsFlyerCfg({
      id: formData.id,
      androidDevKey: android,
      iosDevKey: iosKey,
      iosAppId: iosId,
    })
    if (res.success) {
      ElMessage.success(t('pages.appsFlyerCfg.saveSuccess'))
      await fetchCfg()
    }
  } catch (e) {
    console.error('save apps flyer cfg failed:', e)
    ElMessage.error(t('common.saveFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCfg()
})
</script>

<style scoped>
.tip-alert {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 720px;
}
</style>
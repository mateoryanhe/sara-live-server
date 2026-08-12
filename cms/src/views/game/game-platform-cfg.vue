<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GamePlatformCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.gamePlatformCfg.noteTitle')"
          class="tip-alert"
          show-icon
          type="info"
      >
        <p>{{ t('pages.gamePlatformCfg.tipLine1') }}</p>
        <p>{{ t('pages.gamePlatformCfg.tipLine2') }}</p>
        <p>{{ t('pages.gamePlatformCfg.tipLine3') }}</p>
        <p>{{ t('pages.gamePlatformCfg.tipLine4') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.gamePlatformCfg.vendorUrl')" prop="vendorUrl">
          <el-input v-model="formData.vendorUrl" clearable :placeholder="t('pages.gamePlatformCfg.vendorUrlPlaceholder')"/>
        </el-form-item>

        <el-form-item :label="t('pages.gamePlatformCfg.iconUrl')" prop="iconUrl">
          <el-input v-model="formData.iconUrl" clearable :placeholder="t('pages.gamePlatformCfg.iconUrlPlaceholder')"/>
        </el-form-item>

        <el-form-item :label="t('pages.gamePlatformCfg.token')" prop="token">
          <el-input
              v-model="formData.token"
              clearable
              :placeholder="t('pages.gamePlatformCfg.tokenPlaceholder')"
              show-password
              type="password"
          />
        </el-form-item>

        <el-form-item :label="t('pages.gamePlatformCfg.secretKey')" prop="secretKey">
          <el-input
              v-model="formData.secretKey"
              clearable
              :placeholder="t('pages.gamePlatformCfg.secretKeyPlaceholder')"
              show-password
              type="password"
          />
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.gamePlatformCfg.lastUpdated')">
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
import {gamePlatformApi} from '@/api'
import type {GamePlatformCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  vendorUrl: 'https://gapi.win12.best',
  iconUrl: '',
  token: '',
  secretKey: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  vendorUrl: [
    {required: true, message: t('pages.gamePlatformCfg.vendorUrlRequired'), trigger: 'blur'},
    {max: 255, message: t('pages.gamePlatformCfg.vendorUrlMax'), trigger: 'blur'},
  ],
  iconUrl: [{max: 512, message: t('pages.gamePlatformCfg.iconUrlMax'), trigger: 'blur'}],
  token: [
    {required: true, message: t('pages.gamePlatformCfg.tokenRequired'), trigger: 'blur'},
    {max: 512, message: t('pages.gamePlatformCfg.tokenMax'), trigger: 'blur'},
  ],
  secretKey: [
    {required: true, message: t('pages.gamePlatformCfg.secretKeyRequired'), trigger: 'blur'},
    {max: 255, message: t('pages.gamePlatformCfg.secretKeyMax'), trigger: 'blur'},
  ],
}))

const applyCfg = (cfg: GamePlatformCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.vendorUrl = 'https://gapi.win12.best'
    formData.iconUrl = ''
    formData.token = ''
    formData.secretKey = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.vendorUrl = cfg.vendorUrl || 'https://gapi.win12.best'
  formData.iconUrl = cfg.iconUrl || ''
  formData.token = cfg.token || ''
  formData.secretKey = cfg.secretKey || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await gamePlatformApi.getGamePlatformCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch game platform cfg failed:', error)
    ElMessage.error(t('pages.gamePlatformCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await gamePlatformApi.saveGamePlatformCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      vendorUrl: formData.vendorUrl.trim(),
      iconUrl: formData.iconUrl.trim().replace(/\/+$/, ''),
      token: formData.token.trim(),
      secretKey: formData.secretKey.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.gamePlatformCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.gamePlatformCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save game platform cfg failed:', error)
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

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.cfg-form {
  max-width: 860px;
}

code {
  padding: 0 4px;
  background: #f4f4f5;
  border-radius: 4px;
}
</style>

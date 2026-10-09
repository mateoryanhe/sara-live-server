<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.MetaPixelCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.metaPixel.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.metaPixel.noticeLine1') }}</p>
        <p>{{ t('pages.metaPixel.noticeLine2') }}</p>
        <p>{{ t('pages.metaPixel.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.metaPixel.enabled')" prop="enabled">
          <el-switch v-model="formData.enabled" :active-value="1" :inactive-value="0"/>
        </el-form-item>
        <el-form-item :label="t('pages.metaPixel.pixelId')" prop="pixelId">
          <el-input v-model="formData.pixelId" clearable :placeholder="t('pages.metaPixel.pixelIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.metaPixel.accessToken')" prop="accessToken">
          <el-input
              v-model="formData.accessToken"
              :rows="4"
              :placeholder="t('pages.metaPixel.accessTokenPlaceholder')"
              type="textarea"
          />
        </el-form-item>
        <el-form-item :label="t('pages.metaPixel.testEventCode')" prop="testEventCode">
          <el-input v-model="formData.testEventCode" clearable :placeholder="t('pages.metaPixel.testEventCodePlaceholder')"/>
        </el-form-item>
        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.metaPixel.lastUpdated')">
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
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {metaPixelApi} from '@/api/modules/meta-pixel'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  enabled: 0 as 0 | 1,
  pixelId: '',
  accessToken: '',
  testEventCode: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  pixelId: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.metaPixel.pixelIdRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
  accessToken: [{
    validator: (_r, value, callback) => {
      if (formData.enabled !== 1) {
        callback()
        return
      }
      if (!String(value || '').trim()) {
        callback(new Error(t('pages.metaPixel.accessTokenRequired')))
        return
      }
      callback()
    },
    trigger: 'blur',
  }],
}))

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await metaPixelApi.getMetaPixelCfg()
    const cfg = res.cfg
    if (!cfg) {
      formData.id = '0'
      formData.enabled = 0
      formData.pixelId = ''
      formData.accessToken = ''
      formData.testEventCode = ''
      metaInfo.updatedAt = ''
      return
    }
    formData.id = cfg.id
    formData.enabled = cfg.enabled === 1 ? 1 : 0
    formData.pixelId = cfg.pixelId ?? ''
    formData.accessToken = cfg.accessToken ?? ''
    formData.testEventCode = cfg.testEventCode ?? ''
    metaInfo.updatedAt = cfg.updatedAt ?? ''
  } catch (e) {
    console.error('fetch meta pixel cfg failed:', e)
    ElMessage.error(t('common.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const res = await metaPixelApi.saveMetaPixelCfg({
        id: Number(formData.id) || 0,
        enabled: formData.enabled,
        pixelId: formData.pixelId.trim(),
        accessToken: formData.accessToken.trim(),
        testEventCode: formData.testEventCode.trim(),
      })
      if (res.success) {
        ElMessage.success(t('pages.metaPixel.saveSuccess'))
        await fetchCfg()
      }
    } catch (e) {
      console.error('save meta pixel cfg failed:', e)
      ElMessage.error(t('common.saveFailed'))
    } finally {
      loading.value = false
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
  font-size: 16px;
  font-weight: bold;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.cfg-form {
  max-width: 720px;
}
</style>

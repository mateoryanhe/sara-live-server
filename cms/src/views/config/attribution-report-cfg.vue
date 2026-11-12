<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AttributionReportCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.attrReportCfg.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.attrReportCfg.noticeLine1') }}</p>
        <p>{{ t('pages.attrReportCfg.noticeLine2') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.attrReportCfg.baseApi')" prop="baseApi">
          <el-input
              v-model="formData.baseApi"
              clearable
              :placeholder="t('pages.attrReportCfg.baseApiPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.attrReportCfg.secretKey')" prop="secretKey">
          <el-input
              v-model="formData.secretKey"
              clearable
              show-password
              :placeholder="t('pages.attrReportCfg.secretKeyPlaceholder')"
          />
        </el-form-item>
        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.attrReportCfg.lastUpdated')">
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
import {attrReportCfgApi} from '@/api/modules/attr-report-cfg'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  baseApi: '',
  secretKey: '',
})

const metaInfo = reactive({
  updatedAt: '',
})

const formRules = computed<FormRules>(() => ({
  baseApi: [
    {required: true, message: t('pages.attrReportCfg.baseApiRequired'), trigger: 'blur'},
  ],
  secretKey: [
    {required: true, message: t('pages.attrReportCfg.secretKeyRequired'), trigger: 'blur'},
  ],
}))

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await attrReportCfgApi.getAttrReportCfg()
    const cfg = res.cfg
    if (!cfg) {
      formData.id = '0'
      formData.baseApi = ''
      formData.secretKey = ''
      metaInfo.updatedAt = ''
      return
    }
    formData.id = cfg.id
    formData.baseApi = cfg.baseApi ?? ''
    formData.secretKey = cfg.secretKey ?? ''
    metaInfo.updatedAt = cfg.updatedAt ?? ''
  } catch (e) {
    console.error('fetch attr report cfg failed:', e)
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
      const res = await attrReportCfgApi.saveAttrReportCfg({
        id: Number(formData.id) || 0,
        baseApi: formData.baseApi.trim(),
        secretKey: formData.secretKey.trim(),
      })
      if (res.success) {
        ElMessage.success(t('pages.attrReportCfg.saveSuccess'))
        await fetchCfg()
      }
    } catch (e) {
      console.error('save attr report cfg failed:', e)
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

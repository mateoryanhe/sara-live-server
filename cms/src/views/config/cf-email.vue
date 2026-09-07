<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CfEmailCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.cfEmail.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.cfEmail.noticeLine1') }}</p>
        <p>{{ t('pages.cfEmail.noticeLine2') }}</p>
        <p>{{ t('pages.cfEmail.noticeLine3') }}</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.cfEmail.enabled')">
          <el-switch
              v-model="formData.enabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.accountId')" prop="accountId">
          <el-input v-model="formData.accountId" clearable :placeholder="t('pages.cfEmail.accountIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.apiToken')" prop="apiToken">
          <el-input
              v-model="formData.apiToken"
              clearable
              show-password
              type="password"
              :placeholder="t('pages.cfEmail.apiTokenPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.fromEmail')" prop="fromEmail">
          <el-input v-model="formData.fromEmail" clearable :placeholder="t('pages.cfEmail.fromEmailPlaceholder')"/>
          <span class="form-tip">{{ t('pages.cfEmail.fromEmailTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.cfEmail.lastUpdated')">
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
import {cfEmailApi} from '@/api/modules/cf-email'
import type {CfEmailCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  enabled: false,
  accountId: '',
  apiToken: '',
  fromEmail: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  accountId: [{required: true, message: t('pages.cfEmail.accountIdRequired'), trigger: 'blur'}],
  apiToken: [{required: true, message: t('pages.cfEmail.apiTokenRequired'), trigger: 'blur'}],
  fromEmail: [{required: true, message: t('pages.cfEmail.fromEmailRequired'), trigger: 'blur'}],
}))

const applyCfg = (cfg: CfEmailCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.accountId = ''
    formData.apiToken = ''
    formData.fromEmail = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.accountId = cfg.accountId || ''
  formData.apiToken = cfg.apiToken || ''
  formData.fromEmail = cfg.fromEmail || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await cfEmailApi.getCfEmailCfg()
    applyCfg(response.data?.cfg)
  } catch {
    ElMessage.error(t('pages.cfEmail.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    const idNum = Number(formData.id) || 0
    await cfEmailApi.saveCfEmailCfg({
      id: idNum,
      enabled: formData.enabled,
      accountId: formData.accountId.trim(),
      apiToken: formData.apiToken.trim(),
      fromEmail: formData.fromEmail.trim(),
    })
    ElMessage.success(t('pages.cfEmail.saveSuccess'))
    await fetchCfg()
  } catch {
    ElMessage.error(t('pages.cfEmail.saveFailed'))
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

.tip-alert p {
  margin: 4px 0;
  line-height: 1.5;
}

.cfg-form {
  max-width: 720px;
}

.form-tip {
  display: block;
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}
</style>

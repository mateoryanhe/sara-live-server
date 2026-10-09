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

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.cfEmail.enabled')">
          <el-switch
              v-model="formData.enabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.smtpHost')" prop="smtpHost">
          <el-input v-model="formData.smtpHost" clearable :placeholder="t('pages.cfEmail.smtpHostPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.smtpPort')" prop="smtpPort">
          <el-input-number v-model="formData.smtpPort" :max="65535" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.smtpUsername')" prop="smtpUsername">
          <el-input v-model="formData.smtpUsername" clearable :placeholder="t('pages.cfEmail.smtpUsernamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.smtpPassword')" prop="smtpPassword">
          <el-input
              v-model="formData.smtpPassword"
              clearable
              show-password
              type="password"
              :placeholder="t('pages.cfEmail.smtpPasswordPlaceholder')"
          />
          <span v-if="smtpPasswordConfigured" class="form-tip">{{ t('pages.cfEmail.smtpPasswordKeepHint') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.cfEmail.fromEmail')" prop="fromEmail">
          <el-input v-model="formData.fromEmail" clearable :placeholder="t('pages.cfEmail.fromEmailPlaceholder')"/>
          <span class="form-tip">{{ t('pages.cfEmail.fromEmailTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.cfEmail.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="can('save')" type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>

        <el-divider v-if="can('sendTest')"/>

        <template v-if="can('sendTest')">
          <p class="form-tip test-section-tip">{{ t('pages.cfEmail.testSectionHint') }}</p>
          <el-form-item :label="t('pages.cfEmail.testEmail')" prop="testEmail">
            <el-input
                v-model="testEmail"
                clearable
                :placeholder="t('pages.cfEmail.testEmailPlaceholder')"
            />
          </el-form-item>
          <el-form-item>
            <el-button :loading="testSending" type="success" @click="handleSendTest">
              {{ t('pages.cfEmail.sendTest') }}
            </el-button>
          </el-form-item>
        </template>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {cfEmailApi} from '@/api/modules/cf-email'
import type {CfEmailCfg} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t, locale} = useI18n()
const {can} = usePagePermission('CfEmailCfgManagement')
const loading = ref(false)
const formRef = ref()
const smtpPasswordConfigured = ref(false)
const smtpPasswordTouched = ref(false)
const testEmail = ref('')
const testSending = ref(false)

const TEST_EMAIL_STORAGE_KEY = 'cms.cfEmail.testEmail'

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const formData = reactive({
  id: '0',
  enabled: false,
  smtpHost: '',
  smtpPort: 587,
  smtpUsername: '',
  smtpPassword: '',
  fromEmail: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = computed(() => ({
  smtpHost: [{required: true, message: t('pages.cfEmail.smtpHostRequired'), trigger: 'blur'}],
  smtpPort: [{required: true, message: t('pages.cfEmail.smtpPortRequired'), trigger: 'change'}],
  smtpUsername: [{required: true, message: t('pages.cfEmail.smtpUsernameRequired'), trigger: 'blur'}],
  smtpPassword: smtpPasswordConfigured.value && !smtpPasswordTouched.value
      ? []
      : [{required: true, message: t('pages.cfEmail.smtpPasswordRequired'), trigger: 'blur'}],
  fromEmail: [{required: true, message: t('pages.cfEmail.fromEmailRequired'), trigger: 'blur'}],
}))

const applyCfg = (cfg: CfEmailCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.smtpHost = ''
    formData.smtpPort = 587
    formData.smtpUsername = ''
    formData.smtpPassword = ''
    formData.fromEmail = ''
    smtpPasswordConfigured.value = false
    smtpPasswordTouched.value = false
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.smtpHost = cfg.smtpHost || ''
  formData.smtpPort = cfg.smtpPort > 0 ? cfg.smtpPort : 587
  formData.smtpUsername = cfg.smtpUsername || ''
  formData.smtpPassword = ''
  smtpPasswordConfigured.value = !!cfg.smtpPasswordConfigured
  smtpPasswordTouched.value = false
  formData.fromEmail = cfg.fromEmail || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await cfEmailApi.getCfEmailCfg()
    applyCfg(response?.cfg)
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
      smtpHost: formData.smtpHost.trim(),
      smtpPort: formData.smtpPort || 587,
      smtpUsername: formData.smtpUsername.trim(),
      smtpPassword: smtpPasswordTouched.value ? formData.smtpPassword.trim() : '',
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

const handleSendTest = async () => {
  const addr = testEmail.value.trim()
  if (!addr) {
    ElMessage.warning(t('pages.cfEmail.testEmailRequired'))
    return
  }
  if (!emailPattern.test(addr)) {
    ElMessage.warning(t('pages.cfEmail.testEmailInvalid'))
    return
  }
  testSending.value = true
  try {
    const lang = locale.value.startsWith('zh') ? 'en' : locale.value.split('-')[0]
    await cfEmailApi.sendCfEmailTest({testEmail: addr, lang})
    localStorage.setItem(TEST_EMAIL_STORAGE_KEY, addr)
    ElMessage.success(t('pages.cfEmail.sendTestSuccess'))
  } catch {
    ElMessage.error(t('pages.cfEmail.sendTestFailed'))
  } finally {
    testSending.value = false
  }
}

watch(
  () => formData.smtpPassword,
  () => {
    smtpPasswordTouched.value = true
  },
)

onMounted(() => {
  const saved = localStorage.getItem(TEST_EMAIL_STORAGE_KEY)
  if (saved) {
    testEmail.value = saved
  }
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

.test-section-tip {
  margin: 0 0 16px 180px;
  max-width: 540px;
}
</style>

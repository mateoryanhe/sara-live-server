<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PrivacyPolicyCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.privacyPolicy.apiBase')" prop="apiBase">
          <el-input
              v-model="formData.apiBase"
              clearable
              :placeholder="t('pages.privacyPolicy.apiBasePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.apiBaseTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.privacyPolicyUrl')" prop="privacyPolicyUrl">
          <el-input
              v-model="formData.privacyPolicyUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.privacyPolicyPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.privacyPolicyTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.termsOfServiceUrl')" prop="termsOfServiceUrl">
          <el-input
              v-model="formData.termsOfServiceUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.termsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.termsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.creatorTermsUrl')" prop="creatorTermsUrl">
          <el-input
              v-model="formData.creatorTermsUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.creatorTermsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.creatorTermsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.roomOwnerTermsUrl')" prop="roomOwnerTermsUrl">
          <el-input
              v-model="formData.roomOwnerTermsUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.roomOwnerTermsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.roomOwnerTermsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.vipDescUrl')" prop="vipDescUrl">
          <el-input
              v-model="formData.vipDescUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.vipDescPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.vipDescTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.aboutSiteUrl')" prop="aboutSiteUrl">
          <el-input
              v-model="formData.aboutSiteUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.aboutSitePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.aboutSiteTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.privacyPolicy.safetyCenterUrl')" prop="safetyCenterUrl">
          <el-input
              v-model="formData.safetyCenterUrl"
              clearable
              :placeholder="t('pages.privacyPolicy.safetyCenterPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.privacyPolicy.safetyCenterTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.privacyPolicy.lastUpdated')">
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
import {privacyPolicyApi} from '@/api/modules/privacy-policy'
import type {PrivacyPolicyCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  apiBase: '',
  privacyPolicyUrl: '',
  termsOfServiceUrl: '',
  creatorTermsUrl: '',
  roomOwnerTermsUrl: '',
  vipDescUrl: '',
  aboutSiteUrl: '',
  safetyCenterUrl: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const validateApiBase = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const base = value?.trim()
  if (!base) {
    callback()
    return
  }
  if (!/^https?:\/\//i.test(base)) {
    callback(new Error(t('pages.privacyPolicy.apiBasePrefixInvalid')))
    return
  }
  callback()
}

const validateOptionalUrl = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const url = value?.trim()
  if (!url) {
    callback()
    return
  }
  if (/^https?:\/\//i.test(url)) {
    callback()
    return
  }
  if (formData.apiBase?.trim()) {
    callback()
    return
  }
  callback(new Error(t('pages.privacyPolicy.urlPrefixOrFullRequired')))
}

const urlFieldRules = () => [
  {max: 512, message: t('pages.privacyPolicy.urlMaxLength'), trigger: 'blur'},
  {validator: validateOptionalUrl, trigger: 'blur'},
]

const formRules = computed(() => ({
  apiBase: [
    {max: 512, message: t('pages.privacyPolicy.apiBaseMaxLength'), trigger: 'blur'},
    {validator: validateApiBase, trigger: 'blur'},
  ],
  privacyPolicyUrl: urlFieldRules(),
  termsOfServiceUrl: urlFieldRules(),
  creatorTermsUrl: urlFieldRules(),
  roomOwnerTermsUrl: urlFieldRules(),
  vipDescUrl: urlFieldRules(),
  aboutSiteUrl: urlFieldRules(),
  safetyCenterUrl: urlFieldRules(),
}))

const applyCfg = (cfg: PrivacyPolicyCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.apiBase = ''
    formData.privacyPolicyUrl = ''
    formData.termsOfServiceUrl = ''
    formData.creatorTermsUrl = ''
    formData.roomOwnerTermsUrl = ''
    formData.vipDescUrl = ''
    formData.aboutSiteUrl = ''
    formData.safetyCenterUrl = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.apiBase = cfg.apiBase || ''
  formData.privacyPolicyUrl = cfg.privacyPolicyUrl || ''
  formData.termsOfServiceUrl = cfg.termsOfServiceUrl || ''
  formData.creatorTermsUrl = cfg.creatorTermsUrl || ''
  formData.roomOwnerTermsUrl = cfg.roomOwnerTermsUrl || ''
  formData.vipDescUrl = cfg.vipDescUrl || ''
  formData.aboutSiteUrl = cfg.aboutSiteUrl || ''
  formData.safetyCenterUrl = cfg.safetyCenterUrl || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await privacyPolicyApi.getPrivacyPolicyCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch privacy policy cfg failed:', error)
    ElMessage.error(t('pages.privacyPolicy.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await privacyPolicyApi.savePrivacyPolicyCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      apiBase: formData.apiBase.trim(),
      privacyPolicyUrl: formData.privacyPolicyUrl.trim(),
      termsOfServiceUrl: formData.termsOfServiceUrl.trim(),
      creatorTermsUrl: formData.creatorTermsUrl.trim(),
      roomOwnerTermsUrl: formData.roomOwnerTermsUrl.trim(),
      vipDescUrl: formData.vipDescUrl.trim(),
      aboutSiteUrl: formData.aboutSiteUrl.trim(),
      safetyCenterUrl: formData.safetyCenterUrl.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.privacyPolicy.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.privacyPolicy.saveFailed'))
    }
  } catch (error) {
    console.error('save privacy policy cfg failed:', error)
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
  max-width: 760px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}
</style>

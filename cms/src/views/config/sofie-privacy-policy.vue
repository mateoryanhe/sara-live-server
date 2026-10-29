<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.SofiePrivacyPolicyCfgManagement') }}</span>
          <el-button
              v-if="can('sync')"
              :loading="syncingToTarget"
              type="warning"
              @click="handleSyncToTarget"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.sofiePrivacyPolicy.packageName')" prop="packageName">
          <el-input
              v-model="formData.packageName"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.packageNamePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.packageNameTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.apiBase')" prop="apiBase">
          <el-input
              v-model="formData.apiBase"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.apiBasePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.apiBaseTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.privacyPolicyUrl')" prop="privacyPolicyUrl">
          <el-input
              v-model="formData.privacyPolicyUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.privacyPolicyPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.privacyPolicyTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.termsOfServiceUrl')" prop="termsOfServiceUrl">
          <el-input
              v-model="formData.termsOfServiceUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.termsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.termsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.creatorTermsUrl')" prop="creatorTermsUrl">
          <el-input
              v-model="formData.creatorTermsUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.creatorTermsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.creatorTermsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.roomOwnerTermsUrl')" prop="roomOwnerTermsUrl">
          <el-input
              v-model="formData.roomOwnerTermsUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.roomOwnerTermsPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.roomOwnerTermsTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.vipDescUrl')" prop="vipDescUrl">
          <el-input
              v-model="formData.vipDescUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.vipDescPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.vipDescTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.aboutSiteUrl')" prop="aboutSiteUrl">
          <el-input
              v-model="formData.aboutSiteUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.aboutSitePlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.aboutSiteTip') }}</span>
        </el-form-item>

        <el-form-item :label="t('pages.sofiePrivacyPolicy.safetyCenterUrl')" prop="safetyCenterUrl">
          <el-input
              v-model="formData.safetyCenterUrl"
              clearable
              :placeholder="t('pages.sofiePrivacyPolicy.safetyCenterPlaceholder')"
          />
          <span class="form-tip">{{ t('pages.sofiePrivacyPolicy.safetyCenterTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.sofiePrivacyPolicy.lastUpdated')">
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
import {dataSyncApi} from '@/api/modules/data-sync'
import {sofiePrivacyPolicyApi} from '@/api/modules/sofie-privacy-policy'
import {confirmDataSync} from '@/utils/confirm-data-sync'
import {usePagePermission} from '@/composables/usePagePermission'
import type {PrivacyPolicyCfg} from '@/types/api'

const {t} = useI18n()
const {can} = usePagePermission('SofiePrivacyPolicyCfgManagement')
const loading = ref(false)
const syncingToTarget = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  packageName: '',
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
    callback(new Error(t('pages.sofiePrivacyPolicy.apiBasePrefixInvalid')))
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
  callback(new Error(t('pages.sofiePrivacyPolicy.urlPrefixOrFullRequired')))
}

const urlFieldRules = () => [
  {max: 512, message: t('pages.sofiePrivacyPolicy.urlMaxLength'), trigger: 'blur'},
  {validator: validateOptionalUrl, trigger: 'blur'},
]

const formRules = computed(() => ({
  packageName: [
    {max: 128, message: t('pages.sofiePrivacyPolicy.packageNameMaxLength'), trigger: 'blur'},
  ],
  apiBase: [
    {max: 512, message: t('pages.sofiePrivacyPolicy.apiBaseMaxLength'), trigger: 'blur'},
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
    formData.packageName = ''
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
  formData.packageName = cfg.packageName || ''
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
    const response = await sofiePrivacyPolicyApi.getSofiePrivacyPolicyCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch sofie privacy policy cfg failed:', error)
    ElMessage.error(t('pages.sofiePrivacyPolicy.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await sofiePrivacyPolicyApi.saveSofiePrivacyPolicyCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      packageName: formData.packageName.trim(),
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
      ElMessage.success(t('pages.sofiePrivacyPolicy.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.sofiePrivacyPolicy.saveFailed'))
    }
  } catch (error) {
    console.error('save sofie privacy policy cfg failed:', error)
  } finally {
    loading.value = false
  }
}

const handleSyncToTarget = async () => {
  if (!formData.id || formData.id === '0') {
    ElMessage.warning(t('pages.sofiePrivacyPolicy.syncNeedSave'))
    return
  }
  try {
    await confirmDataSync({
      detail: t('pages.sofiePrivacyPolicy.syncToTargetConfirm'),
      title: t('common.syncData'),
    })
  } catch {
    return
  }
  syncingToTarget.value = true
  try {
    const response = await dataSyncApi.syncSofiePrivacyPolicyCfg()
    if (response?.success) {
      ElMessage.success(response.message || t('pages.sofiePrivacyPolicy.syncToTargetSuccess'))
    } else {
      ElMessage.error(t('pages.sofiePrivacyPolicy.syncToTargetFailed'))
    }
  } catch (error) {
    console.error('sync sofie privacy policy to target failed:', error)
    ElMessage.error(t('pages.sofiePrivacyPolicy.syncToTargetFailedCheckConfig'))
  } finally {
    syncingToTarget.value = false
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
  font-size: 16px;
  font-weight: bold;
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

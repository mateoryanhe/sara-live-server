<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CMSDomainMappingManagement') }}</span>
        </div>
      </template>

      <el-alert
          :title="t('pages.cmsDomainMapping.noticeTitle')"
          :description="t('pages.cmsDomainMapping.noticeBody')"
          type="info"
          show-icon
          :closable="false"
          class="notice"
      />

      <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="cfg-form"
          label-width="170px"
      >
        <el-form-item :label="t('pages.cmsDomainMapping.domain')" prop="domain">
          <el-input
              v-model="formData.domain"
              :disabled="!can('save')"
              clearable
              maxlength="253"
              :placeholder="t('pages.cmsDomainMapping.domainPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.cmsDomainMapping.domainTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.cmsDomainMapping.urlPrefix')">
          <el-input v-model="formData.urlPrefix" disabled/>
          <div class="form-tip">{{ t('pages.cmsDomainMapping.urlPrefixTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('pages.cmsDomainMapping.root')" prop="root">
          <el-input
              v-model="formData.root"
              :disabled="!can('save')"
              clearable
              maxlength="1024"
              :placeholder="t('pages.cmsDomainMapping.rootPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.cmsDomainMapping.rootTip') }}</div>
        </el-form-item>

        <el-form-item v-if="formData.updatedAt" :label="t('pages.cmsDomainMapping.updatedAt')">
          <span>{{ formData.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="can('save')" type="primary" :loading="saving" @click="handleSave">
            {{ t('pages.cmsDomainMapping.save') }}
          </el-button>
          <el-button :disabled="loading" @click="fetchMapping">
            {{ t('pages.cmsDomainMapping.reload') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {cmsDomainMappingApi, type CMSDomainSiteMapping} from '@/api/modules/cms-domain-mapping'
import {usePagePermission} from '@/composables/usePagePermission'

const DEFAULT_ROOT = '/home/ec2-user/cdn/cms'
const DEFAULT_PREFIX = '/cms'

const {t} = useI18n()
const {can} = usePagePermission('CMSDomainMappingManagement')
const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)

const formData = reactive<CMSDomainSiteMapping>({
  id: '0',
  domain: '',
  urlPrefix: DEFAULT_PREFIX,
  root: DEFAULT_ROOT,
  updatedAt: '',
})

function isAbsolutePath(value: string): boolean {
  if (value.startsWith('/')) {
    return true
  }
  if (/^[A-Za-z]:[/\\]/.test(value)) {
    return true
  }
  return value.startsWith('\\\\') || value.startsWith('//')
}

function isSingleDomain(value: string): boolean {
  const domain = value.trim()
  if (!domain || domain.includes(',')) {
    return false
  }
  return !domain.includes('://')
      && !/[\\/?#@\s]/.test(domain)
      && domain.length <= 253
}

const formRules = computed<FormRules>(() => ({
  domain: [
    {required: true, message: t('pages.cmsDomainMapping.domainRequired'), trigger: 'blur'},
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (!isSingleDomain(value || '')) {
          callback(new Error(t('pages.cmsDomainMapping.domainInvalid')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  root: [
    {required: true, message: t('pages.cmsDomainMapping.rootRequired'), trigger: 'blur'},
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (!isAbsolutePath((value || '').trim())) {
          callback(new Error(t('pages.cmsDomainMapping.rootAbsolute')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
}))

function applyMapping(mapping?: CMSDomainSiteMapping | null) {
  formData.id = mapping?.id || '0'
  formData.domain = mapping?.domain || ''
  formData.urlPrefix = mapping?.urlPrefix || DEFAULT_PREFIX
  formData.root = mapping?.root || DEFAULT_ROOT
  formData.updatedAt = mapping?.updatedAt || ''
}

async function fetchMapping() {
  loading.value = true
  try {
    const response = await cmsDomainMappingApi.get()
    applyMapping(response.mapping)
    formRef.value?.clearValidate()
  } catch (error) {
    console.error('fetch CMS domain mapping failed:', error)
    ElMessage.error(t('pages.cmsDomainMapping.fetchFailed'))
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    const response = await cmsDomainMappingApi.save({
      domain: formData.domain.trim(),
      root: formData.root.trim(),
    })
    applyMapping(response.mapping)
    ElMessage.success(t('pages.cmsDomainMapping.saveSuccess'))
  } catch (error) {
    console.error('save CMS domain mapping failed:', error)
    ElMessage.error(t('pages.cmsDomainMapping.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(fetchMapping)
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

.notice {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 920px;
}

.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
  line-height: 1.5;
}
</style>

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.ExtraDomainMappingManagement') }}</span>
      </template>

      <el-alert
          :title="t('pages.extraDomainMapping.noticeTitle')"
          :description="t('pages.extraDomainMapping.noticeBody')"
          type="info"
          show-icon
          :closable="false"
          class="notice"
      />

      <el-form
          v-if="can('save')"
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="add-form"
          label-width="120px"
      >
        <el-form-item :label="t('pages.extraDomainMapping.domain')" prop="domain">
          <el-input
              v-model="formData.domain"
              clearable
              maxlength="253"
              :placeholder="t('pages.extraDomainMapping.domainPlaceholder')"
              style="width: 360px"
          />
          <div class="form-tip">{{ t('pages.extraDomainMapping.domainTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.extraDomainMapping.root')" prop="root">
          <el-input
              v-model="formData.root"
              clearable
              maxlength="1024"
              :placeholder="t('pages.extraDomainMapping.rootPlaceholder')"
              style="width: 480px"
          />
          <div class="form-tip">{{ t('pages.extraDomainMapping.rootTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.extraDomainMapping.certFile')" prop="certFile" required>
          <el-upload
              :auto-upload="false"
              :limit="1"
              :on-change="onCertChange"
              :on-remove="onCertRemove"
              accept=".pem,.crt,.cer,.cert"
              :file-list="certFileList"
          >
            <el-button>{{ t('pages.extraDomainMapping.selectCert') }}</el-button>
          </el-upload>
          <div class="form-tip">{{ t('pages.extraDomainMapping.certTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.extraDomainMapping.keyFile')" prop="keyFile" required>
          <el-upload
              :auto-upload="false"
              :limit="1"
              :on-change="onKeyChange"
              :on-remove="onKeyRemove"
              accept=".pem,.key"
              :file-list="keyFileList"
          >
            <el-button>{{ t('pages.extraDomainMapping.selectKey') }}</el-button>
          </el-upload>
          <div class="form-tip">{{ t('pages.extraDomainMapping.keyTip') }}</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleAdd">
            {{ t('pages.extraDomainMapping.add') }}
          </el-button>
          <span v-if="uploadPercent > 0 && uploadPercent < 100" class="upload-percent">{{ uploadPercent }}%</span>
        </el-form-item>
      </el-form>

      <el-table :data="items" border stripe class="mapping-table">
        <el-table-column prop="domain" :label="t('pages.extraDomainMapping.columnDomain')" min-width="200"/>
        <el-table-column prop="root" :label="t('pages.extraDomainMapping.columnRoot')" min-width="240" show-overflow-tooltip/>
        <el-table-column :label="t('pages.extraDomainMapping.columnTls')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.hasTls ? 'success' : 'info'" size="small">
              {{ row.hasTls ? t('pages.extraDomainMapping.tlsYes') : t('pages.extraDomainMapping.tlsNo') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="t('pages.extraDomainMapping.columnUpdatedAt')" min-width="170"/>
        <el-table-column v-if="can('delete')" width="100" align="center" :label="t('common.delete')">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <span>{{ t('pages.extraDomainMapping.empty') }}</span>
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile, type UploadUserFile} from 'element-plus'
import {extraDomainMappingApi, type ExtraDomainMapping} from '@/api/modules/extra-domain-mapping'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('ExtraDomainMappingManagement')
const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const uploadPercent = ref(0)
const items = ref<ExtraDomainMapping[]>([])
const certFileList = ref<UploadUserFile[]>([])
const keyFileList = ref<UploadUserFile[]>([])

const formData = reactive({
  domain: '',
  root: '',
  certFile: null as File | null,
  keyFile: null as File | null,
})

function isAbsolutePath(value: string): boolean {
  const v = value.trim()
  if (v.startsWith('/')) {
    return true
  }
  return /^[A-Za-z]:[/\\]/.test(v) || v.startsWith('\\\\')
}

function isSingleDomain(value: string): boolean {
  const domain = value.trim()
  return Boolean(domain)
      && !domain.includes(',')
      && !domain.includes('://')
      && !/[\\/?#@\s]/.test(domain)
      && domain.length <= 253
}

const formRules = computed<FormRules>(() => ({
  domain: [
    {required: true, message: t('pages.extraDomainMapping.domainRequired'), trigger: 'blur'},
    {
      validator: (_r, value, callback) => {
        if (!isSingleDomain(value || '')) {
          callback(new Error(t('pages.extraDomainMapping.domainInvalid')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  root: [
    {required: true, message: t('pages.extraDomainMapping.rootRequired'), trigger: 'blur'},
    {
      validator: (_r, value, callback) => {
        if (!isAbsolutePath(value || '')) {
          callback(new Error(t('pages.extraDomainMapping.rootAbsolute')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  certFile: [
    {
      validator: (_r, _v, callback) => {
        if (!formData.certFile) {
          callback(new Error(t('pages.extraDomainMapping.certRequired')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
  keyFile: [
    {
      validator: (_r, _v, callback) => {
        if (!formData.keyFile) {
          callback(new Error(t('pages.extraDomainMapping.keyRequired')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
}))

function onCertChange(file: UploadFile) {
  formData.certFile = file.raw || null
  certFileList.value = file.raw ? [{name: file.name, status: 'ready'}] : []
  formRef.value?.validateField('certFile')
}

function onCertRemove() {
  formData.certFile = null
  certFileList.value = []
}

function onKeyChange(file: UploadFile) {
  formData.keyFile = file.raw || null
  keyFileList.value = file.raw ? [{name: file.name, status: 'ready'}] : []
  formRef.value?.validateField('keyFile')
}

function onKeyRemove() {
  formData.keyFile = null
  keyFileList.value = []
}

async function fetchList() {
  loading.value = true
  try {
    const res = await extraDomainMappingApi.list()
    items.value = res.items || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.extraDomainMapping.fetchFailed'))
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  if (!formData.certFile || !formData.keyFile) {
    ElMessage.warning(t('pages.extraDomainMapping.certRequired'))
    return
  }
  saving.value = true
  uploadPercent.value = 0
  try {
    await extraDomainMappingApi.saveWithTLS(
        {
          domain: formData.domain.trim(),
          root: formData.root.trim(),
          cert: formData.certFile,
          key: formData.keyFile,
        },
        (p) => {
          uploadPercent.value = p
        },
    )
    ElMessage.success(t('pages.extraDomainMapping.addSuccess'))
    formData.domain = ''
    formData.root = ''
    formData.certFile = null
    formData.keyFile = null
    certFileList.value = []
    keyFileList.value = []
    formRef.value?.clearValidate()
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.extraDomainMapping.addFailed'))
  } finally {
    saving.value = false
    uploadPercent.value = 0
  }
}

async function handleDelete(row: ExtraDomainMapping) {
  try {
    await ElMessageBox.confirm(
        t('pages.extraDomainMapping.deleteConfirm', {domain: row.domain}),
        t('common.confirm'),
        {type: 'warning'},
    )
    await extraDomainMappingApi.remove(row.id)
    ElMessage.success(t('pages.extraDomainMapping.deleteSuccess'))
    await fetchList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(t('pages.extraDomainMapping.deleteFailed'))
    }
  }
}

onMounted(fetchList)
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.notice {
  margin-bottom: 16px;
}

.add-form {
  max-width: 760px;
  margin-bottom: 20px;
}

.mapping-table {
  width: 100%;
}

.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
  line-height: 1.5;
}

.upload-percent {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>

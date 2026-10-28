<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.HostTlsMappingManagement') }}</span>
      </template>

      <el-alert
          :title="t('pages.hostTlsMapping.noticeTitle')"
          :description="t('pages.hostTlsMapping.noticeBody')"
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
        <el-form-item :label="t('pages.hostTlsMapping.domain')" prop="domain">
          <el-input
              v-model="formData.domain"
              clearable
              maxlength="253"
              :placeholder="t('pages.hostTlsMapping.domainPlaceholder')"
              style="width: 360px"
          />
          <div class="form-tip">{{ t('pages.hostTlsMapping.domainTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.hostTlsMapping.certFile')" prop="certFile" required>
          <el-upload
              :auto-upload="false"
              :limit="1"
              :on-change="onCertChange"
              :on-remove="onCertRemove"
              accept=".pem,.crt,.cer,.cert"
              :file-list="certFileList"
          >
            <el-button>{{ t('pages.hostTlsMapping.selectCert') }}</el-button>
          </el-upload>
          <div class="form-tip">{{ t('pages.hostTlsMapping.certTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.hostTlsMapping.keyFile')" prop="keyFile" required>
          <el-upload
              :auto-upload="false"
              :limit="1"
              :on-change="onKeyChange"
              :on-remove="onKeyRemove"
              accept=".pem,.key"
              :file-list="keyFileList"
          >
            <el-button>{{ t('pages.hostTlsMapping.selectKey') }}</el-button>
          </el-upload>
          <div class="form-tip">{{ t('pages.hostTlsMapping.keyTip') }}</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleAdd">
            {{ t('pages.hostTlsMapping.add') }}
          </el-button>
          <span v-if="uploadPercent > 0 && uploadPercent < 100" class="upload-percent">{{ uploadPercent }}%</span>
        </el-form-item>
      </el-form>

      <el-table :data="items" border stripe class="mapping-table">
        <el-table-column prop="domain" :label="t('pages.hostTlsMapping.columnDomain')" min-width="220"/>
        <el-table-column :label="t('pages.hostTlsMapping.columnTls')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.hasTls ? 'success' : 'info'" size="small">
              {{ row.hasTls ? t('pages.hostTlsMapping.tlsYes') : t('pages.hostTlsMapping.tlsNo') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="t('pages.hostTlsMapping.columnUpdatedAt')" min-width="170"/>
        <el-table-column v-if="can('delete')" width="100" align="center" :label="t('common.delete')">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <span>{{ t('pages.hostTlsMapping.empty') }}</span>
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile, type UploadUserFile} from 'element-plus'
import {hostTlsMappingApi, type HostTlsMapping} from '@/api/modules/host-tls-mapping'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('HostTlsMappingManagement')
const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const uploadPercent = ref(0)
const items = ref<HostTlsMapping[]>([])
const certFileList = ref<UploadUserFile[]>([])
const keyFileList = ref<UploadUserFile[]>([])

const formData = reactive({
  domain: '',
  certFile: null as File | null,
  keyFile: null as File | null,
})

function isSingleDomain(value: string): boolean {
  const domain = value.trim().toLowerCase()
  if (!domain || domain.includes(',') || domain.includes('://') || /[\\/?#@\s:]/.test(domain) || domain.length > 253) {
    return false
  }
  // 精确 Host,或单级泛域名 *.example.com
  if (domain.startsWith('*.')) {
    const rest = domain.slice(2)
    const labels = rest.split('.')
    return labels.length >= 2 && labels.every((l) => /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(l))
  }
  const labels = domain.split('.')
  return labels.length >= 2 && labels.every((l) => /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(l) || l.length === 1 && /[a-z0-9]/.test(l))
}

const formRules = computed<FormRules>(() => ({
  domain: [
    {required: true, message: t('pages.hostTlsMapping.domainRequired'), trigger: 'blur'},
    {
      validator: (_r, value, callback) => {
        if (!isSingleDomain(value || '')) {
          callback(new Error(t('pages.hostTlsMapping.domainInvalid')))
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
          callback(new Error(t('pages.hostTlsMapping.certRequired')))
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
          callback(new Error(t('pages.hostTlsMapping.keyRequired')))
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
    const res = await hostTlsMappingApi.list()
    items.value = res.items || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.hostTlsMapping.fetchFailed'))
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
    ElMessage.warning(t('pages.hostTlsMapping.certRequired'))
    return
  }
  saving.value = true
  uploadPercent.value = 0
  try {
    await hostTlsMappingApi.saveWithTLS(
        {
          domain: formData.domain.trim(),
          cert: formData.certFile,
          key: formData.keyFile,
        },
        (p) => {
          uploadPercent.value = p
        },
    )
    ElMessage.success(t('pages.hostTlsMapping.addSuccess'))
    formData.domain = ''
    formData.certFile = null
    formData.keyFile = null
    certFileList.value = []
    keyFileList.value = []
    formRef.value?.clearValidate()
    await fetchList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.hostTlsMapping.addFailed'))
  } finally {
    saving.value = false
    uploadPercent.value = 0
  }
}

async function handleDelete(row: HostTlsMapping) {
  try {
    await ElMessageBox.confirm(
        t('pages.hostTlsMapping.deleteConfirm', {domain: row.domain}),
        t('common.confirm'),
        {type: 'warning'},
    )
    await hostTlsMappingApi.remove(row.id)
    ElMessage.success(t('pages.hostTlsMapping.deleteSuccess'))
    await fetchList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(t('pages.hostTlsMapping.deleteFailed'))
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

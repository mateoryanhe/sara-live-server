<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ApkPromotionChannelListManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <el-alert
            :closable="false"
            show-icon
            type="info"
            class="domain-hint"
            :title="shareBaseHost ? t('pages.apkPromotionChannelList.shareBaseHostCurrent', {host: shareBaseHost}) : t('pages.apkPromotionChannelList.shareBaseHostEmpty')"
        >
          <template #default>
            <p class="hint-text">{{ t('pages.apkPromotionChannelList.shareBaseHostHint') }}</p>
          </template>
        </el-alert>

        <div class="table-header">
          <el-input
              v-model="searchKey"
              clearable
              class="search-input"
              :placeholder="t('pages.apkPromotionChannelList.searchPlaceholder')"
              @keyup.enter="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button v-if="can('create')" type="primary" :disabled="!shareBaseHost" @click="handleAdd">{{ t('pages.apkPromotionChannelList.add') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.apkPromotionChannelList.channelName')" prop="channelName" min-width="140"/>
          <el-table-column :label="t('pages.apkPromotionChannelList.channelCode')" prop="channelCode" width="140"/>
          <el-table-column :label="t('pages.apkPromotionChannelList.shareUrl')" min-width="280">
            <template #default="{ row }">
              <div class="share-url-cell">
                <span class="share-url-text" :title="row.shareUrl">{{ row.shareUrl || '-' }}</span>
                <el-button v-if="row.shareUrl" link type="primary" @click="() => copyShareUrl(row.shareUrl)">
                  {{ t('pages.apkPromotionChannelList.copyShareUrl') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.apkPromotionChannelList.apkFileName')" prop="apkFileName" min-width="240">
            <template #default="{ row }">{{ row.apkFileName || '-' }}</template>
          </el-table-column>
          <el-table-column fixed="right" :label="t('common.actions')" width="240">
            <template #default="{ row }">
              <el-button v-if="can('edit')" size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button v-if="can('uploadApk')" size="small" type="primary" @click="openUpload(row)">{{ t('pages.apkPromotionChannelList.uploadApk') }}</el-button>
              <el-button v-if="can('delete')" size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="120px">
        <el-form-item :label="t('pages.apkPromotionChannelList.channelName')" prop="channelName">
          <el-input v-model="currentRow.channelName" maxlength="128" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('pages.apkPromotionChannelList.channelCode')" prop="channelCode">
          <el-input v-model="currentRow.channelCode" maxlength="64" show-word-limit :disabled="!!currentRow.id"/>
          <div class="form-tip">{{ t('pages.apkPromotionChannelList.channelCodeTip') }}</div>
        </el-form-item>
        <el-form-item v-if="previewShareUrl" :label="t('pages.apkPromotionChannelList.shareUrl')">
          <el-input :model-value="previewShareUrl" readonly/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="uploadVisible" :title="t('pages.apkPromotionChannelList.uploadDialogTitle')" width="480px" destroy-on-close>
      <p class="upload-channel-name">{{ uploadChannel?.channelName }} ({{ uploadChannel?.channelCode }})</p>
      <p class="form-tip">{{ t('pages.apkPromotionChannelList.uploadVersionHint') }}</p>
      <el-upload
          :auto-upload="false"
          :limit="1"
          accept=".apk"
          :on-change="handleApkFileChange"
          :on-remove="handleApkRemove"
      >
        <el-button type="primary">{{ t('pages.apkPromotionChannelList.uploadSelectFile') }}</el-button>
      </el-upload>
      <el-progress v-if="uploading" :percentage="uploadPercent" :status="uploadStatus || undefined" style="margin-top: 16px"/>
      <template #footer>
        <el-button @click="uploadVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="uploading" :disabled="!selectedApkFile" @click="handleApkUpload">
          {{ t('pages.apkPromotionChannelList.uploadStart') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile} from 'element-plus'
import {apkPromotionChannelApi} from '@/api/modules/apk-promotion-channel'
import {usePagePermission} from '@/composables/usePagePermission'
import type {ApkPromotionChannel} from '@/types/api'

const CHANNEL_CODE_RE = /^[a-zA-Z0-9_-]{1,64}$/

/** Element Plus 表单 validate 失败时 reject 的是非 Error 对象 */
const isElementPlusValidateError = (error: unknown): boolean => {
  if (!error || typeof error !== 'object' || error instanceof Error) {
    return false
  }
  return !('response' in error)
}

/** axios 业务错误已在 request 拦截器弹过后端 message，此处避免重复 */
const isHandledBizRequestError = (error: unknown): boolean => {
  return error instanceof Error && /^\d+$/.test(error.message)
}

const resolveRequestErrorMessage = (error: unknown): string => {
  if (error instanceof Error && error.message && !/^\d+$/.test(error.message)) {
    return error.message
  }
  const data = (error as { response?: { data?: { message?: string } } })?.response?.data
  return data?.message?.trim() || ''
}

const notifyOperationFailed = (error: unknown, fallbackKey: 'saveFailed' | 'deleteFailed' | 'uploadFailed') => {
  if (isElementPlusValidateError(error) || isHandledBizRequestError(error)) {
    return
  }
  const detail = resolveRequestErrorMessage(error)
  ElMessage.error(detail || t(`pages.apkPromotionChannelList.${fallbackKey}`))
}

interface ChannelForm {
  id: string
  channelName: string
  channelCode: string
}

const {t} = useI18n()
const {can} = usePagePermission('ApkPromotionChannelListManagement')
const loading = ref(false)
const tableData = ref<ApkPromotionChannel[]>([])
const total = ref(0)
const shareBaseHost = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const searchKey = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): ChannelForm => ({id: '', channelName: '', channelCode: ''})
const currentRow = ref<ChannelForm>(defaultForm())
const formRef = ref<FormInstance>()

const formRules = computed<FormRules>(() => ({
  channelName: [{required: true, message: t('pages.apkPromotionChannelList.channelNameRequired'), trigger: 'blur'}],
  channelCode: [
    {required: true, message: t('pages.apkPromotionChannelList.channelCodeRequired'), trigger: 'blur'},
    {
      validator: (_rule, value, callback) => {
        if (!value || !CHANNEL_CODE_RE.test(String(value).trim())) {
          callback(new Error(t('pages.apkPromotionChannelList.channelCodePattern')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
}))

const previewShareUrl = computed(() => {
  const code = currentRow.value.channelCode.trim()
  const host = shareBaseHost.value.trim()
  if (!host || !CHANNEL_CODE_RE.test(code)) {
    return ''
  }
  return `https://${host}/${encodeURIComponent(code)}`
})

const uploadVisible = ref(false)
const uploadChannel = ref<ApkPromotionChannel | null>(null)
const selectedApkFile = ref<File | null>(null)
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref<'success' | 'exception' | ''>('')

const fetchList = async () => {
  loading.value = true
  try {
    const response = await apkPromotionChannelApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
      key: searchKey.value.trim(),
    })
    tableData.value = response.data
    total.value = response.total
    shareBaseHost.value = (response.shareBaseHost || '').trim()
  } catch (error) {
    console.error('fetch apk promotion channel list failed:', error)
    ElMessage.error(t('pages.apkPromotionChannelList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const copyShareUrl = async (link: string) => {
  try {
    await navigator.clipboard.writeText(link)
    ElMessage.success(t('pages.apkPromotionChannelList.copyShareUrlSuccess'))
  } catch {
    ElMessage.error(t('pages.apkPromotionChannelList.copyShareUrlFailed'))
  }
}

const handleSearch = () => {
  currentPage.value = 1
  void fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  void fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  void fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.apkPromotionChannelList.add')
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: ApkPromotionChannel) => {
  dialogTitle.value = t('pages.apkPromotionChannelList.edit')
  currentRow.value = {
    id: row.id,
    channelName: row.channelName || '',
    channelCode: row.channelCode || '',
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const payload = {
      channelName: currentRow.value.channelName.trim(),
      channelCode: currentRow.value.channelCode.trim(),
    }
    if (currentRow.value.id) {
      await apkPromotionChannelApi.update({id: currentRow.value.id, ...payload})
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await apkPromotionChannelApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    void fetchList()
  } catch (error) {
    console.error('save apk promotion channel failed:', error)
    if (isElementPlusValidateError(error)) {
      ElMessage.warning(t('pages.apkPromotionChannelList.formInvalid'))
      return
    }
    notifyOperationFailed(error, 'saveFailed')
  }
}

const handleDelete = async (row: ApkPromotionChannel) => {
  try {
    await ElMessageBox.confirm(
        t('pages.apkPromotionChannelList.deleteConfirm', {name: row.channelName}),
        t('common.confirmDelete'),
        {confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning'},
    )
    await apkPromotionChannelApi.delete(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    void fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete apk promotion channel failed:', error)
      notifyOperationFailed(error, 'deleteFailed')
    }
  }
}

const openUpload = (row: ApkPromotionChannel) => {
  uploadChannel.value = row
  selectedApkFile.value = null
  uploadPercent.value = 0
  uploadStatus.value = ''
  uploadVisible.value = true
}

const handleApkFileChange = (uploadFile: UploadFile) => {
  selectedApkFile.value = uploadFile.raw ?? null
}

const handleApkRemove = () => {
  selectedApkFile.value = null
}

const wait = (ms: number) => new Promise<void>((resolve) => window.setTimeout(resolve, ms))

const uploadApkInChunks = async (file: File, channelId: string, onProgress: (p: number) => void) => {
  const session = await apkPromotionChannelApi.initApkUpload(channelId, file.name, file.size)
  if (!session.uploadId || session.chunkSize <= 0 || session.totalChunk <= 0) {
    throw new Error('invalid upload session')
  }
  try {
    for (let chunkIndex = 0; chunkIndex < session.totalChunk; chunkIndex += 1) {
      const start = chunkIndex * session.chunkSize
      const end = Math.min(start + session.chunkSize, file.size)
      const chunk = file.slice(start, end)
      await apkPromotionChannelApi.uploadApkChunk(session.uploadId, chunkIndex, chunk, file.name, (chunkPercent) => {
        const loaded = start + (chunk.size * chunkPercent) / 100
        onProgress(Math.min(99, Math.round((loaded * 100) / file.size)))
      })
      onProgress(Math.min(99, Math.round((end * 100) / file.size)))
    }
    await apkPromotionChannelApi.completeApkUpload(session.uploadId, channelId)
    onProgress(100)
  } catch (error) {
    try {
      await apkPromotionChannelApi.abortApkUpload(session.uploadId)
    } catch {
      // ignore
    }
    throw error
  }
}

const handleApkUpload = async () => {
  const file = selectedApkFile.value
  const channel = uploadChannel.value
  if (!file || !channel?.id) {
    return
  }
  if (!file.name.toLowerCase().endsWith('.apk')) {
    ElMessage.error(t('pages.apkPromotionChannelList.apkOnly'))
    return
  }
  uploading.value = true
  uploadPercent.value = 0
  uploadStatus.value = ''
  try {
    await uploadApkInChunks(file, channel.id, (p) => {
      uploadPercent.value = p
    })
    uploadStatus.value = 'success'
    ElMessage.success(t('pages.apkPromotionChannelList.uploadSuccess'))
    uploadVisible.value = false
    void fetchList()
  } catch (error) {
    uploadStatus.value = 'exception'
    console.error('upload apk promotion channel failed:', error)
    notifyOperationFailed(error, 'uploadFailed')
  } finally {
    uploading.value = false
  }
}

onMounted(() => {
  void fetchList()
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

.domain-hint {
  margin-bottom: 16px;
}

.hint-text {
  margin: 4px 0 0;
  font-size: 13px;
  line-height: 1.5;
}

.table-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.search-input {
  width: 280px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
}

.share-url-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.share-url-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-channel-name {
  margin-bottom: 12px;
  font-weight: 500;
}
</style>

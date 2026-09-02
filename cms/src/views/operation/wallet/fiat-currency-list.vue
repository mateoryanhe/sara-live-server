<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.FiatCurrencyManagement') }}</span>
        </div>
      </template>

      <div class="toolbar">
        <el-button v-if="can('create')" type="primary" @click="handleAdd">{{ t('pages.fiatCurrencyList.addCurrency') }}</el-button>
        <el-button v-if="can('reloadCfgCache')" @click="handleReloadCfgCache">{{ t('pages.fiatCurrencyList.reloadCfgCache') }}</el-button>
        <el-button v-if="can('reloadRateCache')" @click="handleReloadAllRateCache">{{ t('pages.fiatCurrencyList.reloadRateCache') }}</el-button>
      </div>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="t('pages.fiatCurrencyList.currencyCode')">
          <el-input v-model="searchForm.currencyCode" clearable :placeholder="t('pages.fiatCurrencyList.currencyCodePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.currencyName')">
          <el-input v-model="searchForm.name" clearable :placeholder="t('pages.fiatCurrencyList.currencyNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.currencyType')">
          <el-select v-model="searchForm.typeFilter" clearable style="width: 140px">
            <el-option :label="t('pages.fiatCurrencyList.allTypes')" :value="0"/>
            <el-option :label="t('pages.fiatCurrencyList.typeFiat')" :value="1"/>
            <el-option :label="t('pages.fiatCurrencyList.typeCrypto')" :value="2"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="searchForm.statusFilter" clearable style="width: 140px">
            <el-option :label="t('pages.fiatCurrencyList.allStatus')" :value="0"/>
            <el-option :label="t('pages.fiatCurrencyList.statusDisabled')" :value="1"/>
            <el-option :label="t('pages.fiatCurrencyList.statusEnabled')" :value="2"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="ID" prop="id" width="100"/>
        <el-table-column :label="t('pages.fiatCurrencyList.icon')" width="90">
          <template #default="{ row }">
            <el-image
                v-if="row.icon"
                :preview-src-list="[row.icon]"
                :src="row.icon"
                fit="contain"
                preview-teleported
                style="width: 40px; height: 40px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.fiatCurrencyList.currencyCode')" min-width="110" prop="currencyCode"/>
        <el-table-column :label="t('pages.fiatCurrencyList.currencyType')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.currencyType === 2 ? 'warning' : 'primary'">
              {{ row.currencyType === 2 ? t('pages.fiatCurrencyList.typeCrypto') : t('pages.fiatCurrencyList.typeFiat') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.fiatCurrencyList.currencyName')" min-width="140" prop="name"/>
        <el-table-column :label="t('pages.fiatCurrencyList.symbol')" min-width="90" prop="symbol"/>
        <el-table-column :label="t('pages.fiatCurrencyList.adjustPercent')" align="right" min-width="120" prop="adjustPercent"/>
        <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
        <el-table-column :label="t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? t('pages.fiatCurrencyList.statusEnabled') : t('pages.fiatCurrencyList.statusDisabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" prop="createdAt" width="170"/>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170"/>
        <el-table-column fixed="right" :label="t('common.actions')" width="260">
          <template #default="{ row }">
            <el-button v-if="can('previewRate')" link type="primary" @click="handlePreviewRate(row)">
              {{ t('pages.fiatCurrencyList.previewRate') }}
            </el-button>
            <el-button v-if="can('edit')" link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button v-if="can('delete')" link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="pagination.pageIndex"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchList"
            @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" destroy-on-close width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="140px">
        <el-form-item :label="t('pages.fiatCurrencyList.currencyCode')" prop="currencyCode">
          <el-input
              v-model="currentRow.currencyCode"
              :disabled="!!currentRow.id"
              :placeholder="t('pages.fiatCurrencyList.currencyCodePlaceholder')"
              maxlength="8"
              @input="currentRow.currencyCode = currentRow.currencyCode.toUpperCase()"
          />
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.currencyName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.fiatCurrencyList.currencyNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.symbol')" prop="symbol">
          <el-input v-model="currentRow.symbol" :placeholder="t('pages.fiatCurrencyList.symbolPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.currencyType')" prop="currencyType">
          <el-radio-group v-model="currentRow.currencyType">
            <el-radio :value="1">{{ t('pages.fiatCurrencyList.typeFiat') }}</el-radio>
            <el-radio :value="2">{{ t('pages.fiatCurrencyList.typeCrypto') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.icon')">
          <div class="icon-field">
            <el-upload
                :before-upload="beforeIconUpload"
                :disabled="iconUploading"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="icon-uploader"
            >
              <img v-if="iconPreviewUrl" :src="iconPreviewUrl" alt="icon" class="icon-preview"/>
              <div v-else class="icon-uploader-placeholder">
                <el-icon class="icon-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>{{ t('pages.fiatCurrencyList.clickUploadIcon') }}</span>
              </div>
            </el-upload>
            <el-button
                v-if="iconPreviewUrl || currentRow.icon"
                link
                type="danger"
                @click="clearIcon"
            >
              {{ t('pages.fiatCurrencyList.removeIcon') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.fiatCurrencyList.adjustPercent')" prop="adjustPercent">
          <el-input-number v-model="currentRow.adjustPercent" :precision="4" :step="0.1" controls-position="right"/>
          <div class="form-tip">{{ t('pages.fiatCurrencyList.adjustPercentTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">{{ t('pages.fiatCurrencyList.sortHigherFirst') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :value="1">{{ t('pages.fiatCurrencyList.statusEnabled') }}</el-radio>
            <el-radio :value="0">{{ t('pages.fiatCurrencyList.statusDisabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rateDialogVisible" :title="t('pages.fiatCurrencyList.ratePreviewTitle')" destroy-on-close width="520px">
      <el-descriptions v-if="ratePreview" :column="1" border>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.currencyCode')">{{ ratePreview.quote }}</el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.marketRate')">
          {{ t('pages.fiatCurrencyList.oneUsdEquals') }} {{ formatRate(ratePreview.marketRate) }} {{ ratePreview.quote }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.adjustPercent')">{{ ratePreview.adjustPercent }}%</el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.finalRate')">
          {{ t('pages.fiatCurrencyList.oneUsdEquals') }} {{ formatRate(ratePreview.rate) }} {{ ratePreview.quote }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.inverseRate')">
          1 {{ ratePreview.quote }} = {{ formatRate(ratePreview.inverseRate) }} {{ ratePreview.base }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.rateSource')">{{ ratePreview.source }}</el-descriptions-item>
        <el-descriptions-item :label="t('pages.fiatCurrencyList.rateDate')">{{ ratePreview.rateDate || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('common.status')">
          {{ ratePreview.cached ? t('pages.fiatCurrencyList.cached') : t('pages.fiatCurrencyList.liveFetched') }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button v-if="can('reloadRateCache')" @click="handleReloadRowRateCache">{{ t('pages.fiatCurrencyList.reloadRateCache') }}</el-button>
        <el-button type="primary" @click="rateDialogVisible = false">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {fiatCurrencyApi, uploadApi} from '@/api'
import type {FiatCurrency, FiatExchangeRate} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

interface CurrencyForm {
  id: string
  currencyCode: string
  name: string
  symbol: string
  icon: string
  adjustPercent: number
  currencyType: number
  sort: number
  status: number
}

const {t} = useI18n()
const {can} = usePagePermission('FiatCurrencyManagement')

const loading = ref(false)
const tableData = ref<FiatCurrency[]>([])
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})
const searchForm = reactive({
  currencyCode: '',
  name: '',
  typeFilter: 0,
  statusFilter: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): CurrencyForm => ({
  id: '',
  currencyCode: '',
  name: '',
  symbol: '',
  icon: '',
  adjustPercent: 0,
  currencyType: 1,
  sort: 0,
  status: 1,
})
const currentRow = ref<CurrencyForm>(defaultForm())
const formRef = ref<FormInstance>()

const iconUploading = ref(false)
const iconPreviewUrl = ref('')
let objectPreviewUrl = ''

const rateDialogVisible = ref(false)
const ratePreview = ref<FiatExchangeRate | null>(null)
const previewCurrencyCode = ref('')

const formRules = computed<FormRules>(() => ({
  currencyCode: [
    {required: true, message: t('pages.fiatCurrencyList.currencyCodeRequired'), trigger: 'blur'},
    {min: 3, max: 8, message: t('pages.fiatCurrencyList.currencyCodeLength'), trigger: 'blur'},
  ],
  name: [{required: true, message: t('pages.fiatCurrencyList.currencyNameRequired'), trigger: 'blur'}],
  symbol: [{required: true, message: t('pages.fiatCurrencyList.symbolRequired'), trigger: 'blur'}],
}))

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = ''
  }
}

const setIconPreview = (url: string, isObjectUrl = false) => {
  revokeObjectPreview()
  iconPreviewUrl.value = url
  if (isObjectUrl) {
    objectPreviewUrl = url
  }
}

const clearIcon = () => {
  currentRow.value.icon = ''
  setIconPreview('')
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    iconPreviewUrl.value = ''
  }
})

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.fiatCurrencyList.imageOnly'))
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  iconUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.icon = res.fileName
    setIconPreview(URL.createObjectURL(file), true)
    ElMessage.success(t('pages.fiatCurrencyList.uploadSuccess'))
  } catch (error) {
    console.error('upload fiat currency icon failed:', error)
    ElMessage.error(t('pages.fiatCurrencyList.uploadFailed'))
  } finally {
    iconUploading.value = false
  }
}

const formatRate = (value: number | null | undefined) => {
  if (value == null || Number.isNaN(value)) {
    return '-'
  }
  return new Intl.NumberFormat(undefined, {maximumFractionDigits: 6}).format(value)
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await fiatCurrencyApi.getList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      currencyCode: searchForm.currencyCode.trim(),
      name: searchForm.name.trim(),
      typeFilter: searchForm.typeFilter,
      statusFilter: searchForm.statusFilter,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch fiat currency list failed:', error)
    ElMessage.error(t('pages.fiatCurrencyList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.currencyCode = ''
  searchForm.name = ''
  searchForm.typeFilter = 0
  searchForm.statusFilter = 0
  pagination.pageIndex = 1
  fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.fiatCurrencyList.addCurrency')
  currentRow.value = defaultForm()
  setIconPreview('')
  dialogVisible.value = true
}

const handleEdit = (row: FiatCurrency) => {
  dialogTitle.value = t('pages.fiatCurrencyList.editCurrency')
  currentRow.value = {
    id: row.id,
    currencyCode: row.currencyCode,
    name: row.name,
    symbol: row.symbol,
    icon: row.iconName || '',
    adjustPercent: Number(row.adjustPercent) || 0,
    currencyType: Number(row.currencyType) === 2 ? 2 : 1,
    sort: Number(row.sort) || 0,
    status: Number(row.status) || 0,
  }
  setIconPreview(row.icon || '')
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    const payload = {
      currencyCode: currentRow.value.currencyCode.trim().toUpperCase(),
      name: currentRow.value.name.trim(),
      symbol: currentRow.value.symbol.trim(),
      icon: currentRow.value.icon.trim(),
      adjustPercent: currentRow.value.adjustPercent,
      currencyType: currentRow.value.currencyType,
      sort: currentRow.value.sort,
      status: currentRow.value.status,
    }
    if (currentRow.value.id) {
      await fiatCurrencyApi.update({id: currentRow.value.id, ...payload})
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await fiatCurrencyApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save fiat currency failed:', error)
    if (!(error && typeof error === 'object' && 'message' in error)) {
      ElMessage.error(t('pages.fiatCurrencyList.saveFailed'))
    }
  }
}

const handleDelete = async (row: FiatCurrency) => {
  try {
    await ElMessageBox.confirm(
        t('pages.fiatCurrencyList.deleteConfirm', {code: row.currencyCode}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    await fiatCurrencyApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete fiat currency failed:', error)
    }
  }
}

const handleReloadCfgCache = async () => {
  try {
    await fiatCurrencyApi.reloadCfgCache()
    ElMessage.success(t('pages.fiatCurrencyList.reloadCfgCacheSuccess'))
  } catch (error) {
    console.error('reload fiat currency cfg cache failed:', error)
  }
}

const handleReloadAllRateCache = async () => {
  try {
    await fiatCurrencyApi.reloadRateCache()
    ElMessage.success(t('pages.fiatCurrencyList.reloadRateCacheSuccess'))
  } catch (error) {
    console.error('reload all fiat exchange rate cache failed:', error)
  }
}

const handleReloadRowRateCache = async () => {
  if (!previewCurrencyCode.value) {
    return
  }
  try {
    await fiatCurrencyApi.reloadRateCache(previewCurrencyCode.value)
    ElMessage.success(t('pages.fiatCurrencyList.reloadRateCacheSuccess'))
    await loadRatePreview(previewCurrencyCode.value)
  } catch (error) {
    console.error('reload fiat exchange rate cache failed:', error)
  }
}

const loadRatePreview = async (currencyCode: string) => {
  const response = await fiatCurrencyApi.getExchangeRate(currencyCode)
  ratePreview.value = response
}

const handlePreviewRate = async (row: FiatCurrency) => {
  previewCurrencyCode.value = row.currencyCode
  rateDialogVisible.value = true
  ratePreview.value = null
  try {
    await loadRatePreview(row.currencyCode)
  } catch (error) {
    console.error('preview fiat exchange rate failed:', error)
    rateDialogVisible.value = false
    ElMessage.error(t('pages.fiatCurrencyList.previewRateFailed'))
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.page-container {
  padding: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-tip {
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}

.icon-field {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.icon-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.icon-uploader-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.icon-uploader-icon {
  margin-bottom: 6px;
  font-size: 24px;
}

.icon-preview {
  display: block;
  width: 96px;
  height: 96px;
  object-fit: contain;
}
</style>

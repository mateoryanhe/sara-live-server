<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.RechargeCfgManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.rechargeCfgList.addCfg') }}</el-button>
          <el-button
              v-if="hasButtonPermission('RechargeCfgManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">
          {{ t('common.selectedCount', {count: selectedRows.length}) }}
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.rechargeCfgList.tierName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.rechargeCfgList.nameFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('pages.rechargeCfgList.packageName')">
            <el-input v-model="searchForm.packageName" clearable :placeholder="t('pages.rechargeCfgList.packageNameFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('pages.rechargeCfgList.cfgType')">
            <el-select v-model="searchForm.typeFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option v-for="item in cfgTypeOptions" :key="item.value" :label="item.label" :value="item.value"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('common.status')">
            <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option :value="2" :label="t('common.onlyOnShelf')"/>
              <el-option :value="1" :label="t('common.onlyOffShelf')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table
            v-loading="loading"
            :data="tableData"
            row-key="id"
            style="width: 100%"
            @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="48"/>
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.rechargeCfgList.tierName')" prop="name" min-width="120"/>
          <el-table-column :label="t('pages.rechargeCfgList.packageName')" prop="packageName" min-width="160" show-overflow-tooltip/>
          <el-table-column :label="t('pages.rechargeCfgList.cfgType')" width="100">
            <template #default="{ row }">
              {{ cfgTypeLabel(row.cfgType) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.icon')" width="90">
            <template #default="{ row }">
              <el-image
                  v-if="row.icon"
                  :preview-src-list="[row.icon]"
                  :src="row.icon"
                  fit="cover"
                  preview-teleported
                  style="width: 48px; height: 48px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.rechargeCfgList.baseGold')" prop="gold" width="100"/>
          <el-table-column :label="t('pages.rechargeCfgList.extraGold')" prop="extraGold" width="100"/>
          <el-table-column :label="t('pages.rechargeCfgList.totalGold')" width="100">
            <template #default="{ row }">
              {{ totalGold(row) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.price')" width="120">
            <template #default="{ row }">
              {{ formatPrice(row.price) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.rechargeCfgList.productSku')" prop="productId" min-width="140" show-overflow-tooltip/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.description')" prop="description" min-width="160" show-overflow-tooltip/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="260">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="row.status !== 1"
                  size="small"
                  type="success"
                  @click="handleOnShelf(row)"
              >
                {{ t('common.onShelf') }}
              </el-button>
              <el-button
                  v-else
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                {{ t('common.offShelf') }}
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
        <el-form-item :label="t('pages.rechargeCfgList.tierName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.rechargeCfgList.tierNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.rechargeCfgList.packageName')" prop="packageName">
          <el-input v-model="currentRow.packageName" :placeholder="t('pages.rechargeCfgList.packageNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.rechargeCfgList.cfgType')" prop="cfgType">
          <el-select v-model="currentRow.cfgType" :placeholder="t('pages.rechargeCfgList.selectType')" style="width: 100%">
            <el-option v-for="item in cfgTypeOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.icon')" prop="icon">
          <div class="icon-upload-wrap">
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
                <span>{{ t('pages.rechargeCfgList.clickUploadIcon') }}</span>
              </div>
            </el-upload>
            <el-button
                v-if="iconPreviewUrl || currentRow.icon"
                link
                type="danger"
                @click="clearIcon"
            >
              {{ t('pages.rechargeCfgList.removeIcon') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.rechargeCfgList.baseGold')" prop="gold">
          <el-input-number v-model="currentRow.gold" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.rechargeCfgList.extraGold')" prop="extraGold">
          <el-input-number v-model="currentRow.extraGold" :min="0" controls-position="right"/>
          <div class="form-tip">{{ t('pages.rechargeCfgList.goldCreditTip') }}</div>
        </el-form-item>
        <el-form-item :label="t('common.price')" prop="price">
          <el-input-number v-model="currentRow.price" :min="0.0001" :precision="NUMBER_INPUT_DECIMALS" :step="0.01" controls-position="right"/>
          <div class="form-tip">{{ t('pages.rechargeCfgList.priceTip') }}</div>
        </el-form-item>
        <el-form-item
            :required="isStoreSkuRequired(currentRow.cfgType)"
            :label="t('pages.rechargeCfgList.productSku')"
            prop="productId"
        >
          <el-input
              v-model="currentRow.productId"
              :placeholder="productSkuPlaceholder"
          />
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('common.description')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="descriptionPlaceholder" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {dataSyncApi, rechargeCfgApi, uploadApi} from '@/api'
import type {RechargeCfg} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'
import {formatPrice, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  name: string
  packageName: string
  typeFilter: number
  statusFilter: number
}

interface RechargeCfgForm {
  id: string
  name: string
  packageName: string
  cfgType: number
  icon: string
  gold: number
  extraGold: number
  price: number
  productId: string
  sort: number
  description: string
}

const {t} = useI18n()

const cfgTypeOptions = computed(() => [
  {value: 1, label: t('pages.rechargeCfgList.typeIos')},
  {value: 2, label: t('pages.rechargeCfgList.typeGoogle')},
  {value: 3, label: t('pages.rechargeCfgList.typeChannel')}
])

const cfgTypeLabel = (cfgType: number) => {
  return cfgTypeOptions.value.find((item) => item.value === cfgType)?.label ?? t('pages.rechargeCfgList.unknown')
}

const isStoreSkuRequired = (cfgType: number) => cfgType === 1 || cfgType === 2

const productSkuPlaceholder = computed(() =>
    isStoreSkuRequired(currentRow.value.cfgType)
        ? t('pages.rechargeCfgList.productSkuRequired')
        : t('pages.rechargeCfgList.productSkuOptional')
)

const descriptionPlaceholder = computed(() => `${t('common.pleaseEnter')}${t('common.description')}`)

const loading = ref(false)
const syncing = ref(false)
const selectedRows = ref<RechargeCfg[]>([])
const tableData = ref<RechargeCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  packageName: '',
  typeFilter: 0,
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): RechargeCfgForm => ({
  id: '',
  name: '',
  packageName: '',
  cfgType: 1,
  icon: '',
  gold: 1,
  extraGold: 0,
  price: 0.99,
  productId: '',
  sort: 0,
  description: ''
})
const currentRow = ref<RechargeCfgForm>(defaultForm())
const formRef = ref<FormInstance>()
const iconUploading = ref(false)
const iconPreviewUrl = ref('')
let objectPreviewUrl: string | null = null

const totalGold = (row: RechargeCfg) => {
  return (Number(row.gold) || 0) + (Number(row.extraGold) || 0)
}

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = null
  }
}

const setIconPreview = (url: string, fromObject = false) => {
  revokeObjectPreview()
  iconPreviewUrl.value = url
  if (fromObject) {
    objectPreviewUrl = url
  }
}

const clearIcon = () => {
  currentRow.value.icon = ''
  setIconPreview('')
  formRef.value?.validateField('icon').catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    iconPreviewUrl.value = ''
  }
})

watch(() => currentRow.value.cfgType, () => {
  formRef.value?.validateField('productId').catch(() => undefined)
})

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.rechargeCfgList.imageOnly'))
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
    ElMessage.success(t('pages.rechargeCfgList.uploadSuccess'))
  } catch (error) {
    console.error('upload failed:', error)
    ElMessage.error(t('pages.rechargeCfgList.uploadFailed'))
  } finally {
    iconUploading.value = false
  }
}

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.rechargeCfgList.nameRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.rechargeCfgList.nameLength'), trigger: 'blur'}
  ],
  packageName: [
    {required: true, message: t('pages.rechargeCfgList.packageNameRequired'), trigger: 'blur'},
    {min: 1, max: 128, message: t('pages.rechargeCfgList.packageNameLength'), trigger: 'blur'}
  ],
  cfgType: [{required: true, message: t('pages.rechargeCfgList.typeRequired'), trigger: 'change'}],
  gold: [{required: true, message: t('pages.rechargeCfgList.goldRequired'), trigger: 'change'}],
  price: [{required: true, message: t('pages.rechargeCfgList.priceRequired'), trigger: 'change'}],
  productId: [
    {
      validator: (_rule, value, callback) => {
        if (!isStoreSkuRequired(currentRow.value.cfgType)) {
          callback()
          return
        }
        if (!String(value ?? '').trim()) {
          callback(new Error(t('pages.rechargeCfgList.skuRequired')))
          return
        }
        if (String(value).length > 64) {
          callback(new Error(t('pages.rechargeCfgList.skuMaxLength')))
          return
        }
        callback()
      },
      trigger: 'blur'
    }
  ],
  description: [{max: 255, message: t('pages.rechargeCfgList.descriptionMaxLength'), trigger: 'blur'}]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await rechargeCfgApi.getRechargeCfgList({
      name: searchForm.name,
      packageName: searchForm.packageName,
      typeFilter: searchForm.typeFilter,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch recharge cfg list failed:', error)
    ElMessage.error(t('pages.rechargeCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.rechargeCfgList.addCfg')
  currentRow.value = defaultForm()
  setIconPreview('')
  dialogVisible.value = true
}

const handleEdit = (row: RechargeCfg) => {
  dialogTitle.value = t('pages.rechargeCfgList.editCfg')
  currentRow.value = {
    id: row.id,
    name: row.name,
    packageName: row.packageName || '',
    cfgType: Number(row.cfgType) || 1,
    icon: row.iconName || '',
    gold: Number(row.gold) || 1,
    extraGold: Number(row.extraGold) || 0,
    price: truncateNumber(row.price) || 0.99,
    productId: row.productId || '',
    sort: Number(row.sort) || 0,
    description: row.description || ''
  }
  setIconPreview(row.icon || '')
  dialogVisible.value = true
}

const handleDelete = async (row: RechargeCfg) => {
  try {
    await ElMessageBox.confirm(
        t('pages.rechargeCfgList.deleteConfirm', {name: row.name}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await rechargeCfgApi.deleteRechargeCfg(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete recharge cfg failed:', error)
    }
  }
}

const handleOnShelf = async (row: RechargeCfg) => {
  try {
    await rechargeCfgApi.onShelfRechargeCfg(row.id)
    ElMessage.success(t('pages.rechargeCfgList.onShelfSuccess'))
    fetchList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.rechargeCfgList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: RechargeCfg) => {
  try {
    await ElMessageBox.confirm(
        t('pages.rechargeCfgList.offShelfConfirm', {name: row.name}),
        t('common.confirmOffShelf'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await rechargeCfgApi.offShelfRechargeCfg(row.id)
    ElMessage.success(t('pages.rechargeCfgList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('off shelf failed:', error)
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        name: currentRow.value.name,
        packageName: currentRow.value.packageName,
        cfgType: currentRow.value.cfgType,
        icon: currentRow.value.icon,
        gold: currentRow.value.gold,
        extraGold: currentRow.value.extraGold,
        price: currentRow.value.price,
        productId: currentRow.value.productId,
        sort: currentRow.value.sort,
        description: currentRow.value.description
      }
      if (currentRow.value.id) {
        await rechargeCfgApi.updateRechargeCfg({id: currentRow.value.id, ...payload})
      } else {
        await rechargeCfgApi.createRechargeCfg(payload)
      }
      ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.rechargeCfgList.saveFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.packageName = ''
  searchForm.typeFilter = 0
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
}

const handleSelectionChange = (rows: RechargeCfg[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.rechargeCfgList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.rechargeCfgList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.rechargeCfgList.syncConfirm', {count: ids.length}),
        t('common.syncData'),
        {confirmButtonText: t('common.confirmSync'), cancelButtonText: t('common.cancel'), type: 'warning'}
    )
    syncing.value = true
    const response = await dataSyncApi.syncRechargeCfg({ids})
    if (response?.success) {
      ElMessage.success(
          response.message || t('pages.rechargeCfgList.syncSuccessDetail', {
            rows: response.rowCount,
            files: response.fileCount
          })
      )
    } else {
      ElMessage.error(t('pages.rechargeCfgList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync failed:', error)
      ElMessage.error(t('pages.rechargeCfgList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

onMounted(() => {
  fetchList()
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

.table-header {
  margin-bottom: 20px;
}

.selection-tip {
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--el-color-primary);
}

.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.icon-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.icon-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.icon-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.icon-uploader-placeholder {
  width: 96px;
  height: 96px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.icon-uploader-icon {
  font-size: 28px;
}

.icon-preview {
  width: 96px;
  height: 96px;
  display: block;
  object-fit: cover;
}
</style>

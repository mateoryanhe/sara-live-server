<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GiftManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.giftList.addGift') }}</el-button>
          <el-button
              v-if="hasButtonPermission('GiftManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
          <el-button
              v-if="hasButtonPermission('GiftManagement', 'syncAssets')"
              :disabled="selectedRows.length === 0"
              :loading="syncingAssets"
              type="warning"
              plain
              @click="handleSyncAssets"
          >
            {{ t('pages.giftList.syncAssets') }}
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">
          {{ t('common.selectedCount', {count: selectedRows.length}) }}
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.giftList.giftName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.giftList.giftNameFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('common.category')">
            <el-input v-model="searchForm.category" clearable :placeholder="t('common.category')"/>
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
          <el-table-column :label="t('pages.giftList.giftName')" prop="name" min-width="120"/>
          <el-table-column :label="t('common.nameEn')" prop="nameEn" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameEs')" prop="nameEs" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.namePt')" prop="namePt" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameHi')" prop="nameHi" min-width="120" show-overflow-tooltip/>
          <el-table-column :label="t('common.nameId')" prop="nameId" min-width="120" show-overflow-tooltip/>
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
          <el-table-column :label="t('common.animationResource')" min-width="220">
            <template #default="{ row }">
              <video
                  v-if="isVideoUrl(row.animation)"
                  :src="row.animation"
                  class="table-media-preview"
                  controls
                  preload="metadata"
              />
              <el-image
                  v-else-if="isImageUrl(row.animation)"
                  :preview-src-list="[row.animation]"
                  :src="row.animation"
                  fit="cover"
                  preview-teleported
                  style="width: 48px; height: 48px"
              />
              <span v-else class="media-url-text">{{ row.animationName || row.animation || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.price')" prop="price" width="110">
            <template #default="{ row }">{{ formatPrice(row.price) }}</template>
          </el-table-column>
          <el-table-column :label="t('common.category')" prop="category" width="120"/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.description')" prop="description" min-width="160" show-overflow-tooltip/>
          <el-table-column :label="t('common.publishedAt')" prop="publishedAt" width="160">
            <template #default="{ row }">
              {{ row.publishedAt || '-' }}
            </template>
          </el-table-column>
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
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.giftList.giftName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.giftList.giftNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.nameEn')" prop="nameEn">
          <el-input v-model="currentRow.nameEn" :placeholder="t('pages.giftList.nameEnPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.nameEs')" prop="nameEs">
          <el-input v-model="currentRow.nameEs" :placeholder="t('pages.giftList.nameEsPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.namePt')" prop="namePt">
          <el-input v-model="currentRow.namePt" :placeholder="t('pages.giftList.namePtPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.nameHi')" prop="nameHi">
          <el-input v-model="currentRow.nameHi" :placeholder="t('pages.giftList.nameHiPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.nameId')" prop="nameId">
          <el-input v-model="currentRow.nameId" :placeholder="t('pages.giftList.nameIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.icon')" prop="icon">
          <div class="asset-upload-wrap">
            <el-upload
                :before-upload="beforeIconUpload"
                :disabled="iconUploading"
                :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'icon')"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="icon-uploader"
            >
              <img v-if="iconPreviewUrl" :src="iconPreviewUrl" alt="icon" class="icon-preview"/>
              <div v-else class="asset-uploader-placeholder icon-placeholder">
                <el-icon class="asset-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>{{ t('pages.giftList.clickUploadIcon') }}</span>
              </div>
            </el-upload>
            <el-button
                v-if="iconPreviewUrl || currentRow.icon"
                link
                type="danger"
                @click="clearAsset('icon')"
            >
              {{ t('pages.giftList.removeIcon') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('common.animationResource')" prop="animation">
          <div class="asset-upload-wrap">
            <el-upload
                :before-upload="beforeAnimationUpload"
                :disabled="animationUploading"
                :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'animation')"
                :show-file-list="false"
                accept=".svga,.pag,.json,.lottie,.mp4,.webm,.zip,.gif,.apng,.png,.webp,.jpg,.jpeg,.bmp"
                action="#"
                class="animation-uploader"
            >
              <img
                  v-if="animationPreviewType === 'image'"
                  :src="animationPreviewUrl"
                  alt="animation"
                  class="animation-preview"
              />
              <video
                  v-else-if="animationPreviewType === 'video'"
                  :src="animationPreviewUrl"
                  class="animation-preview"
                  controls
                  preload="metadata"
              />
              <div v-else-if="animationPreviewType === 'file' && animationFileLabel" class="asset-file-label">
                <el-icon class="asset-uploader-icon">
                  <Document/>
                </el-icon>
                <span class="file-name">{{ animationFileLabel }}</span>
              </div>
              <div v-else class="asset-uploader-placeholder animation-placeholder">
                <el-icon class="asset-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>{{ t('pages.giftList.clickUploadAnimation') }}</span>
              </div>
            </el-upload>
            <el-button
                v-if="animationPreviewType !== 'none' || currentRow.animation"
                link
                type="danger"
                @click="clearAsset('animation')"
            >
              {{ t('pages.giftList.removeAnimation') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.giftList.priceDiamond')" prop="price">
          <el-input-number v-model="currentRow.price" :min="0" :precision="NUMBER_INPUT_DECIMALS" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('common.category')" prop="category">
          <el-input v-model="currentRow.category" :placeholder="t('pages.giftList.categoryPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('common.description')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="t('pages.giftList.descriptionPlaceholder')" type="textarea"/>
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Document, Plus} from '@element-plus/icons-vue'
import {dataSyncApi, giftApi, uploadApi} from '@/api'
import type {Gift} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'
import {
  getExt,
  isImageUrl,
  isVideoUrl,
  resolveFilePreviewType,
  resolveMediaPreviewType,
  type MediaPreviewType
} from '@/utils/media-preview'
import {formatPrice, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  name: string
  category: string
  statusFilter: number
}

interface GiftForm {
  id: string
  name: string
  nameEn: string
  nameEs: string
  namePt: string
  nameHi: string
  nameId: string
  icon: string
  animation: string
  price: number
  category: string
  sort: number
  description: string
}

const {t} = useI18n()

const loading = ref(false)
const syncing = ref(false)
const syncingAssets = ref(false)
const selectedRows = ref<Gift[]>([])
const tableData = ref<Gift[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  category: '',
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): GiftForm => ({
  id: '',
  name: '',
  nameEn: '',
  nameEs: '',
  namePt: '',
  nameHi: '',
  nameId: '',
  icon: '',
  animation: '',
  price: 0,
  category: '',
  sort: 0,
  description: ''
})
const currentRow = ref<GiftForm>(defaultForm())

const formRef = ref<FormInstance>()

const iconUploading = ref(false)
const animationUploading = ref(false)
const iconPreviewUrl = ref('')
const animationPreviewUrl = ref('')
const animationPreviewType = ref<MediaPreviewType>('none')
const animationFileLabel = ref('')
const objectPreviewUrls: Partial<Record<'icon' | 'animation', string>> = {}

const revokeObjectPreview = (field: 'icon' | 'animation') => {
  const url = objectPreviewUrls[field]
  if (url) {
    URL.revokeObjectURL(url)
    delete objectPreviewUrls[field]
  }
}

const revokeAllObjectPreviews = () => {
  revokeObjectPreview('icon')
  revokeObjectPreview('animation')
}

const resetAssetPreview = () => {
  revokeAllObjectPreviews()
  iconPreviewUrl.value = ''
  animationPreviewUrl.value = ''
  animationPreviewType.value = 'none'
  animationFileLabel.value = ''
}

const setIconPreview = (url: string, fromObject = false) => {
  revokeObjectPreview('icon')
  iconPreviewUrl.value = url
  if (fromObject && url) {
    objectPreviewUrls.icon = url
  }
}

const setAnimationPreview = (
    url: string,
    fileLabel: string,
    type?: MediaPreviewType,
    fromObject = false
) => {
  revokeObjectPreview('animation')
  const previewType = type ?? resolveMediaPreviewType(url, fileLabel)
  animationPreviewType.value = previewType
  animationPreviewUrl.value = previewType === 'image' || previewType === 'video' ? url : ''
  animationFileLabel.value = previewType === 'file' ? fileLabel : ''
  if (fromObject && url) {
    objectPreviewUrls.animation = url
  }
}

const clearAsset = (field: 'icon' | 'animation') => {
  currentRow.value[field] = ''
  if (field === 'icon') {
    setIconPreview('')
  } else {
    setAnimationPreview('', '', 'none')
  }
  formRef.value?.validateField(field).catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAssetPreview()
  }
})

const allowedAnimationExt = [
  '.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.apng',
  '.svga', '.pag', '.json', '.lottie', '.mp4', '.webm', '.zip'
]

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.giftList.iconImageOnly'))
    return false
  }
  return true
}

const beforeAnimationUpload = (file: File): boolean => {
  const ext = getExt(file.name)
  if (!allowedAnimationExt.includes(ext)) {
    ElMessage.error(t('pages.giftList.unsupportedFileType', {ext: ext || '-'}))
    return false
  }
  return true
}

const doUpload = async (
    options: UploadRequestOptions,
    field: 'icon' | 'animation'
) => {
  const file = options.file as File
  const flag = field === 'icon' ? iconUploading : animationUploading
  flag.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value[field] = res.fileName
    if (field === 'icon') {
      setIconPreview(URL.createObjectURL(file), true)
    } else {
      const previewType = resolveFilePreviewType(file.name)
      if (previewType === 'image' || previewType === 'video') {
        setAnimationPreview(URL.createObjectURL(file), '', previewType, true)
      } else {
        setAnimationPreview('', res.fileName, 'file')
      }
    }
    formRef.value?.validateField(field).catch(() => undefined)
    ElMessage.success(t('pages.giftList.uploadSuccess'))
  } catch (error) {
    console.error('upload failed:', error)
    ElMessage.error(t('pages.giftList.uploadFailed'))
  } finally {
    flag.value = false
  }
}

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.giftList.nameRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.giftList.nameLength'), trigger: 'blur'}
  ],
  nameEn: [{max: 64, message: t('pages.giftList.nameEnMaxLength'), trigger: 'blur'}],
  nameEs: [{max: 64, message: t('pages.giftList.nameEsMaxLength'), trigger: 'blur'}],
  namePt: [{max: 64, message: t('pages.giftList.namePtMaxLength'), trigger: 'blur'}],
  nameHi: [{max: 64, message: t('pages.giftList.nameHiMaxLength'), trigger: 'blur'}],
  nameId: [{max: 64, message: t('pages.giftList.nameIdMaxLength'), trigger: 'blur'}],
  icon: [],
  animation: [],
  category: [
    {max: 32, message: t('pages.giftList.categoryMaxLength'), trigger: 'blur'}
  ],
  description: [
    {max: 255, message: t('pages.giftList.descriptionMaxLength'), trigger: 'blur'}
  ]
}))

const fetchGiftList = async () => {
  loading.value = true
  try {
    const response = await giftApi.getGiftList({
      name: searchForm.name,
      category: searchForm.category,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch gift list failed:', error)
    ElMessage.error(t('pages.giftList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchGiftList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGiftList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGiftList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.giftList.addGift')
  currentRow.value = defaultForm()
  resetAssetPreview()
  dialogVisible.value = true
}

const handleEdit = (row: Gift) => {
  dialogTitle.value = t('pages.giftList.editGift')
  const iconName = row.iconName || ''
  const animationName = row.animationName || ''
  currentRow.value = {
    id: row.id,
    name: row.name,
    nameEn: row.nameEn || '',
    nameEs: row.nameEs || '',
    namePt: row.namePt || '',
    nameHi: row.nameHi || '',
    nameId: row.nameId || '',
    icon: iconName,
    animation: animationName,
    price: truncateNumber(row.price),
    category: row.category,
    sort: Number(row.sort) || 0,
    description: row.description
  }
  setIconPreview(row.icon || '')
  if (animationName) {
    const previewType = resolveMediaPreviewType(row.animation || '', animationName)
    if (previewType === 'image' || previewType === 'video') {
      setAnimationPreview(row.animation || '', '', previewType)
    } else {
      setAnimationPreview('', animationName, 'file')
    }
  } else {
    setAnimationPreview('', '', 'none')
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.deleteConfirm', {name: row.name}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )

    await giftApi.deleteGift(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchGiftList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete gift failed:', error)
    }
  }
}

const handleOnShelf = async (row: Gift) => {
  try {
    await giftApi.onShelfGift(row.id)
    ElMessage.success(t('pages.giftList.onShelfSuccess'))
    fetchGiftList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.giftList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.offShelfConfirm', {name: row.name}),
        t('common.confirmOffShelf'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await giftApi.offShelfGift(row.id)
    ElMessage.success(t('pages.giftList.offShelfSuccess'))
    fetchGiftList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('off shelf failed:', error)
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentRow.value.id) {
          const {id, name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, category, sort, description} = currentRow.value
          await giftApi.updateGift({id, name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, category, sort, description})
        } else {
          const {name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, category, sort, description} = currentRow.value
          await giftApi.createGift({name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, category, sort, description})
        }

        ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
        dialogVisible.value = false
        fetchGiftList()
      } catch (error) {
        console.error(currentRow.value.id ? 'update failed:' : 'create failed:', error)
        ElMessage.error(currentRow.value.id ? t('pages.giftList.updateFailed') : t('pages.giftList.createFailed'))
      }
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.category = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchGiftList()
}

const handleSelectionChange = (rows: Gift[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.giftList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.giftList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.syncConfirm', {count: ids.length}),
        t('common.syncData'),
        {confirmButtonText: t('common.confirmSync'), cancelButtonText: t('common.cancel'), type: 'warning'}
    )
    syncing.value = true
    const response = await dataSyncApi.syncGift({ids})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.giftList.syncSuccessDetail', {
        rows: response.rowCount,
        files: response.fileCount
      }))
    } else {
      ElMessage.error(t('pages.giftList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync failed:', error)
      ElMessage.error(t('pages.giftList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

const handleSyncAssets = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.giftList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.giftList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.giftList.syncAssetsConfirm', {count: ids.length}),
        t('pages.giftList.syncAssets'),
        {confirmButtonText: t('common.confirmSync'), cancelButtonText: t('common.cancel'), type: 'warning'}
    )
    syncingAssets.value = true
    const response = await dataSyncApi.syncGiftAssets({ids})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.giftList.syncAssetsSuccessDetail', {
        files: response.fileCount
      }))
    } else {
      ElMessage.error(t('pages.giftList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync gift assets failed:', error)
      ElMessage.error(t('pages.giftList.syncFailedCheckConfig'))
    }
  } finally {
    syncingAssets.value = false
  }
}

onMounted(() => {
  fetchGiftList()
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

.asset-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.icon-uploader :deep(.el-upload),
.animation-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.icon-uploader :deep(.el-upload:hover),
.animation-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.asset-uploader-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.icon-placeholder {
  width: 96px;
  height: 96px;
}

.animation-placeholder {
  width: 240px;
  height: 120px;
}

.asset-uploader-icon {
  font-size: 28px;
}

.icon-preview {
  width: 96px;
  height: 96px;
  display: block;
  object-fit: cover;
}

.animation-preview {
  width: 240px;
  height: 120px;
  display: block;
  object-fit: cover;
}

.table-media-preview {
  width: 160px;
  max-height: 90px;
  display: block;
  background: #000;
  border-radius: 4px;
}

.media-url-text {
  word-break: break-all;
  line-height: 1.4;
}

.asset-file-label {
  width: 240px;
  height: 120px;
  padding: 12px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-align: center;
}

.asset-file-label .file-name {
  font-size: 13px;
  color: var(--el-text-color-regular);
  word-break: break-all;
  line-height: 1.4;
}
</style>

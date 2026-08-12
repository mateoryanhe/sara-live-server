<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.BannerManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.bannerList.addBanner') }}</el-button>
          <el-button
              v-if="hasButtonPermission('BannerManagement', 'sync')"
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
          <el-form-item :label="t('common.title')">
            <el-input v-model="searchForm.title" clearable :placeholder="t('pages.bannerList.titleFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('pages.bannerList.displayScene')">
            <el-select v-model="searchForm.sceneFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option
                  v-for="item in sceneOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
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
          <el-table-column :label="t('common.title')" prop="title" min-width="140"/>
          <el-table-column :label="t('pages.bannerList.image')" width="100">
            <template #default="{ row }">
              <el-image
                  v-if="row.image"
                  :preview-src-list="[row.image]"
                  :src="row.image"
                  fit="cover"
                  preview-teleported
                  style="width: 72px; height: 40px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.bannerList.link')" prop="link" min-width="200" show-overflow-tooltip/>
          <el-table-column :label="t('pages.bannerList.displayScene')" width="100">
            <template #default="{ row }">
              {{ sceneLabel(row.scene) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.bannerList.displayPosition')" width="110">
            <template #default="{ row }">
              {{ directionLabel(row.direction, row.scene) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('common.title')" prop="title">
          <el-input v-model="currentRow.title" :placeholder="t('pages.bannerList.titleRequired')"/>
        </el-form-item>
        <el-form-item :label="t('pages.bannerList.image')" prop="image">
          <div class="image-upload-wrap">
            <el-upload
                :before-upload="beforeImageUpload"
                :disabled="imageUploading"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="banner-uploader"
            >
              <img v-if="imagePreviewUrl" :src="imagePreviewUrl" alt="banner" class="banner-preview"/>
              <div v-else class="banner-uploader-placeholder">
                <el-icon class="banner-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>{{ t('pages.bannerList.clickUploadImage') }}</span>
              </div>
            </el-upload>
            <el-button
                v-if="imagePreviewUrl"
                class="banner-remove-btn"
                link
                type="danger"
                @click="clearImage"
            >
              {{ t('pages.bannerList.removeImage') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.bannerList.link')" prop="link">
          <el-input v-model="currentRow.link" :placeholder="t('pages.bannerList.linkPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.bannerList.displayScene')" prop="scene">
          <el-radio-group v-model="currentRow.scene">
            <el-radio v-for="item in sceneOptions" :key="item.value" :label="item.value">
              {{ item.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('pages.bannerList.displayPosition')" prop="direction">
          <el-select
              v-model="currentRow.direction"
              :placeholder="t('pages.bannerList.selectDisplayPosition')"
              style="width: 100%"
          >
            <el-option
                v-for="item in currentDirectionOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
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
import {Plus} from '@element-plus/icons-vue'
import {bannerApi, dataSyncApi, uploadApi} from '@/api'
import type {Banner} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'

const {t} = useI18n()

interface SearchForm {
  title: string
  sceneFilter: number
  statusFilter: number
}

interface BannerForm {
  id: string
  title: string
  image: string
  link: string
  scene: number
  direction: number
  sort: number
}

const homeDirectionOptions = computed(() => [
  {value: 1, label: t('pages.bannerList.positionHomeTop')},
  {value: 2, label: t('pages.bannerList.positionHomeMiddle')},
  {value: 3, label: t('pages.bannerList.positionHomeBottom')},
])

const liveRoomDirectionOptions = computed(() => [
  {value: 4, label: t('pages.bannerList.positionLiveHall')},
])

const sceneOptions = computed(() => [
  {value: 1, label: t('pages.bannerList.sceneHome')},
  {value: 2, label: t('pages.bannerList.sceneLiveRoom')},
])

const sceneLabel = (scene: number) => {
  return sceneOptions.value.find((item) => item.value === Number(scene))?.label ?? t('pages.bannerList.sceneHome')
}

const directionLabel = (direction: number, scene = 1) => {
  const options = Number(scene) === 2 ? liveRoomDirectionOptions.value : homeDirectionOptions.value
  return options.find((item) => item.value === direction)?.label ?? t('pages.bannerList.unknown')
}

const currentDirectionOptions = computed(() => {
  return currentRow.value.scene === 2 ? liveRoomDirectionOptions.value : homeDirectionOptions.value
})

const loading = ref(false)
const syncing = ref(false)
const selectedRows = ref<Banner[]>([])
const tableData = ref<Banner[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  title: '',
  sceneFilter: 0,
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): BannerForm => ({
  id: '',
  title: '',
  image: '',
  link: '',
  scene: 1,
  direction: 1,
  sort: 0
})
const currentRow = ref<BannerForm>(defaultForm())
const formRef = ref<FormInstance>()
const imageUploading = ref(false)
const imagePreviewUrl = ref('')
let objectPreviewUrl: string | null = null

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = null
  }
}

const setImagePreview = (url: string, fromObject = false) => {
  revokeObjectPreview()
  imagePreviewUrl.value = url
  if (fromObject) {
    objectPreviewUrl = url
  }
}

const clearImage = () => {
  currentRow.value.image = ''
  setImagePreview('')
  formRef.value?.validateField('image').catch(() => undefined)
}

watch(() => currentRow.value.scene, (scene) => {
  if (scene === 2) {
    currentRow.value.direction = 4
    return
  }
  if (![1, 2, 3].includes(currentRow.value.direction)) {
    currentRow.value.direction = 1
  }
})

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    imagePreviewUrl.value = ''
  }
})

const beforeImageUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.bannerList.imageOnly'))
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  imageUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.image = res.fileName
    setImagePreview(URL.createObjectURL(file), true)
    formRef.value?.validateField('image').catch(() => undefined)
    ElMessage.success(t('pages.bannerList.uploadSuccess'))
  } catch (error) {
    console.error('upload banner image failed:', error)
    ElMessage.error(t('pages.bannerList.uploadFailed'))
  } finally {
    imageUploading.value = false
  }
}

const formRules = computed<FormRules>(() => ({
  title: [
    {required: true, message: t('pages.bannerList.titleRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.bannerList.titleLength'), trigger: 'blur'}
  ],
  image: [{required: true, message: t('pages.bannerList.imageRequired'), trigger: 'change'}],
  link: [{max: 512, message: t('pages.bannerList.linkMaxLength'), trigger: 'blur'}],
  scene: [{required: true, message: t('pages.bannerList.sceneRequired'), trigger: 'change'}],
  direction: [{required: true, message: t('pages.bannerList.positionRequired'), trigger: 'change'}]
}))

const fetchBannerList = async () => {
  loading.value = true
  try {
    const response = await bannerApi.getBannerList({
      title: searchForm.title,
      sceneFilter: searchForm.sceneFilter,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch banner list failed:', error)
    ElMessage.error(t('pages.bannerList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchBannerList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchBannerList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchBannerList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.bannerList.addBanner')
  currentRow.value = defaultForm()
  setImagePreview('')
  dialogVisible.value = true
}

const handleEdit = (row: Banner) => {
  dialogTitle.value = t('pages.bannerList.editBanner')
  currentRow.value = {
    id: row.id,
    title: row.title,
    image: row.imageName || '',
    link: row.link,
    scene: Number(row.scene) || 1,
    direction: Number(row.direction) || 1,
    sort: Number(row.sort) || 0
  }
  setImagePreview(row.image || '')
  dialogVisible.value = true
}

const handleDelete = async (row: Banner) => {
  try {
    await ElMessageBox.confirm(
        t('pages.bannerList.deleteConfirm', {title: row.title}),
        t('common.confirmDelete'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await bannerApi.deleteBanner(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchBannerList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete banner failed:', error)
    }
  }
}

const handleOnShelf = async (row: Banner) => {
  try {
    await bannerApi.onShelfBanner(row.id)
    ElMessage.success(t('pages.bannerList.onShelfSuccess'))
    fetchBannerList()
  } catch (error) {
    console.error('publish banner failed:', error)
    ElMessage.error(t('pages.bannerList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: Banner) => {
  try {
    await ElMessageBox.confirm(
        t('pages.bannerList.offShelfConfirm', {title: row.title}),
        t('common.confirmOffShelf'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await bannerApi.offShelfBanner(row.id)
    ElMessage.success(t('pages.bannerList.offShelfSuccess'))
    fetchBannerList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('unpublish banner failed:', error)
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (currentRow.value.id) {
        await bannerApi.updateBanner(currentRow.value)
      } else {
        const {title, image, link, scene, direction, sort} = currentRow.value
        await bannerApi.createBanner({title, image, link, scene, direction, sort})
      }
      ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchBannerList()
    } catch (error) {
      console.error('save banner failed:', error)
      ElMessage.error(t('pages.bannerList.saveFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.title = ''
  searchForm.sceneFilter = 0
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchBannerList()
}

const handleSelectionChange = (rows: Banner[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.bannerList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.bannerList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.bannerList.syncConfirm', {count: ids.length}),
        t('common.syncData'),
        {confirmButtonText: t('common.confirmSync'), cancelButtonText: t('common.cancel'), type: 'warning'}
    )
    syncing.value = true
    const response = await dataSyncApi.syncBanner({ids})
    if (response?.success) {
      ElMessage.success(
          response.message || t('pages.bannerList.syncSuccessDetail', {
            rows: response.rowCount,
            files: response.fileCount
          })
      )
    } else {
      ElMessage.error(t('pages.bannerList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync banner failed:', error)
      ElMessage.error(t('pages.bannerList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

onMounted(() => {
  fetchBannerList()
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

.image-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.banner-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.banner-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.banner-uploader-placeholder {
  width: 240px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.banner-uploader-icon {
  font-size: 28px;
}

.banner-preview {
  width: 240px;
  height: 120px;
  display: block;
  object-fit: cover;
}

.banner-remove-btn {
  padding: 0;
}
</style>

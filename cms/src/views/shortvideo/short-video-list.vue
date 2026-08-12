<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ShortVideoManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.shortVideoList.addShortVideo') }}</el-button>
          <span class="table-tip">{{ t('pages.shortVideoList.uploadTip') }}</span>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('common.title')">
            <el-input v-model="searchForm.title" clearable :placeholder="t('pages.shortVideoList.titlePlaceholder')"/>
          </el-form-item>
          <el-form-item :label="t('pages.shortVideoList.authorNickname')">
            <el-input v-model="searchForm.authorNickname" clearable :placeholder="t('pages.shortVideoList.authorNicknamePlaceholder')"/>
          </el-form-item>
          <el-form-item :label="t('common.status')">
            <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option :value="2" :label="t('common.onlyOnShelf')"/>
              <el-option :value="1" :label="t('common.onlyOffShelf')"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('common.sort')">
            <el-select
                v-model="searchForm.sortField"
                :placeholder="t('common.all')"
                style="width: 160px"
                @change="handleSearch"
            >
              <el-option value="" :label="t('pages.shortVideoList.sortDefault')"/>
              <el-option value="viewCount" :label="t('pages.shortVideoList.sortViewCount')"/>
              <el-option value="totalDiamondIncome" :label="t('pages.shortVideoList.sortDiamondIncome')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <div v-loading="storageStatLoading" class="list-summary">
          <span>{{ t('pages.shortVideoList.totalCount', {count: storageStat.totalCount}) }}</span>
          <span v-if="storageStat.imageDirPath">
            {{ t('pages.shortVideoList.dirUsage', {path: storageStat.imageDirPath, size: formatBytes(storageStat.imageDirUsedBytes)}) }}
          </span>
          <span
              v-if="storageStat.diskFreeRatio > 0"
              :class="{ 'disk-free-warning': isDiskFreeLow }"
          >
            {{ t('pages.shortVideoList.diskFree', {ratio: formatPercent(storageStat.diskFreeRatio)}) }}
          </span>
          <el-button link type="primary" @click="fetchStorageStat">{{ t('pages.shortVideoList.refreshStat') }}</el-button>
        </div>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('common.title')" prop="title" min-width="140"/>
          <el-table-column :label="t('pages.shortVideoList.cover')" width="100">
            <template #default="{ row }">
              <el-image
                  v-if="row.cover"
                  :preview-src-list="[row.cover]"
                  :src="row.cover"
                  fit="cover"
                  preview-teleported
                  style="width: 72px; height: 40px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.video')" min-width="200">
            <template #default="{ row }">
              <div v-if="row.video" class="table-video-cell">
                <video
                    :key="row.video"
                    :src="row.video"
                    class="table-video-preview"
                    controls
                    preload="metadata"
                />
                <el-button link type="primary" @click="openVideoPreview(row.video)">{{ t('pages.shortVideoList.enlargePreview') }}</el-button>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.authorNickname')" min-width="120">
            <template #default="{ row }">
              {{ row.authorNickname || '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.authorType')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.authorType === 1 ? 'warning' : 'success'">
                {{ row.authorType === 1 ? 'CMS' : 'App' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.authorId')" prop="authorId" width="100"/>
          <el-table-column :label="t('common.sort')" prop="sort" width="80"/>
          <el-table-column :label="t('pages.shortVideoList.isPaid')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.isPaid === 1 ? 'warning' : 'success'">
                {{ row.isPaid === 1 ? t('pages.shortVideoList.paid') : t('pages.shortVideoList.free') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.payDiamond')" width="110">
            <template #default="{ row }">
              {{ row.isPaid === 1 ? formatPrice(row.payDiamond) : '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.freeWatchSeconds')" width="110">
            <template #default="{ row }">
              {{ row.isPaid === 1 ? (row.freeWatchSeconds != null ? row.freeWatchSeconds : 15) : '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.videoCategory')" width="100">
            <template #default="{ row }">
              {{ categoryName(row.categoryId) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.source')" width="90">
            <template #default="{ row }">
              {{ sourceLabel(row.source) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.shortVideoList.likeCount')" prop="likeCount" width="90"/>
          <el-table-column :label="t('pages.shortVideoList.viewCount')" prop="viewCount" width="90"/>
          <el-table-column :label="t('pages.shortVideoList.watchCount')" prop="watchCount" width="90"/>
          <el-table-column :label="t('pages.shortVideoList.totalDiamondIncome')" width="120">
            <template #default="{ row }">
              {{ formatPrice(row.totalDiamondIncome) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.status')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
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

    <el-dialog v-model="dialogVisible" :close-on-click-modal="false" :title="dialogTitle" destroy-on-close width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('common.title')" prop="title">
          <el-input v-model="currentRow.title" :placeholder="t('pages.shortVideoList.titlePlaceholderForm')"/>
        </el-form-item>
        <el-form-item v-if="isCreateMode" :label="t('pages.shortVideoList.video')">
          <div class="video-upload-wrap">
            <el-upload
                :before-upload="beforeVideoSelect"
                :disabled="saving"
                :show-file-list="false"
                accept=".mp4,.webm,.mov,video/mp4,video/webm,video/quicktime"
                action="#"
                class="video-uploader"
            >
              <video
                  v-if="videoPreviewUrl"
                  :key="videoPreviewUrl"
                  :src="videoPreviewUrl"
                  class="dialog-video-preview"
                  controls
                  preload="metadata"
              />
              <div v-else class="video-uploader-placeholder">
                <el-button type="primary">{{ t('pages.shortVideoList.selectVideo') }}</el-button>
              </div>
            </el-upload>
            <span v-if="videoDuration > 0" class="form-tip">{{ t('pages.shortVideoList.videoDuration', {seconds: videoDuration}) }}</span>
            <el-button v-if="videoPreviewUrl" link type="danger" @click="clearCreateVideo">{{ t('pages.shortVideoList.clearVideo') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item v-else-if="videoPreviewUrl" :label="t('pages.shortVideoList.video')">
          <div class="preview-box">
            <video
                :key="videoPreviewUrl"
                :src="videoPreviewUrl"
                class="dialog-video-preview"
                controls
                preload="metadata"
            />
          </div>
          <div class="form-tip">{{ t('pages.shortVideoList.videoReadonlyTip') }}</div>
        </el-form-item>
        <el-form-item v-if="isCreateMode" :label="t('pages.shortVideoList.authorNickname')">
          <el-input v-model="currentRow.authorNickname" clearable maxlength="32" :placeholder="t('pages.shortVideoList.authorNicknamePlaceholderForm')"/>
          <div class="form-tip">{{ t('pages.shortVideoList.authorNicknameTip') }}</div>
        </el-form-item>
        <el-form-item v-else :label="t('pages.shortVideoList.authorNickname')">
          <span>{{ currentRow.authorNickname || '-' }}</span>
          <div v-if="currentRow.authorId" class="form-tip">{{ t('pages.shortVideoList.authorIdTip', {id: currentRow.authorId}) }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.shortVideoList.cover')" prop="cover">
          <div class="upload-wrap">
            <el-upload
                :before-upload="beforeCoverUpload"
                :disabled="coverUploading"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="cover-uploader"
            >
              <img v-if="coverPreviewUrl" :src="coverPreviewUrl" alt="cover" class="cover-preview"/>
              <div v-else class="cover-placeholder">
                <el-icon><Plus/></el-icon>
                <span>{{ t('pages.shortVideoList.uploadCover') }}</span>
              </div>
            </el-upload>
            <el-button v-if="coverPreviewUrl || currentRow.cover" link type="danger" @click="clearAsset">
              {{ t('pages.shortVideoList.removeCover') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('common.sort')" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
        </el-form-item>
        <el-form-item :label="t('pages.shortVideoList.isPaidLabel')" prop="isPaid">
          <el-radio-group v-model="currentRow.isPaid">
            <el-radio :label="0">{{ t('pages.shortVideoList.free') }}</el-radio>
            <el-radio :label="1">{{ t('pages.shortVideoList.paid') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="currentRow.isPaid === 1" :label="t('pages.shortVideoList.payDiamond')" prop="payDiamond">
          <el-input-number
              v-model="currentRow.payDiamond"
              :min="0.0001"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoList.payDiamondTip') }}</span>
        </el-form-item>
        <el-form-item v-if="currentRow.isPaid === 1" :label="t('pages.shortVideoList.freeWatchDuration')" prop="freeWatchSeconds">
          <el-input-number
              v-model="currentRow.freeWatchSeconds"
              :min="0"
              :step="1"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">{{ t('pages.shortVideoList.freeWatchDurationTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.shortVideoList.videoCategoryLabel')" prop="categoryId">
          <el-select v-model="currentRow.categoryId" clearable :placeholder="t('pages.shortVideoList.selectCategory')" style="width: 220px">
            <el-option
                v-for="item in categoryOptions"
                :key="item.id"
                :label="item.name"
                :value="Number(item.id)"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.shortVideoList.videoSource')" prop="source">
          <el-radio-group v-model="currentRow.source">
            <el-radio :label="1">{{ t('pages.shortVideoList.sourceOriginal') }}</el-radio>
            <el-radio :label="2">{{ t('pages.shortVideoList.sourceRepost') }}</el-radio>
            <el-radio :label="3">{{ t('pages.shortVideoList.sourceAi') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="saving" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="videoDialogVisible"
        destroy-on-close
        :title="t('pages.shortVideoList.videoPreview')"
        width="720px"
    >
      <video
          v-if="dialogVideoUrl"
          :key="dialogVideoUrl"
          :src="dialogVideoUrl"
          class="dialog-video-preview"
          controls
          autoplay
          preload="metadata"
      />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {shortVideoApi, uploadApi} from '@/api'
import type {ShortVideo, ShortVideoCategory} from '@/types/api.ts'
import {formatPrice, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  title: string
  authorNickname: string
  statusFilter: number
  sortField: '' | 'viewCount' | 'totalDiamondIncome'
}

interface ShortVideoForm {
  id: string
  title: string
  cover: string
  sort: number
  isPaid: number
  payDiamond: number
  freeWatchSeconds: number
  categoryId: number
  source: number
  authorNickname: string
  authorId: string
}

const {t} = useI18n()

const loading = ref(false)
const storageStatLoading = ref(false)
const saving = ref(false)
const tableData = ref<ShortVideo[]>([])
const categoryOptions = ref<ShortVideoCategory[]>([])
const total = ref(0)
const storageStat = reactive({
  totalCount: 0,
  imageDirPath: '',
  imageDirUsedBytes: 0,
  diskTotalBytes: 0,
  diskFreeBytes: 0,
  diskFreeRatio: 0,
})
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  title: '',
  authorNickname: '',
  statusFilter: 0,
  sortField: '',
})

const dialogVisible = ref(false)
const videoDialogVisible = ref(false)
const dialogVideoUrl = ref('')
const dialogTitle = ref('')
const defaultForm = (): ShortVideoForm => ({
  id: '',
  title: '',
  cover: '',
  sort: 0,
  isPaid: 0,
  payDiamond: 0,
  freeWatchSeconds: 15,
  categoryId: 0,
  source: 1,
  authorNickname: '',
  authorId: '',
})
const currentRow = ref<ShortVideoForm>(defaultForm())
const formRef = ref<FormInstance>()
const coverUploading = ref(false)
const videoPreviewUrl = ref('')
const coverPreviewUrl = ref('')
const videoFile = ref<File | null>(null)
const videoDuration = ref(0)
const maxVideoDuration = ref(60)
let videoObjectPreviewUrl = ''
const objectPreviewUrls = reactive<{ cover: string | null }>({
  cover: null
})

const isCreateMode = computed(() => !currentRow.value.id)

const isDiskFreeLow = computed(() => {
  const ratio = storageStat.diskFreeRatio
  return ratio > 0 && ratio < 30
})

const sourceLabel = (source: number) => {
  const map: Record<number, string> = {
    1: t('pages.shortVideoList.sourceOriginal'),
    2: t('pages.shortVideoList.sourceRepost'),
    3: t('pages.shortVideoList.sourceAi'),
  }
  return map[source] || '-'
}

const formatBytes = (bytes: number) => {
  const value = Number(bytes || 0)
  if (value <= 0) {
    return '0 B'
  }
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = value
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }
  return `${size.toFixed(unitIndex === 0 ? 0 : 2)} ${units[unitIndex]}`
}

const formatPercent = (ratio: number) => `${Number(ratio || 0).toFixed(1)}%`

const categoryName = (categoryId: number) => {
  if (!categoryId) {
    return '-'
  }
  const item = categoryOptions.value.find((c) => Number(c.id) === categoryId)
  return item?.name || String(categoryId)
}

const revokeObjectPreview = (field: 'cover') => {
  if (objectPreviewUrls[field]) {
    URL.revokeObjectURL(objectPreviewUrls[field]!)
    objectPreviewUrls[field] = null
  }
}

const resetAssetPreview = () => {
  revokeObjectPreview('cover')
  revokeVideoObjectPreview()
  videoPreviewUrl.value = ''
  coverPreviewUrl.value = ''
  videoFile.value = null
  videoDuration.value = 0
}

const revokeVideoObjectPreview = () => {
  if (videoObjectPreviewUrl) {
    URL.revokeObjectURL(videoObjectPreviewUrl)
    videoObjectPreviewUrl = ''
  }
}

const setCreateVideoPreview = (file: File) => {
  revokeVideoObjectPreview()
  const url = URL.createObjectURL(file)
  videoObjectPreviewUrl = url
  videoPreviewUrl.value = url
}

const detectVideoDuration = (file: File): Promise<number> => {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(file)
    const video = document.createElement('video')
    video.preload = 'metadata'
    video.onloadedmetadata = () => {
      URL.revokeObjectURL(url)
      resolve(Math.max(1, Math.ceil(video.duration || 0)))
    }
    video.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error('read video duration failed'))
    }
    video.src = url
  })
}

const allowedVideoExt = ['.mp4', '.webm', '.mov']

const getFileExt = (name: string) => {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.slice(idx).toLowerCase() : ''
}

const beforeVideoSelect = async (file: File): Promise<boolean> => {
  const ext = getFileExt(file.name)
  if (!allowedVideoExt.includes(ext)) {
    ElMessage.error(t('pages.shortVideoList.videoFormatError'))
    return false
  }
  try {
    const duration = await detectVideoDuration(file)
    if (maxVideoDuration.value > 0 && duration > maxVideoDuration.value) {
      ElMessage.error(t('pages.shortVideoList.videoDurationExceeded', {max: maxVideoDuration.value}))
      return false
    }
    videoFile.value = file
    videoDuration.value = duration
    setCreateVideoPreview(file)
  } catch (error) {
    console.error('read video failed:', error)
    ElMessage.error(t('pages.shortVideoList.readVideoFailed'))
    return false
  }
  return false
}

const clearCreateVideo = () => {
  videoFile.value = null
  videoDuration.value = 0
  revokeVideoObjectPreview()
  videoPreviewUrl.value = ''
}

const openVideoPreview = (url: string) => {
  dialogVideoUrl.value = url
  videoDialogVisible.value = true
}

const beforeCoverUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.shortVideoList.coverImageOnly'))
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  coverUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.cover = res.fileName
    const previewUrl = res.fileUrl || URL.createObjectURL(file)
    revokeObjectPreview('cover')
    coverPreviewUrl.value = previewUrl
    if (!res.fileUrl) {
      objectPreviewUrls.cover = previewUrl
    }
    formRef.value?.validateField('cover').catch(() => undefined)
    ElMessage.success(t('pages.shortVideoList.uploadSuccess'))
  } catch (error) {
    console.error('upload failed:', error)
    ElMessage.error(t('pages.shortVideoList.uploadFailed'))
  } finally {
    coverUploading.value = false
  }
}

const clearAsset = () => {
  currentRow.value.cover = ''
  revokeObjectPreview('cover')
  coverPreviewUrl.value = ''
  formRef.value?.validateField('cover').catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAssetPreview()
  }
})

watch(() => currentRow.value.isPaid, (paid) => {
  if (paid === 1) {
    if (!currentRow.value.payDiamond || currentRow.value.payDiamond <= 0) {
      currentRow.value.payDiamond = 1
    }
    if (currentRow.value.freeWatchSeconds == null || currentRow.value.freeWatchSeconds < 0) {
      currentRow.value.freeWatchSeconds = 15
    }
    return
  }
  currentRow.value.payDiamond = 0
  currentRow.value.freeWatchSeconds = 0
})

const formRules = computed<FormRules>(() => ({
  title: [
    {required: true, message: t('pages.shortVideoList.titleRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.shortVideoList.titleLength'), trigger: 'blur'}
  ],
  payDiamond: [
    {
      validator: (_rule, value, callback) => {
        if (currentRow.value.isPaid === 1 && (!value || value <= 0)) {
          callback(new Error(t('pages.shortVideoList.payDiamondRequired')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
  freeWatchSeconds: [
    {
      validator: (_rule, value, callback) => {
        if (currentRow.value.isPaid === 1 && (value == null || value < 0)) {
          callback(new Error(t('pages.shortVideoList.freeWatchSecondsMin')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
  source: [{required: true, message: t('pages.shortVideoList.sourceRequired'), trigger: 'change'}],
}))

const fetchStorageStat = async () => {
  storageStatLoading.value = true
  try {
    const response = await shortVideoApi.getShortVideoStorageStat()
    storageStat.totalCount = response.totalCount || 0
    storageStat.imageDirPath = response.imageDirPath || ''
    storageStat.imageDirUsedBytes = response.imageDirUsedBytes || 0
    storageStat.diskTotalBytes = response.diskTotalBytes || 0
    storageStat.diskFreeBytes = response.diskFreeBytes || 0
    storageStat.diskFreeRatio = response.diskFreeRatio || 0
  } catch (error) {
    console.error('fetch storage stat failed:', error)
  } finally {
    storageStatLoading.value = false
  }
}

const fetchShortVideoCfg = async () => {
  try {
    const response = await shortVideoApi.getShortVideoCfg()
    if (response.cfg?.maxDuration) {
      maxVideoDuration.value = Math.max(1, response.cfg.maxDuration)
    }
  } catch (error) {
    console.error('fetch short video cfg failed:', error)
  }
}

const fetchCategoryOptions = async () => {
  try {
    const response = await shortVideoApi.getShortVideoCategoryList({
      pageIndex: 1,
      pageSize: 100,
    })
    categoryOptions.value = response.data || []
  } catch (error) {
    console.error('fetch categories failed:', error)
  }
}

const fetchShortVideoList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoList({
      title: searchForm.title,
      authorNickname: searchForm.authorNickname,
      statusFilter: searchForm.statusFilter,
      sortField: searchForm.sortField,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('fetch short video list failed:', error)
    ElMessage.error(t('pages.shortVideoList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchShortVideoList()
}

const resetSearch = () => {
  searchForm.title = ''
  searchForm.authorNickname = ''
  searchForm.statusFilter = 0
  searchForm.sortField = ''
  currentPage.value = 1
  fetchShortVideoList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchShortVideoList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchShortVideoList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.shortVideoList.addShortVideo')
  currentRow.value = defaultForm()
  resetAssetPreview()
  dialogVisible.value = true
}

const handleEdit = (row: ShortVideo) => {
  dialogTitle.value = t('pages.shortVideoList.editShortVideo')
  currentRow.value = {
    id: row.id,
    title: row.title,
    cover: row.coverName || '',
    sort: Number(row.sort) || 0,
    isPaid: row.isPaid ?? 0,
    payDiamond: row.isPaid === 1 ? (truncateNumber(row.payDiamond) || 1) : 0,
    freeWatchSeconds: row.isPaid === 1 ? (row.freeWatchSeconds != null ? Number(row.freeWatchSeconds) : 15) : 0,
    categoryId: Number(row.categoryId) || 0,
    source: row.source || 1,
    authorId: row.authorId || '',
    authorNickname: row.authorNickname || '',
  }
  videoFile.value = null
  videoDuration.value = Number(row.duration) || 0
  videoPreviewUrl.value = row.video || ''
  coverPreviewUrl.value = row.cover || ''
  dialogVisible.value = true
}

const handleDelete = async (row: ShortVideo) => {
  try {
    await ElMessageBox.confirm(t('pages.shortVideoList.deleteConfirm', {title: row.title}), t('common.confirmDelete'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await shortVideoApi.deleteShortVideo(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchShortVideoList()
  } catch (error) {
    console.error('delete failed:', error)
  }
}

const handleOnShelf = async (row: ShortVideo) => {
  try {
    await shortVideoApi.onShelfShortVideo(row.id)
    ElMessage.success(t('pages.shortVideoList.onShelfSuccess'))
    fetchShortVideoList()
  } catch (error) {
    console.error('on shelf failed:', error)
    ElMessage.error(t('pages.shortVideoList.onShelfFailed'))
  }
}

const handleOffShelf = async (row: ShortVideo) => {
  try {
    await ElMessageBox.confirm(t('pages.shortVideoList.offShelfConfirm', {title: row.title}), t('common.confirmOffShelf'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await shortVideoApi.offShelfShortVideo(row.id)
    ElMessage.success(t('pages.shortVideoList.offShelfSuccess'))
    fetchShortVideoList()
  } catch (error) {
    console.error('off shelf failed:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      const payload = {
        title: currentRow.value.title,
        cover: currentRow.value.cover,
        sort: currentRow.value.sort,
        isPaid: currentRow.value.isPaid,
        payDiamond: currentRow.value.isPaid === 1 ? currentRow.value.payDiamond : 0,
        freeWatchSeconds: currentRow.value.isPaid === 1 ? currentRow.value.freeWatchSeconds : 0,
        categoryId: currentRow.value.categoryId || 0,
        source: currentRow.value.source,
      }
      if (isCreateMode.value) {
        if (!videoFile.value) {
          ElMessage.error(t('pages.shortVideoList.selectVideoRequired'))
          return
        }
        ElMessage.info(t('pages.shortVideoList.uploadingVideo'))
        const videoRes = await uploadApi.uploadFile(videoFile.value)
        await shortVideoApi.createShortVideo({
          video: videoRes.fileName,
          coverName: currentRow.value.cover || undefined,
          title: payload.title,
          sort: payload.sort,
          isPaid: payload.isPaid,
          payDiamond: payload.payDiamond,
          freeWatchSeconds: payload.freeWatchSeconds,
          categoryId: payload.categoryId,
          source: payload.source,
          duration: videoDuration.value,
          authorNickname: currentRow.value.authorNickname || undefined,
        })
        ElMessage.success(t('common.createSuccess'))
      } else {
        await shortVideoApi.updateShortVideo({id: currentRow.value.id, ...payload})
        ElMessage.success(t('common.updateSuccess'))
      }
      dialogVisible.value = false
      fetchShortVideoList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.shortVideoList.saveFailed'))
    } finally {
      saving.value = false
    }
  })
}

onMounted(() => {
  fetchShortVideoCfg()
  fetchCategoryOptions()
  fetchStorageStat()
  fetchShortVideoList()
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

.content {
  min-height: 400px;
}

.table-header {
  margin-bottom: 15px;
}

.table-tip {
  color: #909399;
  font-size: 13px;
}

.search-form {
  margin-bottom: 12px;
}

.list-summary {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  min-height: 32px;
}

.disk-free-warning {
  color: var(--el-color-danger);
  font-weight: 600;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.upload-wrap {
  width: 100%;
}

.preview-box {
  margin-top: 12px;
  width: 100%;
}

.video-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  width: 100%;
}

.video-uploader :deep(.el-upload) {
  display: block;
  cursor: pointer;
}

.video-uploader-placeholder {
  min-width: 120px;
}

.table-video-cell {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.table-video-preview {
  display: block;
  width: 160px;
  min-height: 90px;
  max-height: 90px;
  border-radius: 4px;
  background: #000;
}

.dialog-video-preview {
  display: block;
  width: 100%;
  min-height: 220px;
  max-height: 480px;
  border-radius: 6px;
  background: #000;
}

.file-name {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.form-tip {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.cover-uploader {
  display: inline-block;
}

.cover-preview {
  width: 160px;
  height: 90px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
}

.cover-placeholder {
  width: 160px;
  height: 90px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  cursor: pointer;
}
</style>

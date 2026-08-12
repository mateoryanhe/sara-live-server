<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.BotAnchorManagement') }}</span>
        </div>
      </template>

      <div class="table-header">
        <el-button type="primary" @click="handleAdd">{{ t('pages.botAnchorList.addBotAnchor') }}</el-button>
        <el-button
            :disabled="!canBatchStartLive"
            :loading="batchOperating"
            type="success"
            @click="handleBatchStartLive"
        >
          {{ t('pages.botAnchorList.batchStartLive') }}
        </el-button>
        <el-button
            :disabled="!canBatchStopLive"
            :loading="batchOperating"
            type="danger"
            @click="handleBatchStopLive"
        >
          {{ t('pages.botAnchorList.batchStopLive') }}
        </el-button>
        <span v-if="selectedRows.length" class="selection-tip">
          {{ t('common.selectedCount', {count: selectedRows.length}) }}
        </span>
      </div>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.key"
              clearable
              :placeholder="t('pages.botAnchorList.keywordPlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          ref="tableRef"
          v-loading="loading"
          :data="tableData"
          row-key="id"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column :selectable="isRowSelectable" type="selection" width="48"/>
        <el-table-column :label="t('common.userId')" prop="id" width="180"/>
        <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.avatar"
                :preview-src-list="[row.avatar]"
                :src="row.avatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.guildId')" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.liveRoomId')" prop="roomId" width="180">
          <template #default="{ row }">{{ row.roomId || row.id || '-' }}</template>
        </el-table-column>
        <el-table-column
            :label="t('pages.botAnchorList.roomTitle')"
            min-width="140"
            prop="roomTitle"
            show-overflow-tooltip
        >
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.roomType')" width="100">
          <template #default="{ row }">
            <el-tag :type="categoryTagType(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.liveRoomTag')" min-width="120">
          <template #default="{ row }">{{ row.tagName || (row.tagId ? row.tagId : '-') }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.cloudPlayerVideo')" min-width="200">
          <template #default="{ row }">
            <div v-if="row.cloudPlayerVideo" class="table-video-cell">
              <video
                  :src="row.cloudPlayerVideo"
                  class="table-video-preview"
                  controls
                  preload="metadata"
              />
              <el-button link type="primary" @click="openVideoPreview(row.cloudPlayerVideo)">
                {{ t('pages.botAnchorList.enlargePreview') }}
              </el-button>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.pushStream')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.pushStream ? 'success' : 'info'">
              {{ row.pushStream ? t('common.yes') : t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.testAccount')" width="90">
          <template #default="{ row }">
            <el-tag :type="row.isTest ? 'warning' : 'info'">
              {{ row.isTest ? t('common.yes') : t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.botAnchorStatus === 1 ? 'success' : 'info'">
              {{ row.botAnchorStatus === 1 ? t('common.enabled') : t('common.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.botAnchorList.liveStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? t('common.live') : t('common.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" prop="createdAt" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170">
          <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="280">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button
                v-if="row.botAnchorStatus === 1 && row.liveStatus !== 1"
                link
                type="success"
                @click="handleStartLive(row)"
            >
              {{ t('pages.botAnchorList.startLive') }}
            </el-button>
            <el-button
                v-if="row.botAnchorStatus === 1 && row.liveStatus === 1"
                link
                type="danger"
                @click="handleStopLive(row)"
            >
              {{ t('pages.botAnchorList.stopLive') }}
            </el-button>
            <el-button
                :type="row.botAnchorStatus === 1 ? 'warning' : 'success'"
                link
                @click="toggleStatus(row)"
            >
              {{
                row.botAnchorStatus === 1
                    ? t('pages.botAnchorList.disable')
                    : t('pages.botAnchorList.enable')
              }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
            v-model:current-page="pagination.pageIndex"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <el-dialog
        v-model="dialogVisible"
        :close-on-click-modal="false"
        :title="dialogTitle"
        destroy-on-close
        width="640px"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item v-if="formData.id" :label="t('common.userId')">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')" prop="nickname">
          <el-input
              v-model="formData.nickname"
              maxlength="32"
              :placeholder="t('pages.botAnchorList.enterNickname')"
              show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.roomTitle')" prop="roomTitle">
          <el-input
              v-model="formData.roomTitle"
              maxlength="128"
              :placeholder="t('pages.botAnchorList.enterRoomTitle')"
              show-word-limit
          />
        </el-form-item>
        <el-form-item v-if="!formData.id" :label="t('pages.botAnchorList.guildId')">
          <el-input
              v-model="formData.guildId"
              :placeholder="t('pages.botAnchorList.guildIdOptional')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.roomType')" prop="category">
          <el-select
              v-model="formData.category"
              :placeholder="t('pages.botAnchorList.selectRoomType')"
              style="width: 220px"
          >
            <el-option :value="LIVE_ROOM_CATEGORY_HOT" :label="t('pages.botAnchorList.categoryShow')"/>
            <el-option :value="LIVE_ROOM_CATEGORY_GAME" :label="t('pages.botAnchorList.categoryGame')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.liveRoomTag')" prop="tagId">
          <el-select
              v-model="formData.tagId"
              clearable
              :placeholder="t('pages.botAnchorList.selectTag')"
              style="width: 220px"
          >
            <el-option :value="0" :label="t('pages.botAnchorList.none')"/>
            <el-option
                v-for="item in tagOptions"
                :key="item.id"
                :label="item.name"
                :value="Number(item.id)"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.cloudMp4')" prop="cloudPlayerVideo">
          <div class="video-upload-wrap">
            <el-upload
                action="#"
                :before-upload="beforeVideoUpload"
                :disabled="videoUploading"
                :http-request="doVideoUpload"
                :show-file-list="false"
                accept=".mp4,video/mp4"
                class="video-uploader"
            >
              <video
                  v-if="videoPreviewUrl"
                  :key="videoPreviewUrl"
                  :src="videoPreviewUrl"
                  class="video-preview"
                  controls
                  preload="metadata"
              />
              <div v-else class="video-uploader-placeholder">
                <el-button :loading="videoUploading" type="primary">
                  {{ t('pages.botAnchorList.uploadMp4') }}
                </el-button>
              </div>
            </el-upload>
            <span v-if="formData.cloudPlayerVideo" class="video-file-name">{{ formData.cloudPlayerVideo }}</span>
            <el-progress
                v-if="videoUploading && videoUploadProgress > 0"
                :percentage="videoUploadProgress"
                :stroke-width="8"
                style="width: 100%; max-width: 480px"
            />
            <el-button
                v-if="formData.cloudPlayerVideo || videoPreviewUrl"
                link
                type="danger"
                @click="clearVideo"
            >
              {{ t('pages.botAnchorList.clearVideo') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.pushStream')">
          <el-switch v-model="formData.pushStream"/>
        </el-form-item>
        <el-form-item :label="t('pages.botAnchorList.testAccount')">
          <el-switch v-model="formData.isTest"/>
          <span class="form-tip">{{ t('pages.botAnchorList.isTestTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('common.avatar')" prop="avatar">
          <div class="avatar-upload-wrap">
            <el-upload
                :before-upload="beforeAvatarUpload"
                :disabled="avatarUploading"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                class="avatar-uploader"
            >
              <el-image
                  v-if="avatarPreviewUrl"
                  :src="avatarPreviewUrl"
                  fit="cover"
                  style="width:80px;height:80px;border-radius:50%"
              />
              <div v-else class="avatar-uploader-placeholder">
                <el-icon class="avatar-uploader-icon">
                  <Plus/>
                </el-icon>
              </div>
            </el-upload>
            <el-button v-if="formData.avatar" link type="danger" @click="clearAvatar">
              {{ t('pages.botAnchorList.clearAvatar') }}
            </el-button>
          </div>
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
        :title="t('pages.botAnchorList.cloudVideoPreview')"
        width="720px"
    >
      <video
          v-if="dialogVideoUrl"
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
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type TableInstance, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {botAnchorApi, liveRoomTagApi, uploadApi} from '@/api'
import type {BotAnchorListItem, LiveRoomTag} from '@/types/api'

const {t} = useI18n()
const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_PRIVATE = 3

interface SearchForm {
  key: string
}

interface BotAnchorForm {
  id: string
  nickname: string
  roomTitle: string
  avatar: string
  guildId: string
  category: number
  tagId: number
  cloudPlayerVideo: string
  pushStream: boolean
  isTest: boolean
}

const loading = ref(false)
const saving = ref(false)
const batchOperating = ref(false)
const avatarUploading = ref(false)
const videoUploading = ref(false)
const videoUploadProgress = ref(0)
const tableData = ref<BotAnchorListItem[]>([])
const selectedRows = ref<BotAnchorListItem[]>([])
const tableRef = ref<TableInstance>()
const tagOptions = ref<LiveRoomTag[]>([])
const searchForm = reactive<SearchForm>({key: ''})
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0
})

const dialogVisible = ref(false)
const videoDialogVisible = ref(false)
const dialogVideoUrl = ref('')
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<BotAnchorForm>({
  id: '',
  nickname: '',
  roomTitle: '',
  avatar: '',
  guildId: '',
  category: LIVE_ROOM_CATEGORY_HOT,
  tagId: 0,
  cloudPlayerVideo: '',
  pushStream: false,
  isTest: false
})
const avatarChanged = ref(false)
const videoChanged = ref(false)
let objectPreviewUrl = ''
let objectVideoPreviewUrl = ''
const avatarPreviewUrl = ref('')
const videoPreviewUrl = ref('')

const formRules = computed<FormRules>(() => ({
  nickname: [
    {required: true, message: t('pages.botAnchorList.nicknameRequired'), trigger: 'blur'},
    {min: 1, max: 32, message: t('pages.botAnchorList.nicknameLength'), trigger: 'blur'}
  ],
  roomTitle: [{max: 128, message: t('pages.botAnchorList.roomTitleMaxLength'), trigger: 'blur'}],
  category: [{required: true, message: t('pages.botAnchorList.categoryRequired'), trigger: 'change'}]
}))

const isStartableRow = (row: BotAnchorListItem) => row.botAnchorStatus === 1 && row.liveStatus !== 1
const isStoppableRow = (row: BotAnchorListItem) => row.botAnchorStatus === 1 && row.liveStatus === 1
const isRowSelectable = (row: BotAnchorListItem) => isStartableRow(row) || isStoppableRow(row)

const startableSelectedRows = computed(() => selectedRows.value.filter(isStartableRow))
const stoppableSelectedRows = computed(() => selectedRows.value.filter(isStoppableRow))
const canBatchStartLive = computed(() => startableSelectedRows.value.length > 0 && !batchOperating.value)
const canBatchStopLive = computed(() => stoppableSelectedRows.value.length > 0 && !batchOperating.value)

const handleSelectionChange = (rows: BotAnchorListItem[]) => {
  selectedRows.value = rows
}

const clearSelection = () => {
  tableRef.value?.clearSelection()
  selectedRows.value = []
}

const categoryLabel = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_HOT) return t('pages.botAnchorList.categoryHot')
  if (category === LIVE_ROOM_CATEGORY_GAME) return t('pages.botAnchorList.categoryGame')
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return t('pages.botAnchorList.categoryPrivate')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return 'danger'
  if (category === LIVE_ROOM_CATEGORY_GAME) return 'warning'
  if (category === LIVE_ROOM_CATEGORY_HOT) return 'success'
  return 'info'
}

const fetchTagOptions = async () => {
  try {
    const response = await liveRoomTagApi.getLiveRoomTagList({
      pageIndex: 1,
      pageSize: 200
    })
    tagOptions.value = response.data || []
  } catch (error) {
    console.error('fetchTagOptions failed:', error)
    ElMessage.error(t('pages.botAnchorList.fetchTagsFailed'))
  }
}

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = ''
  }
}

const setAvatarPreview = (url: string, isObjectUrl = false) => {
  revokeObjectPreview()
  if (isObjectUrl) {
    objectPreviewUrl = url
  }
  avatarPreviewUrl.value = url
}

const revokeVideoObjectPreview = () => {
  if (objectVideoPreviewUrl) {
    URL.revokeObjectURL(objectVideoPreviewUrl)
    objectVideoPreviewUrl = ''
  }
}

const setVideoPreview = (url: string, isObjectUrl = false) => {
  revokeVideoObjectPreview()
  if (isObjectUrl) {
    objectVideoPreviewUrl = url
  }
  videoPreviewUrl.value = url
}

const openVideoPreview = (url: string) => {
  dialogVideoUrl.value = url
  videoDialogVisible.value = true
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  if (Number.isNaN(date.getTime())) return dateString
  return date.toLocaleString(undefined, {hour12: false})
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await botAnchorApi.getBotAnchorList({
      key: searchForm.key,
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize
    })
    tableData.value = response.data
    pagination.total = response.total
  } catch (error) {
    console.error('fetchList failed:', error)
    ElMessage.error(t('pages.botAnchorList.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  pagination.pageIndex = 1
  fetchList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  clearSelection()
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  clearSelection()
  fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.botAnchorList.addDialogTitle')
  formData.value = {
    id: '',
    nickname: '',
    roomTitle: '',
    avatar: '',
    guildId: '',
    category: LIVE_ROOM_CATEGORY_HOT,
    tagId: 0,
    cloudPlayerVideo: '',
    pushStream: false,
    isTest: false
  }
  avatarChanged.value = false
  videoChanged.value = false
  setAvatarPreview('')
  setVideoPreview('')
  dialogVisible.value = true
}

const handleEdit = (row: BotAnchorListItem) => {
  dialogTitle.value = t('pages.botAnchorList.editDialogTitle')
  formData.value = {
    id: row.id,
    nickname: row.nickname || '',
    roomTitle: row.roomTitle || '',
    avatar: '',
    guildId: '',
    category: row.category || LIVE_ROOM_CATEGORY_HOT,
    tagId: Number(row.tagId) || 0,
    cloudPlayerVideo: row.cloudPlayerVideoFile || '',
    pushStream: !!row.pushStream,
    isTest: !!row.isTest
  }
  avatarChanged.value = false
  videoChanged.value = false
  setAvatarPreview(row.avatar || '')
  setVideoPreview(row.cloudPlayerVideo || '')
  dialogVisible.value = true
}

const beforeAvatarUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.botAnchorList.imageOnly'))
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  avatarUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    formData.value.avatar = res.fileName
    avatarChanged.value = true
    setAvatarPreview(URL.createObjectURL(file), true)
    ElMessage.success(t('pages.botAnchorList.uploadSuccess'))
  } catch (error) {
    console.error('upload failed:', error)
    ElMessage.error(t('pages.botAnchorList.uploadFailed'))
  } finally {
    avatarUploading.value = false
  }
}

const clearAvatar = () => {
  formData.value.avatar = ''
  avatarChanged.value = true
  setAvatarPreview('')
}

const beforeVideoUpload = (file: File): boolean => {
  const isMp4 = file.type === 'video/mp4' || file.name.toLowerCase().endsWith('.mp4')
  if (!isMp4) {
    ElMessage.error(t('pages.botAnchorList.mp4Only'))
    return false
  }

  return true
}

const doVideoUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  videoUploading.value = true
  videoUploadProgress.value = 0
  try {
    const res = await uploadApi.uploadFile(file, (percent) => {
      videoUploadProgress.value = percent
    })
    formData.value.cloudPlayerVideo = res.fileName
    videoChanged.value = true
    const previewUrl = res.fileUrl || URL.createObjectURL(file)
    setVideoPreview(previewUrl, !res.fileUrl)
    ElMessage.success(t('pages.botAnchorList.videoUploadSuccess'))
  } catch (error) {
    console.error('video upload failed:', error)
    ElMessage.error(t('pages.botAnchorList.videoUploadFailed'))
  } finally {
    videoUploading.value = false
    videoUploadProgress.value = 0
  }
}

const clearVideo = () => {
  formData.value.cloudPlayerVideo = ''
  videoChanged.value = true
  setVideoPreview('')
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      if (formData.value.id) {
        const payload: {
          id: string
          nickname: string
          roomTitle: string
          category: number
          tagId: number
          pushStream: boolean
          isTest: boolean
          avatar?: string
          cloudPlayerVideo?: string
        } = {
          id: formData.value.id,
          nickname: formData.value.nickname,
          roomTitle: formData.value.roomTitle,
          category: formData.value.category,
          tagId: formData.value.tagId || 0,
          pushStream: formData.value.pushStream,
          isTest: formData.value.isTest
        }
        if (avatarChanged.value) {
          payload.avatar = formData.value.avatar
        }
        if (videoChanged.value) {
          payload.cloudPlayerVideo = formData.value.cloudPlayerVideo
        }
        await botAnchorApi.updateBotAnchor(payload)
      } else {
        await botAnchorApi.createBotAnchor({
          nickname: formData.value.nickname,
          roomTitle: formData.value.roomTitle,
          avatar: formData.value.avatar,
          guildId: formData.value.guildId || undefined,
          category: formData.value.category,
          tagId: formData.value.tagId || 0,
          cloudPlayerVideo: formData.value.cloudPlayerVideo || undefined,
          pushStream: formData.value.pushStream,
          isTest: formData.value.isTest
        })
      }
      ElMessage.success(formData.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.botAnchorList.saveFailed'))
    } finally {
      saving.value = false
    }
  })
}

const handleStopLive = async (row: BotAnchorListItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.botAnchorList.stopLiveConfirm', {name: row.nickname || row.id}),
        t('pages.botAnchorList.stopLiveTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await botAnchorApi.stopBotAnchorLive({id: row.id})
    ElMessage.success(t('pages.botAnchorList.stopLiveSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('stopLive failed:', error)
      ElMessage.error(t('pages.botAnchorList.stopLiveFailed'))
    }
  }
}

const handleBatchStartLive = async () => {
  const rows = startableSelectedRows.value
  if (!rows.length) {
    ElMessage.warning(t('pages.botAnchorList.selectStartable'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.botAnchorList.batchStartConfirm', {count: rows.length}),
        t('pages.botAnchorList.batchStartTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    batchOperating.value = true
    const response = await botAnchorApi.batchStartBotAnchorLive({ids: rows.map((row) => row.id)})
    if (response.failCount > 0) {
      ElMessage.warning(
          t('pages.botAnchorList.batchStartPartial', {
            success: response.successCount,
            fail: response.failCount
          })
      )
    } else {
      ElMessage.success(t('pages.botAnchorList.batchStartSuccess', {count: response.successCount}))
    }
    clearSelection()
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('batchStartLive failed:', error)
      ElMessage.error(t('pages.botAnchorList.batchStartFailed'))
    }
  } finally {
    batchOperating.value = false
  }
}

const handleBatchStopLive = async () => {
  const rows = stoppableSelectedRows.value
  if (!rows.length) {
    ElMessage.warning(t('pages.botAnchorList.selectStoppable'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.botAnchorList.batchStopConfirm', {count: rows.length}),
        t('pages.botAnchorList.batchStopTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    batchOperating.value = true
    const response = await botAnchorApi.batchStopBotAnchorLive({ids: rows.map((row) => row.id)})
    if (response.failCount > 0) {
      ElMessage.warning(
          t('pages.botAnchorList.batchStopPartial', {
            success: response.successCount,
            fail: response.failCount
          })
      )
    } else {
      ElMessage.success(t('pages.botAnchorList.batchStopSuccess', {count: response.successCount}))
    }
    clearSelection()
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('batchStopLive failed:', error)
      ElMessage.error(t('pages.botAnchorList.batchStopFailed'))
    }
  } finally {
    batchOperating.value = false
  }
}

const handleStartLive = async (row: BotAnchorListItem) => {
  try {
    await ElMessageBox.confirm(
        t('pages.botAnchorList.startLiveConfirm', {name: row.nickname || row.id}),
        t('pages.botAnchorList.startLiveTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await botAnchorApi.startBotAnchorLive({id: row.id})
    ElMessage.success(t('pages.botAnchorList.startLiveSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('startLive failed:', error)
      ElMessage.error(t('pages.botAnchorList.startLiveFailed'))
    }
  }
}

const toggleStatus = async (row: BotAnchorListItem) => {
  const nextStatus = row.botAnchorStatus === 1 ? 0 : 1
  const actionText = nextStatus === 1
      ? t('pages.botAnchorList.toggleEnable')
      : t('pages.botAnchorList.toggleDisable')
  try {
    await ElMessageBox.confirm(
        t('pages.botAnchorList.toggleConfirm', {action: actionText, name: row.nickname || row.id}),
        t('pages.botAnchorList.toggleTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await botAnchorApi.setBotAnchorStatus({id: row.id, status: nextStatus})
    ElMessage.success(t('pages.botAnchorList.toggleSuccess', {action: actionText}))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('toggleStatus failed:', error)
      ElMessage.error(t('pages.botAnchorList.operationFailed'))
    }
  }
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    revokeVideoObjectPreview()
    avatarPreviewUrl.value = ''
    videoPreviewUrl.value = ''
  }
})

onMounted(() => {
  fetchTagOptions()
  fetchList()
})
</script>

<style scoped>
.page-container {
  padding: 16px;
}

.card-header {
  font-weight: 600;
}

.table-header {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.selection-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.avatar-upload-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 50%;
  cursor: pointer;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.avatar-uploader-placeholder {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-uploader-icon {
  font-size: 24px;
  color: #8c939d;
}

.video-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
}

.video-uploader :deep(.el-upload) {
  display: block;
  cursor: pointer;
}

.video-uploader-placeholder {
  min-width: 120px;
}

.video-preview {
  display: block;
  width: 100%;
  max-width: 480px;
  min-height: 135px;
  max-height: 270px;
  border-radius: 6px;
  background: #000;
}

.table-video-cell {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.table-video-preview {
  width: 160px;
  max-height: 90px;
  border-radius: 4px;
  background: #000;
}

.dialog-video-preview {
  width: 100%;
  max-height: 480px;
  border-radius: 6px;
  background: #000;
}

.video-file-name {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style>

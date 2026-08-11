<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>机器人主播</span>
        </div>
      </template>

      <div class="table-header">
        <el-button type="primary" @click="handleAdd">新增机器人主播</el-button>
        <el-button
            :disabled="!canBatchStartLive"
            :loading="batchOperating"
            type="success"
            @click="handleBatchStartLive"
        >
          批量开播
        </el-button>
        <el-button
            :disabled="!canBatchStopLive"
            :loading="batchOperating"
            type="danger"
            @click="handleBatchStopLive"
        >
          批量下播
        </el-button>
        <span v-if="selectedRows.length" class="selection-tip">已选 {{ selectedRows.length }} 项</span>
      </div>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item label="关键字">
          <el-input v-model="searchForm.key" clearable placeholder="用户ID/昵称"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
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
        <el-table-column label="用户ID" prop="id" width="180"/>
        <el-table-column label="昵称" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="头像" width="80">
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
        <el-table-column label="工会ID" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播间ID" prop="roomId" width="180">
          <template #default="{ row }">{{ row.roomId || row.id || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播间标题" min-width="140" prop="roomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column label="房间类型" width="100">
          <template #default="{ row }">
            <el-tag :type="categoryTagType(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="直播间标签" min-width="120">
          <template #default="{ row }">{{ row.tagName || (row.tagId ? row.tagId : '-') }}</template>
        </el-table-column>
        <el-table-column label="云播视频" min-width="200">
          <template #default="{ row }">
            <div v-if="row.cloudPlayerVideo" class="table-video-cell">
              <video
                  :src="row.cloudPlayerVideo"
                  class="table-video-preview"
                  controls
                  preload="metadata"
              />
              <el-button link type="primary" @click="openVideoPreview(row.cloudPlayerVideo)">放大预览</el-button>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="是否推流" width="100">
          <template #default="{ row }">
            <el-tag :type="row.pushStream ? 'success' : 'info'">
              {{ row.pushStream ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="测试号" width="90">
          <template #default="{ row }">
            <el-tag :type="row.isTest ? 'warning' : 'info'">
              {{ row.isTest ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.botAnchorStatus === 1 ? 'success' : 'info'">
              {{ row.botAnchorStatus === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="直播状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? '直播中' : '未开播' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createdAt" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="更新时间" prop="updatedAt" width="170">
          <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="280">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button
                v-if="row.botAnchorStatus === 1 && row.liveStatus !== 1"
                link
                type="success"
                @click="handleStartLive(row)"
            >
              开播
            </el-button>
            <el-button
                v-if="row.botAnchorStatus === 1 && row.liveStatus === 1"
                link
                type="danger"
                @click="handleStopLive(row)"
            >
              下播
            </el-button>
            <el-button
                :type="row.botAnchorStatus === 1 ? 'warning' : 'success'"
                link
                @click="toggleStatus(row)"
            >
              {{ row.botAnchorStatus === 1 ? '停用' : '启用' }}
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
        <el-form-item v-if="formData.id" label="用户ID">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="formData.nickname" maxlength="32" placeholder="请输入昵称" show-word-limit/>
        </el-form-item>
        <el-form-item label="直播间标题" prop="roomTitle">
          <el-input v-model="formData.roomTitle" maxlength="128" placeholder="请输入直播间标题" show-word-limit/>
        </el-form-item>
        <el-form-item v-if="!formData.id" label="工会ID">
          <el-input v-model="formData.guildId" placeholder="可选,留空表示不绑定工会"/>
        </el-form-item>
        <el-form-item label="房间类型" prop="category">
          <el-select v-model="formData.category" placeholder="请选择房间类型" style="width: 220px">
            <el-option :value="LIVE_ROOM_CATEGORY_HOT" label="秀场"/>
            <el-option :value="LIVE_ROOM_CATEGORY_GAME" label="游戏"/>
          </el-select>
        </el-form-item>
        <el-form-item label="直播间标签" prop="tagId">
          <el-select v-model="formData.tagId" clearable placeholder="请选择标签" style="width: 220px">
            <el-option :value="0" label="无"/>
            <el-option
                v-for="item in tagOptions"
                :key="item.id"
                :label="item.name"
                :value="Number(item.id)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="云播MP4" prop="cloudPlayerVideo">
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
                <el-button :loading="videoUploading" type="primary">上传MP4</el-button>
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
              清除视频
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="是否推流">
          <el-switch v-model="formData.pushStream"/>
        </el-form-item>
        <el-form-item label="测试号">
          <el-switch v-model="formData.isTest"/>
          <span class="form-tip">开启后 App 端直播间信息 isTest 为 true，不参与服务端逻辑</span>
        </el-form-item>
        <el-form-item label="头像" prop="avatar">
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
            <el-button v-if="formData.avatar" link type="danger" @click="clearAvatar">清除头像</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button :loading="saving" type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="videoDialogVisible"
        destroy-on-close
        title="云播视频预览"
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type TableInstance, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {botAnchorApi, liveRoomTagApi, uploadApi} from '@/api'
import type {BotAnchorListItem, LiveRoomTag} from '@/types/api'

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

const formRules: FormRules = {
  nickname: [
    {required: true, message: '请输入昵称', trigger: 'blur'},
    {min: 1, max: 32, message: '昵称长度在1-32个字符', trigger: 'blur'}
  ],
  roomTitle: [{max: 128, message: '直播间标题最长128个字符', trigger: 'blur'}],
  category: [{required: true, message: '请选择房间类型', trigger: 'change'}]
}

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
  if (category === LIVE_ROOM_CATEGORY_HOT) return '热门'
  if (category === LIVE_ROOM_CATEGORY_GAME) return '游戏'
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return '私密'
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
    console.error('获取直播间标签失败:', error)
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
  return date.toLocaleString('zh-CN', {hour12: false})
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
    console.error('获取机器人主播列表失败:', error)
    ElMessage.error('获取机器人主播列表失败')
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
  dialogTitle.value = '新增机器人主播'
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
  dialogTitle.value = '编辑机器人主播'
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
    ElMessage.error('只能上传图片文件')
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
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
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
    ElMessage.error('只能上传MP4视频')
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
    // 优先用服务端 URL 预览(大文件 blob URL 常无法播放)
    const previewUrl = res.fileUrl || URL.createObjectURL(file)
    setVideoPreview(previewUrl, !res.fileUrl)
    ElMessage.success('视频上传成功')
  } catch (error) {
    console.error('视频上传失败:', error)
    ElMessage.error('视频上传失败')
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
      ElMessage.success(formData.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    } finally {
      saving.value = false
    }
  })
}

const handleStopLive = async (row: BotAnchorListItem) => {
  try {
    await ElMessageBox.confirm(`确定要让机器人主播「${row.nickname || row.id}」下播吗？`, '确认下播', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await botAnchorApi.stopBotAnchorLive({id: row.id})
    ElMessage.success('下播成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('下播失败:', error)
      ElMessage.error('下播失败')
    }
  }
}

const handleBatchStartLive = async () => {
  const rows = startableSelectedRows.value
  if (!rows.length) {
    ElMessage.warning('请选择可开播的机器人主播')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要批量开播 ${rows.length} 个机器人主播吗？`, '确认批量开播', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    batchOperating.value = true
    const response = await botAnchorApi.batchStartBotAnchorLive({ids: rows.map((row) => row.id)})
    if (response.failCount > 0) {
      ElMessage.warning(`批量开播完成：成功 ${response.successCount} 个，失败 ${response.failCount} 个`)
    } else {
      ElMessage.success(`批量开播成功，共 ${response.successCount} 个`)
    }
    clearSelection()
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量开播失败:', error)
      ElMessage.error('批量开播失败')
    }
  } finally {
    batchOperating.value = false
  }
}

const handleBatchStopLive = async () => {
  const rows = stoppableSelectedRows.value
  if (!rows.length) {
    ElMessage.warning('请选择可下播的机器人主播')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要批量下播 ${rows.length} 个机器人主播吗？`, '确认批量下播', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    batchOperating.value = true
    const response = await botAnchorApi.batchStopBotAnchorLive({ids: rows.map((row) => row.id)})
    if (response.failCount > 0) {
      ElMessage.warning(`批量下播完成：成功 ${response.successCount} 个，失败 ${response.failCount} 个`)
    } else {
      ElMessage.success(`批量下播成功，共 ${response.successCount} 个`)
    }
    clearSelection()
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量下播失败:', error)
      ElMessage.error('批量下播失败')
    }
  } finally {
    batchOperating.value = false
  }
}

const handleStartLive = async (row: BotAnchorListItem) => {
  try {
    await ElMessageBox.confirm(`确定要让机器人主播「${row.nickname || row.id}」开播吗？`, '确认开播', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await botAnchorApi.startBotAnchorLive({id: row.id})
    ElMessage.success('开播成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('开播失败:', error)
      ElMessage.error('开播失败')
    }
  }
}

const toggleStatus = async (row: BotAnchorListItem) => {
  const nextStatus = row.botAnchorStatus === 1 ? 0 : 1
  const actionText = nextStatus === 1 ? '启用' : '停用'
  try {
    await ElMessageBox.confirm(`确定要${actionText}机器人主播「${row.nickname || row.id}」吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await botAnchorApi.setBotAnchorStatus({id: row.id, status: nextStatus})
    ElMessage.success(`${actionText}成功`)
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error)
      ElMessage.error('操作失败')
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

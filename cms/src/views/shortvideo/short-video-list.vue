<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>短视频管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增短视频</el-button>
          <span class="table-tip">支持 App 端与 CMS 端上传；CMS 上传作者类型为 CMS</span>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="标题">
            <el-input v-model="searchForm.title" clearable placeholder="标题(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="作者昵称">
            <el-input v-model="searchForm.authorNickname" clearable placeholder="作者昵称(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.statusFilter" placeholder="全部" style="width: 140px">
              <el-option :value="0" label="全部"/>
              <el-option :value="2" label="只看上架"/>
              <el-option :value="1" label="只看下架"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="标题" prop="title" min-width="140"/>
          <el-table-column label="封面" width="100">
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
          <el-table-column label="视频" min-width="200">
            <template #default="{ row }">
              <div v-if="row.video" class="table-video-cell">
                <video
                    :key="row.video"
                    :src="row.video"
                    class="table-video-preview"
                    controls
                    preload="metadata"
                />
                <el-button link type="primary" @click="openVideoPreview(row.video)">放大预览</el-button>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="作者昵称" min-width="120">
            <template #default="{ row }">
              {{ row.authorNickname || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="作者类型" width="90">
            <template #default="{ row }">
              <el-tag :type="row.authorType === 1 ? 'warning' : 'success'">
                {{ row.authorType === 1 ? 'CMS' : 'App' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="作者ID" prop="authorId" width="100"/>
          <el-table-column label="排序" prop="sort" width="80"/>
          <el-table-column label="是否付费" width="90">
            <template #default="{ row }">
              <el-tag :type="row.isPaid === 1 ? 'warning' : 'success'">
                {{ row.isPaid === 1 ? '付费' : '免费' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="付费钻石" width="110">
            <template #default="{ row }">
              {{ row.isPaid === 1 ? formatPrice(row.payDiamond) : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="免费时长(秒)" width="110">
            <template #default="{ row }">
              {{ row.isPaid === 1 ? (row.freeWatchSeconds != null ? row.freeWatchSeconds : 15) : '-' }}
            </template>
          </el-table-column>
          <el-table-column label="分类" width="100">
            <template #default="{ row }">
              {{ categoryName(row.categoryId) }}
            </template>
          </el-table-column>
          <el-table-column label="来源" width="90">
            <template #default="{ row }">
              {{ sourceLabel(row.source) }}
            </template>
          </el-table-column>
          <el-table-column label="点赞数" prop="likeCount" width="90"/>
          <el-table-column label="观看人数" prop="viewCount" width="90"/>
          <el-table-column label="观看次数" prop="watchCount" width="90"/>
          <el-table-column label="累计钻石收益" width="120">
            <template #default="{ row }">
              {{ formatPrice(row.totalDiamondIncome) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '上架' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="260">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button
                  v-if="row.status !== 1"
                  size="small"
                  type="success"
                  @click="handleOnShelf(row)"
              >
                上架
              </el-button>
              <el-button
                  v-else
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                下架
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
        <el-form-item label="标题" prop="title">
          <el-input v-model="currentRow.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item v-if="isCreateMode" label="视频">
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
                <el-button type="primary">选择视频</el-button>
              </div>
            </el-upload>
            <span v-if="videoDuration > 0" class="form-tip">视频时长：{{ videoDuration }} 秒</span>
            <el-button v-if="videoPreviewUrl" link type="danger" @click="clearCreateVideo">清除视频</el-button>
          </div>
        </el-form-item>
        <el-form-item v-else-if="videoPreviewUrl" label="视频">
          <div class="preview-box">
            <video
                :key="videoPreviewUrl"
                :src="videoPreviewUrl"
                class="dialog-video-preview"
                controls
                preload="metadata"
            />
          </div>
          <div class="form-tip">视频文件不可修改，仅可编辑元数据</div>
        </el-form-item>
        <el-form-item v-if="isCreateMode" label="作者昵称">
          <el-input v-model="currentRow.authorNickname" clearable maxlength="32" placeholder="可选，留空则从随机昵称库抽取"/>
          <div class="form-tip">留空时从随机昵称库（默认英文）自动分配；系统将自动创建 CMS 短视频作者账号</div>
        </el-form-item>
        <el-form-item v-else label="作者昵称">
          <span>{{ currentRow.authorNickname || '-' }}</span>
          <div v-if="currentRow.authorId" class="form-tip">作者ID：{{ currentRow.authorId }}</div>
        </el-form-item>
        <el-form-item label="封面" prop="cover">
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
                <span>上传封面</span>
              </div>
            </el-upload>
            <el-button v-if="coverPreviewUrl || currentRow.cover" link type="danger" @click="clearAsset">
              移除封面
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
        </el-form-item>
        <el-form-item label="是否付费" prop="isPaid">
          <el-radio-group v-model="currentRow.isPaid">
            <el-radio :label="0">免费</el-radio>
            <el-radio :label="1">付费</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="currentRow.isPaid === 1" label="付费钻石" prop="payDiamond">
          <el-input-number
              v-model="currentRow.payDiamond"
              :min="0.0001"
              :precision="4"
              :step="0.0001"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">付费视频一次性解锁价格</span>
        </el-form-item>
        <el-form-item v-if="currentRow.isPaid === 1" label="免费观看时长" prop="freeWatchSeconds">
          <el-input-number
              v-model="currentRow.freeWatchSeconds"
              :min="0"
              :step="1"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">单位：秒，0 表示无免费时长，默认 15 秒</span>
        </el-form-item>
        <el-form-item label="视频分类" prop="categoryId">
          <el-select v-model="currentRow.categoryId" clearable placeholder="请选择分类" style="width: 220px">
            <el-option
                v-for="item in categoryOptions"
                :key="item.id"
                :label="item.name"
                :value="Number(item.id)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="视频来源" prop="source">
          <el-radio-group v-model="currentRow.source">
            <el-radio :label="1">原创</el-radio>
            <el-radio :label="2">转发</el-radio>
            <el-radio :label="3">AI生成</el-radio>
          </el-radio-group>
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
        title="视频预览"
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {shortVideoApi, uploadApi} from '@/api'
import type {ShortVideo, ShortVideoCategory} from '@/types/api.ts'

interface SearchForm {
  title: string
  authorNickname: string
  statusFilter: number
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

const sourceLabelMap: Record<number, string> = {
  1: '原创',
  2: '转发',
  3: 'AI生成',
}

const loading = ref(false)
const saving = ref(false)
const tableData = ref<ShortVideo[]>([])
const categoryOptions = ref<ShortVideoCategory[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  title: '',
  authorNickname: '',
  statusFilter: 0
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

const formatPrice = (price: number) => Number(price || 0).toFixed(4)

const sourceLabel = (source: number) => sourceLabelMap[source] || '-'

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
      reject(new Error('无法读取视频时长'))
    }
    video.src = url
  })
}

const allowedVideoExt = ['.mp4', '.webm', '.mov']

const getFileExt = (name: string) => {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.slice(idx).toLowerCase() : ''
}

const beforeVideoSelect = async (file: File): boolean | Promise<boolean> => {
  const ext = getFileExt(file.name)
  if (!allowedVideoExt.includes(ext)) {
    ElMessage.error('仅支持 MP4 / WebM / MOV 格式')
    return false
  }
  try {
    const duration = await detectVideoDuration(file)
    if (maxVideoDuration.value > 0 && duration > maxVideoDuration.value) {
      ElMessage.error(`视频时长不能超过 ${maxVideoDuration.value} 秒`)
      return false
    }
    videoFile.value = file
    videoDuration.value = duration
    setCreateVideoPreview(file)
  } catch (error) {
    console.error('读取视频失败:', error)
    ElMessage.error('无法读取视频信息')
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
    ElMessage.error('封面只能上传图片文件')
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
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
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

const formRules: FormRules = {
  title: [
    {required: true, message: '请输入标题', trigger: 'blur'},
    {min: 1, max: 64, message: '标题长度在1-64个字符', trigger: 'blur'}
  ],
  payDiamond: [
    {
      validator: (_rule, value, callback) => {
        if (currentRow.value.isPaid === 1 && (!value || value <= 0)) {
          callback(new Error('付费视频请填写付费钻石'))
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
          callback(new Error('免费观看时长不能小于0'))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
  source: [{required: true, message: '请选择视频来源', trigger: 'change'}],
}

const fetchShortVideoCfg = async () => {
  try {
    const response = await shortVideoApi.getShortVideoCfg()
    if (response.cfg?.maxDuration) {
      maxVideoDuration.value = Math.max(1, response.cfg.maxDuration)
    }
  } catch (error) {
    console.error('获取短视频配置失败:', error)
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
    console.error('获取短视频分类失败:', error)
  }
}

const fetchShortVideoList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoList({
      title: searchForm.title,
      authorNickname: searchForm.authorNickname,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取短视频列表失败:', error)
    ElMessage.error('获取短视频列表失败')
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
  dialogTitle.value = '新增短视频'
  currentRow.value = defaultForm()
  resetAssetPreview()
  dialogVisible.value = true
}

const handleEdit = (row: ShortVideo) => {
  dialogTitle.value = '编辑短视频'
  currentRow.value = {
    id: row.id,
    title: row.title,
    cover: row.coverName || '',
    sort: Number(row.sort) || 0,
    isPaid: row.isPaid ?? 0,
    payDiamond: row.isPaid === 1 ? (Number(row.payDiamond) || 1) : 0,
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
    await ElMessageBox.confirm(`确定要删除短视频 "${row.title}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await shortVideoApi.deleteShortVideo(row.id)
    ElMessage.success('删除成功')
    fetchShortVideoList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: ShortVideo) => {
  try {
    await shortVideoApi.onShelfShortVideo(row.id)
    ElMessage.success('上架成功')
    fetchShortVideoList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: ShortVideo) => {
  try {
    await ElMessageBox.confirm(`确定要下架短视频 "${row.title}" 吗？`, '确认下架', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await shortVideoApi.offShelfShortVideo(row.id)
    ElMessage.success('下架成功')
    fetchShortVideoList()
  } catch (error) {
    console.error('下架失败:', error)
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
          ElMessage.error('请选择视频文件')
          return
        }
        ElMessage.info('正在上传视频，请稍候...')
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
        ElMessage.success('创建成功')
      } else {
        await shortVideoApi.updateShortVideo({id: currentRow.value.id, ...payload})
        ElMessage.success('更新成功')
      }
      dialogVisible.value = false
      fetchShortVideoList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    } finally {
      saving.value = false
    }
  })
}

onMounted(() => {
  fetchShortVideoCfg()
  fetchCategoryOptions()
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
  margin-bottom: 15px;
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

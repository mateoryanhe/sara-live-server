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
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="标题">
            <el-input v-model="searchForm.title" clearable placeholder="标题(模糊匹配)"/>
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
          <el-table-column label="视频" min-width="180">
            <template #default="{ row }">
              <video
                  v-if="row.video"
                  :src="row.video"
                  controls
                  preload="metadata"
                  style="width: 160px; max-height: 90px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="排序" prop="sort" width="80"/>
          <el-table-column label="是否付费" width="90">
            <template #default="{ row }">
              <el-tag :type="row.isPaid === 1 ? 'warning' : 'success'">
                {{ row.isPaid === 1 ? '付费' : '免费' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '上架' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="description" min-width="160" show-overflow-tooltip/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="currentRow.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item label="视频" prop="video">
          <div class="upload-wrap">
            <el-upload
                :before-upload="beforeVideoUpload"
                :disabled="videoUploading"
                :http-request="(options) => doUpload(options, 'video')"
                :show-file-list="false"
                accept="video/*,.mp4,.webm,.mov"
                action="#"
            >
              <el-button :loading="videoUploading" type="primary">上传视频</el-button>
            </el-upload>
            <div v-if="videoPreviewUrl" class="preview-box">
              <video :src="videoPreviewUrl" controls preload="metadata" style="width: 100%; max-height: 220px"/>
              <el-button link type="danger" @click="clearAsset('video')">移除视频</el-button>
            </div>
            <div v-else-if="currentRow.video" class="file-name">{{ currentRow.video }}</div>
          </div>
        </el-form-item>
        <el-form-item label="封面" prop="cover">
          <div class="upload-wrap">
            <el-upload
                :before-upload="beforeCoverUpload"
                :disabled="coverUploading"
                :http-request="(options) => doUpload(options, 'cover')"
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
            <el-button v-if="coverPreviewUrl || currentRow.cover" link type="danger" @click="clearAsset('cover')">
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
        <el-form-item label="描述" prop="description">
          <el-input v-model="currentRow.description" placeholder="请输入描述" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {shortVideoApi, uploadApi} from '@/api'
import type {ShortVideo} from '@/types/api.ts'

interface SearchForm {
  title: string
  statusFilter: number
}

interface ShortVideoForm {
  id: string
  title: string
  video: string
  cover: string
  sort: number
  isPaid: number
  description: string
}

const allowedVideoExt = ['.mp4', '.webm', '.mov']

const loading = ref(false)
const tableData = ref<ShortVideo[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  title: '',
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): ShortVideoForm => ({
  id: '',
  title: '',
  video: '',
  cover: '',
  sort: 0,
  isPaid: 0,
  description: ''
})
const currentRow = ref<ShortVideoForm>(defaultForm())
const formRef = ref<FormInstance>()
const videoUploading = ref(false)
const coverUploading = ref(false)
const videoPreviewUrl = ref('')
const coverPreviewUrl = ref('')
const objectPreviewUrls = reactive<{ video: string | null; cover: string | null }>({
  video: null,
  cover: null
})

const getExt = (name: string): string => {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.slice(idx).toLowerCase() : ''
}

const revokeObjectPreview = (field: 'video' | 'cover') => {
  if (objectPreviewUrls[field]) {
    URL.revokeObjectURL(objectPreviewUrls[field]!)
    objectPreviewUrls[field] = null
  }
}

const resetAssetPreview = () => {
  revokeObjectPreview('video')
  revokeObjectPreview('cover')
  videoPreviewUrl.value = ''
  coverPreviewUrl.value = ''
}

const beforeCoverUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('封面只能上传图片文件')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('封面不能超过5MB')
    return false
  }
  return true
}

const beforeVideoUpload = (file: File): boolean => {
  const ext = getExt(file.name)
  if (!allowedVideoExt.includes(ext) && !file.type.startsWith('video/')) {
    ElMessage.error('视频仅支持 mp4、webm、mov 格式')
    return false
  }
  if (file.size > 100 * 1024 * 1024) {
    ElMessage.error('视频不能超过100MB')
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions, field: 'video' | 'cover') => {
  const file = options.file as File
  const uploading = field === 'video' ? videoUploading : coverUploading
  uploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value[field] = res.fileName
    const objectUrl = URL.createObjectURL(file)
    if (field === 'video') {
      revokeObjectPreview('video')
      videoPreviewUrl.value = objectUrl
      objectPreviewUrls.video = objectUrl
    } else {
      revokeObjectPreview('cover')
      coverPreviewUrl.value = objectUrl
      objectPreviewUrls.cover = objectUrl
    }
    formRef.value?.validateField(field).catch(() => undefined)
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

const clearAsset = (field: 'video' | 'cover') => {
  currentRow.value[field] = ''
  if (field === 'video') {
    revokeObjectPreview('video')
    videoPreviewUrl.value = ''
  } else {
    revokeObjectPreview('cover')
    coverPreviewUrl.value = ''
  }
  formRef.value?.validateField(field).catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAssetPreview()
  }
})

const formRules: FormRules = {
  title: [
    {required: true, message: '请输入标题', trigger: 'blur'},
    {min: 1, max: 64, message: '标题长度在1-64个字符', trigger: 'blur'}
  ],
  video: [{required: true, message: '请上传视频', trigger: 'change'}],
  description: [{max: 255, message: '描述最长255字符', trigger: 'blur'}]
}

const fetchShortVideoList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoList({
      title: searchForm.title,
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
    video: row.videoName || '',
    cover: row.coverName || '',
    sort: Number(row.sort) || 0,
    isPaid: row.isPaid ?? 0,
    description: row.description || ''
  }
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
    try {
      const payload = {
        title: currentRow.value.title,
        video: currentRow.value.video,
        cover: currentRow.value.cover,
        sort: currentRow.value.sort,
        isPaid: currentRow.value.isPaid,
        description: currentRow.value.description
      }
      if (currentRow.value.id) {
        await shortVideoApi.updateShortVideo({id: currentRow.value.id, ...payload})
      } else {
        await shortVideoApi.createShortVideo(payload)
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchShortVideoList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    }
  })
}

onMounted(() => {
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
}

.file-name {
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

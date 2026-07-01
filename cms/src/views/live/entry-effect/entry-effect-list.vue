<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>进场特效</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增进场特效</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="名称">
            <el-input v-model="searchForm.name" clearable placeholder="名称(模糊匹配)"/>
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
          <el-table-column label="名称" prop="name" min-width="140"/>
          <el-table-column label="等级开始" prop="levelStart" width="100"/>
          <el-table-column label="等级结束" prop="levelEnd" width="100"/>
          <el-table-column label="动画资源" min-width="220">
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
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '上架' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入名称"/>
        </el-form-item>
        <el-form-item label="等级开始" prop="levelStart">
          <el-input-number v-model="currentRow.levelStart" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item label="等级结束" prop="levelEnd">
          <el-input-number v-model="currentRow.levelEnd" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item label="动画资源" prop="animation">
          <div class="asset-upload-wrap">
            <el-upload
                :before-upload="beforeAnimationUpload"
                :disabled="animationUploading"
                :http-request="doUpload"
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
                <span>点击上传动画资源</span>
              </div>
            </el-upload>
            <el-button
                v-if="animationPreviewType !== 'none' || currentRow.animation"
                link
                type="danger"
                @click="clearAnimation"
            >
              移除动画
            </el-button>
          </div>
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
import {Document, Plus} from '@element-plus/icons-vue'
import {entryEffectApi, uploadApi} from '@/api'
import type {EntryEffect} from '@/types/api.ts'
import {
  getExt,
  isImageUrl,
  isVideoUrl,
  resolveFilePreviewType,
  resolveMediaPreviewType,
  type MediaPreviewType
} from '@/utils/media-preview'

interface SearchForm {
  name: string
  statusFilter: number
}

interface EntryEffectForm {
  id: string
  name: string
  levelStart: number
  levelEnd: number
  animation: string
}

const loading = ref(false)
const tableData = ref<EntryEffect[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): EntryEffectForm => ({
  id: '',
  name: '',
  levelStart: 1,
  levelEnd: 1,
  animation: ''
})
const currentRow = ref<EntryEffectForm>(defaultForm())
const formRef = ref<FormInstance>()

const animationUploading = ref(false)
const animationPreviewUrl = ref('')
const animationPreviewType = ref<MediaPreviewType>('none')
const animationFileLabel = ref('')
let objectPreviewUrl = ''

const allowedAnimationExt = [
  '.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.apng',
  '.svga', '.pag', '.json', '.lottie', '.mp4', '.webm', '.zip'
]

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = ''
  }
}

const resetAnimationPreview = () => {
  revokeObjectPreview()
  animationPreviewUrl.value = ''
  animationPreviewType.value = 'none'
  animationFileLabel.value = ''
}

const setAnimationPreview = (
    url: string,
    fileLabel: string,
    type?: MediaPreviewType,
    fromObject = false
) => {
  revokeObjectPreview()
  const previewType = type ?? resolveMediaPreviewType(url, fileLabel)
  animationPreviewType.value = previewType
  animationPreviewUrl.value = previewType === 'image' || previewType === 'video' ? url : ''
  animationFileLabel.value = previewType === 'file' ? fileLabel : ''
  if (fromObject && url) {
    objectPreviewUrl = url
  }
}

const clearAnimation = () => {
  currentRow.value.animation = ''
  resetAnimationPreview()
  formRef.value?.validateField('animation').catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAnimationPreview()
  }
})

const beforeAnimationUpload = (file: File): boolean => {
  const ext = getExt(file.name)
  if (!allowedAnimationExt.includes(ext)) {
    ElMessage.error(`不支持的文件类型: ${ext || '未知'}`)
    return false
  }
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.error('文件不能超过50MB')
    return false
  }
  return true
}

const doUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  animationUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.animation = res.fileName
    const previewType = resolveFilePreviewType(file.name)
    if (previewType === 'image' || previewType === 'video') {
      setAnimationPreview(URL.createObjectURL(file), '', previewType, true)
    } else {
      setAnimationPreview('', res.fileName, 'file')
    }
    formRef.value?.validateField('animation').catch(() => undefined)
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    animationUploading.value = false
  }
}

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入名称', trigger: 'blur'},
    {min: 1, max: 64, message: '名称长度在1-64个字符', trigger: 'blur'}
  ],
  levelStart: [{required: true, message: '请输入等级开始', trigger: 'change'}],
  levelEnd: [
    {required: true, message: '请输入等级结束', trigger: 'change'},
    {
      validator: (_rule, value, callback) => {
        if (value < currentRow.value.levelStart) {
          callback(new Error('等级结束不能小于等级开始'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  animation: [{required: true, message: '请上传动画资源', trigger: 'change'}]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await entryEffectApi.getEntryEffectList({
      name: searchForm.name,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取进场特效列表失败:', error)
    ElMessage.error('获取进场特效列表失败')
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
  dialogTitle.value = '新增进场特效'
  currentRow.value = defaultForm()
  resetAnimationPreview()
  dialogVisible.value = true
}

const handleEdit = (row: EntryEffect) => {
  dialogTitle.value = '编辑进场特效'
  const animationName = row.animationName || ''
  currentRow.value = {
    id: row.id,
    name: row.name,
    levelStart: Number(row.levelStart) || 1,
    levelEnd: Number(row.levelEnd) || 1,
    animation: animationName
  }
  if (animationName) {
    const previewType = resolveMediaPreviewType(row.animation || '', animationName)
    if (previewType === 'image' || previewType === 'video') {
      setAnimationPreview(row.animation || '', '', previewType)
    } else {
      setAnimationPreview('', animationName, 'file')
    }
  } else {
    resetAnimationPreview()
  }
  dialogVisible.value = true
}

const handleDelete = async (row: EntryEffect) => {
  try {
    await ElMessageBox.confirm(`确定要删除进场特效 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await entryEffectApi.deleteEntryEffect(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: EntryEffect) => {
  try {
    await entryEffectApi.onShelfEntryEffect(row.id)
    ElMessage.success('上架成功')
    fetchList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: EntryEffect) => {
  try {
    await ElMessageBox.confirm(`确定要下架进场特效 "${row.name}" 吗？`, '确认下架', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await entryEffectApi.offShelfEntryEffect(row.id)
    ElMessage.success('下架成功')
    fetchList()
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
        name: currentRow.value.name,
        levelStart: currentRow.value.levelStart,
        levelEnd: currentRow.value.levelEnd,
        animation: currentRow.value.animation
      }
      if (currentRow.value.id) {
        await entryEffectApi.updateEntryEffect({id: currentRow.value.id, ...payload})
      } else {
        await entryEffectApi.createEntryEffect(payload)
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
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

.animation-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

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

.animation-placeholder {
  width: 240px;
  height: 120px;
}

.asset-uploader-icon {
  font-size: 28px;
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

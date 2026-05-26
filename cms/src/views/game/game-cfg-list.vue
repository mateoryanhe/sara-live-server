<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>游戏配置管理</span>
          <el-button type="primary" @click="handleAdd">新增游戏</el-button>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="游戏名称">
          <el-input v-model="searchForm.name" clearable placeholder="游戏名称(模糊匹配)"/>
        </el-form-item>
        <el-form-item label="游戏编码">
          <el-input v-model="searchForm.code" clearable placeholder="游戏编码(模糊匹配)"/>
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
        <el-table-column label="游戏名称" min-width="140" prop="name"/>
        <el-table-column label="游戏编码" min-width="140" prop="code"/>
        <el-table-column label="直播间封面" min-width="120">
          <template #default="{ row }">
            <el-image
                v-if="row.liveCoverUrl"
                :preview-src-list="[row.liveCoverUrl]"
                :src="row.liveCoverUrl"
                fit="cover"
                preview-teleported
                style="width: 72px; height: 48px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="跳转链接" min-width="180" prop="link" show-overflow-tooltip>
          <template #default="{ row }">{{ row.link || '-' }}</template>
        </el-table-column>
        <el-table-column label="排序" prop="sort" width="80"/>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createdAt" width="160"/>
        <el-table-column label="更新时间" prop="updatedAt" width="160"/>
        <el-table-column fixed="right" label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
        <el-form-item label="游戏名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入游戏名称"/>
        </el-form-item>
        <el-form-item label="游戏编码" prop="code">
          <el-input v-model="currentRow.code" placeholder="请输入唯一游戏编码"/>
        </el-form-item>
        <el-form-item label="直播间封面" prop="liveCover">
          <div class="cover-upload-wrap">
            <el-upload
                :before-upload="beforeCoverUpload"
                :disabled="coverUploading"
                :http-request="doCoverUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="cover-uploader"
            >
              <img v-if="coverPreviewUrl" :src="coverPreviewUrl" alt="cover" class="cover-preview"/>
              <div v-else class="cover-uploader-placeholder">
                <el-icon class="cover-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>点击上传封面</span>
              </div>
            </el-upload>
            <el-button
                v-if="coverPreviewUrl"
                class="cover-remove-btn"
                link
                type="danger"
                @click="clearCover"
            >
              移除封面
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="跳转链接" prop="link">
          <el-input v-model="currentRow.link" placeholder="请输入跳转链接"/>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
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
import {gameCfgApi, uploadApi} from '@/api'
import type {GameCfg} from '@/types/api.ts'

interface SearchForm {
  name: string
  code: string
  statusFilter: number
}

interface GameCfgForm {
  id: string
  name: string
  code: string
  liveCover: string
  link: string
  sort: number
  status: number
}

const loading = ref(false)
const tableData = ref<GameCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  code: '',
  statusFilter: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): GameCfgForm => ({
  id: '',
  name: '',
  code: '',
  liveCover: '',
  link: '',
  sort: 0,
  status: 1,
})
const currentRow = ref<GameCfgForm>(defaultForm())
const formRef = ref<FormInstance>()
const coverUploading = ref(false)
const coverPreviewUrl = ref('')
let objectPreviewUrl: string | null = null

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = null
  }
}

const setCoverPreview = (url: string, fromObject = false) => {
  revokeObjectPreview()
  coverPreviewUrl.value = url
  if (fromObject) {
    objectPreviewUrl = url
  }
}

const clearCover = () => {
  currentRow.value.liveCover = ''
  setCoverPreview('')
  formRef.value?.validateField('liveCover').catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    coverPreviewUrl.value = ''
  }
})

const beforeCoverUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片不能超过5MB')
    return false
  }
  return true
}

const doCoverUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  coverUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.liveCover = res.fileName
    setCoverPreview(URL.createObjectURL(file), true)
    formRef.value?.validateField('liveCover').catch(() => undefined)
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    coverUploading.value = false
  }
}

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入游戏名称', trigger: 'blur'},
    {min: 1, max: 64, message: '游戏名称长度在1-64个字符', trigger: 'blur'},
  ],
  code: [
    {required: true, message: '请输入游戏编码', trigger: 'blur'},
    {min: 1, max: 64, message: '游戏编码长度在1-64个字符', trigger: 'blur'},
  ],
  liveCover: [{required: true, message: '请上传直播间游戏封面', trigger: 'change'}],
  status: [{required: true, message: '请选择状态', trigger: 'change'}],
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await gameCfgApi.getGameCfgList({
      name: searchForm.name,
      code: searchForm.code,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取游戏配置列表失败:', error)
    ElMessage.error('获取游戏配置列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.code = ''
  searchForm.statusFilter = 0
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
  dialogTitle.value = '新增游戏'
  currentRow.value = defaultForm()
  setCoverPreview('')
  dialogVisible.value = true
}

const handleEdit = (row: GameCfg) => {
  dialogTitle.value = '编辑游戏'
  currentRow.value = {
    id: row.id,
    name: row.name,
    code: row.code,
    liveCover: row.liveCover || '',
    link: row.link || '',
    sort: Number(row.sort) || 0,
    status: Number(row.status) || 0,
  }
  setCoverPreview(row.liveCoverUrl || '')
  dialogVisible.value = true
}

const handleDelete = async (row: GameCfg) => {
  try {
    await ElMessageBox.confirm(`确定要删除游戏 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await gameCfgApi.deleteGameCfg(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    try {
      const payload = {
        name: currentRow.value.name.trim(),
        code: currentRow.value.code.trim(),
        liveCover: currentRow.value.liveCover.trim(),
        link: currentRow.value.link.trim(),
        sort: currentRow.value.sort,
        status: currentRow.value.status,
      }
      if (currentRow.value.id) {
        await gameCfgApi.updateGameCfg({id: currentRow.value.id, ...payload})
      } else {
        await gameCfgApi.createGameCfg(payload)
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存失败:', error)
    }
  })
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  font-weight: bold;
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

.cover-upload-wrap {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.cover-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.cover-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.cover-uploader-placeholder {
  width: 160px;
  height: 96px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
}

.cover-uploader-icon {
  font-size: 24px;
  margin-bottom: 6px;
}

.cover-preview {
  width: 160px;
  height: 96px;
  object-fit: cover;
  display: block;
}

.cover-remove-btn {
  padding: 0;
}
</style>

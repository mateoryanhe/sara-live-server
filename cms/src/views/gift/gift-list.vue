<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>礼物管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增礼物</el-button>
        </div>

        <!-- 搜索表单 -->
        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="礼物名称">
            <el-input v-model="searchForm.name" clearable placeholder="礼物名称(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="分类">
            <el-input v-model="searchForm.category" clearable placeholder="分类"/>
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
          <el-table-column label="礼物名称" prop="name" min-width="120"/>
          <el-table-column label="图标" width="90">
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
          <el-table-column label="动画资源" prop="animation" min-width="180" show-overflow-tooltip/>
          <el-table-column label="价格(钻石)" prop="price" width="110"/>
          <el-table-column label="分类" prop="category" width="120"/>
          <el-table-column label="排序" prop="sort" width="80"/>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '上架' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="description" min-width="160" show-overflow-tooltip/>
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

    <!-- 礼物编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item label="礼物名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入礼物名称"/>
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <div class="upload-field">
            <el-input v-model="currentRow.icon" placeholder="点击右侧按钮上传或手动输入文件名"/>
            <el-upload
                :before-upload="beforeIconUpload"
                :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'icon')"
                :show-file-list="false"
                accept="image/*"
                action="#"
            >
              <el-button :loading="iconUploading" type="primary">上传图标</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="动画资源" prop="animation">
          <div class="upload-field">
            <el-input v-model="currentRow.animation" placeholder="点击右侧按钮上传或手动输入文件名"/>
            <el-upload
                :before-upload="beforeAnimationUpload"
                :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'animation')"
                :show-file-list="false"
                accept=".svga,.pag,.json,.lottie,.mp4,.webm,.zip,.gif,.apng,.png,.webp"
                action="#"
            >
              <el-button :loading="animationUploading" type="primary">上传动画</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="价格(钻石)" prop="price">
          <el-input-number v-model="currentRow.price" :min="0" controls-position="right"/>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="currentRow.category" placeholder="请输入分类"/>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {giftApi, uploadApi} from '@/api'
import type {Gift} from '@/types/api'

interface SearchForm {
  name: string
  category: string
  statusFilter: number
}

interface GiftForm {
  id: string
  name: string
  icon: string
  animation: string
  price: number
  category: string
  sort: number
  description: string
}

const loading = ref(false)
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

// CMS后台允许上传的扩展名(与后端 allowedCMSExt 保持一致)
const allowedAnimationExt = [
  '.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.apng',
  '.svga', '.pag', '.json', '.lottie', '.mp4', '.webm', '.zip'
]

const getExt = (name: string): string => {
  const idx = name.lastIndexOf('.')
  return idx >= 0 ? name.slice(idx).toLowerCase() : ''
}

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('图标只能上传图片文件')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图标文件不能超过5MB')
    return false
  }
  return true
}

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

const doUpload = async (
    options: UploadRequestOptions,
    field: 'icon' | 'animation'
) => {
  const flag = field === 'icon' ? iconUploading : animationUploading
  flag.value = true
  try {
    const res = await uploadApi.uploadFile(options.file as File)
    currentRow.value[field] = res.fileName
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    flag.value = false
  }
}

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入礼物名称', trigger: 'blur'},
    {min: 1, max: 64, message: '礼物名称长度在1-64个字符', trigger: 'blur'}
  ],
  icon: [
    {max: 255, message: '图标URL最长255字符', trigger: 'blur'}
  ],
  animation: [
    {max: 255, message: '动画URL最长255字符', trigger: 'blur'}
  ],
  category: [
    {max: 32, message: '分类最长32字符', trigger: 'blur'}
  ],
  description: [
    {max: 255, message: '描述最长255字符', trigger: 'blur'}
  ]
}

// 获取礼物列表
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
    console.error('获取礼物列表失败:', error)
    ElMessage.error('获取礼物列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchGiftList()
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGiftList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGiftList()
}

// 操作处理
const handleAdd = () => {
  dialogTitle.value = '新增礼物'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: Gift) => {
  dialogTitle.value = '编辑礼物'
  currentRow.value = {
    id: row.id,
    name: row.name,
    icon: row.icon,
    animation: row.animation,
    price: Number(row.price) || 0,
    category: row.category,
    sort: Number(row.sort) || 0,
    description: row.description
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(`确定要删除礼物 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await giftApi.deleteGift(row.id)
    ElMessage.success('删除成功')
    fetchGiftList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: Gift) => {
  try {
    await giftApi.onShelfGift(row.id)
    ElMessage.success('上架成功')
    fetchGiftList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: Gift) => {
  try {
    await ElMessageBox.confirm(`确定要下架礼物 "${row.name}" 吗？`, '确认下架', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await giftApi.offShelfGift(row.id)
    ElMessage.success('下架成功')
    fetchGiftList()
  } catch (error) {
    console.error('下架失败:', error)
  }
}

// 保存操作
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentRow.value.id) {
          await giftApi.updateGift(currentRow.value)
        } else {
          const {name, icon, animation, price, category, sort, description} = currentRow.value
          await giftApi.createGift({name, icon, animation, price, category, sort, description})
        }

        ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchGiftList()
      } catch (error) {
        console.error(currentRow.value.id ? '更新失败:' : '创建失败:', error)
        ElMessage.error(currentRow.value.id ? '更新失败' : '创建失败')
      }
    }
  })
}

// 重置搜索
const resetSearch = () => {
  searchForm.name = ''
  searchForm.category = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchGiftList()
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

.upload-field {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.upload-field .el-input {
  flex: 1;
}
</style>

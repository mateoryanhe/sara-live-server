<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>首页 Banner 管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增 Banner</el-button>
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
          <el-table-column label="图片" width="100">
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
          <el-table-column label="跳转链接" prop="link" min-width="200" show-overflow-tooltip/>
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
        <el-form-item label="标题" prop="title">
          <el-input v-model="currentRow.title" placeholder="请输入标题"/>
        </el-form-item>
        <el-form-item label="图片" prop="image">
          <div class="upload-field">
            <el-input v-model="currentRow.image" placeholder="点击右侧按钮上传或手动输入文件名"/>
            <el-upload
                :before-upload="beforeImageUpload"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
            >
              <el-button :loading="imageUploading" type="primary">上传图片</el-button>
            </el-upload>
          </div>
        </el-form-item>
        <el-form-item label="跳转链接" prop="link">
          <el-input v-model="currentRow.link" placeholder="请输入跳转链接"/>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
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
import {bannerApi, uploadApi} from '@/api'
import type {Banner} from '@/types/api'

interface SearchForm {
  title: string
  statusFilter: number
}

interface BannerForm {
  id: string
  title: string
  image: string
  link: string
  sort: number
}

const loading = ref(false)
const tableData = ref<Banner[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  title: '',
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): BannerForm => ({
  id: '',
  title: '',
  image: '',
  link: '',
  sort: 0
})
const currentRow = ref<BannerForm>(defaultForm())
const formRef = ref<FormInstance>()
const imageUploading = ref(false)

const beforeImageUpload = (file: File): boolean => {
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

const doUpload = async (options: UploadRequestOptions) => {
  imageUploading.value = true
  try {
    const res = await uploadApi.uploadFile(options.file as File)
    currentRow.value.image = res.fileName
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    imageUploading.value = false
  }
}

const formRules: FormRules = {
  title: [
    {required: true, message: '请输入标题', trigger: 'blur'},
    {min: 1, max: 64, message: '标题长度在1-64个字符', trigger: 'blur'}
  ],
  image: [{max: 255, message: '图片资源名最长255字符', trigger: 'blur'}],
  link: [{max: 512, message: '跳转链接最长512字符', trigger: 'blur'}]
}

const fetchBannerList = async () => {
  loading.value = true
  try {
    const response = await bannerApi.getBannerList({
      title: searchForm.title,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取Banner列表失败:', error)
    ElMessage.error('获取Banner列表失败')
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
  dialogTitle.value = '新增 Banner'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: Banner) => {
  dialogTitle.value = '编辑 Banner'
  currentRow.value = {
    id: row.id,
    title: row.title,
    image: row.image,
    link: row.link,
    sort: Number(row.sort) || 0
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Banner) => {
  try {
    await ElMessageBox.confirm(`确定要删除 Banner "${row.title}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await bannerApi.deleteBanner(row.id)
    ElMessage.success('删除成功')
    fetchBannerList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: Banner) => {
  try {
    await bannerApi.onShelfBanner(row.id)
    ElMessage.success('上架成功')
    fetchBannerList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: Banner) => {
  try {
    await ElMessageBox.confirm(`确定要下架 Banner "${row.title}" 吗？`, '确认下架', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await bannerApi.offShelfBanner(row.id)
    ElMessage.success('下架成功')
    fetchBannerList()
  } catch (error) {
    console.error('下架失败:', error)
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
        const {title, image, link, sort} = currentRow.value
        await bannerApi.createBanner({title, image, link, sort})
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchBannerList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    }
  })
}

const resetSearch = () => {
  searchForm.title = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchBannerList()
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

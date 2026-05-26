<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>充值配置管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增充值档位</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="档位名称">
            <el-input v-model="searchForm.name" clearable placeholder="名称(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="类型">
            <el-select v-model="searchForm.typeFilter" placeholder="全部" style="width: 140px">
              <el-option :value="0" label="全部"/>
              <el-option v-for="item in cfgTypeOptions" :key="item.value" :label="item.label" :value="item.value"/>
            </el-select>
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
          <el-table-column label="档位名称" prop="name" min-width="120"/>
          <el-table-column label="类型" width="100">
            <template #default="{ row }">
              {{ cfgTypeLabel(row.cfgType) }}
            </template>
          </el-table-column>
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
          <el-table-column label="基础金币" prop="gold" width="100"/>
          <el-table-column label="赠送金币" prop="extraGold" width="100"/>
          <el-table-column label="合计金币" width="100">
            <template #default="{ row }">
              {{ totalGold(row) }}
            </template>
          </el-table-column>
          <el-table-column label="价格" width="120">
            <template #default="{ row }">
              {{ formatPrice(row.price) }}
            </template>
          </el-table-column>
          <el-table-column label="商品SKU" prop="productId" min-width="140" show-overflow-tooltip/>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="640px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
        <el-form-item label="档位名称" prop="name">
          <el-input v-model="currentRow.name" placeholder="请输入档位名称"/>
        </el-form-item>
        <el-form-item label="类型" prop="cfgType">
          <el-select v-model="currentRow.cfgType" placeholder="请选择类型" style="width: 100%">
            <el-option v-for="item in cfgTypeOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <div class="icon-upload-wrap">
            <el-upload
                :before-upload="beforeIconUpload"
                :disabled="iconUploading"
                :http-request="doUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="icon-uploader"
            >
              <img v-if="iconPreviewUrl" :src="iconPreviewUrl" alt="icon" class="icon-preview"/>
              <div v-else class="icon-uploader-placeholder">
                <el-icon class="icon-uploader-icon">
                  <Plus/>
                </el-icon>
                <span>点击上传图标</span>
              </div>
            </el-upload>
            <el-button
                v-if="iconPreviewUrl || currentRow.icon"
                link
                type="danger"
                @click="clearIcon"
            >
              移除图标
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="基础金币" prop="gold">
          <el-input-number v-model="currentRow.gold" :min="1" controls-position="right"/>
        </el-form-item>
        <el-form-item label="赠送金币" prop="extraGold">
          <el-input-number v-model="currentRow.extraGold" :min="0" controls-position="right"/>
          <div class="form-tip">充值成功后玩家实际到账金币 = 基础金币 + 赠送金币</div>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="currentRow.price" :min="0.0001" :precision="4" :step="0.01" controls-position="right"/>
          <div class="form-tip">充值货币固定，支持4位小数</div>
        </el-form-item>
        <el-form-item label="商品SKU" prop="productId">
          <el-input v-model="currentRow.productId" placeholder="第三方商品SKU(可选)"/>
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
import {onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {rechargeCfgApi, uploadApi} from '@/api'
import type {RechargeCfg} from '@/types/api.ts'

interface SearchForm {
  name: string
  typeFilter: number
  statusFilter: number
}

interface RechargeCfgForm {
  id: string
  name: string
  cfgType: number
  icon: string
  gold: number
  extraGold: number
  price: number
  productId: string
  sort: number
  description: string
}

const cfgTypeOptions = [
  {value: 1, label: 'iOS'},
  {value: 2, label: 'Google'},
  {value: 3, label: '渠道'}
]

const cfgTypeLabel = (cfgType: number) => {
  return cfgTypeOptions.find((item) => item.value === cfgType)?.label ?? '未知'
}

const loading = ref(false)
const tableData = ref<RechargeCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: '',
  typeFilter: 0,
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): RechargeCfgForm => ({
  id: '',
  name: '',
  cfgType: 1,
  icon: '',
  gold: 1,
  extraGold: 0,
  price: 0.99,
  productId: '',
  sort: 0,
  description: ''
})
const currentRow = ref<RechargeCfgForm>(defaultForm())
const formRef = ref<FormInstance>()
const iconUploading = ref(false)
const iconPreviewUrl = ref('')
let objectPreviewUrl: string | null = null

const formatPrice = (price: number) => {
  return `${Number(price).toFixed(4)}`
}

const totalGold = (row: RechargeCfg) => {
  return (Number(row.gold) || 0) + (Number(row.extraGold) || 0)
}

const revokeObjectPreview = () => {
  if (objectPreviewUrl) {
    URL.revokeObjectURL(objectPreviewUrl)
    objectPreviewUrl = null
  }
}

const setIconPreview = (url: string, fromObject = false) => {
  revokeObjectPreview()
  iconPreviewUrl.value = url
  if (fromObject) {
    objectPreviewUrl = url
  }
}

const clearIcon = () => {
  currentRow.value.icon = ''
  setIconPreview('')
  formRef.value?.validateField('icon').catch(() => undefined)
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    revokeObjectPreview()
    iconPreviewUrl.value = ''
  }
})

const beforeIconUpload = (file: File): boolean => {
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
  const file = options.file as File
  iconUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value.icon = res.fileName
    setIconPreview(URL.createObjectURL(file), true)
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    iconUploading.value = false
  }
}

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入档位名称', trigger: 'blur'},
    {min: 1, max: 64, message: '名称长度在1-64个字符', trigger: 'blur'}
  ],
  cfgType: [{required: true, message: '请选择类型', trigger: 'change'}],
  gold: [{required: true, message: '请输入基础金币数', trigger: 'change'}],
  price: [{required: true, message: '请输入价格', trigger: 'change'}],
  productId: [{max: 64, message: '商品SKU最长64字符', trigger: 'blur'}],
  description: [{max: 255, message: '描述最长255字符', trigger: 'blur'}]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await rechargeCfgApi.getRechargeCfgList({
      name: searchForm.name,
      typeFilter: searchForm.typeFilter,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取充值配置列表失败:', error)
    ElMessage.error('获取充值配置列表失败')
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
  dialogTitle.value = '新增充值档位'
  currentRow.value = defaultForm()
  setIconPreview('')
  dialogVisible.value = true
}

const handleEdit = (row: RechargeCfg) => {
  dialogTitle.value = '编辑充值档位'
  currentRow.value = {
    id: row.id,
    name: row.name,
    cfgType: Number(row.cfgType) || 1,
    icon: row.iconName || '',
    gold: Number(row.gold) || 1,
    extraGold: Number(row.extraGold) || 0,
    price: Number(row.price) || 0.99,
    productId: row.productId || '',
    sort: Number(row.sort) || 0,
    description: row.description || ''
  }
  setIconPreview(row.icon || '')
  dialogVisible.value = true
}

const handleDelete = async (row: RechargeCfg) => {
  try {
    await ElMessageBox.confirm(`确定要删除充值档位 "${row.name}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await rechargeCfgApi.deleteRechargeCfg(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: RechargeCfg) => {
  try {
    await rechargeCfgApi.onShelfRechargeCfg(row.id)
    ElMessage.success('上架成功')
    fetchList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: RechargeCfg) => {
  try {
    await ElMessageBox.confirm(`确定要下架充值档位 "${row.name}" 吗？`, '确认下架', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await rechargeCfgApi.offShelfRechargeCfg(row.id)
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
        cfgType: currentRow.value.cfgType,
        icon: currentRow.value.icon,
        gold: currentRow.value.gold,
        extraGold: currentRow.value.extraGold,
        price: currentRow.value.price,
        productId: currentRow.value.productId,
        sort: currentRow.value.sort,
        description: currentRow.value.description
      }
      if (currentRow.value.id) {
        await rechargeCfgApi.updateRechargeCfg({id: currentRow.value.id, ...payload})
      } else {
        await rechargeCfgApi.createRechargeCfg(payload)
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
  searchForm.typeFilter = 0
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

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.icon-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.icon-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.icon-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.icon-uploader-placeholder {
  width: 96px;
  height: 96px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.icon-uploader-icon {
  font-size: 28px;
}

.icon-preview {
  width: 96px;
  height: 96px;
  display: block;
  object-fit: cover;
}
</style>

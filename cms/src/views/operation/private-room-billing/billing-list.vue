<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>私密直播间计费</span>
          <span class="card-tip">按分钟计费</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增计费</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
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
          <el-table-column label="每分钟钻石价格" width="160">
            <template #default="{ row }">
              {{ formatPrice(row.pricePerMinute) }}
            </template>
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
              :page-sizes="[10, 20]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="140px">
        <el-form-item label="每分钟钻石价格" prop="pricePerMinute">
          <el-input-number
              v-model="currentRow.pricePerMinute"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="currentRow.sort" controls-position="right"/>
          <div class="form-tip">数值越大越靠前</div>
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
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {privateRoomBillingApi} from '@/api'
import type {PrivateRoomBilling} from '@/types/api.ts'

interface SearchForm {
  statusFilter: number
}

interface BillingForm {
  id: string
  pricePerMinute: number
  sort: number
}

const loading = ref(false)
const tableData = ref<PrivateRoomBilling[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  statusFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): BillingForm => ({
  id: '',
  pricePerMinute: 0,
  sort: 0
})
const currentRow = ref<BillingForm>(defaultForm())
const formRef = ref<FormInstance>()

const formatPrice = (price: number) => Number(price).toFixed(4)

const formRules: FormRules = {
  pricePerMinute: [{required: true, message: '请输入每分钟钻石价格', trigger: 'blur'}]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await privateRoomBillingApi.getBillingList({
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取计费列表失败:', error)
    ElMessage.error('获取计费列表失败')
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
  dialogTitle.value = '新增计费'
  currentRow.value = defaultForm()
  dialogVisible.value = true
}

const handleEdit = (row: PrivateRoomBilling) => {
  dialogTitle.value = '编辑计费'
  currentRow.value = {
    id: row.id,
    pricePerMinute: Number(row.pricePerMinute) || 0,
    sort: Number(row.sort) || 0
  }
  dialogVisible.value = true
}

const handleDelete = async (row: PrivateRoomBilling) => {
  try {
    await ElMessageBox.confirm(`确定要删除计费配置 ID ${row.id} 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await privateRoomBillingApi.deleteBilling(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleOnShelf = async (row: PrivateRoomBilling) => {
  try {
    await privateRoomBillingApi.onShelfBilling(row.id)
    ElMessage.success('上架成功')
    fetchList()
  } catch (error) {
    console.error('上架失败:', error)
    ElMessage.error('上架失败')
  }
}

const handleOffShelf = async (row: PrivateRoomBilling) => {
  try {
    await privateRoomBillingApi.offShelfBilling(row.id)
    ElMessage.success('下架成功')
    fetchList()
  } catch (error) {
    console.error('下架失败:', error)
    ElMessage.error('下架失败')
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        pricePerMinute: currentRow.value.pricePerMinute,
        sort: currentRow.value.sort
      }
      if (currentRow.value.id) {
        await privateRoomBillingApi.updateBilling({id: currentRow.value.id, ...payload})
      } else {
        await privateRoomBillingApi.createBilling(payload)
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  font-weight: bold;
}

.card-tip {
  font-size: 13px;
  font-weight: normal;
  color: var(--el-text-color-secondary);
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
</style>

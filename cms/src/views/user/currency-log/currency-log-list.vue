<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" clearable placeholder="请输入用户ID"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="流水ID" min-width="180" prop="id"/>
        <el-table-column label="用户ID" min-width="180" prop="userId"/>
        <el-table-column label="用户昵称" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="变动类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.action === 1 ? 'success' : 'danger'">
              {{ actionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="amountColumnLabel" prop="amount" width="120">
          <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="变动前" prop="before" width="120">
          <template #default="{ row }">{{ formatAmount(row.before) }}</template>
        </el-table-column>
        <el-table-column label="变动后" prop="after" width="120">
          <template #default="{ row }">{{ formatAmount(row.after) }}</template>
        </el-table-column>
        <el-table-column label="原因" min-width="140" prop="reasonText" show-overflow-tooltip>
          <template #default="{ row }">{{ row.reasonText || '-' }}</template>
        </el-table-column>
        <el-table-column label="时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
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
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {currencyLogApi} from '@/api'
import type {CurrencyLogItem} from '@/types/api'

const route = useRoute()

const currencyType = computed(() => Number(route.meta.currencyType) || 1)
const pageTitle = computed(() => String(route.meta.title || '货币流水'))
const amountColumnLabel = computed(() => currencyType.value === 2 ? '钻石变动' : '金币变动')

const loading = ref(false)
const tableData = ref<CurrencyLogItem[]>([])

const searchForm = reactive({
  userId: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await currencyLogApi.getCurrencyLogList({
      userId: searchForm.userId.trim(),
      currencyType: currencyType.value,
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取流水失败:', error)
    ElMessage.error('获取流水失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.userId = ''
  pagination.pageIndex = 1
  fetchList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const formatAmount = (value: number | null | undefined) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(2)
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString('zh-CN')
  } catch {
    return '-'
  }
}

const actionLabel = (action: number) => {
  return action === 1 ? '增加' : action === 2 ? '减少' : '-'
}

watch(
    () => route.meta.currencyType,
    () => {
      searchForm.userId = ''
      pagination.pageIndex = 1
      fetchList()
    },
)

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
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

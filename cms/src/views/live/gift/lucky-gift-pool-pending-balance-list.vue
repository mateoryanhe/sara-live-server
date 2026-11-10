<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LuckyGiftPoolPendingBalanceList') }}</span>
          <el-button :loading="loading" @click="fetchList">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.luckyGiftPoolPendingBalance.pageDesc') }}</p>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('pages.luckyGiftPoolPendingBalance.status')">
          <el-select v-model="searchForm.status" clearable style="width: 160px">
            <el-option :label="t('pages.luckyGiftPoolPendingBalance.statusAll')" :value="-1"/>
            <el-option :label="t('pages.luckyGiftPoolPendingBalance.statusPending')" :value="0"/>
            <el-option :label="t('pages.luckyGiftPoolPendingBalance.statusProcessed')" :value="1"/>
            <el-option :label="t('pages.luckyGiftPoolPendingBalance.statusProcessing')" :value="2"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
        <el-table-column label="ID" min-width="180" prop="id" show-overflow-tooltip/>
        <el-table-column :label="t('pages.luckyGiftPoolPendingBalance.amount')" min-width="140" prop="amount"/>
        <el-table-column :label="t('pages.luckyGiftPoolPendingBalance.status')" width="120">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.luckyGiftPoolPendingBalance.createdAt')" width="170" prop="createdAt"/>
        <el-table-column :label="t('pages.luckyGiftPoolPendingBalance.updatedAt')" width="170" prop="updatedAt"/>
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
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {luckyGiftPoolPendingBalanceApi} from '@/api/modules/lucky-gift-pool-pending-balance'
import type {LuckyGiftPoolPendingBalanceItem} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<LuckyGiftPoolPendingBalanceItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  status: -1 as number,
})

const statusLabel = (status: number) => {
  if (status === 0) return t('pages.luckyGiftPoolPendingBalance.statusPending')
  if (status === 1) return t('pages.luckyGiftPoolPendingBalance.statusProcessed')
  if (status === 2) return t('pages.luckyGiftPoolPendingBalance.statusProcessing')
  return String(status)
}

const statusTagType = (status: number) => {
  if (status === 0) return 'warning'
  if (status === 1) return 'success'
  if (status === 2) return 'info'
  return ''
}

const formatRowIndex = (index: number) => (currentPage.value - 1) * pageSize.value + index + 1

const buildQueryParams = () => ({
  status: searchForm.status < 0 ? undefined : searchForm.status,
  pageIndex: currentPage.value,
  pageSize: pageSize.value,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await luckyGiftPoolPendingBalanceApi.listPendingBalance(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load lucky gift pool pending balance list:', error)
    ElMessage.error(t('pages.luckyGiftPoolPendingBalance.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleReset = () => {
  searchForm.status = -1
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

const handleCurrentChange = () => {
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
}

.page-desc {
  margin: 0 0 16px;
  color: #606266;
  font-size: 13px;
}

.search-form {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

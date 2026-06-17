<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>直播记录</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item label="主播ID">
          <el-input v-model="searchForm.anchorId" clearable placeholder="请输入主播ID"/>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              range-separator="至"
              start-placeholder="开始日期"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="记录ID" min-width="180" prop="id"/>
        <el-table-column label="主播ID" min-width="180" prop="anchorId"/>
        <el-table-column label="主播昵称" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="开始时间" width="170">
          <template #default="{ row }">{{ formatDate(row.startTime) }}</template>
        </el-table-column>
        <el-table-column label="结束时间" width="170">
          <template #default="{ row }">{{ formatDate(row.endTime) }}</template>
        </el-table-column>
        <el-table-column label="累计观众" prop="totalAudience" width="100"/>
        <el-table-column label="直播时长" width="120">
          <template #default="{ row }">{{ formatDuration(row.totalLiveDuration) }}</template>
        </el-table-column>
        <el-table-column label="总收益" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalIncome) }}</template>
        </el-table-column>
        <el-table-column label="礼物收入" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalGiftIncome) }}</template>
        </el-table-column>
        <el-table-column label="付费弹幕收入" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPaidDanmakuIncome) }}</template>
        </el-table-column>
        <el-table-column label="私密房收入" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalPrivateRoomIncome) }}</template>
        </el-table-column>
        <el-table-column label="送礼人数" prop="totalGiftSender" width="100"/>
        <el-table-column label="新加粉丝" prop="totalNewFollower" width="100"/>
        <el-table-column label="游戏下注总额" width="130">
          <template #default="{ row }">{{ formatAmount(row.totalGameBet) }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {liveRecordApi} from '@/api'
import type {LiveRecordItem} from '@/types/api'

const loading = ref(false)
const tableData = ref<LiveRecordItem[]>([])

const searchForm = reactive({
  anchorId: '',
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
    anchorId: searchForm.anchorId.trim(),
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRecordApi.getLiveRecordList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取直播记录失败:', error)
    ElMessage.error('获取直播记录失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.anchorId = ''
  searchForm.dateRange = []
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

const formatAmount = (value: number | null | undefined) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(2)
}

const formatDuration = (seconds: number | null | undefined) => {
  if (seconds === null || seconds === undefined || seconds <= 0) {
    return '-'
  }
  const total = Math.floor(seconds)
  const h = Math.floor(total / 3600)
  const m = Math.floor((total % 3600) / 60)
  const s = total % 60
  if (h > 0) {
    return `${h}时${m}分${s}秒`
  }
  if (m > 0) {
    return `${m}分${s}秒`
  }
  return `${s}秒`
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

.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

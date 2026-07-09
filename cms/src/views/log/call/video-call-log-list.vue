<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>视频通话日志</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item label="呼叫者ID">
          <el-input v-model="searchForm.callerId" clearable placeholder="请输入呼叫者ID"/>
        </el-form-item>
        <el-form-item label="接收者ID">
          <el-input v-model="searchForm.receiverId" clearable placeholder="请输入接收者ID"/>
        </el-form-item>
        <el-form-item label="来源">
          <el-select v-model="searchForm.source" clearable placeholder="全部" style="width: 140px">
            <el-option :value="0" label="全部"/>
            <el-option :value="1" label="直播间"/>
            <el-option :value="2" label="私信"/>
          </el-select>
        </el-form-item>
        <el-form-item label="呼叫时间">
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
        <el-table-column label="订单ID" min-width="180" prop="id"/>
        <el-table-column label="状态" prop="statusText" width="90"/>
        <el-table-column label="来源" prop="sourceText" width="90"/>
        <el-table-column label="呼叫者ID" min-width="180" prop="callerId"/>
        <el-table-column label="呼叫者昵称" min-width="120">
          <template #default="{ row }">{{ row.callerNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="接收者ID" min-width="180" prop="receiverId"/>
        <el-table-column label="接收者昵称" min-width="120">
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="呼叫时间" width="170">
          <template #default="{ row }">{{ formatDate(row.callStartTime) }}</template>
        </el-table-column>
        <el-table-column label="接听时间" width="170">
          <template #default="{ row }">{{ formatDate(row.answerTime) }}</template>
        </el-table-column>
        <el-table-column label="结束时间" width="170">
          <template #default="{ row }">{{ formatDate(row.orderEndTime) }}</template>
        </el-table-column>
        <el-table-column label="通话时长" width="120">
          <template #default="{ row }">{{ formatDuration(row.callDuration) }}</template>
        </el-table-column>
        <el-table-column label="门票(钻石)" width="110">
          <template #default="{ row }">{{ formatAmount(row.ticketPrice) }}</template>
        </el-table-column>
        <el-table-column label="分钟单价(钻石)" width="130">
          <template #default="{ row }">{{ formatAmount(row.pricePerMinute) }}</template>
        </el-table-column>
        <el-table-column label="计费时长(分钟)" prop="billingDuration" width="120"/>
        <el-table-column label="总费用(钻石)" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalCost) }}</template>
        </el-table-column>
        <el-table-column label="最近扣费时间" width="170">
          <template #default="{ row }">{{ formatDate(row.chargeTime) }}</template>
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
import {videoCallLogApi} from '@/api'
import type {VideoCallLogItem} from '@/types/api'

const loading = ref(false)
const tableData = ref<VideoCallLogItem[]>([])

const searchForm = reactive({
  callerId: '',
  receiverId: '',
  source: 0,
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
    callerId: searchForm.callerId.trim(),
    receiverId: searchForm.receiverId.trim(),
    source: searchForm.source || 0,
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await videoCallLogApi.getVideoCallLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取视频通话日志失败:', error)
    ElMessage.error('获取视频通话日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.callerId = ''
  searchForm.receiverId = ''
  searchForm.source = 0
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

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>直播收益流水</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item label="收益用户ID">
          <el-input v-model="searchForm.receiverId" clearable placeholder="请输入收益用户ID"/>
        </el-form-item>
        <el-form-item label="收益类型">
          <el-select v-model="searchForm.revenueType" clearable placeholder="全部" style="width: 140px">
            <el-option :value="0" label="全部"/>
            <el-option :value="1" label="礼物"/>
            <el-option :value="2" label="付费弹幕"/>
            <el-option :value="3" label="游戏下注"/>
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
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
        <el-table-column label="流水ID" min-width="180" prop="id"/>
        <el-table-column label="收益类型" prop="revenueTypeText" width="100">
          <template #default="{ row }">{{ row.revenueTypeText || formatRevenueType(row.revenueType) }}</template>
        </el-table-column>
        <el-table-column label="直播间ID" min-width="180" prop="roomId"/>
        <el-table-column label="直播记录ID" min-width="180" prop="liveRecordId"/>
        <el-table-column label="付款用户ID" min-width="180" prop="senderId"/>
        <el-table-column label="付款用户昵称" min-width="120" prop="senderNickname">
          <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="收益用户ID" min-width="180" prop="receiverId"/>
        <el-table-column label="收益用户昵称" min-width="120" prop="receiverNickname">
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="业务ID" min-width="180" prop="bizId"/>
        <el-table-column label="业务名称" min-width="120" prop="bizName">
          <template #default="{ row }">{{ row.bizName || '-' }}</template>
        </el-table-column>
        <el-table-column label="数量" prop="count" width="80"/>
        <el-table-column label="单价(钻石)" prop="unitPrice" width="110"/>
        <el-table-column label="流水金额(钻石)" prop="totalAmount" width="120"/>
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
import {liveRevenueLogApi} from '@/api'
import type {LiveRevenueLogItem} from '@/types/api'

const loading = ref(false)
const tableData = ref<LiveRevenueLogItem[]>([])

const searchForm = reactive({
  receiverId: '',
  revenueType: 0,
  dateRange: [] as string[],
})

const revenueTypeMap: Record<number, string> = {
  1: '礼物',
  2: '付费弹幕',
  3: '游戏下注',
}

const formatRevenueType = (type: number) => revenueTypeMap[type] || '未知'

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
    receiverId: searchForm.receiverId.trim(),
    revenueType: searchForm.revenueType || 0,
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRevenueLogApi.getLiveRevenueLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取直播收益流水失败:', error)
    ElMessage.error('获取直播收益流水失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.receiverId = ''
  searchForm.revenueType = 0
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

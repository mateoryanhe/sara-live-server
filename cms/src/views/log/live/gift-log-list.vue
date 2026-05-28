<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>礼物流水</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item label="接收者ID">
          <el-input v-model="searchForm.receiverId" clearable placeholder="请输入接收者ID"/>
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
        <el-table-column label="直播间ID" min-width="180" prop="roomId"/>
        <el-table-column label="直播记录ID" min-width="180" prop="liveRecordId"/>
        <el-table-column label="送礼用户ID" min-width="180" prop="senderId"/>
        <el-table-column label="送礼用户昵称" min-width="120" prop="senderNickname">
          <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="接收者ID" min-width="180" prop="receiverId"/>
        <el-table-column label="接收者昵称" min-width="120" prop="receiverNickname">
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="礼物ID" min-width="180" prop="giftId"/>
        <el-table-column label="礼物名称" min-width="120" prop="giftName">
          <template #default="{ row }">{{ row.giftName || '-' }}</template>
        </el-table-column>
        <el-table-column label="数量" prop="count" width="80"/>
        <el-table-column label="单价(钻石)" prop="unitPrice" width="110"/>
        <el-table-column label="总消耗(钻石)" prop="totalCost" width="120"/>
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
import {liveGiftLogApi} from '@/api'
import type {LiveGiftLogItem} from '@/types/api'

const loading = ref(false)
const tableData = ref<LiveGiftLogItem[]>([])

const searchForm = reactive({
  receiverId: '',
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
    receiverId: searchForm.receiverId.trim(),
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveGiftLogApi.getLiveGiftLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取礼物流水失败:', error)
    ElMessage.error('获取礼物流水失败')
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

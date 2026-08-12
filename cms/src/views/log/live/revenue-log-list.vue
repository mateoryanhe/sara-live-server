<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRevenueLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.revenueLogList.receiverId')">
          <el-input v-model="searchForm.receiverId" clearable :placeholder="t('pages.revenueLogList.enterReceiverId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.revenueLogList.revenueType')">
          <el-select v-model="searchForm.revenueType" clearable :placeholder="t('common.all')" style="width: 140px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.revenueLogList.revenueGift')"/>
            <el-option :value="2" :label="t('pages.revenueLogList.revenuePaidDanmaku')"/>
            <el-option :value="3" :label="t('pages.revenueLogList.revenueGameBet')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.createdAt')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.revenueLogList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.revenueLogList.dateRangeSeparator')"
              :start-placeholder="t('pages.revenueLogList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.revenueLogList.logId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.revenueLogList.revenueType')" prop="revenueTypeText" width="100">
          <template #default="{ row }">{{ row.revenueTypeText || formatRevenueType(row.revenueType) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.liveRoomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.revenueLogList.liveRecordId')" min-width="180" prop="liveRecordId"/>
        <el-table-column :label="t('pages.revenueLogList.payerUserId')" min-width="180" prop="senderId"/>
        <el-table-column :label="t('pages.revenueLogList.payerNickname')" min-width="120" prop="senderNickname">
          <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.receiverUserId')" min-width="180" prop="receiverId"/>
        <el-table-column :label="t('pages.revenueLogList.receiverNickname')" min-width="120" prop="receiverNickname">
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.bizId')" min-width="180" prop="bizId"/>
        <el-table-column :label="t('pages.revenueLogList.bizName')" min-width="120" prop="bizName">
          <template #default="{ row }">{{ row.bizName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.count')" prop="count" width="80"/>
        <el-table-column :label="t('pages.revenueLogList.unitPriceDiamond')" prop="unitPrice" width="110"/>
        <el-table-column :label="t('pages.revenueLogList.totalAmountDiamond')" prop="totalAmount" width="120"/>
        <el-table-column :label="t('common.status')" prop="statusText" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'warning' : 'success'" size="small">
              {{ row.statusText || (row.status === 1 ? t('pages.revenueLogList.refunded') : t('common.normal')) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" width="170">
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
import {useI18n} from 'vue-i18n'
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {liveRevenueLogApi} from '@/api'
import type {LiveRevenueLogItem} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<LiveRevenueLogItem[]>([])

const searchForm = reactive({
  receiverId: '',
  revenueType: 0,
  dateRange: [] as string[],
})

const formatRevenueType = (type: number) => {
  const map: Record<number, string> = {
    1: t('pages.revenueLogList.revenueGift'),
    2: t('pages.revenueLogList.revenuePaidDanmaku'),
    3: t('pages.revenueLogList.revenueGameBet'),
  }
  return map[type] || t('pages.revenueLogList.unknown')
}

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
    console.error('Failed to load revenue logs:', error)
    ElMessage.error(t('pages.revenueLogList.fetchFailed'))
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
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
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

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.VideoCallLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.videoCallLogList.callerId')">
          <el-input v-model="searchForm.callerId" clearable :placeholder="t('pages.videoCallLogList.enterCallerId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.videoCallLogList.receiverId')">
          <el-input v-model="searchForm.receiverId" clearable :placeholder="t('pages.videoCallLogList.enterReceiverId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.videoCallLogList.source')">
          <el-select v-model="searchForm.source" clearable :placeholder="t('common.all')" style="width: 140px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.videoCallLogList.sourceLiveRoom')"/>
            <el-option :value="2" :label="t('pages.videoCallLogList.sourcePrivateMessage')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="searchForm.status" clearable :placeholder="t('common.all')" style="width: 160px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.videoCallLogList.statusCalling')"/>
            <el-option :value="2" :label="t('pages.videoCallLogList.statusAnswered')"/>
            <el-option :value="3" :label="t('pages.videoCallLogList.statusInCall')"/>
            <el-option :value="4" :label="t('pages.videoCallLogList.statusEnded')"/>
            <el-option :value="5" :label="t('pages.videoCallLogList.statusRejected')"/>
            <el-option :value="6" :label="t('pages.videoCallLogList.statusCallTimeout')"/>
            <el-option :value="7" :label="t('pages.videoCallLogList.statusHeartTimeout')"/>
            <el-option :value="8" :label="t('pages.videoCallLogList.statusInsufficientDiamond')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.videoCallLogList.callTime')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.videoCallLogList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.videoCallLogList.dateRangeSeparator')"
              :start-placeholder="t('pages.videoCallLogList.startDate')"
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
        <el-table-column :label="t('pages.rechargeOrderList.orderId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.status')" prop="statusText" width="140"/>
        <el-table-column :label="t('pages.videoCallLogList.source')" prop="sourceText" width="90"/>
        <el-table-column :label="t('pages.videoCallLogList.callerId')" min-width="180" prop="callerId"/>
        <el-table-column :label="t('pages.videoCallLogList.callerNickname')" min-width="120">
          <template #default="{ row }">{{ row.callerNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.receiverId')" min-width="180" prop="receiverId"/>
        <el-table-column :label="t('pages.videoCallLogList.receiverNickname')" min-width="120">
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.callStartTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.answerTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.answerTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callerLastHeart')" width="170">
          <template #default="{ row }">{{ formatDate(row.callerHeartTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.receiverLastHeart')" width="170">
          <template #default="{ row }">{{ formatDate(row.receiverHeartTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.endTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.orderEndTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callDuration')" width="120">
          <template #default="{ row }">{{ formatDuration(row.callDuration) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.ticketDiamond')" width="110">
          <template #default="{ row }">{{ formatAmount(row.ticketPrice) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.pricePerMinuteDiamond')" width="130">
          <template #default="{ row }">{{ formatAmount(row.pricePerMinute) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.billingDurationMinutes')" prop="billingDuration" width="120"/>
        <el-table-column :label="t('pages.videoCallLogList.totalCostDiamond')" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalCost) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.lastChargeTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.chargeTime) }}</template>
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
import {videoCallLogApi} from '@/api'
import type {VideoCallLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<VideoCallLogItem[]>([])

const searchForm = reactive({
  callerId: '',
  receiverId: '',
  source: 0,
  status: 0,
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
    status: searchForm.status || 0,
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
    console.error('Failed to load video call logs:', error)
    ElMessage.error(t('pages.videoCallLogList.fetchFailed'))
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
  searchForm.status = 0
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

const formatDuration = (seconds: number | null | undefined) => {
  if (seconds === null || seconds === undefined || seconds <= 0) return '-'
  const total = Math.floor(seconds)
  const h = Math.floor(total / 3600)
  const m = Math.floor((total % 3600) / 60)
  const s = total % 60
  if (h > 0) {
    return t('pages.videoCallLogList.durationHours', {h, m, s})
  }
  if (m > 0) {
    return t('pages.videoCallLogList.durationMinutes', {m, s})
  }
  return t('pages.videoCallLogList.durationSeconds', {s})
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

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRecordList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.liveRecordList.anchorFilter')">
          <AnchorRemoteSelect
              v-model="searchForm.anchorIds"
              :placeholder="t('pages.liveRecordList.searchAnchor')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.startDate')">
          <el-date-picker
              v-model="searchForm.startDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.startDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.endDate')">
          <el-date-picker
              v-model="searchForm.endDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.endDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('export')" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
        <el-table-column :label="t('pages.liveRecordList.recordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.liveRecordList.anchorId')" min-width="180" prop="anchorId"/>
        <el-table-column :label="t('pages.liveRecordList.anchorNickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.startTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.startTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.endTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.endTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.totalAudience')" prop="totalAudience" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.liveDuration')" width="120">
          <template #default="{ row }">{{ formatDuration(row.totalLiveDuration) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.totalIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.giftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.paidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoTicketIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoBillingIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.giftSenderCount')" prop="totalGiftSender" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.newFollowers')" prop="totalNewFollower" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.totalGameBet')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameBet) }}</span></template>
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
import {liveRecordApi} from '@/api'
import type {LiveRecordItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_RECORD} from '@/utils/cms-async-export'
import {buildLiveRecordCsvColumns} from '@/utils/live-record-csv'
import {formatWalletBalance} from '@/utils/number-format'

const {t} = useI18n()
const {can} = usePagePermission('LiveRecordList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<LiveRecordItem[]>([])

const searchForm = reactive({
  anchorIds: [] as string[],
  startDate: '',
  endDate: '',
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

const buildQueryParams = () => ({
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
  anchorIds: [...searchForm.anchorIds],
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
})

const buildFilterParams = () => ({
  anchorIds: [...searchForm.anchorIds],
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRecordApi.getLiveRecordList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load live records:', error)
    ElMessage.error(t('pages.liveRecordList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.anchorIds = []
  searchForm.startDate = ''
  searchForm.endDate = ''
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

const handleExport = async () => {
  await runExport(
    CMS_EXPORT_TYPE_LIVE_RECORD,
    {
      headers: buildCsvHeaders(buildLiveRecordCsvColumns(t)),
      ...buildFilterParams(),
    },
    `live-record-${Date.now()}.csv`,
  )
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
    return t('pages.liveRecordList.durationHours', {h, m, s})
  }
  if (m > 0) {
    return t('pages.liveRecordList.durationMinutes', {m, s})
  }
  return t('pages.liveRecordList.durationSeconds', {s})
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

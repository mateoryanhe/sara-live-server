<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.guildList.anchorSettlementLogHint')"
        type="info"
    />

    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')">
        <el-input
            v-model="searchForm.roomId"
            clearable
            :placeholder="t('pages.guildAnchorIncomeSettlementLogList.enterRoomId')"
            style="width: 200px"
        />
      </el-form-item>
      <el-form-item :label="t('common.createdAt')">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('pages.guildAnchorIncomeSettlementLogList.endDate')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.guildAnchorIncomeSettlementLogList.dateRangeSeparator')"
            :start-placeholder="t('pages.guildAnchorIncomeSettlementLogList.startDate')"
            style="width: 260px"
            type="daterange"
            value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        <el-button :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.logId')" fixed="left" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120">
        <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalIncome')" min-width="110" prop="totalIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalGiftIncome')" min-width="110" prop="totalGiftIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalPaidDanmakuIncome')" min-width="120" prop="totalPaidDanmakuIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalPrivateRoomTicketIncome')" min-width="130" prop="totalPrivateRoomTicketIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalPrivateRoomWatchIncome')" min-width="130" prop="totalPrivateRoomWatchIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallIncome')" min-width="120" prop="totalVideoCallIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallTicketIncome')" min-width="140" prop="totalVideoCallTicketIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallBillingIncome')" min-width="140" prop="totalVideoCallBillingIncome"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementSalary')" min-width="110" prop="settlementSalary"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.anchorSharePercent')" min-width="110">
        <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementShareAmount')" min-width="120" prop="settlementShareAmount"/>
      <el-table-column :label="t('common.createdAt')" fixed="right" width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildList.noAnchorSettlementLogData')"/>

    <div v-if="pagination.total > 0" class="pagination">
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
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {downloadCsv, fetchAllPagedRows} from '@/utils/csv-export'
import {buildGuildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const exporting = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const loaded = ref(false)

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const formatDateString = (date: Date) => {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

const createDefaultDateRange = () => {
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - 6)
  return [formatDateString(start), formatDateString(end)] as string[]
}

const searchForm = reactive({
  roomId: '',
  dateRange: createDefaultDateRange(),
})

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    guildId: props.guildId,
    roomId: searchForm.roomId.trim(),
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await guildApi.getGuildAnchorIncomeSettlementLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load guild anchor settlement logs:', error)
    ElMessage.error(t('pages.guildList.anchorSettlementLogFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.roomId = ''
  searchForm.dateRange = createDefaultDateRange()
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
  if (!props.guildId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  exporting.value = true
  try {
    const rows = await fetchAllPagedRows((pageIndex, pageSize) =>
      guildApi.getGuildAnchorIncomeSettlementLogList({
        ...buildFilterParams(),
        pageIndex,
        pageSize,
      }),
    )
    if (rows.length === 0) {
      ElMessage.warning(t('common.exportEmpty'))
      return
    }
    downloadCsv(
      `guild-anchor-settlement-log-${props.guildId}-${Date.now()}.csv`,
      buildGuildAnchorSettlementLogCsvColumns(t),
      rows,
    )
    ElMessage.success(t('common.exportSuccess'))
  } catch (error) {
    console.error('Failed to export guild anchor settlement logs:', error)
    ElMessage.error(t('common.exportFailed'))
  } finally {
    exporting.value = false
  }
}

const formatSharePercent = (value: number | null | undefined) => {
  if (value == null || Number.isNaN(value)) return '-'
  return `${value}%`
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  pagination.pageIndex = 1
  pagination.total = 0
  searchForm.roomId = ''
  searchForm.dateRange = createDefaultDateRange()
}

watch(
  () => props.guildId,
  () => {
    resetState()
    if (props.active) {
      fetchList()
    }
  },
)

watch(
  () => props.active,
  (active) => {
    if (active && !loaded.value && props.guildId) {
      fetchList()
    }
  },
  {immediate: true},
)
</script>

<style scoped>
.hint-alert {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.search-form :deep(.el-form-item__label) {
  white-space: nowrap;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

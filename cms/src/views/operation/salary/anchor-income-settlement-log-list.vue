<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AnchorIncomeSettlementLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.anchorIncomeSettlementLogList.roomId')">
          <el-input v-model="searchForm.roomId" clearable :placeholder="t('pages.anchorIncomeSettlementLogList.enterRoomId')"/>
        </el-form-item>
        <el-form-item :label="t('common.createdAt')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.anchorIncomeSettlementLogList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.anchorIncomeSettlementLogList.dateRangeSeparator')"
              :start-placeholder="t('pages.anchorIncomeSettlementLogList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('export')" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.logId')" fixed="left" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomNickname')" min-width="120" prop="roomNickname">
          <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalIncome')" min-width="110" prop="totalIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalGiftIncome')" min-width="110" prop="totalGiftIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPaidDanmakuIncome')" min-width="120" prop="totalPaidDanmakuIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomTicketIncome')" min-width="130" prop="totalPrivateRoomTicketIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomWatchIncome')" min-width="130" prop="totalPrivateRoomWatchIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallIncome')" min-width="120" prop="totalVideoCallIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallTicketIncome')" min-width="140" prop="totalVideoCallTicketIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallBillingIncome')" min-width="140" prop="totalVideoCallBillingIncome"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalLiveDuration')" min-width="120" prop="totalLiveDuration"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementSalary')" min-width="110" prop="settlementSalary"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSharePercent')" min-width="110" prop="anchorSharePercent">
          <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementShareAmount')" min-width="120" prop="settlementShareAmount"/>
        <el-table-column :label="t('common.createdAt')" fixed="right" width="170">
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
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {downloadCsv, fetchAllPagedRows} from '@/utils/csv-export'
import {buildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'

const {t} = useI18n()
const {can} = usePagePermission('AnchorIncomeSettlementLogList')
const loading = ref(false)
const exporting = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])

const searchForm = reactive({
  roomId: '',
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

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
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
  loading.value = true
  try {
    const response = await anchorIncomeSettlementLogApi.getList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch anchor income settlement log list failed:', error)
    ElMessage.error(t('pages.anchorIncomeSettlementLogList.fetchFailed'))
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

const handleExport = async () => {
  exporting.value = true
  try {
    const rows = await fetchAllPagedRows((pageIndex, pageSize) =>
      anchorIncomeSettlementLogApi.getList({
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
      `anchor-income-settlement-log-${Date.now()}.csv`,
      buildAnchorSettlementLogCsvColumns(t),
      rows,
    )
    ElMessage.success(t('common.exportSuccess'))
  } catch (error) {
    console.error('export anchor income settlement log failed:', error)
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
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

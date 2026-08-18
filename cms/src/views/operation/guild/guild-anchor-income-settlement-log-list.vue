<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildProfileAnchorSettlementLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.guildAnchorIncomeSettlementLogList.guildId')">
          <el-select
              v-model="searchForm.guildId"
              clearable
              filterable
              :placeholder="t('pages.guildAnchorIncomeSettlementLogList.selectGuild')"
              style="width: 220px"
          >
            <el-option :label="t('pages.guildAnchorIncomeSettlementLogList.allGuilds')" value=""/>
            <el-option
                v-for="item in guildOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')">
          <el-input v-model="searchForm.roomId" clearable :placeholder="t('pages.guildAnchorIncomeSettlementLogList.enterRoomId')"/>
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
          <el-button v-if="can('export')" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.logId')" fixed="left" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.guildName')" min-width="120" prop="guildName">
          <template #default="{ row }">{{ row.guildName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120" prop="roomNickname">
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
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalLiveDuration')" min-width="120" prop="totalLiveDuration"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementSalary')" min-width="110" prop="settlementSalary"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.anchorSharePercent')" min-width="110" prop="anchorSharePercent">
          <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementShareAmount')" min-width="120" prop="settlementShareAmount"/>
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
import {onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {AnchorIncomeSettlementLogItem, MyGuildProfile} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {downloadCsv, fetchAllPagedRows} from '@/utils/csv-export'
import {buildGuildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'

const {t} = useI18n()
const route = useRoute()
const {can} = usePagePermission('GuildProfileAnchorSettlementLogList')
const loading = ref(false)
const exporting = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const guildOptions = ref<MyGuildProfile[]>([])

const searchForm = reactive({
  guildId: '',
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
    guildId: searchForm.guildId.trim(),
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

const fetchGuildOptions = async () => {
  try {
    const response = await guildApi.getMyGuildProfile()
    guildOptions.value = response?.list ?? []
  } catch (error) {
    console.error('fetch my guild profile failed:', error)
    guildOptions.value = []
    ElMessage.error(t('pages.guildAnchorIncomeSettlementLogList.fetchGuildFailed'))
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getMyGuildAnchorIncomeSettlementLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch my guild anchor income settlement log list failed:', error)
    ElMessage.error(t('pages.guildAnchorIncomeSettlementLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.guildId = ''
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
      guildApi.getMyGuildAnchorIncomeSettlementLogList({
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
      `guild-anchor-income-settlement-log-${Date.now()}.csv`,
      buildGuildAnchorSettlementLogCsvColumns(t),
      rows,
    )
    ElMessage.success(t('common.exportSuccess'))
  } catch (error) {
    console.error('export guild anchor income settlement log failed:', error)
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

const applyRouteGuildFilter = () => {
  const value = route.query.guildId
  const guildId = Array.isArray(value) ? String(value[0] ?? '') : String(value ?? '')
  searchForm.guildId = guildId
}

watch(
  () => route.query.guildId,
  () => {
    applyRouteGuildFilter()
    pagination.pageIndex = 1
    fetchList()
  },
)

onMounted(async () => {
  applyRouteGuildFilter()
  await fetchGuildOptions()
  await fetchList()
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

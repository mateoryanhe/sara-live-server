<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AnchorIncomeSettlementLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('pages.liveRecordList.platformAnchor')">
          <div class="anchor-filter anchor-filter--compact">
            <el-input
                :model-value="platformAnchorInputValue"
                class="anchor-input"
                disabled
                :placeholder="t('pages.liveRecordList.noPlatformAnchorSelected')"
            />
            <el-button size="small" @click="openPlatformAnchorPicker">
              {{ t('pages.liveRecordList.selectPlatformAnchor') }}
            </el-button>
            <el-button v-if="selectedPlatformAnchors.length > 0" link size="small" @click="clearPlatformAnchors">
              {{ t('pages.liveRecordList.clearPlatformAnchor') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.guildAnchor')">
          <div class="anchor-filter anchor-filter--compact">
            <el-input
                :model-value="guildAnchorInputValue"
                class="anchor-input"
                disabled
                :placeholder="t('pages.liveRecordList.noGuildAnchorSelected')"
            />
            <el-button size="small" @click="openGuildAnchorPicker">
              {{ t('pages.liveRecordList.selectGuildAnchor') }}
            </el-button>
            <el-button v-if="selectedGuildAnchors.length > 0" link size="small" @click="clearGuildAnchors">
              {{ t('pages.liveRecordList.clearGuildAnchor') }}
            </el-button>
          </div>
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
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.logId')" fixed="left" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomNickname')" min-width="120" prop="roomNickname">
          <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalGiftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPaidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomTicketIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomWatchIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomWatchIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallTicketIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallBillingIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalLiveDuration')" min-width="120">
          <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementSalary')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementSalary) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSharePercent')" min-width="110" prop="anchorSharePercent">
          <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementShareAmount')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementShareAmount) }}</span></template>
        </el-table-column>
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

    <AnchorPickerDialog
        v-model:visible="platformAnchorPickerVisible"
        :initial-anchors="selectedPlatformAnchors"
        multiple
        platform-only
        @confirm-multiple="handlePlatformAnchorsPicked"
    />

    <GuildAnchorPickerDialog
        v-model:visible="guildAnchorPickerVisible"
        :initial-anchors="selectedGuildAnchors"
        @confirm="handleGuildAnchorsPicked"
    />
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem, AnchorListItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import GuildAnchorPickerDialog from '@/components/GuildAnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_ANCHOR_INCOME_SETTLEMENT_LOG} from '@/utils/cms-async-export'
import {buildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const {t} = useI18n()
const {can} = usePagePermission('AnchorIncomeSettlementLogList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const selectedGuildAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)
const guildAnchorPickerVisible = ref(false)

const searchForm = reactive({
  startDate: '',
  endDate: '',
})

const formatPlatformAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  return `${nickname} (${anchor.id})`
}

const formatGuildAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  if (anchor.guildName) {
    return `[${anchor.guildName}] ${nickname} (${anchor.id})`
  }
  return `${nickname} (${anchor.id})`
}

const platformAnchorInputValue = computed(() => {
  if (selectedPlatformAnchors.value.length === 0) {
    return ''
  }
  if (selectedPlatformAnchors.value.length === 1) {
    return formatPlatformAnchorLabel(selectedPlatformAnchors.value[0])
  }
  return t('pages.liveRecordList.selectedPlatformAnchorsCount', {count: selectedPlatformAnchors.value.length})
})

const guildAnchorInputValue = computed(() => {
  if (selectedGuildAnchors.value.length === 0) {
    return ''
  }
  if (selectedGuildAnchors.value.length === 1) {
    return formatGuildAnchorLabel(selectedGuildAnchors.value[0])
  }
  return t('pages.liveRecordList.selectedGuildAnchorsCount', {count: selectedGuildAnchors.value.length})
})

const openPlatformAnchorPicker = () => {
  platformAnchorPickerVisible.value = true
}

const clearPlatformAnchors = () => {
  selectedPlatformAnchors.value = []
}

const handlePlatformAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedPlatformAnchors.value = anchors
}

const openGuildAnchorPicker = () => {
  guildAnchorPickerVisible.value = true
}

const clearGuildAnchors = () => {
  selectedGuildAnchors.value = []
}

const handleGuildAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedGuildAnchors.value = anchors
}

const buildSelectedAnchorIds = () => {
  const platformIds = selectedPlatformAnchors.value.map(anchor => String(anchor.id))
  const guildIds = selectedGuildAnchors.value.map(anchor => String(anchor.id))
  return [...new Set([...platformIds, ...guildIds])]
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

const buildFilterParams = () => ({
  anchorIds: buildSelectedAnchorIds(),
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
})

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
  selectedPlatformAnchors.value = []
  selectedGuildAnchors.value = []
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
    CMS_EXPORT_TYPE_ANCHOR_INCOME_SETTLEMENT_LOG,
    {
      headers: buildCsvHeaders(buildAnchorSettlementLogCsvColumns(t)),
      ...buildFilterParams(),
    },
    `anchor-income-settlement-log-${Date.now()}.csv`,
  )
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

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 12px;
}

.anchor-filter--compact {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.anchor-input {
  width: 220px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

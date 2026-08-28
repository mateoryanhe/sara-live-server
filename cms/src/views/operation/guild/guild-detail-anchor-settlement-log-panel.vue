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

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width:100%">
      <el-table-column :label="t('common.createdAt')" fixed="left" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.avatar')" fixed="left" width="80">
        <template #default="{ row }">
          <el-image
              v-if="row.roomAvatar"
              :preview-src-list="[row.roomAvatar]"
              :src="row.roomAvatar"
              fit="cover"
              hide-on-click-modal
              preview-teleported
              style="width:40px;height:40px;border-radius:50%"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180">
        <template #default="{ row }">
          <el-button v-if="row.roomId" link type="primary" @click="openAnchorDetail(row.roomId)">
            {{ row.roomId }}
          </el-button>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120">
        <template #default="{ row }">
          <el-button
              v-if="canViewUserDetail && row.roomId && row.roomNickname"
              link
              type="primary"
              @click="openUserDetail(row.roomId)"
          >
            {{ row.roomNickname }}
          </el-button>
          <span v-else>{{ row.roomNickname || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementSalary')" align="right" min-width="140">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementSalary) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.settlementShareAmount')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementShareAmount) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalGiftIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalPaidDanmakuIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallTicketIncome')" align="right" min-width="150">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalVideoCallBillingIncome')" align="right" min-width="150">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalShortVideoIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalShortVideoIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalGameIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.anchorSharePercent')" min-width="120">
        <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
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
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_GUILD_ANCHOR_INCOME_SETTLEMENT_LOG} from '@/utils/cms-async-export'
import {buildGuildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const router = useRouter()
const {canViewUserDetail, openUserDetail} = useUserDetailNav('GuildDetail')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const loaded = ref(false)

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const createDefaultDateRange = () => {
  const end = new Date()
  const start = new Date(end.getTime() - 6 * 86400000)
  return [formatServerDateOnly(start), formatServerDateOnly(end)] as string[]
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
  await runExport(
    CMS_EXPORT_TYPE_GUILD_ANCHOR_INCOME_SETTLEMENT_LOG,
    {
      headers: buildCsvHeaders(buildGuildAnchorSettlementLogCsvColumns(t)),
      ...buildFilterParams(),
    },
    `guild-anchor-settlement-log-${props.guildId}-${Date.now()}.csv`,
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

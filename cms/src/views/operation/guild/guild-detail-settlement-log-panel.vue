<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.guildList.settlementLogHint')"
        type="info"
    />

    <el-form :model="searchForm" class="search-form" inline label-width="100px">
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
        <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width:100%">
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column
          :label="t('pages.anchorList.settlementSalary')"
          align="right"
          min-width="140"
      >
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementSalary) }}</span></template>
      </el-table-column>
      <el-table-column
          :label="t('pages.anchorList.settlementFlowCommission')"
          align="right"
          min-width="120"
      >
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementShareAmount) }}</span></template>
      </el-table-column>
      <el-table-column
          :label="t('pages.anchorList.settlementShareAmountUsd')"
          align="right"
          min-width="130"
      >
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementShareAmountUsd) }}</span></template>
      </el-table-column>
      <el-table-column
          :label="t('pages.anchorList.settlementReceivableUsd')"
          align="right"
          min-width="280"
      >
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalGiftIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalPaidDanmakuIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallTicketIncome')" align="right" min-width="150">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallBillingIncome')" align="right" min-width="150">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalShortVideoIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalShortVideoIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalGameIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.guildSharePercent')" min-width="150">
        <template #default="{ row }">{{ formatSharePercent(row.guildSharePercent) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildList.noSettlementLogData')"/>
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_GUILD_INCOME_SETTLEMENT_LOG} from '@/utils/cms-async-export'
import {buildGuildSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const {can} = usePagePermission('GuildDetail')
const canExport = computed(() => can('exportSettlementLog'))
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<GuildIncomeSettlementLogItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    guildId: props.guildId,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await guildIncomeSettlementLogApi.getList({
      ...buildFilterParams(),
      pageIndex: 1,
      pageSize: 50,
    })
    tableData.value = response.data || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load guild settlement logs:', error)
    ElMessage.error(t('pages.guildList.settlementLogFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchList()
}

const handleReset = () => {
  searchForm.dateRange = []
  fetchList()
}

const handleExport = async () => {
  if (!props.guildId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  await runExport(
    CMS_EXPORT_TYPE_GUILD_INCOME_SETTLEMENT_LOG,
    {
      headers: buildCsvHeaders(buildGuildSettlementLogCsvColumns(t)),
      ...buildFilterParams(),
    },
    `guild-settlement-log-${props.guildId}-${Date.now()}.csv`,
  )
}

const formatSharePercent = (value: number | null | undefined) => {
  if (value == null || Number.isNaN(value)) return '0%'
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
  searchForm.dateRange = []
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
</style>

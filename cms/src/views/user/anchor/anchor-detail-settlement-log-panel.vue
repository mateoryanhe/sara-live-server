<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.anchorList.settlementLogHint')"
        type="info"
    />

    <div class="toolbar">
      <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
    </div>

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.logId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalIncome')" min-width="110" prop="totalIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalGiftIncome')" min-width="110" prop="totalGiftIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPaidDanmakuIncome')" min-width="120" prop="totalPaidDanmakuIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomTicketIncome')" min-width="130" prop="totalPrivateRoomTicketIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPrivateRoomWatchIncome')" min-width="130" prop="totalPrivateRoomWatchIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallIncome')" min-width="120" prop="totalVideoCallIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallTicketIncome')" min-width="140" prop="totalVideoCallTicketIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallBillingIncome')" min-width="140" prop="totalVideoCallBillingIncome"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementSalary')" min-width="110" prop="settlementSalary"/>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSharePercent')" min-width="110" prop="anchorSharePercent">
        <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementShareAmount')" min-width="120" prop="settlementShareAmount"/>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.anchorList.noSettlementLogData')"/>
  </div>
</template>

<script lang="ts" setup>
import {computed, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {downloadCsv, fetchAllPagedRows} from '@/utils/csv-export'
import {buildAnchorSettlementLogCsvColumns} from '@/utils/income-settlement-log-csv'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  anchorId: string
  active: boolean
}>()

const {t} = useI18n()
const {can} = usePagePermission('AnchorDetail')
const canExport = computed(() => can('exportSettlementLog'))
const loading = ref(false)
const exporting = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const loaded = ref(false)

const buildFilterParams = () => ({
  roomId: props.anchorId,
})

const fetchList = async () => {
  if (!props.anchorId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await anchorIncomeSettlementLogApi.getList({
      ...buildFilterParams(),
      pageIndex: 1,
      pageSize: 50,
    })
    tableData.value = response.data || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load anchor settlement logs:', error)
    ElMessage.error(t('pages.anchorList.settlementLogFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  if (!props.anchorId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
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
      `anchor-settlement-log-${props.anchorId}-${Date.now()}.csv`,
      buildAnchorSettlementLogCsvColumns(t),
      rows,
    )
    ElMessage.success(t('common.exportSuccess'))
  } catch (error) {
    console.error('Failed to export anchor settlement logs:', error)
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
}

watch(
  () => props.anchorId,
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
    if (active && !loaded.value && props.anchorId) {
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

.toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

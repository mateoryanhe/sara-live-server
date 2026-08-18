<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.guildList.settlementLogHint')"
        type="info"
    />

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.logId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalIncome')" min-width="110">
        <template #default="{ row }">{{ formatNum(row.totalIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalGiftIncome')" min-width="110">
        <template #default="{ row }">{{ formatNum(row.totalGiftIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalPaidDanmakuIncome')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalPaidDanmakuIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalPrivateRoomTicketIncome')" min-width="130">
        <template #default="{ row }">{{ formatNum(row.totalPrivateRoomTicketIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalPrivateRoomWatchIncome')" min-width="130">
        <template #default="{ row }">{{ formatNum(row.totalPrivateRoomWatchIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallIncome')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalVideoCallIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallTicketIncome')" min-width="140">
        <template #default="{ row }">{{ formatNum(row.totalVideoCallTicketIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalVideoCallBillingIncome')" min-width="140">
        <template #default="{ row }">{{ formatNum(row.totalVideoCallBillingIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalLiveDuration) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.guildSharePercent')" min-width="110">
        <template #default="{ row }">{{ formatSharePercent(row.guildSharePercent) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.guildIncomeSettlementLogList.settlementShareAmount')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.settlementShareAmount) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildList.noSettlementLogData')"/>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<GuildIncomeSettlementLogItem[]>([])
const loaded = ref(false)

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await guildIncomeSettlementLogApi.getList({
      guildId: props.guildId,
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

const formatNum = (value?: number | null) => formatAmount(value, '0')

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
</style>

<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.guildList.anchorDailyEffectiveLiveHint')"
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
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
      <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120">
        <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.liveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.dailyReportedLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.liveIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.giftIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.privateRoomTicketIncome')" align="right" min-width="140">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomTicketIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.privateRoomWatchIncome')" align="right" min-width="140">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomWatchIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.videoCallIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.dailySettled')" min-width="100">
        <template #default="{ row }">
          <el-tag :type="row.settled ? 'success' : 'warning'">
            {{ row.settled ? t('pages.anchorList.dailySettledYes') : t('pages.anchorList.dailySettledNo') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.roomUpdatedAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildList.noAnchorDailyEffectiveLiveData')"/>
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {GuildAnchorDailyEffectiveLiveItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {downloadCsv, fetchAllPagedRows} from '@/utils/csv-export'
import {buildGuildAnchorDailyEffectiveLiveCsvColumns} from '@/utils/daily-effective-live-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const {can} = usePagePermission('GuildDetail')
const canExport = computed(() => can('exportAnchorDailyEffectiveLive'))
const loading = ref(false)
const exporting = ref(false)
const tableData = ref<GuildAnchorDailyEffectiveLiveItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  roomId: '',
})

const DEFAULT_PAGE_SIZE = 8

const buildFilterParams = () => ({
  guildId: props.guildId,
  roomId: searchForm.roomId.trim() || undefined,
  settled: 0,
})

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await guildApi.getGuildAnchorDailyEffectiveLiveList({
      ...buildFilterParams(),
      pageIndex: 1,
      pageSize: DEFAULT_PAGE_SIZE,
    })
    tableData.value = response.data || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load guild anchor daily flow list:', error)
    ElMessage.error(t('pages.guildList.anchorDailyEffectiveLiveFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchList()
}

const handleReset = () => {
  searchForm.roomId = ''
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
      guildApi.getGuildAnchorDailyEffectiveLiveList({
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
      `guild-anchor-daily-flow-${props.guildId}-${Date.now()}.csv`,
      buildGuildAnchorDailyEffectiveLiveCsvColumns(t),
      rows,
    )
    ElMessage.success(t('common.exportSuccess'))
  } catch (error) {
    console.error('Failed to export guild anchor daily flow:', error)
    ElMessage.error(t('common.exportFailed'))
  } finally {
    exporting.value = false
  }
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
  searchForm.roomId = ''
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

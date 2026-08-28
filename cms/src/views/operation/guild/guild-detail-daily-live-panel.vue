<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.guildList.dailyEffectiveLiveHint')"
        type="info"
    />

    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item :label="t('pages.anchorList.dailyLiveDate')">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('pages.anchorList.liveDateEnd')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.anchorList.liveDateRangeSeparator')"
            :start-placeholder="t('pages.anchorList.liveDateStart')"
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
      <el-table-column :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
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

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildList.noDailyEffectiveLiveData')"/>
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {GuildDailyEffectiveLiveItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, buildDailyEffectiveLiveExportLabels, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_GUILD_DAILY_EFFECTIVE_LIVE} from '@/utils/cms-async-export'
import {buildGuildDailyEffectiveLiveCsvColumns} from '@/utils/daily-effective-live-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const {can} = usePagePermission('GuildDetail')
const canExport = computed(() => can('exportDailyEffectiveLive'))
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<GuildDailyEffectiveLiveItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

const DEFAULT_PAGE_SIZE = 8

const buildFilterParams = () => {
  const [liveDateStart, liveDateEnd] = searchForm.dateRange || []
  return {
    guildId: props.guildId,
    liveDateStart: liveDateStart || undefined,
    liveDateEnd: liveDateEnd || undefined,
    settled: 0,
  }
}

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await guildApi.getGuildDailyEffectiveLiveList({
      ...buildFilterParams(),
      pageIndex: 1,
      pageSize: DEFAULT_PAGE_SIZE,
    })
    tableData.value = response.data || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load guild daily flow list:', error)
    ElMessage.error(t('pages.guildList.dailyEffectiveLiveFetchFailed'))
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
    CMS_EXPORT_TYPE_GUILD_DAILY_EFFECTIVE_LIVE,
    {
      headers: buildCsvHeaders(buildGuildDailyEffectiveLiveCsvColumns(t)),
      guildId: Number(props.guildId),
      ...buildFilterParams(),
      ...buildDailyEffectiveLiveExportLabels(t),
    },
    `guild-daily-flow-${props.guildId}-${Date.now()}.csv`,
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

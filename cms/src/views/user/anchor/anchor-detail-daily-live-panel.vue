<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.anchorList.dailyEffectiveLiveHint')"
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
      </el-form-item>
    </el-form>

    <div class="toolbar">
      <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
    </div>

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width:100%">
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.liveDuration, t) }}</template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.dailyReportedLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.liveIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.giftIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.paidDanmakuIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.privateRoomTicketIncome')" align="right" min-width="140">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomTicketIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.privateRoomWatchIncome')" align="right" min-width="140">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomWatchIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.videoCallIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.dailySettled')" min-width="100">
        <template #default="{ row }">
          <el-tag :type="row.settled ? 'success' : 'warning'">
            {{ row.settled ? t('pages.anchorList.dailySettledYes') : t('pages.anchorList.dailySettledNo') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column v-if="!simpleColumns" :label="t('pages.anchorList.roomUpdatedAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.anchorList.noDailyEffectiveLiveData')"/>

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
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {accountApi, guildApi} from '@/api'
import type {AnchorDailyEffectiveLiveItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, buildDailyEffectiveLiveExportLabels, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_ANCHOR_DAILY_EFFECTIVE_LIVE} from '@/utils/cms-async-export'
import {buildAnchorDailyEffectiveLiveCsvColumns} from '@/utils/daily-effective-live-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  anchorId: string
  active?: boolean
  guildId?: string
  permissionPage?: string
  simpleColumns?: boolean
}>()

const {t} = useI18n()
const resolvedPermissionPage = props.permissionPage ?? (props.guildId ? 'GuildProfileManagement' : 'AnchorDetail')
const {can} = usePagePermission(resolvedPermissionPage)
const canExport = computed(() => can('exportDailyEffectiveLive'))
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<AnchorDailyEffectiveLiveItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildFilterParams = () => {
  const [liveDateStart, liveDateEnd] = searchForm.dateRange || []
  const base = {
    liveDateStart: liveDateStart || undefined,
    liveDateEnd: liveDateEnd || undefined,
    settled: 0,
  }
  if (props.guildId) {
    return {
      ...base,
      guildId: props.guildId,
      anchorId: props.anchorId,
    }
  }
  return {
    ...base,
    anchorId: props.anchorId,
  }
}

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  if (!props.anchorId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = props.guildId
        ? await guildApi.getMyGuildAnchorDailyEffectiveLiveList(buildQueryParams())
        : await accountApi.getAnchorDailyEffectiveLiveList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load daily effective live list:', error)
    ElMessage.error(t('pages.anchorList.dailyEffectiveLiveFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
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
  if (!props.anchorId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  await runExport(
    CMS_EXPORT_TYPE_ANCHOR_DAILY_EFFECTIVE_LIVE,
    {
      headers: buildCsvHeaders(buildAnchorDailyEffectiveLiveCsvColumns(t)),
      anchorId: Number(props.anchorId),
      ...buildFilterParams(),
      ...buildDailyEffectiveLiveExportLabels(t),
    },
    `anchor-daily-flow-${props.anchorId}-${Date.now()}.csv`,
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
  pagination.pageIndex = 1
  pagination.total = 0
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
    if (active !== false && !loaded.value && props.anchorId) {
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

.toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

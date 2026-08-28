<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="100px">
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
        <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
      <el-table-column :label="t('pages.liveRecordList.recordId')" min-width="180" prop="id">
        <template #default="{ row }">
          <el-button v-if="row.id" link type="primary" @click="openLiveRecordDetail(row.id)">
            {{ row.id }}
          </el-button>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.startTime')" width="170">
        <template #default="{ row }">{{ formatDate(row.startTime) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.endTime')" width="170">
        <template #default="{ row }">{{ formatDate(row.endTime) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.totalAudience')" prop="totalAudience" width="100"/>
      <el-table-column :label="t('pages.liveRecordList.liveDuration')" width="120">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.totalIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.giftIncome')" align="right" min-width="120">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.paidDanmakuIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.videoTicketIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.videoBillingIncome')" align="right" min-width="150">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.videoCallIncome')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('pages.liveRecordList.giftSenderCount')" prop="totalGiftSender" width="100"/>
      <el-table-column :label="t('pages.liveRecordList.newFollowers')" prop="totalNewFollower" width="100"/>
      <el-table-column :label="t('pages.liveRecordList.totalGameBet')" align="right" min-width="130">
        <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameBet) }}</span></template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" width="170">
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
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {liveRecordApi} from '@/api'
import type {LiveRecordItem} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_RECORD} from '@/utils/cms-async-export'
import {buildLiveRecordCsvColumns} from '@/utils/live-record-csv'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  anchorId: string
  active: boolean
}>()

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('AnchorDetail')
const canExport = computed(() => can('exportLiveRecord'))
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<LiveRecordItem[]>([])
const loaded = ref(false)

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const createDefaultDateRange = () => {
  const end = new Date()
  const start = new Date(end.getTime() - 6 * 86400000)
  return {
    startDate: formatServerDateOnly(start),
    endDate: formatServerDateOnly(end),
  }
}

const searchForm = reactive(createDefaultDateRange())

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildFilterParams = () => ({
  anchorId: props.anchorId,
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
})

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
    const response = await liveRecordApi.getLiveRecordList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load anchor live records:', error)
    ElMessage.error(t('pages.liveRecordList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  Object.assign(searchForm, createDefaultDateRange())
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

const openLiveRecordDetail = (recordId: string | number) => {
  if (!recordId || !props.anchorId) {
    return
  }
  router.push({
    name: 'AnchorLiveRecordDetail',
    query: {
      anchorId: props.anchorId,
      liveRecordId: String(recordId),
    },
  })
}

const handleExport = async () => {
  if (!props.anchorId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  await runExport(
    CMS_EXPORT_TYPE_LIVE_RECORD,
    {
      headers: buildCsvHeaders(buildLiveRecordCsvColumns(t)),
      ...buildFilterParams(),
    },
    `anchor-live-record-${props.anchorId}-${Date.now()}.csv`,
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
  pagination.pageIndex = 1
  pagination.total = 0
  Object.assign(searchForm, createDefaultDateRange())
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
.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

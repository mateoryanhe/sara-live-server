<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.anchorList.backToLiveRecordList') }}</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <el-empty v-if="!loading && !liveRecord" :description="t('pages.liveRecordList.fetchFailed')"/>
        <template v-else-if="liveRecord">
          <el-descriptions :column="2" border class="record-summary">
            <el-descriptions-item :label="t('pages.liveRecordList.recordId')">{{ liveRecord.id }}</el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.anchorId')">{{ liveRecord.anchorId }}</el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.anchorNickname')">{{ liveRecord.nickname || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.totalAudience')">{{ liveRecord.totalAudience }}</el-descriptions-item>
            <el-descriptions-item :label="t('common.startTime')">{{ formatDate(liveRecord.startTime) }}</el-descriptions-item>
            <el-descriptions-item :label="t('common.endTime')">{{ formatDate(liveRecord.endTime) }}</el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.liveDuration')">
              {{ formatLiveDurationMinutes(liveRecord.totalLiveDuration, t) }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.totalIncome')">
              <span class="money-amount">{{ formatWalletBalance(liveRecord.totalIncome) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.giftIncome')">
              <span class="money-amount">{{ formatWalletBalance(liveRecord.totalGiftIncome) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.paidDanmakuIncome')">
              <span class="money-amount">{{ formatWalletBalance(liveRecord.totalPaidDanmakuIncome) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('pages.liveRecordList.videoCallIncome')">
              <span class="money-amount">{{ formatWalletBalance(liveRecord.totalVideoCallIncome) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('common.createdAt')">{{ formatDate(liveRecord.createdAt) }}</el-descriptions-item>
          </el-descriptions>

          <el-divider v-if="canViewRevenue" content-position="left">{{ t('menu.LiveRevenueLogList') }}</el-divider>
          <template v-if="canViewRevenue">
          <el-alert
              :closable="false"
              class="hint-alert"
              show-icon
              :title="t('pages.anchorList.liveRecordRevenueHint')"
              type="info"
          />

          <el-form :model="searchForm" class="search-form" inline label-width="88px">
            <el-form-item :label="t('pages.revenueLogList.revenueType')">
              <el-select v-model="searchForm.revenueType" clearable :placeholder="t('common.all')" style="width: 200px">
                <el-option :value="0" :label="t('common.all')"/>
                <el-option
                    v-for="option in revenueTypeOptions"
                    :key="option.value"
                    :label="t(option.labelKey)"
                    :value="option.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
              <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
              <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
            </el-form-item>
          </el-form>

          <el-table v-loading="revenueLoading || exporting" :data="revenueTableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
            <el-table-column :label="t('pages.revenueLogList.logId')" min-width="180" prop="id"/>
            <el-table-column :label="t('pages.revenueLogList.revenueType')" min-width="140" prop="revenueTypeText">
              <template #default="{ row }">{{ row.revenueTypeText || formatRevenueType(row.revenueType) }}</template>
            </el-table-column>
            <el-table-column :label="t('pages.revenueLogList.payerUserId')" min-width="180" prop="senderId"/>
            <el-table-column :label="t('pages.revenueLogList.payerNickname')" min-width="120" prop="senderNickname">
              <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
            </el-table-column>
            <el-table-column :label="t('pages.revenueLogList.bizId')" min-width="180" prop="bizId"/>
            <el-table-column :label="t('pages.revenueLogList.bizName')" min-width="120" prop="bizName">
              <template #default="{ row }">{{ row.bizName || '-' }}</template>
            </el-table-column>
            <el-table-column :label="t('pages.revenueLogList.count')" prop="count" width="80"/>
            <el-table-column :label="t('pages.revenueLogList.unitPriceDiamond')" prop="unitPrice" width="110"/>
            <el-table-column :label="t('pages.revenueLogList.totalAmountDiamond')" prop="totalAmount" width="120"/>
            <el-table-column :label="t('common.status')" prop="statusText" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'warning' : 'success'" size="small">
                  {{ row.statusText || (row.status === 1 ? t('pages.revenueLogList.refunded') : t('common.normal')) }}
                </el-tag>
              </template>
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
          </template>
        </template>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {liveRecordApi, liveRevenueLogApi} from '@/api'
import type {LiveRecordItem, LiveRevenueLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_REVENUE_LOG} from '@/utils/cms-async-export'
import {buildLiveRevenueLogCsvColumns} from '@/utils/live-revenue-log-csv'
import {createLiveRevenueTypeFormatter, LIVE_REVENUE_TYPE_OPTIONS} from '@/utils/live-revenue-type'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {can} = usePagePermission('AnchorDetail')
const canViewRevenue = computed(() => can('liveRecordRevenue'))
const canExport = computed(() => can('exportLiveRecord'))
const revenueTypeOptions = LIVE_REVENUE_TYPE_OPTIONS
const formatRevenueType = createLiveRevenueTypeFormatter(t)
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()

const loading = ref(false)
const revenueLoading = ref(false)
const liveRecord = ref<LiveRecordItem | null>(null)
const revenueTableData = ref<LiveRevenueLogItem[]>([])

const searchForm = reactive({
  revenueType: 0,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const anchorId = computed(() => {
  const value = route.query.anchorId
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  return value == null ? '' : String(value)
})

const liveRecordId = computed(() => {
  const value = route.query.liveRecordId
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  return value == null ? '' : String(value)
})

const pageTitle = computed(() => {
  if (liveRecordId.value) {
    return t('pages.anchorList.liveRecordDetailTitleWithId', {id: liveRecordId.value})
  }
  return t('pages.anchorList.liveRecordDetailTitle')
})

const buildRevenueQueryParams = () => ({
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
  liveRecordId: liveRecordId.value,
  receiverId: anchorId.value || undefined,
  revenueType: searchForm.revenueType || 0,
})

const buildRevenueFilterParams = () => ({
  liveRecordId: liveRecordId.value,
  receiverId: anchorId.value || undefined,
  revenueType: searchForm.revenueType || 0,
})

const fetchLiveRecord = async () => {
  if (!liveRecordId.value) {
    liveRecord.value = null
    return
  }
  loading.value = true
  try {
    const response = await liveRecordApi.getLiveRecordList({
      pageIndex: 1,
      pageSize: 1,
      anchorId: anchorId.value || undefined,
      liveRecordId: liveRecordId.value,
    })
    liveRecord.value = response.data?.[0] ?? null
  } catch (error) {
    console.error('Failed to load live record detail:', error)
    liveRecord.value = null
    ElMessage.error(t('pages.liveRecordList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchRevenueList = async () => {
  if (!canViewRevenue.value) {
    revenueTableData.value = []
    pagination.total = 0
    return
  }
  if (!liveRecordId.value) {
    revenueTableData.value = []
    pagination.total = 0
    return
  }
  revenueLoading.value = true
  try {
    const response = await liveRevenueLogApi.getLiveRevenueLogList(buildRevenueQueryParams())
    revenueTableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load live record revenue logs:', error)
    revenueTableData.value = []
    pagination.total = 0
    ElMessage.error(t('pages.anchorList.liveRecordRevenueFetchFailed'))
  } finally {
    revenueLoading.value = false
  }
}

const fetchAll = async () => {
  pagination.pageIndex = 1
  await Promise.all([fetchLiveRecord(), fetchRevenueList()])
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchRevenueList()
}

const handleReset = () => {
  searchForm.revenueType = 0
  pagination.pageIndex = 1
  fetchRevenueList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchRevenueList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchRevenueList()
}

const handleExport = async () => {
  if (!liveRecordId.value) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  await runExport(
    CMS_EXPORT_TYPE_LIVE_REVENUE_LOG,
    {
      headers: buildCsvHeaders(buildLiveRevenueLogCsvColumns(t, formatRevenueType)),
      ...buildRevenueFilterParams(),
    },
    `live-record-revenue-${liveRecordId.value}-${Date.now()}.csv`,
  )
}

const goBack = () => {
  if (anchorId.value) {
    router.push({
      name: 'AnchorDetail',
      query: {
        id: anchorId.value,
        tab: 'liveRecord',
      },
    })
    return
  }
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'AnchorListManagement'})
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

watch([anchorId, liveRecordId], () => {
  fetchAll()
})

onMounted(() => {
  fetchAll()
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
  font-size: 16px;
  font-weight: bold;
}

.record-summary {
  margin-bottom: 16px;
}

.hint-alert {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

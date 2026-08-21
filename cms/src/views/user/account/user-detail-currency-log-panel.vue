<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item :label="t('common.createdAt')">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('common.endTime')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.anchorIncomeSettlementLogList.dateRangeSeparator')"
            :start-placeholder="t('common.startTime')"
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

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column :label="t('pages.currencyLogList.logId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.currencyLogList.changeType')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.action === 1 ? 'success' : 'danger'">
            {{ actionLabel(row.action) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="amountColumnLabel" prop="amount" width="120">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.currencyLogList.beforeChange')" prop="before" width="120">
        <template #default="{ row }">{{ formatAmount(row.before) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.currencyLogList.afterChange')" prop="after" width="120">
        <template #default="{ row }">{{ formatAmount(row.after) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.currencyLogList.reason')" min-width="140" prop="reasonText" show-overflow-tooltip>
        <template #default="{ row }">{{ row.reasonText || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.currencyLogList.time')" width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.userList.noData')"/>

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
import {currencyLogApi} from '@/api'
import type {CurrencyLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_CURRENCY_LOG} from '@/utils/cms-async-export'
import {buildCurrencyLogCsvColumns} from '@/utils/currency-log-csv'

const props = defineProps<{
  userId: string
  active: boolean
  currencyType: 1 | 2
  exportPermission: 'exportGoldLog' | 'exportDiamondLog'
}>()

const {t} = useI18n()
const {can} = usePagePermission('UserDetail')
const canExport = computed(() => can(props.exportPermission))
const amountColumnLabel = computed(() =>
    props.currencyType === 2
        ? t('pages.currencyLogList.diamondChange')
        : t('pages.currencyLogList.goldChange'))
const fetchFailedMessage = computed(() =>
    props.currencyType === 2
        ? t('pages.userList.diamondLogFetchFailed')
        : t('pages.userList.goldLogFetchFailed'))
const exportFilePrefix = computed(() =>
    props.currencyType === 2 ? 'user-diamond-log' : 'user-gold-log')

const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<CurrencyLogItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

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

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    userId: props.userId,
    currencyType: props.currencyType,
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
  if (!props.userId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await currencyLogApi.getCurrencyLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load user currency logs:', error)
    ElMessage.error(fetchFailedMessage.value)
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

const formatRowIndex = (index: number) =>
    (pagination.pageIndex - 1) * pagination.pageSize + index + 1

const handleExport = async () => {
  if (!props.userId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  await runExport(
      CMS_EXPORT_TYPE_CURRENCY_LOG,
      {
        headers: buildCsvHeaders(buildCurrencyLogCsvColumns(t, props.currencyType)),
        ...buildFilterParams(),
      },
      `${exportFilePrefix.value}-${props.userId}-${Date.now()}.csv`,
  )
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const actionLabel = (action: number) => {
  if (action === 1) return t('pages.currencyLogList.actionIncrease')
  if (action === 2) return t('pages.currencyLogList.actionDecrease')
  return '-'
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  searchForm.dateRange = []
  pagination.pageIndex = 1
  pagination.total = 0
}

watch(
    () => [props.userId, props.currencyType] as const,
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
      if (active && !loaded.value && props.userId) {
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
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

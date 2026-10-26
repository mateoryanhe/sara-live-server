<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item v-if="showCurrencyTypeFilter" :label="t('pages.currencyLogList.currencyType')">
        <el-select v-model="searchForm.currencyType" style="width: 140px">
          <el-option
              v-if="canViewGold && canViewDiamond"
              :label="t('pages.currencyLogList.currencyTypeAll')"
              :value="0"
          />
          <el-option
              v-if="canViewGold"
              :label="t('pages.currencyLogList.currencyTypeGold')"
              :value="1"
          />
          <el-option
              v-if="canViewDiamond"
              :label="t('pages.currencyLogList.currencyTypeDiamond')"
              :value="2"
          />
        </el-select>
      </el-form-item>
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
      <el-table-column :label="t('pages.currencyLogList.currencyType')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.currencyType === 2 ? 'warning' : 'primary'">
            {{ currencyTypeLabel(row.currencyType) }}
          </el-tag>
        </template>
      </el-table-column>
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
import {formatServerDateTime as formatDate, toServerDayStartUnix, toServerDayEndUnix} from '@/utils/server-datetime'

const props = defineProps<{
  userId: string
  active: boolean
}>()

const {t, locale} = useI18n()
const {can} = usePagePermission('UserDetail')
const canViewGold = computed(() => can('goldLog'))
const canViewDiamond = computed(() => can('diamondLog'))
const canExportGold = computed(() => can('exportGoldLog'))
const canExportDiamond = computed(() => can('exportDiamondLog'))
const showCurrencyTypeFilter = computed(() => canViewGold.value && canViewDiamond.value)

const defaultCurrencyType = (): 0 | 1 | 2 => {
  if (canViewGold.value && canViewDiamond.value) return 0
  if (canViewDiamond.value) return 2
  return 1
}

const resolveQueryCurrencyType = (selected: 0 | 1 | 2): 0 | 1 | 2 => {
  if (selected === 1 && canViewGold.value) return 1
  if (selected === 2 && canViewDiamond.value) return 2
  if (selected === 0 && canViewGold.value && canViewDiamond.value) return 0
  return defaultCurrencyType()
}

const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<CurrencyLogItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  currencyType: defaultCurrencyType() as 0 | 1 | 2,
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const canExport = computed(() => {
  const type = resolveQueryCurrencyType(searchForm.currencyType)
  if (type === 1) return canExportGold.value
  if (type === 2) return canExportDiamond.value
  return canExportGold.value || canExportDiamond.value
})

const amountColumnLabel = computed(() => {
  const type = resolveQueryCurrencyType(searchForm.currencyType)
  if (type === 2) return t('pages.currencyLogList.diamondChange')
  if (type === 1) return t('pages.currencyLogList.goldChange')
  return t('pages.currencyLogList.amountChange')
})

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    userId: props.userId,
    currencyType: resolveQueryCurrencyType(searchForm.currencyType),
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
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
    ElMessage.error(t('pages.userList.currencyLogFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.currencyType = defaultCurrencyType()
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
  const filters = buildFilterParams()
  let exportCurrencyType = filters.currencyType as 0 | 1 | 2
  // 仅有一侧导出权限时，合并视图也只导出有权限的类型
  if (exportCurrencyType === 0) {
    if (canExportGold.value && !canExportDiamond.value) exportCurrencyType = 1
    else if (!canExportGold.value && canExportDiamond.value) exportCurrencyType = 2
    else if (!canExportGold.value && !canExportDiamond.value) return
  } else if (exportCurrencyType === 1 && !canExportGold.value) {
    return
  } else if (exportCurrencyType === 2 && !canExportDiamond.value) {
    return
  }
  await runExport(
      CMS_EXPORT_TYPE_CURRENCY_LOG,
      {
        headers: buildCsvHeaders(buildCurrencyLogCsvColumns(t, exportCurrencyType)),
        lang: String(locale.value || 'en'),
        ...filters,
        currencyType: exportCurrencyType,
      },
      `user-currency-log-${props.userId}-${Date.now()}.csv`,
  )
}

const currencyTypeLabel = (currencyType: number) => {
  if (currencyType === 2) return t('pages.currencyLogList.currencyTypeDiamond')
  if (currencyType === 1) return t('pages.currencyLogList.currencyTypeGold')
  return '-'
}

const actionLabel = (action: number) => {
  if (action === 1) return t('pages.currencyLogList.actionIncrease')
  if (action === 2) return t('pages.currencyLogList.actionDecrease')
  return '-'
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  searchForm.currencyType = defaultCurrencyType()
  searchForm.dateRange = []
  pagination.pageIndex = 1
  pagination.total = 0
}

watch(
    () => props.userId,
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

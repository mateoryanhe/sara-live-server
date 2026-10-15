<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item :label="t('pages.rechargeOrderList.orderId')">
        <el-input v-model="searchForm.orderId" clearable :placeholder="t('pages.rechargeOrderList.enterOrderId')" style="width: 200px"/>
      </el-form-item>
      <el-form-item :label="t('common.createdAt')">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('common.endTime')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.rechargeOrderList.dateRangeSeparator')"
            :start-placeholder="t('common.startTime')"
            style="width: 260px"
            type="daterange"
            value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item :label="t('pages.rechargeOrderList.orderStatus')">
        <el-select v-model="searchForm.statusFilter" clearable style="width: 140px">
          <el-option :value="0" :label="t('common.all')"/>
          <el-option :value="1" :label="t('pages.rechargeOrderList.statusPending')"/>
          <el-option :value="2" :label="t('pages.rechargeOrderList.statusCompleted')"/>
          <el-option :value="3" :label="t('pages.rechargeOrderList.statusCancelled')"/>
          <el-option :value="4" :label="t('pages.rechargeOrderList.statusFailed')"/>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <div class="list-summary">{{ t('pages.rechargeOrderList.totalOrders', {count: pagination.total}) }}</div>

    <el-table v-loading="loading" :data="tableData" style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column :label="t('pages.rechargeOrderList.orderId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.rechargeOrderList.userCountry')" min-width="120" prop="userCountry" show-overflow-tooltip>
        <template #default="{ row }">{{ row.userCountry || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.cfgId')" min-width="100" prop="cfgId">
        <template #default="{ row }">{{ row.cfgId || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.price')" prop="price" width="110">
        <template #default="{ row }">{{ formatAmount(row.price) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.payAmount')" prop="payAmount" width="120">
        <template #default="{ row }">{{ row.payAmount ? formatAmount(row.payAmount) : '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.payCurrency')" prop="currency" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.currency" effect="plain" size="small">{{ row.currency }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.gold')" prop="gold" width="110">
        <template #default="{ row }">{{ formatAmount(row.gold) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.source')" width="100">
        <template #default="{ row }">{{ sourceLabel(row.source) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.thirdOrderId')" min-width="140" prop="thirdOrderId" show-overflow-tooltip>
        <template #default="{ row }">{{ row.thirdOrderId || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.remark')" min-width="120" prop="remark" show-overflow-tooltip>
        <template #default="{ row }">{{ row.remark || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" width="170">
        <template #default="{ row }">{{ formatUnixTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.rechargeOrderList.paidAt')" width="170">
        <template #default="{ row }">{{ formatUnixTime(row.paidAt) }}</template>
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
import {reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {rechargeOrderApi} from '@/api'
import type {RechargeOrder} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {formatServerDateTime, toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

const props = defineProps<{
  userId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<RechargeOrder[]>([])
const loaded = ref(false)

const searchForm = reactive({
  orderId: '',
  dateRange: [] as string[],
  statusFilter: 0,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
    userId: props.userId,
    orderId: searchForm.orderId.trim(),
    statusFilter: searchForm.statusFilter,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  if (!props.userId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await rechargeOrderApi.getRechargeOrderList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load user recharge orders:', error)
    ElMessage.error(t('pages.userList.rechargeOrderFetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.orderId = ''
  searchForm.dateRange = []
  searchForm.statusFilter = 0
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

const formatUnixTime = (ts: number | null | undefined) => {
  if (!ts) return '-'
  return formatServerDateTime(ts)
}

const statusLabel = (status: number) => {
  switch (status) {
    case 1:
      return t('pages.rechargeOrderList.statusCompleted')
    case 2:
      return t('pages.rechargeOrderList.statusCancelled')
    case 3:
      return t('pages.rechargeOrderList.statusFailed')
    default:
      return t('pages.rechargeOrderList.statusPending')
  }
}

const statusTagType = (status: number) => {
  switch (status) {
    case 1:
      return 'success'
    case 2:
      return 'info'
    case 3:
      return 'danger'
    default:
      return 'warning'
  }
}

const sourceLabel = (source: number) => {
  switch (source) {
    case 1:
      return t('pages.rechargeOrderList.sourceAndroid')
    case 2:
      return t('pages.rechargeOrderList.sourceManual')
    case 3:
      return t('pages.rechargeOrderList.sourceH5')
    case 5:
      return t('pages.rechargeOrderList.sourceIOS')
    case 6:
      return t('pages.rechargeOrderList.sourceCoinMerchant')
    default:
      return String(source || '-')
  }
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  searchForm.orderId = ''
  searchForm.dateRange = []
  searchForm.statusFilter = 0
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

.list-summary {
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

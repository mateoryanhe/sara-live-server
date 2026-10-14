<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionChannelDailyStatManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.promotionChannelDailyStatList.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.promotionChannelDailyStatList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.promotionChannelDailyStatList.dateRangeSeparator')"
              :start-placeholder="t('pages.promotionChannelDailyStatList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.promotionChannelDailyStatList.channelCode')">
          <el-input
              v-model="searchForm.channelCode"
              clearable
              :placeholder="t('pages.promotionChannelDailyStatList.channelCodePlaceholder')"
              style="width: 160px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.promotionChannelDailyStatList.statDate')" prop="statDate" width="120"/>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.channelName')" min-width="140" prop="channelName">
          <template #default="{ row }">{{ row.channelName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.channelCode')" prop="channelCode" width="110"/>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.registerCount')" prop="registerCount" width="100"/>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.orderCount')" prop="orderCount" width="100"/>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.payUserCount')" prop="payUserCount" width="100"/>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.orderTotalAmount')" min-width="130" prop="orderTotalAmount">
          <template #default="{ row }">{{ formatUsd(row.orderTotalAmount) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.promotionChannelDailyStatList.rechargeSuccessAmount')" min-width="140" prop="rechargeSuccessAmount">
          <template #default="{ row }">{{ formatUsd(row.rechargeSuccessAmount) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {promotionChannelDailyStatApi} from '@/api/modules/promotion-channel-daily-stat'
import type {PromotionChannelDailyStatRow} from '@/types/api'
import {formatServerDateOnly} from '@/utils/server-datetime'

interface SearchForm {
  dateRange: string[]
  channelCode: string
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<PromotionChannelDailyStatRow[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const createDefaultDateRange = () => {
  const end = new Date()
  const start = new Date(end.getTime() - 6 * 86400000)
  return [formatServerDateOnly(start), formatServerDateOnly(end)] as string[]
}

const searchForm = reactive<SearchForm>({
  dateRange: createDefaultDateRange(),
  channelCode: '',
})

const formatUsd = (val: number | undefined) => {
  const n = Number(val ?? 0)
  return n.toFixed(2)
}

const parseChannelCode = () => {
  const raw = searchForm.channelCode.trim()
  if (!raw) {
    return 0
  }
  const n = Number.parseInt(raw, 10)
  return Number.isFinite(n) && n > 0 ? n : -1
}

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
    startDate: startDate || '',
    endDate: endDate || '',
    channelCode: parseChannelCode(),
  }
}

const fetchList = async () => {
  const [startDate, endDate] = searchForm.dateRange || []
  if (!startDate || !endDate) {
    ElMessage.warning(t('pages.promotionChannelDailyStatList.dateRangeRequired'))
    return
  }
  if (parseChannelCode() < 0) {
    ElMessage.warning(t('pages.promotionChannelDailyStatList.channelCodeInvalid'))
    return
  }
  loading.value = true
  try {
    const params = buildQueryParams()
    const response = await promotionChannelDailyStatApi.getList(params)
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch promotion daily stat failed:', error)
    ElMessage.error(t('pages.promotionChannelDailyStatList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  void fetchList()
}

const handleReset = () => {
  searchForm.dateRange = createDefaultDateRange()
  searchForm.channelCode = ''
  currentPage.value = 1
  void fetchList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  void fetchList()
}

const handleCurrentChange = () => {
  void fetchList()
}

onMounted(() => {
  void fetchList()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

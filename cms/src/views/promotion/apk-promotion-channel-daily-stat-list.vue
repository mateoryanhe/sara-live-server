<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ApkPromotionChannelDailyStatManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.apkPromotionChannelDailyStatList.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.apkPromotionChannelDailyStatList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.apkPromotionChannelDailyStatList.dateRangeSeparator')"
              :start-placeholder="t('pages.apkPromotionChannelDailyStatList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.apkPromotionChannelDailyStatList.channelCode')">
          <div class="channel-code-field">
            <el-input
                v-model="searchForm.channelCode"
                clearable
                :placeholder="t('pages.apkPromotionChannelDailyStatList.channelCodePlaceholder')"
                style="width: 200px"
            />
            <el-button @click="openChannelPicker">{{ t('pages.apkPromotionChannelDailyStatList.selectChannel') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.statDate')" prop="statDate" width="120"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.channelName')" min-width="140" prop="channelName">
          <template #default="{ row }">{{ row.channelName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.channelCode')" prop="channelCode" width="140"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.registerCount')" prop="registerCount" width="100"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.downloadClickCount')" prop="downloadClickCount" width="110"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.orderCount')" prop="orderCount" width="100"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.payUserCount')" prop="payUserCount" width="100"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.orderTotalAmount')" min-width="130" prop="orderTotalAmount">
          <template #default="{ row }">{{ formatUsd(row.orderTotalAmount) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.rechargeSuccessAmount')" min-width="140" prop="rechargeSuccessAmount">
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

    <el-dialog
        v-model="channelPickerVisible"
        :title="t('pages.apkPromotionChannelDailyStatList.channelPickerTitle')"
        width="640px"
        destroy-on-close
        @open="handleChannelPickerOpen"
    >
      <div class="channel-picker-toolbar">
        <el-input
            v-model="channelPickerSearchKey"
            clearable
            :placeholder="t('pages.apkPromotionChannelDailyStatList.channelPickerSearchPlaceholder')"
            style="flex: 1"
            @keyup.enter="fetchChannelPickerList"
        />
        <el-button type="primary" @click="fetchChannelPickerList">{{ t('common.search') }}</el-button>
        <el-button @click="selectChannelCode('')">{{ t('pages.apkPromotionChannelDailyStatList.channelPickerAll') }}</el-button>
      </div>
      <el-table
          v-loading="channelPickerLoading"
          :data="channelPickerRows"
          highlight-current-row
          max-height="420"
          style="width: 100%; margin-top: 12px"
          @row-dblclick="(row: ApkPromotionChannel) => selectChannelCode(row.channelCode)"
      >
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.channelName')" min-width="160" prop="channelName"/>
        <el-table-column :label="t('pages.apkPromotionChannelDailyStatList.channelCode')" prop="channelCode" width="140"/>
        <el-table-column :label="t('common.actions')" width="100">
          <template #default="{ row }">
            <el-button link type="primary" @click="selectChannelCode(row.channelCode)">
              {{ t('pages.apkPromotionChannelDailyStatList.channelPickerSelect') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="channel-picker-pagination">
        <el-pagination
            v-model:current-page="channelPickerPage"
            v-model:page-size="channelPickerPageSize"
            :page-sizes="[10, 20, 50]"
            :total="channelPickerTotal"
            layout="total, prev, pager, next"
            @current-change="fetchChannelPickerList"
            @size-change="handleChannelPickerSizeChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {apkPromotionChannelApi, apkPromotionChannelDailyStatApi} from '@/api/modules/apk-promotion-channel'
import type {ApkPromotionChannel, ApkPromotionChannelDailyStatRow} from '@/types/api'
import {formatServerDateOnly} from '@/utils/server-datetime'

interface SearchForm {
  dateRange: string[]
  channelCode: string
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<ApkPromotionChannelDailyStatRow[]>([])
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

const channelPickerVisible = ref(false)
const channelPickerLoading = ref(false)
const channelPickerRows = ref<ApkPromotionChannel[]>([])
const channelPickerTotal = ref(0)
const channelPickerPage = ref(1)
const channelPickerPageSize = ref(20)
const channelPickerSearchKey = ref('')

const openChannelPicker = () => {
  channelPickerVisible.value = true
}

const handleChannelPickerOpen = () => {
  channelPickerPage.value = 1
  channelPickerSearchKey.value = ''
  void fetchChannelPickerList()
}

const fetchChannelPickerList = async () => {
  channelPickerLoading.value = true
  try {
    const response = await apkPromotionChannelApi.getList({
      pageIndex: channelPickerPage.value,
      pageSize: channelPickerPageSize.value,
      key: channelPickerSearchKey.value.trim(),
    })
    channelPickerRows.value = response.data
    channelPickerTotal.value = response.total
  } catch (error) {
    console.error('fetch channel picker list failed:', error)
    ElMessage.error(t('pages.apkPromotionChannelDailyStatList.fetchChannelListFailed'))
  } finally {
    channelPickerLoading.value = false
  }
}

const handleChannelPickerSizeChange = () => {
  channelPickerPage.value = 1
  void fetchChannelPickerList()
}

const selectChannelCode = (code: string) => {
  searchForm.channelCode = code
  channelPickerVisible.value = false
}

const formatUsd = (val: number | undefined) => Number(val ?? 0).toFixed(2)

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
    startDate: startDate || '',
    endDate: endDate || '',
    channelCode: searchForm.channelCode.trim(),
  }
}

const fetchList = async () => {
  const [startDate, endDate] = searchForm.dateRange || []
  if (!startDate || !endDate) {
    ElMessage.warning(t('pages.apkPromotionChannelDailyStatList.dateRangeRequired'))
    return
  }
  loading.value = true
  try {
    const response = await apkPromotionChannelDailyStatApi.getList(buildQueryParams())
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch apk promotion daily stat failed:', error)
    ElMessage.error(t('pages.apkPromotionChannelDailyStatList.fetchFailed'))
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
.page-container {
  padding: 20px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.channel-code-field {
  display: flex;
  align-items: center;
  gap: 8px;
}

.channel-picker-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.channel-picker-pagination {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
</style>

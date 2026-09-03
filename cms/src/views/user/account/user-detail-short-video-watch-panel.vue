<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="100px">
      <el-form-item :label="dateFilterLabel">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('pages.shortVideoWatchList.endDate')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.shortVideoWatchList.dateRangeSeparator')"
            :start-placeholder="t('pages.shortVideoWatchList.startDate')"
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

    <el-table v-loading="loading" :data="tableData" style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column :label="t('pages.shortVideoWatchList.recordId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.shortVideoWatchList.videoId')" min-width="180" prop="videoId"/>
      <el-table-column :label="t('pages.shortVideoWatchList.videoTitle')" min-width="160" prop="videoTitle" show-overflow-tooltip>
        <template #default="{ row }">{{ row.videoTitle || '-' }}</template>
      </el-table-column>
      <el-table-column v-if="onlyPaid" :label="t('pages.shortVideoList.payDiamond')" width="120" align="right">
        <template #default="{ row }">{{ formatWalletBalance(row.payDiamond) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.shortVideoWatchList.paidTime')" prop="paidTime" width="170">
        <template #default="{ row }">{{ row.paidTime || '-' }}</template>
      </el-table-column>
      <el-table-column v-if="!onlyPaid" :label="t('common.createdAt')" prop="createdAt" width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column v-if="!onlyPaid" :label="t('common.updatedAt')" prop="updatedAt" width="170">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
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
import {shortVideoApi} from '@/api'
import type {ShortVideoWatchRecord} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'

const props = withDefaults(defineProps<{
  userId: string
  active: boolean
  onlyPaid?: boolean
}>(), {
  onlyPaid: false,
})

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<ShortVideoWatchRecord[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const dateFilterLabel = computed(() =>
    props.onlyPaid
        ? t('pages.shortVideoWatchList.paidTime')
        : t('common.updatedAt'),
)

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    userId: props.userId,
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
    onlyPaid: props.onlyPaid,
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
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
    const response = await shortVideoApi.getShortVideoWatchList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load user short video watch/purchase logs:', error)
    ElMessage.error(t('pages.shortVideoWatchList.fetchFailed'))
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

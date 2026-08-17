<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.anchorList.dailyLiveDurationHint')"
        type="info"
    />

    <el-form :model="searchForm" class="search-form" inline label-width="90px">
      <el-form-item :label="t('pages.anchorList.dailyLiveDate')">
        <el-date-picker
            v-model="searchForm.dateRange"
            clearable
            :end-placeholder="t('pages.anchorList.dailyEndDate')"
            format="YYYY-MM-DD"
            :range-separator="t('pages.anchorList.dailyDateRangeSeparator')"
            :start-placeholder="t('pages.anchorList.dailyStartDate')"
            style="width: 260px"
            type="daterange"
            value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item :label="t('pages.anchorList.dailySettled')">
        <el-select v-model="searchForm.settled" :placeholder="t('common.all')" style="width: 140px">
          <el-option :label="t('common.all')" :value="-1"/>
          <el-option :label="t('pages.anchorList.dailySettledNo')" :value="0"/>
          <el-option :label="t('pages.anchorList.dailySettledYes')" :value="1"/>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatDuration(row.liveDuration) }}</template>
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

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.anchorList.noDailyEffectiveLiveData')"/>

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
import {reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorDailyEffectiveLiveItem} from '@/types/api'

const props = defineProps<{
  anchorId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<AnchorDailyEffectiveLiveItem[]>([])
const loaded = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
  settled: -1 as number,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => {
  const [liveDateStart, liveDateEnd] = searchForm.dateRange || []
  return {
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
    anchorId: props.anchorId,
    liveDateStart: liveDateStart || '',
    liveDateEnd: liveDateEnd || '',
    settled: searchForm.settled ?? -1,
  }
}

const fetchList = async () => {
  if (!props.anchorId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await accountApi.getAnchorDailyEffectiveLiveList(buildQueryParams())
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
  searchForm.settled = -1
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

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const formatDuration = (seconds?: number) => {
  if (seconds === undefined || seconds === null) return '-'
  const total = Math.max(0, Math.floor(seconds))
  const hours = Math.floor(total / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const secs = total % 60
  return t('pages.anchorList.durationFormat', {hours, minutes, seconds: secs})
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  pagination.pageIndex = 1
  pagination.total = 0
  searchForm.dateRange = []
  searchForm.settled = -1
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
.hint-alert {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

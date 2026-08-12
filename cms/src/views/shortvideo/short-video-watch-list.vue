<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ShortVideoWatchManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="searchForm.userId" clearable :placeholder="t('pages.shortVideoWatchList.userIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('common.updatedAt')">
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
        <el-table-column :label="t('pages.shortVideoWatchList.recordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.userId')" min-width="180" prop="userId"/>
        <el-table-column :label="t('common.nickname')" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.shortVideoWatchList.videoId')" min-width="180" prop="videoId"/>
        <el-table-column :label="t('pages.shortVideoWatchList.videoTitle')" min-width="160" prop="videoTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.videoTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.shortVideoWatchList.paidTime')" prop="paidTime" width="170">
          <template #default="{ row }">{{ row.paidTime || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" prop="createdAt" width="170"/>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170"/>
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
import {shortVideoApi} from '@/api/modules/shortVideo'
import type {ShortVideoWatchRecord} from '@/types/api'

interface SearchForm {
  userId: string
  dateRange: string[]
}

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<ShortVideoWatchRecord[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  userId: '',
  dateRange: [],
})

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
    userId: searchForm.userId.trim(),
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchWatchList = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoWatchList(buildQueryParams())
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch watch list failed:', error)
    ElMessage.error(t('pages.shortVideoWatchList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchWatchList()
}

const handleReset = () => {
  searchForm.userId = ''
  searchForm.dateRange = []
  currentPage.value = 1
  fetchWatchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchWatchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchWatchList()
}

onMounted(() => {
  fetchWatchList()
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
}

.search-form {
  margin-bottom: 15px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ProcessDownLogListManagement') }}</span>
          <el-button :loading="loading" @click="fetchList">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.processDownLog.pageDesc') }}</p>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('pages.processDownLog.host')">
          <el-input
              v-model="searchForm.host"
              clearable
              :placeholder="t('pages.processDownLog.hostPlaceholder')"
              style="width: 220px"
              @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('pages.processDownLog.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('common.endTime')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.processDownLog.dateRangeSeparator')"
              :start-placeholder="t('common.startTime')"
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
        <el-table-column :label="t('pages.processDownLog.host')" min-width="140" prop="host"/>
        <el-table-column :label="t('pages.processDownLog.appName')" width="140" prop="appName"/>
        <el-table-column :label="t('pages.processDownLog.missingSince')" width="170" prop="missingSince"/>
        <el-table-column :label="t('pages.processDownLog.detectedAt')" width="170" prop="detectedAt"/>
        <el-table-column :label="t('pages.processDownLog.coldStartOK')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.coldStartOK ? 'success' : 'danger'" size="small">
              {{ row.coldStartOK ? t('pages.processDownLog.coldStartSuccess') : t('pages.processDownLog.coldStartFailed') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.processDownLog.newPid')" width="100" prop="newPid"/>
        <el-table-column :label="t('pages.processDownLog.detail')" min-width="220" prop="detail" show-overflow-tooltip/>
        <el-table-column :label="t('pages.processDownLog.createdAt')" width="170" prop="createdAt"/>
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
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {processDownLogApi} from '@/api/modules/process-down-log'
import type {CMSProcessDownLogItem} from '@/types/api'
import {toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<CMSProcessDownLogItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  host: '',
  dateRange: [] as string[],
})

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    host: searchForm.host.trim() || undefined,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await processDownLogApi.listProcessDownLog(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load process down log list:', error)
    ElMessage.error(t('pages.processDownLog.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  void fetchList()
}

const handleReset = () => {
  searchForm.host = ''
  searchForm.dateRange = []
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

const formatRowIndex = (index: number) =>
    (currentPage.value - 1) * pageSize.value + index + 1

onMounted(() => {
  void fetchList()
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

.page-desc {
  margin: 0 0 16px;
  font-size: 13px;
  color: #64748b;
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

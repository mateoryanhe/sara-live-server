<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.SimulatorBlacklistLogListManagement') }}</span>
          <el-button :loading="loading" @click="fetchList">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.simulatorBlacklistLog.pageDesc') }}</p>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('pages.simulatorBlacklistLog.keyword')">
          <el-input
              v-model="searchForm.key"
              clearable
              :placeholder="t('pages.simulatorBlacklistLog.keywordPlaceholder')"
              style="width: 240px"
              @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('pages.simulatorBlacklistLog.loginType')">
          <el-select v-model="searchForm.loginType" clearable :placeholder="t('common.all')" style="width: 160px">
            <el-option :label="t('pages.simulatorBlacklistLog.loginTypeDevice')" value="device"/>
            <el-option :label="t('pages.simulatorBlacklistLog.loginTypeEmail')" value="email"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.simulatorBlacklistLog.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('common.endTime')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.simulatorBlacklistLog.dateRangeSeparator')"
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
        <el-table-column :label="t('pages.simulatorBlacklistLog.ip')" min-width="130" prop="ip"/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.deviceId')" min-width="240" prop="deviceId">
          <template #default="{ row }">
            <div class="device-id-cell">
              <span class="device-id-text" :title="row.deviceId || ''">{{ row.deviceId || '-' }}</span>
              <el-button
                  v-if="row.deviceId"
                  link
                  type="primary"
                  @click="copyDeviceId(row.deviceId)"
              >{{ t('common.copy') }}</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.simulatorBlacklistLog.registerAt')" width="170" prop="registerAt">
          <template #default="{ row }">{{ row.registerAt || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.simulatorBlacklistLog.loginType')" width="120">
          <template #default="{ row }">{{ formatLoginType(row.loginType) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.simulatorBlacklistLog.email')" min-width="180" prop="email" show-overflow-tooltip>
          <template #default="{ row }">{{ row.email || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.simulatorBlacklistLog.accountId')" min-width="160" prop="accountId">
          <template #default="{ row }">{{ row.accountId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.simulatorBlacklistLog.cpuModel')" min-width="160" prop="cpuModel" show-overflow-tooltip/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.deviceType')" width="100" prop="deviceType"/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.deviceModel')" min-width="140" prop="deviceModel" show-overflow-tooltip/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.appVersion')" width="110" prop="appVersion"/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.packageName')" min-width="160" prop="packageName" show-overflow-tooltip/>
        <el-table-column :label="t('pages.simulatorBlacklistLog.rejectedAt')" width="170" prop="createdAt"/>
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
import {simulatorBlacklistLogApi} from '@/api/modules/simulator-blacklist-log'
import type {CMSSimulatorBlacklistLogItem} from '@/types/api'
import {toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<CMSSimulatorBlacklistLogItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  key: '',
  loginType: '' as string,
  dateRange: [] as string[],
})

const formatLoginType = (loginType: string) => {
  if (loginType === 'device') {
    return t('pages.simulatorBlacklistLog.loginTypeDevice')
  }
  if (loginType === 'email') {
    return t('pages.simulatorBlacklistLog.loginTypeEmail')
  }
  return loginType || '-'
}

const copyDeviceId = async (deviceId: string | undefined) => {
  const value = String(deviceId ?? '').trim()
  if (!value) {
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(t('common.operationSuccess'))
  } catch (error) {
    console.error('copy device id failed:', error)
    ElMessage.error(t('common.failed'))
  }
}

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    key: searchForm.key.trim() || undefined,
    loginType: searchForm.loginType || undefined,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await simulatorBlacklistLogApi.listSimulatorBlacklistLog(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load simulator blacklist log list:', error)
    ElMessage.error(t('pages.simulatorBlacklistLog.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.loginType = ''
  searchForm.dateRange = []
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchList()
}

const formatRowIndex = (index: number) => (currentPage.value - 1) * pageSize.value + index + 1

onMounted(fetchList)
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-desc {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.search-form {
  margin-bottom: 12px;
}

.device-id-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

.device-id-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

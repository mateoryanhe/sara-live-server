<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildTransferManagement') }}</span>
        </div>
      </template>

      <el-alert :closable="false" :title="t('pages.guildTransferList.weekHint')" class="week-hint" type="info"/>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.guildTransferList.guildId')">
          <el-input v-model="searchForm.guildId" clearable :placeholder="t('pages.guildTransferList.enterGuildId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildTransferList.status')">
          <el-select v-model="searchForm.status" clearable style="width: 140px" :placeholder="t('pages.guildTransferList.statusAll')">
            <el-option :label="t('pages.guildTransferList.statusPending')" :value="0"/>
            <el-option :label="t('pages.guildTransferList.statusApproved')" :value="1"/>
            <el-option :label="t('pages.guildTransferList.statusTransferred')" :value="2"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('batchApprove')" :disabled="!selectedRows.length" type="success" @click="handleBatchApprove">
            {{ t('pages.guildTransferList.batchApprove') }}
          </el-button>
          <el-button v-if="can('batchTransfer')" :disabled="!selectedRows.length" type="warning" @click="handleBatchTransfer">
            {{ t('pages.guildTransferList.batchTransfer') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="loading"
          :data="tableData"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column fixed="left" type="selection" width="48"/>
        <el-table-column :label="t('common.createdAt')" fixed="left" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.guildId')" min-width="160" prop="guildId">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.guildName')" min-width="120" prop="guildName">
          <template #default="{ row }">{{ row.guildName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.settlementReceivableUsd')" align="right" min-width="140">
          <template #default="{ row }">
            <span class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.status')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferAt')" min-width="170">
          <template #default="{ row }">{{ formatDate(row.transferAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferCurrency')" min-width="90">
          <template #default="{ row }">{{ row.transferCurrency || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferPayeeName')" min-width="120">
          <template #default="{ row }">{{ row.transferPayeeName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferBankName')" min-width="120">
          <template #default="{ row }">{{ row.transferBankName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferAccountNo')" min-width="140">
          <template #default="{ row }">{{ row.transferAccountNo || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferBankCode')" min-width="110">
          <template #default="{ row }">{{ row.transferBankCode || '-' }}</template>
        </el-table-column>
      </el-table>

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
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'
import {
  getServerWeekDateRange,
  toServerDayEndUnix,
  toServerDayStartUnix,
} from '@/utils/server-datetime'

const {t} = useI18n()
const {can} = usePagePermission('GuildTransferManagement')

const loading = ref(false)
const tableData = ref<GuildIncomeSettlementLogItem[]>([])
const selectedRows = ref<GuildIncomeSettlementLogItem[]>([])

const fixedWeekRange = getServerWeekDateRange()

const searchForm = reactive({
  guildId: '',
  status: undefined as number | undefined,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => ({
  guildId: searchForm.guildId.trim(),
  status: searchForm.status === undefined || searchForm.status === null ? undefined : searchForm.status,
  // 固定查本周写入的上周结算，不提供时间筛选
  startTime: toServerDayStartUnix(fixedWeekRange.start),
  endTime: toServerDayEndUnix(fixedWeekRange.end),
  orderByReceivableUsdDesc: true,
  includeTransferInfo: true,
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildIncomeSettlementLogApi.getList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load guild transfer list:', error)
    ElMessage.error(t('pages.guildTransferList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.guildId = ''
  searchForm.status = undefined
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

const handleSelectionChange = (rows: GuildIncomeSettlementLogItem[]) => {
  selectedRows.value = rows || []
}

const handleBatchApprove = async () => {
  if (!selectedRows.value.length) {
    ElMessage.warning(t('pages.guildTransferList.selectRows'))
    return
  }
  const pendingIds = selectedRows.value
      .filter(row => row.status === 0)
      .map(row => row.id)
  if (!pendingIds.length) {
    ElMessage.warning(t('pages.guildTransferList.selectRows'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.approveConfirm', {count: pendingIds.length}),
        t('pages.guildTransferList.batchApprove'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    const res = await guildIncomeSettlementLogApi.batchApprove({ids: pendingIds})
    ElMessage.success(t('pages.guildTransferList.approveSuccess', {
      success: res.successCount || 0,
      fail: res.failCount || 0,
    }))
    await fetchList()
  } catch (error) {
    console.error('Batch approve failed:', error)
    ElMessage.error(t('pages.guildTransferList.approveFailed'))
  }
}

const handleBatchTransfer = async () => {
  if (!selectedRows.value.length) {
    ElMessage.warning(t('pages.guildTransferList.selectRows'))
    return
  }
  try {
    const res = await guildIncomeSettlementLogApi.batchTransfer({
      ids: selectedRows.value.map(row => row.id),
    })
    ElMessage.warning(res.message || t('pages.guildTransferList.transferReserved'))
  } catch (error) {
    console.error('Batch transfer reserved call failed:', error)
    ElMessage.warning(t('pages.guildTransferList.transferReserved'))
  }
}

const statusLabel = (status: number | undefined) => {
  if (status === 1) return t('pages.guildTransferList.statusApproved')
  if (status === 2) return t('pages.guildTransferList.statusTransferred')
  return t('pages.guildTransferList.statusPending')
}

const statusTagType = (status: number | undefined) => {
  if (status === 1) return 'success'
  if (status === 2) return 'info'
  return 'warning'
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.week-hint {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.money-amount {
  font-variant-numeric: tabular-nums;
}
</style>

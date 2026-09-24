<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorPayoutDetailList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.guildTransferList.transferAt')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.guildTransferList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.guildTransferList.dateRangeSeparator')"
              :start-placeholder="t('pages.guildTransferList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.guildTransferList.transferAt')" fixed="left" min-width="170">
          <template #default="{ row }">{{ formatDate(payoutTime(row)) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildIncomeSettlementLogList.logId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.roomAvatar"
                :preview-src-list="[row.roomAvatar]"
                :src="row.roomAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomNickname')" min-width="130">
          <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.status')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.settlementReceivableUsd')" align="right" min-width="150">
          <template #default="{ row }">
            <span class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferOrderId')" min-width="200">
          <template #default="{ row }">{{ row.transferOrderId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.thirdPartyOrderId')" min-width="200">
          <template #default="{ row }">{{ row.transferPlatformNo || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferCurrency')" min-width="100">
          <template #default="{ row }">{{ row.transferCurrency || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferLocalAmount')" align="right" min-width="130">
          <template #default="{ row }">
            <span class="money-amount">{{ formatWalletBalance(row.transferLocalAmount) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferFailMsg')" min-width="220">
          <template #default="{ row }">
            <span class="failure-message">{{ row.transferFailMsg || '-' }}</span>
          </template>
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
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate, toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const {can} = usePagePermission('PlatformAnchorPayoutDetailList')
const loading = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])

const searchForm = reactive({
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    directPayout: true,
    payoutOnly: true,
    transferStartTime: startDate ? toServerDayStartUnix(startDate) : 0,
    transferEndTime: endDate ? toServerDayEndUnix(endDate) : 0,
    pageIndex: pagination.pageIndex,
    pageSize: pagination.pageSize,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await anchorIncomeSettlementLogApi.getList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load platform anchor payout details:', error)
    ElMessage.error(t('pages.guildTransferList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const payoutTime = (row: AnchorIncomeSettlementLogItem) => row.transferAt || row.updatedAt || null

const statusLabel = (status: number | undefined) => {
  const value = Number(status)
  if (value === 1) return t('pages.guildTransferList.statusApproved')
  if (value === 2) return t('pages.guildTransferList.statusTransferred')
  if (value === 3) return t('pages.guildTransferList.statusTransferring')
  return t('pages.guildTransferList.statusPending')
}

const statusTagType = (status: number | undefined) => {
  const value = Number(status)
  if (value === 1) return 'success'
  if (value === 2) return 'info'
  if (value === 3) return ''
  return 'warning'
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

onMounted(fetchList)
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
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.failure-message {
  display: block;
  overflow-wrap: anywhere;
}
</style>

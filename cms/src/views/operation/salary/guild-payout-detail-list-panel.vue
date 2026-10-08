<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
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
        <el-table-column :label="t('pages.guildTransferList.guildId')" min-width="180" prop="guildId"/>
        <el-table-column :label="t('pages.guildTransferList.guildName')" min-width="130">
          <template #default="{ row }">{{ row.guildName || '-' }}</template>
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
import {type GuildPayoutDetailListScope, useGuildPayoutDetailList} from '@/composables/useGuildPayoutDetailList'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const props = defineProps<{
  scope: GuildPayoutDetailListScope
}>()

const {
  t,
  can,
  pageTitle,
  loading,
  tableData,
  searchForm,
  pagination,
  payoutTime,
  statusLabel,
  statusTagType,
  handleSearch,
  handleReset,
  handlePageChange,
  handleSizeChange,
} = useGuildPayoutDetailList(props.scope)
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

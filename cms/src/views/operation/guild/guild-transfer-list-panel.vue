<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-alert :closable="false" :title="listHint" class="week-hint" type="info"/>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.guildTransferList.guildId')">
          <el-input v-model="searchForm.guildId" clearable :placeholder="t('pages.guildTransferList.enterGuildId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildTransferList.status')">
          <el-select v-model="searchForm.status" clearable style="width: 140px" :placeholder="t('pages.guildTransferList.statusAll')">
            <el-option :label="t('pages.guildTransferList.statusAll')" :value="-1"/>
            <el-option :label="t('pages.guildTransferList.statusPending')" :value="0"/>
            <el-option :label="t('pages.guildTransferList.statusApproved')" :value="1"/>
            <el-option :label="t('pages.guildTransferList.statusTransferred')" :value="2"/>
            <el-option :label="t('pages.guildTransferList.statusTransferring')" :value="3"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('batchApprove')" :disabled="selectedPendingCount === 0" type="success" @click="handleBatchApprove">
            {{ t('pages.guildTransferList.batchApprove') }}
          </el-button>
          <el-button v-if="can('batchTransfer')" :disabled="selectedApprovedCount === 0" type="warning" @click="handleBatchTransfer">
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
          <template #default="{ row }">
            <el-button
                v-if="row.guildId && can('viewGuildDetail')"
                link
                type="primary"
                @click="openGuildDetail(row)"
            >
              {{ row.guildId }}
            </el-button>
            <span v-else>{{ row.guildId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.guildName')" min-width="120" prop="guildName">
          <template #default="{ row }">{{ row.guildName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.settlementReceivableUsd')" align="right" min-width="140">
          <template #default="{ row }">
            <el-button v-if="can('viewDetail')" link type="primary" @click="openSettlementDetail(row)">
              <span class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span>
            </el-button>
            <span v-else class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.status')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
            v-if="can('viewDetail') || can('reopenApproval') || can('copyPayout') || can('editReceivable')"
            :label="t('common.actions')"
            align="center"
            fixed="right"
            width="100"
        >
          <template #default="{ row }">
            <el-dropdown
                v-if="hasRowAction(row)"
                trigger="click"
                @command="(cmd: string) => handleRowCommand(row, cmd)"
            >
              <el-button
                  :loading="reopeningId === String(row.id) || copyingId === String(row.id)"
                  size="small"
                  type="primary"
              >
                {{ t('common.actions') }}
                <el-icon class="el-icon--right">
                  <ArrowDown/>
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="can('viewDetail')" command="viewDetail">
                    {{ t('pages.guildTransferList.viewDetail') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                      v-if="settlementStatus(row) === 1 && can('reopenApproval')"
                      command="reopenApproval"
                  >
                    {{ t('pages.guildTransferList.reopenApproval') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                      v-if="settlementStatus(row) === 2 && can('copyPayout')"
                      command="copyPayout"
                  >
                    {{ t('pages.guildTransferList.copyPayout') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                      v-if="settlementStatus(row) === 0 && can('editReceivable')"
                      command="editReceivable"
                  >
                    {{ t('pages.guildTransferList.editReceivable') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <span v-else>-</span>
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

    <el-dialog
        v-model="editReceivableVisible"
        :close-on-click-modal="false"
        :title="t('pages.guildTransferList.editReceivableTitle')"
        width="480px"
    >
      <el-form label-width="130px">
        <el-form-item :label="t('pages.guildTransferList.guildName')">
          <span>{{ editReceivableRow?.guildName || editReceivableRow?.guildId || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.guildTransferList.currentReceivable')">
          <span>{{ formatWalletBalance(editReceivableRow?.settlementReceivableUsd) }} USD</span>
        </el-form-item>
        <el-form-item :label="t('pages.guildTransferList.newReceivable')" required>
          <el-input-number
              v-model="editReceivableAmount"
              :max="999999999999.9999"
              :min="0.0001"
              :precision="4"
              :step="0.01"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :disabled="editReceivableSaving" @click="editReceivableVisible = false">
          {{ t('common.cancel') }}
        </el-button>
        <el-button :loading="editReceivableSaving" type="primary" @click="handleSaveReceivable">
          {{ t('common.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {ArrowDown} from '@element-plus/icons-vue'
import {type GuildTransferListScope, useGuildTransferList} from '@/composables/useGuildTransferList'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const props = defineProps<{
  scope: GuildTransferListScope
}>()

const {
  t,
  can,
  pageTitle,
  listHint,
  loading,
  tableData,
  reopeningId,
  copyingId,
  editReceivableVisible,
  editReceivableSaving,
  editReceivableRow,
  editReceivableAmount,
  searchForm,
  pagination,
  selectedPendingCount,
  selectedApprovedCount,
  settlementStatus,
  hasRowAction,
  handleSearch,
  handleReset,
  handlePageChange,
  handleSizeChange,
  handleSelectionChange,
  handleRowCommand,
  handleSaveReceivable,
  handleBatchApprove,
  handleBatchTransfer,
  statusLabel,
  statusTagType,
  openGuildDetail,
  openSettlementDetail,
} = useGuildTransferList(props.scope)
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

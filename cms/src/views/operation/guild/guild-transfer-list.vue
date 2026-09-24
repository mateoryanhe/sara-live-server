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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ArrowDown} from '@element-plus/icons-vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {guildIncomeSettlementLogApi} from '@/api/modules/guild-income-settlement-log'
import type {GuildIncomeSettlementLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {can} = usePagePermission()

const isCoinMerchantPage = computed(() => route.name === 'CoinMerchantGuildTransferManagement')
const guildType = computed(() => isCoinMerchantPage.value ? 1 : 0)
const pageTitle = computed(() => t(isCoinMerchantPage.value
    ? 'menu.CoinMerchantGuildTransferManagement'
    : 'menu.GuildTransferManagement'))
const listHint = computed(() => t('pages.guildTransferList.historyHint'))

const loading = ref(false)
const tableData = ref<GuildIncomeSettlementLogItem[]>([])
const selectedRows = ref<GuildIncomeSettlementLogItem[]>([])
const reopeningId = ref('')
const copyingId = ref('')
const editReceivableVisible = ref(false)
const editReceivableSaving = ref(false)
const editReceivableRow = ref<GuildIncomeSettlementLogItem | null>(null)
const editReceivableAmount = ref(0)
const settlementStatus = (row: GuildIncomeSettlementLogItem) => Number(row.status)
const selectedPendingRows = computed(() => selectedRows.value.filter(row => settlementStatus(row) === 0))
const selectedApprovedRows = computed(() => selectedRows.value.filter(row => settlementStatus(row) === 1))
const selectedPendingCount = computed(() => selectedPendingRows.value.length)
const selectedApprovedCount = computed(() => selectedApprovedRows.value.length)

const searchForm = reactive({
  guildId: '',
  status: -1 as number | undefined,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => ({
  guildId: searchForm.guildId.trim(),
  guildType: guildType.value,
  status: searchForm.status !== undefined && searchForm.status !== null && searchForm.status >= 0
      ? searchForm.status
      : undefined,
  // 普通工会和币商均保留历史待处理单，只隐藏本周以前已成功的记录。
  hideHistoricalTransferred: true,
  orderByReceivableUsdDesc: true,
  includeDetail: false,
  includeTransfer: true,
  includeTransferInfo: false,
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
  searchForm.status = -1
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

const hasRowAction = (row: GuildIncomeSettlementLogItem) => (
    can('viewDetail')
    || (settlementStatus(row) === 1 && can('reopenApproval'))
    || (settlementStatus(row) === 2 && can('copyPayout'))
    || (settlementStatus(row) === 0 && can('editReceivable'))
)

const handleReopenApproval = async (row: GuildIncomeSettlementLogItem) => {
  const id = String(row.id || '').trim()
  if (!id) return
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.reopenConfirm', {guild: row.guildName || row.guildId || '-'}),
        t('pages.guildTransferList.reopenApproval'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  reopeningId.value = id
  try {
    await guildIncomeSettlementLogApi.reopenApproval({id})
    ElMessage.success(t('pages.guildTransferList.reopenSuccess'))
    await fetchList()
  } catch (error) {
    console.error('Reopen guild settlement approval failed:', error)
    ElMessage.error(t('pages.guildTransferList.reopenFailed'))
  } finally {
    reopeningId.value = ''
  }
}

const handleCopyPayout = async (row: GuildIncomeSettlementLogItem) => {
  const id = String(row.id || '').trim()
  if (!id) return
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.copyPayoutConfirm', {guild: row.guildName || row.guildId || '-'}),
        t('pages.guildTransferList.copyPayout'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  copyingId.value = id
  try {
    const res = await guildIncomeSettlementLogApi.copyPayout({id})
    ElMessage.success(t('pages.guildTransferList.copyPayoutSuccess', {id: res.id || '-'}))
    searchForm.status = -1
    pagination.pageIndex = 1
    selectedRows.value = []
    await fetchList()
    const copiedId = String(res.id || '').trim()
    if (copiedId && !tableData.value.some(item => String(item.id) === copiedId)) {
      tableData.value.unshift({
        ...row,
        id: copiedId,
        status: 0,
        transferAt: undefined,
        transferOrderId: '',
        transferPlatformNo: '',
        transferLocalAmount: 0,
        transferCurrency: '',
        transferFailMsg: '',
        createdAt: new Date().toISOString(),
      })
      pagination.total += 1
    }
  } catch (error) {
    console.error('Copy guild settlement payout failed:', error)
    ElMessage.error(t('pages.guildTransferList.copyPayoutFailed'))
  } finally {
    copyingId.value = ''
  }
}

const openGuildDetail = (row: GuildIncomeSettlementLogItem) => {
  const guildId = String(row.guildId || '').trim()
  if (!guildId) return
  router.push({
    name: 'GuildDetail',
    query: {
      id: guildId,
      name: row.guildName || '',
    },
  })
}

const openSettlementDetail = (row: GuildIncomeSettlementLogItem) => {
  const id = String(row.id || '').trim()
  if (!id) return
  router.push({name: 'GuildTransferDetail', params: {id}, query: {guildType: guildType.value}})
}

const openEditReceivable = (row: GuildIncomeSettlementLogItem) => {
  editReceivableRow.value = row
  editReceivableAmount.value = Number(row.settlementReceivableUsd || 0)
  editReceivableVisible.value = true
}

const handleRowCommand = (row: GuildIncomeSettlementLogItem, command: string) => {
  switch (command) {
    case 'viewDetail':
      openSettlementDetail(row)
      break
    case 'reopenApproval':
      void handleReopenApproval(row)
      break
    case 'copyPayout':
      void handleCopyPayout(row)
      break
    case 'editReceivable':
      openEditReceivable(row)
      break
  }
}

const handleSaveReceivable = async () => {
  const row = editReceivableRow.value
  const id = String(row?.id || '').trim()
  const amount = Number(editReceivableAmount.value)
  if (!id || !Number.isFinite(amount) || amount <= 0) {
    ElMessage.warning(t('pages.guildTransferList.receivableInvalid'))
    return
  }
  editReceivableSaving.value = true
  try {
    await guildIncomeSettlementLogApi.updateReceivableUsd({
      id,
      settlementReceivableUsd: amount,
    })
    ElMessage.success(t('pages.guildTransferList.updateReceivableSuccess'))
    editReceivableVisible.value = false
    await fetchList()
  } catch (error) {
    console.error('Update guild settlement receivable failed:', error)
    ElMessage.error(t('pages.guildTransferList.updateReceivableFailed'))
  } finally {
    editReceivableSaving.value = false
  }
}

const handleBatchApprove = async () => {
  if (!selectedRows.value.length) {
    ElMessage.warning(t('pages.guildTransferList.selectRows'))
    return
  }
  const pendingIds = selectedPendingRows.value
      .map(row => String(row.id || '').trim())
      .filter(Boolean)
  if (!pendingIds.length) {
    ElMessage.warning(t('pages.guildTransferList.selectPendingRows'))
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
  const approvedIds = selectedApprovedRows.value
      .map(row => String(row.id || '').trim())
      .filter(Boolean)
  if (!approvedIds.length) {
    ElMessage.warning(t('pages.guildTransferList.selectApprovedRows'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.transferConfirm', {count: approvedIds.length}),
        t('pages.guildTransferList.batchTransfer'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    const res = await guildIncomeSettlementLogApi.batchTransfer({ids: approvedIds})
    const success = res.successCount || 0
    const fail = res.failCount || 0
    if (success > 0) {
      ElMessage.success(res.message || t('pages.guildTransferList.transferSuccess', {success, fail}))
    } else {
      ElMessage.warning(res.message || t('pages.guildTransferList.transferFailed'))
    }
    await fetchList()
  } catch (error) {
    console.error('Batch transfer failed:', error)
    ElMessage.error(t('pages.guildTransferList.transferFailed'))
  }
}

const statusLabel = (status: number | undefined) => {
  const normalizedStatus = Number(status)
  if (normalizedStatus === 1) return t('pages.guildTransferList.statusApproved')
  if (normalizedStatus === 2) return t('pages.guildTransferList.statusTransferred')
  if (normalizedStatus === 3) return t('pages.guildTransferList.statusTransferring')
  return t('pages.guildTransferList.statusPending')
}

const statusTagType = (status: number | undefined) => {
  const normalizedStatus = Number(status)
  if (normalizedStatus === 1) return 'success'
  if (normalizedStatus === 2) return 'info'
  if (normalizedStatus === 3) return ''
  return 'warning'
}

onMounted(() => {
  fetchList()
})

watch(() => route.name, (next, previous) => {
  if (next === previous) return
  searchForm.guildId = ''
  searchForm.status = -1
  pagination.pageIndex = 1
  selectedRows.value = []
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

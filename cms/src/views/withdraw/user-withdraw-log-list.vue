<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.UserWithdrawLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('pages.userWithdrawLogList.userId')">
          <el-input
              v-model="searchForm.userId"
              clearable
              :placeholder="t('pages.userWithdrawLogList.userIdPlaceholder')"
              style="width: 220px"
          />
        </el-form-item>
        <el-form-item :label="t('pages.userWithdrawLogList.status')">
          <el-select v-model="searchForm.status" clearable style="width: 160px">
            <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value"/>
          </el-select>
        </el-form-item>
        <el-form-item class="search-form-actions">
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <span v-if="showBatchActions" class="action-divider"/>
          <el-button
              v-if="canApprove"
              type="primary"
              plain
              :disabled="!selectedPendingIds.length"
              @click="handleBatchApprove"
          >{{ t('pages.userWithdrawLogList.batchApprove') }}</el-button>
          <el-button
              v-if="canReject"
              plain
              :disabled="!selectedPendingIds.length"
              @click="openRejectDialog()"
          >{{ t('pages.userWithdrawLogList.batchReject') }}</el-button>
          <el-button
              v-if="canSubmitPayout"
              type="success"
              :disabled="!selectedApprovedIds.length"
              @click="handleBatchSubmitPayout"
          >{{ t('pages.userWithdrawLogList.batchSubmitPayout') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table
          v-loading="loading"
          :data="tableData"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48"/>
        <el-table-column :label="t('pages.userWithdrawLogList.id')" min-width="170" prop="id"/>
        <el-table-column :label="t('pages.userWithdrawLogList.userId')" min-width="150" prop="userId">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.userId"
                link
                type="primary"
                @click="openUserDetail(row.userId)"
            >
              {{ row.userId }}
            </el-button>
            <span v-else>{{ row.userId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.userWithdrawLogList.nickname')" min-width="120" prop="nickname"/>
        <el-table-column :label="t('pages.userWithdrawLogList.diamondAmount')" min-width="110">
          <template #default="{ row }">
            <span v-if="row.diamondAmount > 0">{{ row.diamondAmount }}</span>
            <span v-else-if="row.goldAmount > 0" class="legacy-gold">
              {{ row.goldAmount }} ({{ t('pages.userWithdrawLogList.goldAmountLegacy') }})
            </span>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.userWithdrawLogList.withdrawUsd')" min-width="110" prop="withdrawUsd"/>
        <el-table-column :label="t('pages.userWithdrawLogList.status')" min-width="120">
          <template #default="{ row }">{{ statusLabel(row.status) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.userWithdrawLogList.transferCurrency')" min-width="100" prop="transferCurrency"/>
        <el-table-column :label="t('pages.userWithdrawLogList.transferLocalAmount')" min-width="120" prop="transferLocalAmount"/>
        <el-table-column :label="t('pages.userWithdrawLogList.failMsg')" min-width="160" prop="transferFailMsg" show-overflow-tooltip/>
        <el-table-column :label="t('pages.userWithdrawLogList.createdAt')" min-width="160" prop="createdAt"/>
        <el-table-column
            v-if="showRowAuditActions"
            fixed="right"
            :label="t('pages.userWithdrawLogList.actions')"
            width="168"
        >
          <template #default="{ row }">
            <template v-if="canAuditRow(row)">
              <el-button v-if="canApprove" link type="primary" @click="handleRowApprove(row)">
                {{ t('pages.userWithdrawLogList.approve') }}
              </el-button>
              <el-button v-if="canReject" link type="danger" @click="openRejectDialog([row.id])">
                {{ t('pages.userWithdrawLogList.reject') }}
              </el-button>
            </template>
            <el-button
                v-else-if="canRevokeRow(row) && canRevokeApprove"
                link
                type="warning"
                @click="handleRowRevoke(row)"
            >
              {{ t('pages.userWithdrawLogList.revoke') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchList"
            @size-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="rejectDialogVisible" :title="t('pages.userWithdrawLogList.rejectTitle')" width="480px">
      <el-input v-model="rejectReason" type="textarea" :rows="3" :placeholder="t('pages.userWithdrawLogList.rejectReasonPlaceholder')"/>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleReject">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import userWithdrawApi, {type UserWithdrawLogItem} from '@/api/modules/user-withdraw'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {hasButtonPermission} from '@/utils/permission'

const {t} = useI18n()
const {canViewUserDetail, openUserDetail} = useUserDetailNav('UserWithdrawLogList')
const loading = ref(false)
const tableData = ref<UserWithdrawLogItem[]>([])
const selectedRows = ref<UserWithdrawLogItem[]>([])
const rejectDialogVisible = ref(false)
const rejectReason = ref('')
const rejectTargetIds = ref<string[]>([])

const WITHDRAW_STATUS_PENDING = 0
const WITHDRAW_STATUS_APPROVED = 1

const searchForm = reactive<{userId: string; status: number | undefined}>({
  userId: '',
  status: undefined,
})

const pagination = reactive({page: 1, pageSize: 20, total: 0})

const canApprove = computed(() => hasButtonPermission('UserWithdrawLogList', 'batchApprove'))
const canReject = computed(() => hasButtonPermission('UserWithdrawLogList', 'batchReject'))
const canRevokeApprove = computed(() => hasButtonPermission('UserWithdrawLogList', 'batchRevokeApprove'))
const canSubmitPayout = computed(() => hasButtonPermission('UserWithdrawLogList', 'batchSubmitPayout'))
const showBatchActions = computed(() =>
    canApprove.value || canReject.value || canRevokeApprove.value || canSubmitPayout.value)
const showRowAuditActions = computed(() =>
    canApprove.value || canReject.value || canRevokeApprove.value)

const canAuditRow = (row: UserWithdrawLogItem) => Number(row.status) === WITHDRAW_STATUS_PENDING
const canRevokeRow = (row: UserWithdrawLogItem) => Number(row.status) === WITHDRAW_STATUS_APPROVED

const selectedPendingIds = computed(() =>
    selectedRows.value.filter(canAuditRow).map((row) => row.id).filter(Boolean))
const selectedApprovedIds = computed(() =>
    selectedRows.value.filter(canRevokeRow).map((row) => row.id).filter(Boolean))

const statusOptions = computed(() => [
  {value: 0, label: t('pages.userWithdrawLogList.statusPending')},
  {value: 1, label: t('pages.userWithdrawLogList.statusApproved')},
  {value: 2, label: t('pages.userWithdrawLogList.statusTransferred')},
  {value: 3, label: t('pages.userWithdrawLogList.statusTransferring')},
  {value: 4, label: t('pages.userWithdrawLogList.statusRejected')},
])

const statusLabel = (status: number) => statusOptions.value.find((o) => o.value === status)?.label ?? String(status)

const fetchList = async () => {
  loading.value = true
  try {
    const res = await userWithdrawApi.listWithdrawLogs({
      userId: searchForm.userId.trim() || undefined,
      status: searchForm.status,
      page: pagination.page,
      pageSize: pagination.pageSize,
    })
    tableData.value = res.list || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.userWithdrawLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.userId = ''
  searchForm.status = undefined
  handleSearch()
}

const handleSelectionChange = (rows: UserWithdrawLogItem[]) => {
  selectedRows.value = rows
}

const handleBatchApprove = async () => {
  const ids = selectedPendingIds.value
  if (!ids.length) {
    ElMessage.warning(t('pages.userWithdrawLogList.batchNeedPending'))
    return
  }
  try {
    const res = await userWithdrawApi.batchApproveWithdraw({ids})
    ElMessage.success(t('pages.userWithdrawLogList.batchResult', {ok: res.successCount, fail: res.failCount}))
    await fetchList()
  } catch (e) {
    console.error(e)
  }
}

const openRejectDialog = (ids?: string[]) => {
  const targetIds = ids?.length ? ids : selectedPendingIds.value
  if (!targetIds.length) {
    ElMessage.warning(t('pages.userWithdrawLogList.batchNeedPending'))
    return
  }
  rejectTargetIds.value = targetIds
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

const handleReject = async () => {
  rejectDialogVisible.value = false
  const ids = rejectTargetIds.value
  if (!ids.length) {
    return
  }
  try {
    const res = await userWithdrawApi.batchRejectWithdraw({ids, reason: rejectReason.value})
    ElMessage.success(t('pages.userWithdrawLogList.batchResult', {ok: res.successCount, fail: res.failCount}))
    await fetchList()
  } catch (e) {
    console.error(e)
  }
}

const handleRowRevoke = async (row: UserWithdrawLogItem) => {
  if (!row.id || !canRevokeRow(row)) {
    return
  }
  try {
    await ElMessageBox.confirm(t('pages.userWithdrawLogList.revokeConfirm'), t('common.confirm'), {type: 'warning'})
    const res = await userWithdrawApi.batchRevokeWithdrawApproval({ids: [row.id]})
    ElMessage.success(t('pages.userWithdrawLogList.batchResult', {ok: res.successCount, fail: res.failCount}))
    await fetchList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleRowApprove = async (row: UserWithdrawLogItem) => {
  if (!row.id || !canAuditRow(row)) {
    return
  }
  try {
    await ElMessageBox.confirm(t('pages.userWithdrawLogList.approveConfirm'), t('common.confirm'), {type: 'warning'})
    const res = await userWithdrawApi.batchApproveWithdraw({ids: [row.id]})
    ElMessage.success(t('pages.userWithdrawLogList.batchResult', {ok: res.successCount, fail: res.failCount}))
    await fetchList()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleBatchSubmitPayout = async () => {
  const ids = selectedApprovedIds.value
  if (!ids.length) {
    ElMessage.warning(t('pages.userWithdrawLogList.batchNeedApproved'))
    return
  }
  try {
    await ElMessageBox.confirm(t('pages.userWithdrawLogList.submitPayoutConfirm'), t('common.confirm'), {type: 'warning'})
    const res = await userWithdrawApi.batchSubmitPayout({ids})
    ElMessage.success(res.message || t('pages.userWithdrawLogList.batchResult', {ok: res.successCount, fail: res.failCount}))
    await fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
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

.search-form :deep(.el-form-item) {
  margin-bottom: 12px;
  margin-right: 12px;
}

.search-form-actions :deep(.el-form-item__content) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.action-divider {
  display: inline-block;
  width: 1px;
  height: 20px;
  margin: 0 4px;
  background: var(--el-border-color);
  vertical-align: middle;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

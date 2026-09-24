<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorPayoutList') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.guildTransferList.historyHint')"
          class="history-hint"
          type="info"
      />

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('pages.liveRecordList.platformAnchor')">
          <div class="anchor-filter anchor-filter--compact">
            <el-input
                :model-value="platformAnchorInputValue"
                class="anchor-input"
                disabled
                :placeholder="t('pages.liveRecordList.noPlatformAnchorSelected')"
            />
            <el-button size="small" @click="openPlatformAnchorPicker">
              {{ t('pages.liveRecordList.selectPlatformAnchor') }}
            </el-button>
            <el-button v-if="selectedPlatformAnchors.length > 0" link size="small" @click="clearPlatformAnchors">
              {{ t('pages.liveRecordList.clearPlatformAnchor') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.startDate')">
          <el-date-picker
              v-model="searchForm.startDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.startDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.endDate')">
          <el-date-picker
              v-model="searchForm.endDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.endDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildTransferList.status')">
          <el-select v-model="searchForm.status" style="width: 140px">
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
          <el-button
              v-if="can('batchApprove')"
              :disabled="selectedPendingRows.length === 0"
              type="success"
              @click="handleBatchApprove"
          >
            {{ t('pages.guildTransferList.batchApprove') }}
          </el-button>
          <el-button
              v-if="can('batchTransfer')"
              :disabled="selectedApprovedRows.length === 0"
              type="warning"
              @click="handleBatchTransfer"
          >
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
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomId')" min-width="180">
          <template #default="{ row }">
            <el-button v-if="row.roomId" link type="primary" @click="openAnchorDetail(row.roomId)">
              {{ row.roomId }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.roomNickname')" min-width="120">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.roomId && row.roomNickname"
                link
                type="primary"
                @click="openUserDetail(row.roomId)"
            >
              {{ row.roomNickname }}
            </el-button>
            <span v-else>{{ row.roomNickname || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.settlementReceivableUsd')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementReceivableUsd) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.status')" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferCurrency')" min-width="100">
          <template #default="{ row }">
            <el-button
                v-if="row.roomId && can('transferInfo')"
                link
                type="primary"
                @click="openTransferInfo(row)"
            >
              {{ row.transferCurrency || t('common.edit') }}
            </el-button>
            <span v-else>{{ row.transferCurrency || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferLocalAmount')" align="right" min-width="130">
          <template #default="{ row }">{{ row.transferLocalAmount ? formatWalletBalance(row.transferLocalAmount) : '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.thirdPartyOrderId')" min-width="190" show-overflow-tooltip>
          <template #default="{ row }">{{ row.transferPlatformNo || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferAt')" min-width="170">
          <template #default="{ row }">{{ formatDate(row.transferAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.guildTransferList.transferFailMsg')" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">{{ row.transferFailMsg || '-' }}</template>
        </el-table-column>
        <el-table-column
            :label="t('pages.anchorIncomeSettlementLogList.settlementSalary')"
            align="right"
            label-class-name="header-nowrap"
            min-width="150"
        >
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementSalary) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.settlementFlowCommission')" align="right" min-width="110">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.settlementShareAmount) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.hasSalary')" align="center" min-width="100">
          <template #default="{ row }">
            <el-tag :type="row.hasSalary ? 'success' : 'info'">
              {{ row.hasSalary ? t('common.yes') : t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSocialSharePercent')" align="right" min-width="150">
          <template #default="{ row }">{{ formatSharePercent(row.anchorSocialSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSocialShareAmount')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.anchorSocialShareAmount) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.guildSocialSharePercent')" align="right" min-width="150">
          <template #default="{ row }">{{ formatSharePercent(row.guildSocialSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.guildSocialShareAmount')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.guildSocialShareAmount) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorGameSharePercent')" align="right" min-width="150">
          <template #default="{ row }">{{ formatSharePercent(row.anchorGameSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorGameShareAmountGold')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.anchorGameShareAmountGold) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.guildGameSharePercent')" align="right" min-width="150">
          <template #default="{ row }">{{ formatSharePercent(row.guildGameSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.guildGameShareAmountGold')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.guildGameShareAmountGold) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalSocialIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalSocialIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalGiftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalPaidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
        </el-table-column>
		<el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallTicketIncome')" align="right" min-width="150">
		  <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
		</el-table-column>
		<el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalVideoCallBillingIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalShortVideoIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalShortVideoIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalGameIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.totalLiveDuration')" min-width="120">
          <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorIncomeSettlementLogList.anchorSharePercent')" min-width="110" prop="anchorSharePercent">
          <template #default="{ row }">{{ formatSharePercent(row.anchorSharePercent) }}</template>
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

    <AnchorPickerDialog
        v-model:visible="platformAnchorPickerVisible"
        :initial-anchors="selectedPlatformAnchors"
        multiple
        platform-only
        @confirm-multiple="handlePlatformAnchorsPicked"
    />

  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox} from 'element-plus'
import {anchorIncomeSettlementLogApi} from '@/api/modules/anchor-income-settlement-log'
import type {AnchorIncomeSettlementLogItem, AnchorListItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {formatServerDateTime as formatDate, toServerDayStartUnix, toServerDayEndUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('PlatformAnchorPayoutList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('PlatformAnchorPayoutList')
const loading = ref(false)
const tableData = ref<AnchorIncomeSettlementLogItem[]>([])
const selectedRows = ref<AnchorIncomeSettlementLogItem[]>([])
const settlementStatus = (row: AnchorIncomeSettlementLogItem) => Number(row.status)
const selectedPendingRows = computed(() => selectedRows.value.filter(row => row.directPayout && settlementStatus(row) === 0))
const selectedApprovedRows = computed(() => selectedRows.value.filter(row => row.directPayout && settlementStatus(row) === 1))
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)

const searchForm = reactive({
  startDate: '',
  endDate: '',
  status: -1,
})

const formatPlatformAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  return `${nickname} (${anchor.id})`
}

const platformAnchorInputValue = computed(() => {
  const anchor = selectedPlatformAnchors.value[0]
  if (!anchor) {
    return ''
  }
  if (selectedPlatformAnchors.value.length === 1) {
    return formatPlatformAnchorLabel(anchor)
  }
  return t('pages.liveRecordList.selectedPlatformAnchorsCount', {count: selectedPlatformAnchors.value.length})
})

const openPlatformAnchorPicker = () => {
  platformAnchorPickerVisible.value = true
}

const clearPlatformAnchors = () => {
  selectedPlatformAnchors.value = []
}

const handlePlatformAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedPlatformAnchors.value = anchors
}

const buildSelectedAnchorIds = () => {
  return selectedPlatformAnchors.value.map(anchor => String(anchor.id))
}

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildFilterParams = () => ({
  anchorIds: buildSelectedAnchorIds(),
  startTime: searchForm.startDate ? toServerDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toServerDayEndUnix(searchForm.endDate) : 0,
  status: searchForm.status >= 0 ? searchForm.status : undefined,
  directPayout: true,
  hideHistoricalTransferred: true,
  orderByReceivableUsdDesc: true,
  includeTransferInfo: true,
})

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await anchorIncomeSettlementLogApi.getList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch anchor income settlement log list failed:', error)
    ElMessage.error(t('pages.anchorIncomeSettlementLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  selectedPlatformAnchors.value = []
  searchForm.startDate = ''
  searchForm.endDate = ''
  searchForm.status = -1
  pagination.pageIndex = 1
  fetchList()
}

const handleSelectionChange = (rows: AnchorIncomeSettlementLogItem[]) => {
  selectedRows.value = rows || []
}

const handleBatchApprove = async () => {
  const ids = selectedPendingRows.value.map(row => String(row.id || '').trim()).filter(Boolean)
  if (!ids.length) {
    ElMessage.warning(t('pages.guildTransferList.selectPendingRows'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.approveConfirm', {count: ids.length}),
        t('pages.guildTransferList.batchApprove'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    const res = await anchorIncomeSettlementLogApi.batchApprove({ids})
    ElMessage.success(t('pages.guildTransferList.approveSuccess', {
      success: res.successCount || 0,
      fail: res.failCount || 0,
    }))
    await fetchList()
  } catch (error) {
    console.error('Batch approve platform anchor settlement failed:', error)
    ElMessage.error(t('pages.guildTransferList.approveFailed'))
  }
}

const handleBatchTransfer = async () => {
  const ids = selectedApprovedRows.value.map(row => String(row.id || '').trim()).filter(Boolean)
  if (!ids.length) {
    ElMessage.warning(t('pages.guildTransferList.selectApprovedRows'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.guildTransferList.transferConfirm', {count: ids.length}),
        t('pages.guildTransferList.batchTransfer'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    const res = await anchorIncomeSettlementLogApi.batchTransfer({ids})
    const success = res.successCount || 0
    const fail = res.failCount || 0
    if (success > 0) {
      ElMessage.success(res.message || t('pages.guildTransferList.transferSuccess', {success, fail}))
    } else {
      ElMessage.warning(res.message || t('pages.guildTransferList.transferFailed'))
    }
    await fetchList()
  } catch (error) {
    console.error('Batch transfer platform anchor settlement failed:', error)
    ElMessage.error(t('pages.guildTransferList.transferFailed'))
  }
}

const openTransferInfo = (row: AnchorIncomeSettlementLogItem) => {
  const anchorId = String(row.roomId || '').trim()
  if (!anchorId) return
  router.push({
    name: 'PlatformAnchorTransferInfoEdit',
    params: {anchorId},
    query: {
      nickname: row.roomNickname || '',
      from: 'anchorSettlement',
    },
  })
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

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const formatSharePercent = (value: number | null | undefined) => {
  if (value == null || Number.isNaN(value)) return '-'
  return `${value}%`
}

const openAnchorDetail = (anchorId: string | number) => {
  if (!anchorId) {
    return
  }
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(anchorId)},
  })
}

onMounted(() => {
  fetchList()
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

.history-hint {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 12px;
}

.anchor-filter--compact {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.anchor-input {
  width: 220px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

:deep(th.header-nowrap > .cell) {
  white-space: nowrap;
}
</style>

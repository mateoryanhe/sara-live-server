<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.RechargeOrderList') }}</span>
          <el-button v-if="canManualCreateOrder" type="primary" @click="openCreateOrderDialog">{{ t('pages.rechargeOrderList.manualCreateOrder') }}</el-button>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item :label="t('pages.rechargeOrderList.orderId')">
          <el-input v-model="searchForm.orderId" clearable :placeholder="t('pages.rechargeOrderList.enterOrderId')"/>
        </el-form-item>
        <el-form-item :label="t('common.userId')">
          <el-input v-model="searchForm.userId" clearable :placeholder="t('pages.rechargeOrderList.enterUserId')"/>
        </el-form-item>
        <el-form-item :label="t('common.createdAt')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.rechargeOrderList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.rechargeOrderList.dateRangeSeparator')"
              :start-placeholder="t('pages.rechargeOrderList.startDate')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.rechargeOrderList.orderStatus')">
          <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.rechargeOrderList.statusPending')"/>
            <el-option :value="2" :label="t('pages.rechargeOrderList.statusCompleted')"/>
            <el-option :value="3" :label="t('pages.rechargeOrderList.statusCancelled')"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <div class="list-summary">{{ t('pages.rechargeOrderList.totalOrders', {count: total}) }}</div>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
        <el-table-column :label="t('pages.rechargeOrderList.orderId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.userId')" min-width="180" prop="userId">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail && row.userId" link type="primary" @click="openUserDetail(row.userId)">
              {{ row.userId }}
            </el-button>
            <span v-else>{{ row.userId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.userNickname')" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.cfgId')" min-width="120" prop="cfgId">
          <template #default="{ row }">{{ row.cfgId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.price')" prop="price" width="120">
          <template #default="{ row }">{{ formatAmount(row.price) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.gold')" prop="gold" width="120">
          <template #default="{ row }">{{ formatAmount(row.gold) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.source')" width="100">
          <template #default="{ row }">{{ sourceLabel(row.source) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.thirdOrderId')" min-width="160" prop="thirdOrderId" show-overflow-tooltip>
          <template #default="{ row }">{{ row.thirdOrderId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.remark')" min-width="160" prop="remark" show-overflow-tooltip>
          <template #default="{ row }">{{ row.remark || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" width="170">
          <template #default="{ row }">{{ formatUnixTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.rechargeOrderList.paidAt')" width="170">
          <template #default="{ row }">{{ formatUnixTime(row.paidAt) }}</template>
        </el-table-column>
        <el-table-column v-if="canManualRecharge" fixed="right" :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button
                v-if="row.status === 0"
                :loading="manualRechargingId === row.id"
                size="small"
                type="primary"
                @click="handleManualRecharge(row)"
            >
              {{ t('pages.rechargeOrderList.manualRecharge') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
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

    <el-dialog v-model="createOrderDialogVisible" :title="t('pages.rechargeOrderList.createOrderDialogTitle')" width="420px" @closed="resetCreateOrderForm">
      <el-form ref="createOrderFormRef" :model="createOrderForm" :rules="createOrderRules" label-width="90px">
        <el-form-item :label="t('pages.rechargeOrderList.playerId')" prop="userId">
          <el-input v-model="createOrderForm.userId" clearable :placeholder="t('pages.rechargeOrderList.enterPlayerId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.rechargeOrderList.orderAmount')" prop="amount">
          <el-input-number
              v-model="createOrderForm.amount"
              :min="0.01"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="1"
              controls-position="right"
              :placeholder="t('pages.rechargeOrderList.enterOrderAmount')"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createOrderDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="creatingOrder" type="primary" @click="handleCreateOrder">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {rechargeOrderApi} from '@/api'
import type {RechargeOrder} from '@/types/api.ts'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {formatAmount, NUMBER_INPUT_DECIMALS} from '@/utils/number-format'

interface SearchForm {
  orderId: string
  userId: string
  dateRange: string[]
  statusFilter: number
}

const {t} = useI18n()
const {can} = usePagePermission('RechargeOrderList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('RechargeOrderList')
const canManualCreateOrder = computed(() => can('manualCreateOrder'))
const canManualRecharge = computed(() => can('manualRecharge'))
const loading = ref(false)
const manualRechargingId = ref('')
const createOrderDialogVisible = ref(false)
const creatingOrder = ref(false)
const createOrderFormRef = ref<FormInstance>()
const tableData = ref<RechargeOrder[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  orderId: '',
  userId: '',
  dateRange: [],
  statusFilter: 0,
})

const createOrderForm = reactive({
  userId: '',
  amount: undefined as number | undefined,
})

const createOrderRules = computed<FormRules>(() => ({
  userId: [{required: true, message: t('pages.rechargeOrderList.playerIdRequired'), trigger: 'blur'}],
  amount: [{required: true, message: t('pages.rechargeOrderList.orderAmountRequired'), trigger: 'change'}],
}))

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
    orderId: searchForm.orderId.trim(),
    userId: searchForm.userId.trim(),
    statusFilter: searchForm.statusFilter,
    startTime: startDate ? toDayStartUnix(startDate) : 0,
    endTime: endDate ? toDayEndUnix(endDate) : 0,
  }
}

const fetchOrderList = async () => {
  loading.value = true
  try {
    const response = await rechargeOrderApi.getRechargeOrderList(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load recharge orders:', error)
    ElMessage.error(t('pages.rechargeOrderList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchOrderList()
}

const handleReset = () => {
  searchForm.orderId = ''
  searchForm.userId = ''
  searchForm.dateRange = []
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchOrderList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchOrderList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchOrderList()
}

const formatRowIndex = (index: number) =>
    (currentPage.value - 1) * pageSize.value + index + 1

const afterManualRechargeSuccess = async (after: number) => {
  ElMessage.success(t('pages.rechargeOrderList.manualRechargeSuccess', {balance: formatAmount(after)}))
  await fetchOrderList()
}

const handleManualRecharge = async (row: RechargeOrder) => {
  try {
    await ElMessageBox.confirm(
        t('pages.rechargeOrderList.manualRechargeConfirm', {
          orderId: row.id,
          gold: formatAmount(row.gold),
          player: row.nickname || row.userId,
        }),
        t('pages.rechargeOrderList.manualRechargeTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
    )
    manualRechargingId.value = row.id
    const res = await rechargeOrderApi.manualRecharge(row.id)
    await afterManualRechargeSuccess(res.after)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Manual recharge failed:', error)
    }
  } finally {
    manualRechargingId.value = ''
  }
}

const openCreateOrderDialog = () => {
  createOrderDialogVisible.value = true
}

const resetCreateOrderForm = () => {
  createOrderForm.userId = ''
  createOrderForm.amount = undefined
  createOrderFormRef.value?.clearValidate()
}

const handleCreateOrder = async () => {
  if (!createOrderFormRef.value) return
  try {
    await createOrderFormRef.value.validate()
  } catch {
    return
  }

  creatingOrder.value = true
  try {
    const res = await rechargeOrderApi.manualCreateOrder({
      userId: createOrderForm.userId.trim(),
      amount: Number(createOrderForm.amount),
    })
    createOrderDialogVisible.value = false
    ElMessage.success(t('pages.rechargeOrderList.createOrderSuccess', {
      orderId: res.orderId,
      price: formatAmount(res.price),
      gold: formatAmount(res.gold),
    }))
    currentPage.value = 1
    await fetchOrderList()
  } catch (error) {
    console.error('Manual create order failed:', error)
  } finally {
    creatingOrder.value = false
  }
}

const formatUnixTime = (ts: number | null | undefined) => {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString()
}

const statusLabel = (status: number) => {
  switch (status) {
    case 1:
      return t('pages.rechargeOrderList.statusCompleted')
    case 2:
      return t('pages.rechargeOrderList.statusCancelled')
    default:
      return t('pages.rechargeOrderList.statusPending')
  }
}

const statusTagType = (status: number) => {
  switch (status) {
    case 1:
      return 'success'
    case 2:
      return 'info'
    default:
      return 'warning'
  }
}

const sourceLabel = (source: number) => {
  return source === 2 ? t('pages.rechargeOrderList.sourceManual') : t('pages.rechargeOrderList.sourceApp')
}

onMounted(() => {
  fetchOrderList()
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
  font-size: 16px;
  font-weight: bold;
}

.search-form {
  margin-bottom: 12px;
}

.list-summary {
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>

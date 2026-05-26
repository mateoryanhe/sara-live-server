<template>
  <div
      v-loading="pageWaiting"
      class="page-container"
      element-loading-background="rgba(255, 255, 255, 0.75)"
      element-loading-text="pageWaitingText"
  >
    <el-card>
      <template #header>
        <div class="card-header">
          <span>充值订单</span>
          <el-button type="primary" @click="openCreateOrderDialog">人工创建订单</el-button>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="90px">
        <el-form-item label="订单ID">
          <el-input v-model="searchForm.orderId" clearable placeholder="请输入订单ID"/>
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" clearable placeholder="请输入用户ID"/>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              range-separator="至"
              start-placeholder="开始日期"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchForm.statusFilter" placeholder="全部" style="width: 140px">
            <el-option :value="0" label="全部"/>
            <el-option :value="1" label="待支付"/>
            <el-option :value="2" label="已完成"/>
            <el-option :value="3" label="已取消"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="订单ID" min-width="180" prop="id"/>
        <el-table-column label="用户ID" min-width="180" prop="userId"/>
        <el-table-column label="用户昵称" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="档位ID" min-width="120" prop="cfgId">
          <template #default="{ row }">{{ row.cfgId || '-' }}</template>
        </el-table-column>
        <el-table-column label="金额" prop="price" width="120">
          <template #default="{ row }">{{ formatAmount(row.price) }}</template>
        </el-table-column>
        <el-table-column label="金币" prop="gold" width="120">
          <template #default="{ row }">{{ formatAmount(row.gold) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="来源" width="100">
          <template #default="{ row }">{{ sourceLabel(row.source) }}</template>
        </el-table-column>
        <el-table-column label="第三方订单号" min-width="160" prop="thirdOrderId" show-overflow-tooltip>
          <template #default="{ row }">{{ row.thirdOrderId || '-' }}</template>
        </el-table-column>
        <el-table-column label="备注" min-width="160" prop="remark" show-overflow-tooltip>
          <template #default="{ row }">{{ row.remark || '-' }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatUnixTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="支付时间" width="170">
          <template #default="{ row }">{{ formatUnixTime(row.paidAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="{ row }">
            <el-button
                v-if="row.status === 0"
                :loading="manualRechargingId === row.id"
                size="small"
                type="primary"
                @click="handleManualRecharge(row)"
            >
              人工补单
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

    <el-dialog v-model="createOrderDialogVisible" title="人工创建订单" width="420px" @closed="resetCreateOrderForm">
      <el-form ref="createOrderFormRef" :model="createOrderForm" :rules="createOrderRules" label-width="90px">
        <el-form-item label="玩家ID" prop="userId">
          <el-input v-model="createOrderForm.userId" clearable placeholder="请输入玩家ID"/>
        </el-form-item>
        <el-form-item label="订单金额" prop="amount">
          <el-input-number
              v-model="createOrderForm.amount"
              :min="0.01"
              :precision="2"
              :step="1"
              controls-position="right"
              placeholder="请输入订单金额"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createOrderDialogVisible = false">取消</el-button>
        <el-button :loading="creatingOrder" type="primary" @click="handleCreateOrder">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import type {FormInstance, FormRules} from 'element-plus'
import {rechargeOrderApi} from '@/api'
import type {RechargeOrder} from '@/types/api.ts'

interface SearchForm {
  orderId: string
  userId: string
  dateRange: string[]
  statusFilter: number
}

const loading = ref(false)
const pageWaiting = ref(false)
const pageWaitingText = ref('处理中，请稍候...')
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

const createOrderRules: FormRules = {
  userId: [{required: true, message: '请输入玩家ID', trigger: 'blur'}],
  amount: [{required: true, message: '请输入订单金额', trigger: 'change'}],
}

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
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取充值订单列表失败:', error)
    ElMessage.error('获取充值订单列表失败')
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

const afterManualRechargeSuccess = (after: number) => {
  ElMessage.success(`补单成功，玩家当前金币余额：${formatAmount(after)}`)
  pageWaitingText.value = '人工补单处理中，请稍候...'
  pageWaiting.value = true
  setTimeout(() => {
    fetchOrderList().finally(() => {
      pageWaiting.value = false
    })
  }, 5000)
}

const handleManualRecharge = async (row: RechargeOrder) => {
  try {
    await ElMessageBox.confirm(
        `确定对订单 ${row.id} 进行人工补单吗？将发放 ${formatAmount(row.gold)} 金币给玩家 ${row.nickname || row.userId}。`,
        '确认人工补单',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        },
    )
    manualRechargingId.value = row.id
    const res = await rechargeOrderApi.manualRecharge(row.id)
    afterManualRechargeSuccess(res.after)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('人工补单失败:', error)
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
  if (!createOrderFormRef.value) {
    return
  }
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
    ElMessage.success(`创建成功，订单ID：${res.orderId}，金额：${formatAmount(res.price)}，金币：${formatAmount(res.gold)}`)
    currentPage.value = 1
    pageWaitingText.value = '人工创建订单处理中，请稍候...'
    pageWaiting.value = true
    setTimeout(() => {
      fetchOrderList().finally(() => {
        pageWaiting.value = false
      })
    }, 3000)
  } catch (error) {
    console.error('人工创建订单失败:', error)
  } finally {
    creatingOrder.value = false
  }
}

const formatAmount = (val: number | null | undefined) => {
  if (val === null || val === undefined) {
    return '-'
  }
  const n = Number(val)
  if (Number.isNaN(n)) {
    return '-'
  }
  return n.toLocaleString('zh-CN', {maximumFractionDigits: 4})
}

const formatUnixTime = (ts: number | null | undefined) => {
  if (!ts) {
    return '-'
  }
  return new Date(ts * 1000).toLocaleString('zh-CN')
}

const statusLabel = (status: number) => {
  switch (status) {
    case 1:
      return '已完成'
    case 2:
      return '已取消'
    default:
      return '待支付'
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
  switch (source) {
    case 2:
      return '后台手动'
    default:
      return 'App'
  }
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
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>

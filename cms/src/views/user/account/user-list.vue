<template>
  <div
      v-loading="pageWaiting"
      class="page-container"
      element-loading-background="rgba(255, 255, 255, 0.75)"
      :element-loading-text="pageWaitingText"
  >
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
        </div>
      </template>
      <div class="search-form">
        <el-form :model="searchForm" inline label-width="80px">
          <el-form-item label="关键字">
            <el-input v-model="searchForm.key" clearable placeholder="请输入关键字"/>
          </el-form-item>
          <el-form-item label="开始时间">
            <el-date-picker
                v-model="searchForm.startTime"
                clearable
                format="YYYY-MM-DD"
                placeholder="选择开始时间"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item label="结束时间">
            <el-date-picker
                v-model="searchForm.endTime"
                clearable
                format="YYYY-MM-DD"
                placeholder="选择结束时间"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="content">
        <el-table v-loading="loading" :data="userList" height="500" style="width: 100%">
          <el-table-column label="ID" prop="id" width="200"/>
          <el-table-column label="昵称" prop="nickname" width="140">
            <template #default="scope">{{ scope.row.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column label="头像" prop="avatar" width="80">
            <template #default="scope">
              <el-image
                  v-if="scope.row.avatar"
                  :preview-src-list="[scope.row.avatar]"
                  :src="scope.row.avatar"
                  fit="cover"
                  hide-on-click-modal
                  preview-teleported
                  style="width:40px;height:40px;border-radius:50%"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="OpenId" prop="openId" width="200"/>
          <el-table-column label="手机号" prop="phone" width="140">
            <template #default="scope">{{ scope.row.phone || '-' }}</template>
          </el-table-column>
          <el-table-column label="IP" prop="ip" width="150"/>
          <el-table-column label="设备类型" prop="deviceType" width="100">
            <template #default="scope">{{ scope.row.deviceType || '-' }}</template>
          </el-table-column>
          <el-table-column label="渠道" prop="channel" width="100"/>
          <el-table-column label="金币余额" prop="gold" width="120">
            <template #default="scope">{{ formatAmount(scope.row.gold) }}</template>
          </el-table-column>
          <el-table-column label="钻石余额" prop="diamond" width="120">
            <template #default="scope">{{ formatAmount(scope.row.diamond) }}</template>
          </el-table-column>
          <el-table-column label="VIP等级" prop="vipLevel" width="100">
            <template #default="scope">{{ formatVipLevel(scope.row.vipLevel) }}</template>
          </el-table-column>
          <el-table-column label="主播" prop="isAnchor" width="90">
            <template #default="scope">
              <el-tag v-if="scope.row.isAnchor" type="success">是</el-tag>
              <el-tag v-else type="info">否</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="分享码" prop="shareCode" width="140">
            <template #default="scope">{{ scope.row.shareCode || '-' }}</template>
          </el-table-column>
          <el-table-column label="备注" prop="remark" min-width="160" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column label="封号状态" prop="ban" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.ban" type="danger">已封号</el-tag>
              <el-tag v-else type="success">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="注销状态" prop="cancel" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.cancel" type="warning">已注销</el-tag>
              <el-tag v-else type="success">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="封禁时间" prop="banApplyTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.banApplyTime) }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" min-width="660">
            <template #default="scope">
              <el-button
                  v-if="!scope.row.isAnchor"
                  size="small"
                  type="primary"
                  @click="handleSetAnchor(scope.row)"
              >
                设为主播
              </el-button>
              <el-button size="small" type="primary" @click="openCurrencyDialog(scope.row, 'gold', 'add')">
                加金币
              </el-button>
              <el-button size="small" type="warning" @click="openCurrencyDialog(scope.row, 'gold', 'sub')">
                减金币
              </el-button>
              <el-button size="small" type="primary" @click="openCurrencyDialog(scope.row, 'diamond', 'add')">
                加钻石
              </el-button>
              <el-button size="small" type="warning" @click="openCurrencyDialog(scope.row, 'diamond', 'sub')">
                减钻石
              </el-button>
              <el-button
                  :type="scope.row.ban ? 'warning' : 'success'"
                  size="small"
                  @click="toggleBanStatus(scope.row)">
                {{ scope.row.ban ? '解封' : '封号' }}
              </el-button>
              <el-button
                  :type="scope.row.cancel ? 'info' : 'danger'"
                  size="small"
                  @click="toggleCancelStatus(scope.row)">
                {{ scope.row.cancel ? '取消注销' : '注销' }}
              </el-button>
            </template>
          </el-table-column>

        </el-table>
        <div class="pagination">
          <el-pagination
              :current-page="pagination.pageIndex"
              :page-size="pagination.pageSize"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handlePageChange"
              @size-change="handleSizeChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog
        v-model="currencyDialogVisible"
        :title="currencyDialogTitle"
        width="440px"
        @closed="resetCurrencyForm"
    >
      <el-form ref="currencyFormRef" :model="currencyForm" :rules="currencyFormRules" label-width="100px">
        <el-form-item label="用户ID">
          <el-input v-model="currencyForm.userId" disabled/>
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="currencyForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="currencyBalanceLabel">
          <el-input v-model="currencyForm.currentBalanceText" disabled/>
        </el-form-item>
        <el-form-item :label="currencyAmountLabel" prop="amount">
          <el-input-number
              v-model="currencyForm.amount"
              :min="0.0001"
              :precision="4"
              :step="0.0001"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="currencyDialogVisible = false">取消</el-button>
        <el-button :loading="currencySubmitting" type="primary" @click="submitCurrencyChange">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {accountApi, diamondApi, goldApi} from '@/api'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import type {CancelReq, UnBanReq, UnCancelReq, UserInfo} from '@/types/api.ts'

// 用户列表数据
const userList = ref<UserInfo[]>([])
const loading = ref(false)
type CurrencyType = 'gold' | 'diamond'
type CurrencyMode = 'add' | 'sub'

const pageWaiting = ref(false)
const currencyDialogVisible = ref(false)
const currencyType = ref<CurrencyType>('gold')
const currencyMode = ref<CurrencyMode>('add')
const currencySubmitting = ref(false)
const currencyFormRef = ref<FormInstance>()

interface CurrencyForm {
  userId: string
  nickname: string
  currentBalanceText: string
  amount: number
}

const currencyForm = reactive<CurrencyForm>({
  userId: '',
  nickname: '',
  currentBalanceText: '',
  amount: 1
})

const currencyName = computed(() => (currencyType.value === 'gold' ? '金币' : '钻石'))

const currencyDialogTitle = computed(() =>
    `${currencyMode.value === 'add' ? '增加' : '扣减'}${currencyName.value}`
)

const currencyBalanceLabel = computed(() => `当前${currencyName.value}`)

const currencyAmountLabel = computed(() =>
    currencyMode.value === 'add' ? '增加数量' : '扣减数量'
)

const pageWaitingText = computed(() => `${currencyName.value}变更处理中，请稍候...`)

const currencyFormRules: FormRules = {
  amount: [
    {required: true, message: '请输入数量', trigger: 'blur'},
    {
      validator: (_rule, value, callback) => {
        const n = Number(value)
        if (!n || n <= 0) {
          callback(new Error('数量必须大于0'))
          return
        }
        callback()
      },
      trigger: 'blur'
    }
  ]
}

// 搜索表单
const searchForm = reactive({
  key: '',
  startTime: '',
  endTime: ''
})

// 路由
const router = useRouter()
const route = useRoute()

// 分页信息
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0
})

// 监听路由参数变化，当有refresh参数时重新获取数据
watch(() => route.query.refresh, (newRefresh) => {
  if (newRefresh) {
    fetchUserList()
  }
}, {immediate: false})

// 获取用户列表
const fetchUserList = async () => {
  loading.value = true
  try {
    // 构建查询参数
    const params = {
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
      startTime: searchForm.startTime,
      endTime: searchForm.endTime
    }

    const response = await accountApi.getUserInfo(params)

    const result = response.data
    userList.value = result || []
    pagination.total = response?.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 处理查询
const handleSearch = () => {
  // 重置到第一页并查询
  pagination.pageIndex = 1
  fetchUserList()
}

// 处理重置
const handleReset = () => {
  // 重置搜索表单
  searchForm.key = ''
  searchForm.startTime = ''
  searchForm.endTime = ''

  // 重置到第一页并查询
  pagination.pageIndex = 1
  fetchUserList()
}

// 处理分页变化
const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchUserList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchUserList()
}

const resetCurrencyForm = () => {
  currencyForm.userId = ''
  currencyForm.nickname = ''
  currencyForm.currentBalanceText = ''
  currencyForm.amount = 1
  currencyFormRef.value?.clearValidate()
}

const openCurrencyDialog = (row: UserInfo, type: CurrencyType, mode: CurrencyMode) => {
  currencyType.value = type
  currencyMode.value = mode
  currencyForm.userId = String(row.id)
  currencyForm.nickname = row.nickname || '-'
  currencyForm.currentBalanceText = formatAmount(type === 'gold' ? row.gold : row.diamond)
  currencyForm.amount = 1
  currencyDialogVisible.value = true
}

const afterCurrencyChangeSuccess = () => {
  currencyDialogVisible.value = false
  ElMessage.success('操作成功')
  pageWaiting.value = true
  setTimeout(() => {
    fetchUserList().finally(() => {
      pageWaiting.value = false
    })
  }, 1000*16)
}

const submitCurrencyChange = async () => {
  if (!currencyFormRef.value) return
  await currencyFormRef.value.validate(async (valid) => {
    if (!valid) return
    currencySubmitting.value = true
    try {
      const payload = {
        userId: currencyForm.userId,
        amount: currencyForm.amount
      }
      const api = currencyType.value === 'gold' ? goldApi : diamondApi
      if (currencyMode.value === 'add') {
        await api.add(payload)
      } else {
        await api.sub(payload)
      }
      afterCurrencyChangeSuccess()
    } catch (error) {
      console.error(`${currencyName.value}变更失败:`, error)
    } finally {
      currencySubmitting.value = false
    }
  })
}

const handleSetAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        `确定将用户 ${row.id} 设为主播吗？设为主播后不可撤销。`,
        '设为主播',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    await accountApi.setAnchor({accountId: row.id})
    ElMessage.success('已设为主播')
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('设为主播失败:', error)
    }
  }
}

// 切换封号状态
const toggleBanStatus = async (row: UserInfo) => {
  if (row.ban) {
    // 解封操作
    try {
      const result = await ElMessageBox.confirm(
          `确定要解封用户 ${row.id} 吗？`,
          '确认解封',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const unBanData: UnBanReq = {accountId: row.id}
      const response = await accountApi.unBan(unBanData)

      if (response) {
        ElMessage.success('解封成功')
        // 重新加载数据以确保显示最新状态
        setTimeout(() => {
          fetchUserList()
        }, 5000) // 添加短暂延迟以确保后端状态已更新
      } else {
        ElMessage.error('解封失败')
      }
    } catch (error) {
      console.log('取消解封操作')
    }
  } else {
    // 跳转到封号界面
    router.push({
      path: '/user/account/ban-user',
      query: {
        id: row.id,
        openId: row.openId,
        ip: row.ip,
        channel: String(row.channel),
        ban: String(row.ban),
        cancel: String(row.cancel)
      }
    });
  }
}

// 切换注销状态
const toggleCancelStatus = async (row: UserInfo) => {
  if (row.cancel) {
    // 取消注销操作
    try {
      const result = await ElMessageBox.confirm(
          `确定要取消注销用户 ${row.id} 吗？`,
          '确认取消注销',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const unCancelData: UnCancelReq = {accountId: row.id}
      const response = await accountApi.unCancel(unCancelData)

      if (response) {
        ElMessage.success('取消注销成功')
        // 重新加载数据以确保显示最新状态
        setTimeout(() => {
          fetchUserList()
        }, 500) // 添加短暂延迟以确保后端状态已更新
      } else {
        ElMessage.error('取消注销失败')
      }
    } catch (error) {
      console.log('取消取消注销操作')
    }
  } else {
    // 注销操作
    try {
      const result = await ElMessageBox.confirm(
          `确定要注销用户 ${row.id} 吗？注销后用户将无法登录系统`,
          '确认注销',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const cancelData: CancelReq = {accountId: row.id}
      const response = await accountApi.cancel(cancelData)

      if (response) {
        ElMessage.success('注销成功')
        // 重新加载数据以确保显示最新状态
        setTimeout(() => {
          fetchUserList()
        }, 500) // 添加短暂延迟以确保后端状态已更新
      } else {
        ElMessage.error('注销失败')
      }
    } catch (error) {
      console.log('取消注销操作')
    }
  }
}

// 格式化金额(金币/钻石)，保留4位小数
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

const formatVipLevel = (val: number | null | undefined) => {
  if (val === null || val === undefined || val <= 0) {
    return '-'
  }
  return String(val)
}

// 格式化日期函数
const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString('zh-CN')
  } catch {
    return '-'
  }
}

// 页面初始化时获取数据
onMounted(() => {
  fetchUserList()
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

.search-form {
  margin-bottom: 20px;
}

.content {
  min-height: 500px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

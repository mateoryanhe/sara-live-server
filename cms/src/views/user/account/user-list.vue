<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
        </div>
      </template>
      <div class="search-form">
        <el-form :model="searchForm" inline label-width="100px">
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
          <el-form-item label="充值白名单">
            <el-select v-model="searchForm.rechargeWhitelist" clearable placeholder="全部" style="width: 120px">
              <el-option :value="1" label="是"/>
              <el-option :value="0" label="否"/>
            </el-select>
          </el-form-item>
          <el-form-item label="是否主播">
            <el-select v-model="searchForm.isAnchor" clearable placeholder="全部" style="width: 120px">
              <el-option :value="1" label="是"/>
              <el-option :value="0" label="否"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="handleSearch">查询</el-button>
            <el-button v-if="can('search')" @click="handleReset">重置</el-button>
            <el-button v-if="can('batchSetAnchor')" @click="openBatchAnchorDialog('anchor')">批量设普通主播</el-button>
            <el-button v-if="can('batchSetSeniorAnchor')" @click="openBatchAnchorDialog('senior')">批量设高级主播</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="content">
        <el-table v-loading="loading" :data="userList" style="width: 100%">
          <el-table-column label="ID" prop="id" width="200"/>
          <el-table-column label="创建时间" prop="createdAt" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="最后登录时间" prop="lastLoginTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.lastLoginTime) }}
            </template>
          </el-table-column>
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
              <span v-else class="avatar-empty">-</span>
            </template>
          </el-table-column>
          <el-table-column label="区号" prop="phoneAreaCode" width="90">
            <template #default="scope">{{ formatPhoneAreaCode(scope.row.phoneAreaCode) }}</template>
          </el-table-column>
           <el-table-column label="手机号" prop="phone" width="140">
            <template #default="scope">{{ scope.row.phone || '-' }}</template>
          </el-table-column>
          <el-table-column label="金币余额" prop="gold" width="120">
            <template #default="scope">{{ formatAmount(scope.row.gold) }}</template>
          </el-table-column>
          <el-table-column label="钻石余额" prop="diamond" width="120">
            <template #default="scope">{{ formatAmount(scope.row.diamond) }}</template>
          </el-table-column>
          <el-table-column label="VIP等级" prop="vipLevel" width="100">
            <template #default="scope">{{ formatVipLevel(scope.row.vipLevel) }}</template>
          </el-table-column>
          <el-table-column label="用户类型" prop="userType" width="120">
            <template #default="scope">{{ formatUserType(scope.row.userType) }}</template>
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
         
         
          <el-table-column label="IP" prop="ip" width="150"/>
          <el-table-column label="登录国家" prop="loginCountry" width="120">
            <template #default="scope">{{ scope.row.loginCountry || '-' }}</template>
          </el-table-column>
          <el-table-column label="注册IP" prop="registerIp" width="150">
            <template #default="scope">{{ scope.row.registerIp || '-' }}</template>
          </el-table-column>
          <el-table-column label="注册国家" prop="registerCountry" width="120">
            <template #default="scope">{{ scope.row.registerCountry || '-' }}</template>
          </el-table-column>
          <el-table-column label="设备类型" prop="deviceType" width="100">
            <template #default="scope">{{ scope.row.deviceType || '-' }}</template>
          </el-table-column>
           <el-table-column label="包名" prop="packageName" width="180" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.packageName || '-' }}</template>
          </el-table-column>
          <el-table-column label="版本号" prop="appVersion" width="120" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.appVersion || '-' }}</template>
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
          <el-table-column label="充值白名单" prop="rechargeWhitelist" width="110">
            <template #default="scope">
              <el-tag v-if="scope.row.rechargeWhitelist" type="success">是</el-tag>
              <el-tag v-else type="info">否</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="封禁时间" prop="banApplyTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.banApplyTime) }}
            </template>
          </el-table-column>
           <el-table-column label="OpenId" prop="openId" width="200"/>
          <el-table-column label="渠道" prop="channel" width="100"/>
          <el-table-column label="备注" prop="remark" min-width="160" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="100">
            <template #default="scope">
              <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(scope.row, cmd)">
                <el-button size="small" type="primary">
                  操作
                  <el-icon class="el-icon--right">
                    <ArrowDown/>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="can('setAnchor') && !scope.row.isAnchor" command="setAnchor">设为主播</el-dropdown-item>
                    <el-dropdown-item v-if="can('setSeniorAnchor') && !scope.row.isAnchor" command="setSeniorAnchor">设为高级主播</el-dropdown-item>
                    <el-dropdown-item v-if="can('goldAdd')" :divided="!scope.row.isAnchor" command="gold-add">
                      加金币
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('goldSub')" command="gold-sub">减金币</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondAdd')" divided command="diamond-add">加钻石</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondSub')" command="diamond-sub">减钻石</el-dropdown-item>
                    <el-dropdown-item v-if="can('ban')" divided command="ban">
                      {{ scope.row.ban ? '解封' : '封号' }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOff') && scope.row.canRank !== false" command="rank-off">下榜</el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOn') && scope.row.canRank === false" command="rank-on">上榜</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOn') && !scope.row.rechargeWhitelist" command="recharge-whitelist-on">加入充值白名单</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOff') && scope.row.rechargeWhitelist" command="recharge-whitelist-off">移出充值白名单</el-dropdown-item>
                    <el-dropdown-item v-if="can('cancel')" divided command="cancel">
                      {{ scope.row.cancel ? '取消注销' : '注销' }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('setUserType') && !scope.row.isAnchor" command="setUserType">修改用户类型</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="currencyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCurrencyChange">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="banDialogVisible"
        title="封号"
        width="480px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banFormRules" label-width="110px">
        <el-form-item label="用户ID">
          <el-input v-model="banForm.userId" disabled/>
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item label="封号截止时间" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledBanDate"
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择封号截止时间"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">取消</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="userTypeDialogVisible"
        title="修改用户类型"
        width="440px"
        @closed="resetUserTypeForm"
    >
      <el-form ref="userTypeFormRef" :model="userTypeForm" :rules="userTypeFormRules" label-width="100px">
        <el-form-item label="用户ID">
          <el-input v-model="userTypeForm.userId" disabled/>
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="userTypeForm.nickname" disabled/>
        </el-form-item>
        <el-form-item label="用户类型" prop="userType">
          <el-select v-model="userTypeForm.userType" placeholder="请选择用户类型" style="width: 100%">
            <el-option
                v-for="item in userTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userTypeDialogVisible = false">取消</el-button>
        <el-button :loading="userTypeSubmitting" type="primary" @click="submitUserType">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="batchAnchorDialogVisible"
        :title="batchAnchorDialogTitle"
        width="520px"
        @closed="resetBatchAnchorForm"
    >
      <el-form ref="batchAnchorFormRef" :model="batchAnchorForm" :rules="batchAnchorFormRules" label-width="100px">
        <el-form-item label="用户ID" prop="userIdsText">
          <el-input
              v-model="batchAnchorForm.userIdsText"
              :rows="8"
              placeholder="每行一个用户ID，也支持逗号、空格分隔"
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchAnchorDialogVisible = false">取消</el-button>
        <el-button :loading="batchAnchorSubmitting" type="primary" @click="submitBatchAnchor">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {accountApi, diamondApi, goldApi} from '@/api'
import {ArrowDown} from '@element-plus/icons-vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import type {BanReq, CancelReq, UnBanReq, UnCancelReq, UserInfo} from '@/types/api.ts'
import {formatAmount, NUMBER_INPUT_DECIMALS} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {can} = usePagePermission('UserList')

const ROW_ACTION_KEYS = [
  'setAnchor',
  'setSeniorAnchor',
  'goldAdd',
  'goldSub',
  'diamondAdd',
  'diamondSub',
  'ban',
  'rankOff',
  'rankOn',
  'rechargeWhitelistOn',
  'rechargeWhitelistOff',
  'cancel',
  'setUserType',
] as const

const hasRowActions = computed(() => ROW_ACTION_KEYS.some(key => can(key)))

// 用户列表数据
const userList = ref<UserInfo[]>([])
const loading = ref(false)
type CurrencyType = 'gold' | 'diamond'
type CurrencyMode = 'add' | 'sub'

const currencyDialogVisible = ref(false)
const currencyType = ref<CurrencyType>('gold')
const currencyMode = ref<CurrencyMode>('add')
const currencyFormRef = ref<FormInstance>()
const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<FormInstance>()
const userTypeDialogVisible = ref(false)
const userTypeSubmitting = ref(false)
const userTypeFormRef = ref<FormInstance>()
const batchAnchorDialogVisible = ref(false)
const batchAnchorSubmitting = ref(false)
const batchAnchorFormRef = ref<FormInstance>()
const batchAnchorMode = ref<'anchor' | 'senior'>('anchor')

const batchAnchorDialogTitle = computed(() =>
    batchAnchorMode.value === 'senior' ? '批量设高级主播' : '批量设普通主播'
)

interface BatchAnchorForm {
  userIdsText: string
}

const batchAnchorForm = reactive<BatchAnchorForm>({
  userIdsText: ''
})

const parseUserIds = (text: string): string[] => {
  const ids = new Set<string>()
  for (const line of text.split('\n')) {
    for (const part of line.split(/[\s,，;；]+/)) {
      const id = part.trim()
      if (id && /^\d+$/.test(id)) {
        ids.add(id)
      }
    }
  }
  return [...ids]
}

const batchAnchorFormRules: FormRules = {
  userIdsText: [
    {
      validator: (_rule, value, callback) => {
        if (parseUserIds(String(value || '')).length === 0) {
          callback(new Error('请填写至少一个有效的用户ID'))
          return
        }
        callback()
      },
      trigger: 'blur'
    }
  ]
}

const userTypeLabelMap: Record<number, string> = {
  0: '普通用户',
  1: '普通主播',
  2: '机器人主播',
  3: '机器人观众',
  4: '测试人员',
  5: 'CMS短视频作者',
  7: '高级主播',
}

const userTypeOptions = [
  {label: '普通用户', value: 0},
  {label: '测试人员', value: 4}
]

interface UserTypeForm {
  userId: string
  nickname: string
  userType: number
}

const userTypeForm = reactive<UserTypeForm>({
  userId: '',
  nickname: '',
  userType: 0
})

const userTypeFormRules: FormRules = {
  userType: [
    {required: true, message: '请选择用户类型', trigger: 'change'}
  ]
}

interface BanForm {
  userId: string
  nickname: string
  openId: string
  channel: number
  banApplyTime: string
}

const banForm = reactive<BanForm>({
  userId: '',
  nickname: '',
  openId: '',
  channel: 0,
  banApplyTime: ''
})

const banFormRules: FormRules = {
  banApplyTime: [
    {required: true, message: '请选择封号截止时间', trigger: 'change'}
  ]
}

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
  endTime: '',
  rechargeWhitelist: undefined as number | undefined,
  isAnchor: undefined as number | undefined,
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
const fetchUserList = async (silent = false) => {
  if (!silent) {
    loading.value = true
  }
  try {
    // 构建查询参数
    const params: Record<string, unknown> = {
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
      startTime: searchForm.startTime,
      endTime: searchForm.endTime,
    }
    if (searchForm.rechargeWhitelist !== undefined && searchForm.rechargeWhitelist !== null) {
      params.rechargeWhitelist = searchForm.rechargeWhitelist
    }
    if (searchForm.isAnchor !== undefined && searchForm.isAnchor !== null) {
      params.isAnchor = searchForm.isAnchor
    }

    const response = await accountApi.getUserInfo(params as Parameters<typeof accountApi.getUserInfo>[0])

    const result = response.data
    userList.value = result || []
    pagination.total = response?.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    if (!silent) {
      loading.value = false
    }
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
  searchForm.rechargeWhitelist = undefined
  searchForm.isAnchor = undefined

  // 重置到第一页并查询
  pagination.pageIndex = 1
  fetchUserList()
}

// 处理分页变化
const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchUserList()
}

const resolveAccountOpenId = (row: UserInfo): string => {
  let openId = row.openId || ''
  const canceledIdx = openId.indexOf('__canceled__')
  if (canceledIdx > 0) {
    openId = openId.slice(0, canceledIdx)
  }
  return openId
}

const accountActionPayload = (row: UserInfo) => ({
  accountId: row.id,
  openId: resolveAccountOpenId(row),
  channel: row.channel ?? 0,
})

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

const handleRowCommand = (row: UserInfo, command: string) => {
  switch (command) {
    case 'setUserType':
      openUserTypeDialog(row)
      break
    case 'gold-add':
      openCurrencyDialog(row, 'gold', 'add')
      break
    case 'gold-sub':
      openCurrencyDialog(row, 'gold', 'sub')
      break
    case 'diamond-add':
      openCurrencyDialog(row, 'diamond', 'add')
      break
    case 'diamond-sub':
      openCurrencyDialog(row, 'diamond', 'sub')
      break
    case 'ban':
      handleBanAction(row)
      break
    case 'rank-on':
      handleCanRankAction(row, true)
      break
    case 'rank-off':
      handleCanRankAction(row, false)
      break
    case 'recharge-whitelist-on':
      handleRechargeWhitelistAction(row, true)
      break
    case 'recharge-whitelist-off':
      handleRechargeWhitelistAction(row, false)
      break
    case 'cancel':
      toggleCancelStatus(row)
      break
    case 'setAnchor':
      handleSetAnchor(row)
      break
    case 'setSeniorAnchor':
      handleSetSeniorAnchor(row)
      break
  }
}

const resetUserTypeForm = () => {
  userTypeForm.userId = ''
  userTypeForm.nickname = ''
  userTypeForm.userType = 0
  userTypeFormRef.value?.clearValidate()
}

const openUserTypeDialog = (row: UserInfo) => {
  userTypeForm.userId = String(row.id)
  userTypeForm.nickname = row.nickname || '-'
  const currentType = row.userType ?? 0
  userTypeForm.userType = currentType === 4 ? 4 : 0
  userTypeDialogVisible.value = true
}

const submitUserType = async () => {
  if (!userTypeFormRef.value) return
  await userTypeFormRef.value.validate(async (valid) => {
    if (!valid) return
    userTypeSubmitting.value = true
    try {
      await accountApi.setUserType({
        accountId: userTypeForm.userId,
        userType: userTypeForm.userType
      })
      userTypeDialogVisible.value = false
      ElMessage.success('用户类型已更新')
      await fetchUserList()
    } catch (error) {
      console.error('修改用户类型失败:', error)
    } finally {
      userTypeSubmitting.value = false
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
    await accountApi.setAnchor({accountId: String(row.id)})
    ElMessage.success('已设为主播')
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('设为主播失败:', error)
    }
  }
}

const handleSetSeniorAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        `确定将用户 ${row.id} 设为高级主播吗？设为高级主播后不可撤销。`,
        '设为高级主播',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    await accountApi.setSeniorAnchor({accountId: String(row.id)})
    ElMessage.success('已设为高级主播')
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('设为高级主播失败:', error)
    }
  }
}

const resetBatchAnchorForm = () => {
  batchAnchorForm.userIdsText = ''
  batchAnchorFormRef.value?.clearValidate()
}

const openBatchAnchorDialog = (mode: 'anchor' | 'senior') => {
  batchAnchorMode.value = mode
  batchAnchorDialogVisible.value = true
}

const submitBatchAnchor = async () => {
  if (!batchAnchorFormRef.value) return
  await batchAnchorFormRef.value.validate(async (valid) => {
    if (!valid) return
    const ids = parseUserIds(batchAnchorForm.userIdsText)
    const actionLabel = batchAnchorMode.value === 'senior' ? '高级主播' : '普通主播'
    try {
      await ElMessageBox.confirm(
          `确定将 ${ids.length} 个用户设为${actionLabel}吗？设置后不可撤销。`,
          `批量设${actionLabel}`,
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )
      batchAnchorSubmitting.value = true
      const response = batchAnchorMode.value === 'senior'
          ? await accountApi.batchSetSeniorAnchor({ids})
          : await accountApi.batchSetAnchor({ids})
      batchAnchorDialogVisible.value = false
      if (response.failCount > 0) {
        ElMessage.warning(`批量设置完成：成功 ${response.successCount} 个，失败 ${response.failCount} 个`)
      } else {
        ElMessage.success(`批量设置成功，共 ${response.successCount} 个`)
      }
      await fetchUserList()
    } catch (error) {
      if (error !== 'cancel') {
        console.error('批量设主播失败:', error)
      }
    } finally {
      batchAnchorSubmitting.value = false
    }
  })
}

const afterCurrencyChangeSuccess = () => {
  currencyDialogVisible.value = false
  ElMessage.success('操作成功')
  fetchUserList(true)
}

const submitCurrencyChange = async () => {
  if (!currencyFormRef.value) return
  await currencyFormRef.value.validate(async (valid) => {
    if (!valid) return
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
    }
  })
}

const defaultBanApplyTime = () => {
  const d = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const disabledBanDate = (time: Date) => time.getTime() < Date.now()

const resetBanForm = () => {
  banForm.userId = ''
  banForm.nickname = ''
  banForm.banApplyTime = ''
  banFormRef.value?.clearValidate()
}

const openBanDialog = (row: UserInfo) => {
  banForm.userId = String(row.id)
  banForm.nickname = row.nickname || '-'
  banForm.openId = resolveAccountOpenId(row)
  banForm.channel = row.channel ?? 0
  banForm.banApplyTime = defaultBanApplyTime()
  banDialogVisible.value = true
}

const submitBan = async () => {
  if (!banFormRef.value) return
  await banFormRef.value.validate(async (valid) => {
    if (!valid) return
    banSubmitting.value = true
    try {
      const banData: BanReq = {
        accountId: banForm.userId,
        openId: banForm.openId,
        channel: banForm.channel,
        banApplyTime: banForm.banApplyTime
      }
      const response = await accountApi.ban(banData)
      if (response) {
        banDialogVisible.value = false
        ElMessage.success('封号成功')
        await fetchUserList()
      } else {
        ElMessage.error('封号失败')
      }
    } catch (error) {
      console.error('封号失败:', error)
    } finally {
      banSubmitting.value = false
    }
  })
}

const handleCanRankAction = async (row: UserInfo, canRank: boolean) => {
  const actionText = canRank ? '上榜' : '下榜'
  try {
    await ElMessageBox.confirm(
        `确定要将用户 ${row.id} ${actionText}吗？`,
        `确认${actionText}`,
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    const response = await accountApi.setCanRank({
      accountId: row.id,
      canRank
    })
    if (response) {
      ElMessage.success(`${actionText}成功`)
      await fetchUserList()
    } else {
      ElMessage.error(`${actionText}失败`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(`${actionText}失败:`, error)
    }
  }
}

const handleRechargeWhitelistAction = async (row: UserInfo, rechargeWhitelist: boolean) => {
  const actionText = rechargeWhitelist ? '加入充值白名单' : '移出充值白名单'
  try {
    await ElMessageBox.confirm(
        `确定要将用户 ${row.id} ${actionText}吗？${rechargeWhitelist ? '加入后 App 创建充值订单将直接到账。' : ''}`,
        `确认${actionText}`,
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    const response = await accountApi.setRechargeWhitelist({
      accountId: row.id,
      rechargeWhitelist
    })
    if (response) {
      ElMessage.success(`${actionText}成功`)
      await fetchUserList()
    } else {
      ElMessage.error(`${actionText}失败`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(`${actionText}失败:`, error)
    }
  }
}

const handleBanAction = async (row: UserInfo) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
          `确定要解封用户 ${row.id} 吗？`,
          '确认解封',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const response = await accountApi.unBan(accountActionPayload(row))

      if (response) {
        ElMessage.success('解封成功')
        await fetchUserList()
      } else {
        ElMessage.error('解封失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('解封失败:', error)
      }
    }
    return
  }
  openBanDialog(row)
}

const patchUserRow = (userId: number | string, patch: Partial<UserInfo>) => {
  const id = String(userId)
  const idx = userList.value.findIndex(item => String(item.id) === id)
  if (idx >= 0) {
    userList.value[idx] = {...userList.value[idx], ...patch}
  }
}

// 切换注销状态
const toggleCancelStatus = async (row: UserInfo) => {
  if (row.cancel) {
    // 取消注销操作
    try {
      await ElMessageBox.confirm(
          `确定要取消注销用户 ${row.id} 吗？`,
          '确认取消注销',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const unCancelData: UnCancelReq = accountActionPayload(row)
      const response = await accountApi.unCancel(unCancelData)

      if (response) {
        patchUserRow(row.id, {cancel: false})
        ElMessage.success('取消注销成功')
        await fetchUserList(true)
      } else {
        ElMessage.error('取消注销失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('取消注销失败:', error)
      }
    }
  } else {
    // 注销操作
    try {
      await ElMessageBox.confirm(
          `确定要注销用户 ${row.id} 吗？注销后用户将无法登录系统`,
          '确认注销',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
      )

      const cancelData: CancelReq = accountActionPayload(row)
      const response = await accountApi.cancel(cancelData)

      if (response) {
        patchUserRow(row.id, {cancel: true})
        ElMessage.success('注销成功')
        await fetchUserList(true)
      } else {
        ElMessage.error('注销失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('注销失败:', error)
      }
    }
  }
}

const formatPhoneAreaCode = (val?: string) => {
  const code = (val || '').trim()
  if (!code) {
    return '-'
  }
  return code.startsWith('+') ? code : `+${code}`
}

const formatUserType = (val: number | null | undefined) => {
  if (val === null || val === undefined) {
    return '普通用户'
  }
  return userTypeLabelMap[val] || '普通用户'
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

.search-form :deep(.el-form-item__label) {
  white-space: nowrap;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

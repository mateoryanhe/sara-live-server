<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.UserList') }}</span>
        </div>
      </template>
      <div class="search-form">
        <el-form :model="searchForm" inline label-width="100px">
          <el-form-item :label="t('common.keyword')">
            <el-input v-model="searchForm.key" clearable :placeholder="t('pages.userList.keywordPlaceholder')"/>
          </el-form-item>
          <el-form-item :label="t('common.startTime')">
            <el-date-picker
                v-model="searchForm.startTime"
                clearable
                format="YYYY-MM-DD"
                :placeholder="t('common.selectStartTime')"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item :label="t('common.endTime')">
            <el-date-picker
                v-model="searchForm.endTime"
                clearable
                format="YYYY-MM-DD"
                :placeholder="t('common.selectEndTime')"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item :label="t('pages.userList.rechargeWhitelist')">
            <el-select v-model="searchForm.rechargeWhitelist" clearable :placeholder="t('common.all')" style="width: 120px">
              <el-option :value="1" :label="t('common.yes')"/>
              <el-option :value="0" :label="t('common.no')"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('pages.userList.isAnchor')">
            <el-select v-model="searchForm.isAnchor" clearable :placeholder="t('common.all')" style="width: 120px">
              <el-option :value="1" :label="t('common.yes')"/>
              <el-option :value="0" :label="t('common.no')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
            <el-button v-if="can('search')" @click="handleReset">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="content">
        <el-table v-loading="loading" :data="userList" style="width: 100%">
          <el-table-column label="ID" prop="id" width="200"/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.lastLoginTime')" prop="lastLoginTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.lastLoginTime) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.nickname')" prop="nickname" width="140">
            <template #default="scope">{{ scope.row.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.avatar')" prop="avatar" width="80">
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
          <el-table-column :label="t('pages.userList.phoneAreaCode')" prop="phoneAreaCode" width="90">
            <template #default="scope">{{ formatPhoneAreaCode(scope.row.phoneAreaCode) }}</template>
          </el-table-column>
           <el-table-column :label="t('pages.userList.phone')" prop="phone" width="140">
            <template #default="scope">{{ scope.row.phone || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.goldBalance')" prop="gold" width="120">
            <template #default="scope">{{ formatAmount(scope.row.gold) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.diamondBalance')" prop="diamond" width="120">
            <template #default="scope">{{ formatAmount(scope.row.diamond) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.vipLevel')" prop="vipLevel" width="100">
            <template #default="scope">{{ formatVipLevel(scope.row.vipLevel) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.userType')" prop="userType" width="120">
            <template #default="scope">{{ formatUserType(scope.row.userType) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.anchor')" prop="isAnchor" width="90">
            <template #default="scope">
              <el-tag v-if="scope.row.isAnchor" type="success">{{ t('common.yes') }}</el-tag>
              <el-tag v-else type="info">{{ t('common.no') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.shareCode')" prop="shareCode" width="140">
            <template #default="scope">{{ scope.row.shareCode || '-' }}</template>
          </el-table-column>
         
         
          <el-table-column label="IP" prop="ip" width="150"/>
          <el-table-column :label="t('pages.userList.loginCountry')" prop="loginCountry" width="120">
            <template #default="scope">{{ scope.row.loginCountry || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.registerIp')" prop="registerIp" width="150">
            <template #default="scope">{{ scope.row.registerIp || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.registerCountry')" prop="registerCountry" width="120">
            <template #default="scope">{{ scope.row.registerCountry || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.deviceType')" prop="deviceType" width="100">
            <template #default="scope">{{ scope.row.deviceType || '-' }}</template>
          </el-table-column>
           <el-table-column :label="t('pages.userList.packageName')" prop="packageName" width="180" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.packageName || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.appVersion')" prop="appVersion" width="120" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.appVersion || '-' }}</template>
          </el-table-column>
          
          <el-table-column :label="t('pages.userList.banStatus')" prop="ban" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.ban" type="danger">{{ t('pages.userList.banned') }}</el-tag>
              <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.cancelStatus')" prop="cancel" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.cancel" type="warning">{{ t('pages.userList.canceled') }}</el-tag>
              <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.rechargeWhitelist')" prop="rechargeWhitelist" width="110">
            <template #default="scope">
              <el-tag v-if="scope.row.rechargeWhitelist" type="success">{{ t('common.yes') }}</el-tag>
              <el-tag v-else type="info">{{ t('common.no') }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column :label="t('pages.userList.banTime')" prop="banApplyTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.banApplyTime) }}
            </template>
          </el-table-column>
           <el-table-column :label="t('pages.userList.openId')" prop="openId" width="200"/>
          <el-table-column :label="t('pages.userList.channel')" prop="channel" width="100"/>
          <el-table-column :label="t('common.remark')" prop="remark" min-width="160" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column fixed="right" :label="t('common.actions')" width="100">
            <template #default="scope">
              <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(scope.row, cmd)">
                <el-button size="small" type="primary">
                  {{ t('common.actions') }}
                  <el-icon class="el-icon--right">
                    <ArrowDown/>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="can('setAnchor') && !scope.row.isAnchor" command="setAnchor">{{ t('pages.userList.setAnchor') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('setSeniorAnchor') && !scope.row.isAnchor" command="setSeniorAnchor">{{ t('pages.userList.setSeniorAnchor') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('goldAdd')" :divided="!scope.row.isAnchor" command="gold-add">
                      {{ t('pages.userList.addGold') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('goldSub')" command="gold-sub">{{ t('pages.userList.subGold') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondAdd')" divided command="diamond-add">{{ t('pages.userList.addDiamond') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondSub')" command="diamond-sub">{{ t('pages.userList.subDiamond') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('ban')" divided command="ban">
                      {{ scope.row.ban ? t('pages.userList.unban') : t('pages.userList.ban') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOff') && scope.row.canRank !== false" command="rank-off">{{ t('pages.userList.rankOff') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOn') && scope.row.canRank === false" command="rank-on">{{ t('pages.userList.rankOn') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOn') && !scope.row.rechargeWhitelist" command="recharge-whitelist-on">{{ t('pages.userList.rechargeWhitelistOn') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOff') && scope.row.rechargeWhitelist" command="recharge-whitelist-off">{{ t('pages.userList.rechargeWhitelistOff') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('cancel')" divided command="cancel">
                      {{ scope.row.cancel ? t('pages.userList.uncancel') : t('pages.userList.cancelAccount') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('setUserType') && !scope.row.isAnchor" command="setUserType">{{ t('pages.userList.setUserType') }}</el-dropdown-item>
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
        <el-form-item :label="t('common.userId')">
          <el-input v-model="currencyForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
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
        <el-button @click="currencyDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submitCurrencyChange">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="banDialogVisible"
        :title="t('pages.userList.banDialogTitle')"
        width="480px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banFormRules" label-width="110px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="banForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.userList.banUntil')" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledBanDate"
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.userList.selectBanUntil')"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="userTypeDialogVisible"
        :title="t('pages.userList.editUserTypeTitle')"
        width="440px"
        @closed="resetUserTypeForm"
    >
      <el-form ref="userTypeFormRef" :model="userTypeForm" :rules="userTypeFormRules" label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="userTypeForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="userTypeForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.userList.userType')" prop="userType">
          <el-select v-model="userTypeForm.userType" :placeholder="t('pages.userList.selectUserType')" style="width: 100%">
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
        <el-button @click="userTypeDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="userTypeSubmitting" type="primary" @click="submitUserType">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {accountApi, diamondApi, goldApi} from '@/api'
import {ArrowDown} from '@element-plus/icons-vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import type {BanReq, CancelReq, UnBanReq, UnCancelReq, UserInfo} from '@/types/api.ts'
import {formatAmount, NUMBER_INPUT_DECIMALS} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {t, locale} = useI18n()
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

const userTypeLabelMap = computed<Record<number, string>>(() => ({
  0: t('pages.userList.userTypeNormal'),
  1: t('pages.userList.userTypeAnchor'),
  2: t('pages.userList.userTypeBotAnchor'),
  3: t('pages.userList.userTypeBotViewer'),
  4: t('pages.userList.userTypeTester'),
  5: t('pages.userList.userTypeCmsAuthor'),
  7: t('pages.userList.userTypeSeniorAnchor'),
}))

const userTypeOptions = computed(() => [
  {label: t('pages.userList.userTypeNormal'), value: 0},
  {label: t('pages.userList.userTypeTester'), value: 4}
])

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

const userTypeFormRules = computed<FormRules>(() => ({
  userType: [
    {required: true, message: t('pages.userList.selectUserTypeRequired'), trigger: 'change'}
  ]
}))

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

const banFormRules = computed<FormRules>(() => ({
  banApplyTime: [
    {required: true, message: t('pages.userList.selectBanUntilRequired'), trigger: 'change'}
  ]
}))

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

const currencyName = computed(() =>
    currencyType.value === 'gold' ? t('pages.userList.gold') : t('pages.userList.diamond')
)

const currencyDialogTitle = computed(() =>
    currencyMode.value === 'add'
        ? t('pages.userList.addCurrency', {currency: currencyName.value})
        : t('pages.userList.subCurrency', {currency: currencyName.value})
)

const currencyBalanceLabel = computed(() =>
    t('pages.userList.currentBalance', {balance: currencyName.value})
)

const currencyAmountLabel = computed(() =>
    currencyMode.value === 'add' ? t('pages.userList.addAmount') : t('pages.userList.subAmount')
)

const currencyFormRules = computed<FormRules>(() => ({
  amount: [
    {required: true, message: t('pages.userList.amountRequired'), trigger: 'blur'},
    {
      validator: (_rule, value, callback) => {
        const n = Number(value)
        if (!n || n <= 0) {
          callback(new Error(t('pages.userList.amountPositive')))
          return
        }
        callback()
      },
      trigger: 'blur'
    }
  ]
}))

const searchForm = reactive({
  key: '',
  startTime: '',
  endTime: '',
  rechargeWhitelist: undefined as number | undefined,
  isAnchor: undefined as number | undefined,
})

const router = useRouter()
const route = useRoute()

const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0
})

watch(() => route.query.refresh, (newRefresh) => {
  if (newRefresh) {
    fetchUserList()
  }
}, {immediate: false})

const fetchUserList = async (silent = false) => {
  if (!silent) {
    loading.value = true
  }
  try {
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
    console.error('fetchUserList failed:', error)
    ElMessage.error(t('pages.userList.fetchFailed'))
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchUserList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.startTime = ''
  searchForm.endTime = ''
  searchForm.rechargeWhitelist = undefined
  searchForm.isAnchor = undefined

  pagination.pageIndex = 1
  fetchUserList()
}

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
      ElMessage.success(t('pages.userList.userTypeUpdated'))
      await fetchUserList()
    } catch (error) {
      console.error('setUserType failed:', error)
    } finally {
      userTypeSubmitting.value = false
    }
  })
}

const handleSetAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        t('pages.userList.setAnchorConfirm', {id: row.id}),
        t('pages.userList.setAnchorTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await accountApi.setAnchor({accountId: String(row.id)})
    ElMessage.success(t('pages.userList.setAnchorSuccess'))
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('setAnchor failed:', error)
    }
  }
}

const handleSetSeniorAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        t('pages.userList.setSeniorAnchorConfirm', {id: row.id}),
        t('pages.userList.setSeniorAnchorTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await accountApi.setSeniorAnchor({accountId: String(row.id)})
    ElMessage.success(t('pages.userList.setSeniorAnchorSuccess'))
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('setSeniorAnchor failed:', error)
    }
  }
}

const afterCurrencyChangeSuccess = () => {
  currencyDialogVisible.value = false
  ElMessage.success(t('common.operationSuccess'))
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
      console.error('currencyChange failed:', error)
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
        ElMessage.success(t('pages.userList.banSuccess'))
        await fetchUserList()
      } else {
        ElMessage.error(t('pages.userList.banFailed'))
      }
    } catch (error) {
      console.error('ban failed:', error)
    } finally {
      banSubmitting.value = false
    }
  })
}

const handleCanRankAction = async (row: UserInfo, canRank: boolean) => {
  const actionText = canRank ? t('pages.userList.rankOnAction') : t('pages.userList.rankOffAction')
  try {
    await ElMessageBox.confirm(
        t('pages.userList.rankConfirm', {id: row.id, action: actionText}),
        t('pages.userList.rankTitle', {action: actionText}),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    const response = await accountApi.setCanRank({
      accountId: row.id,
      canRank
    })
    if (response) {
      ElMessage.success(t('pages.userList.rankSuccess', {action: actionText}))
      await fetchUserList()
    } else {
      ElMessage.error(t('pages.userList.rankFailed', {action: actionText}))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('canRank action failed:', error)
    }
  }
}

const handleRechargeWhitelistAction = async (row: UserInfo, rechargeWhitelist: boolean) => {
  const actionText = rechargeWhitelist
      ? t('pages.userList.rechargeWhitelistJoin')
      : t('pages.userList.rechargeWhitelistRemove')
  const extra = rechargeWhitelist ? t('pages.userList.rechargeWhitelistExtra') : ''
  try {
    await ElMessageBox.confirm(
        t('pages.userList.rechargeWhitelistConfirm', {id: row.id, action: actionText, extra}),
        t('pages.userList.rankTitle', {action: actionText}),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    const response = await accountApi.setRechargeWhitelist({
      accountId: row.id,
      rechargeWhitelist
    })
    if (response) {
      ElMessage.success(t('pages.userList.rankSuccess', {action: actionText}))
      await fetchUserList()
    } else {
      ElMessage.error(t('pages.userList.rankFailed', {action: actionText}))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('rechargeWhitelist action failed:', error)
    }
  }
}

const handleBanAction = async (row: UserInfo) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.unbanConfirm', {id: row.id}),
          t('pages.userList.unbanTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const response = await accountApi.unBan(accountActionPayload(row))

      if (response) {
        ElMessage.success(t('pages.userList.unbanSuccess'))
        await fetchUserList()
      } else {
        ElMessage.error(t('pages.userList.unbanFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('unban failed:', error)
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

const toggleCancelStatus = async (row: UserInfo) => {
  if (row.cancel) {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.uncancelConfirm', {id: row.id}),
          t('pages.userList.uncancelTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const unCancelData: UnCancelReq = accountActionPayload(row)
      const response = await accountApi.unCancel(unCancelData)

      if (response) {
        patchUserRow(row.id, {cancel: false})
        ElMessage.success(t('pages.userList.uncancelSuccess'))
        await fetchUserList(true)
      } else {
        ElMessage.error(t('pages.userList.uncancelFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('uncancel failed:', error)
      }
    }
  } else {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.cancelConfirm', {id: row.id}),
          t('pages.userList.cancelTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const cancelData: CancelReq = accountActionPayload(row)
      const response = await accountApi.cancel(cancelData)

      if (response) {
        patchUserRow(row.id, {cancel: true})
        ElMessage.success(t('pages.userList.cancelSuccess'))
        await fetchUserList(true)
      } else {
        ElMessage.error(t('pages.userList.cancelFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('cancel failed:', error)
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
    return t('pages.userList.userTypeNormal')
  }
  return userTypeLabelMap.value[val] || t('pages.userList.userTypeNormal')
}

const formatVipLevel = (val: number | null | undefined) => {
  if (val === null || val === undefined || val <= 0) {
    return '-'
  }
  return String(val)
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString(locale.value)
  } catch {
    return '-'
  }
}

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

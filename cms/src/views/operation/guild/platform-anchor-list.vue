<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorList') }}</span>
          <div class="header-actions">
            <el-button
                v-if="can('batchImportSalaryAnchor')"
                type="warning"
                @click="openSalaryImportDialog"
            >
              {{ t('pages.guildList.batchImportSalaryAnchor') }}
            </el-button>
            <el-button
                v-if="can('batchImmediateSettlement')"
                :disabled="selectedAnchorRows.length === 0"
                :loading="immediateSettling"
                type="danger"
                @click="handleBatchImmediateSettlement"
            >
              {{ t('pages.platformAnchorImmediateSettlement.button') }}
            </el-button>
          </div>
        </div>
      </template>

      <div class="search-form">
        <el-form :model="searchForm" class="search-form" inline label-width="100px">
          <el-form-item :label="t('common.keyword')">
            <el-input
                v-model="searchForm.key"
                clearable
                :placeholder="t('pages.anchorList.keywordPlaceholder')"
                style="width: 200px"
            />
          </el-form-item>
          <el-form-item :label="t('pages.anchorList.liveStatus')">
            <el-select v-model="searchForm.liveStatus" style="width: 120px">
              <el-option :value="ALL_LIVE_STATUS" :label="t('common.all')"/>
              <el-option :value="1" :label="t('common.live')"/>
              <el-option :value="0" :label="t('common.offline')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
            <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
          v-loading="loading"
          :data="tableData"
          style="width: 100%"
          @selection-change="handleSelectionChange"
      >
        <el-table-column fixed type="selection" width="48"/>
        <el-table-column :label="t('common.userId')" prop="id" width="180">
          <template #default="{ row }">
            <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
              {{ row.id }}
            </el-button>
            <span v-else>{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.avatar"
                :preview-src-list="[row.avatar]"
                :src="row.avatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? t('common.live') : t('common.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.socialIncomeTotal')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalSocialIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('menu.UserDetail')" width="110">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail" link type="primary" @click="openUserDetail(row.id)">
              {{ t('pages.userList.viewDetail') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.anchorType')" width="110">
          <template #default="{ row }">
            <el-tag :type="anchorTypeTagType(row.userType)">{{ anchorTypeLabel(row.userType) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.loginIp')" min-width="140" prop="ip">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveRoom')" prop="roomId" width="180">
          <template #default="{ row }">{{ row.roomId || row.id || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.roomType')" width="100">
          <template #default="{ row }">
            <el-tag :type="categoryTagType(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.banStatus')" prop="ban" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.ban" type="danger">{{ t('common.banned') }}</el-tag>
            <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.shelfStatus')" prop="status" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">{{ t('common.onShelf') }}</el-tag>
            <el-tag v-else type="info">{{ t('common.offShelf') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.salaryEffectiveStatus')" width="120">
          <template #default="{ row }">
            <el-tooltip
                :content="salaryValidityText(row)"
                placement="top"
            >
              <el-tag :type="row.salaryEffective ? 'success' : 'info'">
                {{ row.salaryEffective
                  ? t('pages.anchorList.salaryEffective')
                  : t('pages.anchorList.salaryInactive') }}
              </el-tag>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="100">
          <template #default="{ row }">
            <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(row, cmd)">
              <el-button size="small" type="primary">
                {{ t('common.actions') }}
                <el-icon class="el-icon--right"><ArrowDown/></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="canViewDetail" command="viewDetail">
                    {{ t('common.detail') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="canSetAnchorType(row)" command="setAnchorType">
                    {{ t('pages.guildMembers.setAnchorType') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="can('transferInfo')" command="transferInfo">
                    {{ t('pages.guildList.transferInfo') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="can('offShelf')" divided command="offShelf">
                    {{ t('common.offShelf') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                      v-if="row.ban ? can('unban') : can('ban')"
                      :command="row.ban ? 'unban' : 'ban'"
                  >
                    {{ row.ban ? t('pages.anchorList.unban') : t('pages.anchorList.ban') }}
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
        v-model="banDialogVisible"
        :close-on-click-modal="false"
        destroy-on-close
        :title="t('pages.anchorList.banDialogTitle')"
        width="520px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banRules" label-width="100px">
        <el-form-item :label="t('pages.anchorList.anchorId')">
          <el-input v-model="banForm.accountId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.banUntil')" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.anchorList.selectBanUntil')"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.banReason')" prop="banReason">
          <el-input
              v-model="banForm.banReason"
              :maxlength="512"
              :rows="4"
              :placeholder="t('pages.anchorList.enterBanReason')"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">{{ t('pages.anchorList.confirmBan') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="anchorTypeDialogVisible"
        :title="anchorTypeDialogTitle"
        width="480px"
        @closed="resetAnchorTypeForm"
    >
      <el-form ref="anchorTypeFormRef" :model="anchorTypeForm" :rules="anchorTypeRules" label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="anchorTypeForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="anchorTypeForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.anchorType')" prop="anchorType">
          <el-select v-model="anchorTypeForm.anchorType" style="width: 100%">
            <el-option :label="t('pages.anchorList.anchorTypeNormal')" :value="1"/>
            <el-option :label="t('pages.anchorList.anchorTypeSenior')" :value="7"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="anchorTypeDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="anchorTypeSubmitting" type="primary" @click="submitAnchorType">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="salaryImportDialogVisible"
        :title="t('pages.guildList.batchImportSalaryAnchorTitle')"
        width="560px"
    >
      <el-alert
          :closable="false"
          :title="t('pages.guildList.batchImportSalaryAnchorValidityTip')"
          class="salary-import-tip"
          show-icon
          type="info"
      />
      <el-form
          ref="salaryImportFormRef"
          :model="salaryImportForm"
          :rules="salaryImportRules"
          label-width="80px"
      >
        <el-form-item :label="t('pages.guildList.importUserIds')" prop="userIdsText">
          <el-input
              v-model="salaryImportForm.userIdsText"
              :autosize="{ minRows: 8, maxRows: 16 }"
              :placeholder="t('pages.guildList.importUserIdsPlaceholder')"
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="salaryImportDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="salaryImporting" type="primary" @click="handleSalaryImportSubmit">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElForm, ElMessage, ElMessageBox, type FormRules} from 'element-plus'
import {ArrowDown} from '@element-plus/icons-vue'
import {accountApi, guildApi} from '@/api'
import type {AnchorListItem, BanAnchorReq, UnBanAnchorReq} from '@/types/api'
import {formatAmount, formatWalletBalance} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {formatServerDateTime, formatServerNowPlusDays} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('PlatformAnchorList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('PlatformAnchorList')
const canViewDetail = computed(() => can('viewDetail'))
const USER_TYPE_ANCHOR = 1
const USER_TYPE_SENIOR_ANCHOR = 7
const hasRowActions = computed(() => canViewDetail.value || [
  'setAnchorType',
  'transferInfo',
  'offShelf',
  'ban',
  'unban',
].some(key => can(key)))
const canSetAnchorType = (row: AnchorListItem) => {
  if (!can('setAnchorType')) {
    return false
  }
  return row.userType === USER_TYPE_ANCHOR || row.userType === USER_TYPE_SENIOR_ANCHOR
}

const loading = ref(false)
const tableData = ref<AnchorListItem[]>([])
const selectedAnchorRows = ref<AnchorListItem[]>([])
const immediateSettling = ref(false)
const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<InstanceType<typeof ElForm>>()
const anchorTypeDialogVisible = ref(false)
const anchorTypeDialogTitle = ref('')
const anchorTypeSubmitting = ref(false)
const anchorTypeFormRef = ref<InstanceType<typeof ElForm>>()
const anchorTypeForm = reactive({
  userId: '',
  nickname: '',
  anchorType: 1 as 1 | 7,
})
const salaryImportDialogVisible = ref(false)
const salaryImporting = ref(false)
const salaryImportFormRef = ref<InstanceType<typeof ElForm>>()
const salaryImportForm = reactive({userIdsText: ''})

const ALL_LIVE_STATUS = -1
const searchForm = reactive({key: '', liveStatus: 1})
const pagination = reactive({pageIndex: 1, pageSize: 10, total: 0})
const banForm = reactive({accountId: '', nickname: '', banApplyTime: '', banReason: ''})

const banRules = computed<FormRules>(() => ({
  banApplyTime: [{required: true, message: t('pages.anchorList.banApplyTimeRequired'), trigger: 'change'}],
  banReason: [
    {required: true, message: t('pages.anchorList.banReasonRequired'), trigger: 'blur'},
    {min: 1, max: 512, message: t('pages.anchorList.banReasonLength'), trigger: 'blur'},
  ],
}))

const anchorTypeRules = computed<FormRules>(() => ({
  anchorType: [
    {required: true, message: t('pages.guildMembers.anchorTypeRequired'), trigger: 'change'},
  ],
}))

const salaryImportRules = computed<FormRules>(() => ({
  userIdsText: [
    {required: true, message: t('pages.guildList.importUserIdsRequired'), trigger: 'blur'},
  ],
}))

const anchorTypeLabel = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return t('pages.anchorList.anchorTypeSenior')
  if (userType === USER_TYPE_ANCHOR) return t('pages.anchorList.anchorTypeNormal')
  return '-'
}

const anchorTypeTagType = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return 'warning'
  if (userType === USER_TYPE_ANCHOR) return 'success'
  return 'info'
}

const categoryLabel = (category?: number) => {
  if (category === 1) return t('pages.anchorList.categoryHot')
  if (category === 2) return t('pages.anchorList.categoryGame')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === 2) return 'success'
  if (category === 1) return 'danger'
  return 'info'
}

const defaultBanApplyTime = () => formatServerNowPlusDays(7)

const disabledDate = (time: Date) => time.getTime() < Date.now()

const fetchList = async () => {
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
      platformOnly: true,
      ...(searchForm.liveStatus >= 0 ? {liveStatus: searchForm.liveStatus} : {}),
    })
    tableData.value = response.data || []
    selectedAnchorRows.value = []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load platform anchor list:', error)
    ElMessage.error(t('pages.anchorList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.liveStatus = 1
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

const resetBanForm = () => {
  banForm.accountId = ''
  banForm.nickname = ''
  banForm.banApplyTime = ''
  banForm.banReason = ''
  banFormRef.value?.clearValidate()
}

const resetAnchorTypeForm = () => {
  anchorTypeForm.userId = ''
  anchorTypeForm.nickname = ''
  anchorTypeForm.anchorType = 1
  anchorTypeFormRef.value?.clearValidate()
}

const openAnchorTypeDialog = (row: AnchorListItem) => {
  anchorTypeDialogTitle.value = t('pages.guildMembers.setAnchorTypeTitle', {id: row.id})
  anchorTypeForm.userId = row.id
  anchorTypeForm.nickname = row.nickname || '-'
  const userType = Number(row.userType)
  anchorTypeForm.anchorType = userType === USER_TYPE_SENIOR_ANCHOR ? 7 : 1
  anchorTypeDialogVisible.value = true
}

const submitAnchorType = async () => {
  if (!anchorTypeFormRef.value) {
    return
  }
  await anchorTypeFormRef.value.validate(async (valid: boolean) => {
    if (!valid) {
      return
    }
    anchorTypeSubmitting.value = true
    try {
      const response = await accountApi.setPlatformAnchorType({
        userId: anchorTypeForm.userId,
        anchorType: anchorTypeForm.anchorType,
      })
      if (response?.success) {
        ElMessage.success(t('pages.guildMembers.setAnchorTypeSuccess'))
        anchorTypeDialogVisible.value = false
        fetchList()
      } else {
        ElMessage.error(t('pages.guildMembers.setAnchorTypeFailed'))
      }
    } catch (error) {
      console.error('Set platform anchor type failed:', error)
      ElMessage.error(t('pages.guildMembers.setAnchorTypeRequestFailed'))
    } finally {
      anchorTypeSubmitting.value = false
    }
  })
}

const openDetail = (row: AnchorListItem) => {
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(row.id)},
  })
}

const handleSelectionChange = (rows: AnchorListItem[]) => {
  selectedAnchorRows.value = rows
}

const handleBatchImmediateSettlement = async () => {
  const rows = selectedAnchorRows.value
  if (rows.length === 0) {
    ElMessage.warning(t('pages.platformAnchorImmediateSettlement.selectRequired'))
    return
  }
  try {
    await ElMessageBox.confirm(
      t('pages.platformAnchorImmediateSettlement.confirmMessage', {count: rows.length}),
      t('pages.platformAnchorImmediateSettlement.confirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    immediateSettling.value = true
    const response = await accountApi.batchImmediateSettlePlatformAnchors({anchorIds: rows.map(row => row.id)})
    const resultMessage = t('pages.platformAnchorImmediateSettlement.result', {
      settled: response.settledCount ?? 0,
      noData: response.noDataCount ?? 0,
      fail: response.failCount ?? 0,
    })
    if ((response.failCount ?? 0) > 0) {
      ElMessage.warning({message: resultMessage, duration: 8000})
    } else {
      ElMessage.success(resultMessage)
    }
    if (response.failAnchorIds?.length) {
      ElMessage.warning({
        message: t('pages.platformAnchorImmediateSettlement.failAnchorIds', {ids: response.failAnchorIds.join(', ')}),
        duration: 8000,
      })
    }
    await fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('batch immediate platform anchor settlement failed:', error)
    ElMessage.error(t('pages.platformAnchorImmediateSettlement.requestFailed'))
  } finally {
    immediateSettling.value = false
  }
}

const salaryValidityText = (row: AnchorListItem) => {
  if (!row.salaryEffectiveStartTime || !row.salaryEffectiveEndTime) return '-'
  return `${formatServerDateTime(row.salaryEffectiveStartTime)} ~ ${formatServerDateTime(row.salaryEffectiveEndTime)}`
}

const openTransferInfo = (row: AnchorListItem) => {
  router.push({
    name: 'PlatformAnchorTransferInfoEdit',
    params: {anchorId: String(row.id)},
    query: {anchorName: row.nickname || ''},
  })
}

const handleRowCommand = (row: AnchorListItem, command: string) => {
  switch (command) {
    case 'viewDetail':
      openDetail(row)
      break
    case 'setAnchorType':
      openAnchorTypeDialog(row)
      break
    case 'transferInfo':
      openTransferInfo(row)
      break
    case 'offShelf':
      handleOffShelf(row)
      break
    case 'ban':
    case 'unban':
      toggleBanStatus(row)
      break
  }
}

const parseImportUserIds = (text: string): string[] => {
  const seen = new Set<string>()
  const ids: string[] = []
  for (const token of text.split(/[\s,，;；]+/)) {
    const id = token.trim()
    if (!/^\d+$/.test(id) || id === '0' || seen.has(id)) continue
    seen.add(id)
    ids.push(id)
  }
  return ids
}

const openSalaryImportDialog = () => {
  salaryImportForm.userIdsText = ''
  salaryImportDialogVisible.value = true
  salaryImportFormRef.value?.clearValidate()
}

const handleSalaryImportSubmit = async () => {
  if (!salaryImportFormRef.value) return
  await salaryImportFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    const ids = parseImportUserIds(salaryImportForm.userIdsText)
    if (ids.length === 0) {
      ElMessage.warning(t('pages.guildList.importEmpty'))
      return
    }
    salaryImporting.value = true
    try {
      const response = await guildApi.batchImportSalaryAnchors({ids})
      salaryImportDialogVisible.value = false
      ElMessage.success(t('pages.guildList.batchImportSalaryAnchorResult', {
        success: response.successCount ?? 0,
        fail: response.failCount ?? 0,
      }))
      const failIds = response.failIds || []
      if (failIds.length > 0) {
        const visibleIds = failIds.slice(0, 50).join(', ')
        ElMessage.warning({
          message: t('pages.guildList.batchImportSalaryAnchorFailIds', {
            ids: `${visibleIds}${failIds.length > 50 ? '…' : ''}`,
          }),
          duration: 8000,
        })
      }
      await fetchList()
    } catch (error) {
      console.error('batch import salary anchors failed:', error)
      ElMessage.error(t('pages.guildList.importFailed'))
    } finally {
      salaryImporting.value = false
    }
  })
}

const openBanDialog = (row: AnchorListItem) => {
  banForm.accountId = row.id
  banForm.nickname = row.nickname || '-'
  banForm.banApplyTime = defaultBanApplyTime()
  banForm.banReason = ''
  banDialogVisible.value = true
}

const submitBan = async () => {
  if (!banFormRef.value) return
  await banFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    banSubmitting.value = true
    try {
      const response = await accountApi.banAnchor({
        accountId: banForm.accountId,
        banApplyTime: banForm.banApplyTime,
        banReason: banForm.banReason.trim(),
      } as BanAnchorReq)
      if (response) {
        ElMessage.success(t('pages.anchorList.banSuccessNotify'))
        banDialogVisible.value = false
        fetchList()
      } else {
        ElMessage.error(t('pages.anchorList.banFailed'))
      }
    } catch (error) {
      console.error('Ban anchor failed:', error)
      ElMessage.error(t('pages.anchorList.banRequestFailed'))
    } finally {
      banSubmitting.value = false
    }
  })
}

const toggleBanStatus = async (row: AnchorListItem) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
        t('pages.anchorList.unbanConfirm', {id: row.id}),
        t('pages.anchorList.unbanTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        },
      )
      const response = await accountApi.unBanAnchor({accountId: row.id} as UnBanAnchorReq)
      if (response) {
        ElMessage.success(t('pages.anchorList.unbanSuccess'))
        fetchList()
      } else {
        ElMessage.error(t('pages.anchorList.unbanFailed'))
      }
    } catch {
      // cancelled
    }
    return
  }
  openBanDialog(row)
}

const handleOffShelf = async (row: AnchorListItem) => {
  try {
    await ElMessageBox.confirm(
      t('pages.anchorList.offShelfConfirm', {id: row.id}),
      t('common.confirmOffShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await accountApi.setLiveRoomStatus({anchorId: row.id, status: 0})
    ElMessage.success(t('pages.anchorList.offShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') return
    console.error('off shelf live room failed:', error)
    ElMessage.error(t('pages.anchorList.offShelfFailed'))
  }
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

.search-form {
  margin-bottom: 20px;
}

.search-form :deep(.el-form-item__label) {
  white-space: nowrap;
}

.search-form :deep(.el-form-item) {
  margin-right: 16px;
  margin-bottom: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.salary-import-tip {
  margin-bottom: 16px;
}
</style>

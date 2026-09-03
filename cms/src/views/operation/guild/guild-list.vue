<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.guildList.addGuild') }}</el-button>
        </div>
        <el-alert
            :closable="false"
            class="import-tip"
            show-icon
            :title="t('pages.guildList.importTip')"
            type="info"
        />

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.guildList.guildName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.guildList.guildNameSearchPlaceholder')"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchGuildList">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table
            v-loading="loading"
            :data="tableData"
            highlight-current-row
            style="width: 100%"
            @current-change="handleCurrentRowChange"
        >
          <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
          <el-table-column label="ID" prop="id" width="190">
            <template #default="{ row }">
              <el-button v-if="canViewDetail" link type="primary" @click="openDetail(row)">
                {{ row.id }}
              </el-button>
              <span v-else>{{ row.id }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.guildName')" prop="name"/>
          <el-table-column :label="t('pages.guildList.unsettledTotalIncome')" width="140">
            <template #default="{ row }">
              <span class="money-amount">{{ formatWalletBalance(row.unsettledTotalIncome) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.leader')" width="140" show-overflow-tooltip>
            <template #default="{ row }">
              {{ formatLeader(row) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.description')" prop="description" show-overflow-tooltip/>
          
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="760">
            <template #default="{ row }">
              <el-button v-if="canViewDetail" size="small" @click="openDetail(row)">
                {{ t('pages.guildList.viewDetail') }}
              </el-button>
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="can('transferInfo')"
                  size="small"
                  @click="openTransferInfoDialog(row)"
              >
                {{ t('pages.guildList.transferInfo') }}
              </el-button>
              <el-button
                  v-if="can('viewMembers')"
                  size="small"
                  @click="handleViewMembers(row)"
              >
                {{ t('pages.guildList.viewMembers') }}
              </el-button>
              <el-button
                  v-if="can('joinGuildAnchor')"
                  size="small"
                  @click="openJoinDialog(row)"
              >
                {{ t('pages.guildList.joinGuildAnchor') }}
              </el-button>
              <el-button
                  v-if="can('batchSetAnchor')"
                  size="small"
                  @click="openImportDialog(row, 1)"
              >
                {{ t('pages.guildList.importNormalAnchor') }}
              </el-button>
              <el-button
                  v-if="can('batchSetSeniorAnchor')"
                  size="small"
                  @click="openImportDialog(row, 7)"
              >
                {{ t('pages.guildList.importSeniorAnchor') }}
              </el-button>
              <el-button
                  v-if="can('offShelf')"
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                {{ t('common.offShelf') }}
              </el-button>
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
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.guildList.guildName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.guildList.guildNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.leader')" prop="leaderId">
          <div class="leader-picker-field">
            <el-input
                :model-value="leaderDisplayText"
                class="leader-picker-input"
                readonly
                :placeholder="t('pages.guildList.leaderPickPlaceholder')"
            />
            <el-button type="primary" @click="openLeaderPicker">{{ t('pages.guildList.selectLeader') }}</el-button>
            <el-button v-if="currentRow.leaderId" link type="danger" @click="clearLeader">{{ t('pages.guildList.clearLeader') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.description')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="t('pages.guildList.descriptionPlaceholder')" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <CmsUserPickerDialog v-model="leaderPickerVisible" @select="handleLeaderSelect"/>

    <el-dialog v-model="joinDialogVisible" :title="joinDialogTitle" width="480px">
      <el-form ref="joinFormRef" :model="joinForm" :rules="joinFormRules" label-width="100px">
        <el-form-item :label="t('pages.guildList.joinUserId')" prop="userId">
          <el-input v-model="joinForm.userId" clearable :placeholder="t('pages.guildList.joinUserIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.memberType')" prop="anchorType">
          <el-select v-model="joinForm.anchorType" style="width: 100%">
            <el-option :label="t('pages.guildList.normalAnchor')" :value="1"/>
            <el-option :label="t('pages.guildList.seniorAnchor')" :value="7"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="joinSubmitting" type="primary" @click="handleJoinSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="importDialogVisible" :title="importDialogTitle" width="560px">
      <el-form ref="importFormRef" :model="importForm" :rules="importFormRules" label-width="80px">
        <el-form-item :label="t('pages.guildList.importUserIds')" prop="userIdsText">
          <el-input
              v-model="importForm.userIdsText"
              :autosize="{ minRows: 8, maxRows: 16 }"
              :placeholder="t('pages.guildList.importUserIdsPlaceholder')"
              type="textarea"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="importing" type="primary" @click="handleImportSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="transferDialogVisible" :title="transferDialogTitle" destroy-on-close width="560px">
      <el-form
          ref="transferFormRef"
          v-loading="transferLoading"
          :model="transferForm"
          :rules="transferFormRules"
          label-width="100px"
      >
        <el-form-item :label="t('pages.guildList.transferCurrency')" prop="currency">
          <el-input
              v-model="transferForm.currency"
              clearable
              maxlength="8"
              :placeholder="t('pages.guildList.transferCurrencyPlaceholder')"
              @input="transferForm.currency = transferForm.currency.toUpperCase()"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferPayeeName')" prop="payeeName">
          <el-input
              v-model="transferForm.payeeName"
              clearable
              :placeholder="t('pages.guildList.transferPayeeNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferBankName')" prop="bankName">
          <el-input
              v-model="transferForm.bankName"
              clearable
              :placeholder="t('pages.guildList.transferBankNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferAccountNo')" prop="accountNo">
          <el-input
              v-model="transferForm.accountNo"
              clearable
              :placeholder="t('pages.guildList.transferAccountNoPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferBankCode')" prop="bankCode">
          <el-input
              v-model="transferForm.bankCode"
              clearable
              :placeholder="t('pages.guildList.transferBankCodePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferRemark')" prop="remark">
          <el-input
              v-model="transferForm.remark"
              :autosize="{ minRows: 2, maxRows: 4 }"
              :placeholder="t('pages.guildList.transferRemarkPlaceholder')"
              type="textarea"
          />
        </el-form-item>
        <el-form-item v-if="transferForm.updatedAt" :label="t('pages.guildList.transferLastUpdated')">
          <span>{{ transferForm.updatedAt }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="transferDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="transferSaving" type="primary" @click="handleTransferSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {guildApi} from '@/api'
import CmsUserPickerDialog from '@/components/CmsUserPickerDialog.vue'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild, GuildAnchorImportResultState, ImportGuildAnchorRow} from '@/types/api.ts'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'

const GUILD_ANCHOR_IMPORT_RESULT_KEY = 'guildAnchorImportResult'

interface SearchForm {
  name: string
}

interface GuildForm {
  id: string
  name: string
  leaderId: string
  description: string
}

interface JoinGuildForm {
  guildId: string
  guildName: string
  userId: string
  anchorType: 1 | 7
}

interface ImportGuildForm {
  guildId: string
  guildName: string
  anchorType: 1 | 7
  userIdsText: string
}

interface TransferInfoForm {
  guildId: string
  guildName: string
  currency: string
  payeeName: string
  bankName: string
  accountNo: string
  bankCode: string
  remark: string
  updatedAt: string
}

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('GuildManagement')
const canViewDetail = computed(() => can('viewDetail'))
const loading = ref(false)
const importing = ref(false)
const leaderPickerVisible = ref(false)
const selectedLeader = ref<CMSUser | null>(null)
const tableData = ref<Guild[]>([])
const selectedGuild = ref<Guild | null>(null)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const importDialogVisible = ref(false)
const importDialogTitle = ref('')
const importFormRef = ref<FormInstance>()
const importForm = ref<ImportGuildForm>({
  guildId: '',
  guildName: '',
  anchorType: 1,
  userIdsText: '',
})

const transferDialogVisible = ref(false)
const transferDialogTitle = ref('')
const transferLoading = ref(false)
const transferSaving = ref(false)
const transferFormRef = ref<FormInstance>()
const transferForm = ref<TransferInfoForm>({
  guildId: '',
  guildName: '',
  currency: '',
  payeeName: '',
  bankName: '',
  accountNo: '',
  bankCode: '',
  remark: '',
  updatedAt: '',
})

const searchForm = reactive<SearchForm>({
  name: ''
})

const joinDialogVisible = ref(false)
const joinDialogTitle = ref('')
const joinSubmitting = ref(false)
const joinFormRef = ref<FormInstance>()
const joinForm = ref<JoinGuildForm>({
  guildId: '',
  guildName: '',
  userId: '',
  anchorType: 1,
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<GuildForm>({
  id: '',
  name: '',
  leaderId: '',
  description: '',
})

const formRef = ref<FormInstance>()

const leaderDisplayText = computed(() => {
  if (selectedLeader.value) {
    return `${selectedLeader.value.name} (${selectedLeader.value.id})`
  }
  if (currentRow.value.leaderId) {
    return currentRow.value.leaderId
  }
  return ''
})

const formatLeader = (row: Guild) => {
  if (row.leaderName) {
    return `${row.leaderName} (${row.leaderId})`
  }
  return row.leaderId || '-'
}

const openLeaderPicker = () => {
  leaderPickerVisible.value = true
}

const clearLeader = () => {
  currentRow.value.leaderId = ''
  selectedLeader.value = null
}

const handleLeaderSelect = (user: CMSUser) => {
  currentRow.value.leaderId = user.id
  selectedLeader.value = user
}

const importFormRules = computed<FormRules>(() => ({
  userIdsText: [
    {required: true, message: t('pages.guildList.importUserIdsRequired'), trigger: 'blur'},
  ],
}))

const transferFormRules = computed<FormRules>(() => ({
  currency: [
    {required: true, message: t('pages.guildList.transferCurrencyRequired'), trigger: 'blur'},
    {min: 3, max: 8, message: t('pages.guildList.transferCurrencyRequired'), trigger: 'blur'},
  ],
}))

const joinFormRules = computed<FormRules>(() => ({
  userId: [
    {required: true, message: t('pages.guildList.joinUserIdRequired'), trigger: 'blur'},
    {pattern: /^\d+$/, message: t('pages.guildList.joinUserIdInvalid'), trigger: 'blur'},
  ],
  anchorType: [
    {required: true, message: t('pages.guildList.joinAnchorTypeRequired'), trigger: 'change'},
  ],
}))

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.guildList.nameRequired'), trigger: 'blur'},
    {min: 2, max: 32, message: t('pages.guildList.nameLength'), trigger: 'blur'}
  ],
  description: [
    {max: 200, message: t('pages.guildList.descriptionMaxLength'), trigger: 'blur'}
  ]
}))

const fetchGuildList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getGuildList({
      name: searchForm.name,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
    if (selectedGuild.value && !tableData.value.some(item => item.id === selectedGuild.value?.id)) {
      selectedGuild.value = null
    }
  } catch (error) {
    console.error('fetch guild list failed:', error)
    ElMessage.error(t('pages.guildList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGuildList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGuildList()
}

const formatRowIndex = (index: number) =>
    (currentPage.value - 1) * pageSize.value + index + 1

const handleCurrentRowChange = (row: Guild | null) => {
  selectedGuild.value = row
}

const handleAdd = () => {
  dialogTitle.value = t('pages.guildList.addGuild')
  currentRow.value = {
    id: '',
    name: '',
    leaderId: '',
    description: '',
  }
  selectedLeader.value = null
  dialogVisible.value = true
}

const handleEdit = (row: Guild) => {
  dialogTitle.value = t('pages.guildList.editGuild')
  const leaderId = row.leaderId && row.leaderId !== '0' ? row.leaderId : ''
  currentRow.value = {
    id: row.id,
    name: row.name,
    leaderId,
    description: row.description,
  }
  selectedLeader.value = leaderId && row.leaderName
      ? {id: leaderId, name: row.leaderName} as CMSUser
      : null
  dialogVisible.value = true
}

const handleOffShelf = async (row: Guild) => {
  try {
    await ElMessageBox.confirm(
      t('pages.guildList.offShelfConfirm', {name: row.name}),
      t('common.confirmOffShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )

    await guildApi.deleteGuild(row.id)
    ElMessage.success(t('pages.guildList.offShelfSuccess'))
    fetchGuildList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('off shelf guild failed:', error)
    ElMessage.error(t('pages.guildList.offShelfFailed'))
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const leaderId = Number(currentRow.value.leaderId) || 0
        if (currentRow.value.id) {
          await guildApi.updateGuild({...currentRow.value, leaderId})
        } else {
          const {name, description} = currentRow.value
          await guildApi.createGuild({name, leaderId, description})
        }

        ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
        dialogVisible.value = false
        fetchGuildList()
      } catch (error) {
        console.error('save guild failed:', error)
        ElMessage.error(currentRow.value.id ? t('pages.guildList.updateFailed') : t('pages.guildList.createFailed'))
      }
    }
  })
}

const formatJoinFailReason = (reason?: number) => {
  switch (reason) {
    case 1:
      return t('pages.guildAnchorImportResult.reasonUserNotFound')
    case 2:
      return t('pages.guildAnchorImportResult.reasonCancelCodeMismatch')
    case 3:
      return t('pages.guildAnchorImportResult.reasonCancelCodeExpired')
    case 4:
      return t('pages.guildAnchorImportResult.reasonAlreadyInGuild')
    case 5:
      return t('pages.guildAnchorImportResult.reasonCannotSetAnchor')
    case 6:
      return t('pages.guildAnchorImportResult.reasonAlreadyHasLiveRoom')
    default:
      return t('pages.guildAnchorImportResult.reasonUnknown')
  }
}

const openJoinDialog = (row: Guild) => {
  joinDialogTitle.value = t('pages.guildList.joinGuildAnchorTitle', {name: row.name})
  joinForm.value = {
    guildId: row.id,
    guildName: row.name,
    userId: '',
    anchorType: 1,
  }
  joinDialogVisible.value = true
  joinFormRef.value?.clearValidate()
}

const handleJoinSubmit = async () => {
  if (!joinFormRef.value) {
    return
  }
  await joinFormRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    joinSubmitting.value = true
    try {
      const response = await guildApi.joinGuildAnchor({
        guildId: joinForm.value.guildId,
        userId: joinForm.value.userId.trim(),
        anchorType: joinForm.value.anchorType,
      })
      if (response.success) {
        ElMessage.success(t('pages.guildList.joinGuildAnchorSuccess'))
        joinDialogVisible.value = false
        return
      }
      const reasonText = formatJoinFailReason(response.reason)
      const nickname = response.nickname?.trim()
      ElMessage.error(
        nickname
          ? t('pages.guildList.joinGuildAnchorFailedWithUser', {nickname, reason: reasonText})
          : t('pages.guildList.joinGuildAnchorFailed', {reason: reasonText}),
      )
    } catch (error) {
      console.error('join guild anchor failed:', error)
      ElMessage.error(t('pages.guildList.joinGuildAnchorRequestFailed'))
    } finally {
      joinSubmitting.value = false
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  fetchGuildList()
}

const parseImportUserIds = (text: string): ImportGuildAnchorRow[] => {
  const normalized = text.replace(/^\ufeff/, '').trim()
  if (!normalized) {
    return []
  }
  const headerPattern = /^(user_id|userid|用户id|id)$/i
  const tokens = normalized.split(/[\s,;\t\r\n]+/).map(item => item.trim()).filter(Boolean)
  const rows: ImportGuildAnchorRow[] = []
  const seen = new Set<string>()
  for (const token of tokens) {
    if (headerPattern.test(token)) {
      continue
    }
    if (!/^\d+$/.test(token)) {
      continue
    }
    if (seen.has(token)) {
      continue
    }
    seen.add(token)
    rows.push({userId: token})
  }
  return rows
}

const openImportDialog = (row: Guild, anchorType: 1 | 7) => {
  selectedGuild.value = row
  importDialogTitle.value = anchorType === 7
      ? t('pages.guildList.importSeniorAnchorTitle', {name: row.name})
      : t('pages.guildList.importNormalAnchorTitle', {name: row.name})
  importForm.value = {
    guildId: row.id,
    guildName: row.name,
    anchorType,
    userIdsText: '',
  }
  importDialogVisible.value = true
  importFormRef.value?.clearValidate()
}

const handleImportSubmit = async () => {
  if (!importFormRef.value) {
    return
  }
  await importFormRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    const rows = parseImportUserIds(importForm.value.userIdsText)
    if (rows.length === 0) {
      ElMessage.warning(t('pages.guildList.importEmpty'))
      return
    }
    importing.value = true
    try {
      const response = await guildApi.importGuildAnchors({
        guildId: importForm.value.guildId,
        anchorType: importForm.value.anchorType,
        rows,
      })
      const state: GuildAnchorImportResultState = {
        guildId: importForm.value.guildId,
        guildName: importForm.value.guildName,
        anchorType: importForm.value.anchorType,
        successCount: response?.successCount ?? 0,
        failCount: response?.failCount ?? 0,
        fails: response?.fails ?? [],
      }
      importDialogVisible.value = false
      sessionStorage.setItem(GUILD_ANCHOR_IMPORT_RESULT_KEY, JSON.stringify(state))
      await router.push({name: 'GuildAnchorImportResult'})
    } catch (error) {
      console.error('import guild anchors failed:', error)
      ElMessage.error(t('pages.guildList.importFailed'))
    } finally {
      importing.value = false
    }
  })
}

const openTransferInfoDialog = async (row: Guild) => {
  selectedGuild.value = row
  transferDialogTitle.value = t('pages.guildList.transferInfoTitle', {name: row.name})
  transferForm.value = {
    guildId: row.id,
    guildName: row.name,
    currency: '',
    payeeName: '',
    bankName: '',
    accountNo: '',
    bankCode: '',
    remark: '',
    updatedAt: '',
  }
  transferDialogVisible.value = true
  transferLoading.value = true
  try {
    const response = await guildApi.getGuildTransferInfo(row.id)
    const info = response?.info
    transferForm.value = {
      guildId: row.id,
      guildName: row.name,
      currency: info?.currency ?? '',
      payeeName: info?.payeeName ?? '',
      bankName: info?.bankName ?? '',
      accountNo: info?.accountNo ?? '',
      bankCode: info?.bankCode ?? '',
      remark: info?.remark ?? '',
      updatedAt: info?.updatedAt ?? '',
    }
  } catch (error) {
    console.error('fetch guild transfer info failed:', error)
    ElMessage.error(t('pages.guildList.transferFetchFailed'))
  } finally {
    transferLoading.value = false
    transferFormRef.value?.clearValidate()
  }
}

const handleTransferSave = async () => {
  if (!transferFormRef.value) {
    return
  }
  try {
    await transferFormRef.value.validate()
  } catch {
    return
  }
  transferSaving.value = true
  try {
    await guildApi.saveGuildTransferInfo({
      guildId: transferForm.value.guildId,
      currency: transferForm.value.currency.trim().toUpperCase(),
      payeeName: transferForm.value.payeeName.trim(),
      bankName: transferForm.value.bankName.trim(),
      accountNo: transferForm.value.accountNo.trim(),
      bankCode: transferForm.value.bankCode.trim(),
      remark: transferForm.value.remark.trim(),
    })
    ElMessage.success(t('pages.guildList.transferSaveSuccess'))
    transferDialogVisible.value = false
  } catch (error) {
    console.error('save guild transfer info failed:', error)
    ElMessage.error(t('pages.guildList.transferSaveFailed'))
  } finally {
    transferSaving.value = false
  }
}

onMounted(() => {
  fetchGuildList()
})

const openDetail = (row: Guild) => {
  router.push({
    name: 'GuildDetail',
    query: {
      id: row.id,
      name: row.name,
      leaderId: row.leaderId,
      leaderName: row.leaderName ?? '',
      description: row.description ?? '',
      status: String(row.status ?? 1),
      createdAt: row.createdAt ?? '',
      updatedAt: row.updatedAt ?? '',
    },
  })
}

const handleViewMembers = (row: Guild) => {
  router.push({
    name: 'GuildMembers',
    query: {
      guildId: row.id,
      guildName: row.name,
    },
  })
}
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.table-header {
  margin-bottom: 12px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.import-tip {
  margin-bottom: 16px;
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

.leader-picker-field {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.leader-picker-input {
  flex: 1;
}

</style>

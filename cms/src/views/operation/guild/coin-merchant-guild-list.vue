<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CoinMerchantGuildListManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.guildList.addGuild') }}</el-button>
          <el-button
              v-if="can('batchImmediateSettlement')"
              :disabled="selectedGuildRows.length === 0"
              :loading="immediateSettling"
              type="danger"
              @click="handleBatchImmediateSettlement"
          >
            {{ t('pages.guildImmediateSettlement.button') }}
          </el-button>
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
            class="coin-merchant-guild-table"
            :data="tableData"
            fit
            highlight-current-row
            style="width: 100%"
            @current-change="handleCurrentRowChange"
            @selection-change="handleSelectionChange"
        >
          <el-table-column align="center" header-align="center" type="selection" width="48"/>
          <el-table-column align="center" header-align="center" label="#" type="index" width="55" :index="formatRowIndex"/>
          <el-table-column align="center" header-align="center" label="ID" prop="id" width="190"/>
          <el-table-column
              class-name="coin-merchant-guild-gift-col"
              :label="t('pages.coinMerchantGuildList.giftIncomeCumulative')"
              align="center"
              header-align="center"
              min-width="160"
          >
            <template #default="{ row }">
              <span class="money-amount">{{ formatWalletBalance(row.unsettledTotalIncome) }}</span>
            </template>
          </el-table-column>
          <el-table-column
              class-name="coin-merchant-guild-name-col"
              :label="t('pages.guildList.guildName')"
              align="center"
              header-align="center"
              prop="name"
              show-overflow-tooltip
              width="120"
          />
          <el-table-column :label="t('pages.guildList.sharePercent')" align="center" header-align="center" width="120">
            <template #default="{ row }">{{ row.sharePercent ?? '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('common.actions')" align="center" header-align="center" width="100">
            <template #default="{ row }">
              <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(row, cmd)">
                <el-button size="small" type="primary">
                  {{ t('common.actions') }}
                  <el-icon class="el-icon--right">
                    <ArrowDown/>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="canViewDetail" command="viewDetail">
                      {{ t('pages.guildList.viewDetail') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('edit')" command="edit">
                      {{ t('common.edit') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('transferInfo')" command="transferInfo">
                      {{ t('pages.guildList.transferInfo') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('viewMembers')" command="viewMembers">
                      {{ t('pages.guildList.viewMembers') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('joinGuildAnchor')" divided command="joinGuildAnchor">
                      {{ t('pages.guildList.joinGuildAnchor') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('batchSetAnchor')" command="batchSetAnchor">
                      {{ t('pages.guildList.importNormalAnchor') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('batchSetSeniorAnchor')" command="batchSetSeniorAnchor">
                      {{ t('pages.guildList.importSeniorAnchor') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('offShelf')" divided command="offShelf">
                      {{ t('common.offShelf') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="140px">
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
        <el-form-item :label="t('pages.guildList.sharePercent')" prop="sharePercent">
          <el-input-number
              v-model="currentRow.sharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-alert
            :closable="false"
            show-icon
            :title="t('pages.guildList.sharePercentTip')"
            type="info"
            style="margin-bottom: 12px"
        />
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

  </div>
</template>

<script lang="ts" setup>
import {computed, nextTick, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {ArrowDown} from '@element-plus/icons-vue'
import {guildApi} from '@/api'
import CmsUserPickerDialog from '@/components/CmsUserPickerDialog.vue'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild, GuildAnchorImportResultState, ImportGuildAnchorRow} from '@/types/api.ts'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'

const GUILD_ANCHOR_IMPORT_RESULT_KEY = 'guildAnchorImportResult'
const GUILD_TYPE_COIN_MERCHANT = 1

interface SearchForm {
  name: string
}

interface GuildForm {
  id: string
  name: string
  leaderId: string
  description: string
  sharePercent: number
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

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('CoinMerchantGuildListManagement')
const canViewDetail = computed(() => can('viewDetail'))
const GUILD_ROW_ACTION_KEYS = [
  'edit',
  'transferInfo',
  'viewMembers',
  'joinGuildAnchor',
  'batchSetAnchor',
  'batchSetSeniorAnchor',
  'offShelf',
] as const
const hasRowActions = computed(() => canViewDetail.value || GUILD_ROW_ACTION_KEYS.some(key => can(key)))

const loading = ref(false)
const importing = ref(false)
const immediateSettling = ref(false)
const selectedGuildRows = ref<Guild[]>([])
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
const searchForm = reactive<SearchForm>({name: ''})

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
  sharePercent: 10,
})
const formRef = ref<FormInstance>()

const isEditingGuild = computed(() => Boolean(currentRow.value.id?.trim()))

const leaderDisplayText = computed(() => {
  if (selectedLeader.value) {
    return `${selectedLeader.value.name} (${selectedLeader.value.id})`
  }
  if (currentRow.value.leaderId) {
    return currentRow.value.leaderId
  }
  return ''
})

const importFormRules = computed<FormRules>(() => ({
  userIdsText: [
    {required: true, message: t('pages.guildList.importUserIdsRequired'), trigger: 'blur'},
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
    {min: 2, max: 32, message: t('pages.guildList.nameLength'), trigger: 'blur'},
  ],
  sharePercent: [
    {required: true, message: t('pages.coinMerchantGuildList.sharePercentRequired'), trigger: 'change'},
  ],
  description: [
    {max: 200, message: t('pages.guildList.descriptionMaxLength'), trigger: 'blur'},
  ],
}))

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

const fetchGuildList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getGuildList({
      name: searchForm.name,
      guildType: GUILD_TYPE_COIN_MERCHANT,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data
    selectedGuildRows.value = []
    total.value = response.total
    if (selectedGuild.value && !tableData.value.some(item => item.id === selectedGuild.value?.id)) {
      selectedGuild.value = null
    }
  } catch (error) {
    console.error('fetch coin merchant guild list failed:', error)
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

const formatRowIndex = (index: number) => (currentPage.value - 1) * pageSize.value + index + 1

const handleCurrentRowChange = (row: Guild | null) => {
  selectedGuild.value = row
}

const handleSelectionChange = (rows: Guild[]) => {
  selectedGuildRows.value = rows
}

const handleBatchImmediateSettlement = async () => {
  const rows = selectedGuildRows.value
  if (rows.length === 0) {
    ElMessage.warning(t('pages.guildImmediateSettlement.selectRequired'))
    return
  }
  try {
    await ElMessageBox.confirm(
      t('pages.guildImmediateSettlement.confirmMessage', {count: rows.length}),
      t('pages.guildImmediateSettlement.confirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    immediateSettling.value = true
    const response = await guildApi.batchImmediateSettleGuilds({guildIds: rows.map(row => row.id)})
    const resultMessage = t('pages.guildImmediateSettlement.result', {
      settled: response.settledCount ?? 0,
      noData: response.noDataCount ?? 0,
      fail: response.failCount ?? 0,
    })
    if ((response.failCount ?? 0) > 0) {
      ElMessage.warning({message: resultMessage, duration: 8000})
    } else {
      ElMessage.success(resultMessage)
    }
    if (response.failGuildIds?.length) {
      ElMessage.warning({
        message: t('pages.guildImmediateSettlement.failGuildIds', {ids: response.failGuildIds.join(', ')}),
        duration: 8000,
      })
    }
    await fetchGuildList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('batch immediate guild settlement failed:', error)
    ElMessage.error(t('pages.guildImmediateSettlement.requestFailed'))
  } finally {
    immediateSettling.value = false
  }
}

const resetSearch = () => {
  searchForm.name = ''
  currentPage.value = 1
  fetchGuildList()
}

const handleRowCommand = (row: Guild, command: string) => {
  switch (command) {
    case 'viewDetail':
      openDetail(row)
      break
    case 'edit':
      handleEdit(row)
      break
    case 'transferInfo':
      openTransferInfoPage(row)
      break
    case 'viewMembers':
      handleViewMembers(row)
      break
    case 'joinGuildAnchor':
      openJoinDialog(row)
      break
    case 'batchSetAnchor':
      openImportDialog(row, 1)
      break
    case 'batchSetSeniorAnchor':
      openImportDialog(row, 7)
      break
    case 'offShelf':
      handleOffShelf(row)
      break
  }
}

const handleAdd = () => {
  dialogTitle.value = t('pages.guildList.addGuild')
  currentRow.value = {
    id: '',
    name: '',
    leaderId: '',
    description: '',
    sharePercent: 10,
  }
  selectedLeader.value = null
  dialogVisible.value = true
  nextTick(() => formRef.value?.clearValidate())
}

const handleEdit = (row: Guild) => {
  dialogTitle.value = t('pages.guildList.editGuild')
  const leaderId = row.leaderId && row.leaderId !== '0' ? row.leaderId : ''
  currentRow.value = {
    id: row.id,
    name: row.name,
    leaderId,
    description: row.description,
    sharePercent: row.sharePercent ?? 10,
  }
  selectedLeader.value = leaderId && row.leaderName ? {id: leaderId, name: row.leaderName} as CMSUser : null
  dialogVisible.value = true
  nextTick(() => formRef.value?.clearValidate())
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
    console.error('off shelf coin merchant guild failed:', error)
    ElMessage.error(t('pages.guildList.offShelfFailed'))
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      return
    }
    const editing = isEditingGuild.value
    try {
      const leaderId = Number(currentRow.value.leaderId) || 0
      const payload = {
        name: currentRow.value.name,
        leaderId,
        description: currentRow.value.description,
        guildType: GUILD_TYPE_COIN_MERCHANT,
        sharePercent: currentRow.value.sharePercent,
      }
      if (editing) {
        await guildApi.updateGuild({...payload, id: currentRow.value.id})
      } else {
        await guildApi.createGuild(payload)
      }
      ElMessage.success(editing ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchGuildList()
    } catch (error) {
      console.error('save coin merchant guild failed:', error)
      ElMessage.error(editing ? t('pages.guildList.updateFailed') : t('pages.guildList.createFailed'))
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

const openTransferInfoPage = (row: Guild) => {
  router.push({
    name: 'GuildTransferInfoEdit',
    params: {guildId: row.id},
    query: {guildName: row.name},
  })
}

const openDetail = (row: Guild) => {
  router.push({
    name: 'CoinMerchantGuildDetail',
    query: {
      id: row.id,
      name: row.name,
      sharePercent: String(row.sharePercent ?? ''),
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

onMounted(() => {
  fetchGuildList()
})
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

:deep(.coin-merchant-guild-name-col .cell) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

</style>

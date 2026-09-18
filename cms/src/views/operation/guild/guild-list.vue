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
          <el-table-column :label="t('pages.guildList.guildType')" width="110">
            <template #default="{ row }">{{ guildTypeLabel(row.guildType) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.sharePercent')" width="120">
            <template #default="{ row }">{{ row.sharePercent ?? '-' }}</template>
          </el-table-column>
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
          <el-table-column :label="t('pages.guildList.creator')" width="140" show-overflow-tooltip>
            <template #default="{ row }">
              {{ formatCreator(row) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.description')" prop="description" show-overflow-tooltip/>
          
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="100">
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
        <el-form-item :label="t('pages.guildList.guildType')" prop="guildType">
          <el-select
              v-model="currentRow.guildType"
              :placeholder="t('pages.guildList.selectGuildType')"
              style="width: 100%"
              @change="onGuildTypeChange"
          >
            <el-option :label="t('pages.guildList.guildTypeNormal')" :value="0"/>
            <el-option :label="t('pages.guildList.guildTypeCoinMerchant')" :value="1"/>
          </el-select>
        </el-form-item>
        <template v-if="isCoinMerchantGuild">
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
        </template>
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

    <el-dialog
        v-model="transferDialogVisible"
        :title="transferDialogTitle"
        destroy-on-close
        width="760px"
        @closed="transferCurrencyDialogVisible = false"
    >
      <el-form
          ref="transferFormRef"
          v-loading="transferLoading"
          :model="transferForm"
          :rules="transferFormRules"
          label-width="110px"
      >
        <el-form-item :label="t('pages.guildList.transferCurrency')" prop="currency">
          <el-input
              :model-value="transferCurrencyDisplay"
              :disabled="transferLoading"
              :placeholder="t('pages.guildList.transferCurrencyPlaceholder')"
              class="transfer-currency-picker"
              readonly
              @click="openTransferCurrencyDialog"
          >
            <template #append>
              <el-button @click.stop="openTransferCurrencyDialog">
                {{ t('pages.guildList.transferCurrencySelect') }}
              </el-button>
            </template>
          </el-input>
          <div class="form-tip">{{ t('pages.guildList.transferCurrencyHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferAccountType')" prop="accountType">
          <el-select
              v-model="transferForm.accountType"
              clearable
              :disabled="!transferForm.currency"
              style="width: 100%"
              :placeholder="t('pages.guildList.transferAccountTypePlaceholder')"
              @change="onTransferAccountTypeChange"
          >
            <el-option
                v-for="accountType in transferAccountTypeOptions"
                :key="accountType"
                :label="accountTypeLabel(accountType)"
                :value="accountType"
            />
          </el-select>
          <div class="form-tip">{{ t('pages.guildList.transferAccountTypeHint') }}</div>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferBankCode')" prop="bankCode">
          <el-select
              v-model="transferForm.bankCode"
              clearable
              filterable
              :disabled="!transferForm.accountType"
              style="width: 100%"
              :placeholder="t('pages.guildList.transferBankCodePlaceholder')"
          >
            <el-option
                v-for="option in transferBankCodeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
            />
          </el-select>
          <div v-if="selectedTransferMethod" class="form-tip">
            {{ selectedTransferMethod.description }} · {{ selectedTransferMethod.limit }} {{ transferForm.currency }}
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferPayeeName')" prop="payeeName">
          <el-input
              v-model="transferForm.payeeName"
              clearable
              :placeholder="t('pages.guildList.transferPayeeNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferPhone')" prop="phone">
          <el-input
              v-model="transferForm.phone"
              clearable
              :placeholder="t('pages.guildList.transferPhonePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferEmail')" prop="email">
          <el-input
              v-model="transferForm.email"
              clearable
              :placeholder="t('pages.guildList.transferEmailPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferAccountNo')" prop="accountNo">
          <el-input
              v-model="transferForm.accountNo"
              clearable
              :placeholder="t('pages.guildList.transferAccountNoPlaceholder')"
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

    <el-dialog
        v-model="transferCurrencyDialogVisible"
        :title="t('pages.guildList.transferCurrencyDialogTitle')"
        append-to-body
        width="880px"
    >
      <div class="transfer-currency-toolbar">
        <el-input
            v-model="transferCurrencyKeyword"
            clearable
            :placeholder="t('pages.guildList.transferCurrencySearchPlaceholder')"
        />
      </div>
      <el-tabs v-model="transferCurrencyRegion" class="transfer-currency-tabs">
        <el-tab-pane
            v-for="region in transferCurrencyRegionOptions"
            :key="region"
            :label="transferCurrencyRegionLabel(region)"
            :name="region"
        />
      </el-tabs>
      <el-table
          :data="filteredTransferCountries"
          highlight-current-row
          max-height="440"
          style="width: 100%"
          @row-click="handleTransferCurrencyPick"
      >
        <el-table-column :label="t('pages.guildList.transferCurrency')" prop="currency" width="100"/>
        <el-table-column :label="t('pages.guildList.transferCurrencyRegion')" width="115">
          <template #default="{ row }">
            {{ transferCurrencyRegionLabel(transferCurrencyRegionFor(row)) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildList.transferCountry')" min-width="245">
          <template #default="{ row }">
            <span class="transfer-country-option">
              <img v-if="row.icon" :src="row.icon" alt="" class="transfer-flag"/>
              <span>
                {{ row.nameZh || row.nameEn || row.countryCode }}
                <span v-if="row.nameZh && row.nameEn" class="transfer-country-en"> / {{ row.nameEn }}</span>
                ({{ row.countryCode }})
              </span>
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildList.transferAccountType')" min-width="230">
          <template #default="{ row }">
            <el-tag
                v-for="accountType in row.accountTypes"
                :key="accountType"
                class="transfer-account-type-tag"
                effect="plain"
            >
              {{ accountTypeLabel(accountType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="90">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="handleTransferCurrencyPick(row)">
              {{ t('pages.guildList.transferCurrencySelect') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {ArrowDown} from '@element-plus/icons-vue'
import {guildApi} from '@/api'
import {liveRevenueShareCfgApi} from '@/api/modules/live-revenue-share-cfg'
import CmsUserPickerDialog from '@/components/CmsUserPickerDialog.vue'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild, GuildAnchorImportResultState, GuildTransferCountryOption, ImportGuildAnchorRow} from '@/types/api.ts'
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
  guildType: number
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

interface TransferInfoForm {
  guildId: string
  guildName: string
  countryCode: string
  currency: string
  accountType: string
  payeeName: string
  phone: string
  email: string
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
const transferCurrencyDialogVisible = ref(false)
const transferCurrencyKeyword = ref('')
const transferCurrencyRegion = ref('ALL')
const transferLoading = ref(false)
const transferSaving = ref(false)
const transferFormRef = ref<FormInstance>()
const transferForm = ref<TransferInfoForm>({
  guildId: '',
  guildName: '',
  countryCode: '',
  currency: '',
  accountType: '',
  payeeName: '',
  phone: '',
  email: '',
  bankName: '',
  accountNo: '',
  bankCode: '',
  remark: '',
  updatedAt: '',
})
const transferCountries = ref<GuildTransferCountryOption[]>([])
const transferCurrencyRegionOrder = [
  'NORTH_AMERICA',
  'EUROPE',
  'SOUTH_AMERICA',
  'ASIA',
  'MIDDLE_EAST',
  'AFRICA',
] as const
const transferCurrencyRegionFor = (item: GuildTransferCountryOption) => {
  return item.region || ''
}
const transferCurrencyRegionOptions = computed(() => {
  const available = new Set(transferCountries.value.map(transferCurrencyRegionFor))
  return ['ALL', ...transferCurrencyRegionOrder.filter(region => available.has(region))]
})
const transferCurrencyRegionLabel = (region: string) => {
  const labelKeys: Record<string, string> = {
    ALL: 'transferCurrencyRegionAll',
    NORTH_AMERICA: 'transferCurrencyRegionNorthAmerica',
    EUROPE: 'transferCurrencyRegionEurope',
    SOUTH_AMERICA: 'transferCurrencyRegionSouthAmerica',
    ASIA: 'transferCurrencyRegionAsia',
    MIDDLE_EAST: 'transferCurrencyRegionMiddleEast',
    AFRICA: 'transferCurrencyRegionAfrica',
  }
  const key = labelKeys[region]
  return key ? t(`pages.guildList.${key}`) : region
}
const normalizeTransferCurrencySearch = (value: string) => {
  return value.toLocaleLowerCase().replace(/[\s_-]+/g, '')
}
const filteredTransferCountries = computed(() => {
  const keywords = transferCurrencyKeyword.value
      .trim()
      .split(/\s+/)
      .map(normalizeTransferCurrencySearch)
      .filter(Boolean)
  return transferCountries.value.filter(item => {
    const region = transferCurrencyRegionFor(item)
    if (transferCurrencyRegion.value !== 'ALL' && region !== transferCurrencyRegion.value) {
      return false
    }
    if (keywords.length === 0) return true
    const searchable = normalizeTransferCurrencySearch([
      item.currency,
      item.countryCode,
      item.nameZh,
      item.nameEn,
      region,
      ...item.accountTypes,
      ...(item.methods ?? []).flatMap(method => [method.accountType, method.bankCode, method.description]),
    ].join(' '))
    return keywords.every(keyword => searchable.includes(keyword))
  })
})
const selectedTransferCurrency = computed(() => {
  const currency = transferForm.value.currency.trim().toUpperCase()
  return transferCountries.value.find(item => item.currency === currency) ?? null
})
const transferCurrencyDisplay = computed(() => {
  const selected = selectedTransferCurrency.value
  if (!selected) return ''
  const countryName = selected.nameZh || selected.nameEn || selected.countryCode
  return `${selected.currency} · ${countryName} (${selected.countryCode})`
})
const transferAccountTypeOptions = computed(() => {
  return selectedTransferCurrency.value?.accountTypes ?? []
})
const transferBankCodeOptions = computed(() => {
  const accountType = transferForm.value.accountType
  if (!accountType) return []
  return (selectedTransferCurrency.value?.methods ?? [])
      .filter(method => method.accountType === accountType)
      .map(method => ({
        label: `${method.description} (${method.bankCode}) · ${method.limit} ${transferForm.value.currency}`,
        value: method.bankCode,
      }))
})
const selectedTransferMethod = computed(() => {
  const accountType = transferForm.value.accountType
  const bankCode = transferForm.value.bankCode
  if (!accountType || !bankCode) return null
  return (selectedTransferCurrency.value?.methods ?? []).find(method => (
      method.accountType === accountType && method.bankCode === bankCode
  )) ?? null
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
  guildType: 0,
  sharePercent: 10,
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

const formatCreator = (row: Guild) => {
  if (row.creatorName) {
    return `${row.creatorName} (${row.creatorId})`
  }
  if (row.creatorId && row.creatorId !== '0') {
    return row.creatorId
  }
  return '-'
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
    {required: true, message: t('pages.guildList.transferCurrencyRequired'), trigger: 'change'},
  ],
  accountType: [
    {required: true, message: t('pages.guildList.transferAccountTypePlaceholder'), trigger: 'change'},
  ],
  payeeName: [
    {required: true, message: t('pages.guildList.transferPayeeNamePlaceholder'), trigger: 'blur'},
  ],
  phone: [
    {required: true, message: t('pages.guildList.transferPhonePlaceholder'), trigger: 'blur'},
  ],
  email: [
    {required: true, message: t('pages.guildList.transferEmailPlaceholder'), trigger: 'blur'},
  ],
  accountNo: [
    {required: true, message: t('pages.guildList.transferAccountNoPlaceholder'), trigger: 'blur'},
  ],
  bankCode: [
    {required: true, message: t('pages.guildList.transferBankCodePlaceholder'), trigger: 'blur'},
  ],
}))

const accountTypeLabel = (accountType: string) => {
  if (accountType === 'BANK_ACCOUNT') {
    return `${t('pages.guildList.transferAccountTypeBank')} (BANK_ACCOUNT)`
  }
  if (accountType === 'EWALLET') {
    return `${t('pages.guildList.transferAccountTypeEWallet')} (EWALLET)`
  }
  return accountType
}

const onTransferCurrencyChange = (currency: string) => {
  const hit = transferCountries.value.find(c => c.currency === currency)
  transferForm.value.countryCode = hit?.countryCode || ''
  const accountTypes = hit?.accountTypes ?? []
  transferForm.value.accountType = accountTypes.length === 1 ? (accountTypes[0] ?? '') : ''
  transferForm.value.bankName = ''
  transferForm.value.bankCode = ''
  if (transferForm.value.accountType) {
    const methods = (hit?.methods ?? []).filter(method => method.accountType === transferForm.value.accountType)
    transferForm.value.bankCode = methods.length === 1 ? (methods[0]?.bankCode ?? '') : ''
  }
}

const openTransferCurrencyDialog = () => {
  if (transferLoading.value) return
  transferCurrencyKeyword.value = ''
  transferCurrencyRegion.value = 'ALL'
  transferCurrencyDialogVisible.value = true
}

const handleTransferCurrencyPick = (row: GuildTransferCountryOption) => {
  if (!row?.currency) return
  onTransferCurrencyChange(row.currency)
  transferCurrencyDialogVisible.value = false
  transferFormRef.value?.clearValidate(['currency', 'accountType'])
}

const onTransferAccountTypeChange = () => {
  transferForm.value.bankName = ''
  transferForm.value.bankCode = ''
  const methods = (selectedTransferCurrency.value?.methods ?? [])
      .filter(method => method.accountType === transferForm.value.accountType)
  transferForm.value.bankCode = methods.length === 1 ? (methods[0]?.bankCode ?? '') : ''
}

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
  guildType: [
    {required: true, message: t('pages.guildList.selectGuildType'), trigger: 'change'},
  ],
  description: [
    {max: 200, message: t('pages.guildList.descriptionMaxLength'), trigger: 'blur'}
  ]
}))

const loadDefaultSharePercent = async () => {
  try {
    const response = await liveRevenueShareCfgApi.getCfg()
    return response.cfg?.guildSharePercent ?? 10
  } catch {
    return 10
  }
}

const isCoinMerchantGuild = computed(() => Number(currentRow.value.guildType) === 1)

const guildTypeLabel = (guildType?: number) => {
  if (guildType === 1) return t('pages.guildList.guildTypeCoinMerchant')
  return t('pages.guildList.guildTypeNormal')
}

const onGuildTypeChange = async (guildType: number) => {
  if (Number(guildType) === 1) {
    return
  }
  currentRow.value.sharePercent = await loadDefaultSharePercent()
}

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

const handleRowCommand = (row: Guild, command: string) => {
  switch (command) {
    case 'viewDetail':
      openDetail(row)
      break
    case 'edit':
      handleEdit(row)
      break
    case 'transferInfo':
      openTransferInfoDialog(row)
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

const handleAdd = async () => {
  dialogTitle.value = t('pages.guildList.addGuild')
  const sharePercent = await loadDefaultSharePercent()
  currentRow.value = {
    id: '',
    name: '',
    leaderId: '',
    description: '',
    guildType: 0,
    sharePercent,
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
    guildType: row.guildType ?? 0,
    sharePercent: row.sharePercent ?? 10,
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
        const payload = {
          name: currentRow.value.name,
          leaderId,
          description: currentRow.value.description,
          guildType: currentRow.value.guildType,
          sharePercent: currentRow.value.sharePercent,
        }
        if (currentRow.value.id) {
          await guildApi.updateGuild({...payload, id: currentRow.value.id})
        } else {
          await guildApi.createGuild(payload)
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
    countryCode: '',
    currency: '',
    accountType: '',
    payeeName: '',
    phone: '',
    email: '',
    bankName: '',
    accountNo: '',
    bankCode: '',
    remark: '',
    updatedAt: '',
  }
  transferCountries.value = []
  transferCurrencyDialogVisible.value = false
  transferDialogVisible.value = true
  transferLoading.value = true
  try {
    const response = await guildApi.getGuildTransferInfo(row.id)
    const info = response?.info
    transferCountries.value = response?.countries ?? []
    transferForm.value = {
      guildId: row.id,
      guildName: row.name,
      countryCode: info?.countryCode ?? '',
      currency: info?.currency ?? '',
      accountType: info?.accountType || '',
      payeeName: info?.payeeName ?? '',
      phone: info?.phone ?? '',
      email: info?.email ?? '',
      bankName: info?.bankName ?? '',
      accountNo: info?.accountNo ?? '',
      bankCode: info?.bankCode ?? '',
      remark: info?.remark ?? '',
      updatedAt: info?.updatedAt ?? '',
    }
    if (!transferForm.value.currency && transferForm.value.countryCode) {
      transferForm.value.currency = transferCountries.value.find(
          item => item.countryCode === transferForm.value.countryCode,
      )?.currency ?? ''
    }
    const selectedCurrency = transferCountries.value.find(
        item => item.currency === transferForm.value.currency,
    )
    if (selectedCurrency) {
      transferForm.value.countryCode = selectedCurrency.countryCode
      if (!selectedCurrency.accountTypes.includes(transferForm.value.accountType)) {
        transferForm.value.accountType = selectedCurrency.accountTypes.length === 1
            ? (selectedCurrency.accountTypes[0] ?? '')
            : ''
      }
      const methodExists = (selectedCurrency.methods ?? []).some(method => (
          method.accountType === transferForm.value.accountType &&
          method.bankCode === transferForm.value.bankCode
      ))
      if (!methodExists) {
        transferForm.value.bankCode = ''
      }
    } else {
      transferForm.value.countryCode = ''
      transferForm.value.currency = ''
      transferForm.value.accountType = ''
      transferForm.value.bankCode = ''
      transferForm.value.bankName = ''
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
      countryCode: transferForm.value.countryCode.trim().toUpperCase(),
      currency: transferForm.value.currency.trim().toUpperCase(),
      accountType: transferForm.value.accountType.trim().toUpperCase(),
      payeeName: transferForm.value.payeeName.trim(),
      phone: transferForm.value.phone.trim(),
      email: transferForm.value.email.trim(),
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
      creatorId: row.creatorId ?? '',
      creatorName: row.creatorName ?? '',
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

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.transfer-currency-picker {
  cursor: pointer;
}

.transfer-currency-picker :deep(.el-input__inner) {
  cursor: pointer;
}

.transfer-currency-toolbar {
  margin-bottom: 2px;
}

.transfer-currency-tabs :deep(.el-tabs__header) {
  margin-bottom: 12px;
}

.transfer-country-option {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.transfer-country-en {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.transfer-account-type-tag + .transfer-account-type-tag {
  margin-left: 6px;
}

.transfer-flag {
  width: 18px;
  height: 12px;
  object-fit: cover;
  border-radius: 2px;
}
</style>

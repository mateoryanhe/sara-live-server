<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-button @click="goBack">{{ t('pages.guildList.back') }}</el-button>
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-form
          ref="transferFormRef"
          v-loading="loading"
          :model="transferForm"
          :rules="transferFormRules"
          class="editor-form"
          label-width="110px"
      >
        <el-form-item :label="t('pages.guildList.transferCurrency')" prop="currency">
          <el-input
              :model-value="transferCurrencyDisplay"
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferCurrencyPlaceholder')"
              class="transfer-currency-picker"
              readonly
              @click="openTransferCurrencyDialog"
          >
            <template #append>
              <el-button :disabled="loading || saving" @click.stop="openTransferCurrencyDialog">
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
              :disabled="loading || saving || !transferForm.currency"
              :placeholder="t('pages.guildList.transferAccountTypePlaceholder')"
              style="width: 100%"
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
              :disabled="loading || saving || !transferForm.accountType"
              :placeholder="t('pages.guildList.transferBankCodePlaceholder')"
              style="width: 100%"
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
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferPayeeNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferPhone')" prop="phone">
          <el-input
              v-model="transferForm.phone"
              clearable
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferPhonePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferEmail')" prop="email">
          <el-input
              v-model="transferForm.email"
              clearable
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferEmailPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferAccountNo')" prop="accountNo">
          <el-input
              v-model="transferForm.accountNo"
              clearable
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferAccountNoPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.guildList.transferRemark')" prop="remark">
          <el-input
              v-model="transferForm.remark"
              :autosize="{minRows: 2, maxRows: 4}"
              :disabled="loading || saving"
              :placeholder="t('pages.guildList.transferRemarkPlaceholder')"
              type="textarea"
          />
        </el-form-item>
        <el-form-item v-if="transferForm.updatedAt" :label="t('pages.guildList.transferLastUpdated')">
          <span>{{ transferForm.updatedAt }}</span>
        </el-form-item>
        <el-form-item>
          <el-button @click="goBack">{{ t('pages.guildList.back') }}</el-button>
          <el-button
              v-if="can('transferInfo')"
              :loading="saving"
              :disabled="loading"
              type="primary"
              @click="handleTransferSave"
          >
            {{ t('common.save') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-dialog
        v-model="transferCurrencyDialogVisible"
        :title="t('pages.guildList.transferCurrencyDialogTitle')"
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
          <template #default="{row}">
            {{ transferCurrencyRegionLabel(transferCurrencyRegionFor(row)) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildList.transferCountry')" min-width="245">
          <template #default="{row}">
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
          <template #default="{row}">
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
          <template #default="{row}">
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
import {computed, nextTick, onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {guildApi} from '@/api'
import type {GuildTransferCountryOption} from '@/types/api.ts'
import {usePagePermission} from '@/composables/usePagePermission'

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

const emptyTransferForm = (): TransferInfoForm => ({
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

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {can} = usePagePermission('GuildManagement')
const guildId = computed(() => String(route.params.guildId || '').trim())
const guildName = computed(() => {
  const value = route.query.guildName
  return String(Array.isArray(value) ? (value[0] ?? '') : (value ?? '')).trim()
})
const pageTitle = computed(() => t('pages.guildList.transferInfoTitle', {
  name: guildName.value || guildId.value,
}))
const loading = ref(false)
const saving = ref(false)
const transferCurrencyDialogVisible = ref(false)
const transferCurrencyKeyword = ref('')
const transferCurrencyRegion = ref('ALL')
const transferFormRef = ref<FormInstance>()
const transferForm = ref<TransferInfoForm>(emptyTransferForm())
const transferCountries = ref<GuildTransferCountryOption[]>([])
const transferCurrencyRegionOrder = [
  'NORTH_AMERICA',
  'EUROPE',
  'SOUTH_AMERICA',
  'ASIA',
  'MIDDLE_EAST',
  'AFRICA',
] as const

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
    {required: true, message: t('pages.guildList.transferBankCodePlaceholder'), trigger: 'change'},
  ],
}))

const transferCurrencyRegionFor = (item: GuildTransferCountryOption) => item.region || ''
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
const normalizeTransferCurrencySearch = (value: string) => (
    value.toLocaleLowerCase().replace(/[\s_-]+/g, '')
)
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
const transferAccountTypeOptions = computed(() => selectedTransferCurrency.value?.accountTypes ?? [])
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

const accountTypeLabel = (accountType: string) => {
  if (accountType === 'BANK_ACCOUNT') {
    return `${t('pages.guildList.transferAccountTypeBank')} (BANK_ACCOUNT)`
  }
  if (accountType === 'EWALLET') {
    return `${t('pages.guildList.transferAccountTypeEWallet')} (EWALLET)`
  }
  return accountType
}

const normalizeLoadedTransferInfo = () => {
  if (!transferForm.value.currency && transferForm.value.countryCode) {
    transferForm.value.currency = transferCountries.value.find(
        item => item.countryCode === transferForm.value.countryCode,
    )?.currency ?? ''
  }
  const selectedCurrency = transferCountries.value.find(
      item => item.currency === transferForm.value.currency,
  )
  if (!selectedCurrency) {
    transferForm.value.countryCode = ''
    transferForm.value.currency = ''
    transferForm.value.accountType = ''
    transferForm.value.bankCode = ''
    transferForm.value.bankName = ''
    return
  }
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
}

const loadTransferInfo = async () => {
  if (!guildId.value) {
    ElMessage.error(t('pages.guildList.transferFetchFailed'))
    return
  }
  loading.value = true
  transferForm.value = {
    ...emptyTransferForm(),
    guildId: guildId.value,
    guildName: guildName.value,
  }
  transferCountries.value = []
  transferCurrencyDialogVisible.value = false
  try {
    const response = await guildApi.getGuildTransferInfo(guildId.value)
    const info = response?.info
    transferCountries.value = response?.countries ?? []
    transferForm.value = {
      guildId: guildId.value,
      guildName: guildName.value,
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
    normalizeLoadedTransferInfo()
  } catch (error) {
    console.error('fetch guild transfer info failed:', error)
    ElMessage.error(t('pages.guildList.transferFetchFailed'))
  } finally {
    loading.value = false
    await nextTick()
    transferFormRef.value?.clearValidate()
  }
}

const onTransferCurrencyChange = (currency: string) => {
  const hit = transferCountries.value.find(item => item.currency === currency)
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
  if (loading.value || saving.value) return
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

const handleTransferSave = async () => {
  if (!transferFormRef.value) return
  try {
    await transferFormRef.value.validate()
  } catch {
    return
  }
  saving.value = true
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
    await loadTransferInfo()
  } catch (error) {
    console.error('save guild transfer info failed:', error)
    ElMessage.error(t('pages.guildList.transferSaveFailed'))
  } finally {
    saving.value = false
  }
}

const goBack = () => {
  router.push({name: 'GuildManagement'})
}

onMounted(loadTransferInfo)
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: bold;
}

.editor-form {
  max-width: 760px;
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

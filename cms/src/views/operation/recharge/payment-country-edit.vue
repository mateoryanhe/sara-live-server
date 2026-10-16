<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-button @click="goBack">{{ t('pages.paymentCountryCfg.backToCountryList') }}</el-button>
            <span>{{ pageTitle }}</span>
          </div>
        </div>
      </template>

      <el-alert :closable="false" :title="detailNotice" class="notice" show-icon type="info"/>
      <el-skeleton v-if="loading" :rows="10" animated/>
      <el-empty v-else-if="!editForm" :description="t('pages.paymentCountryCfg.invalidCountry')"/>

      <div v-else class="editor-body">
        <div class="country-summary">
          <strong>{{ editForm.countryNameZh || editForm.countryNameEn || editForm.countryCode }}</strong>
          <span>{{ editForm.countryNameEn }} ({{ editForm.countryCode }})</span>
        </div>

        <el-tabs v-model="activeEditTab" class="edit-tabs" type="border-card">
          <el-tab-pane name="config" :label="t(isUserWithdrawPayout
              ? 'pages.paymentCountryCfg.withdrawPayoutSettings'
              : isCoinMerchant
                  ? 'pages.paymentCountryCfg.collectionSettings'
                  : 'pages.paymentCountryCfg.basicSettings')">
            <el-form class="config-form" label-width="120px">
              <el-form-item :label="t('pages.paymentCountryCfg.supportedCurrencies')">
                <el-tag v-for="code in editForm.supportedCurrencies" :key="code" class="currency-tag" size="small">
                  {{ code }}
                </el-tag>
              </el-form-item>

              <el-form-item :label="t('pages.paymentCountryCfg.paymentCurrency')" required>
                <el-select v-model="editForm.currencyCode" :disabled="!can('edit') || saving" style="width: 100%">
                  <el-option v-for="code in editForm.supportedCurrencies" :key="code"
                             :label="currencyLabel(code)" :value="code"/>
                </el-select>
              </el-form-item>

              <template v-if="usesCombinedPaymentFields">
                <el-form-item v-if="!fixedWithdrawPayoutPayType"
                              :label="t('pages.paymentCountryCfg.paymentType')" required>
                  <el-select v-model="editForm.selectedPayType" clearable filterable
                             :disabled="!can('edit') || saving"
                             :placeholder="t('pages.paymentCountryCfg.paymentTypePlaceholder')" style="width: 100%">
                    <el-option v-for="type in paymentTypeOptions" :key="type" :label="type" :value="type"/>
                  </el-select>
                </el-form-item>

                <el-form-item :label="t('pages.paymentCountryCfg.paymentCode')" required>
                  <div v-if="fixedWithdrawPayoutPayType" class="fixed-pay-type-row">
                    <el-tag size="large" type="info">{{ fixedWithdrawPayoutPayType }}</el-tag>
                    <span class="form-tip inline-tip">
                      {{ t('pages.paymentCountryCfg.withdrawPayoutFixedPayTypeTip') }}
                    </span>
                  </div>
                  <el-select v-model="editForm.selectedMethodKey" clearable filterable
                             :disabled="!can('edit') || saving || !effectiveSelectedPayType"
                             :placeholder="t('pages.paymentCountryCfg.paymentCodePlaceholder')" style="width: 100%"
                             :class="{ 'pay-code-select': fixedWithdrawPayoutPayType }">
                    <el-option v-for="method in paymentCodeOptions" :key="methodKey(method)"
                               :label="paymentCodeLabel(method)" :value="methodKey(method)"/>
                  </el-select>
                </el-form-item>

                <el-alert v-if="isCoinMerchant" :closable="false" :title="t('pages.paymentCountryCfg.coinMerchantCombinedTip')"
                          class="credential-tip" show-icon type="info"/>
                <el-alert v-else-if="isUserWithdrawPayout" :closable="false" :title="t('pages.paymentCountryCfg.withdrawPayoutCombinedTip')"
                          class="credential-tip" show-icon type="info"/>
              </template>

              <el-form-item :label="t('pages.paymentCountryCfg.appVisible')">
                <el-switch v-model="editForm.enabled" :disabled="!can('edit') || saving"/>
                <span class="visibility-hint">
                  {{ editForm.enabled ? t('pages.paymentCountryCfg.enabled') : t('pages.paymentCountryCfg.disabled') }}
                </span>
              </el-form-item>

              <el-form-item v-if="can('edit')">
                <el-button type="primary" :loading="saving" @click="saveConfig">
                  {{ t(usesCombinedPaymentFields
                    ? 'pages.paymentCountryCfg.saveCollectionConfig'
                    : 'pages.paymentCountryCfg.saveBasic') }}
                </el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <el-tab-pane v-if="isCoinMerchant" name="official" :label="t('pages.paymentCountryCfg.availableMethods')">
            <el-alert :closable="false"
                      :title="t('pages.paymentCountryCfg.currentCurrencyAvailableMethods', {currency: editForm.currencyCode})"
                      class="currency-context" show-icon type="info"/>
            <el-table :data="paymentMethodsForCurrency" border max-height="520" size="small"
                      :empty-text="t('pages.paymentCountryCfg.noAvailableMethods')">
              <el-table-column :label="t('pages.paymentCountryCfg.selected')" width="82" align="center">
                <template #default="{ row }">
                  <el-tag v-if="isMethodSelected(row)" size="small" type="success">
                    {{ t('pages.paymentCountryCfg.selected') }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="t('pages.paymentCountryCfg.paymentType')" min-width="135" prop="payType"/>
              <el-table-column :label="t('pages.paymentCountryCfg.paymentCode')" min-width="190" prop="inBankCode"/>
              <el-table-column :label="t('pages.paymentCountryCfg.limit')" min-width="145">
                <template #default="{ row }">{{ paymentLimitLabel(row) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.paymentCountryCfg.description')" min-width="180" prop="description"/>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {
  coinMerchantPaymentCountryCfgApi,
  normalizePaymentCountryCfgGroups,
  paymentCountryCfgApi,
  userWithdrawPayoutCountryCfgApi,
  type PaymentCountryCfgItem,
  type PaymentCountryPaymentMethod,
} from '@/api/modules/payment-country-cfg'
import {usePagePermission} from '@/composables/usePagePermission'
import {
  usePaymentCountryCfgScope,
  userWithdrawPayoutFixedPayType,
} from '@/composables/usePaymentCountryCfgScope'

type PaymentCountryEditForm = PaymentCountryCfgItem & {
  selectedPayType: string
  selectedMethodKey: string
}

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {
  isCoinMerchant,
  isUserWithdrawPayout,
  usesCombinedPaymentFields,
  rootRouteName,
  continentRouteName,
  permissionPage,
} = usePaymentCountryCfgScope()
const detailNotice = computed(() => {
  if (isUserWithdrawPayout.value) {
    return t('pages.paymentCountryCfg.withdrawPayoutDetailNotice')
  }
  return t(isCoinMerchant.value
      ? 'pages.paymentCountryCfg.coinMerchantDetailNotice'
      : 'pages.paymentCountryCfg.normalDetailNotice')
})
const {can} = usePagePermission(permissionPage)
const loading = ref(false)
const saving = ref(false)
const activeEditTab = ref('config')
const editForm = ref<PaymentCountryEditForm | null>(null)
const persistedRow = ref<PaymentCountryCfgItem | null>(null)
const continent = computed(() => String(route.params.continent || '').trim().toUpperCase())
const countryCode = computed(() => String(route.params.countryCode || '').trim().toUpperCase())
const fixedWithdrawPayoutPayType = computed(() => {
  if (!isUserWithdrawPayout.value) return null
  const code = String(editForm.value?.countryCode || countryCode.value || '').trim().toUpperCase()
  return userWithdrawPayoutFixedPayType(code)
})

const effectiveSelectedPayType = computed(() => (
    fixedWithdrawPayoutPayType.value || editForm.value?.selectedPayType || ''
))
const pageTitle = computed(() => t('pages.paymentCountryCfg.editPageTitle', {
  country: editForm.value?.countryNameZh || editForm.value?.countryNameEn || countryCode.value,
}))

const methodKey = (method: Pick<PaymentCountryPaymentMethod, 'currencyCode' | 'payType' | 'inBankCode'>) => (
    `${method.currencyCode}\u0000${method.payType}\u0000${method.inBankCode}`
)

const paymentMethodsForCurrency = computed(() => {
  const form = editForm.value
  if (!form) return []
  let methods = form.paymentMethods.filter((method) => method.available && method.currencyCode === form.currencyCode)
  const fixed = fixedWithdrawPayoutPayType.value
  if (fixed) {
    methods = methods.filter((method) => method.payType === fixed)
  }
  return methods
})

const paymentTypeOptions = computed(() => {
  const fixed = fixedWithdrawPayoutPayType.value
  if (fixed) return [fixed]
  return Array.from(new Set(paymentMethodsForCurrency.value.map((method) => method.payType)))
})

const paymentCodeOptions = computed(() => {
  const payType = effectiveSelectedPayType.value
  if (!payType) return []
  return paymentMethodsForCurrency.value.filter((method) => method.payType === payType)
})

const selectedPaymentMethod = computed(() => {
  const form = editForm.value
  if (!form?.selectedMethodKey) return null
  return paymentCodeOptions.value.find((method) => methodKey(method) === form.selectedMethodKey) || null
})

const currencyLabel = (code: string) => code === 'USD'
    ? t('pages.paymentCountryCfg.usdCurrency', {code})
    : t('pages.paymentCountryCfg.localCurrency', {code})

const paymentLimitLabel = (method: PaymentCountryPaymentMethod) => (
    method.minAmount && method.maxAmount ? `${method.minAmount}-${method.maxAmount} ${method.currencyCode}` : ''
)

const paymentCodeLabel = (method: PaymentCountryPaymentMethod) => {
  const detail = [method.description, paymentLimitLabel(method)].filter(Boolean).join(' · ')
  return `${method.inBankCode}${detail ? ` — ${detail}` : ''}`
}

const isMethodSelected = (method: PaymentCountryPaymentMethod) => (
    methodKey(method) === editForm.value?.selectedMethodKey
)

const storedMethodForCurrency = (row: PaymentCountryCfgItem, currencyCode: string) => {
  const stored = row.selectedPaymentMethods.find((method) => method.currencyCode === currencyCode)
  if (!stored) return null
  return row.paymentMethods.find((method) => method.available && methodKey(method) === methodKey(stored)) || null
}

const methodsForFormCurrency = (form: PaymentCountryEditForm) => (
    form.paymentMethods.filter((method) => method.available && method.currencyCode === form.currencyCode)
)

const enforceFixedWithdrawPayoutPayType = (form: PaymentCountryEditForm) => {
  const fixed = userWithdrawPayoutFixedPayType(form.countryCode)
  if (!fixed) return
  form.selectedPayType = fixed
  const allowedKeys = new Set(methodsForFormCurrency(form)
      .filter((method) => method.payType === fixed)
      .map(methodKey))
  if (form.selectedMethodKey && !allowedKeys.has(form.selectedMethodKey)) {
    form.selectedMethodKey = ''
  }
}

const applyCurrencyDefaults = (form: PaymentCountryEditForm, currencyCode: string) => {
  const stored = persistedRow.value
  const useStored = !!stored?.id && stored.currencyCode === currencyCode
  const method = useStored ? storedMethodForCurrency(stored, currencyCode) : null
  form.selectedPayType = method?.payType || ''
  form.selectedMethodKey = method ? methodKey(method) : ''
  enforceFixedWithdrawPayoutPayType(form)
}

const buildEditForm = (row: PaymentCountryCfgItem): PaymentCountryEditForm => {
  const form: PaymentCountryEditForm = {
    ...row,
    supportedCurrencies: [...row.supportedCurrencies],
    paymentMethods: [...row.paymentMethods],
    selectedPaymentMethods: (row.selectedPaymentMethods || []).map((method) => ({...method})),
    selectedPayType: '',
    selectedMethodKey: '',
  }
  applyCurrencyDefaults(form, form.currencyCode)
  return form
}

const fetchConfig = async () => {
  loading.value = true
  editForm.value = null
  persistedRow.value = null
  try {
    const response = isUserWithdrawPayout.value
        ? await userWithdrawPayoutCountryCfgApi.getConfig()
        : isCoinMerchant.value
            ? await coinMerchantPaymentCountryCfgApi.getConfig()
            : await paymentCountryCfgApi.getConfig()
    const groups = normalizePaymentCountryCfgGroups(response.continents)
    const group = groups.find((item) => item.continent === continent.value)
    if (!group) {
      ElMessage.error(t('pages.paymentCountryCfg.invalidContinent'))
      await router.replace({name: rootRouteName.value})
      return
    }
    const row = group.countries.find((item) => item.countryCode === countryCode.value)
    if (!row) {
      ElMessage.error(t('pages.paymentCountryCfg.invalidCountry'))
      await router.replace({name: continentRouteName.value, params: {continent: continent.value}})
      return
    }
    persistedRow.value = {
      ...row,
      supportedCurrencies: [...row.supportedCurrencies],
      paymentMethods: [...row.paymentMethods],
      selectedPaymentMethods: (row.selectedPaymentMethods || []).map((method) => ({...method})),
    }
    editForm.value = buildEditForm(row)
  } catch (error) {
    console.error('fetch collection country config failed:', error)
    ElMessage.error(t('pages.paymentCountryCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const saveConfig = async () => {
  const form = editForm.value
  if (!can('edit') || !form || saving.value) return
  saving.value = true
  try {
    let response: {success: boolean; id: string}
    if (usesCombinedPaymentFields.value) {
      if (fixedWithdrawPayoutPayType.value) {
        enforceFixedWithdrawPayoutPayType(form)
      }
      const method = selectedPaymentMethod.value
      if (!method) {
        ElMessage.warning(t('pages.paymentCountryCfg.methodRequired'))
        return
      }
      const payType = fixedWithdrawPayoutPayType.value || method.payType
      const payload = {
        countryCode: form.countryCode,
        currencyCode: form.currencyCode,
        payType,
        inBankCode: method.inBankCode,
        enabled: form.enabled,
      }
      response = isUserWithdrawPayout.value
          ? await userWithdrawPayoutCountryCfgApi.saveConfig(payload)
          : await coinMerchantPaymentCountryCfgApi.saveConfig(payload)
      form.selectedPaymentMethods = [{
        currencyCode: form.currencyCode,
        payType,
        inBankCode: method.inBankCode,
      }]
    } else {
      response = await paymentCountryCfgApi.saveBasic({
        countryCode: form.countryCode,
        currencyCode: form.currencyCode,
        enabled: form.enabled,
      })
    }
    form.id = response.id || form.id
    persistedRow.value = {
      ...form,
      supportedCurrencies: [...form.supportedCurrencies],
      paymentMethods: [...form.paymentMethods],
      selectedPaymentMethods: form.selectedPaymentMethods.map((method) => ({...method})),
    }
    ElMessage.success(t('pages.paymentCountryCfg.saveSuccess', {
      country: form.countryNameZh || form.countryNameEn || form.countryCode,
    }))
  } catch (error) {
    console.error('save collection country config failed:', error)
    ElMessage.error(t('pages.paymentCountryCfg.saveFailed'))
  } finally {
    saving.value = false
  }
}

const goBack = () => router.push({name: continentRouteName.value, params: {continent: continent.value}})

watch(() => editForm.value?.currencyCode, (currency, previousCurrency) => {
  const form = editForm.value
  if (!form || !previousCurrency || !currency || currency === previousCurrency) return
  applyCurrencyDefaults(form, currency)
})

watch(() => editForm.value?.selectedPayType, (payType) => {
  const form = editForm.value
  if (!form) return
  if (fixedWithdrawPayoutPayType.value) {
    enforceFixedWithdrawPayoutPayType(form)
    return
  }
  if (!payType) {
    form.selectedMethodKey = ''
    return
  }
  const allowedKeys = new Set(paymentMethodsForCurrency.value
      .filter((method) => method.payType === payType)
      .map(methodKey))
  if (!allowedKeys.has(form.selectedMethodKey)) form.selectedMethodKey = ''
})

watch(fixedWithdrawPayoutPayType, () => {
  const form = editForm.value
  if (form) enforceFixedWithdrawPayoutPayType(form)
})

// 普通用户和经销商编辑页都会被 keep-alive 缓存，每次打开都重新读取服务器配置。
onActivated(fetchConfig)
</script>

<style scoped>
.card-header,
.header-title {
  display: flex;
  align-items: center;
}

.card-header { justify-content: space-between; }
.header-title { gap: 12px; font-weight: 600; }
.notice { margin-bottom: 18px; }
.editor-body { max-width: 1180px; }
.country-summary { display: flex; flex-direction: column; gap: 4px; margin: 0 0 16px 4px; }
.country-summary span { color: var(--el-text-color-secondary); font-size: 12px; }
.edit-tabs { min-height: 480px; }
.config-form { max-width: 1000px; }
.currency-tag { margin: 2px 4px 2px 0; }
.visibility-hint { margin-left: 10px; color: var(--el-text-color-secondary); }
.form-tip { display: block; margin-top: 6px; color: var(--el-text-color-secondary); font-size: 12px; }
.fixed-pay-type-row { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-bottom: 10px; }
.inline-tip { display: inline; margin-top: 0; flex: 1; min-width: 200px; }
.pay-code-select { width: 100%; }
.credential-tip,
.currency-context { margin: 0 0 18px; }
</style>

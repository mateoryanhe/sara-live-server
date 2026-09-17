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
          <el-tab-pane name="config" :label="t(isCoinMerchant
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

              <template v-if="isCoinMerchant">
                <el-form-item :label="t('pages.paymentCountryCfg.appId')" required>
                  <el-input v-model="editForm.appId" clearable inputmode="numeric" :disabled="!can('edit') || saving"
                            :placeholder="t('pages.paymentCountryCfg.appIdPlaceholder')"/>
                </el-form-item>

                <el-form-item :label="t('pages.paymentCountryCfg.paymentType')" required>
                  <el-select v-model="editForm.selectedPayType" clearable filterable
                             :disabled="!can('edit') || saving"
                             :placeholder="t('pages.paymentCountryCfg.paymentTypePlaceholder')" style="width: 100%">
                    <el-option v-for="type in paymentTypeOptions" :key="type" :label="type" :value="type"/>
                  </el-select>
                </el-form-item>

                <el-form-item :label="t('pages.paymentCountryCfg.paymentCode')" required>
                  <el-select v-model="editForm.selectedMethodKey" clearable filterable
                             :disabled="!can('edit') || saving || !editForm.selectedPayType"
                             :placeholder="t('pages.paymentCountryCfg.paymentCodePlaceholder')" style="width: 100%">
                    <el-option v-for="method in paymentCodeOptions" :key="methodKey(method)"
                               :label="paymentCodeLabel(method)" :value="methodKey(method)"/>
                  </el-select>
                </el-form-item>

                <el-alert :closable="false" :title="t('pages.paymentCountryCfg.coinMerchantCombinedTip')"
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
                  {{ t(isCoinMerchant
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
import {computed, onMounted, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {
  coinMerchantPaymentCountryCfgApi,
  normalizePaymentCountryCfgGroups,
  paymentCountryCfgApi,
  type PaymentCountryCfgItem,
  type PaymentCountryPaymentMethod,
} from '@/api/modules/payment-country-cfg'
import {usePagePermission} from '@/composables/usePagePermission'

type PaymentCountryEditForm = Omit<PaymentCountryCfgItem, 'appId'> & {
  appId: string
  selectedPayType: string
  selectedMethodKey: string
}

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const isCoinMerchant = route.meta.paymentCountryCfgScope === 'coinMerchant'
const rootRouteName = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgManagement' : 'PaymentCountryCfgManagement'
const continentRouteName = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgContinent' : 'PaymentCountryCfgContinent'
const permissionPage = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgManagement' : 'PaymentCountryCfgManagement'
const detailNotice = computed(() => t(isCoinMerchant
    ? 'pages.paymentCountryCfg.coinMerchantDetailNotice'
    : 'pages.paymentCountryCfg.normalDetailNotice'))
const {can} = usePagePermission(permissionPage)
const loading = ref(false)
const saving = ref(false)
const activeEditTab = ref('config')
const editForm = ref<PaymentCountryEditForm | null>(null)
const persistedRow = ref<PaymentCountryCfgItem | null>(null)
const defaultAppIds = ref<Record<string, number>>({})
const continent = computed(() => String(route.params.continent || '').trim().toUpperCase())
const countryCode = computed(() => String(route.params.countryCode || '').trim().toUpperCase())
const pageTitle = computed(() => t('pages.paymentCountryCfg.editPageTitle', {
  country: editForm.value?.countryNameZh || editForm.value?.countryNameEn || countryCode.value,
}))

const methodKey = (method: Pick<PaymentCountryPaymentMethod, 'currencyCode' | 'payType' | 'inBankCode'>) => (
    `${method.currencyCode}\u0000${method.payType}\u0000${method.inBankCode}`
)

const paymentMethodsForCurrency = computed(() => {
  const form = editForm.value
  if (!form) return []
  return form.paymentMethods.filter((method) => method.available && method.currencyCode === form.currencyCode)
})

const paymentTypeOptions = computed(() => Array.from(new Set(
    paymentMethodsForCurrency.value.map((method) => method.payType),
)))

const paymentCodeOptions = computed(() => {
  const form = editForm.value
  if (!form?.selectedPayType) return []
  return paymentMethodsForCurrency.value.filter((method) => method.payType === form.selectedPayType)
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

const applyCurrencyDefaults = (form: PaymentCountryEditForm, currencyCode: string) => {
  const stored = persistedRow.value
  const useStored = !!stored?.id && stored.currencyCode === currencyCode
  const method = useStored ? storedMethodForCurrency(stored, currencyCode) : null
  form.appId = String(useStored && stored.appId > 0 ? stored.appId : (defaultAppIds.value[currencyCode] || ''))
  form.selectedPayType = method?.payType || ''
  form.selectedMethodKey = method ? methodKey(method) : ''
}

const buildEditForm = (row: PaymentCountryCfgItem): PaymentCountryEditForm => {
  const form: PaymentCountryEditForm = {
    ...row,
    appId: '',
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
    const response = isCoinMerchant
        ? await coinMerchantPaymentCountryCfgApi.getConfig()
        : await paymentCountryCfgApi.getConfig()
    defaultAppIds.value = response.defaultAppIds || {}
    const groups = normalizePaymentCountryCfgGroups(response.continents)
    const group = groups.find((item) => item.continent === continent.value)
    if (!group) {
      ElMessage.error(t('pages.paymentCountryCfg.invalidContinent'))
      await router.replace({name: rootRouteName})
      return
    }
    const row = group.countries.find((item) => item.countryCode === countryCode.value)
    if (!row) {
      ElMessage.error(t('pages.paymentCountryCfg.invalidCountry'))
      await router.replace({name: continentRouteName, params: {continent: continent.value}})
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
    if (isCoinMerchant) {
      const appId = Number(form.appId)
      const method = selectedPaymentMethod.value
      if (!Number.isSafeInteger(appId) || appId <= 0) {
        ElMessage.warning(t('pages.paymentCountryCfg.credentialRequired'))
        return
      }
      if (!method) {
        ElMessage.warning(t('pages.paymentCountryCfg.methodRequired'))
        return
      }
      response = await coinMerchantPaymentCountryCfgApi.saveConfig({
        countryCode: form.countryCode,
        currencyCode: form.currencyCode,
        appId,
        payType: method.payType,
        inBankCode: method.inBankCode,
        enabled: form.enabled,
      })
      form.appId = String(appId)
      form.selectedPaymentMethods = [{
        currencyCode: form.currencyCode,
        payType: method.payType,
        inBankCode: method.inBankCode,
      }]
    } else {
      response = await paymentCountryCfgApi.saveBasic({
        countryCode: form.countryCode,
        currencyCode: form.currencyCode,
        appId: 0,
        enabled: form.enabled,
      })
    }
    form.id = response.id || form.id
    persistedRow.value = {
      ...form,
      appId: Number(form.appId) || 0,
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

const goBack = () => router.push({name: continentRouteName, params: {continent: continent.value}})

watch(() => editForm.value?.currencyCode, (currency, previousCurrency) => {
  const form = editForm.value
  if (!form || !previousCurrency || !currency || currency === previousCurrency) return
  applyCurrencyDefaults(form, currency)
})

watch(() => editForm.value?.selectedPayType, (payType) => {
  const form = editForm.value
  if (!form) return
  if (!payType) {
    form.selectedMethodKey = ''
    return
  }
  const allowedKeys = new Set(paymentMethodsForCurrency.value
      .filter((method) => method.payType === payType)
      .map(methodKey))
  if (!allowedKeys.has(form.selectedMethodKey)) form.selectedMethodKey = ''
})

onMounted(fetchConfig)
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
.credential-tip,
.currency-context { margin: 0 0 18px; }
</style>

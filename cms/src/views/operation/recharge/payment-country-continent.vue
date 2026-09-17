<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <el-button @click="goBack">{{ t('pages.paymentCountryCfg.backToContinents') }}</el-button>
            <span>{{ continentLabel(continent) }}</span>
          </div>
        </div>
      </template>

      <el-alert :closable="false" :title="detailNotice" class="notice" show-icon type="info"/>

      <div class="filter-row">
        <el-input v-model="keyword" clearable :placeholder="t('pages.paymentCountryCfg.searchPlaceholder')"
                  style="max-width: 520px"/>
      </div>

      <el-skeleton v-if="loading" :rows="8" animated/>
      <el-empty v-else-if="countries.length === 0" :description="t('pages.paymentCountryCfg.noCountries')"/>
      <el-table v-else :data="countries" row-key="countryCode" stripe>
        <el-table-column :label="t('pages.paymentCountryCfg.country')" min-width="210">
          <template #default="{ row }">
            <div class="country-cell">
              <strong>{{ row.countryNameZh || row.countryNameEn || row.countryCode }}</strong>
              <span>{{ row.countryNameEn }} ({{ row.countryCode }})</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column :label="t('pages.paymentCountryCfg.supportedCurrencies')" min-width="150">
          <template #default="{ row }">
            <el-tag v-for="code in row.supportedCurrencies" :key="code" class="currency-tag" size="small">
              {{ code }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column :label="t('pages.paymentCountryCfg.paymentCurrency')" min-width="120">
          <template #default="{ row }"><el-tag effect="plain" type="primary">{{ row.currencyCode }}</el-tag></template>
        </el-table-column>

        <el-table-column v-if="isCoinMerchant" :label="t('pages.paymentCountryCfg.paymentType')" min-width="190">
          <template #default="{ row }">
            <div v-if="selectedPayTypes(row).length" class="tag-list">
              <el-tag v-for="type in selectedPayTypes(row)" :key="type" size="small">{{ type }}</el-tag>
            </div>
            <span v-else>{{ t('pages.paymentCountryCfg.notSet') }}</span>
          </template>
        </el-table-column>

        <el-table-column v-if="isCoinMerchant" :label="t('pages.paymentCountryCfg.paymentCode')" min-width="230">
          <template #default="{ row }">
            <div v-if="selectedPaymentMethodsForCurrency(row).length" class="tag-list">
              <el-tag v-for="method in selectedPaymentMethodsForCurrency(row)" :key="methodKey(method)" size="small" type="success">
                {{ method.inBankCode }}
              </el-tag>
            </div>
            <span v-else>{{ t('pages.paymentCountryCfg.notSet') }}</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('pages.paymentCountryCfg.appVisible')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? t('pages.paymentCountryCfg.enabled') : t('pages.paymentCountryCfg.disabled') }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column fixed="right" :label="t('common.actions')" width="150" align="center">
          <template #default="{ row }">
            <el-button v-if="can('edit')" link type="primary" @click="openEditPage(row)">
              {{ t('pages.paymentCountryCfg.edit') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="tip">{{ bottomTip }}</div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {
  coinMerchantPaymentCountryCfgApi,
  normalizePaymentCountryCfgGroups,
  paymentCountryCfgApi,
  type PaymentCountryCfgGroup,
  type PaymentCountryCfgItem,
  type PaymentCountryPaymentMethod,
} from '@/api/modules/payment-country-cfg'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const isCoinMerchant = route.meta.paymentCountryCfgScope === 'coinMerchant'
const configApi = isCoinMerchant ? coinMerchantPaymentCountryCfgApi : paymentCountryCfgApi
const rootRouteName = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgManagement' : 'PaymentCountryCfgManagement'
const editRouteName = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgEdit' : 'PaymentCountryCfgEdit'
const permissionPage = isCoinMerchant ? 'CoinMerchantPaymentCountryCfgManagement' : 'PaymentCountryCfgManagement'
const detailNotice = computed(() => t(isCoinMerchant
    ? 'pages.paymentCountryCfg.coinMerchantDetailNotice'
    : 'pages.paymentCountryCfg.normalDetailNotice'))
const bottomTip = computed(() => t(isCoinMerchant
    ? 'pages.paymentCountryCfg.coinMerchantLocalTip'
    : 'pages.paymentCountryCfg.globalCashierTip'))
const {can} = usePagePermission(permissionPage)
const loading = ref(false)
const keyword = ref('')
const groups = ref<PaymentCountryCfgGroup[]>([])
const continent = computed(() => String(route.params.continent || '').trim().toUpperCase())

const continentNames: Record<string, { zh: string; en: string }> = {
  NORTH_AMERICA: {zh: '北美洲', en: 'North America'},
  EUROPE: {zh: '欧洲', en: 'Europe'},
  SOUTH_AMERICA: {zh: '南美洲', en: 'South America'},
  ASIA: {zh: '亚洲', en: 'Asia'},
  MIDDLE_EAST: {zh: '中东', en: 'Middle East'},
  AFRICA: {zh: '非洲', en: 'Africa'},
}

const continentLabel = (value: string) => {
  const item = continentNames[value]
  return item ? `${item.zh} / ${item.en}` : value
}

const methodKey = (method: Pick<PaymentCountryPaymentMethod, 'currencyCode' | 'payType' | 'inBankCode'>) => (
    `${method.currencyCode}\u0000${method.payType}\u0000${method.inBankCode}`
)

const selectedPaymentMethodsForCurrency = (row: PaymentCountryCfgItem) => row.selectedPaymentMethods
    .filter((method) => method.currencyCode === row.currencyCode)

const selectedPayTypes = (row: PaymentCountryCfgItem) => Array.from(new Set(
    selectedPaymentMethodsForCurrency(row).map((method) => method.payType),
))

const activeGroup = computed(() => groups.value.find((group) => group.continent === continent.value))
const countries = computed(() => {
  const rows = activeGroup.value?.countries || []
  const query = keyword.value.trim().toLowerCase()
  if (!query) return rows
  return rows.filter((row) => {
    const selectedMethods = selectedPaymentMethodsForCurrency(row)
    return [
      row.countryCode, row.countryNameZh, row.countryNameEn, row.currencyCode,
      ...row.supportedCurrencies,
      ...selectedMethods.flatMap((method) => [method.payType, method.inBankCode]),
    ].some((value) => String(value || '').toLowerCase().includes(query))
  })
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const response = await configApi.getConfig()
    groups.value = normalizePaymentCountryCfgGroups(response.continents)
    if (!groups.value.some((group) => group.continent === continent.value)) {
      ElMessage.error(t('pages.paymentCountryCfg.invalidContinent'))
      await router.replace({name: rootRouteName})
    }
  } catch (error) {
    console.error('fetch collection country config failed:', error)
    ElMessage.error(t('pages.paymentCountryCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const openEditPage = (row: PaymentCountryCfgItem) => {
  router.push({
    name: editRouteName,
    params: {continent: continent.value, countryCode: row.countryCode},
  })
}

const goBack = () => router.push({name: rootRouteName})

// 页面被 keep-alive 缓存；从编辑页返回时重新读取服务端缓存，立即显示刚保存的币种和可见性。
onActivated(fetchConfig)
</script>

<style scoped>
.card-header,
.header-title,
.filter-row,
.tag-list {
  display: flex;
  align-items: center;
}

.card-header { justify-content: space-between; }
.header-title { gap: 12px; font-weight: 600; }
.notice, .filter-row { margin-bottom: 16px; }
.country-cell { display: flex; flex-direction: column; gap: 4px; }
.country-cell span, .tip { color: var(--el-text-color-secondary); font-size: 12px; }
.currency-tag { margin: 2px 4px 2px 0; }
.tag-list { flex-wrap: wrap; gap: 5px; }
.tip { margin-top: 16px; }
</style>

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="landingNotice"
          class="notice"
          show-icon
          type="info"
      />

      <el-skeleton v-if="loading" :rows="6" animated/>
      <el-empty v-else-if="groups.length === 0" :description="t('pages.paymentCountryCfg.noCountries')"/>
      <div v-else class="continent-grid">
        <el-card
            v-for="group in groups"
            :key="group.continent"
            class="continent-card"
            shadow="hover"
            tabindex="0"
            @click="openContinent(group.continent)"
            @keyup.enter="openContinent(group.continent)"
        >
          <div class="continent-card__body">
            <div>
              <div class="continent-name">{{ continentLabel(group.continent) }}</div>
              <div class="country-count">
                {{ t('pages.paymentCountryCfg.countryCount', {total: group.countries.length}) }}
              </div>
            </div>
            <div class="continent-card__action">
              <el-tag round :type="enabledCount(group) > 0 ? 'success' : 'info'">
                {{ t('pages.paymentCountryCfg.enabledCount', {enabled: enabledCount(group), total: group.countries.length}) }}
              </el-tag>
              <span>{{ t('pages.paymentCountryCfg.openContinent') }}</span>
              <el-icon><ArrowRight/></el-icon>
            </div>
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ArrowRight} from '@element-plus/icons-vue'
import {ElMessage} from 'element-plus'
import {
  coinMerchantPaymentCountryCfgApi,
  normalizePaymentCountryCfgGroups,
  paymentCountryCfgApi,
  type PaymentCountryCfgGroup,
} from '@/api/modules/payment-country-cfg'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const isCoinMerchant = computed(() => route.meta.paymentCountryCfgScope === 'coinMerchant')
const configApi = computed(() => isCoinMerchant.value ? coinMerchantPaymentCountryCfgApi : paymentCountryCfgApi)
const pageTitle = computed(() => t(isCoinMerchant.value
    ? 'menu.CoinMerchantPaymentCountryCfgManagement'
    : 'menu.PaymentCountryCfgManagement'))
const landingNotice = computed(() => t(isCoinMerchant.value
    ? 'pages.paymentCountryCfg.coinMerchantLandingNotice'
    : 'pages.paymentCountryCfg.normalLandingNotice'))
const continentRouteName = computed(() => isCoinMerchant.value
    ? 'CoinMerchantPaymentCountryCfgContinent'
    : 'PaymentCountryCfgContinent')
const loading = ref(false)
const groups = ref<PaymentCountryCfgGroup[]>([])

const continentNames: Record<string, { zh: string; en: string }> = {
  NORTH_AMERICA: {zh: '北美洲', en: 'North America'},
  EUROPE: {zh: '欧洲', en: 'Europe'},
  SOUTH_AMERICA: {zh: '南美洲', en: 'South America'},
  ASIA: {zh: '亚洲', en: 'Asia'},
  MIDDLE_EAST: {zh: '中东', en: 'Middle East'},
  AFRICA: {zh: '非洲', en: 'Africa'},
}

const continentLabel = (continent: string) => {
  const item = continentNames[continent]
  return item ? `${item.zh} / ${item.en}` : continent
}

const enabledCount = (group: PaymentCountryCfgGroup) => group.countries.filter((row) => row.enabled).length

const fetchConfig = async () => {
  loading.value = true
  try {
    const response = await configApi.value.getConfig()
    groups.value = normalizePaymentCountryCfgGroups(response.continents)
  } catch (error) {
    console.error('fetch payment country config failed:', error)
    ElMessage.error(t('pages.paymentCountryCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const openContinent = (continent: string) => {
  router.push({name: continentRouteName.value, params: {continent}})
}

// 页面被 keep-alive 缓存；从国家列表返回时重新读取启用数量。
onActivated(fetchConfig)
</script>

<style scoped>
.card-header,
.continent-card__body,
.continent-card__action {
  display: flex;
  align-items: center;
}

.card-header,
.continent-card__body {
  justify-content: space-between;
}

.notice {
  margin-bottom: 20px;
}

.continent-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.continent-card {
  cursor: pointer;
  transition: transform 0.18s ease, border-color 0.18s ease;
}

.continent-card:hover,
.continent-card:focus-visible {
  border-color: var(--el-color-primary-light-5);
  outline: none;
  transform: translateY(-2px);
}

.continent-name {
  color: var(--el-text-color-primary);
  font-size: 17px;
  font-weight: 600;
}

.country-count {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.continent-card__action {
  gap: 8px;
  color: var(--el-color-primary);
  font-size: 13px;
}
</style>

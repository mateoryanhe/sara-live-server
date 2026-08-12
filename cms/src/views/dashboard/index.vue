<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('pages.dashboard.basicData') }}</span>
          <el-button :loading="loading" @click="fetchSysStat">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <div class="basic-stat-rows">
        <div
            v-for="(row, rowIndex) in basicStatCardRows"
            :key="rowIndex"
            class="basic-stat-row"
        >
          <div
              v-for="card in row"
              :key="card.key"
              class="stat-card"
              :class="card.theme"
          >
            <div class="stat-label">{{ card.label }}</div>
            <div class="stat-value">
              {{ card.format === 'count' ? formatCount(card.value) : formatAmount(card.value) }}
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="user-stat-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('pages.dashboard.userData') }}</span>
          <el-button :loading="trendLoading" @click="fetchUserStatTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-tabs v-model="activePeriod" lazy @tab-change="handleTabChange">
        <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
          <UserStatChart ref="dailyLineChartRef" :data="userStatTrend.daily" :title="t('pages.dashboard.activeRegisterDaily')"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              :metric-tabs="barMetricTabs"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="dailyBarChartRef"
                :data="userStatTrend.daily"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
          <UserStatChart ref="weeklyLineChartRef" :data="userStatTrend.weekly" :title="t('pages.dashboard.activeRegisterWeekly')"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              :metric-tabs="barMetricTabs"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="weeklyBarChartRef"
                :data="userStatTrend.weekly"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
          <UserStatChart ref="monthlyLineChartRef" :data="userStatTrend.monthly" :title="t('pages.dashboard.activeRegisterMonthly')"/>
          <BarMetricSection
              :metric-key="activeBarMetric"
              :metric-tabs="barMetricTabs"
              @update:metric-key="(v) => { activeBarMetric = v; handleBarMetricChange() }"
          >
            <UserStatBarChart
                ref="monthlyBarChartRef"
                :data="userStatTrend.monthly"
                :metric-key="activeBarMetric"
                :title="barChartTitle"
            />
          </BarMetricSection>
        </el-tab-pane>
      </el-tabs>

      <ResourceMetricChart
          ref="onlineChartRef"
          :data="resourceMetricTrend.points"
          metric-type="online"
          :title="t('pages.dashboard.onlineChartTitle')"
          class="online-chart"
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, defineAsyncComponent, nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {RESOURCE_METRIC_MAX_POINTS, sysStatApi} from '@/api'
import type {ResourceMetricTrend, SysStat, UserStatTrend} from '@/types/api'
import BarMetricSection from './components/bar-metric-section.vue'
import {getUserStatBarMetricTabs, USER_STAT_BAR_SERIES} from './user-stat-bar-series'
import {formatAmount} from '@/utils/number-format'

const UserStatChart = defineAsyncComponent(() => import('./components/user-stat-chart.vue'))
const UserStatBarChart = defineAsyncComponent(() => import('./components/user-stat-bar-chart.vue'))
const ResourceMetricChart = defineAsyncComponent(() => import('@/views/config/components/resource-metric-chart.vue'))

const {t, locale} = useI18n()
const loading = ref(false)
const trendLoading = ref(false)
const activePeriod = ref('daily')
const enabledBarSeries = getUserStatBarMetricTabs()
const activeBarMetric = ref(enabledBarSeries[0]?.key ?? 'rechargeUser')

const barMetricTabs = computed(() =>
    enabledBarSeries.map((item) => ({
      key: item.key,
      label: t(`pages.dashboard.${item.labelKey}`),
    })),
)

const barChartPeriodSuffix = computed(() => {
  if (activePeriod.value === 'weekly') {
    return t('pages.dashboard.last12Weeks')
  }
  if (activePeriod.value === 'monthly') {
    return t('pages.dashboard.last12Months')
  }
  return t('pages.dashboard.last30Days')
})

const barChartTitle = computed(() => {
  const metric = USER_STAT_BAR_SERIES.find((item) => item.key === activeBarMetric.value)
  const label = metric ? t(`pages.dashboard.${metric.labelKey}`) : ''
  return t('pages.dashboard.chartTitleWithPeriod', {
    label,
    period: barChartPeriodSuffix.value,
  })
})

const dailyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const weeklyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const monthlyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const dailyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const weeklyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const monthlyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const onlineChartRef = ref<InstanceType<typeof ResourceMetricChart>>()

const sysStat = reactive<SysStat>({
  totalGold: 0,
  totalGoldConsume: 0,
  totalDiamondConsume: 0,
  totalRecharge: 0,
  totalWithdraw: 0,
  totalRegisterUser: 0,
  todayRecharge: 0,
  todayGoldConsume: 0,
  todayDiamondConsume: 0,
  todayRegisterUser: 0,
})

const resourceMetricTrend = reactive<ResourceMetricTrend>({
  points: [],
})

const userStatTrend = reactive<UserStatTrend>({
  daily: [],
  weekly: [],
  monthly: [],
})

type BasicStatCardKey = keyof Pick<
    SysStat,
    | 'totalGold'
    | 'todayRecharge'
    | 'todayGoldConsume'
    | 'todayDiamondConsume'
    | 'todayRegisterUser'
    | 'totalRecharge'
    | 'totalWithdraw'
    | 'totalGoldConsume'
    | 'totalDiamondConsume'
    | 'totalRegisterUser'
>

type BasicStatCardConfig = {
  key: BasicStatCardKey
  labelKey: string
  theme: string
  format: 'amount' | 'count'
}

type BasicStatCard = BasicStatCardConfig & {
  label: string
  value: number | string
}

const BASIC_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'totalGold', labelKey: 'statTotalGold', theme: 'stat-card-gold', format: 'amount'},
  {key: 'todayRecharge', labelKey: 'statTodayRecharge', theme: 'stat-card-today-recharge', format: 'amount'},
  {key: 'todayGoldConsume', labelKey: 'statTodayGoldConsume', theme: 'stat-card-today-gold-consume', format: 'amount'},
  {key: 'todayDiamondConsume', labelKey: 'statTodayDiamondConsume', theme: 'stat-card-today-diamond-consume', format: 'amount'},
  {key: 'todayRegisterUser', labelKey: 'statTodayRegisterUser', theme: 'stat-card-today-register', format: 'count'},
  {key: 'totalRecharge', labelKey: 'statTotalRecharge', theme: 'stat-card-recharge', format: 'amount'},
  {key: 'totalWithdraw', labelKey: 'statTotalWithdraw', theme: 'stat-card-withdraw', format: 'amount'},
  {key: 'totalGoldConsume', labelKey: 'statTotalGoldConsume', theme: 'stat-card-total-gold-consume', format: 'amount'},
  {key: 'totalDiamondConsume', labelKey: 'statTotalDiamondConsume', theme: 'stat-card-total-diamond-consume', format: 'amount'},
  {key: 'totalRegisterUser', labelKey: 'statTotalRegisterUser', theme: 'stat-card-register', format: 'count'},
]

const BASIC_STAT_ROW_SIZE = 5

const basicStatCardRows = computed<BasicStatCard[][]>(() => {
  const cards: BasicStatCard[] = BASIC_STAT_CARD_CONFIG.map((cfg) => ({
    ...cfg,
    label: t(`pages.dashboard.${cfg.labelKey}`),
    value: sysStat[cfg.key] ?? 0,
  }))
  const rows: BasicStatCard[][] = []
  for (let i = 0; i < cards.length; i += BASIC_STAT_ROW_SIZE) {
    rows.push(cards.slice(i, i + BASIC_STAT_ROW_SIZE))
  }
  return rows
})

const fetchSysStat = async () => {
  loading.value = true
  try {
    const data = await sysStatApi.getSysStat()
    sysStat.totalGold = data.totalGold ?? 0
    sysStat.totalGoldConsume = data.totalGoldConsume ?? 0
    sysStat.totalDiamondConsume = data.totalDiamondConsume ?? 0
    sysStat.totalRecharge = data.totalRecharge ?? 0
    sysStat.totalWithdraw = data.totalWithdraw ?? 0
    sysStat.totalRegisterUser = data.totalRegisterUser ?? 0
    sysStat.todayRecharge = data.todayRecharge ?? 0
    sysStat.todayGoldConsume = data.todayGoldConsume ?? 0
    sysStat.todayDiamondConsume = data.todayDiamondConsume ?? 0
    sysStat.todayRegisterUser = data.todayRegisterUser ?? 0
  } catch (error) {
    console.error('fetch sys stat failed:', error)
    ElMessage.error(t('pages.dashboard.fetchSysStatFailed'))
  } finally {
    loading.value = false
  }
}

const fetchResourceMetricTrend = async () => {
  const data = await sysStatApi.getResourceMetricOnlineTrend({limit: RESOURCE_METRIC_MAX_POINTS})
  resourceMetricTrend.points = data.points || []
  await nextTick()
  setTimeout(() => {
    onlineChartRef.value?.resize()
  }, 0)
}

const fetchUserStatTrend = async () => {
  trendLoading.value = true
  try {
    const [data] = await Promise.all([
      sysStatApi.getUserStatTrend(),
      fetchResourceMetricTrend(),
    ])
    userStatTrend.daily = data.daily || []
    userStatTrend.weekly = data.weekly || []
    userStatTrend.monthly = data.monthly || []
    await nextTick()
    setTimeout(() => {
      resizeActiveChart()
    }, 0)
  } catch (error) {
    console.error('fetch user stat trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchUserTrendFailed'))
  } finally {
    trendLoading.value = false
  }
}

const resizeActiveChart = () => {
  onlineChartRef.value?.resize()
  if (activePeriod.value === 'daily') {
    dailyLineChartRef.value?.resize()
    dailyBarChartRef.value?.resize()
    return
  }
  if (activePeriod.value === 'weekly') {
    weeklyLineChartRef.value?.resize()
    weeklyBarChartRef.value?.resize()
    return
  }
  monthlyLineChartRef.value?.resize()
  monthlyBarChartRef.value?.resize()
}

const handleTabChange = async () => {
  await nextTick()
  resizeActiveChart()
}

const handleBarMetricChange = async () => {
  await nextTick()
  resizeActiveChart()
}

const formatCount = (value: string | number | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString(locale.value)
}

onMounted(() => {
  fetchSysStat()
  fetchUserStatTrend()
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

.user-stat-card {
  margin-top: 20px;
}

.online-chart {
  margin-top: 20px;
}

.basic-stat-rows {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.basic-stat-row {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 20px;
}

.stat-card {
  border-radius: 12px;
  padding: 20px 16px;
  color: #fff;
  min-height: 108px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 12px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
  word-break: break-all;
}

@media (max-width: 1200px) {
  .basic-stat-row {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .basic-stat-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .stat-value {
    font-size: 20px;
  }
}

.stat-card-gold {
  background: linear-gradient(135deg, #f6ad55, #ed8936);
}

.stat-card-recharge {
  background: linear-gradient(135deg, #63b3ed, #3182ce);
}

.stat-card-withdraw {
  background: linear-gradient(135deg, #68d391, #38a169);
}

.stat-card-register {
  background: linear-gradient(135deg, #b794f4, #805ad5);
}

.stat-card-today-recharge {
  background: linear-gradient(135deg, #4fd1c5, #319795);
}

.stat-card-today-register {
  background: linear-gradient(135deg, #f687b3, #d53f8c);
}

.stat-card-total-gold-consume {
  background: linear-gradient(135deg, #fbd38d, #dd6b20);
}

.stat-card-today-gold-consume {
  background: linear-gradient(135deg, #fc8181, #c53030);
}

.stat-card-total-diamond-consume {
  background: linear-gradient(135deg, #9f7aea, #6b46c1);
}

.stat-card-today-diamond-consume {
  background: linear-gradient(135deg, #76e4f7, #3182ce);
}
</style>

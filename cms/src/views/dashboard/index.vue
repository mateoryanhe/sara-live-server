<template>
  <div v-loading="pageLoading" class="page-container">
    <el-card>
      <el-tabs v-model="activeMainTab" lazy @tab-change="handleMainTabChange">
        <el-tab-pane :label="t('pages.dashboard.basicData')" name="basic">
          <div class="tab-toolbar">
            <el-button :loading="loading" @click="fetchSysStat">{{ t('common.refresh') }}</el-button>
          </div>
          <div class="basic-stat-board">
            <section
                v-for="section in basicStatSections"
                :key="section.key"
                class="basic-stat-section"
            >
              <div class="basic-stat-section-head">
                <span class="basic-stat-section-title">{{ section.title }}</span>
                <span class="basic-stat-section-line"/>
              </div>
              <div class="basic-stat-grid" :class="section.gridClass">
                <div
                    v-for="card in section.cards"
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
            </section>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.userData')" name="user">
          <div class="tab-toolbar">
            <el-button :loading="trendLoading" @click="fetchUserStatTrend">{{ t('common.refresh') }}</el-button>
          </div>
          <el-tabs v-model="activePeriod" lazy @tab-change="handlePeriodTabChange">
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
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.onlineUsers')" name="online">
          <div class="tab-toolbar">
            <el-button :loading="onlineLoading" @click="fetchResourceMetricTrend">{{ t('common.refresh') }}</el-button>
          </div>
          <ResourceMetricChart
              ref="onlineChartRef"
              :data="resourceMetricTrend.points"
              metric-type="online"
              :title="t('pages.dashboard.onlineChartTitle')"
              class="online-chart"
          />
        </el-tab-pane>
      </el-tabs>
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
const onlineLoading = ref(false)
const activeMainTab = ref('basic')
const activePeriod = ref('daily')
const enabledBarSeries = getUserStatBarMetricTabs()
const activeBarMetric = ref(enabledBarSeries[0]?.key ?? 'rechargeUser')

const pageLoading = computed(() => {
  if (activeMainTab.value === 'basic') {
    return loading.value
  }
  if (activeMainTab.value === 'user') {
    return trendLoading.value
  }
  return onlineLoading.value
})

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

const TODAY_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'todayRecharge', labelKey: 'statTodayRecharge', theme: 'tone-teal', format: 'amount'},
  {key: 'todayGoldConsume', labelKey: 'statTodayGoldConsume', theme: 'tone-rose', format: 'amount'},
  {key: 'todayDiamondConsume', labelKey: 'statTodayDiamondConsume', theme: 'tone-sky', format: 'amount'},
  {key: 'todayRegisterUser', labelKey: 'statTodayRegisterUser', theme: 'tone-pink', format: 'count'},
]

const TOTAL_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'totalGold', labelKey: 'statTotalGold', theme: 'tone-amber', format: 'amount'},
  {key: 'totalRecharge', labelKey: 'statTotalRecharge', theme: 'tone-blue', format: 'amount'},
  {key: 'totalWithdraw', labelKey: 'statTotalWithdraw', theme: 'tone-green', format: 'amount'},
  {key: 'totalGoldConsume', labelKey: 'statTotalGoldConsume', theme: 'tone-orange', format: 'amount'},
  {key: 'totalDiamondConsume', labelKey: 'statTotalDiamondConsume', theme: 'tone-violet', format: 'amount'},
  {key: 'totalRegisterUser', labelKey: 'statTotalRegisterUser', theme: 'tone-indigo', format: 'count'},
]

const toBasicStatCards = (configs: BasicStatCardConfig[]): BasicStatCard[] =>
    configs.map((cfg) => ({
      ...cfg,
      label: t(`pages.dashboard.${cfg.labelKey}`),
      value: sysStat[cfg.key] ?? 0,
    }))

const basicStatSections = computed(() => [
  {
    key: 'today',
    title: t('pages.dashboard.sectionToday'),
    gridClass: 'basic-stat-grid-today',
    cards: toBasicStatCards(TODAY_STAT_CARD_CONFIG),
  },
  {
    key: 'total',
    title: t('pages.dashboard.sectionTotal'),
    gridClass: 'basic-stat-grid-total',
    cards: toBasicStatCards(TOTAL_STAT_CARD_CONFIG),
  },
])

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
  onlineLoading.value = true
  try {
    const data = await sysStatApi.getResourceMetricOnlineTrend({limit: RESOURCE_METRIC_MAX_POINTS})
    resourceMetricTrend.points = data.points || []
    await nextTick()
    setTimeout(() => {
      onlineChartRef.value?.resize()
    }, 0)
  } catch (error) {
    console.error('fetch online trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchOnlineTrendFailed'))
  } finally {
    onlineLoading.value = false
  }
}

const fetchUserStatTrend = async () => {
  trendLoading.value = true
  try {
    const data = await sysStatApi.getUserStatTrend()
    userStatTrend.daily = data.daily || []
    userStatTrend.weekly = data.weekly || []
    userStatTrend.monthly = data.monthly || []
    await nextTick()
    setTimeout(() => {
      resizeUserCharts()
    }, 0)
  } catch (error) {
    console.error('fetch user stat trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchUserTrendFailed'))
  } finally {
    trendLoading.value = false
  }
}

const resizeUserCharts = () => {
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

const handleMainTabChange = async () => {
  await nextTick()
  if (activeMainTab.value === 'user') {
    resizeUserCharts()
    return
  }
  if (activeMainTab.value === 'online') {
    onlineChartRef.value?.resize()
  }
}

const handlePeriodTabChange = async () => {
  await nextTick()
  resizeUserCharts()
}

const handleBarMetricChange = async () => {
  await nextTick()
  resizeUserCharts()
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
  fetchResourceMetricTrend()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.tab-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.online-chart {
  margin-top: 0;
}

.basic-stat-board {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.basic-stat-section-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.basic-stat-section-title {
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: #64748b;
}

.basic-stat-section-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, #e2e8f0, transparent);
}

.basic-stat-grid {
  display: grid;
  gap: 14px;
}

.basic-stat-grid-today {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.basic-stat-grid-total {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.stat-card {
  position: relative;
  overflow: hidden;
  border-radius: 14px;
  padding: 18px 18px 18px 20px;
  min-height: 112px;
  background: #f8fafc;
  border: 1px solid #e8eef5;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.stat-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: var(--stat-accent, #94a3b8);
}

.stat-card:hover {
  border-color: #d7e3f0;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
  transform: translateY(-1px);
}

.stat-label {
  font-size: 13px;
  line-height: 1.4;
  color: #64748b;
  margin-bottom: 14px;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.02em;
  color: #0f172a;
  font-variant-numeric: tabular-nums;
  word-break: break-all;
}

.tone-teal {
  --stat-accent: #0d9488;
  background: linear-gradient(180deg, #f0fdfa 0%, #f8fafc 100%);
}

.tone-rose {
  --stat-accent: #e11d48;
  background: linear-gradient(180deg, #fff1f2 0%, #f8fafc 100%);
}

.tone-sky {
  --stat-accent: #0284c7;
  background: linear-gradient(180deg, #f0f9ff 0%, #f8fafc 100%);
}

.tone-pink {
  --stat-accent: #db2777;
  background: linear-gradient(180deg, #fdf2f8 0%, #f8fafc 100%);
}

.tone-amber {
  --stat-accent: #d97706;
  background: linear-gradient(180deg, #fffbeb 0%, #f8fafc 100%);
}

.tone-blue {
  --stat-accent: #2563eb;
  background: linear-gradient(180deg, #eff6ff 0%, #f8fafc 100%);
}

.tone-green {
  --stat-accent: #16a34a;
  background: linear-gradient(180deg, #f0fdf4 0%, #f8fafc 100%);
}

.tone-orange {
  --stat-accent: #ea580c;
  background: linear-gradient(180deg, #fff7ed 0%, #f8fafc 100%);
}

.tone-violet {
  --stat-accent: #7c3aed;
  background: linear-gradient(180deg, #f5f3ff 0%, #f8fafc 100%);
}

.tone-indigo {
  --stat-accent: #4f46e5;
  background: linear-gradient(180deg, #eef2ff 0%, #f8fafc 100%);
}

@media (max-width: 1200px) {
  .basic-stat-grid-today,
  .basic-stat-grid-total {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .basic-stat-grid-today,
  .basic-stat-grid-total {
    grid-template-columns: 1fr;
  }

  .stat-value {
    font-size: 22px;
  }
}
</style>

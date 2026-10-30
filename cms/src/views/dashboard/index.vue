<template>
  <div class="page-container">
    <el-card>
      <div class="page-toolbar">
        <span class="page-toolbar-label">{{ t('pages.dashboard.autoRefresh') }}</span>
        <el-select v-model="autoRefreshMs" class="auto-refresh-select">
          <el-option
              v-for="item in autoRefreshOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
          />
        </el-select>
        <el-button
            v-if="canExport"
            class="page-toolbar-export-all"
            :loading="exportAllLoading"
            @click="() => void handleExportAllDashboard()"
        >
          {{ t('pages.dashboard.exportAll') }}
        </el-button>
      </div>
      <el-tabs v-model="activeMainTab" lazy @tab-change="handleMainTabChange">
        <el-tab-pane :label="t('pages.dashboard.basicData')" name="basic">
          <div v-loading="loading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar tab-toolbar-basic">
            <el-button
                v-if="can('basicDataReset')"
                type="danger"
                :loading="basicDataResetSubmitting"
                @click="() => void handleBasicDashboardReset()"
            >
              {{ t('pages.dashboard.resetBasicData') }}
            </el-button>
            <el-button :loading="loading" @click="() => void fetchSysStat(false, true)">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <div class="basic-stat-board basic-stat-board--compact">
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
                  <div class="stat-value" :class="card.valueClass">
                    {{ card.format === 'count' ? formatCount(card.value) : formatAmount(card.value) }}
                    <span v-if="card.format === 'amountUsd'" class="stat-unit">USD</span>
                  </div>
                </div>
              </div>
            </section>
          </div>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.onlineUsers')" name="online">
          <div v-loading="onlineLoading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar tab-toolbar-trend-range tab-toolbar-online">
            <el-button :loading="onlineLoading" @click="() => void refreshOnlineTab()">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <div class="dashboard-realtime-row">
            <div class="dashboard-realtime-row-head">
              <span class="basic-stat-section-title">{{ t('pages.dashboard.sectionOnlineRealtime') }}</span>
              <span class="basic-stat-section-line"/>
            </div>
            <div class="dashboard-realtime-grid">
              <div class="stat-card tone-rose dashboard-realtime-card">
                <div class="stat-label">{{ t('pages.dashboard.onlineRealtimeAudienceCount') }}</div>
                <div class="stat-value">{{ formatCount(onlineRealtime.audienceOnlineCount) }}</div>
              </div>
              <div class="stat-card tone-blue dashboard-realtime-card">
                <div class="stat-label">{{ t('pages.dashboard.onlineRealtimeLiveAnchorCount') }}</div>
                <div class="stat-value">{{ formatCount(onlineRealtime.liveAnchorCount) }}</div>
              </div>
            </div>
          </div>
          <p class="game-bet-event-hint">{{ t('pages.dashboard.onlineTrendRetentionHint') }}</p>
          <p v-if="onlineRangeTruncated" class="game-bet-range-truncated-hint">
            {{ t('pages.dashboard.onlineTrendRangeTruncated', {max: DASHBOARD_ONLINE_TREND_POINT_LIMIT}) }}
          </p>
          <ResourceMetricChart
              ref="onlineChartRef"
              :data="resourceMetricTrend.points"
              metric-type="online"
              online-audience-anchor-only
              :title="onlineChartTitle"
              :show-trend-data-zoom="false"
              class="online-chart"
          />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.gameStats')" name="gameStats">
          <div v-loading="gameBetLoading" class="dashboard-tab-pane-body game-stats-pane">
            <div class="game-stats-toolbar">
              <div class="dashboard-realtime-row-head game-stats-toolbar-head">
                <span class="basic-stat-section-title">{{ t('pages.dashboard.sectionGameRealtime') }}</span>
                <span class="basic-stat-section-line"/>
              </div>
              <el-button :loading="gameBetLoading" @click="() => void refreshGameStatsTab()">
                {{ t('common.refresh') }}
              </el-button>
              <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
            </div>
            <div class="dashboard-realtime-grid dashboard-realtime-grid--fit">
              <div class="stat-card tone-violet dashboard-realtime-card dashboard-realtime-card--compact">
                <div class="stat-label">{{ t('pages.dashboard.gameOnlineCount') }}</div>
                <div class="stat-value">{{ formatCount(gameStatsRealtime.gameOnlineCount) }}</div>
              </div>
              <div class="stat-card tone-rose dashboard-realtime-card dashboard-realtime-card--compact">
                <div class="stat-label">{{ t('pages.dashboard.gameBetUserCount') }}</div>
                <div class="stat-value">{{ formatCount(gameStatsRealtime.gameBetUserCount) }}</div>
              </div>
            </div>
            <div class="game-stats-chart-panel">
              <div class="game-stats-chart-head">
                <div class="game-stats-chart-title">{{ gameStatChartTitle }}</div>
                <p class="game-stats-chart-hint">{{ t('pages.dashboard.onlineTrendRetentionHint') }}</p>
                <p v-if="gameStatRangeTruncated" class="game-bet-range-truncated-hint">
                  {{ t('pages.dashboard.onlineTrendRangeTruncated', {max: DASHBOARD_ONLINE_TREND_POINT_LIMIT}) }}
                </p>
              </div>
              <GameStatsTrendChart
                  ref="gameBetDailyChartRef"
                  :data="gameStatTrend.points"
                  :bet-series-name="t('pages.dashboard.gameBetUserCount')"
                  :online-series-name="t('pages.dashboard.gameOnlineCount')"
                  class="game-stats-chart"
              />
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.userEvents')" name="userEvents">
          <div v-loading="userEventsLoading" class="dashboard-tab-pane-body game-stats-pane">
            <div class="game-stats-toolbar">
              <div class="dashboard-realtime-row-head game-stats-toolbar-head">
                <span class="basic-stat-section-title">{{ t('pages.dashboard.userEvents') }}</span>
                <span class="basic-stat-section-line"/>
              </div>
              <el-button :loading="userEventsLoading" @click="() => void refreshUserEventsTab()">
                {{ t('common.refresh') }}
              </el-button>
              <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
            </div>
            <div class="game-stats-chart-panel user-events-chart-panel">
              <div class="game-stats-chart-head">
                <div class="game-stats-chart-title">{{ userEventsChartTitle }}</div>
                <p class="game-stats-chart-hint">{{ t('pages.trackingMinuteLiveActivity.chartHint') }}</p>
                <p v-if="userEventsRangeTruncated" class="game-bet-range-truncated-hint">
                  {{ t('pages.trackingMinuteLiveActivity.rangeTruncated', {max: userEventsTrendMaxPoints}) }}
                </p>
              </div>
              <el-tabs
                  v-model="activeUserEventMetric"
                  class="user-events-metric-tabs"
                  lazy
                  @tab-change="handleUserEventMetricTabChange"
              >
                <el-tab-pane
                    v-for="tab in userEventSeriesTabs"
                    :key="tab.key"
                    :label="tab.label"
                    :name="tab.key"
                >
                  <MinuteLiveActivityTrendChart
                      ref="userEventsChartRef"
                      :data="userEventsTrend.points"
                      :series-key="tab.key"
                      :title="tab.chartTitle"
                      class="game-stats-chart"
                  />
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.userData')" name="user">
          <div v-loading="trendLoading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar">
            <el-button :loading="trendLoading" @click="() => void fetchUserStatTrend(false, true)">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <el-tabs v-model="activePeriod" lazy @tab-change="handlePeriodTabChange">
            <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
              <UserStatChart ref="dailyLineChartRef" :data="userStatTrend.daily" :title="t('pages.dashboard.activeRegisterDaily')"/>
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
              <UserStatChart ref="weeklyLineChartRef" :data="userStatTrend.weekly" :title="t('pages.dashboard.activeRegisterWeekly')"/>
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
              <UserStatChart ref="monthlyLineChartRef" :data="userStatTrend.monthly" :title="t('pages.dashboard.activeRegisterMonthly')"/>
            </el-tab-pane>
          </el-tabs>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.barMetricRechargeUser')" name="rechargeUser">
          <div v-loading="trendLoading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar">
            <el-button :loading="trendLoading" @click="() => void fetchUserStatTrend(false, true)">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <el-tabs v-model="activeRechargePeriod" lazy @tab-change="handleRechargePeriodTabChange">
            <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
              <UserStatBarChart
                  ref="rechargeDailyBarChartRef"
                  :data="userStatTrend.daily"
                  metric-key="rechargeUser"
                  :title="rechargeBarChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
              <UserStatBarChart
                  ref="rechargeWeeklyBarChartRef"
                  :data="userStatTrend.weekly"
                  metric-key="rechargeUser"
                  :title="rechargeBarChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
              <UserStatBarChart
                  ref="rechargeMonthlyBarChartRef"
                  :data="userStatTrend.monthly"
                  metric-key="rechargeUser"
                  :title="rechargeBarChartTitle"
              />
            </el-tab-pane>
          </el-tabs>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.clientRecharge')" name="clientRecharge">
          <div v-loading="trendLoading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar">
            <el-button :loading="trendLoading" @click="() => void fetchUserStatTrend(false, true)">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <p class="game-bet-event-hint">{{ t('pages.dashboard.clientRechargeHint') }}</p>
          <el-tabs v-model="activeClientRechargePeriod" lazy @tab-change="handleClientRechargePeriodTabChange">
            <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
              <ClientRechargeTrendChart
                  ref="clientRechargeDailyChartRef"
                  :data="userStatTrend.daily"
                  :title="clientRechargeChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
              <ClientRechargeTrendChart
                  ref="clientRechargeWeeklyChartRef"
                  :data="userStatTrend.weekly"
                  :title="clientRechargeChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
              <ClientRechargeTrendChart
                  ref="clientRechargeMonthlyChartRef"
                  :data="userStatTrend.monthly"
                  :title="clientRechargeChartTitle"
              />
            </el-tab-pane>
          </el-tabs>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('pages.dashboard.clientLogin')" name="clientLogin">
          <div v-loading="trendLoading" class="dashboard-tab-pane-body">
          <div class="tab-toolbar">
            <el-button :loading="trendLoading" @click="() => void fetchUserStatTrend(false, true)">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExportDashboard">{{ t('common.export') }}</el-button>
          </div>
          <p class="game-bet-event-hint">{{ t('pages.dashboard.clientLoginHint') }}</p>
          <el-tabs v-model="activeClientLoginPeriod" lazy @tab-change="handleClientLoginPeriodTabChange">
            <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
              <ClientLoginTrendChart
                  ref="clientLoginDailyChartRef"
                  :data="userStatTrend.daily"
                  :title="clientLoginChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
              <ClientLoginTrendChart
                  ref="clientLoginWeeklyChartRef"
                  :data="userStatTrend.weekly"
                  :title="clientLoginChartTitle"
              />
            </el-tab-pane>
            <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
              <ClientLoginTrendChart
                  ref="clientLoginMonthlyChartRef"
                  :data="userStatTrend.monthly"
                  :title="clientLoginChartTitle"
              />
            </el-tab-pane>
          </el-tabs>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {
  computed,
  defineAsyncComponent,
  nextTick,
  onActivated,
  onDeactivated,
  onMounted,
  onUnmounted,
  reactive,
  ref,
} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {sysStatApi} from '@/api'
import {trackingEventApi} from '@/api/modules/tracking-event'
import {usePagePermission} from '@/composables/usePagePermission'
import type {
  GameStatTrend,
  MinuteLiveActivityTrendRes,
  ResourceMetricTrend,
  SysStat,
  UserStatTrend,
} from '@/types/api'
import {
  buildDashboardResourceMetricTrendQuery,
  DASHBOARD_ONLINE_TREND_POINT_LIMIT,
  DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS,
} from './dashboard-trend-range'
import {formatAmount} from '@/utils/number-format'
import {
  buildDashboardBasicCsvContent,
  buildDashboardClientLoginCsvContent,
  buildDashboardClientRechargeCsvContent,
  buildDashboardExportFilename,
  buildDashboardExportStamp,
  buildDashboardGameStatsCsvContent,
  buildDashboardOnlineCsvContent,
  buildDashboardRechargeUserBarCsvContent,
  buildDashboardUserEventsCsvContent,
  buildDashboardUserTrendCsvContent,
  downloadDashboardCsvBundle,
  exportDashboardBasicCsv,
  exportDashboardClientLoginCsv,
  exportDashboardClientRechargeCsv,
  exportDashboardGameStatsCsv,
  exportDashboardOnlineCsv,
  exportDashboardRechargeUserBarCsv,
  exportDashboardUserEventsCsv,
  exportDashboardUserTrendCsv,
  type DashboardExportLabels,
} from '@/utils/dashboard-export'
import {
  buildDashboardTrackingExportCsvFiles,
  fetchAllTrackingTrendsForDashboardExport,
} from '@/utils/dashboard-tracking-export'
import {
  MINUTE_LIVE_ACTIVITY_SERIES_KEYS,
  type MinuteLiveActivitySeriesKey,
} from '@/views/tracking/minute-live-activity-series'

const UserStatChart = defineAsyncComponent(() => import('./components/user-stat-chart.vue'))
const UserStatBarChart = defineAsyncComponent(() => import('./components/user-stat-bar-chart.vue'))
const ClientRechargeTrendChart = defineAsyncComponent(() => import('./components/client-recharge-trend-chart.vue'))
const ClientLoginTrendChart = defineAsyncComponent(() => import('./components/client-login-trend-chart.vue'))
const ResourceMetricChart = defineAsyncComponent(() => import('@/views/config/components/resource-metric-chart.vue'))
const GameStatsTrendChart = defineAsyncComponent(() => import('./components/game-stats-trend-chart.vue'))
const MinuteLiveActivityTrendChart = defineAsyncComponent(
    () => import('@/views/tracking/components/minute-live-activity-trend-chart.vue'),
)

/** 3 天 × 1440 分钟，与后端 retention / maxTrendPoints 一致 */
const USER_EVENTS_TREND_MAX_POINTS = 4320

/** 与后端 DashboardCacheTTL 一致，实际自动刷新间隔固定 5 分钟（下拉仅展示，改选项不影响刷新） */
const DASHBOARD_AUTO_REFRESH_MS = 300_000

const AUTO_REFRESH_INTERVALS = [
  {value: 0, labelKey: 'autoRefreshOff'},
  {value: 10_000, labelKey: 'autoRefresh10s'},
  {value: 30_000, labelKey: 'autoRefresh30s'},
  {value: 60_000, labelKey: 'autoRefresh60s'},
  {value: 300_000, labelKey: 'autoRefresh5m'},
] as const

const {t, locale} = useI18n()
const {can} = usePagePermission('Dashboard')
const canExport = computed(() => can('export'))
const exportAllLoading = ref(false)
const basicDataResetSubmitting = ref(false)
const loading = ref(false)
const trendLoading = ref(false)
const onlineLoading = ref(false)
const gameBetLoading = ref(false)
const userEventsLoading = ref(false)
const activeMainTab = ref('basic')
const activePeriod = ref('daily')
const activeRechargePeriod = ref('daily')
const activeClientRechargePeriod = ref('daily')
const activeClientLoginPeriod = ref('daily')
const onlineRangeTruncated = ref(false)
const onlineTrendMaxPoints = ref(0)
const gameStatRangeTruncated = ref(false)
const gameStatTrendMaxPoints = ref(0)
const userEventsRangeTruncated = ref(false)
const userEventsTrendMaxPoints = ref(USER_EVENTS_TREND_MAX_POINTS)
const activeUserEventMetric = ref<MinuteLiveActivitySeriesKey>('joinCount')
/** 下拉展示用；定时器始终 DASHBOARD_AUTO_REFRESH_MS */
const autoRefreshMs = ref(DASHBOARD_AUTO_REFRESH_MS)
const dashboardRouteActive = ref(false)
let autoRefreshTimer: ReturnType<typeof setInterval> | null = null
let dashboardRefreshInFlight = false
let skipNextActivatedRefresh = false

/** 各主 Tab 是否已拉过数据（切换 Tab 时懒加载，避免首屏 4 连请求） */
const tabDataLoaded = reactive({
  sysStat: false,
  userTrend: false,
  gameTrend: false,
  userEventsTrend: false,
  onlineTrend: false,
})

const isDashboardVisible = () =>
    dashboardRouteActive.value && typeof document !== 'undefined' && document.visibilityState === 'visible'

const autoRefreshOptions = computed(() =>
    AUTO_REFRESH_INTERVALS.map((item) => ({
      value: item.value,
      label: t(`pages.dashboard.${item.labelKey}`),
    })),
)

const rechargeBarChartPeriodSuffix = computed(() => {
  if (activeRechargePeriod.value === 'weekly') {
    return t('pages.dashboard.last12Weeks')
  }
  if (activeRechargePeriod.value === 'monthly') {
    return t('pages.dashboard.last12Months')
  }
  return t('pages.dashboard.last30Days')
})

const rechargeBarChartTitle = computed(() =>
    t('pages.dashboard.chartTitleWithPeriod', {
      label: t('pages.dashboard.barMetricRechargeUser'),
      period: rechargeBarChartPeriodSuffix.value,
    }),
)

const clientRechargeChartPeriodSuffix = computed(() => {
  if (activeClientRechargePeriod.value === 'weekly') {
    return t('pages.dashboard.last12Weeks')
  }
  if (activeClientRechargePeriod.value === 'monthly') {
    return t('pages.dashboard.last12Months')
  }
  return t('pages.dashboard.last30Days')
})

const clientRechargeChartTitle = computed(() =>
    t('pages.dashboard.chartTitleWithPeriod', {
      label: t('pages.dashboard.clientRecharge'),
      period: clientRechargeChartPeriodSuffix.value,
    }),
)

const clientLoginChartPeriodSuffix = computed(() => {
  if (activeClientLoginPeriod.value === 'weekly') {
    return t('pages.dashboard.last12Weeks')
  }
  if (activeClientLoginPeriod.value === 'monthly') {
    return t('pages.dashboard.last12Months')
  }
  return t('pages.dashboard.last30Days')
})

const clientLoginChartTitle = computed(() =>
    t('pages.dashboard.chartTitleWithPeriod', {
      label: t('pages.dashboard.clientLogin'),
      period: clientLoginChartPeriodSuffix.value,
    }),
)

const onlineChartTitle = computed(() => t('pages.dashboard.onlineLiveChartTitleLast3Days'))

const gameStatChartTitle = computed(() => t('pages.dashboard.gameStatChartTitleLast3Days'))
const userEventsChartTitle = computed(() => t('pages.trackingMinuteLiveActivity.chartTitleLast3Days'))

const userEventSeriesLabel = (key: MinuteLiveActivitySeriesKey) => {
  switch (key) {
    case 'joinCount':
      return t('pages.trackingMinuteLiveActivity.join')
    case 'leaveCount':
      return t('pages.trackingMinuteLiveActivity.leave')
    case 'giftCount':
      return t('pages.trackingMinuteLiveActivity.gift')
    case 'chatCount':
      return t('pages.trackingMinuteLiveActivity.chat')
    case 'betCount':
      return t('pages.trackingMinuteLiveActivity.bet')
    case 'call1v1Count':
      return t('pages.trackingMinuteLiveActivity.call1v1')
    case 'rechargeOrderCount':
      return t('pages.trackingMinuteLiveActivity.rechargeOrder')
    default:
      return key
  }
}

const userEventSeriesTabs = computed(() =>
    MINUTE_LIVE_ACTIVITY_SERIES_KEYS.map((key) => {
      const label = userEventSeriesLabel(key)
      return {
        key,
        label,
        chartTitle: `${userEventsChartTitle.value} · ${label}`,
      }
    }),
)

const dailyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const weeklyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const monthlyLineChartRef = ref<InstanceType<typeof UserStatChart>>()
const rechargeDailyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const rechargeWeeklyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const rechargeMonthlyBarChartRef = ref<InstanceType<typeof UserStatBarChart>>()
const clientRechargeDailyChartRef = ref<InstanceType<typeof ClientRechargeTrendChart>>()
const clientRechargeWeeklyChartRef = ref<InstanceType<typeof ClientRechargeTrendChart>>()
const clientRechargeMonthlyChartRef = ref<InstanceType<typeof ClientRechargeTrendChart>>()
const clientLoginDailyChartRef = ref<InstanceType<typeof ClientLoginTrendChart>>()
const clientLoginWeeklyChartRef = ref<InstanceType<typeof ClientLoginTrendChart>>()
const clientLoginMonthlyChartRef = ref<InstanceType<typeof ClientLoginTrendChart>>()
const onlineChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const gameBetDailyChartRef = ref<InstanceType<typeof GameStatsTrendChart>>()
const userEventsChartRef = ref<InstanceType<typeof MinuteLiveActivityTrendChart>>()

type BasicTabSysStat = Omit<
    SysStat,
    'onlineCount' | 'gameOnlineCount' | 'gameBetUserCount' | 'audienceOnlineCount' | 'liveAnchorCount'
>

const basicRealtime = reactive({
  audienceOnlineCount: 0,
  liveAnchorCount: 0,
  gameOnlineCount: 0,
  gameBetUserCount: 0,
})

const sysStat = reactive<BasicTabSysStat>({
  totalGold: 0,
  totalGoldConsume: 0,
  totalDiamondConsume: 0,
  totalRecharge: 0,
  totalNormalUserRecharge: 0,
  totalCoinMerchantRecharge: 0,
  totalAnchorPayout: 0,
  weeklyAnchorPayout: 0,
  totalVirtualRecharge: 0,
  totalWithdraw: 0,
  totalRegisterUser: 0,
  todayRecharge: 0,
  todayGoldConsume: 0,
  todayDiamondConsume: 0,
  todayRegisterUser: 0,
})

const onlineRealtime = reactive({
  audienceOnlineCount: 0,
  liveAnchorCount: 0,
})

const gameStatsRealtime = reactive({
  gameOnlineCount: 0,
  gameBetUserCount: 0,
})

const resourceMetricTrend = reactive<ResourceMetricTrend>({
  points: [],
})

const userStatTrend = reactive<UserStatTrend>({
  daily: [],
  weekly: [],
  monthly: [],
})

const gameStatTrend = reactive<GameStatTrend>({
  points: [],
})

const userEventsTrend = reactive<MinuteLiveActivityTrendRes>({
  points: [],
})

type BasicStatCardKey = keyof Pick<
    SysStat,
    | 'totalGold'
    | 'todayRecharge'
    | 'todayGoldConsume'
    | 'todayDiamondConsume'
    | 'todayRegisterUser'
    | 'totalRecharge'
    | 'totalNormalUserRecharge'
    | 'totalCoinMerchantRecharge'
    | 'totalVirtualRecharge'
    | 'totalWithdraw'
    | 'totalGoldConsume'
    | 'totalDiamondConsume'
    | 'totalRegisterUser'
    | 'weeklyAnchorPayout'
    | 'totalAnchorPayout'
>

type BasicStatCardConfig = {
  key: BasicStatCardKey
  labelKey: string
  theme: string
  format: 'amount' | 'amountUsd' | 'count'
}

type BasicStatCard = BasicStatCardConfig & {
  label: string
  value: number | string
  valueClass?: string
}

const TODAY_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'todayRecharge', labelKey: 'statTodayRecharge', theme: 'tone-teal', format: 'amount'},
  {key: 'todayGoldConsume', labelKey: 'statTodayGoldConsume', theme: 'tone-rose', format: 'amount'},
  {key: 'todayDiamondConsume', labelKey: 'statTodayDiamondConsume', theme: 'tone-sky', format: 'amount'},
  {key: 'todayRegisterUser', labelKey: 'statTodayRegisterUser', theme: 'tone-pink', format: 'count'},
]

const FINANCIAL_TOTAL_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'totalRecharge', labelKey: 'statTotalRecharge', theme: 'tone-blue', format: 'amount'},
  {key: 'totalNormalUserRecharge', labelKey: 'statTotalNormalUserRecharge', theme: 'tone-teal', format: 'amount'},
  {key: 'totalCoinMerchantRecharge', labelKey: 'statTotalCoinMerchantRecharge', theme: 'tone-indigo', format: 'amount'},
  {key: 'totalVirtualRecharge', labelKey: 'statTotalVirtualRecharge', theme: 'tone-sky', format: 'amount'},
  {key: 'totalWithdraw', labelKey: 'statTotalWithdraw', theme: 'tone-green', format: 'amount'},
]

const BUSINESS_TOTAL_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'totalGold', labelKey: 'statTotalGold', theme: 'tone-amber', format: 'amount'},
  {key: 'totalGoldConsume', labelKey: 'statTotalGoldConsume', theme: 'tone-orange', format: 'amount'},
  {key: 'totalDiamondConsume', labelKey: 'statTotalDiamondConsume', theme: 'tone-violet', format: 'amount'},
  {key: 'totalRegisterUser', labelKey: 'statTotalRegisterUser', theme: 'tone-indigo', format: 'count'},
]

const ANCHOR_PAYOUT_STAT_CARD_CONFIG: BasicStatCardConfig[] = [
  {key: 'weeklyAnchorPayout', labelKey: 'statWeeklyAnchorPayout', theme: 'tone-teal', format: 'amountUsd'},
  {key: 'totalAnchorPayout', labelKey: 'statTotalAnchorPayout', theme: 'tone-blue', format: 'amountUsd'},
]

type BasicRealtimeCardConfig = {
  key: keyof typeof basicRealtime
  labelKey: string
  theme: string
  valueClass?: string
}

const BASIC_REALTIME_STAT_CARD_CONFIG: BasicRealtimeCardConfig[] = [
  {key: 'audienceOnlineCount', labelKey: 'onlineRealtimeAudienceCount', theme: 'tone-neutral', valueClass: 'stat-value--red'},
  {key: 'liveAnchorCount', labelKey: 'onlineRealtimeLiveAnchorCount', theme: 'tone-neutral', valueClass: 'stat-value--blue'},
  {key: 'gameOnlineCount', labelKey: 'gameOnlineCount', theme: 'tone-neutral', valueClass: 'stat-value--violet'},
  {key: 'gameBetUserCount', labelKey: 'gameBetUserCount', theme: 'tone-neutral', valueClass: 'stat-value--orange'},
]

const toBasicStatCards = (configs: BasicStatCardConfig[]): BasicStatCard[] =>
    configs.map((cfg) => ({
      ...cfg,
      label: t(`pages.dashboard.${cfg.labelKey}`),
      value: sysStat[cfg.key] ?? 0,
    }))

const toBasicRealtimeCards = () =>
    BASIC_REALTIME_STAT_CARD_CONFIG.map((cfg) => ({
      key: cfg.key,
      labelKey: cfg.labelKey,
      theme: cfg.theme,
      valueClass: cfg.valueClass,
      label: t(`pages.dashboard.${cfg.labelKey}`),
      value: basicRealtime[cfg.key] ?? 0,
      format: 'count' as const,
    }))

const basicStatSections = computed(() => [
  {
    key: 'realtimePeople',
    title: t('pages.dashboard.sectionRealtimePeople'),
    gridClass: 'basic-stat-grid-realtime',
    cards: toBasicRealtimeCards(),
  },
  {
    key: 'today',
    title: t('pages.dashboard.sectionToday'),
    gridClass: 'basic-stat-grid-today',
    cards: toBasicStatCards(TODAY_STAT_CARD_CONFIG),
  },
  {
    key: 'financialTotal',
    title: t('pages.dashboard.sectionFinancialTotal'),
    gridClass: 'basic-stat-grid-total',
    cards: toBasicStatCards(FINANCIAL_TOTAL_STAT_CARD_CONFIG),
  },
  {
    key: 'businessTotal',
    title: t('pages.dashboard.sectionBusinessTotal'),
    gridClass: 'basic-stat-grid-total',
    cards: toBasicStatCards(BUSINESS_TOTAL_STAT_CARD_CONFIG),
  },
  {
    key: 'anchorPayout',
    title: t('pages.dashboard.anchorPayout'),
    gridClass: 'basic-stat-grid-anchor-payout',
    cards: toBasicStatCards(ANCHOR_PAYOUT_STAT_CARD_CONFIG),
  },
])

const dashboardExportLabels = computed((): DashboardExportLabels => ({
  section: t('pages.dashboard.exportCsvSection'),
  metric: t('pages.dashboard.exportCsvMetric'),
  value: t('pages.dashboard.exportCsvValue'),
  time: t('pages.dashboard.exportCsvTime'),
  audienceOnline: t('pages.dashboard.onlineRealtimeAudienceCount'),
  liveAnchorOnline: t('pages.dashboard.onlineRealtimeLiveAnchorCount'),
  gameOnline: t('pages.dashboard.gameOnlineCount'),
  gameBetUser: t('pages.dashboard.gameBetUserCount'),
  activeUsers: t('pages.dashboard.activeUsers'),
  registerUsers: t('pages.dashboard.registerUsers'),
  rechargeUser: t('pages.dashboard.barMetricRechargeUser'),
  appRecharge: t('pages.dashboard.clientRechargeApp'),
  h5Recharge: t('pages.dashboard.clientRechargeH5'),
  appLogin: t('pages.dashboard.clientLoginApp'),
  h5Login: t('pages.dashboard.clientLoginH5'),
  join: t('pages.trackingMinuteLiveActivity.join'),
  leave: t('pages.trackingMinuteLiveActivity.leave'),
  gift: t('pages.trackingMinuteLiveActivity.gift'),
  chat: t('pages.trackingMinuteLiveActivity.chat'),
  bet: t('pages.trackingMinuteLiveActivity.bet'),
  call1v1: t('pages.trackingMinuteLiveActivity.call1v1'),
  rechargeOrder: t('pages.trackingMinuteLiveActivity.rechargeOrder'),
}))

const pickUserTrendPoints = (period: string) => {
  if (period === 'weekly') {
    return userStatTrend.weekly
  }
  if (period === 'monthly') {
    return userStatTrend.monthly
  }
  return userStatTrend.daily
}

const handleExportDashboard = () => {
  const labels = dashboardExportLabels.value
  let ok = false
  const tab = activeMainTab.value
  if (tab === 'basic') {
    ok = exportDashboardBasicCsv(
        buildDashboardExportFilename('basic'),
        basicStatSections.value.map((section) => ({
          title: section.title,
          cards: section.cards.map((card) => ({label: card.label, value: card.value})),
        })),
        labels,
    )
  } else if (tab === 'online') {
    ok = exportDashboardOnlineCsv(
        buildDashboardExportFilename('online'),
        {...onlineRealtime},
        resourceMetricTrend.points,
        labels,
    )
  } else if (tab === 'gameStats') {
    ok = exportDashboardGameStatsCsv(
        buildDashboardExportFilename('game-stats'),
        {...gameStatsRealtime},
        gameStatTrend.points,
        labels,
    )
  } else if (tab === 'userEvents') {
    ok = exportDashboardUserEventsCsv(
        buildDashboardExportFilename('user-events'),
        userEventsTrend.points,
        labels,
    )
  } else if (tab === 'user') {
    ok = exportDashboardUserTrendCsv(
        buildDashboardExportFilename('user', activePeriod.value),
        pickUserTrendPoints(activePeriod.value),
        labels,
    )
  } else if (tab === 'rechargeUser') {
    ok = exportDashboardRechargeUserBarCsv(
        buildDashboardExportFilename('recharge-user', activeRechargePeriod.value),
        pickUserTrendPoints(activeRechargePeriod.value),
        'rechargeUser',
        labels,
    )
  } else if (tab === 'clientRecharge') {
    ok = exportDashboardClientRechargeCsv(
        buildDashboardExportFilename('client-recharge', activeClientRechargePeriod.value),
        pickUserTrendPoints(activeClientRechargePeriod.value),
        labels,
    )
  } else if (tab === 'clientLogin') {
    ok = exportDashboardClientLoginCsv(
        buildDashboardExportFilename('client-login', activeClientLoginPeriod.value),
        pickUserTrendPoints(activeClientLoginPeriod.value),
        labels,
    )
  }
  if (ok) {
    ElMessage.success(t('common.exportSuccess'))
  } else {
    ElMessage.warning(t('common.exportEmpty'))
  }
}

const fetchAllDashboardDataForExport = async (): Promise<boolean> => {
  const results = await Promise.all([
    fetchSysStat(true, true),
    fetchUserStatTrend(true, true),
    fetchGameStatTrend(true, true),
    fetchGameStatsRealtime(true, true),
    fetchUserEventsTrend(true, true),
    fetchResourceMetricTrend(true, true),
    fetchOnlineRealtime(true, true),
  ])
  tabDataLoaded.sysStat = results[0] || tabDataLoaded.sysStat
  tabDataLoaded.userTrend = results[1] || tabDataLoaded.userTrend
  tabDataLoaded.gameTrend = results[2] || tabDataLoaded.gameTrend
  tabDataLoaded.userEventsTrend = results[4] || tabDataLoaded.userEventsTrend
  tabDataLoaded.onlineTrend = results[5] || tabDataLoaded.onlineTrend
  return results.some(Boolean)
}

const handleExportAllDashboard = async () => {
  if (exportAllLoading.value) {
    return
  }
  exportAllLoading.value = true
  try {
    const [, trackingTrends] = await Promise.all([
      fetchAllDashboardDataForExport(),
      fetchAllTrackingTrendsForDashboardExport(),
    ])
    const labels = dashboardExportLabels.value
    const stamp = buildDashboardExportStamp()
    const basicSections = basicStatSections.value.map((section) => ({
      title: section.title,
      cards: section.cards.map((card) => ({label: card.label, value: card.value})),
    }))
    const periods = ['daily', 'weekly', 'monthly'] as const
    const trendByPeriod = {
      daily: userStatTrend.daily,
      weekly: userStatTrend.weekly,
      monthly: userStatTrend.monthly,
    }
    const files = [
      {filename: buildDashboardExportFilename('basic', undefined, stamp), content: buildDashboardBasicCsvContent(basicSections, labels)},
      {filename: buildDashboardExportFilename('online', undefined, stamp), content: buildDashboardOnlineCsvContent({...onlineRealtime}, resourceMetricTrend.points, labels)},
      {filename: buildDashboardExportFilename('game-stats', undefined, stamp), content: buildDashboardGameStatsCsvContent({...gameStatsRealtime}, gameStatTrend.points, labels)},
      {filename: buildDashboardExportFilename('user-events', undefined, stamp), content: buildDashboardUserEventsCsvContent(userEventsTrend.points, labels)},
      ...periods.map((period) => ({
        filename: buildDashboardExportFilename('user', period, stamp),
        content: buildDashboardUserTrendCsvContent(trendByPeriod[period], labels),
      })),
      ...periods.map((period) => ({
        filename: buildDashboardExportFilename('recharge-user', period, stamp),
        content: buildDashboardRechargeUserBarCsvContent(trendByPeriod[period], 'rechargeUser', labels),
      })),
      ...periods.map((period) => ({
        filename: buildDashboardExportFilename('client-recharge', period, stamp),
        content: buildDashboardClientRechargeCsvContent(trendByPeriod[period], labels),
      })),
      ...periods.map((period) => ({
        filename: buildDashboardExportFilename('client-login', period, stamp),
        content: buildDashboardClientLoginCsvContent(trendByPeriod[period], labels),
      })),
      ...buildDashboardTrackingExportCsvFiles(
          t,
          stamp,
          t('pages.dashboard.exportCsvTime'),
          trackingTrends,
      ),
    ]
    const count = await downloadDashboardCsvBundle(files)
    if (count > 0) {
      ElMessage.success(t('pages.dashboard.exportAllSuccess', {count}))
    } else {
      ElMessage.warning(t('common.exportEmpty'))
    }
  } catch (error) {
    console.error('export all dashboard data failed:', error)
    ElMessage.error(t('pages.dashboard.exportAllFailed'))
  } finally {
    exportAllLoading.value = false
  }
}

const clearAutoRefreshTimer = () => {
  if (autoRefreshTimer != null) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

let visibilityListenerBound = false

const bindDocumentVisibilityListener = () => {
  if (visibilityListenerBound || typeof document === 'undefined') {
    return
  }
  document.addEventListener('visibilitychange', handleDocumentVisibilityChange)
  visibilityListenerBound = true
}

const unbindDocumentVisibilityListener = () => {
  if (!visibilityListenerBound || typeof document === 'undefined') {
    return
  }
  document.removeEventListener('visibilitychange', handleDocumentVisibilityChange)
  visibilityListenerBound = false
}

const stopDashboardAutoRefresh = () => {
  clearAutoRefreshTimer()
  unbindDocumentVisibilityListener()
}

const restartAutoRefreshTimer = () => {
  clearAutoRefreshTimer()
  if (!isDashboardVisible()) {
    return
  }
  autoRefreshTimer = setInterval(() => {
    if (!isDashboardVisible()) {
      clearAutoRefreshTimer()
      return
    }
    void refreshActiveMainTab(true, true)
  }, DASHBOARD_AUTO_REFRESH_MS)
}

const handleDocumentVisibilityChange = () => {
  if (isDashboardVisible()) {
    void refreshActiveMainTab(true, true)
    restartAutoRefreshTimer()
    return
  }
  // 浏览器标签隐藏 / 最小化：立刻停掉定时刷新
  clearAutoRefreshTimer()
}

const refreshActiveMainTab = async (silent = false, force = false) => {
  // 路由已切走（keep-alive 失活）或浏览器页不可见时，不发起任何刷新
  if (!isDashboardVisible()) {
    clearAutoRefreshTimer()
    return
  }
  if (dashboardRefreshInFlight) {
    return
  }
  dashboardRefreshInFlight = true
  try {
    const tab = activeMainTab.value
    if (tab === 'basic') {
      if (force || !tabDataLoaded.sysStat) {
        const ok = await fetchSysStat(silent)
        if (ok) {
          tabDataLoaded.sysStat = true
        }
      }
      return
    }
    if (tab === 'user' || tab === 'rechargeUser' || tab === 'clientRecharge' || tab === 'clientLogin') {
      if (force || !tabDataLoaded.userTrend) {
        const ok = await fetchUserStatTrend(silent)
        if (ok) {
          tabDataLoaded.userTrend = true
        }
      }
      return
    }
    if (tab === 'gameStats') {
      const needTrend = force || !tabDataLoaded.gameTrend
      const results = await Promise.all([
        needTrend ? fetchGameStatTrend(silent) : Promise.resolve(true),
        fetchGameStatsRealtime(silent),
      ])
      if (needTrend && results[0]) {
        tabDataLoaded.gameTrend = true
      }
      return
    }
    if (tab === 'userEvents') {
      if (force || !tabDataLoaded.userEventsTrend) {
        const ok = await fetchUserEventsTrend(silent)
        if (ok) {
          tabDataLoaded.userEventsTrend = true
        }
      }
      return
    }
    if (tab === 'online') {
      const needTrend = force || !tabDataLoaded.onlineTrend
      const results = await Promise.all([
        needTrend ? fetchResourceMetricTrend(silent) : Promise.resolve(true),
        fetchOnlineRealtime(silent),
      ])
      if (needTrend && results[0]) {
        tabDataLoaded.onlineTrend = true
      }
    }
  } finally {
    dashboardRefreshInFlight = false
  }
}

const handleBasicDashboardReset = async () => {
  if (basicDataResetSubmitting.value) {
    return
  }
  const boxOpts = {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning' as const,
  }
  try {
    await ElMessageBox.confirm(t('pages.dashboard.resetBasicDataConfirm1'), t('pages.dashboard.resetBasicDataTitle'), boxOpts)
    await ElMessageBox.confirm(t('pages.dashboard.resetBasicDataConfirm2'), t('pages.dashboard.resetBasicDataTitleFinal'), boxOpts)
  } catch {
    return
  }
  basicDataResetSubmitting.value = true
  try {
    await sysStatApi.resetBasicDashboardStats()
    ElMessage.success(t('pages.dashboard.resetBasicDataSuccess'))
    await fetchSysStat(true, true)
  } catch (error) {
    console.error('reset basic dashboard stats failed:', error)
    ElMessage.error(t('pages.dashboard.resetBasicDataFailed'))
  } finally {
    basicDataResetSubmitting.value = false
  }
}

const fetchSysStat = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  if (!silent) {
    loading.value = true
  }
  try {
    const data = await sysStatApi.getSysStat({ bypassDashboardCache })
    sysStat.totalGold = data.totalGold ?? 0
    sysStat.totalGoldConsume = data.totalGoldConsume ?? 0
    sysStat.totalDiamondConsume = data.totalDiamondConsume ?? 0
    sysStat.totalRecharge = data.totalRecharge ?? 0
    sysStat.totalNormalUserRecharge = data.totalNormalUserRecharge ?? 0
    sysStat.totalCoinMerchantRecharge = data.totalCoinMerchantRecharge ?? 0
    sysStat.totalAnchorPayout = data.totalAnchorPayout ?? 0
    sysStat.weeklyAnchorPayout = data.weeklyAnchorPayout ?? 0
    sysStat.totalVirtualRecharge = data.totalVirtualRecharge ?? 0
    sysStat.totalWithdraw = data.totalWithdraw ?? 0
    sysStat.totalRegisterUser = data.totalRegisterUser ?? 0
    sysStat.todayRecharge = data.todayRecharge ?? 0
    sysStat.todayGoldConsume = data.todayGoldConsume ?? 0
    sysStat.todayDiamondConsume = data.todayDiamondConsume ?? 0
    sysStat.todayRegisterUser = data.todayRegisterUser ?? 0
    basicRealtime.audienceOnlineCount = Number(data.audienceOnlineCount ?? 0)
    basicRealtime.liveAnchorCount = Number(data.liveAnchorCount ?? 0)
    basicRealtime.gameOnlineCount = Number(data.gameOnlineCount ?? 0)
    basicRealtime.gameBetUserCount = Number(data.gameBetUserCount ?? 0)
    return true
  } catch (error) {
    console.error('fetch sys stat failed:', error)
    ElMessage.error(t('pages.dashboard.fetchSysStatFailed'))
    return false
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

const fetchOnlineRealtime = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  try {
    const data = await sysStatApi.getSysStat({ bypassDashboardCache })
    onlineRealtime.audienceOnlineCount = Number(data.audienceOnlineCount ?? 0)
    onlineRealtime.liveAnchorCount = Number(data.liveAnchorCount ?? 0)
    return true
  } catch (error) {
    console.error('fetch online realtime failed:', error)
    if (!silent) {
      ElMessage.error(t('pages.dashboard.fetchSysStatFailed'))
    }
    return false
  }
}

const refreshOnlineTab = async () => {
  onlineLoading.value = true
  try {
    await Promise.all([
      fetchResourceMetricTrend(true, true),
      fetchOnlineRealtime(true, true),
    ])
    tabDataLoaded.onlineTrend = true
  } finally {
    onlineLoading.value = false
  }
}

const fetchGameStatsRealtime = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  try {
    const data = await sysStatApi.getSysStat({ bypassDashboardCache })
    gameStatsRealtime.gameOnlineCount = Number(data.gameOnlineCount ?? 0)
    gameStatsRealtime.gameBetUserCount = Number(data.gameBetUserCount ?? 0)
    return true
  } catch (error) {
    console.error('fetch game stats realtime failed:', error)
    if (!silent) {
      ElMessage.error(t('pages.dashboard.fetchSysStatFailed'))
    }
    return false
  }
}

const refreshGameStatsTab = async () => {
  gameBetLoading.value = true
  try {
    await Promise.all([
      fetchGameStatTrend(true, true),
      fetchGameStatsRealtime(true, true),
    ])
    tabDataLoaded.gameTrend = true
  } finally {
    gameBetLoading.value = false
  }
}

const refreshUserEventsTab = async () => {
  const ok = await fetchUserEventsTrend(false, true)
  if (ok) {
    tabDataLoaded.userEventsTrend = true
  }
}

const fetchUserEventsTrend = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  if (!silent) {
    userEventsLoading.value = true
  }
  try {
    const data = await trackingEventApi.getMinuteLiveActivityTrend({
      startTime: 0,
      endTime: 0,
      limit: USER_EVENTS_TREND_MAX_POINTS,
      bypassDashboardCache,
    })
    userEventsTrend.points = data.points || []
    userEventsRangeTruncated.value = Boolean(data.rangeTruncated)
    userEventsTrendMaxPoints.value = data.trendMaxPoints || USER_EVENTS_TREND_MAX_POINTS
    await nextTick()
    setTimeout(() => {
      userEventsChartRef.value?.resize()
    }, 0)
    return true
  } catch (error) {
    console.error('fetch user events trend failed:', error)
    ElMessage.error(t('pages.trackingMinuteLiveActivity.fetchFailed'))
    return false
  } finally {
    if (!silent) {
      userEventsLoading.value = false
    }
  }
}

const fetchResourceMetricTrend = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  if (!silent) {
    onlineLoading.value = true
  }
  try {
    const data = await sysStatApi.getResourceMetricOnlineTrend(
        buildDashboardResourceMetricTrendQuery(DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS, bypassDashboardCache),
    )
    resourceMetricTrend.points = data.points || []
    onlineRangeTruncated.value = Boolean(data.rangeTruncated)
    onlineTrendMaxPoints.value = data.trendMaxPoints ?? 0
    await nextTick()
    setTimeout(() => {
      onlineChartRef.value?.resize()
    }, 0)
    return true
  } catch (error) {
    console.error('fetch online trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchOnlineTrendFailed'))
    return false
  } finally {
    if (!silent) {
      onlineLoading.value = false
    }
  }
}

const fetchGameStatTrend = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  if (!silent) {
    gameBetLoading.value = true
  }
  try {
    const data = await sysStatApi.getResourceMetricGameStatTrend(
        buildDashboardResourceMetricTrendQuery(DASHBOARD_RESOURCE_METRIC_RETENTION_DAYS, bypassDashboardCache),
    )
    gameStatTrend.points = (data.points || []).map((item) => ({
      time: item.time,
      count: Number(item.count || 0),
      onlineCount: Number(item.onlineCount || 0),
    }))
    gameStatRangeTruncated.value = Boolean(data.rangeTruncated)
    gameStatTrendMaxPoints.value = data.trendMaxPoints ?? 0
    await nextTick()
    setTimeout(() => {
      resizeGameBetCharts()
    }, 0)
    return true
  } catch (error) {
    console.error('fetch game stat trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchGameBetUserTrendFailed'))
    return false
  } finally {
    if (!silent) {
      gameBetLoading.value = false
    }
  }
}

const fetchUserStatTrend = async (silent = false, bypassDashboardCache = false): Promise<boolean> => {
  if (!silent) {
    trendLoading.value = true
  }
  try {
    const data = await sysStatApi.getUserStatTrend({ bypassDashboardCache })
    userStatTrend.daily = data.daily || []
    userStatTrend.weekly = data.weekly || []
    userStatTrend.monthly = data.monthly || []
    await nextTick()
    setTimeout(() => {
      resizeUserCharts()
      resizeRechargeCharts()
      resizeClientRechargeCharts()
      resizeClientLoginCharts()
    }, 0)
    return true
  } catch (error) {
    console.error('fetch user stat trend failed:', error)
    ElMessage.error(t('pages.dashboard.fetchUserTrendFailed'))
    return false
  } finally {
    if (!silent) {
      trendLoading.value = false
    }
  }
}

const resizeUserCharts = () => {
  if (activePeriod.value === 'daily') {
    dailyLineChartRef.value?.resize()
    return
  }
  if (activePeriod.value === 'weekly') {
    weeklyLineChartRef.value?.resize()
    return
  }
  monthlyLineChartRef.value?.resize()
}

const resizeRechargeCharts = () => {
  if (activeRechargePeriod.value === 'daily') {
    rechargeDailyBarChartRef.value?.resize()
    return
  }
  if (activeRechargePeriod.value === 'weekly') {
    rechargeWeeklyBarChartRef.value?.resize()
    return
  }
  rechargeMonthlyBarChartRef.value?.resize()
}

const resizeClientRechargeCharts = () => {
  if (activeClientRechargePeriod.value === 'daily') {
    clientRechargeDailyChartRef.value?.resize()
    return
  }
  if (activeClientRechargePeriod.value === 'weekly') {
    clientRechargeWeeklyChartRef.value?.resize()
    return
  }
  clientRechargeMonthlyChartRef.value?.resize()
}

const resizeClientLoginCharts = () => {
  if (activeClientLoginPeriod.value === 'daily') {
    clientLoginDailyChartRef.value?.resize()
    return
  }
  if (activeClientLoginPeriod.value === 'weekly') {
    clientLoginWeeklyChartRef.value?.resize()
    return
  }
  clientLoginMonthlyChartRef.value?.resize()
}

const resizeGameBetCharts = () => {
  gameBetDailyChartRef.value?.resize()
}

const handleMainTabChange = async () => {
  await refreshActiveMainTab(false)
  await nextTick()
  if (activeMainTab.value === 'user') {
    resizeUserCharts()
    return
  }
  if (activeMainTab.value === 'rechargeUser') {
    resizeRechargeCharts()
    return
  }
  if (activeMainTab.value === 'clientRecharge') {
    resizeClientRechargeCharts()
    return
  }
  if (activeMainTab.value === 'clientLogin') {
    resizeClientLoginCharts()
    return
  }
  if (activeMainTab.value === 'gameStats') {
    resizeGameBetCharts()
    return
  }
  if (activeMainTab.value === 'userEvents') {
    userEventsChartRef.value?.resize()
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

const handleRechargePeriodTabChange = async () => {
  await nextTick()
  resizeRechargeCharts()
}

const handleClientRechargePeriodTabChange = async () => {
  await nextTick()
  resizeClientRechargeCharts()
}

const handleClientLoginPeriodTabChange = async () => {
  await nextTick()
  resizeClientLoginCharts()
}

const handleUserEventMetricTabChange = async () => {
  await nextTick()
  userEventsChartRef.value?.resize()
}

const formatCount = (value: string | number | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString(locale.value)
}

onMounted(() => {
  dashboardRouteActive.value = true
  skipNextActivatedRefresh = true
  bindDocumentVisibilityListener()
  // 首屏仅拉「基础数据」；其余主 Tab 在 @tab-change 时懒加载
  void refreshActiveMainTab(false)
  restartAutoRefreshTimer()
})

onActivated(() => {
  dashboardRouteActive.value = true
  bindDocumentVisibilityListener()
  if (skipNextActivatedRefresh) {
    skipNextActivatedRefresh = false
    restartAutoRefreshTimer()
    return
  }
  void refreshActiveMainTab(true, true)
  restartAutoRefreshTimer()
})

onDeactivated(() => {
  dashboardRouteActive.value = false
  stopDashboardAutoRefresh()
})

onUnmounted(() => {
  dashboardRouteActive.value = false
  stopDashboardAutoRefresh()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.dashboard-tab-pane-body {
  min-height: 280px;
}

.page-toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  margin-bottom: 12px;
}

.page-toolbar-label {
  font-size: 13px;
  color: #64748b;
}

.page-toolbar-export-all {
  margin-left: 8px;
}

.auto-refresh-select {
  width: 140px;
}

.tab-toolbar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 16px;
}

.tab-toolbar-basic {
  justify-content: flex-start;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 10px;
}

.online-chart {
  margin-top: 0;
}

.basic-stat-board {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.basic-stat-board--compact {
  gap: 14px;
}

.basic-stat-board--compact .basic-stat-section-head {
  margin-bottom: 8px;
}

.basic-stat-board--compact .stat-card {
  min-height: 80px;
  padding: 12px 14px 12px 16px;
}

.basic-stat-board--compact .stat-label {
  margin-bottom: 6px;
  font-size: 12px;
}

.basic-stat-board--compact .stat-value {
  font-size: 22px;
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

.basic-stat-board--compact .basic-stat-grid {
  gap: 10px;
}

.basic-stat-grid-realtime {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.basic-stat-grid-today {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.basic-stat-grid-total {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.basic-stat-grid-anchor-payout {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.dashboard-realtime-row {
  margin-bottom: 18px;
}

.dashboard-realtime-row-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.dashboard-realtime-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  width: 100%;
}

.dashboard-realtime-grid--fit {
  grid-template-columns: repeat(2, minmax(160px, 240px));
  justify-content: start;
  margin-bottom: 16px;
}

.dashboard-realtime-card {
  min-height: 112px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.dashboard-realtime-card--compact {
  min-height: 88px;
  padding: 14px 16px 14px 18px;
}

.dashboard-realtime-card .stat-label {
  margin-bottom: 10px;
}

.dashboard-realtime-card--compact .stat-label {
  margin-bottom: 6px;
  font-size: 12px;
}

.dashboard-realtime-card .stat-value {
  font-size: 32px;
}

.dashboard-realtime-card--compact .stat-value {
  font-size: 28px;
  line-height: 1.15;
}

.game-stats-pane {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.game-stats-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.game-stats-toolbar-head {
  flex: 1;
  min-width: 0;
  margin-bottom: 0;
}

.game-stats-chart-panel {
  border: 1px solid #e8eef5;
  border-radius: 14px;
  background: #fff;
  padding: 14px 16px 14px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
}

.game-stats-chart-head {
  margin-bottom: 4px;
}

.game-stats-chart-title {
  font-size: 14px;
  font-weight: 600;
  color: #334155;
  margin-bottom: 4px;
}

.game-stats-chart-hint {
  margin: 0;
  font-size: 12px;
  line-height: 1.45;
  color: #94a3b8;
}

.game-stats-chart {
  margin-top: 4px;
}

.user-events-metric-tabs :deep(.el-tabs__header) {
  margin-bottom: 8px;
}

.user-events-metric-tabs :deep(.el-tabs__content) {
  padding-top: 4px;
}

.game-bet-summary-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin-bottom: 16px;
}

.game-bet-event-hint {
  margin: 0 0 14px;
  font-size: 13px;
  line-height: 1.5;
  color: #64748b;
}

.tab-toolbar-trend-range {
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
}

.game-bet-range-truncated-hint {
  margin: 0 0 12px;
  font-size: 13px;
  color: #b45309;
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

.stat-value--red {
  color: #dc2626;
}

.stat-value--blue {
  color: #2563eb;
}

.stat-value--violet {
  color: #7c3aed;
}

.stat-value--orange {
  color: #ea580c;
}

.stat-unit {
  margin-left: 6px;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: #64748b;
}

.tone-neutral {
  --stat-accent: #cbd5e1;
  background: #f8fafc;
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
  .basic-stat-grid-total,
  .basic-stat-grid-anchor-payout,
  .basic-stat-grid-realtime {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .basic-stat-grid-today,
  .basic-stat-grid-total,
  .basic-stat-grid-anchor-payout,
  .basic-stat-grid-realtime {
    grid-template-columns: 1fr;
  }

  .dashboard-realtime-grid {
    grid-template-columns: 1fr;
  }

  .stat-value {
    font-size: 22px;
  }
}
</style>

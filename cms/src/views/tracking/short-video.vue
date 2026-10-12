<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingShortVideoManagement') }}</span>
          <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingShortVideo.tabPlayStart')" name="playStart">
          <p class="page-desc">{{ t('pages.trackingShortVideo.playStartDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingShortVideo.tabPaywallClick')" name="paywallClick">
          <p class="page-desc">{{ t('pages.trackingShortVideo.paywallClickDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingShortVideo.tabPaySuccess')" name="paySuccess">
          <p class="page-desc">{{ t('pages.trackingShortVideo.paySuccessDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingShortVideo.tabWatchEnd')" name="watchEnd">
          <p class="page-desc">{{ t('pages.trackingShortVideo.watchEndDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingShortVideo.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingShortVideo.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingShortVideo.monthTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.monthCount) }}</div>
        </div>
      </div>

      <el-tabs v-model="activePeriod" lazy @tab-change="handlePeriodChange">
        <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
          <CountTrendChart
              ref="dailyChartRef"
              :data="activeTrend.daily"
              :series-name="activeSeriesName"
              :title="activeChartDaily"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
          <CountTrendChart
              ref="weeklyChartRef"
              :data="activeTrend.weekly"
              :series-name="activeSeriesName"
              :title="activeChartWeekly"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
          <CountTrendChart
              ref="monthlyChartRef"
              :data="activeTrend.monthly"
              :series-name="activeSeriesName"
              :title="activeChartMonthly"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {trackingEventApi} from '@/api/modules/tracking-event'
import type {HotLiveRoomJoinTrendRes} from '@/types/api'
import CountTrendChart from './components/count-trend-chart.vue'

type TrackingEventKey = 'playStart' | 'paySuccess' | 'paywallClick' | 'watchEnd'

const {t, locale} = useI18n()
const loading = ref(false)
const activeEvent = ref<TrackingEventKey>('playStart')
const activePeriod = ref<'daily' | 'weekly' | 'monthly'>('daily')
const dailyChartRef = ref<InstanceType<typeof CountTrendChart>>()
const weeklyChartRef = ref<InstanceType<typeof CountTrendChart>>()
const monthlyChartRef = ref<InstanceType<typeof CountTrendChart>>()

const emptyTrend = (): HotLiveRoomJoinTrendRes => ({
  todayCount: 0,
  weekCount: 0,
  monthCount: 0,
  daily: [],
  weekly: [],
  monthly: [],
})

const playStartTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const paySuccessTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const paywallClickTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const watchEndTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const trendByEvent: Record<TrackingEventKey, HotLiveRoomJoinTrendRes> = {
  playStart: playStartTrend,
  paySuccess: paySuccessTrend,
  paywallClick: paywallClickTrend,
  watchEnd: watchEndTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const activeSeriesName = computed(() => {
  const key = activeEvent.value
  if (key === 'paySuccess') {
    return t('pages.trackingShortVideo.paySuccessSeriesName')
  }
  if (key === 'paywallClick') {
    return t('pages.trackingShortVideo.paywallClickSeriesName')
  }
  if (key === 'watchEnd') {
    return t('pages.trackingShortVideo.watchEndSeriesName')
  }
  return t('pages.trackingShortVideo.playStartSeriesName')
})

const activeChartDaily = computed(() => {
  const key = activeEvent.value
  if (key === 'paySuccess') {
    return t('pages.trackingShortVideo.paySuccessChartDaily')
  }
  if (key === 'paywallClick') {
    return t('pages.trackingShortVideo.paywallClickChartDaily')
  }
  if (key === 'watchEnd') {
    return t('pages.trackingShortVideo.watchEndChartDaily')
  }
  return t('pages.trackingShortVideo.playStartChartDaily')
})

const activeChartWeekly = computed(() => {
  const key = activeEvent.value
  if (key === 'paySuccess') {
    return t('pages.trackingShortVideo.paySuccessChartWeekly')
  }
  if (key === 'paywallClick') {
    return t('pages.trackingShortVideo.paywallClickChartWeekly')
  }
  if (key === 'watchEnd') {
    return t('pages.trackingShortVideo.watchEndChartWeekly')
  }
  return t('pages.trackingShortVideo.playStartChartWeekly')
})

const activeChartMonthly = computed(() => {
  const key = activeEvent.value
  if (key === 'paySuccess') {
    return t('pages.trackingShortVideo.paySuccessChartMonthly')
  }
  if (key === 'paywallClick') {
    return t('pages.trackingShortVideo.paywallClickChartMonthly')
  }
  if (key === 'watchEnd') {
    return t('pages.trackingShortVideo.watchEndChartMonthly')
  }
  return t('pages.trackingShortVideo.playStartChartMonthly')
})

const formatCount = (value: number | string | null | undefined) => {
  if (value == null || value === '') {
    return '0'
  }
  return Number(value).toLocaleString(locale.value)
}

const assignTrend = (target: HotLiveRoomJoinTrendRes, data: HotLiveRoomJoinTrendRes) => {
  Object.assign(target, {
    todayCount: data.todayCount ?? 0,
    weekCount: data.weekCount ?? 0,
    monthCount: data.monthCount ?? 0,
    daily: data.daily || [],
    weekly: data.weekly || [],
    monthly: data.monthly || [],
  })
}

const resizeActiveChart = () => {
  void nextTick(() => {
    if (activePeriod.value === 'daily') {
      dailyChartRef.value?.resize()
    } else if (activePeriod.value === 'weekly') {
      weeklyChartRef.value?.resize()
    } else {
      monthlyChartRef.value?.resize()
    }
  })
}

const handlePeriodChange = () => {
  resizeActiveChart()
}

const handleEventChange = () => {
  void fetchActiveTrend()
}

const fetchPlayStartTrend = async () => {
  const data = await trackingEventApi.getShortVideoPlayStartTrend()
  assignTrend(playStartTrend, data)
}

const fetchPaywallClickTrend = async () => {
  const data = await trackingEventApi.getShortVideoPaywallClickTrend()
  assignTrend(paywallClickTrend, data)
}

const fetchPaySuccessTrend = async () => {
  const data = await trackingEventApi.getShortVideoPaySuccessTrend()
  assignTrend(paySuccessTrend, data)
}

const fetchWatchEndTrend = async () => {
  const data = await trackingEventApi.getShortVideoWatchEndTrend()
  assignTrend(watchEndTrend, data)
}

const fetchActiveTrend = async () => {
  loading.value = true
  try {
    if (activeEvent.value === 'paySuccess') {
      await fetchPaySuccessTrend()
    } else if (activeEvent.value === 'paywallClick') {
      await fetchPaywallClickTrend()
    } else if (activeEvent.value === 'watchEnd') {
      await fetchWatchEndTrend()
    } else {
      await fetchPlayStartTrend()
    }
    resizeActiveChart()
  } catch (error) {
    console.error('fetch short video tracking trend failed:', error)
    ElMessage.error(t('pages.trackingShortVideo.fetchFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchActiveTrend()
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

.event-tabs {
  margin-bottom: 8px;
}

.page-desc {
  margin: 0 0 16px;
  font-size: 13px;
  color: #64748b;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 10px;
  padding: 16px 18px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.stat-label {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #0f172a;
}

.tone-blue {
  border-color: #bfdbfe;
  background: #eff6ff;
}

.tone-teal {
  border-color: #99f6e4;
  background: #f0fdfa;
}

.tone-purple {
  border-color: #ddd6fe;
  background: #f5f3ff;
}
</style>

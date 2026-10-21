<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingMiniGameRoundStartManagement') }}</span>
          <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingMiniGameRoundStart.tabSpinClick')" name="spinClick">
          <p class="page-desc">{{ t('pages.trackingMiniGameRoundStart.spinClickDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingMiniGameRoundStart.tabPlusClick')" name="plusClick">
          <p class="page-desc">{{ t('pages.trackingMiniGameRoundStart.plusClickDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingMiniGameRoundStart.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingMiniGameRoundStart.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingMiniGameRoundStart.monthTotal') }}</div>
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

type TrackingEventKey = 'spinClick' | 'plusClick'

const {t, locale} = useI18n()
const loading = ref(false)
const activeEvent = ref<TrackingEventKey>('spinClick')
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

const spinClickTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const plusClickTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const trendByEvent: Record<TrackingEventKey, HotLiveRoomJoinTrendRes> = {
  spinClick: spinClickTrend,
  plusClick: plusClickTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const activeSeriesName = computed(() => {
  if (activeEvent.value === 'plusClick') {
    return t('pages.trackingMiniGameRoundStart.plusClickSeriesName')
  }
  return t('pages.trackingMiniGameRoundStart.spinClickSeriesName')
})

const activeChartDaily = computed(() => {
  if (activeEvent.value === 'plusClick') {
    return t('pages.trackingMiniGameRoundStart.plusClickChartDaily')
  }
  return t('pages.trackingMiniGameRoundStart.spinClickChartDaily')
})

const activeChartWeekly = computed(() => {
  if (activeEvent.value === 'plusClick') {
    return t('pages.trackingMiniGameRoundStart.plusClickChartWeekly')
  }
  return t('pages.trackingMiniGameRoundStart.spinClickChartWeekly')
})

const activeChartMonthly = computed(() => {
  if (activeEvent.value === 'plusClick') {
    return t('pages.trackingMiniGameRoundStart.plusClickChartMonthly')
  }
  return t('pages.trackingMiniGameRoundStart.spinClickChartMonthly')
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

const fetchActiveTrend = async () => {
  loading.value = true
  try {
    if (activeEvent.value === 'plusClick') {
      assignTrend(plusClickTrend, await trackingEventApi.getMiniGamePlusClickTrend())
    } else {
      assignTrend(spinClickTrend, await trackingEventApi.getMiniGameSpinClickTrend())
    }
    resizeActiveChart()
  } catch (error) {
    console.error('fetch mini game tracking trend failed:', error)
    ElMessage.error(t('pages.trackingMiniGameRoundStart.fetchFailed'))
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

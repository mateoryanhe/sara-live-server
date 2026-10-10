<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingCall1v1InitiateManagement') }}</span>
          <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingCall1v1Initiate.tabLiveRoom')" name="liveRoom">
          <p class="page-desc">{{ t('pages.trackingCall1v1Initiate.desc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingCall1v1Initiate.tabOneToOneRoom')" name="oneToOneRoom">
          <p class="page-desc">{{ t('pages.trackingCall1v1Initiate.oneToOneRoomDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingCall1v1Initiate.tabConnectSuccess')" name="connectSuccess">
          <p class="page-desc">{{ t('pages.trackingCall1v1Initiate.connectSuccessDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.monthTotal') }}</div>
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

type TrackingEventKey = 'liveRoom' | 'oneToOneRoom' | 'connectSuccess'

const {t, locale} = useI18n()
const loading = ref(false)
const activeEvent = ref<TrackingEventKey>('liveRoom')
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

const liveRoomTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const oneToOneRoomTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const connectSuccessTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const trendByEvent: Record<TrackingEventKey, HotLiveRoomJoinTrendRes> = {
  liveRoom: liveRoomTrend,
  oneToOneRoom: oneToOneRoomTrend,
  connectSuccess: connectSuccessTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const activeSeriesName = computed(() => {
  const key = activeEvent.value
  if (key === 'oneToOneRoom') {
    return t('pages.trackingCall1v1Initiate.oneToOneRoomSeriesName')
  }
  if (key === 'connectSuccess') {
    return t('pages.trackingCall1v1Initiate.connectSuccessSeriesName')
  }
  return t('pages.trackingCall1v1Initiate.seriesName')
})

const activeChartDaily = computed(() => {
  const key = activeEvent.value
  if (key === 'oneToOneRoom') {
    return t('pages.trackingCall1v1Initiate.oneToOneRoomChartDaily')
  }
  if (key === 'connectSuccess') {
    return t('pages.trackingCall1v1Initiate.connectSuccessChartDaily')
  }
  return t('pages.trackingCall1v1Initiate.chartDaily')
})

const activeChartWeekly = computed(() => {
  const key = activeEvent.value
  if (key === 'oneToOneRoom') {
    return t('pages.trackingCall1v1Initiate.oneToOneRoomChartWeekly')
  }
  if (key === 'connectSuccess') {
    return t('pages.trackingCall1v1Initiate.connectSuccessChartWeekly')
  }
  return t('pages.trackingCall1v1Initiate.chartWeekly')
})

const activeChartMonthly = computed(() => {
  const key = activeEvent.value
  if (key === 'oneToOneRoom') {
    return t('pages.trackingCall1v1Initiate.oneToOneRoomChartMonthly')
  }
  if (key === 'connectSuccess') {
    return t('pages.trackingCall1v1Initiate.connectSuccessChartMonthly')
  }
  return t('pages.trackingCall1v1Initiate.chartMonthly')
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

const fetchLiveRoomTrend = async () => {
  const data = await trackingEventApi.getCall1v1InitiateTrend()
  assignTrend(liveRoomTrend, data)
}

const fetchOneToOneRoomTrend = async () => {
  const data = await trackingEventApi.getCall1v1RoomCallTrend()
  assignTrend(oneToOneRoomTrend, data)
}

const fetchConnectSuccessTrend = async () => {
  const data = await trackingEventApi.getCall1v1ConnectSuccessTrend()
  assignTrend(connectSuccessTrend, data)
}

const fetchActiveTrend = async () => {
  loading.value = true
  try {
    if (activeEvent.value === 'oneToOneRoom') {
      await fetchOneToOneRoomTrend()
    } else if (activeEvent.value === 'connectSuccess') {
      await fetchConnectSuccessTrend()
    } else {
      await fetchLiveRoomTrend()
    }
    resizeActiveChart()
  } catch (error) {
    console.error('fetch call tracking trend failed:', error)
    ElMessage.error(t('pages.trackingCall1v1Initiate.fetchFailed'))
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

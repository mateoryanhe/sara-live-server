<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingVoiceChatRoomManagement') }}</span>
          <div class="card-header-actions">
            <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
            <el-button v-if="canExport" @click="handleExport">{{ t('common.export') }}</el-button>
          </div>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingVoiceChatRoom.tabEnter')" name="enter">
          <p class="page-desc">{{ t('pages.trackingVoiceChatRoom.enterDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingVoiceChatRoom.tabMicRequest')" name="micRequest">
          <p class="page-desc">{{ t('pages.trackingVoiceChatRoom.micRequestDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingVoiceChatRoom.tabGiftClick')" name="giftClick">
          <p class="page-desc">{{ t('pages.trackingVoiceChatRoom.giftClickDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingVoiceChatRoom.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingVoiceChatRoom.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingVoiceChatRoom.monthTotal') }}</div>
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
import {useTrackingTrendExport} from '@/composables/useTrackingTrendExport'
import type {HotLiveRoomJoinTrendRes} from '@/types/api'
import {pickTrackingTrendPeriodPoints} from '@/utils/tracking-trend-export'
import CountTrendChart from './components/count-trend-chart.vue'

type TrackingEventKey = 'enter' | 'micRequest' | 'giftClick'

const {t, locale} = useI18n()
const loading = ref(false)
const activeEvent = ref<TrackingEventKey>('enter')
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

const enterTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const micRequestTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const giftClickTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const trendByEvent: Record<TrackingEventKey, HotLiveRoomJoinTrendRes> = {
  enter: enterTrend,
  micRequest: micRequestTrend,
  giftClick: giftClickTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const labelKey = (base: string) => {
  const key = activeEvent.value
  if (key === 'micRequest') {
    return `pages.trackingVoiceChatRoom.micRequest${base}`
  }
  if (key === 'giftClick') {
    return `pages.trackingVoiceChatRoom.giftClick${base}`
  }
  return `pages.trackingVoiceChatRoom.enter${base}`
}

const activeSeriesName = computed(() => t(labelKey('SeriesName')))
const activeChartDaily = computed(() => t(labelKey('ChartDaily')))
const activeChartWeekly = computed(() => t(labelKey('ChartWeekly')))
const activeChartMonthly = computed(() => t(labelKey('ChartMonthly')))

const {canExport, handleExport} = useTrackingTrendExport('TrackingVoiceChatRoomManagement', () => {
  const trend = activeTrend.value
  const period = activePeriod.value
  return {
    fileSlug: 'voice-chat-room',
    eventSlug: activeEvent.value,
    period,
    todayCount: Number(trend.todayCount) || 0,
    weekCount: Number(trend.weekCount) || 0,
    monthCount: Number(trend.monthCount) || 0,
    points: pickTrackingTrendPeriodPoints(trend, period),
    countColumnHeader: activeSeriesName.value,
    summaryLabels: {
      today: t('pages.trackingVoiceChatRoom.todayTotal'),
      week: t('pages.trackingVoiceChatRoom.weekTotal'),
      month: t('pages.trackingVoiceChatRoom.monthTotal'),
    },
  }
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

const fetchEnterTrend = async () => {
  assignTrend(enterTrend, await trackingEventApi.getVoiceRoomEnterTrend())
}

const fetchMicRequestTrend = async () => {
  assignTrend(micRequestTrend, await trackingEventApi.getVoiceRoomMicRequestTrend())
}

const fetchGiftClickTrend = async () => {
  assignTrend(giftClickTrend, await trackingEventApi.getVoiceRoomGiftClickTrend())
}

const fetchActiveTrend = async () => {
  loading.value = true
  try {
    if (activeEvent.value === 'micRequest') {
      await fetchMicRequestTrend()
    } else if (activeEvent.value === 'giftClick') {
      await fetchGiftClickTrend()
    } else {
      await fetchEnterTrend()
    }
    resizeActiveChart()
  } catch (error) {
    console.error('fetch voice chat room tracking trend failed:', error)
    ElMessage.error(t('pages.trackingVoiceChatRoom.fetchFailed'))
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

.card-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
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

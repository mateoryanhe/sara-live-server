<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingDownloadManagement') }}</span>
          <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingDownload.tabH5')" name="h5">
          <p class="page-desc">{{ t('pages.trackingDownload.h5Desc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingDownload.tabOfficial')" name="official">
          <p class="page-desc">{{ t('pages.trackingDownload.officialDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingDownload.tabShareShort')" name="shareShort">
          <p class="page-desc">{{ t('pages.trackingDownload.shareShortDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingDownload.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingDownload.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingDownload.monthTotal') }}</div>
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

type TrackingEventKey = 'h5' | 'official' | 'shareShort'

const {t, locale} = useI18n()
const loading = ref(false)
const activeEvent = ref<TrackingEventKey>('h5')
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

const h5Trend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const officialTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())
const shareShortTrend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const trendByEvent: Record<TrackingEventKey, HotLiveRoomJoinTrendRes> = {
  h5: h5Trend,
  official: officialTrend,
  shareShort: shareShortTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const labelKey = (suffix: string) => {
  if (activeEvent.value === 'official') {
    return `pages.trackingDownload.official${suffix}`
  }
  if (activeEvent.value === 'shareShort') {
    return `pages.trackingDownload.shareShort${suffix}`
  }
  return `pages.trackingDownload.h5${suffix}`
}

const activeSeriesName = computed(() => t(labelKey('SeriesName')))
const activeChartDaily = computed(() => t(labelKey('ChartDaily')))
const activeChartWeekly = computed(() => t(labelKey('ChartWeekly')))
const activeChartMonthly = computed(() => t(labelKey('ChartMonthly')))

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

const fetchH5Trend = async () => {
  assignTrend(h5Trend, await trackingEventApi.getH5DownloadTrend())
}

const fetchOfficialTrend = async () => {
  assignTrend(officialTrend, await trackingEventApi.getOfficialSiteDownloadTrend())
}

const fetchShareShortTrend = async () => {
  assignTrend(shareShortTrend, await trackingEventApi.getShareShortLinkClickTrend())
}

const fetchActiveTrend = async () => {
  loading.value = true
  try {
    if (activeEvent.value === 'official') {
      await fetchOfficialTrend()
    } else if (activeEvent.value === 'shareShort') {
      await fetchShareShortTrend()
    } else {
      await fetchH5Trend()
    }
    resizeActiveChart()
  } catch (error) {
    console.error('fetch download tracking trend failed:', error)
    ElMessage.error(t('pages.trackingDownload.fetchFailed'))
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

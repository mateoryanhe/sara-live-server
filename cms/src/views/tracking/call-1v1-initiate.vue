<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingCall1v1InitiateManagement') }}</span>
          <el-button :loading="loading" @click="fetchTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.trackingCall1v1Initiate.desc') }}</p>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(trend.todayCount) }}</div>
        </div>
        <div class="stat-card tone-teal">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.weekTotal') }}</div>
          <div class="stat-value">{{ formatCount(trend.weekCount) }}</div>
        </div>
        <div class="stat-card tone-purple">
          <div class="stat-label">{{ t('pages.trackingCall1v1Initiate.monthTotal') }}</div>
          <div class="stat-value">{{ formatCount(trend.monthCount) }}</div>
        </div>
      </div>

      <el-tabs v-model="activePeriod" lazy @tab-change="handlePeriodChange">
        <el-tab-pane :label="t('pages.dashboard.periodDaily')" name="daily">
          <CountTrendChart
              ref="dailyChartRef"
              :data="trend.daily"
              :series-name="t('pages.trackingCall1v1Initiate.seriesName')"
              :title="t('pages.trackingCall1v1Initiate.chartDaily')"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodWeekly')" name="weekly">
          <CountTrendChart
              ref="weeklyChartRef"
              :data="trend.weekly"
              :series-name="t('pages.trackingCall1v1Initiate.seriesName')"
              :title="t('pages.trackingCall1v1Initiate.chartWeekly')"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.dashboard.periodMonthly')" name="monthly">
          <CountTrendChart
              ref="monthlyChartRef"
              :data="trend.monthly"
              :series-name="t('pages.trackingCall1v1Initiate.seriesName')"
              :title="t('pages.trackingCall1v1Initiate.chartMonthly')"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {trackingEventApi} from '@/api/modules/tracking-event'
import type {HotLiveRoomJoinTrendRes} from '@/types/api'
import CountTrendChart from './components/count-trend-chart.vue'

const {t, locale} = useI18n()
const loading = ref(false)
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

const trend = reactive<HotLiveRoomJoinTrendRes>(emptyTrend())

const formatCount = (value: number | string | null | undefined) => {
  if (value == null || value === '') {
    return '0'
  }
  return Number(value).toLocaleString(locale.value)
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

const fetchTrend = async () => {
  loading.value = true
  try {
    const data = await trackingEventApi.getCall1v1InitiateTrend()
    Object.assign(trend, {
      todayCount: data.todayCount ?? 0,
      weekCount: data.weekCount ?? 0,
      monthCount: data.monthCount ?? 0,
      daily: data.daily || [],
      weekly: data.weekly || [],
      monthly: data.monthly || [],
    })
    resizeActiveChart()
  } catch (error) {
    console.error('fetch call 1v1 initiate trend failed:', error)
    ElMessage.error(t('pages.trackingCall1v1Initiate.fetchFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchTrend()
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

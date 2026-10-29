<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingMinuteLiveActivityManagement') }}</span>
          <el-button :loading="loading" @click="fetchTrend">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.trackingMinuteLiveActivity.pageDesc') }}</p>

      <div class="chart-panel">
        <div class="chart-head">
          <div class="chart-title">{{ t('pages.trackingMinuteLiveActivity.chartTitleLast3Days') }}</div>
          <p class="chart-hint">{{ t('pages.trackingMinuteLiveActivity.chartHint') }}</p>
          <p v-if="rangeTruncated" class="range-truncated-hint">
            {{ t('pages.trackingMinuteLiveActivity.rangeTruncated', {max: trendMaxPoints}) }}
          </p>
        </div>
        <MinuteLiveActivityTrendChart
            ref="chartRef"
            :data="trend.points"
            :title="t('pages.trackingMinuteLiveActivity.chartTitleLast3Days')"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {trackingEventApi} from '@/api/modules/tracking-event'
import type {MinuteLiveActivityTrendRes} from '@/types/api'
import MinuteLiveActivityTrendChart from './components/minute-live-activity-trend-chart.vue'

/** 3 天 × 1440 分钟，与后端 retention / maxTrendPoints 一致 */
const MINUTE_TREND_MAX_POINTS = 4320

const {t} = useI18n()
const loading = ref(false)
const chartRef = ref<InstanceType<typeof MinuteLiveActivityTrendChart>>()
const rangeTruncated = ref(false)
const trendMaxPoints = ref(MINUTE_TREND_MAX_POINTS)

const trend = reactive<MinuteLiveActivityTrendRes>({points: []})

const fetchTrend = async () => {
  loading.value = true
  try {
    const res = await trackingEventApi.getMinuteLiveActivityTrend({
      startTime: 0,
      endTime: 0,
      limit: MINUTE_TREND_MAX_POINTS,
    })
    trend.points = res.points || []
    rangeTruncated.value = Boolean(res.rangeTruncated)
    trendMaxPoints.value = res.trendMaxPoints || MINUTE_TREND_MAX_POINTS
    await nextTick()
    setTimeout(() => chartRef.value?.resize(), 0)
  } catch (error) {
    console.error('fetch minute live activity trend failed:', error)
    ElMessage.error(t('pages.trackingMinuteLiveActivity.fetchFailed'))
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
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.chart-panel {
  margin-top: 4px;
}

.chart-head {
  margin-bottom: 8px;
}

.chart-title {
  font-size: 15px;
  font-weight: 600;
  color: #334155;
}

.chart-hint {
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.range-truncated-hint {
  margin: 6px 0 0;
  color: var(--el-color-warning);
  font-size: 13px;
}
</style>

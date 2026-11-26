<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingGoogleRechargeManagement') }}</span>
          <div class="card-header-actions">
            <el-button :loading="loading" @click="fetchActiveTrend">{{ t('common.refresh') }}</el-button>
          </div>
        </div>
      </template>

      <el-tabs v-model="activeEvent" class="event-tabs" @tab-change="handleEventChange">
        <el-tab-pane :label="t('pages.trackingGoogleRecharge.tabList')" name="list">
          <p class="page-desc">{{ t('pages.trackingGoogleRecharge.listDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingGoogleRecharge.tabOrder')" name="order">
          <p class="page-desc">{{ t('pages.trackingGoogleRecharge.orderDesc') }}</p>
        </el-tab-pane>
        <el-tab-pane :label="t('pages.trackingGoogleRecharge.tabResult')" name="result">
          <p class="page-desc">{{ t('pages.trackingGoogleRecharge.resultDesc') }}</p>
        </el-tab-pane>
      </el-tabs>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingGoogleRecharge.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(activeTrend.todayCount) }}</div>
        </div>
      </div>

      <CountTrendChart
          ref="dailyChartRef"
          :data="activeTrend.daily"
          :series-name="activeSeriesName"
          :title="t('pages.trackingGoogleRecharge.chartDaily')"
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {trackingEventApi} from '@/api/modules/tracking-event'
import type {GoogleRechargeDailyTrendRes} from '@/types/api'
import CountTrendChart from './components/count-trend-chart.vue'

type EventKey = 'list' | 'order' | 'result'

const {t} = useI18n()
const loading = ref(false)
const activeEvent = ref<EventKey>('list')
const dailyChartRef = ref<InstanceType<typeof CountTrendChart>>()
const loaded = reactive<Record<EventKey, boolean>>({list: false, order: false, result: false})

const emptyTrend = (): GoogleRechargeDailyTrendRes => ({todayCount: 0, daily: []})

const listTrend = reactive<GoogleRechargeDailyTrendRes>(emptyTrend())
const orderTrend = reactive<GoogleRechargeDailyTrendRes>(emptyTrend())
const resultTrend = reactive<GoogleRechargeDailyTrendRes>(emptyTrend())

const trendByEvent: Record<EventKey, GoogleRechargeDailyTrendRes> = {
  list: listTrend,
  order: orderTrend,
  result: resultTrend,
}

const activeTrend = computed(() => trendByEvent[activeEvent.value])

const activeSeriesName = computed(() => {
  const key = activeEvent.value
  if (key === 'order') return t('pages.trackingGoogleRecharge.orderSeriesName')
  if (key === 'result') return t('pages.trackingGoogleRecharge.resultSeriesName')
  return t('pages.trackingGoogleRecharge.listSeriesName')
})

function formatCount(n: number | string | undefined): string {
  const v = Number(n ?? 0)
  if (!Number.isFinite(v)) return '0'
  return v.toLocaleString()
}

function applyTrend(target: GoogleRechargeDailyTrendRes, res: GoogleRechargeDailyTrendRes) {
  target.todayCount = res.todayCount ?? 0
  target.daily = res.daily ?? []
}

async function fetchEventTrend(key: EventKey) {
  loading.value = true
  try {
    let res: GoogleRechargeDailyTrendRes
    if (key === 'order') {
      res = await trackingEventApi.getGoogleRechargeOrderTrend()
    } else if (key === 'result') {
      res = await trackingEventApi.getGoogleRechargeResultTrend()
    } else {
      res = await trackingEventApi.getGoogleRechargeListTrend()
    }
    applyTrend(trendByEvent[key], res)
    loaded[key] = true
    await nextTick()
    dailyChartRef.value?.resize?.()
  } catch {
    ElMessage.error(t('pages.trackingGoogleRecharge.fetchFailed'))
  } finally {
    loading.value = false
  }
}

async function fetchActiveTrend() {
  await fetchEventTrend(activeEvent.value)
}

function handleEventChange(name: string | number) {
  const key = String(name) as EventKey
  if (!loaded[key]) {
    void fetchEventTrend(key)
  } else {
    void nextTick(() => dailyChartRef.value?.resize?.())
  }
}

onMounted(() => {
  void fetchEventTrend('list')
})

watch(activeEvent, () => {
  void nextTick(() => dailyChartRef.value?.resize?.())
})
</script>

<style scoped>
.page-desc {
  margin: 0 0 16px;
  color: #64748b;
  font-size: 14px;
  line-height: 1.6;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
  padding: 16px;
  background: #f8fafc;
}

.stat-label {
  color: #64748b;
  font-size: 13px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  margin-top: 8px;
}

.tone-blue {
  border-left: 4px solid #3b82f6;
}
</style>
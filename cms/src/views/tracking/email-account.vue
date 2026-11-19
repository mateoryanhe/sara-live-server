<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TrackingEmailAccountManagement') }}</span>
          <div class="card-header-actions">
            <el-button :loading="loading" @click="fetchTrend">{{ t('common.refresh') }}</el-button>
          </div>
        </div>
      </template>

      <p class="page-desc">{{ t('pages.trackingEmailAccount.pageDesc') }}</p>

      <div class="summary-grid">
        <div class="stat-card tone-blue">
          <div class="stat-label">{{ t('pages.trackingEmailAccount.todayTotal') }}</div>
          <div class="stat-value">{{ formatCount(trend.todayVal) }}</div>
        </div>
      </div>

      <CountTrendChart
          ref="dailyChartRef"
          :data="trend.daily"
          :series-name="t('pages.trackingEmailAccount.seriesName')"
          :title="t('pages.trackingEmailAccount.chartDaily')"
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {trackingEventApi} from '@/api/modules/tracking-event'
import type {EmailAccountTrendRes} from '@/types/api'
import CountTrendChart from './components/count-trend-chart.vue'

const {t} = useI18n()
const loading = ref(false)
const dailyChartRef = ref<InstanceType<typeof CountTrendChart>>()

const trend = reactive<EmailAccountTrendRes>({
  todayVal: 0,
  daily: [],
})

function formatCount(n: number | string | undefined): string {
  const v = Number(n ?? 0)
  if (!Number.isFinite(v)) return '0'
  return v.toLocaleString()
}

async function fetchTrend() {
  loading.value = true
  try {
    const res = await trackingEventApi.getEmailAccountTrend()
    trend.todayVal = res.todayVal ?? 0
    trend.daily = res.daily ?? []
    await nextTick()
    dailyChartRef.value?.resize?.()
  } catch {
    ElMessage.error(t('pages.trackingEmailAccount.fetchFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchTrend()
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
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px 20px;
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
</style>
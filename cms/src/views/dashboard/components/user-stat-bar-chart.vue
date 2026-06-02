<template>
  <div ref="chartRef" class="user-stat-bar-chart"></div>
</template>

<script lang="ts" setup>
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import * as echarts from 'echarts'
import type {UserStatTrendPoint} from '@/types/api'
import {USER_STAT_BAR_SERIES, getEnabledUserStatBarSeries} from '../user-stat-bar-series'

const props = defineProps<{
  data: UserStatTrendPoint[]
  title?: string
  /** 当前展示的指标 key(与 barMetrics 字段一致) */
  metricKey: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const activeMetric = computed(() => {
  const found = USER_STAT_BAR_SERIES.find((s) => s.key === props.metricKey && s.enabled)
  if (found) {
    return found
  }
  return getEnabledUserStatBarSeries()[0]
})

const buildOption = (points: UserStatTrendPoint[]): echarts.EChartsOption => {
  const cfg = activeMetric.value
  if (!cfg) {
    return {}
  }

  const times = points.map((item) => item.time)
  const hasTitle = Boolean(props.title)

  return {
    title: hasTitle
        ? {
          text: props.title,
          left: 'center',
          textStyle: {fontSize: 14, fontWeight: 500},
        }
        : undefined,
    tooltip: {
      trigger: 'axis',
      axisPointer: {type: 'shadow'},
    },
    grid: {
      left: 48,
      right: 24,
      top: hasTitle ? 48 : 24,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      data: times,
      axisLabel: {
        rotate: times.length > 10 ? 35 : 0,
      },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
    },
    series: [
      {
        name: cfg.label,
        type: 'bar',
        barMaxWidth: 40,
        data: points.map((item) => Number(item.barMetrics?.[cfg.key] ?? 0)),
        itemStyle: {color: cfg.color},
      },
    ],
  }
}

const renderChart = async () => {
  await nextTick()
  if (!chartRef.value) {
    return
  }
  if (!chartInstance) {
    chartInstance = echarts.init(chartRef.value)
  }
  chartInstance.setOption(buildOption(props.data || []), true)
  chartInstance.resize()
}

const handleResize = () => {
  chartInstance?.resize()
}

watch(
    () => [props.data, props.metricKey, activeMetric.value],
    () => {
      renderChart()
    },
    {deep: true},
)

onMounted(() => {
  renderChart()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({
  resize: () => chartInstance?.resize(),
})
</script>

<style scoped>
.user-stat-bar-chart {
  width: 100%;
  height: 360px;
}
</style>

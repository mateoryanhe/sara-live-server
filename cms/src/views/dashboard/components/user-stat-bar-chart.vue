<template>
  <div ref="chartRef" class="user-stat-bar-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {UserStatTrendPoint} from '@/types/api'
import {USER_STAT_BAR_SERIES, getEnabledUserStatBarSeries} from '../user-stat-bar-series'

const props = defineProps<{
  data: UserStatTrendPoint[]
  title?: string
  metricKey: string
}>()

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const activeMetric = computed(() => {
  const found = USER_STAT_BAR_SERIES.find((s) => s.key === props.metricKey && s.enabled)
  if (found) {
    return found
  }
  return getEnabledUserStatBarSeries()[0]
})

const activeMetricLabel = computed(() => {
  const cfg = activeMetric.value
  return cfg ? t(`pages.dashboard.${cfg.labelKey}`) : ''
})

const buildOption = (points: UserStatTrendPoint[]): EChartsOption => {
  const cfg = activeMetric.value
  if (!cfg) {
    return {}
  }

  const times = points.map((item) => item.time)
  const hasTitle = Boolean(props.title)
  const seriesName = activeMetricLabel.value

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
    },
    grid: {
      left: 48,
      right: 24,
      top: hasTitle ? 48 : 24,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
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
        name: seriesName,
        type: 'line',
        smooth: true,
        showSymbol: times.length <= 30,
        data: points.map((item) => Number(item.barMetrics?.[cfg.key] ?? 0)),
        itemStyle: {color: cfg.color},
        lineStyle: {width: 2, color: cfg.color},
        areaStyle: {opacity: 0.08, color: cfg.color},
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
    () => [props.data, props.metricKey, props.title, activeMetric.value, locale.value],
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

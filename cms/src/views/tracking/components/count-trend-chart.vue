<template>
  <div ref="chartRef" class="count-trend-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {TrackingEventTrendPoint} from '@/types/api'

const props = defineProps<{
  data: TrackingEventTrendPoint[]
  title?: string
  seriesName: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const buildOption = (points: TrackingEventTrendPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const counts = points.map(item => Number(item.count || 0))

  return {
    title: props.title
        ? {
          text: props.title,
          left: 'center',
          textStyle: {fontSize: 14, fontWeight: 500},
        }
        : undefined,
    tooltip: {trigger: 'axis'},
    legend: {
      data: [props.seriesName],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 48,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 10 ? 35 : 0},
    },
    yAxis: {type: 'value', minInterval: 1},
    series: [
      {
        name: props.seriesName,
        type: 'line',
        smooth: true,
        data: counts,
        itemStyle: {color: '#409EFF'},
      },
    ],
  }
}

const renderChart = () => {
  if (!chartRef.value) {
    return
  }
  if (!chartInstance) {
    chartInstance = echarts.init(chartRef.value)
  }
  chartInstance.setOption(buildOption(props.data || []), true)
}

watch(
    () => props.data,
    () => {
      void nextTick(renderChart)
    },
    {deep: true},
)

onMounted(() => {
  void nextTick(renderChart)
  window.addEventListener('resize', renderChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', renderChart)
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({resize: renderChart})
</script>

<style scoped>
.count-trend-chart {
  width: 100%;
  height: 360px;
}
</style>

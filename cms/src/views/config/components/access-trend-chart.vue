<template>
  <div ref="chartRef" class="access-trend-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import * as echarts from 'echarts'
import type {AccessTrendData} from '@/types/api'

const props = defineProps<{
  data: AccessTrendData | null
  loading?: boolean
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const buildOption = (trend: AccessTrendData | null): echarts.EChartsOption => {
  const points = trend?.points || []
  const times = points.map(item => item.time)
  const counts = points.map(item => Number(item.count || 0))
  const intervalMinutes = trend?.intervalMinutes || 1
  const peakIndex = points.findIndex(item => item.time === trend?.peakTime)

  return {
    title: {
      text: `访问量趋势 (每${intervalMinutes}分钟)`,
      left: 'center',
      textStyle: {fontSize: 14, fontWeight: 500},
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return [
          point?.time || items[0].name,
          `访问量: ${point?.count ?? 0}`,
        ].join('<br/>')
      },
    },
    grid: {
      left: 56,
      right: 24,
      top: 56,
      bottom: 48,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 16 ? 35 : 0},
    },
    yAxis: {
      type: 'value',
      name: '访问量',
      min: 0,
      minInterval: 1,
    },
    series: [
      {
        name: '访问量',
        type: 'line',
        smooth: true,
        showSymbol: times.length <= 48,
        data: counts,
        areaStyle: {opacity: 0.08},
        markPoint: peakIndex >= 0 ? {
          symbol: 'pin',
          symbolSize: 48,
          data: [{
            name: '峰值',
            coord: [peakIndex, counts[peakIndex]],
            value: counts[peakIndex],
          }],
        } : undefined,
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
  chartInstance.setOption(buildOption(props.data), true)
}

const resize = () => {
  chartInstance?.resize()
}

watch(
    () => [props.data, props.loading],
    async () => {
      await nextTick()
      renderChart()
    },
    {deep: true},
)

onMounted(async () => {
  await nextTick()
  renderChart()
  window.addEventListener('resize', resize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({resize})
</script>

<style scoped>
.access-trend-chart {
  width: 100%;
  height: 320px;
}
</style>

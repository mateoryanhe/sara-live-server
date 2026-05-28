<template>
  <div ref="chartRef" class="user-stat-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import * as echarts from 'echarts'
import type {UserStatTrendPoint} from '@/types/api'

const props = defineProps<{
  data: UserStatTrendPoint[]
  title?: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const buildOption = (points: UserStatTrendPoint[]): echarts.EChartsOption => {
  const times = points.map(item => item.time)
  const activeSeries = points.map(item => Number(item.activeUserCount || 0))
  const registerSeries = points.map(item => Number(item.registerUserCount || 0))

  return {
    title: props.title
        ? {
          text: props.title,
          left: 'center',
          textStyle: {
            fontSize: 14,
            fontWeight: 500,
          },
        }
        : undefined,
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: ['活跃用户数', '新注册用户数'],
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
        name: '活跃用户数',
        type: 'line',
        smooth: true,
        data: activeSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: '新注册用户数',
        type: 'line',
        smooth: true,
        data: registerSeries,
        itemStyle: {color: '#67C23A'},
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
    () => props.data,
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
.user-stat-chart {
  width: 100%;
  height: 360px;
}
</style>

<template>
  <div ref="chartRef" class="game-stats-trend-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {GameBetUserTrendPoint} from '@/types/api'
import {buildTrendChartXAxisLabel, trendChartGridBottom} from '@/utils/trend-chart-zoom'

const props = defineProps<{
  data: GameBetUserTrendPoint[]
  title?: string
  betSeriesName: string
  onlineSeriesName: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const buildOption = (points: GameBetUserTrendPoint[]): EChartsOption => {
  const times = points.map((item) => item.time)
  const betCounts = points.map((item) => Number(item.count || 0))
  const onlineCounts = points.map((item) => Number(item.onlineCount || 0))
  const gridBottom = trendChartGridBottom(times.length, false)

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
      data: [props.betSeriesName, props.onlineSeriesName],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 48,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: gridBottom,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: buildTrendChartXAxisLabel(times.length),
    },
    yAxis: {type: 'value', minInterval: 1},
    series: [
      {
        name: props.betSeriesName,
        type: 'line',
        smooth: true,
        data: betCounts,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: props.onlineSeriesName,
        type: 'line',
        smooth: true,
        data: onlineCounts,
        itemStyle: {color: '#67C23A'},
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

watch(
    () => [props.betSeriesName, props.onlineSeriesName, props.title],
    () => {
      void nextTick(renderChart)
    },
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
.game-stats-trend-chart {
  width: 100%;
  height: 400px;
}
</style>

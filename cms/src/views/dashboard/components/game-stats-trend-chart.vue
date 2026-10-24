<template>
  <div ref="chartRef" class="game-stats-trend-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {GameBetUserTrendPoint} from '@/types/api'
import {buildTrendChartXAxisLabel, trendChartGridBottom} from '@/utils/trend-chart-zoom'

/** 点太多时按桶压缩，保证仍覆盖整段 3 天并铺满 X 轴 */
const MAX_DISPLAY_POINTS = 360
const ONLINE_COLOR = '#7c3aed'
const BET_COLOR = '#e11d48'

const props = defineProps<{
  data: GameBetUserTrendPoint[]
  title?: string
  betSeriesName: string
  onlineSeriesName: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null
let resizeObserver: ResizeObserver | null = null

/** 点少原样铺满；点多按时间比例压成固定桶数，整段仍铺满 X 轴 */
function toDisplayPoints(points: GameBetUserTrendPoint[]): GameBetUserTrendPoint[] {
  if (!points.length || points.length <= MAX_DISPLAY_POINTS) {
    return points
  }
  const bucketSize = points.length / MAX_DISPLAY_POINTS
  const out: GameBetUserTrendPoint[] = []
  for (let i = 0; i < MAX_DISPLAY_POINTS; i++) {
    const start = Math.floor(i * bucketSize)
    const end = Math.max(start + 1, Math.floor((i + 1) * bucketSize))
    let countMax = 0
    let onlineMax = 0
    let mid = points[start]
    for (let j = start; j < end && j < points.length; j++) {
      const row = points[j]
      countMax = Math.max(countMax, Number(row.count || 0))
      onlineMax = Math.max(onlineMax, Number(row.onlineCount || 0))
      if (j === Math.floor((start + end - 1) / 2)) {
        mid = row
      }
    }
    out.push({
      time: mid.time,
      count: countMax,
      onlineCount: onlineMax,
    })
  }
  return out
}

const buildOption = (rawPoints: GameBetUserTrendPoint[]): EChartsOption => {
  const points = toDisplayPoints(rawPoints || [])
  const times = points.map((item) => item.time)
  const betCounts = points.map((item) => Number(item.count || 0))
  const onlineCounts = points.map((item) => Number(item.onlineCount || 0))
  const gridBottom = trendChartGridBottom(times.length, false)

  return {
    animation: times.length < 500,
    title: props.title
        ? {
          text: props.title,
          left: 'center',
          textStyle: {fontSize: 14, fontWeight: 600, color: '#334155'},
        }
        : undefined,
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.92)',
      borderWidth: 0,
      textStyle: {color: '#f8fafc', fontSize: 12},
      axisPointer: {type: 'line', lineStyle: {color: '#94a3b8', type: 'dashed'}},
    },
    legend: {
      data: [props.onlineSeriesName, props.betSeriesName],
      top: props.title ? 30 : 4,
      itemWidth: 14,
      itemHeight: 8,
      textStyle: {color: '#64748b', fontSize: 12},
    },
    grid: {
      left: 52,
      right: 28,
      top: props.title ? 72 : 44,
      bottom: gridBottom,
      containLabel: false,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: buildTrendChartXAxisLabel(times.length),
      axisLine: {lineStyle: {color: '#e2e8f0'}},
      axisTick: {show: false},
      splitLine: {show: false},
    },
    yAxis: {
      type: 'value',
      min: 0,
      minInterval: 1,
      axisLabel: {color: '#94a3b8', fontSize: 11},
      splitLine: {lineStyle: {color: '#f1f5f9', type: 'dashed'}},
    },
    series: [
      {
        name: props.onlineSeriesName,
        type: 'line',
        smooth: true,
        showSymbol: times.length <= 48,
        symbolSize: 6,
        data: onlineCounts,
        itemStyle: {color: ONLINE_COLOR},
        lineStyle: {width: 2.5, color: ONLINE_COLOR},
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {offset: 0, color: 'rgba(124, 58, 237, 0.28)'},
            {offset: 1, color: 'rgba(124, 58, 237, 0.02)'},
          ]),
        },
      },
      {
        name: props.betSeriesName,
        type: 'line',
        smooth: true,
        showSymbol: times.length <= 48,
        symbolSize: 6,
        data: betCounts,
        itemStyle: {color: BET_COLOR},
        lineStyle: {width: 2.5, color: BET_COLOR},
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {offset: 0, color: 'rgba(225, 29, 72, 0.22)'},
            {offset: 1, color: 'rgba(225, 29, 72, 0.02)'},
          ]),
        },
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
  chartInstance.resize()
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
  if (chartRef.value && typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => {
      chartInstance?.resize()
    })
    resizeObserver.observe(chartRef.value)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', renderChart)
  resizeObserver?.disconnect()
  resizeObserver = null
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({resize: renderChart})
</script>

<style scoped>
.game-stats-trend-chart {
  width: 100%;
  height: 460px;
}
</style>

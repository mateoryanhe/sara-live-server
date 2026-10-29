<template>
  <div ref="chartRef" class="minute-live-activity-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {MinuteLiveActivityTrendPoint} from '@/types/api'
import {buildTrendChartXAxisLabel, trendChartGridBottom} from '@/utils/trend-chart-zoom'

/** 与仪表盘游戏统计一致：点多按桶压缩，仍覆盖整段 3 天 */
const MAX_DISPLAY_POINTS = 360

const props = defineProps<{
  data: MinuteLiveActivityTrendPoint[]
  title?: string
}>()

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null
let resizeObserver: ResizeObserver | null = null

const seriesDefs = () => [
  {key: 'joinCount' as const, label: t('pages.trackingMinuteLiveActivity.join'), color: '#409EFF'},
  {key: 'leaveCount' as const, label: t('pages.trackingMinuteLiveActivity.leave'), color: '#67C23A'},
  {key: 'giftCount' as const, label: t('pages.trackingMinuteLiveActivity.gift'), color: '#E6A23C'},
  {key: 'chatCount' as const, label: t('pages.trackingMinuteLiveActivity.chat'), color: '#64748b'},
  {key: 'betCount' as const, label: t('pages.trackingMinuteLiveActivity.bet'), color: '#F56C6C'},
  {key: 'call1v1Count' as const, label: t('pages.trackingMinuteLiveActivity.call1v1'), color: '#9C27B0'},
  {key: 'rechargeOrderCount' as const, label: t('pages.trackingMinuteLiveActivity.rechargeOrder'), color: '#00BCD4'},
]

function toDisplayPoints(points: MinuteLiveActivityTrendPoint[]): MinuteLiveActivityTrendPoint[] {
  if (!points.length || points.length <= MAX_DISPLAY_POINTS) {
    return points
  }
  const bucketSize = points.length / MAX_DISPLAY_POINTS
  const out: MinuteLiveActivityTrendPoint[] = []
  for (let i = 0; i < MAX_DISPLAY_POINTS; i++) {
    const start = Math.floor(i * bucketSize)
    const end = Math.max(start + 1, Math.floor((i + 1) * bucketSize))
    let mid = points[start]
    const agg: MinuteLiveActivityTrendPoint = {
      time: mid.time,
      joinCount: 0,
      leaveCount: 0,
      giftCount: 0,
      chatCount: 0,
      betCount: 0,
      call1v1Count: 0,
      rechargeOrderCount: 0,
    }
    for (let j = start; j < end && j < points.length; j++) {
      const row = points[j]
      agg.joinCount = Math.max(agg.joinCount, Number(row.joinCount || 0))
      agg.leaveCount = Math.max(agg.leaveCount, Number(row.leaveCount || 0))
      agg.giftCount = Math.max(agg.giftCount, Number(row.giftCount || 0))
      agg.chatCount = Math.max(agg.chatCount, Number(row.chatCount || 0))
      agg.betCount = Math.max(agg.betCount, Number(row.betCount || 0))
      agg.call1v1Count = Math.max(agg.call1v1Count, Number(row.call1v1Count || 0))
      agg.rechargeOrderCount = Math.max(agg.rechargeOrderCount, Number(row.rechargeOrderCount || 0))
      if (j === Math.floor((start + end - 1) / 2)) {
        mid = row
      }
    }
    agg.time = mid.time
    out.push(agg)
  }
  return out
}

const buildOption = (rawPoints: MinuteLiveActivityTrendPoint[]): EChartsOption => {
  const points = toDisplayPoints(rawPoints || [])
  const times = points.map(item => item.time)
  const defs = seriesDefs()
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
      data: defs.map(item => item.label),
      top: props.title ? 30 : 4,
      itemWidth: 14,
      itemHeight: 8,
      textStyle: {color: '#64748b', fontSize: 12},
    },
    grid: {
      left: 52,
      right: 28,
      top: props.title ? 88 : 56,
      bottom: gridBottom,
      containLabel: false,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {
        ...buildTrendChartXAxisLabel(times.length),
        color: '#1e293b',
        fontSize: 12,
        fontWeight: 600,
      },
      axisLine: {lineStyle: {color: '#94a3b8', width: 1}},
      axisTick: {show: true, alignWithLabel: true, lineStyle: {color: '#64748b'}},
      splitLine: {show: false},
    },
    yAxis: {
      type: 'value',
      min: 0,
      minInterval: 1,
      axisLabel: {color: '#1e293b', fontSize: 12, fontWeight: 600},
      axisLine: {show: true, lineStyle: {color: '#94a3b8'}},
      axisTick: {show: true, lineStyle: {color: '#64748b'}},
      splitLine: {lineStyle: {color: '#e2e8f0', type: 'dashed'}},
    },
    series: defs.map(item => ({
      name: item.label,
      type: 'line',
      smooth: true,
      showSymbol: times.length <= 48,
      symbolSize: 6,
      data: points.map(p => Number(p[item.key] || 0)),
      itemStyle: {color: item.color},
      lineStyle: {width: 2, color: item.color},
    })),
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

watch(() => [props.data, props.title, locale.value], () => {
  void nextTick(renderChart)
}, {deep: true})

onMounted(() => {
  void nextTick(renderChart)
  window.addEventListener('resize', resize)
  if (chartRef.value && typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => chartInstance?.resize())
    resizeObserver.observe(chartRef.value)
  }
})

const resize = () => {
  chartInstance?.resize()
}

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  resizeObserver?.disconnect()
  resizeObserver = null
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({resize})
</script>

<style scoped>
.minute-live-activity-chart {
  width: 100%;
  height: 440px;
}
</style>

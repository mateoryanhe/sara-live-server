<template>
  <div ref="chartRef" class="access-trend-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {AccessTrendData} from '@/types/api'

const props = defineProps<{
  data: AccessTrendData | null
  loading?: boolean
}>()

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const buildOption = (trend: AccessTrendData | null): EChartsOption => {
  const points = trend?.points || []
  const times = points.map(item => item.time)
  const counts = points.map(item => Number(item.count || 0))
  const intervalMinutes = trend?.intervalMinutes || 1
  const peakIndex = points.findIndex(item => item.time === trend?.peakTime)

  return {
    title: {
      text: t('pages.resourceMonitor.accessTrendTitle', {minutes: intervalMinutes}),
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
          t('pages.resourceMonitor.accessCountTooltip', {value: point?.count ?? 0}),
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
      name: t('pages.resourceMonitor.accessCount'),
      min: 0,
      minInterval: 1,
    },
    series: [
      {
        name: t('pages.resourceMonitor.accessCount'),
        type: 'line',
        smooth: true,
        showSymbol: times.length <= 48,
        data: counts,
        areaStyle: {opacity: 0.08},
        markPoint: peakIndex >= 0 ? {
          symbol: 'pin',
          symbolSize: 48,
          data: [{
            name: t('pages.resourceMonitor.peak'),
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
    () => [props.data, props.loading, locale.value],
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

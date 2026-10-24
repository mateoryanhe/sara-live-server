<template>
  <div ref="chartRef" class="client-login-trend-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {UserStatTrendPoint} from '@/types/api'
import {buildTrendChartXAxisLabel, trendChartGridBottom} from '@/utils/trend-chart-zoom'

const props = defineProps<{
  data: UserStatTrendPoint[]
  title?: string
}>()

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const buildOption = (points: UserStatTrendPoint[]): EChartsOption => {
  const times = points.map((item) => item.time)
  const appSeries = points.map((item) => Number(item.appDeviceLoginCount || 0))
  const h5Series = points.map((item) => Number(item.h5DeviceLoginCount || 0))
  const gridBottom = trendChartGridBottom(times.length, false)
  const appName = t('pages.dashboard.clientLoginApp')
  const h5Name = t('pages.dashboard.clientLoginH5')

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
      data: [appName, h5Name],
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
    yAxis: {
      type: 'value',
      minInterval: 1,
    },
    series: [
      {
        name: appName,
        type: 'line',
        smooth: true,
        data: appSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: h5Name,
        type: 'line',
        smooth: true,
        data: h5Series,
        itemStyle: {color: '#E6A23C'},
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
    () => [props.data, props.title, locale.value],
    () => {
      void renderChart()
    },
    {deep: true},
)

onMounted(() => {
  void renderChart()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({
  resize: handleResize,
})
</script>

<style scoped>
.client-login-trend-chart {
  width: 100%;
  height: 360px;
}
</style>

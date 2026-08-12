<template>
  <div ref="chartRef" class="user-stat-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {UserStatTrendPoint} from '@/types/api'

const props = defineProps<{
  data: UserStatTrendPoint[]
  title?: string
}>()

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const buildOption = (points: UserStatTrendPoint[]): EChartsOption => {
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
      data: [t('pages.dashboard.activeUsers'), t('pages.dashboard.registerUsers')],
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
        name: t('pages.dashboard.activeUsers'),
        type: 'line',
        smooth: true,
        data: activeSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: t('pages.dashboard.registerUsers'),
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
    () => [props.data, props.title, locale.value],
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

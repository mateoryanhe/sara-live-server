<template>
  <div ref="chartRef" class="recharge-type-trend-chart"></div>
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
  const googleSeries = points.map((item) => Number(item.googleRechargeAmount || 0))
  const haiPaySeries = points.map((item) => Number(item.haiPayRechargeAmount || 0))
  const gridBottom = trendChartGridBottom(times.length, false)
  const googleName = t('pages.dashboard.rechargeTypeGoogle')
  const haiPayName = t('pages.dashboard.rechargeTypeHaiPay')

  return {
    title: props.title
        ? {
          text: props.title,
          left: 'center',
          textStyle: {fontSize: 14, fontWeight: 500},
        }
        : undefined,
    tooltip: {
      trigger: 'axis',
      valueFormatter: (value) => {
        const n = Number(value)
        if (!Number.isFinite(n)) {
          return '-'
        }
        return `${n.toLocaleString(locale.value, {maximumFractionDigits: 4})} USD`
      },
    },
    legend: {
      data: [googleName, haiPayName],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
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
      minInterval: 0.01,
    },
    series: [
      {
        name: googleName,
        type: 'line',
        smooth: true,
        data: googleSeries,
        itemStyle: {color: '#67C23A'},
      },
      {
        name: haiPayName,
        type: 'line',
        smooth: true,
        data: haiPaySeries,
        itemStyle: {color: '#F56C6C'},
      },
    ],
  }
}

const renderChart = async () => {
  await nextTick()
  if (!chartRef.value) {
    return
  }
  const {clientWidth, clientHeight} = chartRef.value
  if (clientWidth === 0 || clientHeight === 0) {
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
  resize: () => {
    void renderChart()
  },
})
</script>

<style scoped>
.recharge-type-trend-chart {
  width: 100%;
  height: 420px;
}
</style>
<template>
  <div ref="chartRef" class="resource-metric-chart"></div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import type {EChartsType} from 'echarts/core'
import echarts, {type EChartsOption} from '@/utils/echarts'
import type {ResourceMetricPoint} from '@/types/api'
import {
  buildTrendChartDataZoom,
  buildTrendChartXAxisLabel,
  buildTrendChartXAxisLabelZoomFriendly,
  TREND_CHART_DENSE_VISIBLE_POINTS,
  trendChartGridBottom,
} from '@/utils/trend-chart-zoom'

type MetricType = 'memory' | 'heap' | 'ratio' | 'cpu' | 'online'

const props = defineProps<{
  data: ResourceMetricPoint[]
  metricType: MetricType
  title?: string
  /** 是否显示底部趋势滑块（仪表盘在线人数关闭） */
  showTrendDataZoom?: boolean
  /** 滑块初始拉满显示全量 */
  dataZoomDefaultFullRange?: boolean
  /** 在线趋势：观众 + 开播主播两条线（仪表盘）；默认仅在线总数一条线（资源监控页） */
  onlineAudienceAnchorOnly?: boolean
}>()

const showTrendDataZoom = computed(() => props.showTrendDataZoom !== false)

const {t, locale} = useI18n()
const chartRef = ref<HTMLDivElement>()
let chartInstance: EChartsType | null = null

const formatNumber = (value: number | null | undefined, digits = 2) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(digits)
}

const formatCount = (value: number | string | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString()
}

const buildMemoryOption = (points: ResourceMetricPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const procMemSeries = points.map(item => Number(item.procMemMb || 0))
  const sysMemSeries = points.map(item => Number(item.sysMemUsedMb || 0))

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
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return [
          items[0].name,
          t('pages.resourceMonitor.procRssTooltip', {value: formatNumber(point?.procMemMb)}),
          t('pages.resourceMonitor.sysMemTooltip', {
            used: formatNumber(point?.sysMemUsedMb),
            total: formatNumber(point?.sysMemTotalMb),
            percent: formatNumber(point?.sysMemUsedPercent),
          }),
        ].join('<br/>')
      },
    },
    legend: {
      data: [t('pages.resourceMonitor.procRss'), t('pages.resourceMonitor.sysMemUsed')],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 12 ? 35 : 0},
    },
    yAxis: {
      type: 'value',
      name: t('pages.resourceMonitor.memoryMb'),
      min: 0,
    },
    series: [
      {
        name: t('pages.resourceMonitor.procRss'),
        type: 'line',
        smooth: true,
        data: procMemSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: t('pages.resourceMonitor.sysMemUsed'),
        type: 'line',
        smooth: true,
        data: sysMemSeries,
        itemStyle: {color: '#67C23A'},
      },
    ],
  }
}

const buildHeapOption = (points: ResourceMetricPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const heapAllocSeries = points.map(item => Number(item.procHeapAllocMb || 0))
  const heapInuseSeries = points.map(item => Number(item.procHeapInuseMb || 0))
  const heapSysSeries = points.map(item => Number(item.procHeapSysMb || 0))

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
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return [
          items[0].name,
          t('pages.resourceMonitor.heapAllocTooltip', {value: formatNumber(point?.procHeapAllocMb)}),
          t('pages.resourceMonitor.heapInuseTooltip', {value: formatNumber(point?.procHeapInuseMb)}),
          t('pages.resourceMonitor.heapSysTooltip', {value: formatNumber(point?.procHeapSysMb)}),
        ].join('<br/>')
      },
    },
    legend: {
      data: [
        t('pages.resourceMonitor.heapAlloc'),
        t('pages.resourceMonitor.heapInuse'),
        t('pages.resourceMonitor.heapSys'),
      ],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 12 ? 35 : 0},
    },
    yAxis: {
      type: 'value',
      name: t('pages.resourceMonitor.heapMemoryMb'),
      min: 0,
    },
    series: [
      {
        name: t('pages.resourceMonitor.heapAlloc'),
        type: 'line',
        smooth: true,
        data: heapAllocSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: t('pages.resourceMonitor.heapInuse'),
        type: 'line',
        smooth: true,
        data: heapInuseSeries,
        itemStyle: {color: '#E6A23C'},
      },
      {
        name: t('pages.resourceMonitor.heapSys'),
        type: 'line',
        smooth: true,
        data: heapSysSeries,
        itemStyle: {color: '#909399'},
      },
    ],
  }
}

const buildHeapRatioOption = (points: ResourceMetricPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const heapUsedPercentSeries = points.map(item => Number(item.procHeapUsedPercent || 0))
  const heapIdlePercentSeries = points.map(item => Number(item.procHeapIdlePercent || 0))

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
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return [
          items[0].name,
          t('pages.resourceMonitor.heapUsedPercentTooltip', {value: formatNumber(point?.procHeapUsedPercent)}),
          t('pages.resourceMonitor.heapIdlePercentTooltip', {value: formatNumber(point?.procHeapIdlePercent)}),
        ].join('<br/>')
      },
    },
    legend: {
      data: [t('pages.resourceMonitor.heapUsedPercent'), t('pages.resourceMonitor.heapIdlePercent')],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 12 ? 35 : 0},
    },
    yAxis: {
      type: 'value',
      name: t('pages.resourceMonitor.ratioPercent'),
      min: 0,
      max: 100,
    },
    series: [
      {
        name: t('pages.resourceMonitor.heapUsedPercent'),
        type: 'line',
        smooth: true,
        data: heapUsedPercentSeries,
        itemStyle: {color: '#67C23A'},
      },
      {
        name: t('pages.resourceMonitor.heapIdlePercent'),
        type: 'line',
        smooth: true,
        data: heapIdlePercentSeries,
        itemStyle: {color: '#F56C6C'},
      },
    ],
  }
}

const resolveAudienceCount = (point: ResourceMetricPoint | undefined) => {
  if (!point) {
    return 0
  }
  if (point.audienceCount !== undefined && point.audienceCount !== '') {
    return Number(point.audienceCount)
  }
  const total = Number(point.onlineCount || 0)
  const anchors = Number(point.liveAnchorCount || 0)
  return total > anchors ? total - anchors : 0
}

const buildOnlineOption = (points: ResourceMetricPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const anchorSeries = points.map(item => Number(item.liveAnchorCount || 0))
  const audienceSeries = points.map(item => resolveAudienceCount(item))
  const audienceAnchorOnly = props.onlineAudienceAnchorOnly === true
  const totalLabel = t('pages.resourceMonitor.onlineCount')
  const anchorLabel = t('pages.resourceMonitor.liveAnchorOnlineCount')
  const audienceLabel = t('pages.resourceMonitor.audienceOnlineCount')
  const withSlider = showTrendDataZoom.value
  const gridBottom = trendChartGridBottom(
      withSlider ? Math.min(times.length, TREND_CHART_DENSE_VISIBLE_POINTS) : times.length,
      withSlider,
  )

  const legendData = audienceAnchorOnly
      ? [audienceLabel, anchorLabel]
      : [totalLabel]

  const tooltipLines = (point: ResourceMetricPoint | undefined, axisName: string) => {
    const lines = [axisName]
    if (audienceAnchorOnly) {
      lines.push(
          t('pages.resourceMonitor.audienceOnlineCountTooltip', {value: formatCount(resolveAudienceCount(point))}),
          t('pages.resourceMonitor.liveAnchorOnlineCountTooltip', {value: formatCount(point?.liveAnchorCount)}),
      )
    } else {
      lines.push(t('pages.resourceMonitor.onlineCountTooltip', {value: formatCount(point?.onlineCount)}))
    }
    return lines.join('<br/>')
  }

  const audienceLineColor = '#dc2626'
  const anchorLineColor = '#2563eb'

  const series: EChartsOption['series'] = audienceAnchorOnly
      ? [
        {
          name: audienceLabel,
          type: 'line',
          smooth: true,
          data: audienceSeries,
          itemStyle: {color: audienceLineColor},
          lineStyle: {color: audienceLineColor},
        },
        {
          name: anchorLabel,
          type: 'line',
          smooth: true,
          data: anchorSeries,
          itemStyle: {color: anchorLineColor},
          lineStyle: {color: anchorLineColor},
        },
      ]
      : [
        {
          name: totalLabel,
          type: 'line',
          smooth: true,
          data: points.map(item => Number(item.onlineCount || 0)),
          itemStyle: {color: '#409EFF'},
        },
      ]

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
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return tooltipLines(point, String(items[0].name ?? ''))
      },
    },
    legend: {
      data: legendData,
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
      right: 24,
      top: props.title ? 88 : 64,
      bottom: gridBottom,
    },
    dataZoom: withSlider
        ? buildTrendChartDataZoom(times.length, {
          defaultFullRange: props.dataZoomDefaultFullRange,
          alignEnd: props.dataZoomDefaultFullRange !== true,
          initialWindowPoints: TREND_CHART_DENSE_VISIBLE_POINTS,
        })
        : undefined,
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: withSlider
          ? buildTrendChartXAxisLabelZoomFriendly()
          : buildTrendChartXAxisLabel(times.length),
    },
    yAxis: {
      type: 'value',
      name: t('pages.resourceMonitor.peopleCount'),
      min: 0,
      minInterval: 1,
    },
    series,
  }
}

const buildCpuOption = (points: ResourceMetricPoint[]): EChartsOption => {
  const times = points.map(item => item.time)
  const procCpuSeries = points.map(item => Number(item.procCpuPercent || 0))
  const sysCpuSeries = points.map(item => Number(item.sysCpuPercent || 0))

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
      formatter: (params) => {
        const items = Array.isArray(params) ? params : [params]
        if (!items.length) {
          return ''
        }
        const index = items[0].dataIndex ?? 0
        const point = points[index]
        return [
          items[0].name,
          t('pages.resourceMonitor.procCpuTooltip', {value: formatNumber(point?.procCpuPercent)}),
          t('pages.resourceMonitor.sysCpuTooltip', {value: formatNumber(point?.sysCpuPercent)}),
        ].join('<br/>')
      },
    },
    legend: {
      data: [t('pages.resourceMonitor.procCpu'), t('pages.resourceMonitor.sysCpu')],
      top: props.title ? 28 : 0,
    },
    grid: {
      left: 56,
      right: 24,
      top: props.title ? 72 : 48,
      bottom: 32,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: times,
      axisLabel: {rotate: times.length > 12 ? 35 : 0},
    },
    yAxis: {
      type: 'value',
      name: t('pages.resourceMonitor.cpuPercent'),
      min: 0,
      max: 100,
    },
    series: [
      {
        name: t('pages.resourceMonitor.procCpu'),
        type: 'line',
        smooth: true,
        data: procCpuSeries,
        itemStyle: {color: '#E6A23C'},
      },
      {
        name: t('pages.resourceMonitor.sysCpu'),
        type: 'line',
        smooth: true,
        data: sysCpuSeries,
        itemStyle: {color: '#F56C6C'},
      },
    ],
  }
}

const buildOption = (points: ResourceMetricPoint[]): EChartsOption => {
  if (props.metricType === 'cpu') {
    return buildCpuOption(points)
  }
  if (props.metricType === 'online') {
    return buildOnlineOption(points)
  }
  if (props.metricType === 'heap') {
    return buildHeapOption(points)
  }
  if (props.metricType === 'ratio') {
    return buildHeapRatioOption(points)
  }
  return buildMemoryOption(points)
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
    () => [props.data, props.metricType, props.title, props.onlineAudienceAnchorOnly, locale.value],
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
.resource-metric-chart {
  width: 100%;
  height: 420px;
}
</style>

<template>
  <div ref="chartRef" class="resource-metric-chart"></div>
</template>

<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import * as echarts from 'echarts'
import type {ResourceMetricPoint} from '@/types/api'

type MetricType = 'memory' | 'heap' | 'ratio' | 'cpu' | 'online'

const props = defineProps<{
  data: ResourceMetricPoint[]
  metricType: MetricType
  title?: string
}>()

const chartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const formatNumber = (value: number | null | undefined, digits = 2) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toFixed(digits)
}

const buildMemoryOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
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
          `进程RSS: ${formatNumber(point?.procMemMb)} MB`,
          `系统已用内存: ${formatNumber(point?.sysMemUsedMb)} MB / ${formatNumber(point?.sysMemTotalMb)} MB (${formatNumber(point?.sysMemUsedPercent)}%)`,
        ].join('<br/>')
      },
    },
    legend: {
      data: ['进程RSS', '系统已用内存'],
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
      name: '内存(MB)',
      min: 0,
    },
    series: [
      {
        name: '进程RSS',
        type: 'line',
        smooth: true,
        data: procMemSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: '系统已用内存',
        type: 'line',
        smooth: true,
        data: sysMemSeries,
        itemStyle: {color: '#67C23A'},
      },
    ],
  }
}

const buildHeapOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
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
          `堆已分配: ${formatNumber(point?.procHeapAllocMb)} MB`,
          `堆使用中: ${formatNumber(point?.procHeapInuseMb)} MB`,
          `堆系统占用: ${formatNumber(point?.procHeapSysMb)} MB`,
        ].join('<br/>')
      },
    },
    legend: {
      data: ['堆已分配', '堆使用中', '堆系统占用'],
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
      name: '堆内存(MB)',
      min: 0,
    },
    series: [
      {
        name: '堆已分配',
        type: 'line',
        smooth: true,
        data: heapAllocSeries,
        itemStyle: {color: '#409EFF'},
      },
      {
        name: '堆使用中',
        type: 'line',
        smooth: true,
        data: heapInuseSeries,
        itemStyle: {color: '#E6A23C'},
      },
      {
        name: '堆系统占用',
        type: 'line',
        smooth: true,
        data: heapSysSeries,
        itemStyle: {color: '#909399'},
      },
    ],
  }
}

const buildHeapRatioOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
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
          `堆使用比例: ${formatNumber(point?.procHeapUsedPercent)}%`,
          `堆空闲比例: ${formatNumber(point?.procHeapIdlePercent)}%`,
        ].join('<br/>')
      },
    },
    legend: {
      data: ['堆使用比例', '堆空闲比例'],
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
      name: '比例(%)',
      min: 0,
      max: 100,
    },
    series: [
      {
        name: '堆使用比例',
        type: 'line',
        smooth: true,
        data: heapUsedPercentSeries,
        itemStyle: {color: '#67C23A'},
      },
      {
        name: '堆空闲比例',
        type: 'line',
        smooth: true,
        data: heapIdlePercentSeries,
        itemStyle: {color: '#F56C6C'},
      },
    ],
  }
}

const formatCount = (value: number | string | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString('zh-CN')
}

const buildOnlineOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
  const times = points.map(item => item.time)
  const onlineSeries = points.map(item => Number(item.onlineCount || 0))

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
          `在线人数: ${formatCount(point?.onlineCount)}`,
        ].join('<br/>')
      },
    },
    legend: {
      data: ['在线人数'],
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
      name: '人数',
      min: 0,
      minInterval: 1,
    },
    series: [
      {
        name: '在线人数',
        type: 'line',
        smooth: true,
        data: onlineSeries,
        itemStyle: {color: '#409EFF'},
      },
    ],
  }
}

const buildCpuOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
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
          `进程CPU: ${formatNumber(point?.procCpuPercent)}%`,
          `系统CPU: ${formatNumber(point?.sysCpuPercent)}%`,
        ].join('<br/>')
      },
    },
    legend: {
      data: ['进程CPU', '系统CPU'],
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
      name: 'CPU(%)',
      min: 0,
      max: 100,
    },
    series: [
      {
        name: '进程CPU',
        type: 'line',
        smooth: true,
        data: procCpuSeries,
        itemStyle: {color: '#E6A23C'},
      },
      {
        name: '系统CPU',
        type: 'line',
        smooth: true,
        data: sysCpuSeries,
        itemStyle: {color: '#F56C6C'},
      },
    ],
  }
}

const buildOption = (points: ResourceMetricPoint[]): echarts.EChartsOption => {
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
    () => [props.data, props.metricType],
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
  height: 400px;
}
</style>

<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ResourceMonitor') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.startTime')">
          <el-date-picker
              v-model="searchForm.startTime"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('common.startTime')"
              style="width: 200px"
              teleported
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item :label="t('common.endTime')">
          <el-date-picker
              v-model="searchForm.endTime"
              clearable
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('common.endTime')"
              style="width: 200px"
              teleported
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">{{ t('common.query') }}</el-button>
          <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
        <span class="search-tip">{{ t('pages.resourceMonitor.searchTip', {max: RESOURCE_METRIC_MAX_POINTS}) }}</span>
      </el-form>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange" @tab-click="handleTabClick">
        <el-tab-pane :label="t('pages.resourceMonitor.tabMemory')" lazy name="memory">
          <ResourceMetricChart
              ref="memoryChartRef"
              :data="tabPoints.memory"
              metric-type="memory"
              :title="chartTitle(t('pages.resourceMonitor.chartProcRssSysMem'))"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.resourceMonitor.tabHeap')" lazy name="heap">
          <ResourceMetricChart
              ref="heapChartRef"
              :data="tabPoints.heap"
              metric-type="heap"
              :title="chartTitle(t('pages.resourceMonitor.chartProcHeap'))"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.resourceMonitor.tabRatio')" lazy name="ratio">
          <ResourceMetricChart
              ref="ratioChartRef"
              :data="tabPoints.ratio"
              metric-type="ratio"
              :title="chartTitle(t('pages.resourceMonitor.chartHeapRatio'))"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.resourceMonitor.tabCpu')" lazy name="cpu">
          <ResourceMetricChart
              ref="cpuChartRef"
              :data="tabPoints.cpu"
              metric-type="cpu"
              :title="chartTitle(t('pages.resourceMonitor.chartProcSysCpu'))"
          />
        </el-tab-pane>
        <el-tab-pane :label="t('pages.resourceMonitor.tabOnline')" lazy name="online">
          <ResourceMetricChart
              ref="onlineChartRef"
              :data="tabPoints.online"
              metric-type="online"
              :title="chartTitle(t('pages.resourceMonitor.chartOnlineCount'))"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {nextTick, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {RESOURCE_METRIC_MAX_POINTS, sysStatApi} from '@/api'
import type {ResourceMetricPoint} from '@/types/api'
import ResourceMetricChart from './components/resource-metric-chart.vue'
import {formatServerDateTimeFromDate} from '@/utils/server-datetime'

type MetricTab = 'memory' | 'heap' | 'ratio' | 'cpu' | 'online'

const {t} = useI18n()
const loading = ref(false)
const activeTab = ref<MetricTab>('memory')
const queryKey = ref('')

const memoryChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const heapChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const ratioChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const cpuChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const onlineChartRef = ref<InstanceType<typeof ResourceMetricChart>>()

const searchForm = reactive({
  startTime: '',
  endTime: '',
})

const tabPoints = reactive<Record<MetricTab, ResourceMetricPoint[]>>({
  memory: [],
  heap: [],
  ratio: [],
  cpu: [],
  online: [],
})

const loadedTabs = ref(new Set<MetricTab>())

const tabFetchers: Record<MetricTab, (params: {
  startTime: string
  endTime: string
  limit: number
}) => Promise<{ points: ResourceMetricPoint[] }>> = {
  memory: sysStatApi.getResourceMetricMemoryTrend,
  heap: sysStatApi.getResourceMetricHeapTrend,
  ratio: sysStatApi.getResourceMetricRatioTrend,
  cpu: sysStatApi.getResourceMetricCpuTrend,
  online: sysStatApi.getResourceMetricOnlineTrend,
}

const pad = (value: number) => String(value).padStart(2, '0')

const formatDateTime = (date: Date) => {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const buildDefaultSearchForm = () => {
  const end = new Date()
  const start = new Date(end.getTime() - 3 * 24 * 60 * 60 * 1000)
  searchForm.startTime = formatDateTime(start)
  searchForm.endTime = formatDateTime(end)
}

const buildQueryKey = () => `${searchForm.startTime}|${searchForm.endTime}`

const buildQueryParams = () => ({
  startTime: searchForm.startTime,
  endTime: searchForm.endTime,
  limit: RESOURCE_METRIC_MAX_POINTS,
})

const chartTitle = (label: string) => {
  if (searchForm.startTime && searchForm.endTime) {
    return t('pages.resourceMonitor.chartTitleRange', {
      label,
      start: searchForm.startTime,
      end: searchForm.endTime,
    })
  }
  return label
}

const clearTabData = () => {
  tabPoints.memory = []
  tabPoints.heap = []
  tabPoints.ratio = []
  tabPoints.cpu = []
  tabPoints.online = []
  loadedTabs.value.clear()
}

const resizeActiveChart = () => {
  if (activeTab.value === 'cpu') {
    cpuChartRef.value?.resize()
    return
  }
  if (activeTab.value === 'online') {
    onlineChartRef.value?.resize()
    return
  }
  if (activeTab.value === 'heap') {
    heapChartRef.value?.resize()
    return
  }
  if (activeTab.value === 'ratio') {
    ratioChartRef.value?.resize()
    return
  }
  memoryChartRef.value?.resize()
}

const fetchTabData = async (tab: MetricTab, force = false) => {
  const currentKey = buildQueryKey()
  if (!force && loadedTabs.value.has(tab) && queryKey.value === currentKey) {
    return
  }
  loading.value = true
  try {
    const data = await tabFetchers[tab](buildQueryParams())
    tabPoints[tab] = data.points || []
    loadedTabs.value.add(tab)
    queryKey.value = currentKey
    await nextTick()
    if (activeTab.value === tab) {
      setTimeout(() => {
        resizeActiveChart()
      }, 0)
    }
  } catch (error) {
    console.error('fetch resource trend failed:', error)
    ElMessage.error(t('pages.resourceMonitor.fetchTrendFailed'))
  } finally {
    loading.value = false
  }
}

const handleQuery = async () => {
  clearTabData()
  await fetchTabData(activeTab.value, true)
}

const resetSearch = () => {
  buildDefaultSearchForm()
  queryKey.value = ''
  clearTabData()
}

const handleTabClick = async (pane: { paneName?: string | number }) => {
  const tab = String(pane.paneName || activeTab.value) as MetricTab
  await fetchTabData(tab)
}

const handleTabChange = async (tabName: string | number) => {
  activeTab.value = String(tabName) as MetricTab
  await nextTick()
  setTimeout(() => {
    resizeActiveChart()
  }, 0)
}

buildDefaultSearchForm()
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-form {
  margin-bottom: 12px;
}

.search-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>

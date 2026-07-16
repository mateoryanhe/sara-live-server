<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统资源监控</span>
          <el-button :loading="loading" @click="fetchResourceMetricTrend">刷新</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="内存" name="memory">
          <ResourceMetricChart
              ref="memoryChartRef"
              :data="resourceMetricTrend.points"
              metric-type="memory"
              title="进程RSS与系统内存 (最近3天, 每10分钟采样)"
          />
        </el-tab-pane>
        <el-tab-pane label="堆内存" name="heap">
          <ResourceMetricChart
              ref="heapChartRef"
              :data="resourceMetricTrend.points"
              metric-type="heap"
              title="进程堆内存 (最近3天, 每10分钟采样)"
          />
        </el-tab-pane>
        <el-tab-pane label="比例" name="ratio">
          <ResourceMetricChart
              ref="ratioChartRef"
              :data="resourceMetricTrend.points"
              metric-type="ratio"
              title="堆使用/空闲比例 (最近3天, 每10分钟采样)"
          />
        </el-tab-pane>
        <el-tab-pane label="CPU" name="cpu">
          <ResourceMetricChart
              ref="cpuChartRef"
              :data="resourceMetricTrend.points"
              metric-type="cpu"
              title="进程/系统CPU (最近3天, 每10分钟采样)"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {sysStatApi} from '@/api'
import type {ResourceMetricTrend} from '@/types/api'
import ResourceMetricChart from './components/resource-metric-chart.vue'

const loading = ref(false)
const activeTab = ref('memory')
const memoryChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const heapChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const ratioChartRef = ref<InstanceType<typeof ResourceMetricChart>>()
const cpuChartRef = ref<InstanceType<typeof ResourceMetricChart>>()

const resourceMetricTrend = reactive<ResourceMetricTrend>({
  points: [],
})

const resizeActiveChart = () => {
  if (activeTab.value === 'cpu') {
    cpuChartRef.value?.resize()
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

const handleTabChange = async () => {
  await nextTick()
  setTimeout(() => {
    resizeActiveChart()
  }, 0)
}

const fetchResourceMetricTrend = async () => {
  loading.value = true
  try {
    const data = await sysStatApi.getResourceMetricTrend()
    resourceMetricTrend.points = data.points || []
    await nextTick()
    setTimeout(() => {
      resizeActiveChart()
    }, 0)
  } catch (error) {
    console.error('获取系统资源趋势失败:', error)
    ElMessage.error('获取系统资源趋势失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchResourceMetricTrend()
})
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
</style>

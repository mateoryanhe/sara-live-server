<template>
  <div class="bar-metric-section">
    <el-tabs
        :model-value="metricKey"
        class="bar-metric-tabs"
        type="card"
        @tab-change="onTabChange"
    >
      <el-tab-pane
          v-for="item in metricTabs"
          :key="item.key"
          :label="item.label"
          :name="item.key"
      />
    </el-tabs>
    <slot/>
  </div>
</template>

<script lang="ts" setup>
import {getUserStatBarMetricTabs} from '../user-stat-bar-series'

defineProps<{
  metricKey: string
}>()

const emit = defineEmits<{
  'update:metricKey': [value: string]
}>()

const metricTabs = getUserStatBarMetricTabs()

const onTabChange = (name: string | number) => {
  emit('update:metricKey', String(name))
}
</script>

<style scoped>
.bar-metric-section {
  margin-top: 24px;
}

.bar-metric-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}
</style>

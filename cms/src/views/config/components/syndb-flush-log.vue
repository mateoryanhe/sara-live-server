<template>
  <div class="syndb-flush-log" :class="{ compact }">
    <div class="syndb-flush-summary">
      <el-tag v-if="flush.reason" :type="syndbFlushReasonTagType(flush.reason) || undefined" size="small">
        {{ formatReason(flush.reason) }}
      </el-tag>
      <span class="syndb-flush-metric">{{ t('pages.dataSync.syndbRows', {count: flush.rows}) }}</span>
      <span class="syndb-flush-metric">{{ t('pages.dataSync.syndbQueues', {count: flush.queues}) }}</span>
      <span class="syndb-flush-metric">{{ flush.costMs }} ms</span>
      <template v-if="!compact">
        <span v-if="flush.batchLimit" class="syndb-flush-metric">{{ t('pages.dataSync.syndbBatchLimit', {count: flush.batchLimit}) }}</span>
        <span v-if="flush.idleQueues" class="syndb-flush-metric">{{ t('pages.dataSync.syndbIdleQueues', {count: flush.idleQueues}) }}</span>
        <span v-if="flush.forceQueues" class="syndb-flush-metric">{{ t('pages.dataSync.syndbForceQueues', {count: flush.forceQueues}) }}</span>
      </template>
      <span v-if="flush.sysCpu !== undefined" class="syndb-flush-metric">{{ t('pages.dataSync.syndbSysCpu', {value: flush.sysCpu.toFixed(1)}) }}</span>
      <span v-if="flush.cpuIdle !== undefined" class="syndb-flush-metric">{{ t('pages.dataSync.syndbCpuIdle', {value: flush.cpuIdle.toFixed(1)}) }}</span>
      <span v-if="flush.idleThreshold !== undefined" class="syndb-flush-metric">{{ t('pages.dataSync.syndbThreshold', {value: flush.idleThreshold.toFixed(0)}) }}</span>
    </div>

    <el-table
        v-if="!compact && flush.details.length"
        :data="flush.details"
        border
        class="syndb-flush-table"
        size="small"
        stripe
    >
      <el-table-column :label="t('pages.dataSync.syndbTable')" min-width="140" prop="table" show-overflow-tooltip/>
      <el-table-column :label="t('pages.dataSync.syndbColumn')" min-width="120" prop="col" show-overflow-tooltip/>
      <el-table-column align="right" :label="t('pages.dataSync.syndbRowCount')" prop="rows" width="72"/>
      <el-table-column :label="t('pages.dataSync.syndbTrigger')" width="96">
        <template #default="{ row }">
          <el-tag :type="syndbFlushReasonTagType(row.reason) || undefined" size="small">
            {{ formatReason(row.reason) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="right" :label="t('pages.dataSync.syndbWait')" width="88">
        <template #default="{ row }">{{ row.waitMs }} ms</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import type {SyndbFlushLog} from '@/types/api'
import {syndbFlushReasonTagType} from '@/utils/logParsers'

defineProps<{
  flush: SyndbFlushLog
  compact?: boolean
}>()

const {t} = useI18n()

const formatReason = (reason: string) => {
  const keyMap: Record<string, string> = {
    cpu_idle: 'pages.dataSync.syndbReasonCpuIdle',
    force_wait: 'pages.dataSync.syndbReasonForceWait',
    shutdown: 'pages.dataSync.syndbReasonShutdown',
  }
  const key = keyMap[reason]
  return key ? t(key) : reason
}
</script>

<style scoped>
.syndb-flush-log {
  width: 100%;
}

.syndb-flush-log.compact .syndb-flush-summary {
  gap: 6px;
}

.syndb-flush-summary {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  color: #606266;
  font-size: 13px;
}

.syndb-flush-metric {
  white-space: nowrap;
}

.syndb-flush-table {
  margin-top: 10px;
  width: 100%;
}
</style>

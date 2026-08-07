<template>
  <div class="syndb-flush-log" :class="{ compact }">
    <div class="syndb-flush-summary">
      <el-tag v-if="flush.reason" :type="syndbFlushReasonTagType(flush.reason) || undefined" size="small">
        {{ formatSyndbFlushReason(flush.reason) }}
      </el-tag>
      <span class="syndb-flush-metric">{{ flush.rows }} 行</span>
      <span class="syndb-flush-metric">{{ flush.queues }} 队列</span>
      <span class="syndb-flush-metric">{{ flush.costMs }} ms</span>
      <template v-if="!compact">
        <span v-if="flush.batchLimit" class="syndb-flush-metric">批次上限 {{ flush.batchLimit }}</span>
        <span v-if="flush.idleQueues" class="syndb-flush-metric">空闲队列 {{ flush.idleQueues }}</span>
        <span v-if="flush.forceQueues" class="syndb-flush-metric">强制队列 {{ flush.forceQueues }}</span>
      </template>
      <span v-if="flush.sysCpu !== undefined" class="syndb-flush-metric">系统 CPU {{ flush.sysCpu.toFixed(1) }}%</span>
      <span v-if="flush.cpuIdle !== undefined" class="syndb-flush-metric">CPU 空闲 {{ flush.cpuIdle.toFixed(1) }}%</span>
      <span v-if="flush.idleThreshold !== undefined" class="syndb-flush-metric">阈值 {{ flush.idleThreshold.toFixed(0) }}%</span>
    </div>

    <el-table
        v-if="!compact && flush.details.length"
        :data="flush.details"
        border
        class="syndb-flush-table"
        size="small"
        stripe
    >
      <el-table-column label="表" min-width="140" prop="table" show-overflow-tooltip/>
      <el-table-column label="列" min-width="120" prop="col" show-overflow-tooltip/>
      <el-table-column align="right" label="行数" prop="rows" width="72"/>
      <el-table-column label="触发" width="96">
        <template #default="{ row }">
          <el-tag :type="syndbFlushReasonTagType(row.reason) || undefined" size="small">
            {{ formatSyndbFlushReason(row.reason) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="right" label="等待" width="88">
        <template #default="{ row }">{{ row.waitMs }} ms</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts" setup>
import type {SyndbFlushLog} from '@/types/api'
import {formatSyndbFlushReason, syndbFlushReasonTagType} from '@/utils/logParsers'

defineProps<{
  flush: SyndbFlushLog
  compact?: boolean
}>()
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

<template>
  <div class="dashboard-trend-range-toolbar">
    <span class="dashboard-trend-range-label">{{ t('pages.dashboard.gameBetTrendRange') }}</span>
    <el-date-picker
        :model-value="modelValue"
        type="daterange"
        value-format="YYYY-MM-DD"
        :start-placeholder="t('common.startTime')"
        :end-placeholder="t('common.endTime')"
        :disabled-date="disableFutureDate"
        class="dashboard-trend-date-range"
        :clearable="clearable"
        @update:model-value="onDateRangeChange"
    />
    <el-button-group class="dashboard-trend-presets">
      <el-button
          v-for="days in presetDays"
          :key="days"
          size="small"
          :type="activePresetDays === days ? 'primary' : 'default'"
          @click="emit('preset', days)"
      >
        {{ presetLabel(days) }}
      </el-button>
    </el-button-group>
    <el-button type="primary" :loading="loading" @click="emit('search')">{{ t('common.search') }}</el-button>
    <el-button @click="emit('reset')">{{ t('pages.dashboard.gameBetTrendResetRange') }}</el-button>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {
  clampDashboardTrendDateRange,
  DASHBOARD_TREND_MAX_CALENDAR_DAYS,
  DASHBOARD_TREND_PRESET_DAYS,
} from '../dashboard-trend-range'

const presetDays = DASHBOARD_TREND_PRESET_DAYS

const props = withDefaults(
    defineProps<{
      modelValue: [string, string] | null
      loading?: boolean
      activePresetDays?: number | null
      maxCalendarDays?: number
      clearable?: boolean
    }>(),
    {
      activePresetDays: null,
      maxCalendarDays: DASHBOARD_TREND_MAX_CALENDAR_DAYS,
      clearable: false,
    },
)

const emit = defineEmits<{
  'update:modelValue': [value: [string, string] | null]
  search: []
  reset: []
  preset: [days: number]
  clamped: []
}>()

const {t} = useI18n()

const presetLabel = (days: number) => {
  if (days === 1) {
    return t('pages.dashboard.dashboardTrendPreset1d')
  }
  if (days === 7) {
    return t('pages.dashboard.dashboardTrendPreset7d')
  }
  return t('pages.dashboard.dashboardTrendPreset30d')
}

const onDateRangeChange = (value: [string, string] | null) => {
  if (!value || props.maxCalendarDays <= 0) {
    emit('update:modelValue', value)
    return
  }
  const {range, clamped} = clampDashboardTrendDateRange(value, props.maxCalendarDays)
  emit('update:modelValue', range)
  if (clamped) {
    emit('clamped')
  }
}

const disableFutureDate = (date: Date) => date.getTime() > Date.now()
</script>

<style scoped>
.dashboard-trend-range-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.dashboard-trend-range-label {
  font-size: 13px;
  color: #64748b;
}

.dashboard-trend-date-range {
  width: 260px;
  max-width: 100%;
}
</style>

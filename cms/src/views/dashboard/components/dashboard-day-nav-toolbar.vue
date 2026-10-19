<template>
  <div class="dashboard-day-nav-toolbar">
    <el-button size="small" :disabled="loading" @click="emit('prev')">
      {{ t('pages.dashboard.dashboardDayPrev') }}
    </el-button>
    <el-date-picker
        :model-value="focusDay"
        type="date"
        value-format="YYYY-MM-DD"
        :placeholder="t('pages.dashboard.dashboardDayPick')"
        :disabled-date="disableFutureDate"
        :disabled="loading"
        :clearable="false"
        class="dashboard-day-nav-picker"
        @update:model-value="onPickDay"
    />
    <el-button size="small" :disabled="loading || !canNext" @click="emit('next')">
      {{ t('pages.dashboard.dashboardDayNext') }}
    </el-button>
    <el-button
        v-if="!isToday"
        size="small"
        link
        type="primary"
        :disabled="loading"
        @click="emit('today')"
    >
      {{ t('pages.dashboard.dashboardDayToday') }}
    </el-button>
  </div>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import {formatServerDateStringFromDate} from '@/utils/server-datetime'
import {
  canGoDashboardNextDay,
  clampDashboardFocusDay,
  getDashboardTodayDateString,
  isDashboardFocusDayToday,
} from '../dashboard-day-nav'

const props = defineProps<{
  focusDay: string
  loading?: boolean
}>()

const emit = defineEmits<{
  prev: []
  next: []
  today: []
  pick: [day: string]
}>()

const {t} = useI18n()

const isToday = computed(() => isDashboardFocusDayToday(props.focusDay))
const canNext = computed(() => canGoDashboardNextDay(props.focusDay))

const disableFutureDate = (date: Date) =>
    formatServerDateStringFromDate(date) > getDashboardTodayDateString()

const onPickDay = (value: string | null) => {
  if (!value) {
    return
  }
  const day = clampDashboardFocusDay(value)
  if (day === props.focusDay) {
    return
  }
  emit('pick', day)
}
</script>

<style scoped>
.dashboard-day-nav-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.dashboard-day-nav-picker {
  width: 148px;
}
</style>

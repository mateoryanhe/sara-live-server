<template>
  <div>
    <el-alert
        :closable="false"
        class="hint-alert"
        show-icon
        :title="t('pages.anchorList.dailyLiveDurationHint')"
        type="info"
    />

    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
      <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
        <template #default="{ row }">{{ formatLiveDurationMinutes(row.liveDuration, t) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.dailySettled')" min-width="100">
        <template #default="{ row }">
          <el-tag :type="row.settled ? 'success' : 'warning'">
            {{ row.settled ? t('pages.anchorList.dailySettledYes') : t('pages.anchorList.dailySettledNo') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.roomUpdatedAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.anchorList.noDailyEffectiveLiveData')"/>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {accountApi} from '@/api'
import type {AnchorDailyEffectiveLiveItem} from '@/types/api'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const props = defineProps<{
  anchorId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<AnchorDailyEffectiveLiveItem[]>([])
const loaded = ref(false)

const fetchList = async () => {
  if (!props.anchorId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await accountApi.getAnchorDailyEffectiveLiveList({
      anchorId: props.anchorId,
      pageIndex: 1,
      pageSize: 8,
      settled: 0,
    })
    tableData.value = response.data || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load daily effective live list:', error)
    ElMessage.error(t('pages.anchorList.dailyEffectiveLiveFetchFailed'))
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
}

watch(
  () => props.anchorId,
  () => {
    resetState()
    if (props.active) {
      fetchList()
    }
  },
)

watch(
  () => props.active,
  (active) => {
    if (active && !loaded.value && props.anchorId) {
      fetchList()
    }
  },
  {immediate: true},
)
</script>

<style scoped>
.hint-alert {
  margin-bottom: 16px;
}
</style>

<template>
  <div>
    <el-table v-loading="loading" :data="tableData" style="width:100%">
      <el-table-column :label="t('pages.anchorList.archiveId')" min-width="180" prop="id"/>
      <el-table-column :label="t('pages.anchorList.liveIncome')" min-width="110">
        <template #default="{ row }">{{ formatNum(row.totalIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.giftIncome')" min-width="110">
        <template #default="{ row }">{{ formatNum(row.totalGiftIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalPaidDanmakuIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.privateRoomTicketIncome')" min-width="130">
        <template #default="{ row }">{{ formatNum(row.totalPrivateRoomTicketIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.privateRoomWatchIncome')" min-width="130">
        <template #default="{ row }">{{ formatNum(row.totalPrivateRoomWatchIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.videoCallIncome')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalVideoCallIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('pages.anchorList.totalLiveDuration')" min-width="120">
        <template #default="{ row }">{{ formatNum(row.totalLiveDuration) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" min-width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.anchorList.noArchiveData')"/>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {GuildIncomeArchiveItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const props = defineProps<{
  guildId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<GuildIncomeArchiveItem[]>([])
const loaded = ref(false)

const fetchList = async () => {
  if (!props.guildId) {
    tableData.value = []
    return
  }
  loading.value = true
  try {
    const response = await guildApi.getGuildIncomeArchives(props.guildId)
    tableData.value = response.list || []
    loaded.value = true
  } catch (error) {
    console.error('Failed to load guild income archives:', error)
    ElMessage.error(t('pages.guildList.archiveFetchFailed'))
  } finally {
    loading.value = false
  }
}

const formatNum = (value?: number | null) => formatAmount(value, '0')

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
  () => props.guildId,
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
    if (active && !loaded.value && props.guildId) {
      fetchList()
    }
  },
  {immediate: true},
)
</script>

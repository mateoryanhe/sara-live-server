<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveWeeklyUnsettledLiveList') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.liveWeeklyUnsettledLiveList.weekRangeHint', {start: weekRange.start, end: weekRange.end})"
          class="week-hint"
          type="info"
      />

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.keyword"
              clearable
              :placeholder="t('pages.liveWeeklyUnsettledLiveList.keywordPlaceholder')"
              style="width: 220px"
          />
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.platformAnchor')">
          <div class="anchor-filter anchor-filter--compact">
            <el-input
                :model-value="platformAnchorInputValue"
                class="anchor-input"
                disabled
                :placeholder="t('pages.liveRecordList.noPlatformAnchorSelected')"
            />
            <el-button size="small" @click="openPlatformAnchorPicker">
              {{ t('pages.liveRecordList.selectPlatformAnchor') }}
            </el-button>
            <el-button v-if="selectedPlatformAnchors.length > 0" link size="small" @click="clearPlatformAnchors">
              {{ t('pages.liveRecordList.clearPlatformAnchor') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.guildAnchor')">
          <div class="anchor-filter anchor-filter--compact">
            <el-input
                :model-value="guildAnchorInputValue"
                class="anchor-input"
                disabled
                :placeholder="t('pages.liveRecordList.noGuildAnchorSelected')"
            />
            <el-button size="small" @click="openGuildAnchorPicker">
              {{ t('pages.liveRecordList.selectGuildAnchor') }}
            </el-button>
            <el-button v-if="selectedGuildAnchors.length > 0" link size="small" @click="clearGuildAnchors">
              {{ t('pages.liveRecordList.clearGuildAnchor') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('export')" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
        <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId">
          <template #default="{ row }">
            <el-button v-if="row.roomId" link type="primary" @click="openAnchorDetail(row.roomId)">
              {{ row.roomId }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.roomAvatar"
                :preview-src-list="[row.roomAvatar]"
                :src="row.roomAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.roomId && row.roomNickname"
                link
                type="primary"
                @click="openUserDetail(row.roomId)"
            >
              {{ row.roomNickname }}
            </el-button>
            <span v-else>{{ row.roomNickname || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.liveWeeklyUnsettledLiveList.weeklyUnsettledTotalIncome')" align="right" min-width="160">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.weeklyUnsettledTotalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
          <template #default="{ row }">{{ formatLiveDurationMinutes(row.liveDuration, t) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveDailyEffectiveLiveList.dailyLiveIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
            v-model:current-page="pagination.pageIndex"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handlePageChange"
            @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <AnchorPickerDialog
        v-model:visible="platformAnchorPickerVisible"
        :initial-anchors="selectedPlatformAnchors"
        multiple
        platform-only
        @confirm-multiple="handlePlatformAnchorsPicked"
    />

    <GuildAnchorPickerDialog
        v-model:visible="guildAnchorPickerVisible"
        :initial-anchors="selectedGuildAnchors"
        @confirm="handleGuildAnchorsPicked"
    />
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {liveRecordApi} from '@/api'
import type {AnchorListItem, WeeklyUnsettledLiveItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import GuildAnchorPickerDialog from '@/components/GuildAnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_WEEKLY_UNSETTLED_LIVE} from '@/utils/cms-async-export'
import {buildLiveWeeklyUnsettledLiveListCsvColumns} from '@/utils/weekly-unsettled-live-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {
  getServerWeekDateRange,
  toServerDayStartUnix,
  toServerDayEndUnix,
} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('LiveWeeklyUnsettledLiveList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('LiveWeeklyUnsettledLiveList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()

const loading = ref(false)
const tableData = ref<WeeklyUnsettledLiveItem[]>([])
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const selectedGuildAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)
const guildAnchorPickerVisible = ref(false)

const searchForm = reactive({
  keyword: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const weekRange = computed(() => getServerWeekDateRange())

const formatPlatformAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  return `${nickname} (${anchor.id})`
}

const formatGuildAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  if (anchor.guildName) {
    return `[${anchor.guildName}] ${nickname} (${anchor.id})`
  }
  return `${nickname} (${anchor.id})`
}

const platformAnchorInputValue = computed(() => {
  if (selectedPlatformAnchors.value.length === 0) {
    return ''
  }
  if (selectedPlatformAnchors.value.length === 1) {
    return formatPlatformAnchorLabel(selectedPlatformAnchors.value[0])
  }
  return t('pages.liveRecordList.selectedPlatformAnchorsCount', {count: selectedPlatformAnchors.value.length})
})

const guildAnchorInputValue = computed(() => {
  if (selectedGuildAnchors.value.length === 0) {
    return ''
  }
  if (selectedGuildAnchors.value.length === 1) {
    return formatGuildAnchorLabel(selectedGuildAnchors.value[0])
  }
  return t('pages.liveRecordList.selectedGuildAnchorsCount', {count: selectedGuildAnchors.value.length})
})

const openPlatformAnchorPicker = () => {
  platformAnchorPickerVisible.value = true
}

const clearPlatformAnchors = () => {
  selectedPlatformAnchors.value = []
}

const handlePlatformAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedPlatformAnchors.value = anchors
}

const openGuildAnchorPicker = () => {
  guildAnchorPickerVisible.value = true
}

const clearGuildAnchors = () => {
  selectedGuildAnchors.value = []
}

const handleGuildAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedGuildAnchors.value = anchors
}

const buildSelectedAnchorIds = () => {
  const platformIds = selectedPlatformAnchors.value.map(anchor => String(anchor.id))
  const guildIds = selectedGuildAnchors.value.map(anchor => String(anchor.id))
  return [...new Set([...platformIds, ...guildIds])]
}

const buildFilterParams = () => ({
  anchorIds: buildSelectedAnchorIds(),
  keyword: searchForm.keyword.trim() || undefined,
})

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRecordApi.getWeeklyUnsettledLiveList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load weekly unsettled live list:', error)
    ElMessage.error(t('pages.liveWeeklyUnsettledLiveList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  selectedPlatformAnchors.value = []
  selectedGuildAnchors.value = []
  searchForm.keyword = ''
  pagination.pageIndex = 1
  fetchList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const openAnchorDetail = (anchorId: string | number) => {
  if (!anchorId) {
    return
  }
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(anchorId)},
  })
}

const handleExport = async () => {
  await runExport(
    CMS_EXPORT_TYPE_LIVE_WEEKLY_UNSETTLED_LIVE,
    {
      headers: buildCsvHeaders(buildLiveWeeklyUnsettledLiveListCsvColumns(t)),
      ...buildFilterParams(),
    },
    `weekly-unsettled-live-${Date.now()}.csv`,
  )
}

onMounted(() => {
  fetchList()
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

.week-hint {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 12px;
}

.anchor-filter--compact {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.anchor-input {
  width: 220px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

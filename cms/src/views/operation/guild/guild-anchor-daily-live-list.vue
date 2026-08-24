<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildAnchorDailyLiveManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="hint-alert"
          show-icon
          :title="t('pages.guildAnchorDailyLiveList.hint')"
          type="info"
      />

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.guildAnchorDailyLiveList.anchor')">
          <div class="anchor-filter">
            <div class="anchor-selection">
              <el-input
                  :model-value="selectedAnchorsLabel"
                  class="anchor-input"
                  disabled
                  :placeholder="t('pages.guildAnchorDailyLiveList.noAnchorSelected')"
              />
              <div v-if="selectedAnchors.length > 0" class="anchor-tags">
                <el-tag
                    v-for="anchor in selectedAnchors"
                    :key="anchor.id"
                    closable
                    @close="removeSelectedAnchor(anchor)"
                >
                  {{ formatAnchorLabel(anchor) }}
                </el-tag>
              </div>
            </div>
            <el-button @click="openAnchorPicker">
              {{ t('pages.guildAnchorDailyLiveList.selectAnchor') }}
            </el-button>
            <el-button v-if="selectedAnchors.length > 0" @click="clearSelectedAnchors">
              {{ t('pages.guildAnchorDailyLiveList.clearAnchor') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.dailyLiveDate')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.anchorList.liveDateEnd')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.anchorList.liveDateRangeSeparator')"
              :start-placeholder="t('pages.anchorList.liveDateStart')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
          <el-button v-if="can('export')" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
        <el-table-column :label="t('pages.anchorList.dailyRecordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.guildAnchorIncomeSettlementLogList.roomNickname')" min-width="120">
          <template #default="{ row }">{{ row.roomNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.dailyLiveDate')" min-width="120" prop="liveDate"/>
        <el-table-column :label="t('pages.anchorList.dailyLiveDuration')" min-width="150">
          <template #default="{ row }">{{ formatLiveDurationMinutes(row.liveDuration, t) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.dailyReportedLiveDuration')" min-width="150">
          <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.liveIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.giftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.privateRoomTicketIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.privateRoomWatchIncome')" align="right" min-width="140">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPrivateRoomWatchIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.anchorList.videoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
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

      <div class="pagination-container">
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

      <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildAnchorDailyLiveList.noData')"/>
    </el-card>

    <AnchorPickerDialog
        v-model:visible="anchorPickerVisible"
        :initial-anchors="selectedAnchors"
        multiple
        owned-guild-only
        @confirm-multiple="handleAnchorsPicked"
    />
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {AnchorListItem, GuildAnchorDailyEffectiveLiveItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, buildDailyEffectiveLiveExportLabels, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_MY_GUILD_ANCHOR_DAILY_EFFECTIVE_LIVE} from '@/utils/cms-async-export'
import {buildGuildAnchorDailyEffectiveLiveCsvColumns} from '@/utils/daily-effective-live-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'

const {t} = useI18n()
const {can} = usePagePermission('GuildAnchorDailyLiveManagement')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()

const loading = ref(false)
const tableData = ref<GuildAnchorDailyEffectiveLiveItem[]>([])
const selectedAnchors = ref<AnchorListItem[]>([])
const anchorPickerVisible = ref(false)

const searchForm = reactive({
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const formatAnchorLabel = (anchor: AnchorListItem) => {
  const nickname = anchor.nickname || '-'
  return `${nickname} (${anchor.id})`
}

const selectedAnchorsLabel = computed(() => {
  if (selectedAnchors.value.length === 0) {
    return ''
  }
  if (selectedAnchors.value.length === 1) {
    return formatAnchorLabel(selectedAnchors.value[0])
  }
  return t('pages.guildAnchorDailyLiveList.selectedAnchorsCount', {count: selectedAnchors.value.length})
})

const openAnchorPicker = () => {
  anchorPickerVisible.value = true
}

const clearSelectedAnchors = () => {
  selectedAnchors.value = []
}

const removeSelectedAnchor = (anchor: AnchorListItem) => {
  selectedAnchors.value = selectedAnchors.value.filter(item => item.id !== anchor.id)
}

const handleAnchorsPicked = (anchors: AnchorListItem[]) => {
  selectedAnchors.value = anchors
}

const buildFilterParams = () => {
  const [liveDateStart, liveDateEnd] = searchForm.dateRange || []
  const roomIds = selectedAnchors.value.map(anchor => String(anchor.id))
  return {
    roomIds: roomIds.length > 0 ? roomIds : undefined,
    liveDateStart: liveDateStart || undefined,
    liveDateEnd: liveDateEnd || undefined,
    settled: 0,
  }
}

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getMyOwnedGuildAnchorDailyEffectiveLiveList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch guild anchor daily live list failed:', error)
    tableData.value = []
    pagination.total = 0
    ElMessage.error(t('pages.guildAnchorDailyLiveList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.dateRange = []
  selectedAnchors.value = []
  pagination.pageIndex = 1
  pagination.total = 0
  tableData.value = []
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

const handleExport = async () => {
  await runExport(
    CMS_EXPORT_TYPE_MY_GUILD_ANCHOR_DAILY_EFFECTIVE_LIVE,
    {
      headers: buildCsvHeaders(buildGuildAnchorDailyEffectiveLiveCsvColumns(t)),
      ...buildFilterParams(),
      ...buildDailyEffectiveLiveExportLabels(t),
    },
    `my-guild-anchor-daily-flow-${Date.now()}.csv`,
  )
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.hint-alert {
  margin-bottom: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.anchor-filter {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.anchor-selection {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.anchor-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-width: 420px;
}

.anchor-input {
  width: 260px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.money-amount {
  font-variant-numeric: tabular-nums;
}
</style>

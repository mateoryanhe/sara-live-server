<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRecordList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.keyword"
              clearable
              :placeholder="t('pages.liveRecordList.keywordPlaceholder')"
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
        <el-form-item :label="t('pages.liveRecordList.startDate')">
          <el-date-picker
              v-model="searchForm.startDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.startDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.liveRecordList.endDate')">
          <el-date-picker
              v-model="searchForm.endDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.liveRecordList.endDate')"
              style="width: 160px"
              type="date"
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
        <el-table-column :label="t('pages.liveRecordList.recordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.liveRecordList.anchorId')" min-width="180" prop="anchorId">
          <template #default="{ row }">
            <el-button v-if="row.anchorId" link type="primary" @click="openAnchorDetail(row.anchorId)">
              {{ row.anchorId }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.avatar"
                :preview-src-list="[row.avatar]"
                :src="row.avatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.anchorNickname')" min-width="120" prop="nickname">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.anchorId && row.nickname"
                link
                type="primary"
                @click="openUserDetail(row.anchorId)"
            >
              {{ row.nickname }}
            </el-button>
            <span v-else>{{ row.nickname || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.startTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.startTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.endTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.endTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.totalAudience')" prop="totalAudience" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.liveDuration')" width="120">
          <template #default="{ row }">{{ formatDuration(row.totalLiveDuration) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.totalIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.giftIncome')" align="right" min-width="120">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGiftIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.paidDanmakuIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalPaidDanmakuIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoTicketIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallTicketIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoBillingIncome')" align="right" min-width="150">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallBillingIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.videoCallIncome')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalVideoCallIncome) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRecordList.giftSenderCount')" prop="totalGiftSender" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.newFollowers')" prop="totalNewFollower" width="100"/>
        <el-table-column :label="t('pages.liveRecordList.totalGameBet')" align="right" min-width="130">
          <template #default="{ row }"><span class="money-amount">{{ formatWalletBalance(row.totalGameBet) }}</span></template>
        </el-table-column>
        <el-table-column :label="t('common.createdAt')" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
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
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {liveRecordApi} from '@/api'
import type {AnchorListItem, LiveRecordItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import GuildAnchorPickerDialog from '@/components/GuildAnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_RECORD} from '@/utils/cms-async-export'
import {buildLiveRecordCsvColumns} from '@/utils/live-record-csv'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate, toServerDayStartUnix, toServerDayEndUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('LiveRecordList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('LiveRecordList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<LiveRecordItem[]>([])
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const selectedGuildAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)
const guildAnchorPickerVisible = ref(false)

const searchForm = reactive({
  keyword: '',
  startDate: '',
  endDate: '',
})

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

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => ({
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
  anchorIds: buildSelectedAnchorIds(),
  keyword: searchForm.keyword.trim() || undefined,
  startTime: searchForm.startDate ? toServerDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toServerDayEndUnix(searchForm.endDate) : 0,
})

const buildFilterParams = () => ({
  anchorIds: buildSelectedAnchorIds(),
  keyword: searchForm.keyword.trim() || undefined,
  startTime: searchForm.startDate ? toServerDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toServerDayEndUnix(searchForm.endDate) : 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRecordApi.getLiveRecordList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load live records:', error)
    ElMessage.error(t('pages.liveRecordList.fetchFailed'))
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
  searchForm.startDate = ''
  searchForm.endDate = ''
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
    CMS_EXPORT_TYPE_LIVE_RECORD,
    {
      headers: buildCsvHeaders(buildLiveRecordCsvColumns(t)),
      ...buildFilterParams(),
    },
    `live-record-${Date.now()}.csv`,
  )
}

const formatDuration = (seconds: number | null | undefined) => {
  if (seconds === null || seconds === undefined || seconds <= 0) return '-'
  const total = Math.floor(seconds)
  const h = Math.floor(total / 3600)
  const m = Math.floor((total % 3600) / 60)
  const s = total % 60
  if (h > 0) {
    return t('pages.liveRecordList.durationHours', {h, m, s})
  }
  if (m > 0) {
    return t('pages.liveRecordList.durationMinutes', {m, s})
  }
  return t('pages.liveRecordList.durationSeconds', {s})
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

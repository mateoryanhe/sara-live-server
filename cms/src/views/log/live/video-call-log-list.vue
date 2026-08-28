<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.VideoCallLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.videoCallLogList.callerId')">
          <el-input v-model="searchForm.callerId" clearable :placeholder="t('pages.videoCallLogList.enterCallerId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.videoCallLogList.receiverId')">
          <el-input v-model="searchForm.receiverId" clearable :placeholder="t('pages.videoCallLogList.enterReceiverId')"/>
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
        <el-form-item :label="t('pages.videoCallLogList.source')">
          <el-select v-model="searchForm.source" clearable :placeholder="t('common.all')" style="width: 140px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.videoCallLogList.sourceLiveRoom')"/>
            <el-option :value="2" :label="t('pages.videoCallLogList.sourcePrivateMessage')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model="searchForm.status" clearable :placeholder="t('common.all')" style="width: 160px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.videoCallLogList.statusCalling')"/>
            <el-option :value="2" :label="t('pages.videoCallLogList.statusAnswered')"/>
            <el-option :value="3" :label="t('pages.videoCallLogList.statusInCall')"/>
            <el-option :value="4" :label="t('pages.videoCallLogList.statusEnded')"/>
            <el-option :value="5" :label="t('pages.videoCallLogList.statusRejected')"/>
            <el-option :value="6" :label="t('pages.videoCallLogList.statusCallTimeout')"/>
            <el-option :value="7" :label="t('pages.videoCallLogList.statusHeartTimeout')"/>
            <el-option :value="8" :label="t('pages.videoCallLogList.statusInsufficientDiamond')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.videoCallLogList.callTime')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('pages.videoCallLogList.endDate')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.videoCallLogList.dateRangeSeparator')"
              :start-placeholder="t('pages.videoCallLogList.startDate')"
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
        <el-table-column :label="t('common.createdAt')" fixed="left" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callDuration')" width="120">
          <template #default="{ row }">{{ formatDuration(row.callDuration) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.totalCostDiamond')" width="120">
          <template #default="{ row }">{{ formatAmount(row.totalCost) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.billingDurationMinutes')" prop="billingDuration" width="120"/>
        <el-table-column :label="t('common.status')" prop="statusText" width="140"/>
        <el-table-column :label="t('common.avatar')" align="center" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.callerAvatar"
                :preview-src-list="[row.callerAvatar]"
                :src="row.callerAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callerId')" min-width="180" prop="callerId">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail && row.callerId" link type="primary" @click="openUserDetail(row.callerId)">
              {{ row.callerId }}
            </el-button>
            <span v-else>{{ row.callerId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callerNickname')" min-width="120">
          <template #default="{ row }">{{ row.callerNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.receiverId')" min-width="180" prop="receiverId">
          <template #default="{ row }">
            <el-button v-if="canOpenReceiverDetail(row)" link type="primary" @click="openReceiverDetail(row)">
              {{ row.receiverId }}
            </el-button>
            <span v-else>{{ row.receiverId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.receiverNickname')" min-width="120">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.receiverId && row.receiverNickname"
                link
                type="primary"
                @click="openUserDetail(row.receiverId)"
            >
              {{ row.receiverNickname }}
            </el-button>
            <span v-else>{{ row.receiverNickname || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.callTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.callStartTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.answerTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.answerTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.endTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.orderEndTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.ticketDiamond')" width="110">
          <template #default="{ row }">{{ formatAmount(row.ticketPrice) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.pricePerMinuteDiamond')" width="130">
          <template #default="{ row }">{{ formatAmount(row.pricePerMinute) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.lastChargeTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.chargeTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.source')" prop="sourceText" width="90"/>
        <el-table-column :label="t('pages.videoCallLogList.callerLastHeart')" width="170">
          <template #default="{ row }">{{ formatDate(row.callerHeartTime) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.videoCallLogList.receiverLastHeart')" width="170">
          <template #default="{ row }">{{ formatDate(row.receiverHeartTime) }}</template>
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
import {videoCallLogApi} from '@/api'
import type {AnchorListItem, VideoCallLogItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import GuildAnchorPickerDialog from '@/components/GuildAnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_VIDEO_CALL_LOG} from '@/utils/cms-async-export'
import {formatAmount} from '@/utils/number-format'
import {buildVideoCallLogCsvColumns} from '@/utils/video-call-log-csv'
import {formatServerDateTime as formatDate, toServerDayStartUnix, toServerDayEndUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('VideoCallLogList')
const {canViewUserDetail, openUserDetail} = useUserDetailNav('VideoCallLogList')
const canViewAnchorDetail = computed(() => can('viewAnchorDetail'))
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<VideoCallLogItem[]>([])
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const selectedGuildAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)
const guildAnchorPickerVisible = ref(false)

const searchForm = reactive({
  callerId: '',
  receiverId: '',
  source: 0,
  status: 0,
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
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

const buildSelectedReceiverIds = () => {
  const platformIds = selectedPlatformAnchors.value.map(anchor => String(anchor.id))
  const guildIds = selectedGuildAnchors.value.map(anchor => String(anchor.id))
  return [...new Set([...platformIds, ...guildIds])]
}

const buildFilterParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    callerId: searchForm.callerId.trim(),
    receiverId: searchForm.receiverId.trim(),
    receiverIds: buildSelectedReceiverIds(),
    source: searchForm.source || 0,
    status: searchForm.status || 0,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
  }
}

const buildQueryParams = () => ({
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
  ...buildFilterParams(),
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await videoCallLogApi.getVideoCallLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load video call logs:', error)
    ElMessage.error(t('pages.videoCallLogList.fetchFailed'))
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
  searchForm.callerId = ''
  searchForm.receiverId = ''
  searchForm.source = 0
  searchForm.status = 0
  searchForm.dateRange = []
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

const handleExport = async () => {
  await runExport(
    CMS_EXPORT_TYPE_VIDEO_CALL_LOG,
    {
      headers: buildCsvHeaders(buildVideoCallLogCsvColumns(t)),
      ...buildFilterParams(),
    },
    `video-call-log-${Date.now()}.csv`,
  )
}

const formatDuration = (seconds: number | null | undefined) => {
  if (seconds === null || seconds === undefined || seconds <= 0) return '-'
  const total = Math.floor(seconds)
  const h = Math.floor(total / 3600)
  const m = Math.floor((total % 3600) / 60)
  const s = total % 60
  if (h > 0) {
    return t('pages.videoCallLogList.durationHours', {h, m, s})
  }
  if (m > 0) {
    return t('pages.videoCallLogList.durationMinutes', {m, s})
  }
  return t('pages.videoCallLogList.durationSeconds', {s})
}

const canOpenReceiverDetail = (row: VideoCallLogItem) => {
  if (!row.receiverId) {
    return false
  }
  if (row.receiverIsAnchor) {
    return canViewAnchorDetail.value
  }
  return canViewUserDetail.value
}

const openReceiverDetail = (row: VideoCallLogItem) => {
  if (!row.receiverId) {
    return
  }
  if (row.receiverIsAnchor) {
    router.push({
      path: '/user/anchor/anchor-detail',
      query: {id: String(row.receiverId)},
    })
    return
  }
  openUserDetail(row.receiverId)
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

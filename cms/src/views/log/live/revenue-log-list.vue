<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRevenueLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="88px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.keyword"
              clearable
              :placeholder="t('pages.revenueLogList.keywordPlaceholder')"
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
        <el-form-item :label="t('pages.revenueLogList.revenueType')">
          <el-select v-model="searchForm.revenueType" clearable :placeholder="t('common.all')" style="width: 200px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option
                v-for="option in revenueTypeOptions"
                :key="option.value"
                :label="t(option.labelKey)"
                :value="option.value"
            />
          </el-select>
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
        <el-table-column :label="t('common.createdAt')" min-width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.liveRoomId')" min-width="180" prop="roomId">
          <template #default="{ row }">
            <el-button v-if="row.roomId" link type="primary" @click="openAnchorDetail(row.roomId)">
              {{ row.roomId }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.receiverUserId')" min-width="180" prop="receiverId">
          <template #default="{ row }">
            <el-button v-if="row.receiverId" link type="primary" @click="openAnchorDetail(row.receiverId)">
              {{ row.receiverId }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.avatar')" width="80">
          <template #default="{ row }">
            <el-image
                v-if="row.receiverAvatar"
                :preview-src-list="[row.receiverAvatar]"
                :src="row.receiverAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.receiverNickname')" min-width="120" prop="receiverNickname">
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
        <el-table-column :label="t('pages.revenueLogList.payerUserId')" min-width="180" prop="senderId"/>
        <el-table-column
            :label="t('pages.revenueLogList.payerUserAvatar')"
            class-name="col-nowrap"
            label-class-name="col-nowrap"
            width="100"
        >
          <template #default="{ row }">
            <el-image
                v-if="row.senderAvatar"
                :preview-src-list="[row.senderAvatar]"
                :src="row.senderAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:40px;height:40px;border-radius:50%"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.payerNickname')" min-width="120" prop="senderNickname">
          <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.bizId')" min-width="180" prop="bizId"/>
        <el-table-column :label="t('pages.revenueLogList.bizName')" min-width="120" prop="bizName">
          <template #default="{ row }">{{ row.bizName || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.revenueType')" min-width="140" prop="revenueTypeText">
          <template #default="{ row }">{{ row.revenueTypeText || formatRevenueType(row.revenueType) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.count')" prop="count" width="80"/>
        <el-table-column :label="t('pages.revenueLogList.unitPriceDiamond')" prop="unitPrice" width="110"/>
        <el-table-column :label="t('pages.revenueLogList.totalAmountDiamond')" prop="totalAmount" width="120"/>
        <el-table-column :label="t('common.status')" prop="statusText" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'warning' : 'success'" size="small">
              {{ row.statusText || (row.status === 1 ? t('pages.revenueLogList.refunded') : t('common.normal')) }}
            </el-tag>
          </template>
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
import {liveRevenueLogApi} from '@/api'
import type {AnchorListItem, LiveRevenueLogItem} from '@/types/api'
import AnchorPickerDialog from '@/components/AnchorPickerDialog.vue'
import GuildAnchorPickerDialog from '@/components/GuildAnchorPickerDialog.vue'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_REVENUE_LOG} from '@/utils/cms-async-export'
import {buildLiveRevenueLogCsvColumns} from '@/utils/live-revenue-log-csv'
import {createLiveRevenueTypeFormatter, LIVE_REVENUE_TYPE_OPTIONS} from '@/utils/live-revenue-type'
import {formatServerDateTime as formatDate, toServerDayStartUnix, toServerDayEndUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('LiveRevenueLogList')
const {openUserDetail, canViewUserDetail} = useUserDetailNav('LiveRevenueLogList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const revenueTypeOptions = LIVE_REVENUE_TYPE_OPTIONS
const formatRevenueType = createLiveRevenueTypeFormatter(t)
const loading = ref(false)
const tableData = ref<LiveRevenueLogItem[]>([])
const selectedPlatformAnchors = ref<AnchorListItem[]>([])
const selectedGuildAnchors = ref<AnchorListItem[]>([])
const platformAnchorPickerVisible = ref(false)
const guildAnchorPickerVisible = ref(false)

const searchForm = reactive({
  keyword: '',
  revenueType: 0,
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

const buildSelectedReceiverIds = () => {
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
  receiverIds: buildSelectedReceiverIds(),
  keyword: searchForm.keyword.trim() || undefined,
  revenueType: searchForm.revenueType || 0,
  startTime: searchForm.startDate ? toServerDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toServerDayEndUnix(searchForm.endDate) : 0,
})

const buildFilterParams = () => ({
  receiverIds: buildSelectedReceiverIds(),
  keyword: searchForm.keyword.trim() || undefined,
  revenueType: searchForm.revenueType || 0,
  startTime: searchForm.startDate ? toServerDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toServerDayEndUnix(searchForm.endDate) : 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await liveRevenueLogApi.getLiveRevenueLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load revenue logs:', error)
    ElMessage.error(t('pages.revenueLogList.fetchFailed'))
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
  searchForm.revenueType = 0
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
    CMS_EXPORT_TYPE_LIVE_REVENUE_LOG,
    {
      headers: buildCsvHeaders(buildLiveRevenueLogCsvColumns(t, formatRevenueType)),
      ...buildFilterParams(),
    },
    `live-revenue-log-${Date.now()}.csv`,
  )
}


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

:deep(.col-nowrap .cell) {
  white-space: nowrap;
}
</style>

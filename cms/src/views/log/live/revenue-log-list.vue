<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRevenueLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('pages.revenueLogList.platformAnchor')">
          <AnchorRemoteSelect
              v-model="searchForm.platformAnchorIds"
              mode="platform"
              :placeholder="t('pages.revenueLogList.searchPlatformAnchor')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.revenueLogList.guildAnchor')">
          <AnchorRemoteSelect
              v-model="searchForm.guildAnchorIds"
              mode="guild"
              :placeholder="t('pages.revenueLogList.searchGuildAnchor')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.revenueLogList.revenueType')">
          <el-select v-model="searchForm.revenueType" clearable :placeholder="t('common.all')" style="width: 140px">
            <el-option :value="0" :label="t('common.all')"/>
            <el-option :value="1" :label="t('pages.revenueLogList.revenueGift')"/>
            <el-option :value="2" :label="t('pages.revenueLogList.revenuePaidDanmaku')"/>
            <el-option :value="3" :label="t('pages.revenueLogList.revenueGameBet')"/>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.revenueLogList.startDate')">
          <el-date-picker
              v-model="searchForm.startDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.revenueLogList.startDate')"
              style="width: 160px"
              type="date"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item :label="t('pages.revenueLogList.endDate')">
          <el-date-picker
              v-model="searchForm.endDate"
              clearable
              format="YYYY-MM-DD"
              :placeholder="t('pages.revenueLogList.endDate')"
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
        <el-table-column :label="t('pages.revenueLogList.logId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.revenueLogList.revenueType')" prop="revenueTypeText" width="100">
          <template #default="{ row }">{{ row.revenueTypeText || formatRevenueType(row.revenueType) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.liveRoomId')" min-width="180" prop="roomId"/>
        <el-table-column :label="t('pages.revenueLogList.liveRecordId')" min-width="180" prop="liveRecordId"/>
        <el-table-column :label="t('pages.revenueLogList.payerUserId')" min-width="180" prop="senderId"/>
        <el-table-column :label="t('pages.revenueLogList.payerNickname')" min-width="120" prop="senderNickname">
          <template #default="{ row }">{{ row.senderNickname || '-' }}</template>
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
          <template #default="{ row }">{{ row.receiverNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('menu.UserDetail')" width="110">
          <template #default="{ row }">
            <el-button v-if="row.receiverId" link type="primary" @click="openUserDetail(row.receiverId)">
              {{ t('pages.userList.viewDetail') }}
            </el-button>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.revenueLogList.bizId')" min-width="180" prop="bizId"/>
        <el-table-column :label="t('pages.revenueLogList.bizName')" min-width="120" prop="bizName">
          <template #default="{ row }">{{ row.bizName || '-' }}</template>
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
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {onMounted, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {liveRevenueLogApi} from '@/api'
import type {LiveRevenueLogItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_LIVE_REVENUE_LOG} from '@/utils/cms-async-export'
import {buildLiveRevenueLogCsvColumns} from '@/utils/live-revenue-log-csv'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('LiveRevenueLogList')
const {openUserDetail} = useUserDetailNav('LiveRevenueLogList')
const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<LiveRevenueLogItem[]>([])

const searchForm = reactive({
  platformAnchorIds: [] as string[],
  guildAnchorIds: [] as string[],
  revenueType: 0,
  startDate: '',
  endDate: '',
})

const buildSelectedReceiverIds = () => {
  return [...new Set([...searchForm.platformAnchorIds, ...searchForm.guildAnchorIds])]
}

const formatRevenueType = (type: number) => {
  const map: Record<number, string> = {
    1: t('pages.revenueLogList.revenueGift'),
    2: t('pages.revenueLogList.revenuePaidDanmaku'),
    3: t('pages.revenueLogList.revenueGameBet'),
  }
  return map[type] || t('pages.revenueLogList.unknown')
}

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const toDayStartUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T00:00:00`).getTime() / 1000)
}

const toDayEndUnix = (dateStr: string): number => {
  return Math.floor(new Date(`${dateStr}T23:59:59`).getTime() / 1000)
}

const buildQueryParams = () => ({
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
  receiverIds: buildSelectedReceiverIds(),
  revenueType: searchForm.revenueType || 0,
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
})

const buildFilterParams = () => ({
  receiverIds: buildSelectedReceiverIds(),
  revenueType: searchForm.revenueType || 0,
  startTime: searchForm.startDate ? toDayStartUnix(searchForm.startDate) : 0,
  endTime: searchForm.endDate ? toDayEndUnix(searchForm.endDate) : 0,
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
  searchForm.platformAnchorIds = []
  searchForm.guildAnchorIds = []
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

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
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

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

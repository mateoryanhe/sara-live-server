<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="80px">
      <el-form-item :label="t(`${ns}.gameCode`)">
        <el-input v-model="searchForm.gameCode" clearable :placeholder="t(`${ns}.gameCodePlaceholder`)"/>
      </el-form-item>
      <el-form-item :label="t(`${ns}.orderId`)">
        <el-input v-model="searchForm.orderId" clearable :placeholder="t(`${ns}.orderIdPlaceholder`)"/>
      </el-form-item>
      <el-form-item :label="t(`${ns}.platform`)">
        <el-input v-model="searchForm.platformType" clearable :placeholder="t(`${ns}.platformPlaceholder`)"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        <el-button v-if="canExport" :loading="exporting" @click="handleExport">{{ t('common.export') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading || exporting" :data="tableData" :element-loading-text="exportStatusTip || undefined" style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column :label="t(`${ns}.recordId`)" min-width="180" prop="id"/>
      <el-table-column :label="t(`${ns}.gameCode`)" min-width="120" prop="gameCode"/>
      <el-table-column :label="t(`${ns}.nameEn`)" min-width="140" prop="nameEn" show-overflow-tooltip>
        <template #default="{ row }">{{ row.nameEn || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t(`${ns}.cover`)" min-width="100">
        <template #default="{ row }">
          <el-image
              v-if="row.cover"
              :preview-src-list="[row.cover]"
              :src="row.cover"
              fit="cover"
              preview-teleported
              style="width: 48px; height: 48px; border-radius: 4px"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="amountColumnLabel" prop="amount" width="120">
        <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="t(`${ns}.platform`)" prop="platformType" width="100"/>
      <template v-if="logType === 'bet'">
        <el-table-column :label="t(`${ns}.liveRoomId`)" min-width="180" prop="liveRoomId">
          <template #default="{ row }">{{ row.liveRoomId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t(`${ns}.liveRoomTitle`)" min-width="160" prop="liveRoomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.liveRoomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t(`${ns}.anchorNickname`)" min-width="140" prop="anchorNickname">
          <template #default="{ row }">{{ row.anchorNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t(`${ns}.liveRecordId`)" min-width="180" prop="liveRecordId">
          <template #default="{ row }">{{ row.liveRecordId || '-' }}</template>
        </el-table-column>
      </template>
      <el-table-column :label="t(`${ns}.orderId`)" min-width="200" prop="orderId" show-overflow-tooltip/>
      <el-table-column :label="t(`${ns}.time`)" width="170">
        <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.userList.noData')"/>

    <div v-if="pagination.total > 0" class="pagination">
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
  </div>
</template>

<script lang="ts" setup>
import {computed, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {gameBetLogApi, gameWinLogApi} from '@/api'
import type {GameBetLogItem, GameWinLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {buildCsvHeaders, useCmsAsyncExport} from '@/composables/useCmsAsyncExport'
import {CMS_EXPORT_TYPE_GAME_BET_LOG, CMS_EXPORT_TYPE_GAME_WIN_LOG} from '@/utils/cms-async-export'
import {buildGameBetLogCsvColumns, buildGameWinLogCsvColumns} from '@/utils/game-log-csv'

const props = defineProps<{
  userId: string
  active: boolean
  logType: 'bet' | 'win'
  exportPermission: 'exportGameBetLog' | 'exportGameWinLog'
}>()

const {t} = useI18n()
const {can} = usePagePermission('UserDetail')
const canExport = computed(() => can(props.exportPermission))
const ns = computed(() => props.logType === 'bet' ? 'pages.gameBetLogList' : 'pages.gameWinLogList')
const amountColumnLabel = computed(() =>
    props.logType === 'bet'
        ? t('pages.gameBetLogList.betAmount')
        : t('pages.gameWinLogList.winAmount'))
const fetchFailedMessage = computed(() =>
    props.logType === 'bet'
        ? t('pages.userList.gameBetLogFetchFailed')
        : t('pages.userList.gameWinLogFetchFailed'))
const exportFilePrefix = computed(() =>
    props.logType === 'bet' ? 'user-game-bet-log' : 'user-game-win-log')

const {exporting, exportStatusTip, runExport} = useCmsAsyncExport()
const loading = ref(false)
const tableData = ref<Array<GameBetLogItem | GameWinLogItem>>([])
const loaded = ref(false)

const searchForm = reactive({
  gameCode: '',
  orderId: '',
  platformType: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildFilterParams = () => ({
  userId: props.userId,
  gameCode: searchForm.gameCode.trim(),
  orderId: searchForm.orderId.trim(),
  platformType: searchForm.platformType.trim(),
})

const buildQueryParams = () => ({
  ...buildFilterParams(),
  pageIndex: pagination.pageIndex,
  pageSize: pagination.pageSize,
})

const fetchList = async () => {
  if (!props.userId) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = props.logType === 'bet'
        ? await gameBetLogApi.getGameBetLogList(buildQueryParams())
        : await gameWinLogApi.getGameWinLogList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load user game logs:', error)
    ElMessage.error(fetchFailedMessage.value)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.gameCode = ''
  searchForm.orderId = ''
  searchForm.platformType = ''
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

const formatRowIndex = (index: number) =>
    (pagination.pageIndex - 1) * pagination.pageSize + index + 1

const handleExport = async () => {
  if (!props.userId) {
    ElMessage.warning(t('common.exportEmpty'))
    return
  }
  const columns = props.logType === 'bet'
      ? buildGameBetLogCsvColumns(t)
      : buildGameWinLogCsvColumns(t)
  await runExport(
      props.logType === 'bet' ? CMS_EXPORT_TYPE_GAME_BET_LOG : CMS_EXPORT_TYPE_GAME_WIN_LOG,
      {
        headers: buildCsvHeaders(columns),
        ...buildFilterParams(),
      },
      `${exportFilePrefix.value}-${props.userId}-${Date.now()}.csv`,
  )
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const resetState = () => {
  loaded.value = false
  tableData.value = []
  searchForm.gameCode = ''
  searchForm.orderId = ''
  searchForm.platformType = ''
  pagination.pageIndex = 1
  pagination.total = 0
}

watch(
    () => [props.userId, props.logType] as const,
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
      if (active && !loaded.value && props.userId) {
        fetchList()
      }
    },
    {immediate: true},
)
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

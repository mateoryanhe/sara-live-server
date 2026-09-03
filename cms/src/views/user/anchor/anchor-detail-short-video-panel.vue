<template>
  <div>
    <el-form :model="searchForm" class="search-form" inline label-width="80px">
      <el-form-item :label="t('common.title')">
        <el-input v-model="searchForm.title" clearable :placeholder="t('pages.shortVideoList.titlePlaceholder')"/>
      </el-form-item>
      <el-form-item :label="t('common.status')">
        <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
          <el-option :value="0" :label="t('common.all')"/>
          <el-option :value="2" :label="t('common.onlyOnShelf')"/>
          <el-option :value="1" :label="t('common.onlyOffShelf')"/>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
        <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column label="ID" min-width="160" prop="id"/>
      <el-table-column :label="t('common.title')" min-width="140" prop="title" show-overflow-tooltip/>
      <el-table-column :label="t('pages.shortVideoList.cover')" width="100">
        <template #default="{ row }">
          <el-image
              v-if="row.cover"
              :preview-src-list="[row.cover]"
              :src="row.cover"
              fit="cover"
              preview-teleported
              style="width: 72px; height: 40px"
          />
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.shortVideoList.video')" min-width="200">
        <template #default="{ row }">
          <div v-if="row.video" class="table-video-cell">
            <video
                :key="row.video"
                :src="row.video"
                class="table-video-preview"
                controls
                preload="metadata"
            />
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.shortVideoList.isPaid')" width="90">
        <template #default="{ row }">
          <el-tag :type="row.isPaid === 1 ? 'warning' : 'success'">
            {{ row.isPaid === 1 ? t('pages.shortVideoList.paid') : t('pages.shortVideoList.free') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.shortVideoList.payDiamond')" width="110">
        <template #default="{ row }">
          {{ row.isPaid === 1 ? formatAmount(row.payDiamond) : '-' }}
        </template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('pages.shortVideoList.likeCount')" width="90" prop="likeCount"/>
      <el-table-column :label="t('pages.shortVideoList.viewCount')" width="100" prop="viewCount"/>
      <el-table-column :label="t('pages.shortVideoList.watchCount')" width="100" prop="watchCount"/>
      <el-table-column :label="t('pages.shortVideoList.totalDiamondIncome')" align="right" min-width="130">
        <template #default="{ row }">{{ formatAmount(row.totalDiamondIncome) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.createdAt')" width="170">
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
import {reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {shortVideoApi} from '@/api'
import type {ShortVideo} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const props = defineProps<{
  userId: string
  active: boolean
}>()

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<ShortVideo[]>([])
const loaded = ref(false)

const searchForm = reactive({
  title: '',
  statusFilter: 0,
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildQueryParams = () => ({
  authorId: props.userId,
  title: searchForm.title.trim(),
  statusFilter: searchForm.statusFilter,
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
    const response = await shortVideoApi.getShortVideoList(buildQueryParams())
    tableData.value = response.data || []
    pagination.total = response.total || 0
    loaded.value = true
  } catch (error) {
    console.error('Failed to load user short videos:', error)
    ElMessage.error(t('pages.shortVideoList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.title = ''
  searchForm.statusFilter = 0
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
  searchForm.title = ''
  searchForm.statusFilter = 0
  pagination.pageIndex = 1
  pagination.total = 0
}

watch(
    () => props.userId,
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

.table-video-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.table-video-preview {
  width: 160px;
  max-height: 90px;
  background: #000;
  border-radius: 4px;
}
</style>

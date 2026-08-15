<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.LiveRoomRecycleBinManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.key"
              clearable
              :placeholder="t('pages.liveRoomRecycleBin.keywordPlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('common.userId')" prop="id" width="180"/>
        <el-table-column :label="t('common.nickname')" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
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
        <el-table-column :label="t('common.phone')" min-width="130" prop="phone">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRoomRecycleBin.guildId')" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRoomRecycleBin.roomTitle')" min-width="140" prop="roomTitle">
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.liveRoomRecycleBin.roomType')" width="100">
          <template #default="{ row }">
            <el-tag :type="categoryTagType(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="170">
          <template #default="{ row }">{{ formatDate(row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button
                v-if="can('onShelf')"
                link
                type="success"
                @click="handleOnShelf(row)"
            >
              {{ t('common.onShelf') }}
            </el-button>
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
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {accountApi} from '@/api'
import type {OffShelfLiveRoomItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_PRIVATE = 3

const {t} = useI18n()
const {can} = usePagePermission('LiveRoomRecycleBinManagement')

const loading = ref(false)
const tableData = ref<OffShelfLiveRoomItem[]>([])
const searchForm = reactive({key: ''})
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const formatDate = (value?: string | null) => {
  if (!value) return '-'
  return value.replace('T', ' ').slice(0, 19)
}

const categoryLabel = (category?: number) => {
  switch (category) {
    case LIVE_ROOM_CATEGORY_HOT:
      return t('pages.liveRoomRecycleBin.categoryHot')
    case LIVE_ROOM_CATEGORY_GAME:
      return t('pages.liveRoomRecycleBin.categoryGame')
    case LIVE_ROOM_CATEGORY_PRIVATE:
      return t('pages.liveRoomRecycleBin.categoryPrivate')
    default:
      return '-'
  }
}

const categoryTagType = (category?: number) => {
  switch (category) {
    case LIVE_ROOM_CATEGORY_HOT:
      return 'danger'
    case LIVE_ROOM_CATEGORY_GAME:
      return 'success'
    case LIVE_ROOM_CATEGORY_PRIVATE:
      return 'warning'
    default:
      return 'info'
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await accountApi.getOffShelfLiveRoomList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key.trim() || undefined,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch off-shelf live room list failed:', error)
    ElMessage.error(t('pages.liveRoomRecycleBin.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.key = ''
  pagination.pageIndex = 1
  fetchList()
}

const handleSizeChange = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleCurrentChange = () => {
  fetchList()
}

const handleOnShelf = async (row: OffShelfLiveRoomItem) => {
  try {
    await ElMessageBox.confirm(
      t('pages.liveRoomRecycleBin.onShelfConfirm', {id: row.id}),
      t('common.onShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await accountApi.setLiveRoomStatus({
      anchorId: row.id,
      status: 1,
    })
    ElMessage.success(t('pages.liveRoomRecycleBin.onShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('on shelf live room failed:', error)
    ElMessage.error(t('pages.liveRoomRecycleBin.onShelfFailed'))
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
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

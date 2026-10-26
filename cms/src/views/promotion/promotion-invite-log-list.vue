<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionInviteLogListManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="100px">
        <el-form-item :label="t('common.keyword')">
          <el-input
              v-model="searchForm.keyword"
              clearable
              :placeholder="t('pages.promotionInviteLogList.keywordPlaceholder')"
              style="width: 280px"
              @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('pages.promotionInviteLogList.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              :end-placeholder="t('common.endTime')"
              format="YYYY-MM-DD"
              :range-separator="t('pages.promotionInviteLogList.dateRangeSeparator')"
              :start-placeholder="t('common.startTime')"
              style="width: 260px"
              type="daterange"
              value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
        <el-table-column :label="t('pages.promotionInviteLogList.invitee')" min-width="280">
          <template #default="{ row }">
            <div class="user-cell">
              <el-image
                  v-if="row.avatar"
                  :preview-src-list="[row.avatar]"
                  :src="row.avatar"
                  fit="cover"
                  hide-on-click-modal
                  preview-teleported
                  class="user-avatar"
              />
              <div class="user-meta">
                <el-button link type="primary" @click="openUserDetail(row.userId)">
                  {{ row.userId }}
                </el-button>
                <div class="user-nickname">{{ row.nickname || '-' }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.promotionInviteLogList.inviter')" min-width="280">
          <template #default="{ row }">
            <div class="user-cell">
              <el-image
                  v-if="row.inviterAvatar"
                  :preview-src-list="[row.inviterAvatar]"
                  :src="row.inviterAvatar"
                  fit="cover"
                  hide-on-click-modal
                  preview-teleported
                  class="user-avatar"
              />
              <div class="user-meta">
                <el-button link type="primary" @click="openUserDetail(row.inviterId)">
                  {{ row.inviterId }}
                </el-button>
                <div class="user-nickname">{{ row.inviterNickname || '-' }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.promotionInviteLogList.inviteTime')" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
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
import {useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {accountApi} from '@/api'
import type {CMSInviteLogItem} from '@/types/api'
import {formatServerDateTime as formatDate, toServerDayEndUnix, toServerDayStartUnix} from '@/utils/server-datetime'

const {t} = useI18n()
const router = useRouter()
const loading = ref(false)
const tableData = ref<CMSInviteLogItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  keyword: '',
  dateRange: [] as string[],
})

const buildQueryParams = () => {
  const [startDate, endDate] = searchForm.dateRange || []
  return {
    keyword: searchForm.keyword.trim() || undefined,
    startTime: startDate ? toServerDayStartUnix(startDate) : 0,
    endTime: endDate ? toServerDayEndUnix(endDate) : 0,
    pageIndex: currentPage.value,
    pageSize: pageSize.value,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await accountApi.getInviteLogList(buildQueryParams())
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load invite log list:', error)
    ElMessage.error(t('pages.promotionInviteLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  void fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.dateRange = []
  currentPage.value = 1
  void fetchList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  void fetchList()
}

const handleCurrentChange = () => {
  void fetchList()
}

const formatRowIndex = (index: number) =>
    (currentPage.value - 1) * pageSize.value + index + 1

const openUserDetail = (userId: string | number | null | undefined) => {
  if (userId == null || userId === '' || String(userId) === '0') {
    return
  }
  router.push({
    name: 'UserDetail',
    query: {id: String(userId)},
  })
}

onMounted(() => {
  void fetchList()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  flex-shrink: 0;
}

.user-meta {
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}

.user-nickname {
  color: var(--el-text-color-regular);
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}
</style>

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CoinMerchantTransferLogList') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="96px">
        <el-form-item :label="t('pages.coinMerchantTransferLogList.coinMerchantId')">
          <el-input
              v-model="searchForm.coinMerchantUserId"
              clearable
              :placeholder="t('pages.coinMerchantTransferLogList.coinMerchantIdPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.coinMerchantTransferLogList.dateRange')">
          <el-date-picker
              v-model="searchForm.dateRange"
              clearable
              type="daterange"
              value-format="YYYY-MM-DD"
              :range-separator="t('pages.coinMerchantTransferLogList.dateRangeSeparator')"
              :start-placeholder="t('pages.coinMerchantTransferLogList.startDate')"
              :end-placeholder="t('pages.coinMerchantTransferLogList.endDate')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.coinMerchantTransferLogList.recordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.coinMerchantId')" min-width="180" prop="coinMerchantUserId">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.coinMerchantUserId"
                link
                type="primary"
                @click="openUserDetail(row.coinMerchantUserId)"
            >
              {{ row.coinMerchantUserId }}
            </el-button>
            <span v-else>{{ row.coinMerchantUserId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.coinMerchantAvatar')" width="88">
          <template #default="{ row }">
            <el-image
                v-if="row.coinMerchantAvatar"
                :preview-src-list="[row.coinMerchantAvatar]"
                :src="row.coinMerchantAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                class="user-avatar"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column
            :label="t('pages.coinMerchantTransferLogList.coinMerchantNickname')"
            min-width="140"
            prop="coinMerchantNickname"
        >
          <template #default="{ row }">{{ row.coinMerchantNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.targetUserId')" min-width="180" prop="targetUserId">
          <template #default="{ row }">
            <el-button
                v-if="canViewUserDetail && row.targetUserId"
                link
                type="primary"
                @click="openUserDetail(row.targetUserId)"
            >
              {{ row.targetUserId }}
            </el-button>
            <span v-else>{{ row.targetUserId || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.targetAvatar')" width="120">
          <template #default="{ row }">
            <el-image
                v-if="row.targetAvatar"
                :preview-src-list="[row.targetAvatar]"
                :src="row.targetAvatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                class="user-avatar"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column
            :label="t('pages.coinMerchantTransferLogList.targetNickname')"
            min-width="140"
            prop="targetNickname"
        >
          <template #default="{ row }">{{ row.targetNickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.amount')" width="120" prop="amount">
          <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.merchantBalance')" min-width="190">
          <template #default="{ row }">
            {{ formatAmount(row.senderGoldBefore) }} → {{ formatAmount(row.senderGoldAfter) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.targetBalance')" min-width="190">
          <template #default="{ row }">
            {{ formatAmount(row.targetGoldBefore) }} → {{ formatAmount(row.targetGoldAfter) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.coinMerchantTransferLogList.transferTime')" width="170">
          <template #default="{ row }">{{ formatServerDateTime(row.createdAt) }}</template>
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {coinMerchantTransferLogApi} from '@/api'
import type {CoinMerchantTransferLogItem} from '@/types/api'
import {useUserDetailNav} from '@/composables/useUserDetailNav'
import {formatAmount} from '@/utils/number-format'
import {
  formatServerDateTime,
  toServerDayEndUnix,
  toServerDayStartUnix,
} from '@/utils/server-datetime'

const {t} = useI18n()
const {canViewUserDetail, openUserDetail} = useUserDetailNav('CoinMerchantTransferLogList')
const loading = ref(false)
const tableData = ref<CoinMerchantTransferLogItem[]>([])

const searchForm = reactive({
  coinMerchantUserId: '',
  dateRange: [] as string[],
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const buildFilterParams = () => {
  const startDate = searchForm.dateRange[0] || ''
  const endDate = searchForm.dateRange[1] || ''
  const hasDateRange = Boolean(startDate && endDate)
  return {
    coinMerchantUserId: searchForm.coinMerchantUserId.trim(),
    startTime: hasDateRange ? toServerDayStartUnix(startDate) : 0,
    endTime: hasDateRange ? toServerDayEndUnix(endDate) : 0,
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await coinMerchantTransferLogApi.getList({
      ...buildFilterParams(),
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load coin merchant transfer logs:', error)
    ElMessage.error(t('pages.coinMerchantTransferLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.coinMerchantUserId = ''
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

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>

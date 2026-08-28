<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GameWinLogListManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="searchForm.userId" clearable :placeholder="t('common.pleaseEnter') + ' ' + t('common.userId')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameWinLogList.gameCode')">
          <el-input v-model="searchForm.gameCode" clearable :placeholder="t('pages.gameWinLogList.gameCodePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameWinLogList.orderId')">
          <el-input v-model="searchForm.orderId" clearable :placeholder="t('pages.gameWinLogList.orderIdPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameWinLogList.platform')">
          <el-input v-model="searchForm.platformType" clearable :placeholder="t('pages.gameWinLogList.platformPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.gameWinLogList.recordId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.userId')" min-width="180" prop="userId"/>
        <el-table-column :label="t('pages.gameWinLogList.userNickname')" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameWinLogList.gameCode')" min-width="120" prop="gameCode"/>
        <el-table-column :label="t('pages.gameWinLogList.nameEn')" min-width="140" prop="nameEn" show-overflow-tooltip>
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameWinLogList.cover')" min-width="100">
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
        <el-table-column :label="t('pages.gameWinLogList.winAmount')" prop="amount" width="120">
          <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameWinLogList.platform')" prop="platformType" width="100"/>
        <el-table-column :label="t('pages.gameWinLogList.orderId')" min-width="200" prop="orderId" show-overflow-tooltip/>
        <el-table-column :label="t('pages.gameWinLogList.time')" width="170">
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
import {ElMessage} from 'element-plus'
import {gameWinLogApi} from '@/api'
import type {GameWinLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const loading = ref(false)
const tableData = ref<GameWinLogItem[]>([])

const searchForm = reactive({
  userId: '',
  gameCode: '',
  orderId: '',
  platformType: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await gameWinLogApi.getGameWinLogList({
      userId: searchForm.userId.trim(),
      gameCode: searchForm.gameCode.trim(),
      orderId: searchForm.orderId.trim(),
      platformType: searchForm.platformType.trim(),
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch game reward log failed:', error)
    ElMessage.error(t('pages.gameWinLogList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.userId = ''
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

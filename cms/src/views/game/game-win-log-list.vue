<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>派彩记录</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.userId" clearable placeholder="请输入用户ID"/>
        </el-form-item>
        <el-form-item label="游戏编码">
          <el-input v-model="searchForm.gameCode" clearable placeholder="游戏编码"/>
        </el-form-item>
        <el-form-item label="订单ID">
          <el-input v-model="searchForm.orderId" clearable placeholder="交易/订单ID"/>
        </el-form-item>
        <el-form-item label="平台">
          <el-input v-model="searchForm.platformType" clearable placeholder="如 ZY"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="记录ID" min-width="180" prop="id"/>
        <el-table-column label="用户ID" min-width="180" prop="userId"/>
        <el-table-column label="用户昵称" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="游戏编码" min-width="120" prop="gameCode"/>
        <el-table-column label="英文名称" min-width="140" prop="nameEn" show-overflow-tooltip>
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column label="封面" min-width="100">
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
        <el-table-column label="派彩金额" prop="amount" width="120">
          <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="平台" prop="platformType" width="100"/>
        <el-table-column label="订单ID" min-width="200" prop="orderId" show-overflow-tooltip/>
        <el-table-column label="时间" width="170">
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {gameWinLogApi} from '@/api'
import type {GameWinLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

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
    console.error('获取派彩记录失败:', error)
    ElMessage.error('获取派彩记录失败')
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

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString('zh-CN')
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
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

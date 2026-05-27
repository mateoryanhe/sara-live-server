<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>主播列表</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item label="关键字">
          <el-input v-model="searchForm.key" clearable placeholder="用户ID/昵称/手机号/分享码"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="用户ID" prop="id" width="180"/>
        <el-table-column label="昵称" min-width="120" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column label="头像" width="80">
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
        <el-table-column label="手机号" min-width="130" prop="phone">
          <template #default="{ row }">{{ row.phone || '-' }}</template>
        </el-table-column>
        <el-table-column label="工会ID" prop="guildId" width="120">
          <template #default="{ row }">{{ row.guildId || '-' }}</template>
        </el-table-column>
        <el-table-column label="登录IP" min-width="140" prop="ip">
          <template #default="{ row }">{{ row.ip || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播间标题" min-width="140" prop="roomTitle" show-overflow-tooltip>
          <template #default="{ row }">{{ row.roomTitle || '-' }}</template>
        </el-table-column>
        <el-table-column label="直播状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.liveStatus === 1 ? 'success' : 'info'">
              {{ row.liveStatus === 1 ? '直播中' : '未开播' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" prop="registeredAt" width="170">
          <template #default="{ row }">{{ formatDate(row.registeredAt) }}</template>
        </el-table-column>
        <el-table-column label="资料更新时间" prop="createdAt" width="170">
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
import {accountApi} from '@/api'
import type {AnchorListItem} from '@/types/api'

const loading = ref(false)
const tableData = ref<AnchorListItem[]>([])

const searchForm = reactive({
  key: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取主播列表失败:', error)
    ElMessage.error('获取主播列表失败')
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

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildRecycleBinManagement') }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item :label="t('pages.guildRecycleBin.guildName')">
          <el-input
              v-model="searchForm.name"
              clearable
              :placeholder="t('pages.guildRecycleBin.guildNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column label="ID" prop="id" width="190"/>
        <el-table-column :label="t('pages.guildRecycleBin.guildName')" prop="name"/>
        <el-table-column :label="t('pages.guildRecycleBin.leader')" width="160" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatLeader(row) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('pages.guildRecycleBin.description')" prop="description" show-overflow-tooltip/>
        <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
        <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
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
import {guildApi} from '@/api'
import type {Guild} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('GuildRecycleBinManagement')

const loading = ref(false)
const tableData = ref<Guild[]>([])
const searchForm = reactive({name: ''})
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const formatLeader = (row: Guild) => {
  if (row.leaderName) {
    return `${row.leaderName} (${row.leaderId})`
  }
  return row.leaderId || '-'
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getOffShelfGuildList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      name: searchForm.name.trim() || undefined,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch off-shelf guild list failed:', error)
    ElMessage.error(t('pages.guildRecycleBin.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchList()
}

const handleReset = () => {
  searchForm.name = ''
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

const handleOnShelf = async (row: Guild) => {
  try {
    await ElMessageBox.confirm(
      t('pages.guildRecycleBin.onShelfConfirm', {name: row.name}),
      t('common.onShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await guildApi.onShelfGuild(row.id)
    ElMessage.success(t('pages.guildRecycleBin.onShelfSuccess'))
    fetchList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('on shelf guild failed:', error)
    ElMessage.error(t('pages.guildRecycleBin.onShelfFailed'))
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

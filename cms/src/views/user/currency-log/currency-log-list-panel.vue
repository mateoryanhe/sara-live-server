<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-form :model="searchForm" class="search-form" inline label-width="80px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="searchForm.userId" clearable :placeholder="t('pages.currencyLogList.userIdPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.currencyLogList.logId')" min-width="180" prop="id"/>
        <el-table-column :label="t('common.userId')" min-width="180" prop="userId">
          <template #default="{ row }">
            <el-button v-if="canViewUserDetail && row.userId" link type="primary" @click="openUserDetail(row.userId)">
              {{ row.userId }}
            </el-button>
            <span v-else>{{ row.userId || '-' }}</span>
          </template>
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
        <el-table-column :label="t('pages.currencyLogList.userNickname')" min-width="140" prop="nickname">
          <template #default="{ row }">{{ row.nickname || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.currencyLogList.changeType')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.action === 1 ? 'success' : 'danger'">
              {{ actionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="amountColumnLabel" prop="amount" width="120">
          <template #default="{ row }">{{ formatAmount(row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.currencyLogList.beforeChange')" prop="before" width="120">
          <template #default="{ row }">{{ formatAmount(row.before) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.currencyLogList.afterChange')" prop="after" width="120">
          <template #default="{ row }">{{ formatAmount(row.after) }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.currencyLogList.reason')" min-width="140" prop="reasonText" show-overflow-tooltip>
          <template #default="{ row }">{{ row.reasonText || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.currencyLogList.time')" width="170">
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
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {currencyLogApi} from '@/api'
import type {CurrencyLogItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'
import {useUserDetailNav} from '@/composables/useUserDetailNav'

const props = defineProps<{
  currencyType: 1 | 2
  pageName: string
}>()

const {t} = useI18n()
const {canViewUserDetail, openUserDetail} = useUserDetailNav(props.pageName)

const pageTitle = computed(() =>
    props.currencyType === 2 ? t('menu.DiamondCurrencyLogList') : t('menu.GoldCurrencyLogList'))
const amountColumnLabel = computed(() =>
    props.currencyType === 2 ? t('pages.currencyLogList.diamondChange') : t('pages.currencyLogList.goldChange'))

const loading = ref(false)
const tableData = ref<CurrencyLogItem[]>([])

const searchForm = reactive({
  userId: '',
})

const pagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const fetchList = async () => {
  loading.value = true
  try {
    const response = await currencyLogApi.getCurrencyLogList({
      userId: searchForm.userId.trim(),
      currencyType: props.currencyType,
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('Failed to load currency logs:', error)
    ElMessage.error(t('pages.currencyLogList.fetchFailed'))
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

const actionLabel = (action: number) => {
  if (action === 1) return t('pages.currencyLogList.actionIncrease')
  if (action === 2) return t('pages.currencyLogList.actionDecrease')
  return '-'
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

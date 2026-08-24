<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.guildProfileMembers.back') }}</el-button>
        </div>
      </template>

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('common.userId')" prop="id" width="180"/>
        <el-table-column :label="t('common.nickname')" min-width="140" prop="nickname" show-overflow-tooltip>
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
        <el-table-column :label="t('pages.guildProfileMembers.unsettledTotalIncome')" align="right" min-width="140">
          <template #default="{ row }">
            <span class="money-amount">{{ formatWalletBalance(row.unsettledTotalIncome) }}</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button
                v-if="can('dailyEffectiveLive')"
                link
                type="primary"
                @click="handleViewDailyLive(row)"
            >
              {{ t('pages.guildProfileMembers.viewDailyEffectiveLive') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
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

      <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildProfileMembers.noAnchors')"/>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import type {MyGuildAnchorListItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const {can} = usePagePermission('GuildProfileManagement')

const loading = ref(false)
const tableData = ref<MyGuildAnchorListItem[]>([])
const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0,
})

const guildId = computed(() => {
  const value = route.query.guildId
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
})

const guildName = computed(() => {
  const value = route.query.guildName
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
})

const pageTitle = computed(() => {
  if (guildName.value) {
    return t('pages.guildProfileMembers.titleWithName', {name: guildName.value})
  }
  return t('pages.guildProfileMembers.title')
})

const fetchList = async () => {
  if (!guildId.value) {
    tableData.value = []
    pagination.total = 0
    return
  }
  loading.value = true
  try {
    const response = await guildApi.getMyGuildAnchorList({
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      guildId: guildId.value,
    })
    tableData.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('fetch guild profile members failed:', error)
    tableData.value = []
    pagination.total = 0
    ElMessage.error(t('pages.guildProfileMembers.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchList()
}

const handleCurrentChange = (page: number) => {
  pagination.pageIndex = page
  fetchList()
}

const goBack = () => {
  router.push('/operation/guild/guild-profile')
}

const handleViewDailyLive = (row: MyGuildAnchorListItem) => {
  const id = row?.id != null ? String(row.id) : ''
  if (!id || !guildId.value) {
    return
  }
  router.push({
    path: '/operation/guild/guild-profile-anchor-daily-live',
    query: {
      guildId: guildId.value,
      guildName: guildName.value,
      anchorId: id,
      nickname: row.nickname || '',
    },
  })
}

watch(guildId, () => {
  pagination.pageIndex = 1
  fetchList()
})

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
  font-size: 16px;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.money-amount {
  font-variant-numeric: tabular-nums;
}
</style>

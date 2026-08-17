<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildProfileManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.guildProfile.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.guildProfile.tipLine1') }}</p>
        <p>{{ t('pages.guildProfile.tipLine2') }}</p>
      </el-alert>

      <div class="content">
        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="190"/>
          <el-table-column :label="t('pages.guildProfile.guildName')" min-width="140" prop="name"/>
          <el-table-column :label="t('pages.guildProfile.description')" min-width="180" prop="description" show-overflow-tooltip/>
          <el-table-column :label="t('pages.guildProfile.lastUpdated')" prop="updatedAt" width="170"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="260">
            <template #default="{ row }">
              <el-button
                  v-if="can('viewAnchors')"
                  link
                  type="primary"
                  @click="handleViewAnchors(row)"
              >
                {{ t('pages.guildProfile.viewAnchors') }}
              </el-button>
              <el-button
                  v-if="can('viewAnchorSettlementLogs')"
                  link
                  type="primary"
                  @click="handleViewSettlementLogs(row)"
              >
                {{ t('pages.guildProfile.viewAnchorSettlementLogs') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildProfile.noGuild')"/>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {guildApi} from '@/api'
import type {MyGuildProfile} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('GuildProfileManagement')
const loading = ref(false)
const tableData = ref<MyGuildProfile[]>([])

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getMyGuildProfile()
    tableData.value = response?.list ?? []
  } catch (error) {
    console.error('fetch guild profile list failed:', error)
    tableData.value = []
  } finally {
    loading.value = false
  }
}

const handleViewAnchors = (row: MyGuildProfile) => {
  const id = row?.id != null ? String(row.id) : ''
  if (!id) {
    return
  }
  router.push({
    path: '/operation/guild/guild-profile-members',
    query: {
      guildId: id,
      guildName: row.name || '',
    },
  })
}

const handleViewSettlementLogs = (row: MyGuildProfile) => {
  const id = row?.id != null ? String(row.id) : ''
  router.push({
    path: '/operation/guild/guild-anchor-income-settlement-log-list',
    query: id ? {guildId: id} : {},
  })
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
  font-size: 16px;
  font-weight: bold;
}

.tip-alert {
  margin-bottom: 20px;
}

.content {
  margin-top: 4px;
}
</style>

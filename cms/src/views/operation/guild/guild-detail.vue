<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.guildList.back') }}</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <el-empty v-if="!guildBasic" :description="t('pages.guildList.detailNotFound')"/>
        <el-tabs v-else v-model="activeTab">
          <el-tab-pane :label="t('pages.guildList.tabBasic')" name="basic">
            <el-descriptions :column="1" border>
              <el-descriptions-item label="ID">{{ guildBasic.id }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildList.guildName')">{{ guildBasic.name || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildList.leader')">{{ formatLeader(guildBasic) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildList.description')">{{ guildBasic.description || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.guildList.shelfStatus')">
                <el-tag v-if="guildBasic.status === 1" type="success">{{ t('common.onShelf') }}</el-tag>
                <el-tag v-else type="info">{{ t('common.offShelf') }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.createdAt')">{{ formatDate(guildBasic.createdAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('common.updatedAt')">{{ formatDate(guildBasic.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildList.tabIncomeUnsettled')" name="incomeUnsettled">
            <IncomePanel :data="incomeData?.incomeUnsettled" :updated-at="incomeData?.incomeUnsettled?.updatedAt"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildList.tabIncomeSettled')" name="incomeSettled">
            <IncomePanel
                :data="incomeData?.incomeSettled"
                show-settlement-share
                :settlement-share-amount="incomeData?.incomeSettled?.settlementShareAmount"
                :updated-at="incomeData?.incomeSettled?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildList.tabIncomeTotal')" name="incomeTotal">
            <IncomePanel
                :data="incomeData?.incomeTotal"
                show-settlement-share
                :settlement-share-amount="incomeData?.incomeTotal?.settlementShareAmount"
                :updated-at="incomeData?.incomeTotal?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildList.tabIncomeArchive')" name="incomeArchive">
            <ArchivePanel :active="activeTab === 'incomeArchive'" :guild-id="guildId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.guildList.tabSettlementLog')" name="settlementLog">
            <SettlementLogPanel :active="activeTab === 'settlementLog'" :guild-id="guildId"/>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {guildApi} from '@/api'
import IncomePanel from './guild-detail-income-panel.vue'
import ArchivePanel from './guild-detail-archive-panel.vue'
import SettlementLogPanel from './guild-detail-settlement-log-panel.vue'
import type {Guild, GuildDetailIncome} from '@/types/api'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const incomeData = ref<GuildDetailIncome | null>(null)
const activeTab = ref('basic')

const parseQueryValue = (key: string) => {
  const value = route.query[key]
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
}

const guildId = computed(() => parseQueryValue('id'))

const guildBasic = computed<Guild | null>(() => {
  if (!guildId.value) {
    return null
  }
  const status = Number(parseQueryValue('status'))
  return {
    id: guildId.value,
    name: parseQueryValue('name'),
    leaderId: parseQueryValue('leaderId'),
    leaderName: parseQueryValue('leaderName'),
    description: parseQueryValue('description'),
    status: Number.isNaN(status) ? 0 : status,
    createdAt: parseQueryValue('createdAt'),
    updatedAt: parseQueryValue('updatedAt'),
  }
})

const pageTitle = computed(() => {
  if (guildBasic.value?.name) {
    return t('pages.guildList.detailTitleWithName', {name: guildBasic.value.name})
  }
  if (guildId.value) {
    return t('pages.guildList.detailTitleWithId', {id: guildId.value})
  }
  return t('pages.guildList.detailTitle')
})

const formatLeader = (guild: Guild) => {
  if (guild.leaderName) {
    return `${guild.leaderName} (${guild.leaderId})`
  }
  return guild.leaderId || '-'
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const fetchIncome = async () => {
  if (!guildId.value) {
    incomeData.value = null
    return
  }
  loading.value = true
  try {
    incomeData.value = await guildApi.getGuildDetail(guildId.value)
  } catch (error) {
    console.error('Failed to load guild income detail:', error)
    incomeData.value = null
    ElMessage.error(t('pages.guildList.detailFetchFailed'))
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'GuildManagement'})
}

watch(guildId, (_id, prev) => {
  if (prev !== undefined) {
    activeTab.value = 'basic'
    fetchIncome()
  }
})

onActivated(() => {
  fetchIncome()
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
</style>

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('pages.guildAnchorImportResult.title') }}</span>
          <el-button @click="goBack">{{ t('pages.guildAnchorImportResult.back') }}</el-button>
        </div>
      </template>

      <template v-if="result">
        <div class="summary">
          <el-tag type="info">{{ t('pages.guildAnchorImportResult.guild') }}: {{ result.guildName }} ({{ result.guildId }})</el-tag>
          <el-tag>
            {{ t('pages.guildAnchorImportResult.anchorType') }}:
            {{
              result.anchorType === 7
                  ? t('pages.guildAnchorImportResult.seniorAnchor')
                  : t('pages.guildAnchorImportResult.normalAnchor')
            }}
          </el-tag>
          <el-tag type="success">{{ t('pages.guildAnchorImportResult.successCount') }}: {{ result.successCount }}</el-tag>
          <el-tag type="danger">{{ t('pages.guildAnchorImportResult.failCount') }}: {{ result.failCount }}</el-tag>
        </div>

        <el-table :data="result.fails" style="width: 100%">
          <el-table-column :label="t('pages.guildAnchorImportResult.userId')" min-width="160" prop="userId"/>
          <el-table-column :label="t('pages.guildAnchorImportResult.nickname')" min-width="160" prop="nickname" show-overflow-tooltip/>
          <el-table-column :label="t('menu.UserDetail')" width="110">
            <template #default="{ row }">
              <el-button v-if="canViewUserDetail && row.userId" link type="primary" @click="openUserDetail(row.userId)">
                {{ t('pages.userList.viewDetail') }}
              </el-button>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildAnchorImportResult.reason')" min-width="200">
            <template #default="{ row }">
              {{ formatReason(row.reason) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="result.fails.length === 0" :description="t('pages.guildAnchorImportResult.empty')"/>
      </template>
      <el-empty v-else :description="t('pages.guildAnchorImportResult.noData')"/>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import type {GuildAnchorImportResultState} from '@/types/api'
import {useUserDetailNav} from '@/composables/useUserDetailNav'

const GUILD_ANCHOR_IMPORT_RESULT_KEY = 'guildAnchorImportResult'

const {t} = useI18n()
const router = useRouter()
const {canViewUserDetail, openUserDetail} = useUserDetailNav('GuildManagement')
const result = ref<GuildAnchorImportResultState | null>(null)

const formatReason = (reason: number) => {
  switch (reason) {
    case 1:
      return t('pages.guildAnchorImportResult.reasonUserNotFound')
    case 2:
      return t('pages.guildAnchorImportResult.reasonCancelCodeMismatch')
    case 3:
      return t('pages.guildAnchorImportResult.reasonCancelCodeExpired')
    case 4:
      return t('pages.guildAnchorImportResult.reasonAlreadyInGuild')
    case 5:
      return t('pages.guildAnchorImportResult.reasonCannotSetAnchor')
    case 6:
      return t('pages.guildAnchorImportResult.reasonAlreadyHasLiveRoom')
    default:
      return t('pages.guildAnchorImportResult.reasonUnknown')
  }
}

const goBack = () => {
  router.push({name: 'GuildManagement'})
}

onMounted(() => {
  const raw = sessionStorage.getItem(GUILD_ANCHOR_IMPORT_RESULT_KEY)
  if (!raw) {
    return
  }
  try {
    result.value = JSON.parse(raw) as GuildAnchorImportResultState
  } catch (error) {
    console.error('parse import result failed:', error)
    result.value = null
  }
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

.summary {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}
</style>

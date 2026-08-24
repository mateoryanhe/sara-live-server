<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.guildProfileMembers.back') }}</el-button>
        </div>
      </template>

      <DailyLivePanel
          :anchor-id="anchorId"
          :guild-id="guildId"
          permission-page="GuildProfileManagement"
          simple-columns
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import DailyLivePanel from '@/views/user/anchor/anchor-detail-daily-live-panel.vue'

const {t} = useI18n()
const router = useRouter()
const currentRoute = useRoute()

const readQuery = (key: string) => {
  const value = currentRoute.query[key]
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
}

const guildId = computed(() => readQuery('guildId'))
const guildName = computed(() => readQuery('guildName'))
const anchorId = computed(() => readQuery('anchorId'))
const anchorNickname = computed(() => readQuery('nickname'))

const pageTitle = computed(() => {
  const name = anchorNickname.value || anchorId.value
  if (name) {
    return t('pages.guildProfileMembers.dailyLiveTitleWithName', {name})
  }
  return t('pages.guildProfileMembers.dailyLiveTitle')
})

const goBack = () => {
  router.push({
    path: '/operation/guild/guild-profile-members',
    query: {
      guildId: guildId.value,
      guildName: guildName.value,
    },
  })
}
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

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.gameShelfList.back') }}</el-button>
        </div>
      </template>

      <div v-loading="loading" class="config-body">
        <el-empty v-if="!loading && !configUrl" :description="t('pages.gameShelfList.vendorConfigEmpty')"/>
        <iframe
            v-else-if="configUrl"
            :src="configUrl"
            class="vendor-config-iframe"
            title="vendor-config"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {gamePlatformApi} from '@/api/modules/gamePlatform'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const configUrl = ref('')

const gameCode = computed(() => (typeof route.query.gameCode === 'string' ? route.query.gameCode.trim() : ''))
const platform = computed(() => (typeof route.query.platform === 'string' ? route.query.platform.trim() : ''))
const gameName = computed(() => (typeof route.query.name === 'string' ? route.query.name.trim() : ''))

const pageTitle = computed(() => {
  const name = gameName.value || gameCode.value || t('pages.gameShelfList.vendorConfig')
  return t('pages.gameShelfList.vendorConfigTitle', {name})
})

const fetchConfigUrl = async () => {
  if (!gameCode.value || !platform.value) {
    configUrl.value = ''
    return
  }
  loading.value = true
  configUrl.value = ''
  try {
    const response = await gamePlatformApi.getMultiplayerConfigUrl({
      gameCode: gameCode.value,
      platform: platform.value,
    })
    const url = response?.configUrl?.trim()
    if (!url) {
      ElMessage.error(t('pages.gameShelfList.vendorConfigEmpty'))
      return
    }
    configUrl.value = url
  } catch (error) {
    console.error('get multiplayer config url failed:', error)
    ElMessage.error(t('pages.gameShelfList.vendorConfigFailed'))
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'GameShelfListManagement'})
}

watch([gameCode, platform], () => {
  fetchConfigUrl()
})

onMounted(() => {
  fetchConfigUrl()
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.config-body {
  min-height: calc(100vh - 180px);
}

.vendor-config-iframe {
  width: 100%;
  height: calc(100vh - 180px);
  border: none;
  display: block;
}
</style>

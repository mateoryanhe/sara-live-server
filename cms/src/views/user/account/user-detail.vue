<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.userList.back') }}</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <el-empty v-if="!loading && !detail?.account" :description="t('pages.userList.detailNotFound')"/>
        <el-tabs v-else-if="detail" v-model="activeTab">
          <el-tab-pane :label="t('pages.userList.tabBasic')" name="basic">
            <el-tabs v-model="basicSubTab" class="basic-sub-tabs">
              <el-tab-pane :label="t('pages.userList.basicSubTabProfile')" name="profile">
                <el-descriptions :column="1" border>
                  <el-descriptions-item :label="t('common.userId')">{{ detail.account?.id || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.openId')">{{ detail.account?.openId || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.phoneAreaCode')">{{ detail.account?.phoneAreaCode || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.phone')">{{ detail.profile?.phone || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.avatar')">
                    <el-image
                        v-if="detail.profile?.avatar"
                        :preview-src-list="[detail.profile.avatar]"
                        :src="detail.profile.avatar"
                        fit="cover"
                        hide-on-click-modal
                        preview-teleported
                        style="width:48px;height:48px;border-radius:50%"
                    />
                    <span v-else>-</span>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('common.nickname')">{{ detail.profile?.nickname || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('common.remark')">{{ detail.profile?.remark || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.shareCode')">{{ detail.profile?.shareCode || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.gender')">{{ genderLabel(detail.profile?.gender) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.birthday')">{{ formatDate(detail.profile?.birthday) }}</el-descriptions-item>
                </el-descriptions>
              </el-tab-pane>

              <el-tab-pane :label="t('pages.userList.basicSubTabIdentity')" name="identity">
                <el-descriptions :column="1" border>
                  <el-descriptions-item :label="t('pages.userList.userType')">
                    <el-tag>{{ userTypeLabel(detail.profile?.userType) }}</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.isAnchor')">
                    <el-tag :type="detail.profile?.isAnchor ? 'success' : 'info'">
                      {{ detail.profile?.isAnchor ? t('common.yes') : t('common.no') }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.inviterId')">{{ detail.profile?.inviterId || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.vipLevel')">{{ detail.profile?.vipLevel ?? '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.guildId')">{{ detail.profile?.guildId || '-' }}</el-descriptions-item>
                  <el-descriptions-item v-if="showBotAnchorStatus" :label="t('pages.userList.botAnchorStatus')">
                    {{ botAnchorStatusLabel(detail.profile?.botAnchorStatus) }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.liveRoomId')">{{ detail.profile?.liveRoomId || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.liveRoomVer')">{{ detail.profile?.liveRoomVer || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.channel')">{{ channelLabel(detail.account?.channel) }}</el-descriptions-item>
                </el-descriptions>
              </el-tab-pane>

              <el-tab-pane :label="t('pages.userList.basicSubTabLogin')" name="login">
                <el-descriptions :column="1" border>
                  <el-descriptions-item :label="t('pages.userList.lastLoginTime')">{{ formatDate(detail.profile?.lastLoginTime) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.ip')">{{ detail.account?.ip || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.loginCountry')">{{ detail.account?.loginCountry || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.registerIp')">{{ detail.account?.registerIp || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.registerCountry')">{{ detail.account?.registerCountry || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.registeredAt')">{{ formatDate(detail.account?.createdAt) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.profileUpdatedAt')">{{ formatDate(detail.profile?.updatedAt) }}</el-descriptions-item>
                </el-descriptions>
              </el-tab-pane>

              <el-tab-pane :label="t('pages.userList.basicSubTabStatus')" name="status">
                <el-descriptions :column="1" border>
                  <el-descriptions-item :label="t('pages.userList.banStatus')">
                    <el-tag v-if="detail.account?.ban" type="danger">{{ t('pages.userList.banned') }}</el-tag>
                    <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.banTime')">{{ formatDate(detail.account?.banApplyTime) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.banAppliedAt')">{{ formatDate(detail.account?.banTime) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.userList.cancelStatus')">
                    <el-tag v-if="detail.account?.cancel" type="warning">{{ t('pages.userList.canceled') }}</el-tag>
                    <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
                  </el-descriptions-item>
                </el-descriptions>
              </el-tab-pane>
            </el-tabs>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabWallet')" name="wallet">
            <div class="wallet-panel">
              <div class="wallet-balance-card wallet-balance-gold">
                <div class="wallet-balance-label">{{ t('pages.userList.goldBalance') }}</div>
                <div class="wallet-balance-value">{{ formatWalletBalance(detail.wallet?.gold) }}</div>
              </div>
              <div class="wallet-balance-card wallet-balance-diamond">
                <div class="wallet-balance-label">{{ t('pages.userList.diamondBalance') }}</div>
                <div class="wallet-balance-value">{{ formatWalletBalance(detail.wallet?.diamond) }}</div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabUserExt')" name="userExt">
            <el-empty v-if="!detail.userExt" :description="t('pages.userList.noData')"/>
            <el-descriptions v-else :column="1" border>
              <el-descriptions-item :label="t('pages.userList.canRank')">
                <el-tag :type="detail.userExt.canRank !== false ? 'success' : 'info'">
                  {{ detail.userExt.canRank !== false ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.prettyId')">{{ detail.userExt.prettyId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.packageName')">{{ detail.userExt.packageName || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.appVersion')">{{ detail.userExt.appVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.followCount')">{{ detail.userExt.followCount ?? '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.followerCount')">{{ detail.userExt.followerCount ?? '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.rechargeWhitelist')">
                <el-tag :type="detail.userExt.rechargeWhitelist ? 'success' : 'info'">
                  {{ detail.userExt.rechargeWhitelist ? t('common.yes') : t('common.no') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.shortVideoUnsettledIncome')">
                <span class="money-amount">{{ formatWalletBalance(detail.userExt.shortVideoUnsettledIncome) }}</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.cancelCode')">{{ detail.userExt.cancelCode || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.cancelCodeExpireAt')">{{ formatDate(detail.userExt.cancelCodeExpireAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.extUpdatedAt')">{{ formatDate(detail.userExt.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabCumulativeStat')" name="cumulativeStat">
            <div>
              <div class="wallet-panel">
                <div class="wallet-balance-card stat-recharge">
                  <div class="wallet-balance-label">{{ t('pages.userList.totalRecharge') }}</div>
                  <div class="wallet-balance-value">{{ formatWalletBalance(cumulativeStat.totalRecharge) }}</div>
                </div>
                <div class="wallet-balance-card stat-withdraw">
                  <div class="wallet-balance-label">{{ t('pages.userList.totalWithdraw') }}</div>
                  <div class="wallet-balance-value">{{ formatWalletBalance(cumulativeStat.totalWithdraw) }}</div>
                </div>
                <div class="wallet-balance-card stat-pay-count">
                  <div class="wallet-balance-label">{{ t('pages.userList.totalPayCount') }}</div>
                  <div class="wallet-balance-value">{{ formatStatCount(cumulativeStat.totalPayCount) }}</div>
                </div>
                <div class="wallet-balance-card stat-diamond-consume">
                  <div class="wallet-balance-label">{{ t('pages.userList.totalDiamondConsume') }}</div>
                  <div class="wallet-balance-value">{{ formatWalletBalance(cumulativeStat.totalDiamondConsume) }}</div>
                </div>
                <div class="wallet-balance-card stat-gold-consume">
                  <div class="wallet-balance-label">{{ t('pages.userList.totalGoldConsume') }}</div>
                  <div class="wallet-balance-value">{{ formatWalletBalance(cumulativeStat.totalGoldConsume) }}</div>
                </div>
              </div>
              <div class="stat-updated-at">
                {{ t('pages.userList.statUpdatedAt') }}：{{ formatDate(cumulativeStat.updatedAt) }}
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabLoginDevice')" name="loginDevice">
            <el-empty v-if="!detail.loginDevice" :description="t('pages.userList.noData')"/>
            <el-descriptions v-else :column="1" border>
              <el-descriptions-item :label="t('pages.userList.deviceType')">{{ detail.loginDevice.deviceType || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.deviceModel')">{{ detail.loginDevice.deviceModel || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.cpuModel')">{{ detail.loginDevice.cpuModel || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.osVersion')">{{ detail.loginDevice.osVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.appVersion')">{{ detail.loginDevice.appVersion || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.deviceId')">{{ detail.loginDevice.deviceId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.deviceUpdatedAt')">{{ formatDate(detail.loginDevice.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane v-if="canViewGoldLog" :label="t('pages.userList.tabGoldLog')" name="goldLog">
            <CurrencyLogPanel
                :active="activeTab === 'goldLog'"
                :currency-type="1"
                export-permission="exportGoldLog"
                :user-id="userId"
            />
          </el-tab-pane>

          <el-tab-pane v-if="canViewDiamondLog" :label="t('pages.userList.tabDiamondLog')" name="diamondLog">
            <CurrencyLogPanel
                :active="activeTab === 'diamondLog'"
                :currency-type="2"
                export-permission="exportDiamondLog"
                :user-id="userId"
            />
          </el-tab-pane>

          <el-tab-pane v-if="canViewGameBetLog" :label="t('pages.userList.tabGameBetLog')" name="gameBetLog">
            <GameLogPanel
                :active="activeTab === 'gameBetLog'"
                export-permission="exportGameBetLog"
                log-type="bet"
                :user-id="userId"
            />
          </el-tab-pane>

          <el-tab-pane v-if="canViewGameWinLog" :label="t('pages.userList.tabGameWinLog')" name="gameWinLog">
            <GameLogPanel
                :active="activeTab === 'gameWinLog'"
                export-permission="exportGameWinLog"
                log-type="win"
                :user-id="userId"
            />
          </el-tab-pane>

          <el-tab-pane v-if="canViewShortVideo" :label="t('pages.userList.tabShortVideo')" name="shortVideo">
            <ShortVideoPanel :active="activeTab === 'shortVideo'" :user-id="userId"/>
          </el-tab-pane>

          <el-tab-pane v-if="canViewShortVideoWatch" :label="t('pages.userList.tabShortVideoWatch')" name="shortVideoWatch">
            <ShortVideoWatchPanel :active="activeTab === 'shortVideoWatch'" :user-id="userId"/>
          </el-tab-pane>

          <el-tab-pane v-if="canViewShortVideoPurchase" :label="t('pages.userList.tabShortVideoPurchase')" name="shortVideoPurchase">
            <ShortVideoWatchPanel :active="activeTab === 'shortVideoPurchase'" only-paid :user-id="userId"/>
          </el-tab-pane>

          <el-tab-pane
              v-if="canViewShortVideoAuthorSettlementLog"
              :label="t('pages.userList.tabShortVideoAuthorSettlementLog')"
              name="shortVideoAuthorSettlementLog"
          >
            <ShortVideoAuthorSettlementLogPanel
                :active="activeTab === 'shortVideoAuthorSettlementLog'"
                :user-id="userId"
            />
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
import {accountApi} from '@/api'
import CurrencyLogPanel from './user-detail-currency-log-panel.vue'
import GameLogPanel from './user-detail-game-log-panel.vue'
import ShortVideoAuthorSettlementLogPanel from './user-detail-short-video-author-settlement-log-panel.vue'
import ShortVideoWatchPanel from './user-detail-short-video-watch-panel.vue'
import ShortVideoPanel from '../anchor/anchor-detail-short-video-panel.vue'
import type {UserCumulativeStatDetailItem, UserDetail} from '@/types/api.ts'
import {formatStatCount, formatWalletBalance} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('UserDetail')
const canViewGoldLog = computed(() => can('goldLog'))
const canViewDiamondLog = computed(() => can('diamondLog'))
const canViewGameBetLog = computed(() => can('gameBetLog'))
const canViewGameWinLog = computed(() => can('gameWinLog'))
const canViewShortVideo = computed(() => can('shortVideo'))
const canViewShortVideoWatch = computed(() => can('shortVideoWatch'))
const canViewShortVideoPurchase = computed(() => can('shortVideoPurchase'))
const canViewShortVideoAuthorSettlementLog = computed(() => can('shortVideoAuthorSettlementLog'))
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<UserDetail | null>(null)
const activeTab = ref('basic')
const basicSubTab = ref('profile')

const USER_TYPE_BOT_ANCHOR = 2

const userId = computed(() => String(route.query.id || '').trim())

const showBotAnchorStatus = computed(() => detail.value?.profile?.userType === USER_TYPE_BOT_ANCHOR)

const cumulativeStat = computed<UserCumulativeStatDetailItem>(() => detail.value?.cumulativeStat ?? {})

const pageTitle = computed(() => {
  if (detail.value?.profile?.nickname) {
    return t('pages.userList.detailTitleWithName', {name: detail.value.profile.nickname})
  }
  if (userId.value) {
    return t('pages.userList.detailTitleWithId', {id: userId.value})
  }
  return t('pages.userList.detailTitle')
})

const userTypeLabelMap = computed<Record<number, string>>(() => ({
  0: t('pages.userList.userTypeNormal'),
  1: t('pages.userList.userTypeAnchor'),
  2: t('pages.userList.userTypeBotAnchor'),
  3: t('pages.userList.userTypeBotViewer'),
  4: t('pages.userList.userTypeTester'),
  5: t('pages.userList.userTypeCmsAuthor'),
  7: t('pages.userList.userTypeSeniorAnchor'),
}))

const userTypeLabel = (userType?: number) => {
  if (userType === undefined || userType === null) return '-'
  return userTypeLabelMap.value[userType] || t('pages.userList.userTypeNormal')
}

const channelLabelMap = computed<Record<number, string>>(() => ({
  1: t('pages.userList.channelTest'),
  2: t('pages.userList.channelPhone'),
  3: t('pages.userList.channelBotAnchor'),
  4: t('pages.userList.channelDevice'),
  5: t('pages.userList.channelShortVideoAuthor'),
  6: t('pages.userList.channelH5Device'),
}))

const channelLabel = (channel?: number) => {
  if (channel === undefined || channel === null) return '-'
  return channelLabelMap.value[channel] || t('pages.userList.channelUnknown')
}

const genderLabel = (gender?: number) => {
  if (gender === 1) return t('pages.userList.genderMale')
  if (gender === 2) return t('pages.userList.genderFemale')
  if (gender === 0) return t('pages.userList.genderUnknown')
  return '-'
}

const botAnchorStatusLabel = (status?: number) => {
  if (status === 1) return t('pages.userList.botAnchorEnabled')
  if (status === 0) return t('pages.userList.botAnchorDisabled')
  return '-'
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const isUserDetailRoute = () => route.name === 'UserDetail'

const fetchDetail = async () => {
  if (!isUserDetailRoute()) {
    return
  }
  if (!userId.value) {
    detail.value = null
    return
  }
  loading.value = true
  try {
    detail.value = await accountApi.getUserDetail(userId.value)
    if (!detail.value?.account) {
      ElMessage.warning(t('pages.userList.detailNotFound'))
    }
  } catch (error) {
    console.error('Failed to load user detail:', error)
    detail.value = null
    ElMessage.error(t('pages.userList.detailFetchFailed'))
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'UserList'})
}

watch(userId, (_id, prev) => {
  if (!isUserDetailRoute()) {
    return
  }
  if (prev !== undefined) {
    activeTab.value = 'basic'
    basicSubTab.value = 'profile'
    fetchDetail()
  }
})

onActivated(() => {
  if (!isUserDetailRoute()) {
    return
  }
  fetchDetail()
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.basic-sub-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.wallet-panel {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.wallet-balance-card {
  flex: 1 1 240px;
  max-width: 360px;
  padding: 20px 24px;
  border-radius: 8px;
  border: 1px solid var(--el-border-color-light);
  background: var(--el-fill-color-blank);
}

.wallet-balance-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-bottom: 12px;
}

.wallet-balance-value {
  font-size: 28px;
  font-weight: 600;
  line-height: 1.2;
  font-variant-numeric: tabular-nums;
  word-break: break-all;
}

.wallet-balance-gold .wallet-balance-value {
  color: #d48806;
}

.wallet-balance-diamond .wallet-balance-value {
  color: #1677ff;
}

.stat-recharge .wallet-balance-value {
  color: #389e0d;
}

.stat-withdraw .wallet-balance-value {
  color: #cf1322;
}

.stat-pay-count .wallet-balance-value {
  color: #531dab;
}

.stat-diamond-consume .wallet-balance-value {
  color: #1677ff;
}

.stat-gold-consume .wallet-balance-value {
  color: #d48806;
}

.stat-updated-at {
  margin-top: 16px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
</style>

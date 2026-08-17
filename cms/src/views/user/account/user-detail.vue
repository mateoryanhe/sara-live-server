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
              <el-descriptions-item :label="t('pages.userList.gender')">{{ genderLabel(detail.profile?.gender) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.birthday')">{{ formatDate(detail.profile?.birthday) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.botAnchorStatus')">
                {{ botAnchorStatusLabel(detail.profile?.botAnchorStatus) }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.liveRoomId')">{{ detail.profile?.liveRoomId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.liveRoomVer')">{{ detail.profile?.liveRoomVer || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.lastLoginTime')">{{ formatDate(detail.profile?.lastLoginTime) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.ip')">{{ detail.account?.ip || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.loginCountry')">{{ detail.account?.loginCountry || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.registerIp')">{{ detail.account?.registerIp || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.registerCountry')">{{ detail.account?.registerCountry || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.channel')">{{ detail.account?.channel ?? '-' }}</el-descriptions-item>
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
              <el-descriptions-item :label="t('pages.userList.registeredAt')">{{ formatDate(detail.account?.createdAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.profileUpdatedAt')">{{ formatDate(detail.profile?.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabWallet')" name="wallet">
            <el-descriptions :column="1" border>
              <el-descriptions-item :label="t('pages.userList.goldBalance')">{{ formatAmount(detail.wallet?.gold) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.diamondBalance')">{{ formatAmount(detail.wallet?.diamond) }}</el-descriptions-item>
            </el-descriptions>
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
              <el-descriptions-item :label="t('pages.userList.cancelCode')">{{ detail.userExt.cancelCode || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.cancelCodeExpireAt')">{{ formatDate(detail.userExt.cancelCodeExpireAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.extUpdatedAt')">{{ formatDate(detail.userExt.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.userList.tabCumulativeStat')" name="cumulativeStat">
            <el-empty v-if="!detail.cumulativeStat" :description="t('pages.userList.noData')"/>
            <el-descriptions v-else :column="1" border>
              <el-descriptions-item :label="t('pages.userList.totalRecharge')">{{ formatAmount(detail.cumulativeStat.totalRecharge) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.totalWithdraw')">{{ formatAmount(detail.cumulativeStat.totalWithdraw) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.totalPayCount')">{{ detail.cumulativeStat.totalPayCount ?? '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.totalDiamondConsume')">{{ formatAmount(detail.cumulativeStat.totalDiamondConsume) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.totalGoldConsume')">{{ formatAmount(detail.cumulativeStat.totalGoldConsume) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.totalLiveDuration')">{{ formatDuration(detail.cumulativeStat.totalLiveDuration) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.userList.statUpdatedAt')">{{ formatDate(detail.cumulativeStat.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
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
import type {UserDetail} from '@/types/api.ts'
import {formatAmount} from '@/utils/number-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<UserDetail | null>(null)
const activeTab = ref('basic')

const userId = computed(() => String(route.query.id || '').trim())

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

const formatDuration = (seconds?: number) => {
  if (seconds === undefined || seconds === null) return '-'
  const total = Math.max(0, Math.floor(seconds))
  const hours = Math.floor(total / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const secs = total % 60
  return t('pages.userList.durationFormat', {hours, minutes, seconds: secs})
}

const fetchDetail = async () => {
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
  if (prev !== undefined) {
    activeTab.value = 'basic'
    fetchDetail()
  }
})

onActivated(() => {
  fetchDetail()
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>

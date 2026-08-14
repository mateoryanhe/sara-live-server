<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button @click="goBack">{{ t('pages.anchorList.back') }}</el-button>
        </div>
      </template>

      <div v-loading="loading">
        <el-descriptions v-if="detail" :column="1" border>
          <el-descriptions-item :label="t('common.userId')">{{ detail.id }}</el-descriptions-item>
          <el-descriptions-item :label="t('common.nickname')">{{ detail.nickname || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('common.avatar')">
            <el-image
                v-if="detail.avatar"
                :preview-src-list="[detail.avatar]"
                :src="detail.avatar"
                fit="cover"
                hide-on-click-modal
                preview-teleported
                style="width:48px;height:48px;border-radius:50%"
            />
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('common.phone')">{{ detail.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.guildId')">{{ detail.guildId || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.anchorType')">
            <el-tag :type="anchorTypeTagType(detail.userType)">{{ anchorTypeLabel(detail.userType) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.loginIp')">{{ detail.ip || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.liveRoom')">{{ detail.roomId || detail.id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.roomType')">
            <el-tag :type="categoryTagType(detail.category)">{{ categoryLabel(detail.category) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.privateInviteType')">
            <template v-if="detail.category === LIVE_ROOM_CATEGORY_HOT">
              <el-tag :type="privateInviteTagType(detail.privateInviteType)">
                {{ privateInviteLabel(detail.privateInviteType) }}
              </el-tag>
            </template>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.ticketPrice')">
            {{ isPrivateRoom(detail.category) ? formatAmount(detail.ticket) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.billingPricePerMinute')">
            {{ isPrivateRoom(detail.category) ? formatAmount(detail.billing) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.roomTitle')">{{ detail.roomTitle || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.liveStatus')">
            <el-tag :type="detail.liveStatus === 1 ? 'success' : 'info'">
              {{ detail.liveStatus === 1 ? t('common.live') : t('common.offline') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.liveIncome')">{{ formatAmount(detail.totalIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.giftIncome')">{{ formatAmount(detail.totalGiftIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.paidDanmakuIncome')">{{ formatAmount(detail.totalPaidDanmakuIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.videoTicketIncome')">{{ formatAmount(detail.totalVideoCallTicketIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.videoBillingIncome')">{{ formatAmount(detail.totalVideoCallBillingIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.videoCallIncome')">{{ formatAmount(detail.totalVideoCallIncome) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.banStatus')">
            <el-tag v-if="detail.ban" type="danger">{{ t('common.banned') }}</el-tag>
            <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.banUntil')">{{ formatDate(detail.banApplyTime) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.banReason')">{{ detail.banReason || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.registeredAt')">{{ formatDate(detail.registeredAt) }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.anchorList.profileUpdatedAt')">{{ formatDate(detail.createdAt) }}</el-descriptions-item>
        </el-descriptions>
        <el-empty v-else-if="!loading" :description="t('pages.anchorList.detailNotFound')"/>
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
import type {AnchorListItem} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<AnchorListItem | null>(null)

const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_PRIVATE = 3
const LIVE_ROOM_PRIVATE_INVITE_ALL = 1
const LIVE_ROOM_PRIVATE_INVITE_VIP = 2
const LIVE_ROOM_PRIVATE_INVITE_REJECT = 3
const USER_TYPE_ANCHOR = 1
const USER_TYPE_SENIOR_ANCHOR = 7

const anchorId = computed(() => {
  const value = route.query.id
  if (Array.isArray(value)) {
    return String(value[0] ?? '')
  }
  if (value == null || value === '') {
    return ''
  }
  return String(value)
})

const pageTitle = computed(() => {
  if (detail.value?.nickname) {
    return t('pages.anchorList.detailTitleWithName', {name: detail.value.nickname})
  }
  if (anchorId.value) {
    return t('pages.anchorList.detailTitleWithId', {id: anchorId.value})
  }
  return t('pages.anchorList.detailTitle')
})

const isPrivateRoom = (category?: number) => category === LIVE_ROOM_CATEGORY_PRIVATE

const anchorTypeLabel = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return t('pages.anchorList.anchorTypeSenior')
  if (userType === USER_TYPE_ANCHOR) return t('pages.anchorList.anchorTypeNormal')
  return '-'
}

const anchorTypeTagType = (userType?: number) => {
  if (userType === USER_TYPE_SENIOR_ANCHOR) return 'warning'
  if (userType === USER_TYPE_ANCHOR) return 'success'
  return 'info'
}

const privateInviteLabel = (type?: number) => {
  if (type === LIVE_ROOM_PRIVATE_INVITE_VIP) return t('pages.anchorList.privateInviteVipOnly')
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return t('pages.anchorList.privateInviteRejectAll')
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return t('pages.anchorList.privateInviteAcceptAll')
  return '-'
}

const privateInviteTagType = (type?: number) => {
  if (type === LIVE_ROOM_PRIVATE_INVITE_VIP) return 'warning'
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return 'danger'
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return 'success'
  return 'info'
}

const categoryLabel = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_HOT) return t('pages.anchorList.categoryHot')
  if (category === LIVE_ROOM_CATEGORY_GAME) return t('pages.anchorList.categoryGame')
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return t('pages.anchorList.categoryPrivate')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_PRIVATE) return 'warning'
  if (category === LIVE_ROOM_CATEGORY_GAME) return 'success'
  if (category === LIVE_ROOM_CATEGORY_HOT) return 'danger'
  return 'info'
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return '-'
  }
}

const fetchDetail = async () => {
  if (!anchorId.value) {
    detail.value = null
    return
  }
  loading.value = true
  try {
    const response = await accountApi.getAnchorList({
      pageIndex: 1,
      pageSize: 50,
      key: anchorId.value,
    })
    const list = response.data || []
    detail.value = list.find(item => String(item.id) === anchorId.value) || null
    if (!detail.value) {
      ElMessage.warning(t('pages.anchorList.detailNotFound'))
    }
  } catch (error) {
    console.error('Failed to load anchor detail:', error)
    detail.value = null
    ElMessage.error(t('pages.anchorList.detailFetchFailed'))
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
    return
  }
  router.push({name: 'AnchorListManagement'})
}

watch(anchorId, (_id, prev) => {
  if (prev !== undefined) {
    fetchDetail()
  }
})

onActivated(() => {
  fetchDetail()
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

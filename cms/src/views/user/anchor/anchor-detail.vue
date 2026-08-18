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
        <el-empty v-if="!loading && !detail" :description="t('pages.anchorList.detailNotFound')"/>
        <el-tabs v-else-if="detail" v-model="activeTab">
          <el-tab-pane :label="t('pages.anchorList.tabBasic')" name="basic">
            <el-descriptions v-if="detail.anchor" :column="1" border>
              <el-descriptions-item :label="t('common.userId')">{{ detail.anchor.id }}</el-descriptions-item>
              <el-descriptions-item :label="t('common.nickname')">{{ detail.anchor.nickname || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('common.avatar')">
                <el-image
                    v-if="detail.anchor.avatar"
                    :preview-src-list="[detail.anchor.avatar]"
                    :src="detail.anchor.avatar"
                    fit="cover"
                    hide-on-click-modal
                    preview-teleported
                    style="width:48px;height:48px;border-radius:50%"
                />
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('common.phone')">{{ detail.anchor.phone || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.guildId')">{{ detail.anchor.guildId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.anchorType')">
                <el-tag :type="anchorTypeTagType(detail.anchor.userType)">{{ anchorTypeLabel(detail.anchor.userType) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.loginIp')">{{ detail.anchor.ip || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banStatus')">
                <el-tag v-if="detail.anchor.ban" type="danger">{{ t('common.banned') }}</el-tag>
                <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banUntil')">{{ formatDate(detail.anchor.banApplyTime) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banReason')">{{ detail.anchor.banReason || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.registeredAt')">{{ formatDate(detail.anchor.registeredAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.profileUpdatedAt')">{{ formatDate(detail.anchor.createdAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabLiveRoom')" name="liveRoom">
            <el-descriptions v-if="detail.liveRoom" :column="1" border>
              <el-descriptions-item :label="t('pages.anchorList.liveRoom')">{{ detail.liveRoom.id || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.guildId')">{{ detail.liveRoom.guildId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomTitle')">{{ detail.liveRoom.title || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomCover')">
                <el-image
                    v-if="detail.liveRoom.cover"
                    :preview-src-list="[detail.liveRoom.cover]"
                    :src="detail.liveRoom.cover"
                    fit="cover"
                    hide-on-click-modal
                    preview-teleported
                    style="width:80px;height:80px;border-radius:4px"
                />
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomNotice')">{{ detail.liveRoom.notice || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.liveRecordId')">{{ detail.liveRoom.liveRecordId || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.heartTime')">{{ formatDate(detail.liveRoom.heartTime) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomType')">
                <el-tag :type="categoryTagType(detail.liveRoom.category)">{{ categoryLabel(detail.liveRoom.category) }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.privateInviteType')">
                <template v-if="detail.liveRoom.category === LIVE_ROOM_CATEGORY_HOT">
                  <el-tag :type="privateInviteTagType(detail.liveRoom.privateInviteType)">
                    {{ privateInviteLabel(detail.liveRoom.privateInviteType) }}
                  </el-tag>
                </template>
                <span v-else>-</span>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.ticketPrice')">
                {{ isPrivateRoom(detail.liveRoom.category) ? formatAmount(detail.liveRoom.ticket) : '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.billingPricePerMinute')">
                {{ isPrivateRoom(detail.liveRoom.category) ? formatAmount(detail.liveRoom.billing) : '-' }}
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.liveStatus')">
                <el-tag :type="detail.liveRoom.liveStatus === 1 ? 'success' : 'info'">
                  {{ detail.liveRoom.liveStatus === 1 ? t('common.live') : t('common.offline') }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.shelfStatus')">
                <el-tag v-if="detail.liveRoom.status === 1" type="success">{{ t('common.onShelf') }}</el-tag>
                <el-tag v-else type="info">{{ t('common.offShelf') }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banStatus')">
                <el-tag v-if="detail.liveRoom.ban" type="danger">{{ t('common.banned') }}</el-tag>
                <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banUntil')">{{ formatDate(detail.liveRoom.banApplyTime) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.banReason')">{{ detail.liveRoom.banReason || '-' }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomCreatedAt')">{{ formatDate(detail.liveRoom.createdAt) }}</el-descriptions-item>
              <el-descriptions-item :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(detail.liveRoom.updatedAt) }}</el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeUnsettled')" name="incomeUnsettled">
            <IncomePanel :data="detail.incomeUnsettled" :updated-at="detail.incomeUnsettled?.updatedAt"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeSettled')" name="incomeSettled">
            <IncomePanel
                :data="detail.incomeSettled"
                :settlement-salary="detail.incomeSettled?.settlementSalary"
                :settlement-share-amount="detail.incomeSettled?.settlementShareAmount"
                :updated-at="detail.incomeSettled?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeTotal')" name="incomeTotal">
            <IncomePanel
                :data="detail.incomeTotal"
                :settlement-salary="detail.incomeTotal?.settlementSalary"
                :settlement-share-amount="detail.incomeTotal?.settlementShareAmount"
                :updated-at="detail.incomeTotal?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabDailyEffectiveLive')" name="dailyEffectiveLive">
            <DailyLivePanel :active="activeTab === 'dailyEffectiveLive'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabLiveRecord')" name="liveRecord">
            <LiveRecordPanel :active="activeTab === 'liveRecord'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabSettlementLog')" name="settlementLog">
            <SettlementLogPanel :active="activeTab === 'settlementLog'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeArchive')" name="incomeArchive">
            <el-table v-if="detail.incomeArchives?.length" :data="detail.incomeArchives" style="width:100%">
              <el-table-column :label="t('pages.anchorList.archiveId')" min-width="180" prop="id"/>
              <el-table-column :label="t('pages.anchorList.guildId')" min-width="120" prop="guildId"/>
              <el-table-column :label="t('pages.anchorList.liveIncome')" min-width="110">
                <template #default="{ row }">{{ formatAmount(row.totalIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.giftIncome')" min-width="110">
                <template #default="{ row }">{{ formatAmount(row.totalGiftIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" min-width="120">
                <template #default="{ row }">{{ formatAmount(row.totalPaidDanmakuIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.privateRoomTicketIncome')" min-width="130">
                <template #default="{ row }">{{ formatAmount(row.totalPrivateRoomTicketIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.privateRoomWatchIncome')" min-width="130">
                <template #default="{ row }">{{ formatAmount(row.totalPrivateRoomWatchIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.videoCallIncome')" min-width="120">
                <template #default="{ row }">{{ formatAmount(row.totalVideoCallIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.totalLiveDuration')" min-width="120">
                <template #default="{ row }">{{ row.totalLiveDuration ?? '-' }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.settlementSalary')" min-width="110">
                <template #default="{ row }">{{ formatAmount(row.settlementSalary) }}</template>
              </el-table-column>
              <el-table-column :label="t('common.createdAt')" min-width="170">
                <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
              </el-table-column>
            </el-table>
            <el-empty v-else :description="t('pages.anchorList.noArchiveData')"/>
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
import IncomePanel from './anchor-detail-income-panel.vue'
import DailyLivePanel from './anchor-detail-daily-live-panel.vue'
import LiveRecordPanel from './anchor-detail-live-record-panel.vue'
import SettlementLogPanel from './anchor-detail-settlement-log-panel.vue'
import type {AnchorDetail} from '@/types/api'
import {formatAmount} from '@/utils/number-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<AnchorDetail | null>(null)
const activeTab = ref('basic')

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
  if (detail.value?.anchor?.nickname) {
    return t('pages.anchorList.detailTitleWithName', {name: detail.value.anchor.nickname})
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

const isAnchorDetailRoute = () => route.name === 'AnchorDetail'

const fetchDetail = async () => {
  if (!isAnchorDetailRoute()) {
    return
  }
  if (!anchorId.value) {
    detail.value = null
    return
  }
  loading.value = true
  try {
    detail.value = await accountApi.getAnchorDetail(anchorId.value)
    if (!detail.value?.anchor) {
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
  if (!isAnchorDetailRoute()) {
    return
  }
  if (prev !== undefined) {
    activeTab.value = 'basic'
    fetchDetail()
  }
})

onActivated(() => {
  if (!isAnchorDetailRoute()) {
    return
  }
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

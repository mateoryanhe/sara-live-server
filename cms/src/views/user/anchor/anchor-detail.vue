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
        <el-tabs v-else-if="detail" :key="anchorId" v-model="activeTab" class="anchor-detail-tabs">
          <el-tab-pane :label="t('pages.anchorList.tabBasic')" name="basic">
            <div class="detail-tab-content">
              <div class="detail-overview">
                <el-image
                    v-if="detail.anchor?.avatar"
                    :preview-src-list="[detail.anchor.avatar]"
                    :src="detail.anchor.avatar"
                    class="overview-avatar"
                    fit="cover"
                    hide-on-click-modal
                    preview-teleported
                />
                <div v-else class="overview-avatar overview-avatar-placeholder">{{ avatarPlaceholder }}</div>
                <div class="detail-overview-main">
                  <div class="detail-overview-title">{{ detail.anchor?.nickname || '-' }}</div>
                  <div class="detail-overview-id">{{ t('common.userId') }}：{{ detail.anchor?.id ?? '-' }}</div>
                  <div class="detail-overview-tags">
                    <el-tag :type="anchorTypeTagType(detail.anchor?.userType)">
                      {{ anchorTypeLabel(detail.anchor?.userType) }}
                    </el-tag>
                    <el-tag v-if="detail.anchor?.ban" type="danger">{{ t('common.banned') }}</el-tag>
                    <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
                  </div>
                </div>
              </div>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.profileSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                  <el-descriptions-item :label="t('common.phone')">{{ detail.anchor?.phone || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.guildId')">{{ detail.anchor?.guildId || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.loginIp')">{{ detail.anchor?.ip || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.registeredAt')">{{ formatDate(detail.anchor?.registeredAt) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.profileUpdatedAt')">{{ formatDate(detail.anchor?.createdAt) }}</el-descriptions-item>
                  <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                </el-descriptions>
              </section>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.accountStatusSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                  <el-descriptions-item :label="t('pages.anchorList.banStatus')">
                    <el-tag v-if="detail.anchor?.ban" type="danger">{{ t('common.banned') }}</el-tag>
                    <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.banUntil')">{{ formatDate(detail.anchor?.banApplyTime) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.banReason')">{{ detail.anchor?.banReason || '-' }}</el-descriptions-item>
                </el-descriptions>
              </section>
            </div>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabLiveRoom')" lazy name="liveRoom">
            <div v-if="detail.liveRoom" class="detail-tab-content">
              <div class="detail-overview">
                <el-image
                    v-if="detail.liveRoom.cover"
                    :preview-src-list="[detail.liveRoom.cover]"
                    :src="detail.liveRoom.cover"
                    class="overview-room-cover"
                    fit="cover"
                    hide-on-click-modal
                    preview-teleported
                />
                <div v-else class="overview-room-cover overview-cover-placeholder">-</div>
                <div class="detail-overview-main">
                  <div class="detail-overview-title">{{ detail.liveRoom.title || '-' }}</div>
                  <div class="detail-overview-id">{{ t('pages.anchorList.liveRoom') }}：{{ detail.liveRoom.id || '-' }}</div>
                  <div class="detail-overview-tags">
                    <el-tag :type="categoryTagType(detail.liveRoom.category)">{{ categoryLabel(detail.liveRoom.category) }}</el-tag>
                    <el-tag :type="detail.liveRoom.liveStatus === 1 ? 'success' : 'info'">
                      {{ detail.liveRoom.liveStatus === 1 ? t('common.live') : t('common.offline') }}
                    </el-tag>
                    <el-tag :type="detail.liveRoom.status === 1 ? 'success' : 'info'">
                      {{ detail.liveRoom.status === 1 ? t('common.onShelf') : t('common.offShelf') }}
                    </el-tag>
                  </div>
                </div>
              </div>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.liveRoomProfileSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                  <el-descriptions-item :label="t('pages.anchorList.liveRoom')">{{ detail.liveRoom.id || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.guildId')">{{ detail.liveRoom.guildId || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.roomTitle')">{{ detail.liveRoom.title || '-' }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.roomNotice')" :span="3">
                    <span class="long-text">{{ detail.liveRoom.notice || '-' }}</span>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.roomCreatedAt')">{{ formatDate(detail.liveRoom.createdAt) }}</el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.roomUpdatedAt')">{{ formatDate(detail.liveRoom.updatedAt) }}</el-descriptions-item>
                  <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                </el-descriptions>
              </section>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.liveRoomOperationSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
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
                  <el-descriptions-item :label="t('pages.anchorList.billingPricePerMinute')">
                    {{ formatAmount(detail.liveRoom.billing) }}
                  </el-descriptions-item>
                  <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                </el-descriptions>
              </section>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.salaryInfoSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
                  <el-descriptions-item :label="t('pages.anchorList.hasSalary')">
                    <el-tag :type="detail.liveRoom.hasSalary ? 'success' : 'info'">
                      {{ detail.liveRoom.hasSalary ? t('common.yes') : t('common.no') }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.salaryEffectiveStatus')">
                    <el-tag :type="detail.liveRoom.salaryEffective ? 'success' : 'info'">
                      {{ detail.liveRoom.salaryEffective ? t('pages.anchorList.salaryEffective') : t('pages.anchorList.salaryInactive') }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.salaryEffectiveStartTime')">
                    {{ formatDate(detail.liveRoom.salaryEffectiveStartTime) }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="t('pages.anchorList.salaryEffectiveEndTime')">
                    {{ formatDate(detail.liveRoom.salaryEffectiveEndTime) }}
                  </el-descriptions-item>
                  <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                  <el-descriptions-item class-name="detail-placeholder-cell" label-class-name="detail-placeholder-cell"/>
                </el-descriptions>
              </section>

              <section class="detail-section">
                <div class="detail-section-title">{{ t('pages.anchorList.roomStatusSection') }}</div>
                <el-descriptions :column="3" :label-width="detailLabelWidth" border class="detail-descriptions">
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
                  <el-descriptions-item :label="t('pages.anchorList.banReason')" :span="2">
                    <span class="long-text">{{ detail.liveRoom.banReason || '-' }}</span>
                  </el-descriptions-item>
                </el-descriptions>
              </section>
            </div>
            <el-empty v-else :description="t('pages.anchorList.noLiveRoomData')"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeUnsettled')" lazy name="incomeUnsettled">
            <IncomePanel :data="detail.incomeUnsettled" :updated-at="detail.incomeUnsettled?.updatedAt"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeSettled')" lazy name="incomeSettled">
            <IncomePanel
                :data="detail.incomeSettled"
                :settlement-salary="detail.incomeSettled?.settlementSalary"
                :settlement-share-amount="detail.incomeSettled?.settlementShareAmount"
                :settlement-share-amount-usd="detail.incomeSettled?.settlementShareAmountUsd"
                :updated-at="detail.incomeSettled?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeTotal')" lazy name="incomeTotal">
            <IncomePanel
                :data="detail.incomeTotal"
                :settlement-salary="detail.incomeTotal?.settlementSalary"
                :settlement-share-amount="detail.incomeTotal?.settlementShareAmount"
                :settlement-share-amount-usd="detail.incomeTotal?.settlementShareAmountUsd"
                :updated-at="detail.incomeTotal?.updatedAt"
            />
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabDailyEffectiveLive')" lazy name="dailyEffectiveLive">
            <DailyLivePanel :active="activeTab === 'dailyEffectiveLive'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabLiveRecord')" lazy name="liveRecord">
            <LiveRecordPanel :active="activeTab === 'liveRecord'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabSettlementLog')" lazy name="settlementLog">
            <SettlementLogPanel :active="activeTab === 'settlementLog'" :anchor-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane v-if="canViewShortVideo" :label="t('pages.anchorList.tabShortVideo')" lazy name="shortVideo">
            <ShortVideoPanel :active="activeTab === 'shortVideo'" :user-id="anchorId"/>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.anchorList.tabIncomeArchive')" lazy name="incomeArchive">
            <el-table v-if="detail.incomeArchives?.length" :data="detail.incomeArchives" style="width:100%">
              <el-table-column :label="t('pages.anchorList.archiveId')" min-width="180" prop="id"/>
              <el-table-column :label="t('pages.anchorList.guildId')" min-width="120" prop="guildId"/>
              <el-table-column :label="t('pages.anchorList.liveIncome')" min-width="110">
                <template #default="{ row }">{{ formatAmount(row.totalIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.socialIncomeTotal')" min-width="130">
                <template #default="{ row }">{{ formatAmount(row.totalSocialIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.giftIncome')" min-width="110">
                <template #default="{ row }">{{ formatAmount(row.totalGiftIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.paidDanmakuIncome')" min-width="120">
                <template #default="{ row }">{{ formatAmount(row.totalPaidDanmakuIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.videoCallIncome')" min-width="120">
                <template #default="{ row }">{{ formatAmount(row.totalVideoCallIncome) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.anchorList.totalLiveDuration')" min-width="120">
                <template #default="{ row }">{{ formatLiveDurationMinutes(row.totalLiveDuration, t) }}</template>
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
import ShortVideoPanel from './anchor-detail-short-video-panel.vue'
import type {AnchorDetail} from '@/types/api'
import {formatAmount} from '@/utils/number-format'
import {formatLiveDurationMinutes} from '@/utils/live-duration-format'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const {t} = useI18n()
const {can} = usePagePermission('AnchorDetail')
const canViewShortVideo = computed(() => can('shortVideo'))
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const detail = ref<AnchorDetail | null>(null)

const LIVE_ROOM_CATEGORY_HOT = 1
const LIVE_ROOM_CATEGORY_GAME = 2
const LIVE_ROOM_CATEGORY_ONE_TO_ONE = 4
const LIVE_ROOM_PRIVATE_INVITE_ALL = 1
const LIVE_ROOM_PRIVATE_INVITE_REJECT = 3
const USER_TYPE_ANCHOR = 1
const USER_TYPE_SENIOR_ANCHOR = 7

const ANCHOR_DETAIL_TAB_NAMES = new Set([
  'basic',
  'liveRoom',
  'incomeUnsettled',
  'incomeSettled',
  'incomeTotal',
  'dailyEffectiveLive',
  'liveRecord',
  'settlementLog',
  'shortVideo',
  'incomeArchive',
])

const resolveActiveTab = (tab: unknown) => {
  if (typeof tab === 'string' && ANCHOR_DETAIL_TAB_NAMES.has(tab)) {
    return tab
  }
  return 'basic'
}

const activeTab = ref(resolveActiveTab(route.query.tab))

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

const detailLabelWidth = 170
const avatarPlaceholder = computed(() => String(detail.value?.anchor?.nickname || '?').trim().slice(0, 1).toUpperCase())

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
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return t('pages.anchorList.privateInviteRejectAll')
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return t('pages.anchorList.privateInviteAcceptAll')
  return '-'
}

const privateInviteTagType = (type?: number) => {
  if (type === LIVE_ROOM_PRIVATE_INVITE_REJECT) return 'danger'
  if (type === LIVE_ROOM_PRIVATE_INVITE_ALL) return 'success'
  return 'info'
}

const categoryLabel = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_HOT) return t('pages.anchorList.categoryHot')
  if (category === LIVE_ROOM_CATEGORY_GAME) return t('pages.anchorList.categoryGame')
  if (category === LIVE_ROOM_CATEGORY_ONE_TO_ONE) return t('pages.anchorList.categoryOneToOne')
  return '-'
}

const categoryTagType = (category?: number) => {
  if (category === LIVE_ROOM_CATEGORY_ONE_TO_ONE) return 'warning'
  if (category === LIVE_ROOM_CATEGORY_GAME) return 'success'
  if (category === LIVE_ROOM_CATEGORY_HOT) return 'danger'
  return 'info'
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

watch(anchorId, () => {
  if (!isAnchorDetailRoute()) {
    return
  }
  activeTab.value = resolveActiveTab(route.query.tab)
  void fetchDetail()
}, {immediate: true})

watch(
  () => route.query.tab,
  (tab) => {
    if (!isAnchorDetailRoute()) {
      return
    }
    activeTab.value = resolveActiveTab(tab)
  },
)

onActivated(() => {
  if (!isAnchorDetailRoute()) {
    return
  }
  void fetchDetail()
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

.detail-tab-content {
  padding-top: 8px;
}

.detail-overview {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-bottom: 20px;
  padding: 18px 20px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  background: var(--el-fill-color-extra-light);
}

.overview-avatar,
.overview-room-cover {
  flex: 0 0 auto;
  width: 72px;
  height: 72px;
}

.overview-avatar {
  border-radius: 50%;
}

.overview-room-cover {
  border-radius: 6px;
}

.overview-avatar-placeholder,
.overview-cover-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-dark);
  color: var(--el-text-color-secondary);
  font-size: 24px;
}

.detail-overview-main {
  min-width: 0;
}

.detail-overview-title {
  overflow: hidden;
  color: var(--el-text-color-primary);
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-overview-id {
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.detail-overview-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.detail-section + .detail-section {
  margin-top: 20px;
}

.detail-section-title {
  margin-bottom: 10px;
  padding-left: 10px;
  border-left: 3px solid var(--el-color-primary);
  color: var(--el-text-color-primary);
  font-size: 15px;
  font-weight: 600;
  line-height: 20px;
}

.detail-descriptions :deep(.el-descriptions__table) {
  table-layout: fixed;
}

.detail-descriptions :deep(.el-descriptions__content) {
  overflow-wrap: anywhere;
}

.detail-descriptions :deep(.detail-placeholder-cell) {
  background: var(--el-fill-color-blank);
}

.long-text {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

@media (max-width: 1280px) {
  .detail-descriptions :deep(.el-descriptions__label) {
    width: 140px !important;
  }
}
</style>

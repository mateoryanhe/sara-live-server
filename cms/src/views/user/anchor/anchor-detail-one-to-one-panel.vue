<template>
  <div class="detail-tab-content">
    <section class="detail-section">
      <div class="detail-section-title">{{ t('pages.anchorList.oneToOneRoomStatusSection') }}</div>
      <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
        <el-descriptions-item :label="t('pages.anchorList.oneToOneRoomEnabled')">
          <el-tag :type="enabled ? 'success' : 'info'">
            {{ enabled ? t('pages.anchorList.oneToOneRoomEnabledYes') : t('pages.anchorList.oneToOneRoomEnabledNo') }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </section>

    <template v-if="room">
      <div class="detail-overview">
        <el-image
            v-if="room.cover"
            :preview-src-list="[room.cover]"
            :src="room.cover"
            class="overview-room-cover"
            fit="cover"
            hide-on-click-modal
            preview-teleported
        />
        <div v-else class="overview-room-cover overview-cover-placeholder">-</div>
        <div class="detail-overview-main">
          <div class="detail-overview-title">{{ room.title || t('pages.anchorList.oneToOneRoomUntitled') }}</div>
          <div class="detail-overview-id">{{ t('common.userId') }}：{{ room.userId }}</div>
          <div class="detail-overview-tags">
            <el-tag :type="room.status === 1 ? 'success' : 'info'">
              {{ room.status === 1 ? t('pages.anchorList.oneToOneShelfOn') : t('pages.anchorList.oneToOneShelfOff') }}
            </el-tag>
            <el-tag :type="room.liveRoomStatus === 1 ? 'success' : 'info'">
              {{ t('pages.oneToOneRoomList.liveRoomStatus') }}：
              {{ room.liveRoomStatus === 1 ? t('common.onShelf') : t('common.offShelf') }}
            </el-tag>
          </div>
        </div>
      </div>

      <section class="detail-section">
        <div class="detail-section-title">{{ t('pages.anchorList.oneToOneRoomDetailSection') }}</div>
        <el-descriptions :column="2" :label-width="detailLabelWidth" border class="detail-descriptions">
          <el-descriptions-item :label="t('pages.oneToOneRoomList.oneToOneTitle')">{{ room.title || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.oneToOneRoomList.billing')">
            <span class="money-amount">{{ formatWalletBalance(room.billing) }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('pages.oneToOneRoomList.tagId')">{{ room.tagId || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('pages.guildList.guildId')">{{ room.guildId || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="t('common.updatedAt')">{{ formatDate(room.updatedAt) }}</el-descriptions-item>
        </el-descriptions>
      </section>
    </template>

    <el-empty v-else :description="t('pages.anchorList.oneToOneRoomNotOpenedHint')"/>
  </div>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import type {OneToOneRoomItem} from '@/types/api'
import {formatWalletBalance} from '@/utils/number-format'
import {formatServerDateTime as formatDate} from '@/utils/server-datetime'

const props = defineProps<{
  oneToOneRoom?: OneToOneRoomItem | null
}>()

const {t} = useI18n()
const detailLabelWidth = 170
const room = computed(() => props.oneToOneRoom ?? null)
const enabled = computed(() => room.value != null)
</script>

<style scoped>
.detail-tab-content {
  padding-top: 8px;
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

.detail-overview {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
  align-items: flex-start;
}

.overview-room-cover {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  flex-shrink: 0;
}

.overview-cover-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
}

.detail-overview-main {
  flex: 1;
  min-width: 0;
}

.detail-overview-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 6px;
}

.detail-overview-id {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  margin-bottom: 10px;
}

.detail-overview-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.detail-descriptions :deep(.el-descriptions__table) {
  table-layout: fixed;
}
</style>

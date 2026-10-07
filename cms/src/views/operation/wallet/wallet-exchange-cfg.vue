<template>
  <div class="page-container">
    <h2 class="page-title">{{ t('menu.WalletExchangeCfgManagement') }}</h2>

    <el-card class="config-card">
      <el-tabs v-model="activeTab" class="config-tabs">
        <el-tab-pane
            v-if="canViewWallet"
            :label="t('pages.effectiveLiveCfg.walletSectionTitle')"
            name="wallet"
        >
          <div v-loading="walletLoading" class="tab-content">
            <el-alert
                :closable="false"
                class="tip-alert"
                show-icon
                :title="t('pages.walletExchangeCfg.tipTitle')"
                type="info"
            >
              <p>{{ t('pages.walletExchangeCfg.tipLine1') }}</p>
              <p>{{ t('pages.walletExchangeCfg.tipLine2') }}</p>
            </el-alert>

            <el-form :model="walletForm" class="cfg-form" label-width="190px">
              <el-form-item :label="t('pages.walletExchangeCfg.goldToDiamondRate')">
                <el-input-number
                    v-model="walletForm.goldToDiamondRate"
                    :min="1"
                    :step="1"
                    controls-position="right"
                />
              </el-form-item>

              <el-form-item :label="t('pages.walletExchangeCfg.usdToGoldRate')">
                <el-input-number
                    v-model="walletForm.usdToGoldRate"
                    :min="1"
                    :step="1"
                    controls-position="right"
                />
              </el-form-item>

              <el-form-item :label="t('pages.walletExchangeCfg.exchangeFeePercent')">
                <el-input-number
                    v-model="walletForm.exchangeFeePercent"
                    :min="0"
                    :precision="2"
                    :step="0.1"
                    controls-position="right"
                />
                <span class="field-tip">{{ t('pages.walletExchangeCfg.exchangeFeeTip') }}</span>
              </el-form-item>

              <el-form-item v-if="walletMeta.updatedAt" :label="t('pages.walletExchangeCfg.lastUpdated')">
                <span>{{ walletMeta.updatedAt }}</span>
              </el-form-item>

              <el-form-item>
                <el-button v-if="canSaveWallet" type="primary" @click="handleSaveWallet">
                  {{ t('common.saveConfig') }}
                </el-button>
                <el-button @click="fetchWalletCfg">{{ t('common.refresh') }}</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane
            v-if="canViewEffectiveLive"
            :label="t('pages.effectiveLiveCfg.sectionTitle')"
            name="effectiveLive"
        >
          <div v-loading="effectiveLiveLoading" class="tab-content">
            <el-alert
                :closable="false"
                class="tip-alert"
                show-icon
                :title="t('pages.effectiveLiveCfg.tipTitle')"
                type="info"
            >
              <p>{{ t('pages.effectiveLiveCfg.tipLine1') }}</p>
              <p>{{ t('pages.effectiveLiveCfg.tipLine2') }}</p>
              <p>{{ t('pages.effectiveLiveCfg.tipLine3') }}</p>
            </el-alert>

            <el-form :model="effectiveLiveForm" class="cfg-form" label-width="190px">
              <el-form-item :label="t('pages.effectiveLiveCfg.minSessionMinutes')">
                <el-input-number
                    v-model="effectiveLiveForm.minSessionMinutes"
                    :max="1440"
                    :min="1"
                    :step="1"
                    controls-position="right"
                />
                <span class="field-unit">{{ t('pages.effectiveLiveCfg.minutes') }}</span>
              </el-form-item>

              <el-form-item :label="t('pages.effectiveLiveCfg.dailyAccumulatedMinutes')">
                <el-input-number
                    v-model="effectiveLiveForm.dailyAccumulatedMinutes"
                    :max="1440"
                    :min="0"
                    :step="1"
                    controls-position="right"
                />
                <span class="field-unit">{{ t('pages.effectiveLiveCfg.minutes') }}</span>
              </el-form-item>

              <el-form-item v-if="effectiveLiveMeta.updatedAt" :label="t('pages.effectiveLiveCfg.lastUpdated')">
                <span>{{ effectiveLiveMeta.updatedAt }}</span>
              </el-form-item>

              <el-form-item>
                <el-button v-if="canSaveEffectiveLive" type="primary" @click="handleSaveEffectiveLive">
                  {{ t('common.saveConfig') }}
                </el-button>
                <el-button @click="fetchEffectiveLiveCfg">{{ t('common.refresh') }}</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane
            v-if="canViewPlatformAnchorSettlement"
            :label="t('pages.platformAnchorSettlementCfg.tabTitle')"
            name="platformAnchorSettlement"
        >
          <div v-loading="platformAnchorSettlementLoading" class="tab-content">
            <el-alert
                :closable="false"
                class="tip-alert"
                show-icon
                :title="t('pages.platformAnchorSettlementCfg.tipTitle')"
                type="info"
            >
              <p>{{ t('pages.platformAnchorSettlementCfg.tipLine1') }}</p>
              <p>{{ t('pages.platformAnchorSettlementCfg.tipLine2') }}</p>
            </el-alert>

            <el-form :model="platformAnchorSettlementForm" class="cfg-form" label-width="190px">
              <el-form-item :label="t('pages.platformAnchorSettlementCfg.minimumSettlementUsd')">
                <el-input-number
                    v-model="platformAnchorSettlementForm.minimumSettlementUsd"
                    :max="1000000"
                    :min="0.01"
                    :precision="2"
                    :step="0.5"
                    controls-position="right"
                />
                <span class="field-unit">{{ t('pages.platformAnchorSettlementCfg.usd') }}</span>
              </el-form-item>

              <el-form-item
                  v-if="platformAnchorSettlementMeta.updatedAt"
                  :label="t('pages.platformAnchorSettlementCfg.lastUpdated')"
              >
                <span>{{ platformAnchorSettlementMeta.updatedAt }}</span>
              </el-form-item>

              <el-form-item>
                <el-button
                    v-if="canSavePlatformAnchorSettlement"
                    type="primary"
                    @click="handleSavePlatformAnchorSettlement"
                >
                  {{ t('common.saveConfig') }}
                </el-button>
                <el-button @click="fetchPlatformAnchorSettlementCfg">{{ t('common.refresh') }}</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {walletApi} from '@/api/modules/wallet'
import {effectiveLiveCfgApi} from '@/api/modules/effective-live-cfg'
import {platformAnchorSettlementCfgApi} from '@/api/modules/platform-anchor-settlement-cfg'
import {hasButtonPermission} from '@/utils/permission'
import type {EffectiveLiveCfg, PlatformAnchorSettlementCfg, WalletExchangeCfg} from '@/types/api'

const {t} = useI18n()
const pageName = 'WalletExchangeCfgManagement'

const canViewWallet = computed(() => hasButtonPermission(pageName, 'view'))
const canSaveWallet = computed(() => hasButtonPermission(pageName, 'save'))
const canViewEffectiveLive = computed(() => hasButtonPermission(pageName, 'viewEffectiveLive'))
const canSaveEffectiveLive = computed(() => hasButtonPermission(pageName, 'saveEffectiveLive'))
const canViewPlatformAnchorSettlement = computed(() => hasButtonPermission(pageName, 'viewPlatformAnchorSettlement'))
const canSavePlatformAnchorSettlement = computed(() => hasButtonPermission(pageName, 'savePlatformAnchorSettlement'))

const activeTab = ref('wallet')
const walletLoading = ref(false)
const effectiveLiveLoading = ref(false)
const platformAnchorSettlementLoading = ref(false)

const walletForm = reactive({
  id: '0',
  goldToDiamondRate: 100,
  usdToGoldRate: 100,
  exchangeFeePercent: 3,
})
const walletMeta = reactive({createdAt: '', updatedAt: ''})

const effectiveLiveForm = reactive({
  id: '0',
  minSessionMinutes: 30,
  dailyAccumulatedMinutes: 0,
})
const effectiveLiveMeta = reactive({createdAt: '', updatedAt: ''})

const platformAnchorSettlementForm = reactive({
  id: '0',
  minimumSettlementUsd: 5,
})
const platformAnchorSettlementMeta = reactive({createdAt: '', updatedAt: ''})

const applyWalletCfg = (cfg: WalletExchangeCfg | null | undefined) => {
  walletForm.id = cfg?.id || '0'
  walletForm.goldToDiamondRate = cfg?.goldToDiamondRate || 100
  walletForm.usdToGoldRate = cfg?.usdToGoldRate || 100
  walletForm.exchangeFeePercent = cfg?.exchangeFeePercent ?? 3
  walletMeta.createdAt = cfg?.createdAt || ''
  walletMeta.updatedAt = cfg?.updatedAt || ''
}

const applyEffectiveLiveCfg = (cfg: EffectiveLiveCfg | null | undefined) => {
  effectiveLiveForm.id = cfg?.id || '0'
  effectiveLiveForm.minSessionMinutes = cfg?.minSessionMinutes || 30
  effectiveLiveForm.dailyAccumulatedMinutes = cfg?.dailyAccumulatedMinutes ?? 0
  effectiveLiveMeta.createdAt = cfg?.createdAt || ''
  effectiveLiveMeta.updatedAt = cfg?.updatedAt || ''
}

const applyPlatformAnchorSettlementCfg = (cfg: PlatformAnchorSettlementCfg | null | undefined) => {
  platformAnchorSettlementForm.id = cfg?.id || '0'
  platformAnchorSettlementForm.minimumSettlementUsd = cfg?.minimumSettlementUsd || 5
  platformAnchorSettlementMeta.createdAt = cfg?.createdAt || ''
  platformAnchorSettlementMeta.updatedAt = cfg?.updatedAt || ''
}

const fetchWalletCfg = async () => {
  walletLoading.value = true
  try {
    applyWalletCfg((await walletApi.getWalletExchangeCfg()).cfg)
  } catch (error) {
    console.error('fetch wallet exchange config failed:', error)
    ElMessage.error(t('pages.walletExchangeCfg.fetchFailed'))
  } finally {
    walletLoading.value = false
  }
}

const fetchEffectiveLiveCfg = async () => {
  effectiveLiveLoading.value = true
  try {
    applyEffectiveLiveCfg((await effectiveLiveCfgApi.getEffectiveLiveCfg()).cfg)
  } catch (error) {
    console.error('fetch effective live config failed:', error)
    ElMessage.error(t('pages.effectiveLiveCfg.fetchFailed'))
  } finally {
    effectiveLiveLoading.value = false
  }
}

const fetchPlatformAnchorSettlementCfg = async () => {
  platformAnchorSettlementLoading.value = true
  try {
    applyPlatformAnchorSettlementCfg((await platformAnchorSettlementCfgApi.getPlatformAnchorSettlementCfg()).cfg)
  } catch (error) {
    console.error('fetch platform anchor settlement config failed:', error)
    ElMessage.error(t('pages.platformAnchorSettlementCfg.fetchFailed'))
  } finally {
    platformAnchorSettlementLoading.value = false
  }
}

const handleSaveWallet = async () => {
  if (walletForm.goldToDiamondRate <= 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.rateMustPositive'))
    return
  }
  if (walletForm.usdToGoldRate <= 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.usdRateMustPositive'))
    return
  }
  if (walletForm.exchangeFeePercent < 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.feeCannotNegative'))
    return
  }

  walletLoading.value = true
  try {
    const response = await walletApi.saveWalletExchangeCfg({
      id: walletForm.id === '0' ? 0 : Number(walletForm.id),
      goldToDiamondRate: walletForm.goldToDiamondRate,
      usdToGoldRate: walletForm.usdToGoldRate,
      exchangeFeePercent: walletForm.exchangeFeePercent,
    })
    if (!response?.success) {
      ElMessage.error(t('pages.walletExchangeCfg.saveFailed'))
      return
    }
    ElMessage.success(t('common.saveConfig'))
    if (response.id) {
      walletForm.id = response.id
    }
    await fetchWalletCfg()
  } catch (error) {
    console.error('save wallet exchange config failed:', error)
    ElMessage.error(t('pages.walletExchangeCfg.saveFailed'))
  } finally {
    walletLoading.value = false
  }
}

const handleSaveEffectiveLive = async () => {
  if (effectiveLiveForm.minSessionMinutes < 1 || effectiveLiveForm.minSessionMinutes > 1440) {
    ElMessage.warning(t('pages.effectiveLiveCfg.rangeWarning'))
    return
  }
  if (effectiveLiveForm.dailyAccumulatedMinutes < 0 || effectiveLiveForm.dailyAccumulatedMinutes > 1440) {
    ElMessage.warning(t('pages.effectiveLiveCfg.dailyRangeWarning'))
    return
  }

  effectiveLiveLoading.value = true
  try {
    const response = await effectiveLiveCfgApi.saveEffectiveLiveCfg({
      id: effectiveLiveForm.id === '0' ? 0 : Number(effectiveLiveForm.id),
      minSessionMinutes: effectiveLiveForm.minSessionMinutes,
      dailyAccumulatedMinutes: effectiveLiveForm.dailyAccumulatedMinutes,
    })
    if (!response?.success) {
      ElMessage.error(t('pages.effectiveLiveCfg.saveFailed'))
      return
    }
    ElMessage.success(t('pages.effectiveLiveCfg.saveSuccess'))
    if (response.id) {
      effectiveLiveForm.id = response.id
    }
    await fetchEffectiveLiveCfg()
  } catch (error) {
    console.error('save effective live config failed:', error)
    ElMessage.error(t('pages.effectiveLiveCfg.saveFailed'))
  } finally {
    effectiveLiveLoading.value = false
  }
}

const handleSavePlatformAnchorSettlement = async () => {
  const minimumSettlementUsd = Number(platformAnchorSettlementForm.minimumSettlementUsd)
  if (!Number.isFinite(minimumSettlementUsd) || minimumSettlementUsd <= 0 || minimumSettlementUsd > 1000000) {
    ElMessage.warning(t('pages.platformAnchorSettlementCfg.rangeWarning'))
    return
  }

  platformAnchorSettlementLoading.value = true
  try {
    const response = await platformAnchorSettlementCfgApi.savePlatformAnchorSettlementCfg({
      id: platformAnchorSettlementForm.id === '0' ? 0 : Number(platformAnchorSettlementForm.id),
      minimumSettlementUsd,
    })
    if (!response?.success) {
      ElMessage.error(t('pages.platformAnchorSettlementCfg.saveFailed'))
      return
    }
    ElMessage.success(t('pages.platformAnchorSettlementCfg.saveSuccess'))
    if (response.id) {
      platformAnchorSettlementForm.id = response.id
    }
    await fetchPlatformAnchorSettlementCfg()
  } catch (error) {
    console.error('save platform anchor settlement config failed:', error)
    ElMessage.error(t('pages.platformAnchorSettlementCfg.saveFailed'))
  } finally {
    platformAnchorSettlementLoading.value = false
  }
}

// 页面由布局 keep-alive 缓存，每次打开都重新读取三个配置选项卡。
onActivated(() => {
  if (!canViewWallet.value) {
    if (canViewEffectiveLive.value) {
      activeTab.value = 'effectiveLive'
    } else if (canViewPlatformAnchorSettlement.value) {
      activeTab.value = 'platformAnchorSettlement'
    }
  }
  if (canViewWallet.value) {
    void fetchWalletCfg()
  }
  if (canViewEffectiveLive.value) {
    void fetchEffectiveLiveCfg()
  }
  if (canViewPlatformAnchorSettlement.value) {
    void fetchPlatformAnchorSettlementCfg()
  }
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.page-title {
  margin: 0 0 16px;
  font-size: 20px;
  font-weight: 600;
}

.config-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
}

.tab-content {
  min-height: 300px;
  padding: 0 4px;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 4px 0;
}

.cfg-form {
  max-width: 760px;
}

.field-tip,
.field-unit {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.WalletExchangeCfgManagement') }}</span>
        </div>
      </template>

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

      <el-form :model="formData" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.walletExchangeCfg.goldToDiamondRate')">
          <el-input-number
              v-model="formData.goldToDiamondRate"
              :min="1"
              :step="1"
              controls-position="right"
          />
        </el-form-item>

        <el-form-item :label="t('pages.walletExchangeCfg.usdToGoldRate')">
          <el-input-number
              v-model="formData.usdToGoldRate"
              :min="1"
              :step="1"
              controls-position="right"
          />
        </el-form-item>

        <el-form-item :label="t('pages.walletExchangeCfg.exchangeFeePercent')">
          <el-input-number
              v-model="formData.exchangeFeePercent"
              :min="0"
              :precision="2"
              :step="0.1"
              controls-position="right"
          />
          <span class="field-tip">{{ t('pages.walletExchangeCfg.exchangeFeeTip') }}</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.walletExchangeCfg.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {walletApi} from '@/api/modules/wallet'
import type {WalletExchangeCfg} from '@/types/api'

const {t} = useI18n()

const loading = ref(false)

const formData = reactive({
  id: '0',
  goldToDiamondRate: 100,
  usdToGoldRate: 100,
  exchangeFeePercent: 3,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: WalletExchangeCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.goldToDiamondRate = 100
    formData.usdToGoldRate = 100
    formData.exchangeFeePercent = 3
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.goldToDiamondRate = cfg.goldToDiamondRate || 100
  formData.usdToGoldRate = cfg.usdToGoldRate || 100
  formData.exchangeFeePercent = cfg.exchangeFeePercent ?? 3
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await walletApi.getWalletExchangeCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch wallet exchange config failed:', error)
    ElMessage.error(t('pages.walletExchangeCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (formData.goldToDiamondRate <= 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.rateMustPositive'))
    return
  }
  if (formData.usdToGoldRate <= 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.usdRateMustPositive'))
    return
  }
  if (formData.exchangeFeePercent < 0) {
    ElMessage.warning(t('pages.walletExchangeCfg.feeCannotNegative'))
    return
  }

  loading.value = true
  try {
    const response = await walletApi.saveWalletExchangeCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      goldToDiamondRate: formData.goldToDiamondRate,
      usdToGoldRate: formData.usdToGoldRate,
      exchangeFeePercent: formData.exchangeFeePercent,
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.walletExchangeCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save wallet exchange config failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCfg()
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

.field-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>

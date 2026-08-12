<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>金币兑换钻石配置</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>配置 1 金币兑换钻石数；App 端手动兑换时，若手续费比例大于 0 则额外扣金币。</p>
        <p>送礼、付费弹幕等业务内自动兑换始终免手续费，仅按兑换比例扣金币。</p>
      </el-alert>

      <el-form :model="formData" class="cfg-form" label-width="180px">
        <el-form-item label="1金币兑换钻石数">
          <el-input-number
              v-model="formData.goldToDiamondRate"
              :min="1"
              :step="1"
              controls-position="right"
          />
        </el-form-item>

        <el-form-item label="App兑换手续费(%)">
          <el-input-number
              v-model="formData.exchangeFeePercent"
              :min="0"
              :precision="2"
              :step="0.1"
              controls-position="right"
          />
          <span class="field-tip">设为 0 表示 App 手动兑换不收手续费</span>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" label="最近更新">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">保存配置</el-button>
          <el-button @click="fetchCfg">刷新</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {walletApi} from '@/api/modules/wallet'
import type {WalletExchangeCfg} from '@/types/api'

const loading = ref(false)

const formData = reactive({
  id: '0',
  goldToDiamondRate: 100,
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
    formData.exchangeFeePercent = 3
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.goldToDiamondRate = cfg.goldToDiamondRate || 100
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
    console.error('获取金币兑换配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (formData.goldToDiamondRate <= 0) {
    ElMessage.warning('兑换比例必须大于 0')
    return
  }
  if (formData.exchangeFeePercent < 0) {
    ElMessage.warning('手续费不能为负数')
    return
  }

  loading.value = true
  try {
    const response = await walletApi.saveWalletExchangeCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      goldToDiamondRate: formData.goldToDiamondRate,
      exchangeFeePercent: formData.exchangeFeePercent,
    })
    if (response?.success) {
      ElMessage.success('保存成功')
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error('保存失败')
    }
  } catch (error) {
    console.error('保存金币兑换配置失败:', error)
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

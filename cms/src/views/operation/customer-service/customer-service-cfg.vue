<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CustomerServiceCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.customerServiceCfg.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.customerServiceCfg.tipLine1') }}</p>
        <p>{{ t('pages.customerServiceCfg.tipLine2') }}</p>
      </el-alert>

      <el-form :model="formData" class="cfg-form" label-width="160px">
        <el-form-item label="Telegram">
          <el-input
              v-model="formData.telegramUrl"
              clearable
              :placeholder="t('pages.customerServiceCfg.telegramPlaceholder')"
          />
        </el-form-item>

        <el-form-item label="Facebook">
          <el-input
              v-model="formData.facebookUrl"
              clearable
              :placeholder="t('pages.customerServiceCfg.facebookPlaceholder')"
          />
        </el-form-item>

        <el-form-item label="WhatsApp">
          <el-input
              v-model="formData.whatsappUrl"
              clearable
              :placeholder="t('pages.customerServiceCfg.whatsappPlaceholder')"
          />
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.customerServiceCfg.lastUpdated')">
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
import {useI18n} from 'vue-i18n'
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {customerServiceApi} from '@/api/modules/customer-service'
import type {CustomerServiceCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)

const formData = reactive({
  id: '0',
  telegramUrl: '',
  facebookUrl: '',
  whatsappUrl: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: CustomerServiceCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.telegramUrl = ''
    formData.facebookUrl = ''
    formData.whatsappUrl = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.telegramUrl = cfg.telegramUrl || ''
  formData.facebookUrl = cfg.facebookUrl || ''
  formData.whatsappUrl = cfg.whatsappUrl || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await customerServiceApi.getCustomerServiceCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch customer service config failed:', error)
    ElMessage.error(t('pages.customerServiceCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  loading.value = true
  try {
    const response = await customerServiceApi.saveCustomerServiceCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      telegramUrl: formData.telegramUrl.trim(),
      facebookUrl: formData.facebookUrl.trim(),
      whatsappUrl: formData.whatsappUrl.trim(),
    })
    if (response?.success) {
      ElMessage.success(t('common.saveConfig'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.customerServiceCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save customer service config failed:', error)
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
</style>

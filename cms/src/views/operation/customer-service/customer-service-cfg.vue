<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>客服联系配置</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>配置 App 端展示的客服联系方式。App 通过 POST /customerService/cfg 获取（无需登录）。</p>
        <p>联系方式可为完整链接或任意文本；留空表示不展示对应入口。</p>
      </el-alert>

      <el-form :model="formData" class="cfg-form" label-width="160px">
        <el-form-item label="Telegram">
          <el-input
              v-model="formData.telegramUrl"
              clearable
              placeholder="如 https://t.me/your_support"
          />
        </el-form-item>

        <el-form-item label="Facebook">
          <el-input
              v-model="formData.facebookUrl"
              clearable
              placeholder="如 https://facebook.com/your_page"
          />
        </el-form-item>

        <el-form-item label="WhatsApp">
          <el-input
              v-model="formData.whatsappUrl"
              clearable
              placeholder="如 https://wa.me/1234567890"
          />
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
import {customerServiceApi} from '@/api/modules/customer-service'
import type {CustomerServiceCfg} from '@/types/api'

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
    console.error('获取客服联系配置失败:', error)
    ElMessage.error('获取配置失败')
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
      ElMessage.success('保存成功')
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error('保存失败')
    }
  } catch (error) {
    console.error('保存客服联系配置失败:', error)
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

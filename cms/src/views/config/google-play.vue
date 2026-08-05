<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>Google Play 充值配置</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>配置 Google Play RTDN 回调与验单参数。Pub/Sub Push 地址：<code>POST /webhook/googlePlay/rtdn</code></p>
        <p>App 创建订单后，发起 Google 支付时需将返回的 <code>obfuscatedAccountId</code> 传入 BillingFlowParams.setObfuscatedAccountId。</p>
        <p>充值档位请在「运营 → 充值配置」中维护 Google 类型商品 SKU（productId）。发货成功后会自动 consume 消耗型商品。</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="180px">
        <el-form-item label="服务账号 JSON" prop="serviceAccountJson">
          <el-input
              v-model="formData.serviceAccountJson"
              :rows="8"
              placeholder="粘贴 Google Cloud 服务账号 JSON 全文"
              type="textarea"
          />
          <span class="form-tip">需具备 Google Play Android Developer API 权限，并在 Play Console 关联该服务账号</span>
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
import {googlePlayApi} from '@/api/modules/google-play'
import type {GooglePlayCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  serviceAccountJson: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const validateServiceAccountJson = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const json = value?.trim()
  if (!json) {
    callback(new Error('请填写服务账号 JSON'))
    return
  }
  try {
    JSON.parse(json)
  } catch {
    callback(new Error('服务账号 JSON 格式无效'))
    return
  }
  callback()
}

const formRules = reactive({
  serviceAccountJson: [{validator: validateServiceAccountJson, trigger: 'blur'}],
})

const applyCfg = (cfg: GooglePlayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.serviceAccountJson = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.serviceAccountJson = cfg.serviceAccountJson || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await googlePlayApi.getGooglePlayCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取 Google Play 配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await googlePlayApi.saveGooglePlayCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: true,
      serviceAccountJson: formData.serviceAccountJson.trim(),
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
    console.error('保存 Google Play 配置失败:', error)
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
  max-width: 860px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}

code {
  padding: 0 4px;
  background: #f4f4f5;
  border-radius: 4px;
}
</style>

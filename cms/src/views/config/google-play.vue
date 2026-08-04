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
        <el-form-item label="启用 RTDN 到账" prop="enabled">
          <el-switch v-model="formData.enabled"/>
        </el-form-item>

        <el-form-item label="Android 包名" prop="packageName">
          <el-input
              v-model="formData.packageName"
              clearable
              placeholder="如 com.example.app"
          />
        </el-form-item>

        <el-form-item label="服务账号 JSON" prop="serviceAccountJson">
          <el-input
              v-model="formData.serviceAccountJson"
              :rows="8"
              placeholder="粘贴 Google Cloud 服务账号 JSON 全文"
              type="textarea"
          />
          <span class="form-tip">需具备 Google Play Android Developer API 权限，并在 Play Console 关联该服务账号</span>
        </el-form-item>

        <el-form-item label="RTDN JWT Audience" prop="rtdnAudience">
          <el-input
              v-model="formData.rtdnAudience"
              clearable
              placeholder="如 https://your-domain.com/webhook/googlePlay/rtdn"
          />
          <span class="form-tip">Pub/Sub Push 订阅 JWT 校验 aud；留空则不校验</span>
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
  enabled: false,
  packageName: '',
  serviceAccountJson: '',
  rtdnAudience: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const validateOptionalUrl = (_: unknown, value: string, callback: (e?: Error) => void) => {
  const url = value?.trim()
  if (!url) {
    callback()
    return
  }
  if (!/^https?:\/\//i.test(url)) {
    callback(new Error('URL 需以 http:// 或 https:// 开头'))
    return
  }
  callback()
}

const validateServiceAccountJson = (_: unknown, value: string, callback: (e?: Error) => void) => {
  if (!formData.enabled) {
    callback()
    return
  }
  const json = value?.trim()
  if (!json) {
    callback(new Error('启用时请填写服务账号 JSON'))
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
  packageName: [
    {max: 128, message: '包名最长 128 字符', trigger: 'blur'},
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error('启用时请填写包名'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  serviceAccountJson: [{validator: validateServiceAccountJson, trigger: 'blur'}],
  rtdnAudience: [
    {max: 512, message: '长度不能超过 512', trigger: 'blur'},
    {validator: validateOptionalUrl, trigger: 'blur'},
  ],
})

const applyCfg = (cfg: GooglePlayCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.packageName = ''
    formData.serviceAccountJson = ''
    formData.rtdnAudience = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.packageName = cfg.packageName || ''
  formData.serviceAccountJson = cfg.serviceAccountJson || ''
  formData.rtdnAudience = cfg.rtdnAudience || ''
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
      enabled: formData.enabled,
      packageName: formData.packageName.trim(),
      serviceAccountJson: formData.serviceAccountJson.trim(),
      rtdnAudience: formData.rtdnAudience.trim(),
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

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>游戏平台接入配置</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>配置第三方游戏平台 API 凭证，用于后续拉取游戏列表、启动游戏等接口。</p>
        <p>厂家 URL 为平台 API 根地址；Token 对应请求头 <code>x-token</code>；SecretKey 用于平台回调验签（verify / balance / transfer）。</p>
        <p>IconUrl 为游戏封面 CDN 根地址，第三方游戏列表返回的 cover 为相对路径，展示时会自动拼接。</p>
        <p>文档参考：<code>https://admin.win12.best/#/layout/integrator/docs</code></p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="厂家 URL" prop="vendorUrl">
          <el-input v-model="formData.vendorUrl" clearable placeholder="如 https://gapi.win12.best"/>
        </el-form-item>

        <el-form-item label="IconUrl" prop="iconUrl">
          <el-input v-model="formData.iconUrl" clearable placeholder="游戏封面 CDN 根地址，如 https://cdn.example.com/images"/>
        </el-form-item>

        <el-form-item label="Token" prop="token">
          <el-input
              v-model="formData.token"
              clearable
              placeholder="平台接入 Token (x-token)"
              show-password
              type="password"
          />
        </el-form-item>

        <el-form-item label="SecretKey" prop="secretKey">
          <el-input
              v-model="formData.secretKey"
              clearable
              placeholder="平台 SecretKey"
              show-password
              type="password"
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
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {gamePlatformApi} from '@/api'
import type {GamePlatformCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  vendorUrl: 'https://gapi.win12.best',
  iconUrl: '',
  token: '',
  secretKey: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules: FormRules = {
  vendorUrl: [
    {required: true, message: '请填写厂家 URL', trigger: 'blur'},
    {max: 255, message: '厂家 URL 最长 255 字符', trigger: 'blur'},
  ],
  iconUrl: [{max: 512, message: 'IconUrl 最长 512 字符', trigger: 'blur'}],
  token: [
    {required: true, message: '请填写 Token', trigger: 'blur'},
    {max: 512, message: 'Token 最长 512 字符', trigger: 'blur'},
  ],
  secretKey: [
    {required: true, message: '请填写 SecretKey', trigger: 'blur'},
    {max: 255, message: 'SecretKey 最长 255 字符', trigger: 'blur'},
  ],
}

const applyCfg = (cfg: GamePlatformCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.vendorUrl = 'https://gapi.win12.best'
    formData.iconUrl = ''
    formData.token = ''
    formData.secretKey = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.vendorUrl = cfg.vendorUrl || 'https://gapi.win12.best'
  formData.iconUrl = cfg.iconUrl || ''
  formData.token = cfg.token || ''
  formData.secretKey = cfg.secretKey || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await gamePlatformApi.getGamePlatformCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取游戏平台配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await gamePlatformApi.saveGamePlatformCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      vendorUrl: formData.vendorUrl.trim(),
      iconUrl: formData.iconUrl.trim().replace(/\/+$/, ''),
      token: formData.token.trim(),
      secretKey: formData.secretKey.trim(),
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
    console.error('保存游戏平台配置失败:', error)
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

code {
  padding: 0 4px;
  background: #f4f4f5;
  border-radius: 4px;
}
</style>

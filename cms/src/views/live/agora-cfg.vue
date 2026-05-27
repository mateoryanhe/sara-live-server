<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>声网配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="AppId" prop="appId">
          <el-input v-model="formData.appId" clearable placeholder="请输入声网 AppId"/>
        </el-form-item>

        <el-form-item label="AppCertificate" prop="appCertificate">
          <el-input
              v-model="formData.appCertificate"
              clearable
              placeholder="请输入声网 AppCertificate"
              show-password
              type="password"
          />
        </el-form-item>

        <el-form-item label="REST CustomerId" prop="restCustomerId">
          <el-input v-model="formData.restCustomerId" clearable placeholder="请输入声网 REST CustomerId"/>
          <span class="form-tip">可选，查询用户在线状态时需要配置</span>
        </el-form-item>

        <el-form-item label="REST CustomerSecret" prop="restCustomerSecret">
          <el-input
              v-model="formData.restCustomerSecret"
              clearable
              placeholder="请输入声网 REST CustomerSecret"
              show-password
              type="password"
          />
          <span class="form-tip">可选，查询用户在线状态时需要配置</span>
        </el-form-item>

        <el-form-item label="Token有效期" prop="tokenExpireSeconds">
          <el-input-number
              v-model="formData.tokenExpireSeconds"
              :min="60"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">秒，默认 6 小时(21600 秒)，最小 60 秒</span>
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
import {agoraApi} from '@/api/modules/agora'
import type {AgoraCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  appId: '',
  appCertificate: '',
  restCustomerId: '',
  restCustomerSecret: '',
  tokenExpireSeconds: 21600,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  appId: [
    {required: true, message: '请输入 AppId', trigger: 'blur'},
    {min: 1, max: 64, message: 'AppId 长度在 1-64 个字符', trigger: 'blur'},
  ],
  appCertificate: [
    {required: true, message: '请输入 AppCertificate', trigger: 'blur'},
    {min: 1, max: 128, message: 'AppCertificate 长度在 1-128 个字符', trigger: 'blur'},
  ],
  restCustomerId: [
    {max: 64, message: 'REST CustomerId 长度不能超过 64 个字符', trigger: 'blur'},
  ],
  restCustomerSecret: [
    {max: 128, message: 'REST CustomerSecret 长度不能超过 128 个字符', trigger: 'blur'},
  ],
  tokenExpireSeconds: [
    {required: true, message: '请输入 Token 有效期', trigger: 'blur'},
    {type: 'number', min: 60, message: 'Token 有效期不能小于 60 秒', trigger: 'blur'},
  ],
})

const applyCfg = (cfg: AgoraCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.appId = ''
    formData.appCertificate = ''
    formData.restCustomerId = ''
    formData.restCustomerSecret = ''
    formData.tokenExpireSeconds = 21600
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.appId = cfg.appId || ''
  formData.appCertificate = cfg.appCertificate || ''
  formData.restCustomerId = cfg.restCustomerId || ''
  formData.restCustomerSecret = cfg.restCustomerSecret || ''
  formData.tokenExpireSeconds = cfg.tokenExpireSeconds || 21600
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await agoraApi.getAgoraCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取声网配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await agoraApi.saveAgoraCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      appId: formData.appId.trim(),
      appCertificate: formData.appCertificate.trim(),
      restCustomerId: formData.restCustomerId.trim(),
      restCustomerSecret: formData.restCustomerSecret.trim(),
      tokenExpireSeconds: formData.tokenExpireSeconds,
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
    console.error('保存声网配置失败:', error)
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

.cfg-form {
  max-width: 720px;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>

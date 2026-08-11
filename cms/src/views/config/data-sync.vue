<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>数据同步配置</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>配置跨环境数据同步的目标 API 与 Token。测试服与正式服填写相同的 Token，同步请求通过 Header 携带校验。</p>
      </el-alert>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="目标 API 根地址" prop="targetApiBase">
          <el-input
              v-model="formData.targetApiBase"
              clearable
              placeholder="例如 https://api.example.com"
          />
        </el-form-item>

        <el-form-item label="同步 Token" prop="token">
          <div class="token-input">
            <el-input
                v-model="formData.token"
                clearable
                placeholder="测试服与正式服保持一致"
                show-password
                type="password"
            />
            <el-button @click="generateSyncToken">随机生成</el-button>
            <el-button :disabled="!formData.token" @click="copySyncToken">复制</el-button>
          </div>
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
import {dataSyncApi} from '@/api/modules/data-sync'
import type {DataSyncCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  targetApiBase: '',
  token: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  targetApiBase: [{required: true, message: '请填写目标 API 根地址', trigger: 'blur'}],
  token: [{required: true, message: '请填写同步 Token', trigger: 'blur'}],
})

const generateSyncToken = () => {
  formData.token = crypto.randomUUID().replace(/-/g, '')
}

const copySyncToken = async () => {
  const value = formData.token.trim()
  if (!value) {
    ElMessage.warning('Token 为空，无法复制')
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    console.error('复制 Token 失败:', error)
    ElMessage.error('复制失败')
  }
}

const applyCfg = (cfg: DataSyncCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.targetApiBase = ''
    formData.token = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.targetApiBase = cfg.targetApiBase || ''
  formData.token = cfg.token || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await dataSyncApi.getDataSyncCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取数据同步配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await dataSyncApi.saveDataSyncCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      targetApiBase: formData.targetApiBase.trim(),
      token: formData.token.trim(),
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
    console.error('保存数据同步配置失败:', error)
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

.token-input {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.token-input .el-input {
  flex: 1;
}
</style>

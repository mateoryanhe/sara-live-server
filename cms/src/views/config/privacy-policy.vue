<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>隐私政策配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item label="隐私政策 URL" prop="privacyPolicyUrl">
          <el-input
              v-model="formData.privacyPolicyUrl"
              clearable
              placeholder="如 https://cdn.example.com/privacy.html"
          />
          <span class="form-tip">App 将通过系统配置接口读取该地址，用于展示隐私政策页面</span>
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
import {privacyPolicyApi} from '@/api/modules/privacy-policy'
import type {PrivacyPolicyCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  privacyPolicyUrl: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  privacyPolicyUrl: [
    {max: 512, message: 'URL 长度不能超过 512', trigger: 'blur'},
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
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
      },
      trigger: 'blur',
    },
  ],
})

const applyCfg = (cfg: PrivacyPolicyCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.privacyPolicyUrl = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.privacyPolicyUrl = cfg.privacyPolicyUrl || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await privacyPolicyApi.getPrivacyPolicyCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取隐私政策配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await privacyPolicyApi.savePrivacyPolicyCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      privacyPolicyUrl: formData.privacyPolicyUrl.trim(),
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
    console.error('保存隐私政策配置失败:', error)
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
  max-width: 760px;
}

.form-tip {
  display: block;
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}
</style>

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>资源与头像配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="资源域名" prop="resourceDomain">
          <el-input
              v-model="formData.resourceDomain"
              clearable
              placeholder="如 http://127.0.0.1 或 https://cdn.example.com"
          />
          <span class="form-tip">留空保存后使用默认 http://127.0.0.1；可只填 127.0.0.1，会自动补 http://</span>
        </el-form-item>

        <el-form-item label="默认头像 URL" prop="defaultAvatarUrl">
          <el-input
              v-model="formData.defaultAvatarUrl"
              clearable
              placeholder="用户未上传头像时展示的完整图片地址"
          />
          <span class="form-tip">留空则使用系统内置默认头像</span>
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
import {uploadResourceApi} from '@/api/modules/upload-resource'
import type {UploadResourceCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  resourceDomain: 'http://127.0.0.1',
  defaultAvatarUrl: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  resourceDomain: [{max: 256, message: '域名长度不能超过 256', trigger: 'blur'}],
  defaultAvatarUrl: [{max: 512, message: 'URL 长度不能超过 512', trigger: 'blur'}],
})

const applyCfg = (cfg: UploadResourceCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.resourceDomain = 'http://127.0.0.1'
    formData.defaultAvatarUrl = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.resourceDomain = cfg.resourceDomain || 'http://127.0.0.1'
  formData.defaultAvatarUrl = cfg.defaultAvatarUrl || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await uploadResourceApi.getUploadResourceCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取资源配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await uploadResourceApi.saveUploadResourceCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      resourceDomain: formData.resourceDomain.trim(),
      defaultAvatarUrl: formData.defaultAvatarUrl.trim(),
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
    console.error('保存资源配置失败:', error)
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
  display: block;
  margin-top: 6px;
  margin-left: 0;
  color: #909399;
  font-size: 13px;
  line-height: 1.4;
}
</style>

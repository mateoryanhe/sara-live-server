<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>预热配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="200px">
        <el-form-item label="最近登录用户预热数量" prop="recentLoginLimit">
          <el-input-number v-model="formData.recentLoginLimit" :min="1" :max="10000" controls-position="right"/>
          <div class="form-tip">
            服务启动时按 user_infos.last_login_time 倒序预热最近 N 个用户的缓存，默认 100；修改后需重启服务生效
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
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import preloadCfgApi from '@/api/modules/preload-cfg'
import type {PreloadCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  recentLoginLimit: 100,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules: FormRules = {
  recentLoginLimit: [
    {required: true, message: '请输入预热数量', trigger: 'change'},
    {type: 'number', min: 1, max: 10000, message: '预热数量需在 1-10000 之间', trigger: 'change'},
  ],
}

const applyCfg = (cfg: PreloadCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.recentLoginLimit = Number(cfg?.recentLoginLimit) || 100
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await preloadCfgApi.getPreloadCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取预热配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const response = await preloadCfgApi.savePreloadCfg({
        id: formData.id === '0' ? 0 : Number(formData.id),
        recentLoginLimit: formData.recentLoginLimit,
      })
      if (response?.success) {
        ElMessage.success('保存成功，重启服务后生效')
        if (response.id) {
          formData.id = response.id
        }
        await fetchCfg()
      } else {
        ElMessage.error('保存失败')
      }
    } catch (error) {
      console.error('保存预热配置失败:', error)
    }
  })
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
}

.cfg-form {
  max-width: 720px;
}

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}
</style>

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>直播配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="付费弹幕价格" prop="paidDanmakuPrice">
          <el-input-number
              v-model="formData.paidDanmakuPrice"
              :min="0"
              :precision="4"
              :step="0.0001"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">钻石，保留 4 位小数</span>
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
import {liveCfgApi} from '@/api/modules/liveCfg'
import type {LiveCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  paidDanmakuPrice: 0,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  paidDanmakuPrice: [
    {required: true, message: '请输入付费弹幕价格', trigger: 'blur'},
    {type: 'number', min: 0, message: '价格不能小于 0', trigger: 'blur'},
  ],
})

const applyCfg = (cfg: LiveCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.paidDanmakuPrice = 0
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.paidDanmakuPrice = cfg.paidDanmakuPrice ?? 0
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await liveCfgApi.getLiveCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取直播配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await liveCfgApi.saveLiveCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      paidDanmakuPrice: formData.paidDanmakuPrice,
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
    console.error('保存直播配置失败:', error)
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

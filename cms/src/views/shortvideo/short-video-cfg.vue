<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>短视频配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item label="最大文件大小" prop="maxFileSizeMB">
          <el-input-number
              v-model="formData.maxFileSizeMB"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">MB，保存后按字节写入数据库</span>
        </el-form-item>

        <el-form-item label="最大时长" prop="maxDuration">
          <el-input-number
              v-model="formData.maxDuration"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">秒</span>
        </el-form-item>

        <el-form-item label="免费观看时长" prop="freeWatchSeconds">
          <el-input-number
              v-model="formData.freeWatchSeconds"
              :min="0"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">秒，观看前该时长免费</span>
        </el-form-item>

        <el-form-item label="短视频入口" prop="entryEnabled">
          <el-switch
              v-model="formData.entryEnabled"
              :active-value="1"
              :inactive-value="0"
              active-text="开启"
              inactive-text="关闭"
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
import {shortVideoApi} from '@/api/modules/shortVideo'
import type {ShortVideoCfg} from '@/types/api'

const MB = 1024 * 1024

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  maxFileSizeMB: 100,
  maxDuration: 60,
  freeWatchSeconds: 7,
  entryEnabled: 1,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  maxFileSizeMB: [
    {required: true, message: '请输入最大文件大小', trigger: 'blur'},
    {type: 'number', min: 1, message: '最大文件大小必须大于0', trigger: 'blur'},
  ],
  maxDuration: [
    {required: true, message: '请输入最大时长', trigger: 'blur'},
    {type: 'number', min: 1, message: '最大时长必须大于0', trigger: 'blur'},
  ],
  freeWatchSeconds: [
    {required: true, message: '请输入免费观看时长', trigger: 'blur'},
    {type: 'number', min: 0, message: '免费观看时长不能小于0', trigger: 'blur'},
  ],
  entryEnabled: [
    {required: true, message: '请选择入口开关', trigger: 'change'},
  ],
})

const applyCfg = (cfg: ShortVideoCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.maxFileSizeMB = 100
    formData.maxDuration = 60
    formData.freeWatchSeconds = 7
    formData.entryEnabled = 1
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.maxFileSizeMB = Math.max(1, Math.round(cfg.maxFileSize / MB))
  formData.maxDuration = cfg.maxDuration || 60
  formData.freeWatchSeconds = cfg.freeWatchSeconds ?? 7
  formData.entryEnabled = cfg.entryEnabled ?? 1
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await shortVideoApi.getShortVideoCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取短视频配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await shortVideoApi.saveShortVideoCfg({
      id: formData.id,
      maxFileSize: formData.maxFileSizeMB * MB,
      maxDuration: formData.maxDuration,
      freeWatchSeconds: formData.freeWatchSeconds,
      entryEnabled: formData.entryEnabled,
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
    console.error('保存短视频配置失败:', error)
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

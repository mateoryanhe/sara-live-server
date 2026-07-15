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
          <span class="form-tip">云播放器接口需要配置</span>
        </el-form-item>

        <el-form-item label="REST CustomerSecret" prop="restCustomerSecret">
          <el-input
              v-model="formData.restCustomerSecret"
              clearable
              placeholder="请输入声网 REST CustomerSecret"
              show-password
              type="password"
          />
          <span class="form-tip">云播放器接口需要配置</span>
        </el-form-item>

        <el-form-item label="云播放器区域" prop="cloudPlayerRegion">
          <el-select v-model="formData.cloudPlayerRegion" placeholder="请选择区域" style="width: 220px">
            <el-option label="中国大陆 (cn)" value="cn"/>
            <el-option label="亚太 (ap)" value="ap"/>
            <el-option label="欧洲 (eu)" value="eu"/>
            <el-option label="北美 (na)" value="na"/>
          </el-select>
        </el-form-item>

        <el-form-item label="Token有效期" prop="tokenExpireHours">
          <el-input-number
              v-model="formData.tokenExpireHours"
              :max="TOKEN_EXPIRE_MAX_HOURS"
              :min="TOKEN_EXPIRE_MIN_HOURS"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">小时，范围 {{ TOKEN_EXPIRE_MIN_HOURS }}-{{ TOKEN_EXPIRE_MAX_HOURS }} 小时，默认 {{ TOKEN_EXPIRE_DEFAULT_HOURS }} 小时</span>
        </el-form-item>

        <el-form-item label="提前刷新阈值" prop="tokenRefreshHours">
          <el-input-number
              v-model="formData.tokenRefreshHours"
              :max="maxTokenRefreshHours"
              :min="TOKEN_REFRESH_MIN_HOURS"
              :precision="0"
              controls-position="right"
              style="width: 220px"
          />
          <span class="form-tip">小时，范围 {{ TOKEN_REFRESH_MIN_HOURS }}-{{ maxTokenRefreshHours }} 小时，且比 Token 有效期至少少 {{ TOKEN_REFRESH_AHEAD_GAP_HOURS }} 小时</span>
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {agoraApi} from '@/api/modules/agora'
import type {AgoraCfg} from '@/types/api'

const TOKEN_EXPIRE_MIN_HOURS = 4
const TOKEN_EXPIRE_MAX_HOURS = 24
const TOKEN_EXPIRE_DEFAULT_HOURS = 24
const TOKEN_REFRESH_MIN_HOURS = 2
const TOKEN_REFRESH_AHEAD_GAP_HOURS = 2
const TOKEN_REFRESH_DEFAULT_HOURS = TOKEN_EXPIRE_DEFAULT_HOURS - TOKEN_REFRESH_AHEAD_GAP_HOURS
const SECONDS_PER_HOUR = 3600

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  appId: '',
  appCertificate: '',
  restCustomerId: '',
  restCustomerSecret: '',
  cloudPlayerRegion: 'cn',
  tokenExpireHours: TOKEN_EXPIRE_DEFAULT_HOURS,
  tokenRefreshHours: TOKEN_REFRESH_DEFAULT_HOURS,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const maxTokenRefreshHours = computed(() => {
  return Math.max(TOKEN_REFRESH_MIN_HOURS, formData.tokenExpireHours - TOKEN_REFRESH_AHEAD_GAP_HOURS)
})

const clampTokenExpireHours = (hours: number) => {
  if (!Number.isFinite(hours) || hours <= 0) {
    return TOKEN_EXPIRE_DEFAULT_HOURS
  }
  return Math.min(TOKEN_EXPIRE_MAX_HOURS, Math.max(TOKEN_EXPIRE_MIN_HOURS, Math.round(hours)))
}

const clampTokenRefreshHours = (hours: number, expireHours: number) => {
  const maxHours = Math.max(TOKEN_REFRESH_MIN_HOURS, expireHours - TOKEN_REFRESH_AHEAD_GAP_HOURS)
  if (!Number.isFinite(hours)) {
    return maxHours
  }
  return Math.min(maxHours, Math.max(TOKEN_REFRESH_MIN_HOURS, Math.round(hours)))
}

const secondsToHours = (seconds: number, fallback: number) => {
  if (!Number.isFinite(seconds) || seconds < 0) {
    return fallback
  }
  return Math.round(seconds / SECONDS_PER_HOUR)
}

const hoursToSeconds = (hours: number) => Math.round(hours) * SECONDS_PER_HOUR

const syncTokenRefreshHours = () => {
  formData.tokenRefreshHours = clampTokenRefreshHours(formData.tokenRefreshHours, formData.tokenExpireHours)
}

watch(() => formData.tokenExpireHours, () => {
  syncTokenRefreshHours()
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
  tokenExpireHours: [
    {required: true, message: '请输入 Token 有效期', trigger: 'blur'},
    {
      type: 'number',
      min: TOKEN_EXPIRE_MIN_HOURS,
      max: TOKEN_EXPIRE_MAX_HOURS,
      message: `Token 有效期需在 ${TOKEN_EXPIRE_MIN_HOURS}-${TOKEN_EXPIRE_MAX_HOURS} 小时之间`,
      trigger: 'blur',
    },
  ],
  tokenRefreshHours: [
    {required: true, message: '请输入提前刷新阈值', trigger: 'blur'},
    {
      validator: (_rule: unknown, value: number, callback: (error?: Error) => void) => {
        const max = maxTokenRefreshHours.value
        if (!Number.isFinite(value) || value < TOKEN_REFRESH_MIN_HOURS || value > max) {
          callback(new Error(`提前刷新阈值需在 ${TOKEN_REFRESH_MIN_HOURS}-${max} 小时之间`))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
})

const applyCfg = (cfg: AgoraCfg | null | undefined) => {
  if (!cfg) {
    formData.id = '0'
    formData.appId = ''
    formData.appCertificate = ''
    formData.restCustomerId = ''
    formData.restCustomerSecret = ''
    formData.cloudPlayerRegion = 'cn'
    formData.tokenExpireHours = TOKEN_EXPIRE_DEFAULT_HOURS
    formData.tokenRefreshHours = TOKEN_REFRESH_DEFAULT_HOURS
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.appId = cfg.appId || ''
  formData.appCertificate = cfg.appCertificate || ''
  formData.restCustomerId = cfg.restCustomerId || ''
  formData.restCustomerSecret = cfg.restCustomerSecret || ''
  formData.cloudPlayerRegion = cfg.cloudPlayerRegion || 'cn'
  formData.tokenExpireHours = clampTokenExpireHours(secondsToHours(cfg.tokenExpireSeconds, TOKEN_EXPIRE_DEFAULT_HOURS))
  formData.tokenRefreshHours = clampTokenRefreshHours(
      secondsToHours(cfg.tokenRefreshSeconds, TOKEN_REFRESH_DEFAULT_HOURS),
      formData.tokenExpireHours,
  )
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
      cloudPlayerRegion: formData.cloudPlayerRegion || 'cn',
      tokenExpireSeconds: hoursToSeconds(formData.tokenExpireHours),
      tokenRefreshSeconds: hoursToSeconds(formData.tokenRefreshHours),
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

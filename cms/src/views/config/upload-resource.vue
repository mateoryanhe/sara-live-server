<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>资源与头像配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-divider content-position="left">静态资源</el-divider>

        <el-form-item label="资源域名" prop="resourceDomain">
          <el-input
              v-model="formData.resourceDomain"
              clearable
              placeholder="如 http://127.0.0.1 或 https://cdn.example.com"
          />
          <span class="form-tip">留空默认 http://127.0.0.1；开启图片审核时须填写阿里云可公网访问的域名（不能是 127.0.0.1/内网 IP）</span>
        </el-form-item>

        <el-form-item label="默认头像 URL" prop="defaultAvatarUrl">
          <el-input
              v-model="formData.defaultAvatarUrl"
              clearable
              placeholder="用户未上传头像时展示的完整图片地址"
          />
          <span class="form-tip">留空则使用系统内置默认头像</span>
        </el-form-item>

        <el-divider content-position="left">App 图片审核（阿里云）</el-divider>

        <el-form-item label="开启图片审核">
          <el-switch v-model="formData.imageModerationEnabled" active-text="开启" inactive-text="关闭"/>
          <span class="form-tip">仅对 App 端上传的图片生效（如头像），CMS 后台上传不审核；图片先存本地，再通过资源域名拼 imageUrl 调用阿里云 ImageModeration 检测</span>
        </el-form-item>

        <template v-if="formData.imageModerationEnabled">
          <el-form-item label="AccessKey ID" prop="imageModerationAccessKeyId">
            <el-input v-model="formData.imageModerationAccessKeyId" clearable placeholder="RAM 用户 AccessKey ID"/>
          </el-form-item>

          <el-form-item label="AccessKey Secret" prop="imageModerationAccessKeySecret">
            <el-input
                v-model="formData.imageModerationAccessKeySecret"
                clearable
                placeholder="留空表示不修改已保存的 Secret"
                show-password
                type="password"
            />
          </el-form-item>

          <el-form-item label="地域 RegionId" prop="imageModerationRegionId">
            <el-input v-model="formData.imageModerationRegionId" clearable placeholder="cn-shanghai"/>
          </el-form-item>

          <el-form-item label="接入点 Endpoint" prop="imageModerationEndpoint">
            <el-input
                v-model="formData.imageModerationEndpoint"
                clearable
                placeholder="green-cip.cn-shanghai.aliyuncs.com"
            />
          </el-form-item>

          <el-form-item label="审核 Service" prop="imageModerationService">
            <el-input
                v-model="formData.imageModerationService"
                clearable
                placeholder="profilePhotoCheck（头像）"
            />
            <span class="form-tip">头像推荐 profilePhotoCheck；通用可用 baselineCheck</span>
          </el-form-item>
        </template>

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
import {onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {uploadResourceApi} from '@/api/modules/upload-resource'
import type {UploadResourceCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()
const imageSecretTouched = ref(false)

const formData = reactive({
  id: '0',
  resourceDomain: 'http://127.0.0.1',
  defaultAvatarUrl: '',
  imageModerationEnabled: false,
  imageModerationAccessKeyId: '',
  imageModerationAccessKeySecret: '',
  imageModerationRegionId: 'cn-shanghai',
  imageModerationEndpoint: 'green-cip.cn-shanghai.aliyuncs.com',
  imageModerationService: 'profilePhotoCheck',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  resourceDomain: [{max: 256, message: '域名长度不能超过 256', trigger: 'blur'}],
  defaultAvatarUrl: [{max: 512, message: 'URL 长度不能超过 512', trigger: 'blur'}],
  imageModerationAccessKeyId: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.imageModerationEnabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error('开启图片审核时请填写 AccessKey ID'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
})

watch(
    () => formData.imageModerationAccessKeySecret,
    () => {
      imageSecretTouched.value = true
    },
)

const applyCfg = (cfg: UploadResourceCfg | null | undefined) => {
  imageSecretTouched.value = false
  if (!cfg) {
    formData.id = '0'
    formData.resourceDomain = 'http://127.0.0.1'
    formData.defaultAvatarUrl = ''
    formData.imageModerationEnabled = false
    formData.imageModerationAccessKeyId = ''
    formData.imageModerationAccessKeySecret = ''
    formData.imageModerationRegionId = 'cn-shanghai'
    formData.imageModerationEndpoint = 'green-cip.cn-shanghai.aliyuncs.com'
    formData.imageModerationService = 'profilePhotoCheck'
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.resourceDomain = cfg.resourceDomain || 'http://127.0.0.1'
  formData.defaultAvatarUrl = cfg.defaultAvatarUrl || ''
  formData.imageModerationEnabled = !!cfg.imageModerationEnabled
  formData.imageModerationAccessKeyId = cfg.imageModerationAccessKeyId || ''
  formData.imageModerationAccessKeySecret = ''
  formData.imageModerationRegionId = cfg.imageModerationRegionId || 'cn-shanghai'
  formData.imageModerationEndpoint = cfg.imageModerationEndpoint || 'green-cip.cn-shanghai.aliyuncs.com'
  formData.imageModerationService = cfg.imageModerationService || 'profilePhotoCheck'
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
      imageModerationEnabled: formData.imageModerationEnabled,
      imageModerationAccessKeyId: formData.imageModerationAccessKeyId.trim(),
      imageModerationAccessKeySecret: imageSecretTouched.value
          ? formData.imageModerationAccessKeySecret.trim()
          : '',
      imageModerationRegionId: formData.imageModerationRegionId.trim(),
      imageModerationEndpoint: formData.imageModerationEndpoint.trim(),
      imageModerationService: formData.imageModerationService.trim(),
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
  max-width: 760px;
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

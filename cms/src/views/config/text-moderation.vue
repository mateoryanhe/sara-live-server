<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>敏感词过滤（阿里云）</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item label="开启过滤">
          <el-switch v-model="formData.enabled" active-text="开启" inactive-text="关闭"/>
          <span class="form-tip">关闭后 App 文本不再调用阿里云审核（直接放行）</span>
        </el-form-item>

        <template v-if="formData.enabled">
          <el-form-item label="AccessKey ID" prop="accessKeyId">
            <el-input v-model="formData.accessKeyId" clearable placeholder="RAM 用户 AccessKey ID"/>
          </el-form-item>

          <el-form-item label="AccessKey Secret" prop="accessKeySecret">
            <el-input
                v-model="formData.accessKeySecret"
                clearable
                placeholder="留空表示不修改已保存的 Secret"
                show-password
                type="password"
            />
          </el-form-item>

          <el-form-item label="地域 RegionId" prop="regionId">
            <el-input v-model="formData.regionId" clearable placeholder="如 cn-shanghai"/>
          </el-form-item>

          <el-form-item label="接入点 Endpoint" prop="endpoint">
            <el-input v-model="formData.endpoint" clearable placeholder="如 green-cip.cn-shanghai.aliyuncs.com"/>
          </el-form-item>

          <el-form-item label="公聊/私信 Service" prop="chatService">
            <el-input v-model="formData.chatService" clearable placeholder="默认 chat_detection"/>
          </el-form-item>

          <el-form-item label="昵称 Service" prop="nicknameService">
            <el-input v-model="formData.nicknameService" clearable placeholder="默认 nickname_detection"/>
          </el-form-item>

          <el-form-item label="评论/公告 Service" prop="commentService">
            <el-input v-model="formData.commentService" clearable placeholder="默认 comment_detection"/>
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
import {textModerationApi} from '@/api/modules/text-moderation'
import type {TextModerationCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()
const secretTouched = ref(false)

const formData = reactive({
  id: '0',
  enabled: false,
  accessKeyId: '',
  accessKeySecret: '',
  regionId: 'cn-shanghai',
  endpoint: 'green-cip.cn-shanghai.aliyuncs.com',
  chatService: 'chat_detection',
  nicknameService: 'nickname_detection',
  commentService: 'comment_detection',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const formRules = reactive({
  accessKeyId: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error('开启过滤时请填写 AccessKey ID'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  regionId: [{required: true, message: '请填写 RegionId', trigger: 'blur'}],
  endpoint: [{required: true, message: '请填写 Endpoint', trigger: 'blur'}],
})

watch(
    () => formData.accessKeySecret,
    () => {
      secretTouched.value = true
    },
)

const applyCfg = (cfg: TextModerationCfg | null | undefined) => {
  secretTouched.value = false
  if (!cfg) {
    formData.id = '0'
    formData.enabled = false
    formData.accessKeyId = ''
    formData.accessKeySecret = ''
    formData.regionId = 'cn-shanghai'
    formData.endpoint = 'green-cip.cn-shanghai.aliyuncs.com'
    formData.chatService = 'chat_detection'
    formData.nicknameService = 'nickname_detection'
    formData.commentService = 'comment_detection'
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.enabled = !!cfg.enabled
  formData.accessKeyId = cfg.accessKeyId || ''
  formData.accessKeySecret = ''
  formData.regionId = cfg.regionId || 'cn-shanghai'
  formData.endpoint = cfg.endpoint || 'green-cip.cn-shanghai.aliyuncs.com'
  formData.chatService = cfg.chatService || 'chat_detection'
  formData.nicknameService = cfg.nicknameService || 'nickname_detection'
  formData.commentService = cfg.commentService || 'comment_detection'
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await textModerationApi.getTextModerationCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取文本审核配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await textModerationApi.saveTextModerationCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      enabled: formData.enabled,
      accessKeyId: formData.accessKeyId.trim(),
      accessKeySecret: secretTouched.value ? formData.accessKeySecret.trim() : '',
      regionId: formData.regionId.trim(),
      endpoint: formData.endpoint.trim(),
      chatService: formData.chatService.trim(),
      nicknameService: formData.nicknameService.trim(),
      commentService: formData.commentService.trim(),
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
    console.error('保存文本审核配置失败:', error)
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
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>

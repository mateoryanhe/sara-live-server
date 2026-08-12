<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.TextModerationCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.textModeration.enableFilter')">
          <el-switch
              v-model="formData.enabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <span class="form-tip">{{ t('pages.textModeration.enableFilterTip') }}</span>
        </el-form-item>

        <template v-if="formData.enabled">
          <el-form-item label="AccessKey ID" prop="accessKeyId">
            <el-input v-model="formData.accessKeyId" clearable :placeholder="t('pages.textModeration.accessKeyIdPlaceholder')"/>
          </el-form-item>

          <el-form-item label="AccessKey Secret" prop="accessKeySecret">
            <el-input
                v-model="formData.accessKeySecret"
                clearable
                :placeholder="t('pages.textModeration.accessKeySecretPlaceholder')"
                show-password
                type="password"
            />
          </el-form-item>

          <el-form-item :label="t('pages.textModeration.regionId')" prop="regionId">
            <el-input v-model="formData.regionId" clearable :placeholder="t('pages.textModeration.regionIdPlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.textModeration.endpoint')" prop="endpoint">
            <el-input v-model="formData.endpoint" clearable :placeholder="t('pages.textModeration.endpointPlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.textModeration.chatService')" prop="chatService">
            <el-input v-model="formData.chatService" clearable :placeholder="t('pages.textModeration.chatServicePlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.textModeration.nicknameService')" prop="nicknameService">
            <el-input v-model="formData.nicknameService" clearable :placeholder="t('pages.textModeration.nicknameServicePlaceholder')"/>
          </el-form-item>

          <el-form-item :label="t('pages.textModeration.commentService')" prop="commentService">
            <el-input v-model="formData.commentService" clearable :placeholder="t('pages.textModeration.commentServicePlaceholder')"/>
          </el-form-item>
        </template>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.textModeration.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {textModerationApi} from '@/api/modules/text-moderation'
import type {TextModerationCfg} from '@/types/api'

const {t} = useI18n()
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

const formRules = computed(() => ({
  accessKeyId: [
    {
      validator: (_: unknown, value: string, callback: (e?: Error) => void) => {
        if (!formData.enabled) {
          callback()
          return
        }
        if (!value?.trim()) {
          callback(new Error(t('pages.textModeration.accessKeyIdRequired')))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
  regionId: [{required: true, message: t('pages.textModeration.regionIdRequired'), trigger: 'blur'}],
  endpoint: [{required: true, message: t('pages.textModeration.endpointRequired'), trigger: 'blur'}],
}))

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
    console.error('fetch text moderation cfg failed:', error)
    ElMessage.error(t('pages.textModeration.fetchCfgFailed'))
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
      ElMessage.success(t('pages.textModeration.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.textModeration.saveFailed'))
    }
  } catch (error) {
    console.error('save text moderation cfg failed:', error)
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

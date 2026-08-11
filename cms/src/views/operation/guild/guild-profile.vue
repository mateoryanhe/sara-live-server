<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>工会基础信息</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>仅可编辑当前登录 CMS 账号作为会长（leader_id）关联的工会信息。</p>
        <p>若未匹配到工会，请联系管理员在「工会管理」中设置会长。</p>
      </el-alert>

      <el-empty v-if="!loading && !hasGuild" description="未找到与当前账号关联的工会"/>

      <el-form
          v-else
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="profile-form"
          label-width="100px"
      >
        <el-form-item label="工会ID">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <el-form-item label="工会名称" prop="name">
          <el-input v-model="formData.name" maxlength="32" placeholder="请输入工会名称" show-word-limit/>
        </el-form-item>
        <el-form-item label="银行卡" prop="bankCard">
          <el-input v-model="formData.bankCard" maxlength="64" placeholder="请输入银行卡信息" show-word-limit/>
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="formData.contact" maxlength="64" placeholder="请输入联系方式" show-word-limit/>
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input
              v-model="formData.description"
              :rows="4"
              maxlength="255"
              placeholder="请输入工会简介"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
        <el-form-item v-if="formData.updatedAt" label="最近更新">
          <span>{{ formData.updatedAt }}</span>
        </el-form-item>
        <el-form-item>
          <el-button v-if="can('save')" :loading="saving" type="primary" @click="handleSave">保存</el-button>
          <el-button @click="fetchProfile">刷新</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {guildApi} from '@/api'
import type {MyGuildProfile} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {can} = usePagePermission('GuildProfileManagement')

const loading = ref(false)
const saving = ref(false)
const hasGuild = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '',
  name: '',
  bankCard: '',
  contact: '',
  description: '',
  updatedAt: '',
})

const formRules: FormRules = {
  name: [
    {required: true, message: '请输入工会名称', trigger: 'blur'},
    {min: 2, max: 32, message: '工会名称长度在2-32个字符', trigger: 'blur'},
  ],
  bankCard: [
    {max: 64, message: '银行卡信息不能超过64个字符', trigger: 'blur'},
  ],
  contact: [
    {max: 64, message: '联系方式不能超过64个字符', trigger: 'blur'},
  ],
  description: [
    {max: 255, message: '简介长度不能超过255个字符', trigger: 'blur'},
  ],
}

const applyProfile = (profile: MyGuildProfile | null | undefined) => {
  if (!profile?.id) {
    hasGuild.value = false
    formData.id = ''
    formData.name = ''
    formData.bankCard = ''
    formData.contact = ''
    formData.description = ''
    formData.updatedAt = ''
    return
  }
  hasGuild.value = true
  formData.id = profile.id
  formData.name = profile.name || ''
  formData.bankCard = profile.bankCard || ''
  formData.contact = profile.contact || ''
  formData.description = profile.description || ''
  formData.updatedAt = profile.updatedAt || ''
}

const fetchProfile = async () => {
  loading.value = true
  try {
    const response = await guildApi.getMyGuildProfile()
    applyProfile(response)
  } catch (error) {
    console.error('获取工会基础信息失败:', error)
    applyProfile(null)
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value || !hasGuild.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      const response = await guildApi.updateMyGuildProfile({
        name: formData.name.trim(),
        bankCard: formData.bankCard.trim(),
        contact: formData.contact.trim(),
        description: formData.description.trim(),
      })
      if (response?.success) {
        ElMessage.success('保存成功')
        await fetchProfile()
      } else {
        ElMessage.error('保存失败')
      }
    } catch (error) {
      console.error('保存工会基础信息失败:', error)
    } finally {
      saving.value = false
    }
  })
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.tip-alert {
  margin-bottom: 20px;
}

.profile-form {
  max-width: 640px;
}
</style>

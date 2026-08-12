<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildProfileManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.guildProfile.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.guildProfile.tipLine1') }}</p>
        <p>{{ t('pages.guildProfile.tipLine2') }}</p>
      </el-alert>

      <el-empty v-if="!loading && !hasGuild" :description="t('pages.guildProfile.noGuild')"/>

      <el-form
          v-else
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="profile-form"
          label-width="100px"
      >
        <el-form-item :label="t('pages.guildProfile.guildId')">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.guildProfile.guildName')" prop="name">
          <el-input v-model="formData.name" maxlength="32" :placeholder="t('pages.guildProfile.guildNamePlaceholder')" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('pages.guildProfile.bankCard')" prop="bankCard">
          <el-input v-model="formData.bankCard" maxlength="64" :placeholder="t('pages.guildProfile.bankCardPlaceholder')" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('pages.guildProfile.contact')" prop="contact">
          <el-input v-model="formData.contact" maxlength="64" :placeholder="t('pages.guildProfile.contactPlaceholder')" show-word-limit/>
        </el-form-item>
        <el-form-item :label="t('pages.guildProfile.description')" prop="description">
          <el-input
              v-model="formData.description"
              :rows="4"
              maxlength="255"
              :placeholder="t('pages.guildProfile.descriptionPlaceholder')"
              show-word-limit
              type="textarea"
          />
        </el-form-item>
        <el-form-item v-if="formData.updatedAt" :label="t('pages.guildProfile.lastUpdated')">
          <span>{{ formData.updatedAt }}</span>
        </el-form-item>
        <el-form-item>
          <el-button v-if="can('save')" :loading="saving" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
          <el-button @click="fetchProfile">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import {guildApi} from '@/api'
import type {MyGuildProfile} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {can} = usePagePermission('GuildProfileManagement')

const {t} = useI18n()
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

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.guildProfile.nameRequired'), trigger: 'blur'},
    {min: 2, max: 32, message: t('pages.guildProfile.nameLength'), trigger: 'blur'},
  ],
  bankCard: [
    {max: 64, message: t('pages.guildProfile.bankCardMaxLength'), trigger: 'blur'},
  ],
  contact: [
    {max: 64, message: t('pages.guildProfile.contactMaxLength'), trigger: 'blur'},
  ],
  description: [
    {max: 255, message: t('pages.guildProfile.descriptionMaxLength'), trigger: 'blur'},
  ],
}))

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
    console.error('fetch guild profile failed:', error)
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
        ElMessage.success(t('common.saveConfig'))
        await fetchProfile()
      } else {
        ElMessage.error(t('pages.guildProfile.saveFailed'))
      }
    } catch (error) {
      console.error('save guild profile failed:', error)
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

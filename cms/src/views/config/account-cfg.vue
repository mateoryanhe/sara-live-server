<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AccountCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="200px">
        <el-form-item :label="t('pages.accountCfg.cancelAccountByCodeEnabled')">
          <el-switch
              v-model="formData.cancelAccountByCodeEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">
            {{ t('pages.accountCfg.cancelAccountByCodeTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.accountCfg.simulatorLoginEnabled')">
          <el-switch
              v-model="formData.simulatorLoginEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">
            {{ t('pages.accountCfg.simulatorLoginTip') }}
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.accountCfg.lastUpdated')">
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
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import accountCfgApi from '@/api/modules/account-cfg'
import type {AccountCfg} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  cancelAccountByCodeEnabled: false,
  simulatorLoginEnabled: true,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: AccountCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.cancelAccountByCodeEnabled = !!cfg?.cancelAccountByCodeEnabled
  formData.simulatorLoginEnabled = cfg?.simulatorLoginEnabled !== false
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await accountCfgApi.getAccountCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch account cfg failed:', error)
    ElMessage.error(t('pages.accountCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    const response = await accountCfgApi.saveAccountCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      cancelAccountByCodeEnabled: formData.cancelAccountByCodeEnabled,
      simulatorLoginEnabled: formData.simulatorLoginEnabled,
    })
    if (response?.success) {
      ElMessage.success(t('pages.accountCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.accountCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save account cfg failed:', error)
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

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>账号配置</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="200px">
        <el-form-item label="注销码销户(官网)">
          <el-switch
              v-model="formData.cancelAccountByCodeEnabled"
              active-text="开启"
              inactive-text="关闭"
          />
          <div class="form-tip">
            控制官网公开接口 POST /userInfo/cancelAccountByCode 是否可用，默认关闭
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
import {ElMessage} from 'element-plus'
import accountCfgApi from '@/api/modules/account-cfg'
import type {AccountCfg} from '@/types/api'

const loading = ref(false)
const formRef = ref()

const formData = reactive({
  id: '0',
  cancelAccountByCodeEnabled: false,
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const applyCfg = (cfg: AccountCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.cancelAccountByCodeEnabled = !!cfg?.cancelAccountByCodeEnabled
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await accountCfgApi.getAccountCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('获取账号配置失败:', error)
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  try {
    const response = await accountCfgApi.saveAccountCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      cancelAccountByCodeEnabled: formData.cancelAccountByCodeEnabled,
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
    console.error('保存账号配置失败:', error)
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

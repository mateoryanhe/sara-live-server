<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <span>{{ t('menu.AnchorNoSalaryShareCfgManagement') }}</span>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.anchorNoSalaryShareCfg.tip')"
          class="tip-alert"
          show-icon
          type="info"
      />

      <el-form :model="formData" class="cfg-form" label-width="180px">
        <el-form-item :label="t('pages.anchorNoSalaryShareCfg.anchorSocialSharePercent')">
          <el-input-number
              v-model="formData.anchorSocialSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
          />
        </el-form-item>

        <el-form-item :label="t('pages.anchorNoSalaryShareCfg.guildSocialSharePercent')">
          <el-input-number
              v-model="formData.guildSocialSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
          />
        </el-form-item>

        <el-form-item v-if="updatedAt" :label="t('pages.anchorNoSalaryShareCfg.lastUpdated')">
          <span>{{ updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="can('edit')" type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {anchorNoSalaryShareCfgApi} from '@/api/modules/anchor-no-salary-share-cfg'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('AnchorNoSalaryShareCfgManagement')
const loading = ref(false)
const updatedAt = ref('')
const formData = reactive({
  id: '0',
  anchorSocialSharePercent: 10,
  guildSocialSharePercent: 10,
})

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await anchorNoSalaryShareCfgApi.getCfg()
    const cfg = response.cfg
    formData.id = cfg?.id || '0'
    formData.anchorSocialSharePercent = cfg?.anchorSocialSharePercent ?? 10
    formData.guildSocialSharePercent = cfg?.guildSocialSharePercent ?? 10
    updatedAt.value = cfg?.updatedAt || ''
  } catch (error) {
    console.error('fetch anchor no salary share config failed:', error)
    ElMessage.error(t('pages.anchorNoSalaryShareCfg.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const validPercent = (value: number) => Number.isFinite(value) && value >= 0 && value <= 100

const handleSave = async () => {
  if (!validPercent(formData.anchorSocialSharePercent) || !validPercent(formData.guildSocialSharePercent)) {
    ElMessage.warning(t('pages.anchorNoSalaryShareCfg.percentRangeInvalid'))
    return
  }
  loading.value = true
  try {
    const response = await anchorNoSalaryShareCfgApi.saveCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      anchorSocialSharePercent: formData.anchorSocialSharePercent,
      guildSocialSharePercent: formData.guildSocialSharePercent,
    })
    if (!response.success) {
      ElMessage.error(t('pages.anchorNoSalaryShareCfg.saveFailed'))
      return
    }
    ElMessage.success(t('common.updateSuccess'))
    await fetchCfg()
  } catch (error) {
    console.error('save anchor no salary share config failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchCfg)
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.tip-alert {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 760px;
}
</style>

<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AppVersionCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" class="cfg-form" label-width="200px">
        <el-form-item :label="t('pages.appVersionCfg.versionQueryEnabled')">
          <el-switch
              v-model="formData.versionQueryEnabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
          <div class="form-tip">
            {{ t('pages.appVersionCfg.versionQueryEnabledTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.appVersionCfg.version')" prop="version">
          <el-input
              v-model="formData.version"
              maxlength="32"
              show-word-limit
              clearable
              :placeholder="t('pages.appVersionCfg.versionTip')"
              style="max-width: 360px"
          />
          <div class="form-tip">
            {{ t('pages.appVersionCfg.versionTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.appVersionCfg.downloadUrl')" prop="downloadUrl">
          <el-input
              v-model="formData.downloadUrl"
              maxlength="512"
              show-word-limit
              clearable
              :placeholder="t('pages.appVersionCfg.downloadUrlTip')"
              style="max-width: 560px"
          />
          <div class="form-tip">
            {{ t('pages.appVersionCfg.downloadUrlTip') }}
          </div>
        </el-form-item>

        <el-form-item :label="t('pages.appVersionCfg.updateDetails')">
          <div class="detail-panel">
            <div class="detail-toolbar">
              <el-button type="primary" @click="addUpdateDetail">{{ t('pages.appVersionCfg.addUpdateDetail') }}</el-button>
            </div>
            <div class="form-tip">
              {{ t('pages.appVersionCfg.updateDetailsTip') }}
            </div>
            <el-table :data="formData.updateDetails" empty-text=" " style="width: 100%">
              <template #empty>
                <el-empty :description="t('pages.appVersionCfg.updateDetailEmpty')"/>
              </template>
              <el-table-column :label="t('pages.appVersionCfg.updateDetailContent')" min-width="420">
                <template #default="{ row, $index }">
                  <el-input
                      v-model="row.content"
                      :placeholder="t('pages.appVersionCfg.updateDetailContent')"
                      maxlength="512"
                      show-word-limit
                      type="textarea"
                      :rows="2"
                      @blur="validateDetailRow($index)"
                  />
                </template>
              </el-table-column>
              <el-table-column :label="t('common.operation')" width="100">
                <template #default="{ $index }">
                  <el-button
                      link
                      type="danger"
                      @click="removeUpdateDetail($index)"
                  >
                    {{ t('common.remove') }}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.appVersionCfg.lastUpdated')">
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
import {ElMessage, type FormInstance} from 'element-plus'
import appVersionCfgApi from '@/api/modules/app-version-cfg'
import type {AppVersionCfg, AppVersionUpdateDetailItem} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  versionQueryEnabled: false,
  version: '',
  downloadUrl: '',
  updateDetails: [] as AppVersionUpdateDetailItem[],
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const rules = {}

const applyCfg = (cfg: AppVersionCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.versionQueryEnabled = !!cfg?.versionQueryEnabled
  formData.version = cfg?.version || ''
  formData.downloadUrl = cfg?.downloadUrl || ''
  formData.updateDetails = (cfg?.updateDetails || []).map((item, index) => ({
    content: item.content || '',
    sort: item.sort || index + 1,
  }))
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
}

const addUpdateDetail = () => {
  formData.updateDetails.push({
    content: '',
    sort: formData.updateDetails.length + 1,
  })
}

const removeUpdateDetail = (index: number) => {
  formData.updateDetails.splice(index, 1)
  formData.updateDetails.forEach((item, idx) => {
    item.sort = idx + 1
  })
}

const validateDetailRow = (index: number) => {
  const row = formData.updateDetails[index]
  if (!row) {
    return
  }
  row.content = String(row.content || '').trim()
}

const buildSaveDetails = () => {
  return formData.updateDetails
      .map((item, index) => ({
        content: String(item.content || '').trim(),
        sort: index + 1,
      }))
      .filter((item) => item.content)
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await appVersionCfgApi.getAppVersionCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch app version cfg failed:', error)
    ElMessage.error(t('pages.appVersionCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  try {
    await formRef.value.validate()
    const response = await appVersionCfgApi.saveAppVersionCfg({
      id: formData.id === '0' ? 0 : Number(formData.id),
      versionQueryEnabled: formData.versionQueryEnabled,
      version: formData.version.trim(),
      downloadUrl: formData.downloadUrl.trim(),
      updateDetails: buildSaveDetails(),
    })
    if (response?.success) {
      ElMessage.success(t('pages.appVersionCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.appVersionCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save app version cfg failed:', error)
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
  max-width: 960px;
}

.form-tip {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}

.detail-panel {
  width: 100%;
}

.detail-toolbar {
  margin-bottom: 8px;
}
</style>

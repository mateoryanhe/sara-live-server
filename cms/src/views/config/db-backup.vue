<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.DbBackupCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.dbBackup.noticeTitle')"
          type="info"
      >
        <p>{{ t('pages.dbBackup.noticeLine1') }}</p>
        <p>{{ t('pages.dbBackup.noticeLine2') }}</p>
        <p>{{ t('pages.dbBackup.noticeLine3') }}</p>
        <p>{{ t('pages.dbBackup.noticeLine4') }}</p>
      </el-alert>

      <el-form class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.dbBackup.enabled')">
          <el-switch
              v-model="formData.enabled"
              :active-text="t('common.open')"
              :inactive-text="t('common.close')"
          />
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.retainDays')">
          <el-input-number v-model="formData.retainDays" :min="1" :max="365"/>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.storagePrefix')">
          <span>{{ formData.storagePrefix || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.sourceDatabase')">
          <span>{{ formData.sourceDatabase || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.s3Enabled')">
          <el-tag :type="formData.s3Enabled ? 'success' : 'danger'">
            {{ formData.s3Enabled ? t('pages.dbBackup.s3On') : t('pages.dbBackup.s3Off') }}
          </el-tag>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.lastSuccessAt')">
          <span>{{ formData.lastSuccessAt || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.lastObjectKey')">
          <span class="mono">{{ formData.lastObjectKey || '-' }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.lastFileSize')">
          <span>{{ formatSize(formData.lastFileSize) }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.lastError')">
          <span class="error-text">{{ formData.lastError || '-' }}</span>
        </el-form-item>

        <el-form-item>
          <el-button v-if="can('save')" type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="refreshAll">{{ t('common.refresh') }}</el-button>
          <el-button
              v-if="can('runNow')"
              :loading="running"
              type="warning"
              @click="handleRunNow"
          >
            {{ t('pages.dbBackup.runNow') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-loading="filesLoading" class="restore-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('pages.dbBackup.restoreTitle') }}</span>
          <el-button v-if="can('list')" link type="primary" @click="fetchFiles">
            {{ t('pages.dbBackup.refreshFiles') }}
          </el-button>
        </div>
      </template>

      <el-form class="cfg-form" label-width="160px">
        <el-form-item :label="t('pages.dbBackup.restoreFile')">
          <el-select
              v-model="restoreForm.objectKey"
              clearable
              filterable
              style="width: 100%; max-width: 640px"
              :placeholder="t('pages.dbBackup.restoreFilePlaceholder')"
          >
            <el-option
                v-for="item in fileList"
                :key="item.objectKey"
                :label="`${item.fileName} (${formatSize(item.size)}) ${item.lastModified}`"
                :value="item.objectKey"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.dbBackup.targetDatabase')">
          <el-input
              v-model="restoreForm.targetDatabase"
              clearable
              style="max-width: 360px"
              :placeholder="t('pages.dbBackup.targetDatabasePlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button
              v-if="can('restore')"
              :disabled="!restoreForm.objectKey || !restoreForm.targetDatabase.trim()"
              :loading="restoring"
              type="danger"
              @click="handleRestore"
          >
            {{ t('pages.dbBackup.restoreTitle') }}
          </el-button>
        </el-form-item>
        <el-empty v-if="!filesLoading && fileList.length === 0" :description="t('pages.dbBackup.emptyFiles')"/>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {dbBackupApi} from '@/api/modules/db-backup'
import type {DbBackupCfg, DbBackupFileItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('DbBackupCfgManagement')

const loading = ref(false)
const filesLoading = ref(false)
const running = ref(false)
const restoring = ref(false)
const fileList = ref<DbBackupFileItem[]>([])

const formData = reactive({
  id: '0',
  enabled: false,
  retainDays: 1,
  storagePrefix: '',
  sourceDatabase: '',
  s3Enabled: false,
  lastSuccessAt: '',
  lastError: '',
  lastObjectKey: '',
  lastFileSize: 0,
})

const restoreForm = reactive({
  objectKey: '',
  targetDatabase: '',
})

const formatSize = (size?: number) => {
  const n = Number(size || 0)
  if (n <= 0) return '-'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(2)} MB`
}

const applyCfg = (cfg: DbBackupCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.enabled = !!cfg?.enabled
  formData.retainDays = cfg?.retainDays && cfg.retainDays > 0 ? cfg.retainDays : 1
  formData.storagePrefix = cfg?.storagePrefix || ''
  formData.sourceDatabase = cfg?.sourceDatabase || ''
  formData.s3Enabled = !!cfg?.s3Enabled
  formData.lastSuccessAt = cfg?.lastSuccessAt || ''
  formData.lastError = cfg?.lastError || ''
  formData.lastObjectKey = cfg?.lastObjectKey || ''
  formData.lastFileSize = cfg?.lastFileSize || 0
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await dbBackupApi.getCfg()
    applyCfg(response?.cfg)
  } catch (error) {
    console.error('fetch db backup cfg failed:', error)
    ElMessage.error(t('pages.dbBackup.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const fetchFiles = async () => {
  filesLoading.value = true
  try {
    const response = await dbBackupApi.listFiles()
    fileList.value = response?.list || []
  } catch (error) {
    console.error('list db backup files failed:', error)
    ElMessage.error(t('pages.dbBackup.listFailed'))
  } finally {
    filesLoading.value = false
  }
}

const refreshAll = async () => {
  await fetchCfg()
  await fetchFiles()
}

const handleSave = async () => {
  loading.value = true
  try {
    const response = await dbBackupApi.saveCfg({
      id: Number(formData.id) || 0,
      enabled: formData.enabled,
      retainDays: formData.retainDays,
    })
    if (response?.success) {
      ElMessage.success(t('common.updateSuccess'))
      await refreshAll()
    } else {
      ElMessage.error(t('pages.dbBackup.saveFailed'))
    }
  } catch (error) {
    console.error('save db backup cfg failed:', error)
    ElMessage.error(t('pages.dbBackup.saveFailed'))
  } finally {
    loading.value = false
  }
}

const handleRunNow = async () => {
  try {
    await ElMessageBox.confirm(
        t('pages.dbBackup.noticeLine1'),
        t('pages.dbBackup.runNow'),
        {confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning'},
    )
  } catch {
    return
  }
  running.value = true
  try {
    const response = await dbBackupApi.runNow()
    if (response?.success) {
      ElMessage.success(response.message || t('pages.dbBackup.runNowSuccess'))
      await refreshAll()
    } else {
      ElMessage.error(t('pages.dbBackup.runNowFailed'))
    }
  } catch (error) {
    console.error('run db backup failed:', error)
    ElMessage.error(t('pages.dbBackup.runNowFailed'))
    await fetchCfg()
  } finally {
    running.value = false
  }
}

const handleRestore = async () => {
  const objectKey = restoreForm.objectKey
  const targetDatabase = restoreForm.targetDatabase.trim()
  if (!objectKey || !targetDatabase) {
    ElMessage.warning(t('pages.dbBackup.targetDatabaseRequired'))
    return
  }
  const fileName = fileList.value.find((x) => x.objectKey === objectKey)?.fileName || objectKey
  try {
    if (formData.sourceDatabase && targetDatabase === formData.sourceDatabase) {
      await ElMessageBox.confirm(
          t('pages.dbBackup.restoreSameDbWarn'),
          t('pages.dbBackup.restoreTitle'),
          {confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'error'},
      )
    }
    await ElMessageBox.confirm(
        t('pages.dbBackup.restoreConfirm', {file: fileName, db: targetDatabase}),
        t('pages.dbBackup.restoreTitle'),
        {confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning'},
    )
  } catch {
    return
  }
  restoring.value = true
  try {
    const response = await dbBackupApi.restore({objectKey, targetDatabase})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.dbBackup.restoreSuccess'))
    } else {
      ElMessage.error(t('pages.dbBackup.restoreFailed'))
    }
  } catch (error) {
    console.error('restore db backup failed:', error)
    ElMessage.error(t('pages.dbBackup.restoreFailed'))
  } finally {
    restoring.value = false
  }
}

onMounted(() => {
  refreshAll()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tip-alert {
  margin-bottom: 20px;
}

.tip-alert p {
  margin: 0 0 4px;
}

.cfg-form {
  max-width: 900px;
}

.mono {
  word-break: break-all;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 13px;
}

.error-text {
  color: var(--el-color-danger);
  word-break: break-all;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.restore-card {
  margin-top: 0;
}
</style>

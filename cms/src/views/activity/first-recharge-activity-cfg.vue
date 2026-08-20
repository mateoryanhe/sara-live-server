<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.FirstRechargeActivityManagement') }}</span>
          <el-button
              v-if="can('sync')"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" class="cfg-form" label-width="160px">
        <el-tabs v-model="mainTab">
          <el-tab-pane :label="t('pages.firstRechargeActivityCfg.tabBasic')" name="basic">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.enabled')">
              <el-switch
                  v-model="formData.enabled"
                  :active-text="t('common.open')"
                  :inactive-text="t('common.close')"
              />
              <div class="form-tip">{{ t('pages.firstRechargeActivityCfg.enabledTip') }}</div>
            </el-form-item>

            <el-form-item :label="t('pages.firstRechargeActivityCfg.icon')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="mainIconUploading"
                    :http-request="doMainIconUpload"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="mainIconPreviewUrl" :src="mainIconPreviewUrl" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon">
                      <Plus/>
                    </el-icon>
                  </div>
                </el-upload>
                <el-button
                    v-if="mainIconPreviewUrl || formData.icon"
                    link
                    type="danger"
                    @click="clearMainIcon"
                >
                  {{ t('common.remove') }}
                </el-button>
              </div>
              <div class="form-tip">{{ t('pages.firstRechargeActivityCfg.iconTip') }}</div>
            </el-form-item>

            <el-tabs v-model="basicLangTab" class="inner-lang-tabs">
              <el-tab-pane label="English" name="en">
                <el-form-item :label="t('pages.firstRechargeActivityCfg.titleEn')">
                  <el-input v-model="formData.titleEn" clearable maxlength="128" show-word-limit/>
                </el-form-item>
                <el-form-item :label="t('pages.firstRechargeActivityCfg.rechargeBtnTextEn')">
                  <el-input v-model="formData.rechargeBtnTextEn" clearable maxlength="64" show-word-limit/>
                </el-form-item>
              </el-tab-pane>
              <el-tab-pane label="Español" name="es">
                <el-form-item :label="t('pages.firstRechargeActivityCfg.titleEs')">
                  <el-input v-model="formData.titleEs" clearable maxlength="128" show-word-limit/>
                </el-form-item>
                <el-form-item :label="t('pages.firstRechargeActivityCfg.rechargeBtnTextEs')">
                  <el-input v-model="formData.rechargeBtnTextEs" clearable maxlength="64" show-word-limit/>
                </el-form-item>
              </el-tab-pane>
              <el-tab-pane label="Português" name="pt">
                <el-form-item :label="t('pages.firstRechargeActivityCfg.titlePt')">
                  <el-input v-model="formData.titlePt" clearable maxlength="128" show-word-limit/>
                </el-form-item>
                <el-form-item :label="t('pages.firstRechargeActivityCfg.rechargeBtnTextPt')">
                  <el-input v-model="formData.rechargeBtnTextPt" clearable maxlength="64" show-word-limit/>
                </el-form-item>
              </el-tab-pane>
              <el-tab-pane label="हिन्दी" name="hi">
                <el-form-item :label="t('pages.firstRechargeActivityCfg.titleHi')">
                  <el-input v-model="formData.titleHi" clearable maxlength="128" show-word-limit/>
                </el-form-item>
                <el-form-item :label="t('pages.firstRechargeActivityCfg.rechargeBtnTextHi')">
                  <el-input v-model="formData.rechargeBtnTextHi" clearable maxlength="64" show-word-limit/>
                </el-form-item>
              </el-tab-pane>
              <el-tab-pane label="Bahasa" name="id">
                <el-form-item :label="t('pages.firstRechargeActivityCfg.titleId')">
                  <el-input v-model="formData.titleId" clearable maxlength="128" show-word-limit/>
                </el-form-item>
                <el-form-item :label="t('pages.firstRechargeActivityCfg.rechargeBtnTextId')">
                  <el-input v-model="formData.rechargeBtnTextId" clearable maxlength="64" show-word-limit/>
                </el-form-item>
              </el-tab-pane>
            </el-tabs>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.firstRechargeActivityCfg.tabPrivileges')" name="privileges">
            <div class="privilege-toolbar">
              <el-button type="primary" @click="openPrivilegeDialog()">{{ t('pages.firstRechargeActivityCfg.addPrivilege') }}</el-button>
            </div>

            <el-table :data="privilegeList" empty-text=" " style="width: 100%">
              <template #empty>
                <el-empty :description="t('pages.firstRechargeActivityCfg.privilegeEmpty')"/>
              </template>
              <el-table-column :label="t('pages.firstRechargeActivityCfg.privilegeIcon')" width="100">
                <template #default="{ row }">
                  <el-image
                      v-if="row.icon"
                      :preview-src-list="[row.icon]"
                      :src="row.icon"
                      fit="cover"
                      preview-teleported
                      style="width: 48px; height: 48px"
                  />
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column :label="t('pages.firstRechargeActivityCfg.privilegeDescPreview')" min-width="240" prop="descEn" show-overflow-tooltip/>
              <el-table-column fixed="right" :label="t('common.actions')" width="160">
                <template #default="{ row, $index }">
                  <el-button link type="primary" @click="openPrivilegeDialog(row, $index)">{{ t('common.edit') }}</el-button>
                  <el-button link type="danger" @click="removePrivilege($index)">{{ t('common.delete') }}</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.firstRechargeActivityCfg.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-dialog
        v-model="privilegeDialogVisible"
        :title="privilegeDialogTitle"
        destroy-on-close
        width="640px"
        @closed="resetPrivilegeDialog"
    >
      <el-form label-width="140px">
        <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeIcon')">
          <div class="icon-upload-wrap">
            <el-upload
                :before-upload="beforeImageUpload"
                :disabled="privilegeIconUploading"
                :http-request="doPrivilegeIconUpload"
                :show-file-list="false"
                accept="image/*"
                action="#"
                class="icon-uploader"
            >
              <img v-if="privilegeIconPreviewUrl" :src="privilegeIconPreviewUrl" alt="icon" class="icon-preview"/>
              <div v-else class="icon-uploader-placeholder">
                <el-icon class="icon-uploader-icon">
                  <Plus/>
                </el-icon>
              </div>
            </el-upload>
            <el-button
                v-if="privilegeIconPreviewUrl || privilegeForm.iconName"
                link
                type="danger"
                @click="clearPrivilegeIcon"
            >
              {{ t('common.remove') }}
            </el-button>
          </div>
        </el-form-item>

        <el-tabs v-model="privilegeLangTab">
          <el-tab-pane label="English" name="en">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeDescEn')">
              <el-input v-model="privilegeForm.descEn" clearable maxlength="256" show-word-limit type="textarea" :rows="3"/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Español" name="es">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeDescEs')">
              <el-input v-model="privilegeForm.descEs" clearable maxlength="256" show-word-limit type="textarea" :rows="3"/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Português" name="pt">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeDescPt')">
              <el-input v-model="privilegeForm.descPt" clearable maxlength="256" show-word-limit type="textarea" :rows="3"/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="हिन्दी" name="hi">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeDescHi')">
              <el-input v-model="privilegeForm.descHi" clearable maxlength="256" show-word-limit type="textarea" :rows="3"/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Bahasa" name="id">
            <el-form-item :label="t('pages.firstRechargeActivityCfg.privilegeDescId')">
              <el-input v-model="privilegeForm.descId" clearable maxlength="256" show-word-limit type="textarea" :rows="3"/>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <template #footer>
        <el-button @click="privilegeDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmPrivilegeDialog">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import firstRechargeActivityApi from '@/api/modules/first-recharge-activity'
import dataSyncApi from '@/api/modules/data-sync'
import uploadApi from '@/api/modules/upload'
import type {FirstRechargeActivityCfg, FirstRechargePrivilegeItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('FirstRechargeActivityManagement')
const loading = ref(false)
const syncing = ref(false)
const mainIconUploading = ref(false)
const privilegeIconUploading = ref(false)
const mainIconPreviewUrl = ref('')
const privilegeIconPreviewUrl = ref('')
const mainTab = ref('basic')
const basicLangTab = ref('en')
const privilegeLangTab = ref('en')
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '0',
  enabled: false,
  icon: '',
  titleEn: '',
  titleEs: '',
  titlePt: '',
  titleHi: '',
  titleId: '',
  rechargeBtnTextEn: '',
  rechargeBtnTextEs: '',
  rechargeBtnTextPt: '',
  rechargeBtnTextHi: '',
  rechargeBtnTextId: '',
})

const privilegeList = ref<FirstRechargePrivilegeItem[]>([])
const privilegeDialogVisible = ref(false)
const privilegeEditIndex = ref<number | null>(null)
const privilegeForm = reactive({
  iconName: '',
  descEn: '',
  descEs: '',
  descPt: '',
  descHi: '',
  descId: '',
})

const metaInfo = reactive({
  createdAt: '',
  updatedAt: '',
})

const privilegeDialogTitle = computed(() => (
    privilegeEditIndex.value == null
        ? t('pages.firstRechargeActivityCfg.addPrivilege')
        : t('pages.firstRechargeActivityCfg.editPrivilege')
))

const clonePrivilegeList = (list?: FirstRechargePrivilegeItem[]) => {
  return (list || []).map(item => ({
    icon: item.icon || '',
    iconName: item.iconName || item.icon || '',
    descEn: item.descEn || '',
    descEs: item.descEs || '',
    descPt: item.descPt || '',
    descHi: item.descHi || '',
    descId: item.descId || '',
  }))
}

const applyCfg = (cfg: FirstRechargeActivityCfg | null | undefined) => {
  formData.id = cfg?.id || '0'
  formData.enabled = !!cfg?.enabled
  formData.icon = cfg?.iconName || cfg?.icon || ''
  formData.titleEn = cfg?.titleEn || ''
  formData.titleEs = cfg?.titleEs || ''
  formData.titlePt = cfg?.titlePt || ''
  formData.titleHi = cfg?.titleHi || ''
  formData.titleId = cfg?.titleId || ''
  formData.rechargeBtnTextEn = cfg?.rechargeBtnTextEn || ''
  formData.rechargeBtnTextEs = cfg?.rechargeBtnTextEs || ''
  formData.rechargeBtnTextPt = cfg?.rechargeBtnTextPt || ''
  formData.rechargeBtnTextHi = cfg?.rechargeBtnTextHi || ''
  formData.rechargeBtnTextId = cfg?.rechargeBtnTextId || ''
  privilegeList.value = clonePrivilegeList(cfg?.privileges)
  metaInfo.createdAt = cfg?.createdAt || ''
  metaInfo.updatedAt = cfg?.updatedAt || ''
  mainIconPreviewUrl.value = cfg?.icon || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const response = await firstRechargeActivityApi.getFirstRechargeActivityCfg()
    applyCfg(response.cfg)
  } catch (error) {
    console.error('fetch first recharge activity cfg failed:', error)
    ElMessage.error(t('pages.firstRechargeActivityCfg.fetchCfgFailed'))
  } finally {
    loading.value = false
  }
}

const beforeImageUpload = (file: File) => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.firstRechargeActivityCfg.imageTypeInvalid'))
    return false
  }
  return true
}

const doMainIconUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  mainIconUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    formData.icon = res.fileName
    mainIconPreviewUrl.value = res.fileUrl || ''
  } catch (error) {
    console.error('upload icon failed:', error)
    ElMessage.error(t('pages.firstRechargeActivityCfg.uploadIconFailed'))
  } finally {
    mainIconUploading.value = false
  }
}

const clearMainIcon = () => {
  formData.icon = ''
  mainIconPreviewUrl.value = ''
}

const doPrivilegeIconUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  privilegeIconUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    privilegeForm.iconName = res.fileName
    privilegeIconPreviewUrl.value = res.fileUrl || ''
  } catch (error) {
    console.error('upload privilege icon failed:', error)
    ElMessage.error(t('pages.firstRechargeActivityCfg.uploadIconFailed'))
  } finally {
    privilegeIconUploading.value = false
  }
}

const clearPrivilegeIcon = () => {
  privilegeForm.iconName = ''
  privilegeIconPreviewUrl.value = ''
}

const resetPrivilegeDialog = () => {
  privilegeEditIndex.value = null
  privilegeForm.iconName = ''
  privilegeForm.descEn = ''
  privilegeForm.descEs = ''
  privilegeForm.descPt = ''
  privilegeForm.descHi = ''
  privilegeForm.descId = ''
  privilegeIconPreviewUrl.value = ''
  privilegeLangTab.value = 'en'
}

const openPrivilegeDialog = (row?: FirstRechargePrivilegeItem, index?: number) => {
  resetPrivilegeDialog()
  if (row != null && index != null) {
    privilegeEditIndex.value = index
    privilegeForm.iconName = row.iconName || row.icon || ''
    privilegeForm.descEn = row.descEn || ''
    privilegeForm.descEs = row.descEs || ''
    privilegeForm.descPt = row.descPt || ''
    privilegeForm.descHi = row.descHi || ''
    privilegeForm.descId = row.descId || ''
    privilegeIconPreviewUrl.value = row.icon || ''
  }
  privilegeDialogVisible.value = true
}

const confirmPrivilegeDialog = () => {
  const item: FirstRechargePrivilegeItem = {
    iconName: privilegeForm.iconName.trim(),
    icon: privilegeIconPreviewUrl.value,
    descEn: privilegeForm.descEn.trim(),
    descEs: privilegeForm.descEs.trim(),
    descPt: privilegeForm.descPt.trim(),
    descHi: privilegeForm.descHi.trim(),
    descId: privilegeForm.descId.trim(),
  }
  if (privilegeEditIndex.value == null) {
    privilegeList.value.push(item)
  } else {
    privilegeList.value.splice(privilegeEditIndex.value, 1, item)
  }
  privilegeDialogVisible.value = false
}

const removePrivilege = async (index: number) => {
  try {
    await ElMessageBox.confirm(
        t('pages.firstRechargeActivityCfg.deletePrivilegeConfirm'),
        t('common.confirmDelete'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  privilegeList.value.splice(index, 1)
}

const buildSavePayload = () => ({
  id: formData.id === '0' ? 0 : Number(formData.id),
  enabled: formData.enabled,
  icon: formData.icon.trim(),
  titleEn: formData.titleEn.trim(),
  titleEs: formData.titleEs.trim(),
  titlePt: formData.titlePt.trim(),
  titleHi: formData.titleHi.trim(),
  titleId: formData.titleId.trim(),
  rechargeBtnTextEn: formData.rechargeBtnTextEn.trim(),
  rechargeBtnTextEs: formData.rechargeBtnTextEs.trim(),
  rechargeBtnTextPt: formData.rechargeBtnTextPt.trim(),
  rechargeBtnTextHi: formData.rechargeBtnTextHi.trim(),
  rechargeBtnTextId: formData.rechargeBtnTextId.trim(),
  privileges: privilegeList.value.map(item => ({
    iconName: (item.iconName || item.icon || '').trim(),
    descEn: (item.descEn || '').trim(),
    descEs: (item.descEs || '').trim(),
    descPt: (item.descPt || '').trim(),
    descHi: (item.descHi || '').trim(),
    descId: (item.descId || '').trim(),
  })),
})

const handleSave = async () => {
  try {
    const response = await firstRechargeActivityApi.saveFirstRechargeActivityCfg(buildSavePayload())
    if (response?.success) {
      ElMessage.success(t('pages.firstRechargeActivityCfg.saveSuccess'))
      if (response.id) {
        formData.id = response.id
      }
      await fetchCfg()
    } else {
      ElMessage.error(t('pages.firstRechargeActivityCfg.saveFailed'))
    }
  } catch (error) {
    console.error('save first recharge activity cfg failed:', error)
  }
}

const handleSyncData = async () => {
  if (!formData.id || formData.id === '0') {
    ElMessage.warning(t('pages.firstRechargeActivityCfg.syncNeedSave'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.firstRechargeActivityCfg.syncConfirm'),
        t('common.syncData'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  syncing.value = true
  try {
    const response = await dataSyncApi.syncFirstRechargeActivityCfg()
    if (response?.success) {
      ElMessage.success(response.message || t('pages.firstRechargeActivityCfg.syncSuccess'))
    } else {
      ElMessage.error(t('pages.firstRechargeActivityCfg.syncFailed'))
    }
  } catch (error) {
    console.error('sync first recharge activity cfg failed:', error)
    ElMessage.error(t('pages.firstRechargeActivityCfg.syncFailed'))
  } finally {
    syncing.value = false
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
  font-size: 16px;
  font-weight: bold;
}

.cfg-form {
  max-width: 900px;
}

.inner-lang-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}

.privilege-toolbar {
  margin-bottom: 16px;
}

.form-tip {
  width: 100%;
  margin-top: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}

.icon-upload-wrap {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.icon-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
}

.icon-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.icon-uploader-placeholder {
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-light);
}

.icon-uploader-icon {
  font-size: 24px;
  color: var(--el-text-color-secondary);
}

.icon-preview {
  width: 72px;
  height: 72px;
  object-fit: cover;
  display: block;
}
</style>

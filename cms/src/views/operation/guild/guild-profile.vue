<template>
  <div class="page-container">
    <el-card>
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

      <div class="content">
        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.guildProfile.guildName')" min-width="140" prop="name"/>
          <el-table-column :label="t('pages.guildProfile.description')" min-width="180" prop="description" show-overflow-tooltip/>
          <el-table-column :label="t('pages.guildProfile.lastUpdated')" prop="updatedAt" width="170"/>
          <el-table-column :label="t('common.actions')" width="100" fixed="right">
            <template #default="{ row }">
              <el-button v-if="can('save')" size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-empty v-if="!loading && tableData.length === 0" :description="t('pages.guildProfile.noGuild')"/>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="t('pages.guildProfile.editGuild')" destroy-on-close width="640px">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.guildProfile.guildId')">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.guildProfile.guildName')" prop="name">
          <el-input v-model="formData.name" maxlength="32" :placeholder="t('pages.guildProfile.guildNamePlaceholder')" show-word-limit/>
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
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="saving" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
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
const dialogVisible = ref(false)
const tableData = ref<MyGuildProfile[]>([])
const formRef = ref<FormInstance>()

const formData = reactive({
  id: '',
  name: '',
  description: '',
})

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.guildProfile.nameRequired'), trigger: 'blur'},
    {min: 2, max: 32, message: t('pages.guildProfile.nameLength'), trigger: 'blur'},
  ],
  description: [
    {max: 255, message: t('pages.guildProfile.descriptionMaxLength'), trigger: 'blur'},
  ],
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getMyGuildProfile()
    tableData.value = response?.list ?? []
  } catch (error) {
    console.error('fetch guild profile list failed:', error)
    tableData.value = []
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: MyGuildProfile) => {
  formData.id = row.id
  formData.name = row.name || ''
  formData.description = row.description || ''
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!formRef.value || !formData.id) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      const response = await guildApi.updateMyGuildProfile({
        id: formData.id,
        name: formData.name.trim(),
        description: formData.description.trim(),
      })
      if (response?.success) {
        ElMessage.success(t('common.saveConfig'))
        dialogVisible.value = false
        await fetchList()
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
  fetchList()
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

.content {
  margin-top: 4px;
}
</style>

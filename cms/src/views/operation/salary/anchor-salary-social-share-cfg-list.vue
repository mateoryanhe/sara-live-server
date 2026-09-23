<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AnchorSalarySocialShareCfgManagement') }}</span>
          <el-button v-if="can('create')" type="primary" @click="handleAdd">
            {{ t('pages.anchorSalarySocialShareCfgList.addTier') }}
          </el-button>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.anchorSalarySocialShareCfgList.dynamicTierHint')"
          class="tier-hint"
          type="info"
      />

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column
            :label="t('pages.anchorSalarySocialShareCfgList.level')"
            align="center"
            prop="level"
            width="120"
        />
        <el-table-column
            :label="t('pages.anchorSalarySocialShareCfgList.socialTotalDiamondRevenue')"
            align="right"
            min-width="200"
            prop="socialTotalDiamondRevenue"
        >
          <template #default="{row}">
            <span class="money-amount">{{ formatWalletBalance(row.socialTotalDiamondRevenue) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            :label="t('pages.anchorSalarySocialShareCfgList.anchorSocialSharePercent')"
            align="right"
            min-width="180"
            prop="anchorSocialSharePercent"
        >
          <template #default="{row}">{{ formatPercent(row.anchorSocialSharePercent) }}</template>
        </el-table-column>
        <el-table-column
            :label="t('pages.anchorSalarySocialShareCfgList.guildSocialSharePercent')"
            align="right"
            min-width="180"
            prop="guildSocialSharePercent"
        >
          <template #default="{row}">{{ formatPercent(row.guildSocialSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.updatedAt')" min-width="180" prop="updatedAt">
          <template #default="{row}">{{ row.updatedAt || '-' }}</template>
        </el-table-column>
        <el-table-column v-if="can('edit') || can('delete')" :label="t('common.actions')" fixed="right" width="150">
          <template #default="{row}">
            <el-button v-if="can('edit')" link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button v-if="can('delete')" link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="180px">
        <el-form-item :label="t('pages.anchorSalarySocialShareCfgList.level')" prop="level">
          <el-input-number v-model="formData.level" :min="1" :precision="0" controls-position="right" style="width: 100%"/>
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorSalarySocialShareCfgList.socialTotalDiamondRevenue')"
            prop="socialTotalDiamondRevenue"
        >
          <el-input-number
              v-model="formData.socialTotalDiamondRevenue"
              :min="0"
              :precision="4"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorSalarySocialShareCfgList.anchorSocialSharePercent')"
            prop="anchorSocialSharePercent"
        >
          <el-input-number
              v-model="formData.anchorSocialSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorSalarySocialShareCfgList.guildSocialSharePercent')"
            prop="guildSocialSharePercent"
        >
          <el-input-number
              v-model="formData.guildSocialSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :disabled="saving" @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="saving" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {anchorSalarySocialShareCfgApi} from '@/api/modules/anchor-salary-social-share-cfg'
import type {AnchorSalarySocialShareCfg} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'

const {t} = useI18n()
const {can} = usePagePermission('AnchorSalarySocialShareCfgManagement')

const loading = ref(false)
const saving = ref(false)
const tableData = ref<AnchorSalarySocialShareCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = reactive({
  id: '',
  level: 0,
  socialTotalDiamondRevenue: 0,
  anchorSocialSharePercent: 0,
  guildSocialSharePercent: 0,
})

const formRules = computed<FormRules>(() => ({
  level: [
    {required: true, type: 'number', min: 1, message: t('pages.anchorSalarySocialShareCfgList.levelRequired'), trigger: 'change'},
  ],
  socialTotalDiamondRevenue: [
    {required: true, message: t('pages.anchorSalarySocialShareCfgList.socialTotalRevenueRequired'), trigger: 'change'},
  ],
  anchorSocialSharePercent: [
    {required: true, message: t('pages.anchorSalarySocialShareCfgList.anchorSocialSharePercentRequired'), trigger: 'change'},
    {
      validator: (_rule, value, callback) => {
        const percent = Number(value)
        if (!Number.isFinite(percent) || percent < 0 || percent > 100) {
          callback(new Error(t('pages.anchorSalarySocialShareCfgList.sharePercentInvalid')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
  guildSocialSharePercent: [
    {required: true, message: t('pages.anchorSalarySocialShareCfgList.guildSocialSharePercentRequired'), trigger: 'change'},
    {
      validator: (_rule, value, callback) => {
        const percent = Number(value)
        if (!Number.isFinite(percent) || percent < 0 || percent > 100) {
          callback(new Error(t('pages.anchorSalarySocialShareCfgList.sharePercentInvalid')))
          return
        }
        callback()
      },
      trigger: 'change',
    },
  ],
}))

const formatPercent = (value: number | null | undefined) => `${Number(value || 0).toFixed(2)}%`

const fetchList = async () => {
  loading.value = true
  try {
    const response = await anchorSalarySocialShareCfgApi.getList({
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = Number(response.total || 0)
  } catch (error) {
    console.error('fetch anchor salary social share config failed:', error)
    ElMessage.error(t('pages.anchorSalarySocialShareCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.id = ''
  formData.level = 1
  formData.socialTotalDiamondRevenue = 0
  formData.anchorSocialSharePercent = 0
  formData.guildSocialSharePercent = 0
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = t('pages.anchorSalarySocialShareCfgList.addTier')
  dialogVisible.value = true
}

const handleEdit = (row: AnchorSalarySocialShareCfg) => {
  formData.id = String(row.id || '')
  formData.level = Number(row.level || 0)
  formData.socialTotalDiamondRevenue = Number(row.socialTotalDiamondRevenue || 0)
  formData.anchorSocialSharePercent = Number(row.anchorSocialSharePercent || 0)
  formData.guildSocialSharePercent = Number(row.guildSocialSharePercent || 0)
  dialogTitle.value = t('pages.anchorSalarySocialShareCfgList.editTier')
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!await formRef.value?.validate()) return
  saving.value = true
  try {
    const payload = {
      level: formData.level,
      socialTotalDiamondRevenue: formData.socialTotalDiamondRevenue,
      anchorSocialSharePercent: formData.anchorSocialSharePercent,
      guildSocialSharePercent: formData.guildSocialSharePercent,
    }
    if (formData.id) {
      await anchorSalarySocialShareCfgApi.update({id: formData.id, ...payload})
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await anchorSalarySocialShareCfgApi.create(payload)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    await fetchList()
  } catch (error) {
    console.error('update anchor salary social share config failed:', error)
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row: AnchorSalarySocialShareCfg) => {
  try {
    await ElMessageBox.confirm(
      t('pages.anchorSalarySocialShareCfgList.deleteConfirm', {level: row.level}),
      t('common.confirmDelete'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await anchorSalarySocialShareCfgApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    if (tableData.value.length === 1 && currentPage.value > 1) {
      currentPage.value -= 1
    }
    await fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete anchor salary social share config failed:', error)
    }
  }
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

onMounted(fetchList)
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

.tier-hint {
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>

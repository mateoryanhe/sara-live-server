<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ pageTitle }}</span>
          <el-button v-if="can('create')" type="primary" @click="handleAdd">
            {{ t('pages.anchorGameShareCfgList.addTier') }}
          </el-button>
        </div>
      </template>

      <el-alert
          :closable="false"
          :title="t('pages.anchorGameShareCfgList.dynamicTierHint')"
          class="tier-hint"
          type="info"
      />

      <el-table v-loading="loading" :data="tableData" style="width: 100%">
        <el-table-column :label="t('pages.anchorGameShareCfgList.level')" align="center" prop="level" width="100"/>
        <el-table-column
            :label="t('pages.anchorGameShareCfgList.gameTotalGoldRevenue')"
            align="right"
            min-width="190"
            prop="gameTotalGoldRevenue"
        >
          <template #default="{row}">
            <span class="money-amount">{{ formatWalletBalance(row.gameTotalGoldRevenue) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            :label="t('pages.anchorGameShareCfgList.anchorGameSharePercent')"
            align="right"
            min-width="190"
            prop="anchorGameSharePercent"
        >
          <template #default="{row}">{{ formatPercent(row.anchorGameSharePercent) }}</template>
        </el-table-column>
        <el-table-column
            :label="t('pages.anchorGameShareCfgList.guildGameSharePercent')"
            align="right"
            min-width="190"
            prop="guildGameSharePercent"
        >
          <template #default="{row}">{{ formatPercent(row.guildGameSharePercent) }}</template>
        </el-table-column>
        <el-table-column :label="t('common.updatedAt')" min-width="180" prop="updatedAt">
          <template #default="{row}">{{ row.updatedAt || '-' }}</template>
        </el-table-column>
        <el-table-column
            v-if="can('edit') || can('delete')"
            :label="t('common.actions')"
            fixed="right"
            width="150"
        >
          <template #default="{row}">
            <el-button v-if="can('edit')" link type="primary" @click="handleEdit(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button v-if="can('delete')" link type="danger" @click="handleDelete(row)">
              {{ t('common.delete') }}
            </el-button>
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
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="190px">
        <el-form-item :label="t('pages.anchorGameShareCfgList.level')" prop="level">
          <el-input-number
              v-model="formData.level"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorGameShareCfgList.gameTotalGoldRevenue')"
            prop="gameTotalGoldRevenue"
        >
          <el-input-number
              v-model="formData.gameTotalGoldRevenue"
              :min="0"
              :precision="4"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorGameShareCfgList.anchorGameSharePercent')"
            prop="anchorGameSharePercent"
        >
          <el-input-number
              v-model="formData.anchorGameSharePercent"
              :max="100"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item
            :label="t('pages.anchorGameShareCfgList.guildGameSharePercent')"
            prop="guildGameSharePercent"
        >
          <el-input-number
              v-model="formData.guildGameSharePercent"
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {anchorGameShareCfgApi} from '@/api/modules/anchor-game-share-cfg'
import type {AnchorGameShareCfg} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import {formatWalletBalance} from '@/utils/number-format'

const props = defineProps<{salaryType: 1 | 2}>()
const {t} = useI18n()
const {can} = usePagePermission()
const pageTitle = computed(() => props.salaryType === 1
  ? t('pages.anchorGameShareCfgList.withSalaryTitle')
  : t('pages.anchorGameShareCfgList.noSalaryTitle'))

const loading = ref(false)
const saving = ref(false)
const tableData = ref<AnchorGameShareCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = reactive({
  id: '',
  level: 1,
  gameTotalGoldRevenue: 0,
  anchorGameSharePercent: 0,
  guildGameSharePercent: 0,
})

const percentValidator = (_rule: unknown, value: unknown, callback: (error?: Error) => void) => {
  const percent = Number(value)
  if (!Number.isFinite(percent) || percent < 0 || percent > 100) {
    callback(new Error(t('pages.anchorGameShareCfgList.sharePercentInvalid')))
    return
  }
  callback()
}

const formRules = computed<FormRules>(() => ({
  level: [
    {required: true, type: 'number', min: 1, message: t('pages.anchorGameShareCfgList.levelRequired'), trigger: 'change'},
  ],
  gameTotalGoldRevenue: [
    {required: true, message: t('pages.anchorGameShareCfgList.gameTotalGoldRevenueRequired'), trigger: 'change'},
  ],
  anchorGameSharePercent: [
    {required: true, message: t('pages.anchorGameShareCfgList.anchorGameSharePercentRequired'), trigger: 'change'},
    {validator: percentValidator, trigger: 'change'},
  ],
  guildGameSharePercent: [
    {required: true, message: t('pages.anchorGameShareCfgList.guildGameSharePercentRequired'), trigger: 'change'},
    {validator: percentValidator, trigger: 'change'},
  ],
}))

const formatPercent = (value: number | null | undefined) => `${Number(value || 0).toFixed(2)}%`

const fetchList = async () => {
  loading.value = true
  try {
    const response = await anchorGameShareCfgApi.getList({
      salaryType: props.salaryType,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = Number(response.total || 0)
  } catch (error) {
    console.error('fetch anchor game share config failed:', error)
    ElMessage.error(t('pages.anchorGameShareCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.id = ''
  formData.level = 1
  formData.gameTotalGoldRevenue = 0
  formData.anchorGameSharePercent = 0
  formData.guildGameSharePercent = 0
}

const handleAdd = () => {
  resetForm()
  dialogTitle.value = t('pages.anchorGameShareCfgList.addTier')
  dialogVisible.value = true
}

const handleEdit = (row: AnchorGameShareCfg) => {
  formData.id = String(row.id || '')
  formData.level = Number(row.level || 1)
  formData.gameTotalGoldRevenue = Number(row.gameTotalGoldRevenue || 0)
  formData.anchorGameSharePercent = Number(row.anchorGameSharePercent || 0)
  formData.guildGameSharePercent = Number(row.guildGameSharePercent || 0)
  dialogTitle.value = t('pages.anchorGameShareCfgList.editTier')
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!await formRef.value?.validate()) return
  saving.value = true
  try {
    const payload = {
      level: formData.level,
      gameTotalGoldRevenue: formData.gameTotalGoldRevenue,
      anchorGameSharePercent: formData.anchorGameSharePercent,
      guildGameSharePercent: formData.guildGameSharePercent,
    }
    if (formData.id) {
      await anchorGameShareCfgApi.update({id: formData.id, ...payload})
      ElMessage.success(t('common.updateSuccess'))
    } else {
      await anchorGameShareCfgApi.create({salaryType: props.salaryType, ...payload})
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    await fetchList()
  } catch (error) {
    console.error('save anchor game share config failed:', error)
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row: AnchorGameShareCfg) => {
  try {
    await ElMessageBox.confirm(
      t('pages.anchorGameShareCfgList.deleteConfirm', {level: row.level}),
      t('common.confirmDelete'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )
    await anchorGameShareCfgApi.remove(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    if (tableData.value.length === 1 && currentPage.value > 1) {
      currentPage.value -= 1
    }
    await fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete anchor game share config failed:', error)
    }
  }
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

watch(() => props.salaryType, () => {
  currentPage.value = 1
  fetchList()
})
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

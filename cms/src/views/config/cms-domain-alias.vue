<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.CMSDomainAliasManagement') }}</span>
          <el-button :disabled="loading" @click="fetchList">{{ t('common.refresh') }}</el-button>
        </div>
      </template>

      <el-alert
          :title="t('pages.cmsDomainAlias.noticeTitle')"
          :description="t('pages.cmsDomainAlias.noticeBody')"
          type="info"
          show-icon
          :closable="false"
          class="notice"
      />

      <div class="summary-row">
        <div class="summary-item">
          <span class="summary-label">{{ t('pages.cmsDomainAlias.primaryDomain') }}</span>
          <span class="summary-value">{{ primaryDomain || '—' }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">{{ t('pages.cmsDomainAlias.deployRoot') }}</span>
          <span class="summary-value">{{ deployRoot || '—' }}</span>
        </div>
      </div>

      <div v-if="can('save')" class="add-row">
        <el-input
            v-model="newDomain"
            clearable
            maxlength="253"
            class="add-input"
            :placeholder="t('pages.cmsDomainAlias.domainPlaceholder')"
            @keyup.enter="handleAdd"
        />
        <el-button type="primary" :loading="adding" @click="handleAdd">
          {{ t('pages.cmsDomainAlias.addDomain') }}
        </el-button>
      </div>

      <el-table :data="items" stripe class="alias-table">
        <el-table-column prop="domain" :label="t('pages.cmsDomainAlias.columnDomain')" min-width="220"/>
        <el-table-column prop="updatedAt" :label="t('pages.cmsDomainAlias.columnUpdatedAt')" min-width="180"/>
        <el-table-column v-if="can('delete')" :label="t('common.actions')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" link @click="() => handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <span>{{ t('pages.cmsDomainAlias.empty') }}</span>
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox} from 'element-plus'
import {cmsDomainAliasApi, type CMSDomainAliasItem} from '@/api/modules/cms-domain-alias'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('CMSDomainAliasManagement')
const loading = ref(false)
const adding = ref(false)
const newDomain = ref('')
const primaryDomain = ref('')
const deployRoot = ref('')
const items = ref<CMSDomainAliasItem[]>([])

function isSingleDomain(value: string): boolean {
  const domain = value.trim()
  if (!domain || domain.includes(',')) {
    return false
  }
  return !domain.includes('://') && !/[\\/?#@\s]/.test(domain) && domain.length <= 253
}

async function fetchList() {
  loading.value = true
  try {
    const data = await cmsDomainAliasApi.list()
    primaryDomain.value = data.primaryDomain || ''
    deployRoot.value = data.root || ''
    items.value = data.items || []
  } catch (error) {
    console.error('fetch CMS domain alias list failed:', error)
    ElMessage.error(t('pages.cmsDomainAlias.fetchFailed'))
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  const domain = newDomain.value.trim()
  if (!isSingleDomain(domain)) {
    ElMessage.warning(t('pages.cmsDomainAlias.domainInvalid'))
    return
  }
  adding.value = true
  try {
    await cmsDomainAliasApi.add(domain)
    newDomain.value = ''
    ElMessage.success(t('pages.cmsDomainAlias.addSuccess'))
    await fetchList()
  } catch (error) {
    console.error('add CMS domain alias failed:', error)
    ElMessage.error(t('pages.cmsDomainAlias.addFailed'))
  } finally {
    adding.value = false
  }
}

async function handleDelete(row: CMSDomainAliasItem) {
  try {
    await ElMessageBox.confirm(
        t('pages.cmsDomainAlias.deleteConfirm', {domain: row.domain}),
        t('common.confirm'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  try {
    await cmsDomainAliasApi.remove(row.id)
    ElMessage.success(t('pages.cmsDomainAlias.deleteSuccess'))
    await fetchList()
  } catch (error) {
    console.error('delete CMS domain alias failed:', error)
    ElMessage.error(t('pages.cmsDomainAlias.deleteFailed'))
  }
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

.notice {
  margin-bottom: 16px;
}

.summary-row {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 20px;
  font-size: 14px;
}

.summary-label {
  color: #64748b;
  margin-right: 8px;
}

.summary-value {
  color: #0f172a;
  font-weight: 600;
}

.add-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  max-width: 640px;
}

.add-input {
  flex: 1;
}

.alias-table {
  width: 100%;
}
</style>

<template>
  <el-dialog
      v-model="visible"
      :title="title || t('pages.guildList.leaderPickerTitle')"
      destroy-on-close
      width="860px"
      @closed="handleClosed"
      @open="handleOpen"
  >
    <el-form :model="searchForm" class="search-form" inline>
      <el-form-item :label="t('pages.cmsUserList.username')">
        <el-input v-model="searchForm.name" clearable :placeholder="t('pages.cmsUserList.usernamePlaceholder')"/>
      </el-form-item>
      <el-form-item :label="t('pages.cmsUserList.role')">
        <el-select
            v-model="searchForm.roleId"
            clearable
            filterable
            :loading="roleLoading"
            :placeholder="t('pages.cmsUserList.selectRole')"
            style="width: 180px"
        >
          <el-option
              v-for="role in roleOptions"
              :key="String(role.id)"
              :label="role.name"
              :value="String(role.id)"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
        <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table v-loading="loading" :data="tableData" highlight-current-row style="width: 100%">
      <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
      <el-table-column label="ID" prop="id" width="100"/>
      <el-table-column :label="t('pages.cmsUserList.username')" min-width="140" prop="name"/>
      <el-table-column :label="t('pages.cmsUserList.role')" min-width="120" prop="roleName">
        <template #default="{ row }">{{ row.roleName || '-' }}</template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.remark')" min-width="160" prop="remark" show-overflow-tooltip>
        <template #default="{ row }">{{ row.remark || '-' }}</template>
      </el-table-column>
      <el-table-column fixed="right" :label="t('common.actions')" width="100">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleSelect(row)">{{ t('pages.guildList.pickLeader') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @current-change="fetchList"
          @size-change="handleSizeChange"
      />
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import {computed, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {cmsUserApi, roleApi} from '@/api'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Role} from '@/api/modules/role'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title?: string
}>(), {
  title: '',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  select: [user: CMSUser]
}>()

const {t} = useI18n()
const loading = ref(false)
const roleLoading = ref(false)
const tableData = ref<CMSUser[]>([])
const roleOptions = ref<Role[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive({
  name: '',
  roleId: '',
})

const visible = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit('update:modelValue', value),
})

const fetchRoleOptions = async () => {
  if (roleOptions.value.length > 0) {
    return
  }
  roleLoading.value = true
  try {
    const response = await roleApi.getAllRoles()
    roleOptions.value = response.data || []
  } catch (error) {
    console.error('fetch role list failed:', error)
  } finally {
    roleLoading.value = false
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await cmsUserApi.getCMSUserList({
      name: searchForm.name.trim(),
      roleId: searchForm.roleId || undefined,
      status: 1,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('fetch cms user list failed:', error)
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.roleId = ''
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchList()
}

const formatRowIndex = (index: number) =>
    (currentPage.value - 1) * pageSize.value + index + 1

const handleSelect = (row: CMSUser) => {
  if (row.status !== 1) {
    return
  }
  emit('select', row)
  visible.value = false
}

const handleOpen = async () => {
  await fetchRoleOptions()
  await fetchList()
}

const handleClosed = () => {
  searchForm.name = ''
  searchForm.roleId = ''
  currentPage.value = 1
  pageSize.value = 10
  tableData.value = []
  total.value = 0
}
</script>

<style scoped>
.search-form {
  margin-bottom: 16px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

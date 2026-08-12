<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.guildList.addGuild') }}</el-button>
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.guildList.guildName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.guildList.guildNameSearchPlaceholder')"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchGuildList">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table v-loading="loading" :data="tableData" style="width: 100%">
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('pages.guildList.guildName')" prop="name"/>
          <el-table-column :label="t('pages.guildList.leader')" min-width="180" show-overflow-tooltip>
            <template #default="{ row }">
              {{ formatLeader(row) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.contact')" prop="contact" width="160"/>
          <el-table-column :label="t('pages.guildList.description')" prop="description" show-overflow-tooltip/>
          <el-table-column :label="t('common.status')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? t('common.enabled') : t('common.disabled') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column :label="t('common.actions')" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
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
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="100px">
        <el-form-item :label="t('pages.guildList.guildName')" prop="name">
          <el-input v-model="currentRow.name" :placeholder="t('pages.guildList.guildNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.leader')" prop="leaderId">
          <el-select
              v-model="currentRow.leaderId"
              :loading="cmsUserLoading"
              :remote-method="searchCmsUsers"
              clearable
              filterable
              :placeholder="t('pages.guildList.leaderSearchPlaceholder')"
              remote
              reserve-keyword
              style="width: 100%"
              @focus="loadInitialCmsUsers"
          >
            <el-option
                v-for="item in cmsUserOptions"
                :key="item.id"
                :label="formatCmsUserOption(item)"
                :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.contact')" prop="contact">
          <el-input v-model="currentRow.contact" :placeholder="t('pages.guildList.contactPlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.guildList.description')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="t('pages.guildList.descriptionPlaceholder')" type="textarea"/>
        </el-form-item>
        <el-form-item :label="t('common.status')" prop="status">
          <el-radio-group v-model="currentRow.status">
            <el-radio :label="1">{{ t('common.enabled') }}</el-radio>
            <el-radio :label="0">{{ t('common.disabled') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {cmsUserApi, guildApi} from '@/api'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild} from '@/types/api.ts'

interface SearchForm {
  name: string
}

interface GuildForm {
  id: string
  name: string
  leaderId: string
  contact: string
  description: string
  status: number
}

const {t} = useI18n()
const loading = ref(false)
const cmsUserLoading = ref(false)
const cmsUserOptions = ref<CMSUser[]>([])
const tableData = ref<Guild[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  name: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<GuildForm>({
  id: '',
  name: '',
  leaderId: '',
  contact: '',
  description: '',
  status: 1
})

const formRef = ref<FormInstance>()
let cmsUserSearchTimer: ReturnType<typeof setTimeout> | undefined

const formatCmsUserOption = (item: CMSUser) => `${item.name} (${item.id})`

const formatLeader = (row: Guild) => {
  if (row.leaderName) {
    return `${row.leaderName} (${row.leaderId})`
  }
  return row.leaderId || '-'
}

const mergeCmsUserOptions = (users: CMSUser[]) => {
  const map = new Map<string, CMSUser>()
  for (const item of cmsUserOptions.value) {
    map.set(item.id, item)
  }
  for (const item of users) {
    map.set(item.id, item)
  }
  cmsUserOptions.value = [...map.values()]
}

const fetchCmsUserOptions = async (key = '') => {
  cmsUserLoading.value = true
  try {
    const response = await cmsUserApi.getCMSUserList({
      key: key.trim(),
      pageIndex: 1,
      pageSize: 20
    })
    mergeCmsUserOptions(response.data)
  } catch (error) {
    console.error('fetch cms user list failed:', error)
  } finally {
    cmsUserLoading.value = false
  }
}

const searchCmsUsers = (query: string) => {
  if (cmsUserSearchTimer) {
    clearTimeout(cmsUserSearchTimer)
  }
  cmsUserSearchTimer = setTimeout(() => {
    fetchCmsUserOptions(query)
  }, 300)
}

const loadInitialCmsUsers = () => {
  if (cmsUserOptions.value.length === 0) {
    fetchCmsUserOptions('')
  }
}

const ensureLeaderOption = async (leaderId: string) => {
  const id = leaderId.trim()
  if (!id || id === '0') {
    return
  }
  if (cmsUserOptions.value.some(item => item.id === id)) {
    return
  }
  await fetchCmsUserOptions(id)
}

const formRules = computed<FormRules>(() => ({
  name: [
    {required: true, message: t('pages.guildList.nameRequired'), trigger: 'blur'},
    {min: 2, max: 32, message: t('pages.guildList.nameLength'), trigger: 'blur'}
  ],
  description: [
    {max: 200, message: t('pages.guildList.descriptionMaxLength'), trigger: 'blur'}
  ]
}))

const fetchGuildList = async () => {
  loading.value = true
  try {
    const response = await guildApi.getGuildList({
      name: searchForm.name,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch guild list failed:', error)
    ElMessage.error(t('pages.guildList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchGuildList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchGuildList()
}

const handleAdd = async () => {
  dialogTitle.value = t('pages.guildList.addGuild')
  currentRow.value = {
    id: '',
    name: '',
    leaderId: '',
    contact: '',
    description: '',
    status: 1
  }
  cmsUserOptions.value = []
  dialogVisible.value = true
  await fetchCmsUserOptions('')
}

const handleEdit = async (row: Guild) => {
  dialogTitle.value = t('pages.guildList.editGuild')
  const leaderId = row.leaderId && row.leaderId !== '0' ? row.leaderId : ''
  currentRow.value = {
    id: row.id,
    name: row.name,
    leaderId,
    contact: row.contact,
    description: row.description,
    status: row.status
  }
  cmsUserOptions.value = []
  if (leaderId) {
    if (row.leaderName) {
      cmsUserOptions.value = [{id: leaderId, name: row.leaderName} as CMSUser]
    }
    await ensureLeaderOption(leaderId)
  } else {
    await fetchCmsUserOptions('')
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Guild) => {
  try {
    await ElMessageBox.confirm(t('pages.guildList.deleteConfirm', {name: row.name}), t('common.confirmDelete'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await guildApi.deleteGuild(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchGuildList()
  } catch (error) {
    console.error('delete failed:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const leaderId = Number(currentRow.value.leaderId) || 0
        if (currentRow.value.id) {
          await guildApi.updateGuild({...currentRow.value, leaderId})
        } else {
          const {name, contact, description, status} = currentRow.value
          await guildApi.createGuild({name, leaderId, contact, description, status})
        }

        ElMessage.success(currentRow.value.id ? t('common.updateSuccess') : t('common.createSuccess'))
        dialogVisible.value = false
        fetchGuildList()
      } catch (error) {
        console.error('save guild failed:', error)
        ElMessage.error(currentRow.value.id ? t('pages.guildList.updateFailed') : t('pages.guildList.createFailed'))
      }
    }
  })
}

const resetSearch = () => {
  searchForm.name = ''
  fetchGuildList()
}

onMounted(() => {
  fetchGuildList()
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

.table-header {
  margin-bottom: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>

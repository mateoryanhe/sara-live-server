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
          <el-button @click="downloadTxtTemplate">{{ t('pages.guildList.downloadTemplate') }}</el-button>
        </div>
        <el-alert
            :closable="false"
            class="import-tip"
            show-icon
            :title="t('pages.guildList.importTip')"
            type="info"
        />

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.guildList.guildName')">
            <el-input v-model="searchForm.name" clearable :placeholder="t('pages.guildList.guildNameSearchPlaceholder')"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchGuildList">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table
            v-loading="loading"
            :data="tableData"
            highlight-current-row
            style="width: 100%"
            @current-change="handleCurrentRowChange"
        >
          <el-table-column label="ID" prop="id" width="190"/>
          <el-table-column :label="t('pages.guildList.guildName')" prop="name"/>
          <el-table-column :label="t('pages.guildList.leader')" width="140" show-overflow-tooltip>
            <template #default="{ row }">
              {{ formatLeader(row) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.guildList.description')" prop="description" show-overflow-tooltip/>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="500">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="can('viewMembers')"
                  size="small"
                  @click="handleViewMembers(row)"
              >
                {{ t('pages.guildList.viewMembers') }}
              </el-button>
              <el-button
                  v-if="can('batchSetAnchor')"
                  :loading="importing"
                  size="small"
                  @click="triggerCsvImport(row, 1)"
              >
                {{ t('pages.guildList.importNormalAnchor') }}
              </el-button>
              <el-button
                  v-if="can('batchSetSeniorAnchor')"
                  :loading="importing"
                  size="small"
                  @click="triggerCsvImport(row, 7)"
              >
                {{ t('pages.guildList.importSeniorAnchor') }}
              </el-button>
              <el-button
                  v-if="can('offShelf')"
                  size="small"
                  type="warning"
                  @click="handleOffShelf(row)"
              >
                {{ t('common.offShelf') }}
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
        <el-form-item :label="t('pages.guildList.description')" prop="description">
          <el-input v-model="currentRow.description" :placeholder="t('pages.guildList.descriptionPlaceholder')" type="textarea"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>

    <input
        ref="csvInputRef"
        accept=".txt,.csv,text/plain,text/csv"
        class="hidden-file-input"
        type="file"
        @change="onCsvFileSelected"
    />
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRouter} from 'vue-router'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules} from 'element-plus'
import {cmsUserApi, guildApi} from '@/api'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild, GuildAnchorImportResultState, ImportGuildAnchorRow} from '@/types/api.ts'
import {usePagePermission} from '@/composables/usePagePermission'

const GUILD_ANCHOR_IMPORT_RESULT_KEY = 'guildAnchorImportResult'

interface SearchForm {
  name: string
}

interface GuildForm {
  id: string
  name: string
  leaderId: string
  description: string
}

const {t} = useI18n()
const router = useRouter()
const {can} = usePagePermission('GuildManagement')
const loading = ref(false)
const importing = ref(false)
const cmsUserLoading = ref(false)
const cmsUserOptions = ref<CMSUser[]>([])
const tableData = ref<Guild[]>([])
const selectedGuild = ref<Guild | null>(null)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const pendingAnchorType = ref<1 | 7>(1)
const csvInputRef = ref<HTMLInputElement | null>(null)

const searchForm = reactive<SearchForm>({
  name: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const currentRow = ref<GuildForm>({
  id: '',
  name: '',
  leaderId: '',
  description: '',
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
    if (selectedGuild.value && !tableData.value.some(item => item.id === selectedGuild.value?.id)) {
      selectedGuild.value = null
    }
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

const handleCurrentRowChange = (row: Guild | null) => {
  selectedGuild.value = row
}

const handleAdd = async () => {
  dialogTitle.value = t('pages.guildList.addGuild')
  currentRow.value = {
    id: '',
    name: '',
    leaderId: '',
    description: '',
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
    description: row.description,
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

const handleOffShelf = async (row: Guild) => {
  try {
    await ElMessageBox.confirm(
      t('pages.guildList.offShelfConfirm', {name: row.name}),
      t('common.confirmOffShelf'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
    )

    await guildApi.deleteGuild(row.id)
    ElMessage.success(t('pages.guildList.offShelfSuccess'))
    fetchGuildList()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    console.error('off shelf guild failed:', error)
    ElMessage.error(t('pages.guildList.offShelfFailed'))
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
          const {name, description} = currentRow.value
          await guildApi.createGuild({name, leaderId, description})
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

const parseCsvLine = (line: string): string[] => {
  const result: string[] = []
  let current = ''
  let inQuotes = false
  for (let i = 0; i < line.length; i++) {
    const ch = line[i]
    if (ch === '"') {
      if (inQuotes && line[i + 1] === '"') {
        current += '"'
        i++
      } else {
        inQuotes = !inQuotes
      }
      continue
    }
    if ((ch === ',' || ch === '\t') && !inQuotes) {
      result.push(current.trim())
      current = ''
      continue
    }
    current += ch
  }
  result.push(current.trim())
  return result
}

const normalizeHeader = (value: string) => value.trim().toLowerCase().replace(/^\ufeff/, '')

const isHeaderRow = (cols: string[]) => {
  if (cols.length < 1) return false
  const a = normalizeHeader(cols[0])
  const userHeaders = new Set(['user_id', 'userid', '用户id', 'id'])
  return userHeaders.has(a)
}

const parseGuildAnchorCsv = (text: string): ImportGuildAnchorRow[] => {
  const lines = text.split(/\r?\n/).map(line => line.trim()).filter(Boolean)
  if (lines.length === 0) {
    return []
  }
  let start = 0
  const firstCols = parseCsvLine(lines[0])
  if (isHeaderRow(firstCols)) {
    start = 1
  }
  const rows: ImportGuildAnchorRow[] = []
  const seen = new Set<string>()
  for (let i = start; i < lines.length; i++) {
    const cols = parseCsvLine(lines[i])
    if (cols.length < 1) continue
    const userId = cols[0].replace(/^\ufeff/, '').trim()
    if (!userId || !/^\d+$/.test(userId)) continue
    if (seen.has(userId)) continue
    seen.add(userId)
    rows.push({userId})
  }
  return rows
}

const downloadTxtTemplate = () => {
  // 用 txt，避免 Excel 把大整数改成科学计数法
  const content = 'user_id\n10001\n'
  const blob = new Blob(['\uFEFF' + content], {type: 'text/plain;charset=utf-8;'})
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'guild_anchor_import_template.txt'
  link.click()
  URL.revokeObjectURL(url)
}

const triggerCsvImport = (row: Guild, anchorType: 1 | 7) => {
  selectedGuild.value = row
  pendingAnchorType.value = anchorType
  if (csvInputRef.value) {
    csvInputRef.value.value = ''
    csvInputRef.value.click()
  }
}

const onCsvFileSelected = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !selectedGuild.value) {
    return
  }
  importing.value = true
  try {
    const text = await file.text()
    const rows = parseGuildAnchorCsv(text)
    if (rows.length === 0) {
      ElMessage.warning(t('pages.guildList.csvEmpty'))
      return
    }
    const response = await guildApi.importGuildAnchors({
      guildId: selectedGuild.value.id,
      anchorType: pendingAnchorType.value,
      rows,
    })
    const state: GuildAnchorImportResultState = {
      guildId: selectedGuild.value.id,
      guildName: selectedGuild.value.name,
      anchorType: pendingAnchorType.value,
      successCount: response?.successCount ?? 0,
      failCount: response?.failCount ?? 0,
      fails: response?.fails ?? [],
    }
    sessionStorage.setItem(GUILD_ANCHOR_IMPORT_RESULT_KEY, JSON.stringify(state))
    await router.push({name: 'GuildAnchorImportResult'})
  } catch (error) {
    console.error('import guild anchors failed:', error)
    ElMessage.error(t('pages.guildList.importFailed'))
  } finally {
    importing.value = false
    if (csvInputRef.value) {
      csvInputRef.value.value = ''
    }
  }
}

onMounted(() => {
  fetchGuildList()
})

const handleViewMembers = (row: Guild) => {
  router.push({
    name: 'GuildMembers',
    query: {
      guildId: row.id,
      guildName: row.name,
    },
  })
}
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
  margin-bottom: 12px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.import-tip {
  margin-bottom: 16px;
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

.hidden-file-input {
  display: none;
}
</style>

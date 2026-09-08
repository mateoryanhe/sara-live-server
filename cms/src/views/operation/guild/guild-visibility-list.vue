<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildVisibilityManagement') }}</span>
        </div>
      </template>

      <div class="content">
        <el-alert
            :closable="false"
            class="tip-alert"
            show-icon
            :title="t('pages.guildVisibilityList.tip')"
            type="info"
        />

        <el-form class="search-form" inline>
          <el-form-item :label="t('pages.guildVisibilityList.selectCmsUser')">
            <div class="cms-user-picker-field">
              <el-input
                  :model-value="selectedCmsUserText"
                  class="cms-user-picker-input"
                  readonly
                  :placeholder="t('pages.guildVisibilityList.selectCmsUserPlaceholder')"
              />
              <el-button type="primary" @click="cmsUserPickerVisible = true">
                {{ t('pages.guildVisibilityList.pickCmsUser') }}
              </el-button>
              <el-button
                  v-if="selectedCmsUserId"
                  link
                  type="danger"
                  @click="clearCmsUser"
              >
                {{ t('pages.guildVisibilityList.clearCmsUser') }}
              </el-button>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :disabled="!selectedCmsUserId" :loading="loading" @click="reloadForUser">
              {{ t('common.search') }}
            </el-button>
          </el-form-item>
        </el-form>

        <el-empty v-if="!selectedCmsUserId" :description="t('pages.guildVisibilityList.emptyHint')"/>

        <div v-else v-loading="loading" class="panel-wrap">
          <div class="guild-panel">
            <div class="panel-header">
              <span>{{ t('pages.guildVisibilityList.leftTitle') }} ({{ availableGuilds.length }})</span>
              <el-input
                  v-model="leftFilter"
                  clearable
                  class="panel-filter"
                  :placeholder="t('pages.guildVisibilityList.filterPlaceholder')"
              />
            </div>
            <el-table
                ref="leftTableRef"
                :data="filteredAvailableGuilds"
                height="100%"
                row-key="id"
                @selection-change="onLeftSelectionChange"
            >
              <el-table-column type="selection" width="42"/>
              <el-table-column label="ID" prop="id" width="120" show-overflow-tooltip/>
              <el-table-column :label="t('pages.guildList.guildName')" min-width="140" prop="name" show-overflow-tooltip/>
              <el-table-column :label="t('pages.guildList.guildType')" width="100">
                <template #default="{ row }">{{ guildTypeLabel(row.guildType) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.guildList.leader')" min-width="140" show-overflow-tooltip>
                <template #default="{ row }">{{ formatLeader(row) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.guildList.creator')" min-width="140" show-overflow-tooltip>
                <template #default="{ row }">{{ formatCreator(row) }}</template>
              </el-table-column>
              <el-table-column :label="t('common.createdAt')" width="160" prop="createdAt"/>
            </el-table>
          </div>

          <div class="panel-actions">
            <el-button
                type="primary"
                :disabled="!leftSelected.length || syncing"
                :loading="syncing"
                @click="grantSelected"
            >
              {{ t('pages.guildVisibilityList.toRight') }} →
            </el-button>
            <el-button
                :disabled="!rightSelected.length || syncing"
                :loading="syncing"
                @click="revokeSelected"
            >
              ← {{ t('pages.guildVisibilityList.toLeft') }}
            </el-button>
          </div>

          <div class="guild-panel">
            <div class="panel-header">
              <span>{{ t('pages.guildVisibilityList.rightTitle') }} ({{ grantedGuilds.length }})</span>
              <el-input
                  v-model="rightFilter"
                  clearable
                  class="panel-filter"
                  :placeholder="t('pages.guildVisibilityList.filterPlaceholder')"
              />
            </div>
            <el-table
                ref="rightTableRef"
                :data="filteredGrantedGuilds"
                height="100%"
                row-key="id"
                @selection-change="onRightSelectionChange"
            >
              <el-table-column type="selection" width="42"/>
              <el-table-column label="ID" prop="id" width="120" show-overflow-tooltip/>
              <el-table-column :label="t('pages.guildList.guildName')" min-width="140" prop="name" show-overflow-tooltip/>
              <el-table-column :label="t('pages.guildList.guildType')" width="100">
                <template #default="{ row }">{{ guildTypeLabel(row.guildType) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.guildList.leader')" min-width="140" show-overflow-tooltip>
                <template #default="{ row }">{{ formatLeader(row) }}</template>
              </el-table-column>
              <el-table-column :label="t('pages.guildList.creator')" min-width="140" show-overflow-tooltip>
                <template #default="{ row }">{{ formatCreator(row) }}</template>
              </el-table-column>
              <el-table-column :label="t('common.createdAt')" width="160" prop="createdAt"/>
            </el-table>
          </div>
        </div>
      </div>
    </el-card>

    <CmsUserPickerDialog
        v-model="cmsUserPickerVisible"
        :title="t('pages.guildVisibilityList.pickerTitle')"
        :select-text="t('pages.guildVisibilityList.pickCmsUser')"
        @select="handleCmsUserSelect"
    />
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox, type TableInstance} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {guildApi} from '@/api/modules/guild'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {Guild} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import CmsUserPickerDialog from '@/components/CmsUserPickerDialog.vue'

const {t} = useI18n()
const {can} = usePagePermission('GuildVisibilityManagement')

const loading = ref(false)
const syncing = ref(false)
const guildOptions = ref<Guild[]>([])
const grantedIds = ref<string[]>([])
const selectedCmsUserId = ref('')
const selectedCmsUserName = ref('')
const cmsUserPickerVisible = ref(false)
const leftFilter = ref('')
const rightFilter = ref('')
const leftSelected = ref<Guild[]>([])
const rightSelected = ref<Guild[]>([])
const leftTableRef = ref<TableInstance>()
const rightTableRef = ref<TableInstance>()

const selectedCmsUserText = computed(() => {
  if (selectedCmsUserName.value) {
    return `${selectedCmsUserName.value} (${selectedCmsUserId.value})`
  }
  return selectedCmsUserId.value
})

const grantedIdSet = computed(() => new Set(grantedIds.value))

const availableGuilds = computed(() =>
    guildOptions.value.filter((g) => !grantedIdSet.value.has(g.id)),
)

const grantedGuilds = computed(() =>
    guildOptions.value.filter((g) => grantedIdSet.value.has(g.id)),
)

const matchGuild = (g: Guild, keyword: string) => {
  const key = keyword.trim().toLowerCase()
  if (!key) {
    return true
  }
  return [g.id, g.name, g.leaderId, g.leaderName, g.creatorId, g.creatorName]
      .filter(Boolean)
      .some((v) => String(v).toLowerCase().includes(key))
}

const filteredAvailableGuilds = computed(() =>
    availableGuilds.value.filter((g) => matchGuild(g, leftFilter.value)),
)

const filteredGrantedGuilds = computed(() =>
    grantedGuilds.value.filter((g) => matchGuild(g, rightFilter.value)),
)

const guildTypeLabel = (guildType?: number) => {
  if (guildType === 1) {
    return t('pages.guildList.guildTypeCoinMerchant')
  }
  return t('pages.guildList.guildTypeNormal')
}

const formatLeader = (row: Guild) => {
  if (row.leaderName) {
    return `${row.leaderName} (${row.leaderId})`
  }
  return row.leaderId && row.leaderId !== '0' ? row.leaderId : '-'
}

const formatCreator = (row: Guild) => {
  if (row.creatorName) {
    return `${row.creatorName} (${row.creatorId})`
  }
  if (row.creatorId && row.creatorId !== '0') {
    return row.creatorId
  }
  return '-'
}

const clearTableSelection = () => {
  leftTableRef.value?.clearSelection()
  rightTableRef.value?.clearSelection()
  leftSelected.value = []
  rightSelected.value = []
}

const onLeftSelectionChange = (rows: Guild[]) => {
  leftSelected.value = rows
}

const onRightSelectionChange = (rows: Guild[]) => {
  rightSelected.value = rows
}

const fetchGuildOptions = async () => {
  try {
    const res = await guildApi.getGuildListForVisibility({pageIndex: 1, pageSize: 500, name: ''})
    guildOptions.value = res.data || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.fetchGuildFailed'))
  }
}

const fetchGrantedForUser = async () => {
  if (!selectedCmsUserId.value) {
    grantedIds.value = []
    return
  }
  loading.value = true
  try {
    const res = await guildApi.getGuildVisibilityByUserList(selectedCmsUserId.value)
    grantedIds.value = (res.list || []).map((item) => item.guildId)
    clearTableSelection()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const reloadForUser = async () => {
  if (!selectedCmsUserId.value) {
    ElMessage.warning(t('pages.guildVisibilityList.cmsUserRequired'))
    return
  }
  await fetchGrantedForUser()
}

const handleCmsUserSelect = async (user: CMSUser) => {
  selectedCmsUserId.value = user.id
  selectedCmsUserName.value = user.name || ''
  leftFilter.value = ''
  rightFilter.value = ''
  await fetchGrantedForUser()
}

const clearCmsUser = () => {
  selectedCmsUserId.value = ''
  selectedCmsUserName.value = ''
  grantedIds.value = []
  leftFilter.value = ''
  rightFilter.value = ''
  clearTableSelection()
}

const grantSelected = async () => {
  const movedKeys = leftSelected.value.map((g) => g.id)
  if (!movedKeys.length || !selectedCmsUserId.value || syncing.value) {
    return
  }
  if (!can('grant')) {
    ElMessage.warning(t('pages.guildVisibilityList.noGrantPermission'))
    return
  }
  syncing.value = true
  try {
    const res = await guildApi.batchGrantGuildVisibility({
      cmsUserId: selectedCmsUserId.value,
      guildIds: movedKeys,
    })
    grantedIds.value = [...new Set([...grantedIds.value, ...movedKeys])]
    clearTableSelection()
    ElMessage.success(`${t('pages.guildVisibilityList.grantSuccess')} (${res.grantedCount ?? movedKeys.length})`)
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.grantFailed'))
  } finally {
    syncing.value = false
  }
}

const revokeSelected = async () => {
  const movedKeys = rightSelected.value.map((g) => g.id)
  if (!movedKeys.length || !selectedCmsUserId.value || syncing.value) {
    return
  }
  if (!can('revoke')) {
    ElMessage.warning(t('pages.guildVisibilityList.noRevokePermission'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.guildVisibilityList.revokeConfirm'),
        t('common.confirm'),
        {type: 'warning'},
    )
  } catch {
    return
  }

  syncing.value = true
  try {
    const res = await guildApi.batchRevokeGuildVisibility({
      cmsUserId: selectedCmsUserId.value,
      guildIds: movedKeys,
    })
    const revokeSet = new Set(movedKeys)
    grantedIds.value = grantedIds.value.filter((id) => !revokeSet.has(id))
    clearTableSelection()
    ElMessage.success(`${t('pages.guildVisibilityList.revokeSuccess')} (${res.revokedCount ?? movedKeys.length})`)
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.revokeFailed'))
  } finally {
    syncing.value = false
  }
}

onMounted(async () => {
  await fetchGuildOptions()
})
</script>

<style scoped>
.tip-alert {
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 12px;
}

.cms-user-picker-field {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cms-user-picker-input {
  width: 280px;
}

.panel-wrap {
  display: flex;
  align-items: stretch;
  gap: 12px;
  width: 100%;
  min-height: calc(100vh - 280px);
  height: calc(100vh - 280px);
}

.guild-panel {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border-bottom: 1px solid var(--el-border-color);
  background: var(--el-fill-color-light);
  font-weight: 500;
}

.panel-filter {
  width: 220px;
}

.guild-panel :deep(.el-table) {
  flex: 1;
}

.panel-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 4px;
}

.panel-actions .el-button {
  margin: 0;
  min-width: 96px;
}
</style>

<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PlatformAnchorVisibilityManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.platformAnchorVisibilityList.tip')"
          type="info"
      />

      <el-form class="search-form" inline>
        <el-form-item :label="t('pages.platformAnchorVisibilityList.selectCmsUser')">
          <div class="cms-user-picker-field">
            <el-input
                :model-value="selectedCmsUserText"
                class="cms-user-picker-input"
                readonly
                :placeholder="t('pages.platformAnchorVisibilityList.selectCmsUserPlaceholder')"
            />
            <el-button type="primary" @click="cmsUserPickerVisible = true">
              {{ t('pages.platformAnchorVisibilityList.pickCmsUser') }}
            </el-button>
            <el-button v-if="selectedCmsUserId" link type="danger" @click="clearCmsUser">
              {{ t('pages.platformAnchorVisibilityList.clearCmsUser') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :disabled="!selectedCmsUserId" :loading="loading" @click="reloadForUser">
            {{ t('common.search') }}
          </el-button>
        </el-form-item>
      </el-form>

      <el-empty v-if="!selectedCmsUserId" :description="t('pages.platformAnchorVisibilityList.emptyHint')"/>

      <div v-else v-loading="loading" class="panel-wrap">
        <div class="anchor-panel">
          <div class="panel-header">
            <span>{{ t('pages.platformAnchorVisibilityList.leftTitle') }} ({{ availableAnchors.length }})</span>
            <el-input
                v-model="leftFilter"
                clearable
                class="panel-filter"
                :placeholder="t('pages.platformAnchorVisibilityList.filterPlaceholder')"
            />
          </div>
          <AnchorVisibilityTable
              ref="leftTableRef"
              :data="filteredAvailableAnchors"
              @selection-change="onLeftSelectionChange"
          />
        </div>

        <div class="panel-actions">
          <el-button
              type="primary"
              :disabled="!leftSelected.length || syncing"
              :loading="syncing"
              @click="grantSelected"
          >
            {{ t('pages.platformAnchorVisibilityList.toRight') }} →
          </el-button>
          <el-button
              :disabled="!rightSelected.length || syncing"
              :loading="syncing"
              @click="revokeSelected"
          >
            ← {{ t('pages.platformAnchorVisibilityList.toLeft') }}
          </el-button>
        </div>

        <div class="anchor-panel">
          <div class="panel-header">
            <span>{{ t('pages.platformAnchorVisibilityList.rightTitle') }} ({{ grantedAnchors.length }})</span>
            <el-input
                v-model="rightFilter"
                clearable
                class="panel-filter"
                :placeholder="t('pages.platformAnchorVisibilityList.filterPlaceholder')"
            />
          </div>
          <AnchorVisibilityTable
              ref="rightTableRef"
              :data="filteredGrantedAnchors"
              @selection-change="onRightSelectionChange"
          />
        </div>
      </div>
    </el-card>

    <CmsUserPickerDialog
        v-model="cmsUserPickerVisible"
        :title="t('pages.platformAnchorVisibilityList.pickerTitle')"
        :select-text="t('pages.platformAnchorVisibilityList.pickCmsUser')"
        @select="handleCmsUserSelect"
    />
  </div>
</template>

<script lang="ts" setup>
import {computed, defineComponent, h, onMounted, ref} from 'vue'
import {ElAvatar, ElMessage, ElMessageBox, ElTable, ElTableColumn, ElTag, type TableInstance} from 'element-plus'
import {useI18n} from 'vue-i18n'
import accountApi from '@/api/modules/account'
import type {CMSUser} from '@/api/modules/cmsuser'
import type {AnchorListItem} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'
import CmsUserPickerDialog from '@/components/CmsUserPickerDialog.vue'

const {t} = useI18n()
const {can} = usePagePermission('PlatformAnchorVisibilityManagement')

const AnchorVisibilityTable = defineComponent({
  name: 'AnchorVisibilityTable',
  props: {
    data: {type: Array as () => AnchorListItem[], required: true},
  },
  emits: ['selection-change'],
  setup(props, {emit, expose}) {
    const tableRef = ref<TableInstance>()
    expose({clearSelection: () => tableRef.value?.clearSelection()})
    return () => h(ElTable, {
      ref: tableRef,
      data: props.data,
      height: '100%',
      rowKey: 'id',
      onSelectionChange: (rows: AnchorListItem[]) => emit('selection-change', rows),
    }, {
      default: () => [
        h(ElTableColumn, {type: 'selection', width: 42}),
        h(ElTableColumn, {label: 'ID', prop: 'id', width: 130, showOverflowTooltip: true}),
        h(ElTableColumn, {label: t('pages.platformAnchorVisibilityList.anchor'), minWidth: 180}, {
          default: ({row}: {row: AnchorListItem}) => h('div', {class: 'anchor-cell'}, [
            h(ElAvatar, {size: 32, src: row.avatar || ''}),
            h('div', {class: 'anchor-text'}, [
              h('span', row.nickname || '-'),
              h('small', row.phone || ''),
            ]),
          ]),
        }),
        h(ElTableColumn, {label: t('pages.platformAnchorVisibilityList.anchorType'), width: 105}, {
          default: ({row}: {row: AnchorListItem}) => row.userType === 7
              ? t('pages.platformAnchorVisibilityList.seniorAnchor')
              : t('pages.platformAnchorVisibilityList.normalAnchor'),
        }),
        h(ElTableColumn, {label: t('pages.platformAnchorVisibilityList.liveStatus'), width: 95}, {
          default: ({row}: {row: AnchorListItem}) => h(ElTag, {
            type: row.liveStatus === 1 ? 'success' : 'info',
            size: 'small',
          }, () => row.liveStatus === 1
              ? t('pages.platformAnchorVisibilityList.living')
              : t('pages.platformAnchorVisibilityList.offline')),
        }),
      ],
    })
  },
})

const loading = ref(false)
const syncing = ref(false)
const anchorOptions = ref<AnchorListItem[]>([])
const grantedIds = ref<string[]>([])
const selectedCmsUserId = ref('')
const selectedCmsUserName = ref('')
const cmsUserPickerVisible = ref(false)
const leftFilter = ref('')
const rightFilter = ref('')
const leftSelected = ref<AnchorListItem[]>([])
const rightSelected = ref<AnchorListItem[]>([])
const leftTableRef = ref<{clearSelection: () => void}>()
const rightTableRef = ref<{clearSelection: () => void}>()

const selectedCmsUserText = computed(() => selectedCmsUserName.value
    ? `${selectedCmsUserName.value} (${selectedCmsUserId.value})`
    : selectedCmsUserId.value)
const grantedIdSet = computed(() => new Set(grantedIds.value))
const availableAnchors = computed(() => anchorOptions.value.filter((item) => !grantedIdSet.value.has(String(item.id))))
const grantedAnchors = computed(() => anchorOptions.value.filter((item) => grantedIdSet.value.has(String(item.id))))

const matchAnchor = (item: AnchorListItem, keyword: string) => {
  const key = keyword.trim().toLowerCase()
  if (!key) return true
  return [item.id, item.nickname, item.phone]
      .filter(Boolean)
      .some((value) => String(value).toLowerCase().includes(key))
}
const filteredAvailableAnchors = computed(() => availableAnchors.value.filter((item) => matchAnchor(item, leftFilter.value)))
const filteredGrantedAnchors = computed(() => grantedAnchors.value.filter((item) => matchAnchor(item, rightFilter.value)))

const clearTableSelection = () => {
  leftTableRef.value?.clearSelection()
  rightTableRef.value?.clearSelection()
  leftSelected.value = []
  rightSelected.value = []
}
const onLeftSelectionChange = (rows: AnchorListItem[]) => { leftSelected.value = rows }
const onRightSelectionChange = (rows: AnchorListItem[]) => { rightSelected.value = rows }

const fetchAnchorOptions = async () => {
  try {
    const res = await accountApi.getPlatformAnchorListForVisibility({pageIndex: 1, pageSize: 500, key: ''})
    anchorOptions.value = res.data || []
  } catch (error) {
    console.error(error)
    ElMessage.error(t('pages.platformAnchorVisibilityList.fetchAnchorFailed'))
  }
}

const fetchGrantedForUser = async () => {
  if (!selectedCmsUserId.value) {
    grantedIds.value = []
    return
  }
  loading.value = true
  try {
    const res = await accountApi.getPlatformAnchorVisibilityByUserList(selectedCmsUserId.value)
    grantedIds.value = (res.list || []).map((item) => String(item.anchorId))
    clearTableSelection()
  } catch (error) {
    console.error(error)
    ElMessage.error(t('pages.platformAnchorVisibilityList.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const reloadForUser = async () => {
  if (!selectedCmsUserId.value) {
    ElMessage.warning(t('pages.platformAnchorVisibilityList.cmsUserRequired'))
    return
  }
  await Promise.all([fetchAnchorOptions(), fetchGrantedForUser()])
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
  const anchorIds = leftSelected.value.map((item) => item.id)
  if (!anchorIds.length || !selectedCmsUserId.value || syncing.value) return
  if (!can('grant')) {
    ElMessage.warning(t('pages.platformAnchorVisibilityList.noGrantPermission'))
    return
  }
  syncing.value = true
  try {
    const res = await accountApi.batchGrantPlatformAnchorVisibility({cmsUserId: selectedCmsUserId.value, anchorIds})
    grantedIds.value = [...new Set([...grantedIds.value, ...anchorIds.map(String)])]
    clearTableSelection()
    ElMessage.success(`${t('pages.platformAnchorVisibilityList.grantSuccess')} (${res.grantedCount ?? anchorIds.length})`)
  } catch (error) {
    console.error(error)
    ElMessage.error(t('pages.platformAnchorVisibilityList.grantFailed'))
  } finally {
    syncing.value = false
  }
}

const revokeSelected = async () => {
  const anchorIds = rightSelected.value.map((item) => item.id)
  if (!anchorIds.length || !selectedCmsUserId.value || syncing.value) return
  if (!can('revoke')) {
    ElMessage.warning(t('pages.platformAnchorVisibilityList.noRevokePermission'))
    return
  }
  try {
    await ElMessageBox.confirm(t('pages.platformAnchorVisibilityList.revokeConfirm'), t('common.confirm'), {type: 'warning'})
  } catch {
    return
  }
  syncing.value = true
  try {
    const res = await accountApi.batchRevokePlatformAnchorVisibility({cmsUserId: selectedCmsUserId.value, anchorIds})
    const removed = new Set(anchorIds.map(String))
    grantedIds.value = grantedIds.value.filter((id) => !removed.has(id))
    clearTableSelection()
    ElMessage.success(`${t('pages.platformAnchorVisibilityList.revokeSuccess')} (${res.revokedCount ?? anchorIds.length})`)
  } catch (error) {
    console.error(error)
    ElMessage.error(t('pages.platformAnchorVisibilityList.revokeFailed'))
  } finally {
    syncing.value = false
  }
}

onMounted(fetchAnchorOptions)
</script>

<style scoped>
.tip-alert { margin-bottom: 12px; }
.search-form { margin-bottom: 12px; }
.cms-user-picker-field { display: flex; align-items: center; gap: 8px; }
.cms-user-picker-input { width: 280px; }
.panel-wrap { display: flex; align-items: stretch; gap: 12px; width: 100%; min-height: calc(100vh - 280px); height: calc(100vh - 280px); }
.anchor-panel { flex: 1; min-width: 0; display: flex; flex-direction: column; border: 1px solid var(--el-border-color); border-radius: 4px; overflow: hidden; }
.panel-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 10px 12px; border-bottom: 1px solid var(--el-border-color); background: var(--el-fill-color-light); font-weight: 500; }
.panel-filter { width: 230px; }
.anchor-panel :deep(.el-table) { flex: 1; }
.panel-actions { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; padding: 0 4px; }
.panel-actions .el-button { margin: 0; min-width: 96px; }
:deep(.anchor-cell) { display: flex; align-items: center; gap: 8px; min-width: 0; }
:deep(.anchor-text) { display: flex; flex-direction: column; min-width: 0; overflow: hidden; }
:deep(.anchor-text span), :deep(.anchor-text small) { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
:deep(.anchor-text small) { color: var(--el-text-color-secondary); }
</style>

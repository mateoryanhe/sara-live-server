<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-title">
            <span>{{ t('menu.ModuleList') }}</span>
            <span v-if="roleName" class="role-name">{{ t('pages.moduleList.configuringRole', {name: roleName}) }}</span>
          </div>
          <el-button type="default" @click="handleHeaderBack">
            {{ navStack.length > 1 ? t('pages.moduleList.backLevel') : t('pages.moduleList.back') }}
          </el-button>
        </div>
      </template>

      <div class="content">
        <div class="toolbar">
          <el-button type="primary" @click="handleSave">{{ t('pages.moduleList.savePermissions') }}</el-button>
          <span class="selected-count">{{ t('pages.moduleList.selectedCount', {count: checkedSet.size}) }}</span>
          <el-input
              v-model="filterText"
              clearable
              class="search-input"
              :placeholder="t('pages.moduleList.searchPlaceholder')"
          />
        </div>

        <el-breadcrumb class="breadcrumb" separator="/">
          <el-breadcrumb-item
              v-for="(frame, index) in navStack"
              :key="`${frame.title}-${index}`"
          >
            <a
                v-if="index < navStack.length - 1"
                class="breadcrumb-link"
                href="#"
                @click.prevent="goToLevel(index)"
            >
              {{ frame.title }}
            </a>
            <span v-else>{{ frame.title }}</span>
          </el-breadcrumb-item>
        </el-breadcrumb>

        <p class="hint">{{ t('pages.moduleList.permissionHint') }}</p>

        <div v-loading="loading" class="perm-list">
          <div
              v-for="item in filteredList"
              :key="item.id"
              class="perm-row"
          >
            <el-checkbox
                class="perm-checkbox"
                :indeterminate="isNodeIndeterminate(item)"
                :model-value="isNodeChecked(item)"
                @change="(val: boolean) => toggleNode(item, val)"
                @click.stop
            />
            <div
                class="perm-main"
                :class="{ clickable: hasChildren(item) }"
                @click="enterNode(item)"
            >
              <span class="perm-name" :title="item.name">{{ item.name }}</span>
              <span v-if="hasChildren(item)" class="perm-meta">
                {{ t('pages.moduleList.childCount', {count: countPermissionLeaves(item)}) }}
              </span>
            </div>
            <el-button
                v-if="hasChildren(item)"
                class="perm-enter"
                link
                type="primary"
                @click="enterNode(item)"
            >
              {{ t('pages.moduleList.enter') }}
              <el-icon><ArrowRight/></el-icon>
            </el-button>
          </div>

          <el-empty v-if="!loading && filteredList.length === 0" :description="t('pages.moduleList.emptyList')"/>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onActivated, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {ArrowRight} from '@element-plus/icons-vue'
import router from '@/router'
import {useRoute} from 'vue-router'
import {roleApi} from '@/api/modules/role'
import {buildPermissionModuleTree, type PermissionModuleNode} from '@/config/permission-menu-tree'
import {getPermissionApiPath} from '@/config/permission-api-paths'
import {isButtonPermissionKey, pageHasGranularButtons} from '@/config/page-buttons'

interface NavFrame {
  title: string
  nodes: PermissionModuleNode[]
}

const {t} = useI18n()
const route = useRoute()
const loading = ref(false)
const moduleTreeData = ref<PermissionModuleNode[]>([])
const checkedSet = ref(new Set<string>())
const filterText = ref('')

const roleId = ref(0)
const roleName = ref('')
const loadedRoleKey = ref('')

const parseRoleId = (): number => {
  const raw = route.query.roleId
  if (raw == null || raw === '') {
    return 0
  }
  const id = Number(Array.isArray(raw) ? raw[0] : raw)
  return Number.isFinite(id) ? id : 0
}

const parseRoleName = (): string => {
  const raw = route.query.roleName
  return typeof raw === 'string' ? raw : ''
}

const buildRoleKey = () => `${parseRoleId()}::${parseRoleName()}`

const resetNavStack = () => {
  navStack.value = [{
    title: t('pages.moduleList.rootTitle'),
    nodes: moduleTreeData.value,
  }]
  filterText.value = ''
}

const loadRolePermissions = async () => {
  const roleKey = buildRoleKey()
  if (roleKey === loadedRoleKey.value) {
    return
  }

  const nextRoleId = parseRoleId()
  roleId.value = nextRoleId
  roleName.value = parseRoleName()
  loadedRoleKey.value = roleKey

  if (!nextRoleId) {
    checkedSet.value = new Set<string>()
    resetNavStack()
    return
  }

  loading.value = true
  try {
    if (moduleTreeData.value.length === 0) {
      moduleTreeData.value = buildPermissionModuleTree(key => t(key))
    }
    resetNavStack()

    const permissions = await roleApi.getRolePermissionList(nextRoleId)
    const next = new Set<string>()
    permissions
        .map(p => p.module)
        .filter(mod => isPermissionId(mod))
        .forEach(mod => next.add(mod))
    checkedSet.value = next
  } finally {
    loading.value = false
  }
}

const navStack = ref<NavFrame[]>([{title: '', nodes: []}])

const currentList = computed(() => navStack.value[navStack.value.length - 1]?.nodes ?? [])

const filteredList = computed(() => {
  const keyword = filterText.value.trim().toLowerCase()
  if (!keyword) {
    return currentList.value
  }
  return currentList.value.filter(item => item.name.toLowerCase().includes(keyword))
})

const isFolderId = (id: string) => id.startsWith('module_') || id.startsWith('permission_slice_')

const isPermissionId = (id: string) => !isFolderId(id)

const hasChildren = (node: PermissionModuleNode) => (node.children?.length ?? 0) > 0

const collectPermissionIds = (node: PermissionModuleNode): string[] => {
  if (!hasChildren(node)) {
    return isPermissionId(node.id) ? [node.id] : []
  }
  return node.children!.flatMap(child => collectPermissionIds(child))
}

const countPermissionLeaves = (node: PermissionModuleNode): number => collectPermissionIds(node).length

/** 库中可能存页面级 module（如 LiveRoomRecycleBinManagement），需与按钮子权限一并参与勾选判断 */
const hasPageLevelPermission = (node: PermissionModuleNode): boolean =>
    isPermissionId(node.id) && checkedSet.value.has(node.id)

const isNodeChecked = (node: PermissionModuleNode): boolean => {
  if (!hasChildren(node)) {
    return checkedSet.value.has(node.id)
  }
  const ids = collectPermissionIds(node)
  return ids.length > 0 && ids.every(id => checkedSet.value.has(id))
}

const isNodeIndeterminate = (node: PermissionModuleNode): boolean => {
  if (!hasChildren(node)) {
    return false
  }
  const ids = collectPermissionIds(node)
  const checkedCount = ids.filter(id => checkedSet.value.has(id)).length
  if (checkedCount > 0 && checkedCount < ids.length) {
    return true
  }
  return hasPageLevelPermission(node) && checkedCount < ids.length
}

/** 有按钮子权限的页面只持久化按钮级 module，去掉遗留的页面级 module */
const normalizeModulesForSave = (modules: string[]): string[] =>
    modules.filter(moduleId => isButtonPermissionKey(moduleId) || !pageHasGranularButtons(moduleId))

const toggleNode = (node: PermissionModuleNode, checked: boolean) => {
  const next = new Set(checkedSet.value)
  const leafIds = hasChildren(node)
      ? collectPermissionIds(node).filter(isPermissionId)
      : [node.id].filter(isPermissionId)
  if (checked) {
    leafIds.forEach(id => next.add(id))
  } else {
    if (isPermissionId(node.id)) {
      next.delete(node.id)
    }
    leafIds.forEach(id => next.delete(id))
  }
  checkedSet.value = next
}

const enterNode = (node: PermissionModuleNode) => {
  if (!hasChildren(node)) {
    return
  }
  filterText.value = ''
  navStack.value.push({
    title: node.name,
    nodes: node.children!,
  })
}

const goToLevel = (index: number) => {
  if (index < 0 || index >= navStack.value.length - 1) {
    return
  }
  filterText.value = ''
  navStack.value = navStack.value.slice(0, index + 1)
}

const handleHeaderBack = () => {
  if (navStack.value.length > 1) {
    navStack.value.pop()
    filterText.value = ''
    return
  }
  router.go(-1)
}

const handleSave = async () => {
  try {
    const selectedModules = normalizeModulesForSave([...checkedSet.value])

    const permissionData = selectedModules.map(moduleId => ({
      id: 0,
      module: moduleId,
      roleId: roleId.value,
      apiPath: getPermissionApiPath(moduleId),
    }))

    const response = await roleApi.createPermission(permissionData)

    if (response) {
      ElMessage.success(t('pages.moduleList.saveSuccess', {count: selectedModules.length}))
      loadedRoleKey.value = ''
    } else {
      ElMessage.error(t('pages.moduleList.saveFailed'))
    }
  } catch (error) {
    console.error('Failed to save permissions:', error)
    ElMessage.error(t('pages.moduleList.saveFailed'))
  }
}

// keep-alive 按 path 缓存，query 变化不会重新挂载
watch(
    () => [route.query.roleId, route.query.roleName] as const,
    () => {
      void loadRolePermissions()
    },
    {immediate: true},
)

onActivated(() => {
  void loadRolePermissions()
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

.header-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.role-name {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  font-weight: normal;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.selected-count {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.search-input {
  width: 260px;
  margin-left: auto;
}

.breadcrumb {
  margin-bottom: 12px;
}

.breadcrumb-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.breadcrumb-link:hover {
  text-decoration: underline;
}

.hint {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  margin: 0 0 16px;
}

.perm-list {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  overflow: hidden;
}

.perm-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.perm-checkbox {
  flex-shrink: 0;
}

.perm-main {
  flex: 0 1 auto;
  max-width: min(360px, 50vw);
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.perm-enter {
  flex-shrink: 0;
  margin-left: 4px;
}

.perm-name {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.perm-row:last-child {
  border-bottom: none;
}

.perm-row:hover {
  background: var(--el-fill-color-light);
}

.perm-main.clickable {
  cursor: pointer;
}

.perm-meta {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.content {
  padding: 20px 0;
}
</style>

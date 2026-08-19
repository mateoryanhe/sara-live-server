<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.RoleManagement') }}</span>
          <el-button style="margin-left: auto;" type="default" @click="handleBack">{{ t('pages.moduleList.back') }}</el-button>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleSave">{{ t('pages.moduleList.savePermissions') }}</el-button>
          <span class="hint">{{ t('pages.moduleList.permissionHint') }}</span>
        </div>

        <el-tree
            ref="treeRef"
            :data="moduleTreeData"
            :default-checked-keys="checkedModules"
            :props="treeProps"
            default-expand-all
            node-key="id"
            show-checkbox
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {onMounted, ref} from 'vue'
import {ElMessage} from 'element-plus'
import router from '@/router'
import {useRoute} from 'vue-router'
import {roleApi} from '@/api/modules/role'
import {buildPermissionModuleTree} from '@/config/permission-menu-tree'
import {getPermissionApiPath} from '@/config/permission-api-paths'

interface ModuleNode {
  id: string
  name: string
  children?: ModuleNode[]
}

const {t} = useI18n()
const route = useRoute()
const treeRef = ref()
const moduleTreeData = ref<ModuleNode[]>([])
const checkedModules = ref<string[]>([])

const roleId = ref(Number(route.query.roleId) || 0)

const treeProps = {
  children: 'children',
  label: 'name',
}

const isPermissionModuleKey = (key: string) => {
  const value = String(key)
  return value.startsWith('module_') || value.startsWith('permission_slice_')
}

const collectSelectedPermissions = (): string[] => {
  const checkedLeaves = treeRef.value.getCheckedKeys(true) as string[]
  const checkedAll = treeRef.value.getCheckedKeys(false) as string[]

  const pageKeys = checkedAll.filter(key => {
    const value = String(key)
    return !isPermissionModuleKey(value) && !value.includes(':')
  })
  const buttonKeys = checkedLeaves.filter(key => String(key).includes(':'))

  return [...new Set([...pageKeys, ...buttonKeys])]
}

const handleSave = async () => {
  try {
    const selectedModules = collectSelectedPermissions()

    const permissionData = selectedModules.map(moduleId => ({
      id: 0,
      module: moduleId,
      roleId: roleId.value,
      apiPath: getPermissionApiPath(moduleId),
    }))

    const response = await roleApi.createPermission(permissionData)

    if (response) {
      ElMessage.success(t('pages.moduleList.saveSuccess', {count: selectedModules.length}))
    } else {
      ElMessage.error(t('pages.moduleList.saveFailed'))
    }
  } catch (error) {
    console.error('Failed to save permissions:', error)
    ElMessage.error(t('pages.moduleList.saveFailed'))
  }
}

const handleBack = () => {
  router.go(-1)
}

onMounted(async () => {
  moduleTreeData.value = buildPermissionModuleTree(t)

  const permissions = await roleApi.getRolePermissionList(roleId.value)
  const permissionModules = permissions.map(p => p.module)

  checkedModules.value = permissionModules.filter(
      mod => !mod.startsWith('module_'),
  )
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
}

.table-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.hint {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.content {
  padding: 20px 0;
}
</style>

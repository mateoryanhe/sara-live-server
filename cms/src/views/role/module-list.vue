<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>权限配置</span>
          <el-button style="margin-left: auto;" type="default" @click="handleBack">返回</el-button>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleSave">保存权限配置</el-button>
          <span class="hint">勾选页面可授予整页权限；展开后可单独勾选各按钮</span>
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
import {onMounted, ref} from 'vue'
import {ElMessage} from 'element-plus'
import router from '@/router'
import {layoutRouteGroups} from '@/router/routes'
import {useRoute} from 'vue-router'
import {roleApi} from '@/api/modules/role'
import {
  buttonPermissionKey,
  getPageButtons,
  type PageButtonDef,
} from '@/config/page-buttons'

interface ModuleNode {
  id: string
  name: string
  children?: ModuleNode[]
}

const route = useRoute()
const treeRef = ref()
const moduleTreeData = ref<ModuleNode[]>([])
const checkedModules = ref<string[]>([])

const roleId = ref(Number(route.query.roleId) || 0)

const treeProps = {
  children: 'children',
  label: 'name',
}

const buildButtonNodes = (pageName: string, metaButtons?: PageButtonDef[]): ModuleNode[] => {
  return getPageButtons(pageName, metaButtons).map(btn => ({
    id: buttonPermissionKey(pageName, btn.key),
    name: btn.label,
  }))
}

const generateModuleTreeFromRoutes = () => {
  const routeModules: ModuleNode[] = []

  layoutRouteGroups.forEach(group => {
    if (group.children && group.meta) {
      const module: ModuleNode = {
        id: `module_${group.path.replace('/', '')}`,
        name: (group.meta.title as string) || group.path,
        children: [],
      }

      group.children.forEach(child => {
        if (child.name && child.meta && !child.meta.hidden) {
          const pageName = child.name as string
          const metaButtons = child.meta.buttons as PageButtonDef[] | undefined
          const buttonNodes = buildButtonNodes(pageName, metaButtons)
          module.children?.push({
            id: pageName,
            name: (child.meta.title as string) || pageName,
            children: buttonNodes.length > 0 ? buttonNodes : undefined,
          })
        }
      })

      if (module.children && module.children.length > 0) {
        routeModules.push(module)
      }
    }
  })

  return routeModules
}

const collectSelectedPermissions = (): string[] => {
  const checkedLeaves = treeRef.value.getCheckedKeys(true) as string[]
  const checkedAll = treeRef.value.getCheckedKeys(false) as string[]

  const pageKeys = checkedAll.filter(
      key => !String(key).startsWith('module_') && !String(key).includes(':'),
  )
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
    }))

    const response = await roleApi.createPermission(permissionData)

    if (response) {
      ElMessage.success(`已保存 ${selectedModules.length} 项权限`)
    } else {
      ElMessage.error('保存权限配置失败')
    }
  } catch (error) {
    console.error('保存权限配置失败:', error)
    ElMessage.error('保存权限配置失败')
  }
}

const handleBack = () => {
  router.go(-1)
}

onMounted(async () => {
  moduleTreeData.value = generateModuleTreeFromRoutes()

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

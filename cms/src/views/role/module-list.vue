<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>模块管理</span>
          <el-button style="margin-left: auto;" type="default" @click="handleBack">返回</el-button>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleSave">保存模块配置</el-button>
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
import {useRoute} from 'vue-router'
// 从role API模块导入Permission接口
import {roleApi} from '@/api/modules/role'

interface ModuleNode {
  id: string
  name: string
  children?: ModuleNode[]
}

const route = useRoute()
const treeRef = ref()
// 从路由配置中动态生成模块数据
const moduleTreeData = ref<ModuleNode[]>([])

const checkedModules = ref<string[]>([])

// 当前角色ID，从路由参数获取或默认为0
const roleId = ref(Number(route.query.roleId) || 0)

const treeProps = {
  children: 'children',
  label: 'name'
}

// 从路由配置中提取模块数据
const generateModuleTreeFromRoutes = () => {
  const routeModules: ModuleNode[] = []

  // 遍历路由配置，提取模块信息
  router.options.routes.forEach(route => {
    // 跳过登录页面等特殊路由
    if (route.path === '/login' || route.path === '/') {
      return
    }

    // 检查路由是否有meta信息和子路由
    if (route.children && route.meta) {
      const module: ModuleNode = {
        id: `module_${route.path.replace('/', '')}`,
        name: route.meta.title as string || route.path,
        children: []
      }

      // 遍历子路由，添加为模块的子节点
      route.children.forEach(child => {
        if (child.name && child.meta && !child.meta.hidden) {
          module.children?.push({
            id: child.name as string,
            name: child.meta.title as string || child.name as string
          })
        }
      })

      // 只有当模块有子节点时才添加
      if (module.children && module.children.length > 0) {
        routeModules.push(module)
      }
    }
  })

  return routeModules
}

const handleSave = async () => {
  try {
    const checkedKeys = treeRef.value.getCheckedKeys(false)
    const halfCheckedKeys = treeRef.value.getHalfCheckedKeys()
    const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]

    // 过滤掉模块节点，只保留具体功能节点
    const selectedModules = allCheckedKeys.filter(key => !key.startsWith('module_'))

    // 准备保存数据到后端
    const permissionData = selectedModules.map(moduleId => ({
      id: 0, // 新增权限时ID为0
      module: moduleId,
      roleId: roleId.value
    }))
    
    console.log('选中的模块:', selectedModules)

    // 调用API创建或更新角色权限
    const response = await roleApi.createPermission(permissionData)

    console.log('权限已更新:', selectedModules)

    if (response) {
      ElMessage.success(`已保存 ${selectedModules.length} 个模块权限`)
    } else {
      ElMessage.error('保存模块配置失败')
    }
  } catch (error) {
    console.error('保存模块配置失败:', error)
    ElMessage.error('保存模块配置失败')
  }
}

const handleBack = () => {
  // 返回到角色管理页面或其他上级页面
  router.go(-1) // 或者使用 router.push('/role') 返回角色管理页面
}

onMounted(async () => {
  // 从路由配置中动态生成模块树
  moduleTreeData.value = generateModuleTreeFromRoutes()

  // 从接口获取当前角色的权限列表
  const response = await roleApi.getRolePermissionList(roleId.value)
  const permissions = response

  // 根据接口返回的权限列表设置默认勾选项
  // ModuleNode.id 必须在响应结果里面才显示勾选，使用响应结果的 module 字段比较
  checkedModules.value = []

  // 提取接口返回的模块ID列表
  const permissionModules = permissions.map(p => p.module)

  // 遍历模块树数据，只有当模块ID在权限列表中时才设置为选中
  moduleTreeData.value.forEach(module => {
    if (module.children) {
      module.children.forEach(child => {
        if (permissionModules.includes(child.id)) {
          checkedModules.value.push(child.id)
        }
      })
    }
  })
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

.content {
  padding: 20px 0;
}
</style>
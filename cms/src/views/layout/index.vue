<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside class="aside" width="220px">
      <div class="logo">XR Game Server</div>
      <el-menu
          :collapse="isCollapse"
          :default-active="activeMenu"
          :router="true"
          :unique-opened="true"
          class="sidebar-menu"
      >
        <!-- 移除仪表盘菜单项 -->
        <el-sub-menu
            v-if="hasMenuPermission('UserList') || hasMenuPermission('GuildManagement')"
            index="/account">
          <template #title>
            <el-icon>
              <User/>
            </el-icon>
            <span>运营管理</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('UserList')" index="/account/user-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>用户列表</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GuildManagement')" index="/guild/guild-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>工会管理</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('GiftManagement')" index="/gift/gift-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>礼物管理</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu v-if="hasMenuPermission('GlobalConfig')" index="/config">
          <template #title>
            <el-icon>
              <Setting/>
            </el-icon>
            <span>系统配置</span>
          </template>
          <el-menu-item index="/config/global">
            <el-icon>
              <Monitor/>
            </el-icon>
            <span>全局配置</span>
          </el-menu-item>
        </el-sub-menu>
        <!-- 角色权限管理菜单 -->
        <el-sub-menu
            v-if="hasMenuPermission('RoleManagement') || hasMenuPermission('ModuleList') || hasMenuPermission('CMSUserManagement')"
            index="/role">
          <template #title>
            <el-icon>
              <Lock/>
            </el-icon>
            <span>角色权限</span>
          </template>
          <el-menu-item v-if="hasMenuPermission('RoleManagement')" index="/role/role-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>角色权限管理</span>
          </el-menu-item>
          <el-menu-item v-if="hasMenuPermission('CMSUserManagement')" index="/role/cmsuser-list">
            <el-icon>
              <User/>
            </el-icon>
            <span>CMS用户管理</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 主内容区域 -->
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button class="collapse-btn" @click="toggleCollapse">
            <el-icon>
              <Fold v-if="!isCollapse"/>
              <Expand v-else/>
            </el-icon>
          </el-button>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ username }}
              <el-icon class="el-icon--right">
                <arrow-down/>
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view/>
      </el-main>
      <el-footer class="footer">
        <div class="footer-content">
          XR Game Server 管理系统 &copy; {{ new Date().getFullYear() }}
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import {computed, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ArrowDown, Expand, Fold, Lock, Monitor, Present, Setting, User} from '@element-plus/icons-vue'
import {clearPermissions, getIsAdmin, hasPermission} from '@/utils/permission'

const route = useRoute()
const router = useRouter()
const isCollapse = ref(false)

const activeMenu = computed(() => {
  const {path} = route
  return path
})

const username = computed(() => {
  // 从localStorage获取用户名，如果不存在则显示默认值
  return localStorage.getItem('username') || '管理员'
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const logout = () => {
  // 清除登录信息
  localStorage.removeItem('token')
  localStorage.removeItem('authId')
  localStorage.removeItem('username')

  // 清除权限信息
  clearPermissions()

  // 跳转到登录页
  router.push('/login')
}

// 检查菜单项是否有权限
const hasMenuPermission = (moduleName: string) => {
  // 管理员拥有所有权限
  if (getIsAdmin()) {
    return true
  }

  // 检查是否有访问该模块的权限
  return hasPermission(moduleName)
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  background-color: #ffffff;
  color: #333;
  height: 100vh;
  overflow: hidden;
  box-shadow: 1px 0 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #333;
  font-size: 18px;
  font-weight: 600;
  border-bottom: 1px solid #f0f0f0;
  background: #fff;
  transition: all 0.3s ease;
}

.sidebar-menu {
  border: none;
  height: calc(100% - 60px);
  background-color: #ffffff;
  overflow-y: auto;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  margin-right: 20px;
  font-size: 16px;
  color: #409eff;
}

.main-content {
  background-color: #f5f7f8;
  padding: 20px;
  overflow-y: auto;
}

.footer {
  background-color: #fafafa;
  border-top: 1px solid #e6e6e6;
  padding: 0;
  height: 50px;
}

.footer-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  display: flex;
  align-items: center;
  font-size: 14px;
}

/* 简约风格的菜单样式 */
.el-menu {
  border-right: none;
  background-color: #ffffff;
}

.el-sub-menu__title {
  color: #333;
  font-size: 14px;
  height: 48px;
  line-height: 48px;
  display: flex;
  align-items: center;
  padding-left: 20px !important;
  background-color: #fff;
}

.el-sub-menu__title:hover {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

.el-menu-item {
  color: #666;
  font-size: 14px;
  height: 42px;
  line-height: 42px;
  display: flex;
  align-items: center;
  padding-left: 50px !important;
  background-color: #fff;
}

.el-menu-item:hover {
  background-color: #f5f7fa !important;
  color: #409eff !important;
}

.el-menu-item.is-active {
  background-color: #ecf5ff !important;
  color: #409eff !important;
  border-left: 3px solid #409eff;
  font-weight: 500;
}

/* 子菜单项样式 */
.el-sub-menu .el-menu-item {
  color: #666;
  padding-left: 65px !important;
  height: 40px;
  line-height: 40px;
}

.el-sub-menu .el-menu-item:hover {
  color: #409eff !important;
}

.el-sub-menu .el-menu-item.is-active {
  color: #409eff !important;
}

/* 菜单收起时的样式 */
.el-aside:not(.el-menu--collapse) {
  width: 220px;
  height: 100%;
}

/* 滚动条样式 */
.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background-color: #e0e0e0;
  border-radius: 2px;
}

.sidebar-menu::-webkit-scrollbar-track {
  background-color: #fff;
}
</style>
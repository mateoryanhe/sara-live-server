<template>
  <el-container
      v-loading="layoutBootstrapping"
      class="layout-container"
      element-loading-background="rgba(255, 255, 255, 0.92)"
      :element-loading-text="t('common.loading')"
  >
    <component
        :is="layoutSidebarComponent"
        v-if="layoutShellReady && layoutSidebarComponent"
        ref="layoutSidebarRef"
        :active-menu="activeMenu"
        :is-collapse="isCollapse"
    />

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
          <LanguageSwitcher class="header-lang" compact show-label/>
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ username }}
              <el-icon class="el-icon--right">
                <arrow-down/>
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="logout">{{ t('common.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <LayoutTabs/>
      <el-main
          v-loading="routeNavigating"
          class="main-content"
          element-loading-background="rgba(255, 255, 255, 0.72)"
          :element-loading-text="t('common.loading')"
      >
        <router-view v-slot="{ Component, route: currentRoute }">
          <keep-alive :max="15">
            <component :is="Component" v-if="Component" :key="currentRoute.fullPath"/>
          </keep-alive>
        </router-view>
      </el-main>
      <el-footer class="footer">
        <div class="footer-content">
          {{ t('common.footer') }} &copy; {{ new Date().getFullYear() }}
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import LayoutTabs from '@/components/layout/LayoutTabs.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'
import {useLayoutTabs} from '@/composables/useLayoutTabs'
import {ArrowDown, Expand, Fold} from '@element-plus/icons-vue'
import {clearAuthSession} from '@/utils/auth'
import {routeNavigating} from '@/utils/route-loading'
import {
  layoutBootstrapping,
  layoutShellReady,
  layoutSidebarComponent,
  runLayoutBootstrap,
} from '@/utils/layout-bootstrap'

const route = useRoute()
const router = useRouter()
const {t} = useI18n()
const isCollapse = ref(false)
const {addTab} = useLayoutTabs()

type LayoutSidebarExpose = { openSubmenusForPath: (path: string) => Promise<void> }
const layoutSidebarRef = ref<LayoutSidebarExpose | null>(null)

const syncRouteParentMenu = async () => {
  await layoutSidebarRef.value?.openSubmenusForPath(route.path)
}

watch(
    () => route.fullPath,
    () => {
      addTab(route)
      if (layoutShellReady.value) {
        void syncRouteParentMenu()
      }
    },
    {immediate: true},
)

watch(layoutShellReady, (ready) => {
  if (ready) {
    void syncRouteParentMenu()
  }
})

onMounted(() => {
  void runLayoutBootstrap()
})

const activeMenu = computed(() => route.path)

const username = computed(() => localStorage.getItem('username') || t('common.admin'))

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const logout = () => {
  clearAuthSession()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
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

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-lang {
  margin-right: 4px;
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
</style>


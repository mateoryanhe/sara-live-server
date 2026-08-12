<template>
  <div class="layout-tabs">
    <el-tabs
        v-model="activeTab"
        class="layout-tabs__list"
        type="card"
        @contextmenu.prevent="handleTabContextMenu"
        @tab-click="handleTabClick"
        @tab-remove="handleTabRemove"
    >
      <el-tab-pane
          v-for="tab in tabs"
          :key="tab.path"
          :closable="tab.path !== '/dashboard'"
          :label="tabLabel(tab)"
          :name="tab.path"
      />
    </el-tabs>
    <div class="layout-tabs__actions">
      <el-dropdown trigger="click" @command="handleCommand">
        <el-button class="layout-tabs__menu-btn" size="small" text>
          {{ t('tabs.actions') }}
          <el-icon class="el-icon--right">
            <ArrowDown/>
          </el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="closeOthers">{{ t('tabs.closeOthers') }}</el-dropdown-item>
            <el-dropdown-item command="closeAll">{{ t('tabs.closeAll') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <ul
        v-show="contextMenu.visible"
        class="tab-context-menu"
        :style="{left: `${contextMenu.x}px`, top: `${contextMenu.y}px`}"
        @click.stop
    >
      <li
          :class="{ disabled: !canCloseContextTab }"
          @click="handleContextAction('close')"
      >
        {{ t('tabs.close') }}
      </li>
      <li
          :class="{ disabled: !canCloseOthers }"
          @click="handleContextAction('closeOthers')"
      >
        {{ t('tabs.closeOthers') }}
      </li>
      <li @click="handleContextAction('closeAll')">{{ t('tabs.closeAll') }}</li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ArrowDown} from '@element-plus/icons-vue'
import type {TabsPaneContext} from 'element-plus'
import {useLayoutTabs, type LayoutTab} from '@/composables/useLayoutTabs'

const DASHBOARD_PATH = '/dashboard'

const route = useRoute()
const router = useRouter()
const {t, te} = useI18n()
const {tabs, removeTab, closeOtherTabs, closeAllTabs} = useLayoutTabs()

const tabLabel = (tab: LayoutTab) => {
  if (tab.name && te(`menu.${tab.name}`)) {
    return t(`menu.${tab.name}`)
  }
  return tab.title
}

const contextTabPath = ref('')
const contextMenu = reactive({
  visible: false,
  x: 0,
  y: 0,
})

const activeTab = computed({
  get: () => route.path,
  set: (path: string) => {
    if (path !== route.path) {
      void router.push(path)
    }
  },
})

const canCloseContextTab = computed(() => contextTabPath.value !== DASHBOARD_PATH)

const canCloseOthers = computed(() => {
  if (!contextTabPath.value) {
    return false
  }
  return tabs.value.some(
      tab => tab.path !== contextTabPath.value && tab.path !== DASHBOARD_PATH,
  )
})

const hideContextMenu = () => {
  contextMenu.visible = false
}

const resolveTabPathFromEvent = (event: MouseEvent): string => {
  const tabEl = (event.target as HTMLElement).closest('.el-tabs__item')
  if (!tabEl) {
    return ''
  }

  const tabId = tabEl.getAttribute('id')
  if (!tabId?.startsWith('tab-')) {
    return ''
  }

  return tabId.slice(4)
}

const handleTabContextMenu = (event: MouseEvent) => {
  const path = resolveTabPathFromEvent(event)
  if (!path) {
    return
  }

  contextTabPath.value = path
  contextMenu.x = event.clientX
  contextMenu.y = event.clientY
  contextMenu.visible = true
}

const handleTabClick = (pane: TabsPaneContext) => {
  hideContextMenu()
  const path = String(pane.paneName)
  if (path !== route.path) {
    void router.push(path)
  }
}

const handleTabRemove = (path: string | number) => {
  removeTab(String(path), router)
}

const handleCommand = (command: 'closeOthers' | 'closeAll') => {
  if (command === 'closeOthers') {
    closeOtherTabs(route.path, router)
    return
  }
  closeAllTabs(router)
}

const handleContextAction = (action: 'close' | 'closeOthers' | 'closeAll') => {
  const path = contextTabPath.value
  hideContextMenu()

  if (!path) {
    return
  }

  if (action === 'close') {
    if (path !== DASHBOARD_PATH) {
      removeTab(path, router)
    }
    return
  }

  if (action === 'closeOthers') {
    if (canCloseOthers.value) {
      closeOtherTabs(path, router)
    }
    return
  }

  closeAllTabs(router)
}

const handleDocumentClick = () => {
  hideContextMenu()
}

const handleDocumentContextMenu = (event: MouseEvent) => {
  const menuEl = (event.target as HTMLElement).closest('.tab-context-menu')
  if (!menuEl) {
    hideContextMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
  document.addEventListener('contextmenu', handleDocumentContextMenu)
})

onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick)
  document.removeEventListener('contextmenu', handleDocumentContextMenu)
})
</script>

<style scoped>
.layout-tabs {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px 0;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
}

.layout-tabs__list {
  flex: 1;
  min-width: 0;
}

.layout-tabs__actions {
  flex-shrink: 0;
  padding-bottom: 4px;
}

.layout-tabs__menu-btn {
  color: #606266;
}

.tab-context-menu {
  position: fixed;
  z-index: 3000;
  min-width: 120px;
  margin: 0;
  padding: 4px 0;
  list-style: none;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.12);
}

.tab-context-menu li {
  padding: 8px 16px;
  font-size: 13px;
  color: #606266;
  cursor: pointer;
  user-select: none;
}

.tab-context-menu li:hover:not(.disabled) {
  background-color: #f5f7fa;
  color: #409eff;
}

.tab-context-menu li.disabled {
  color: #c0c4cc;
  cursor: not-allowed;
}

:deep(.el-tabs__header) {
  margin: 0;
}

:deep(.el-tabs__nav-wrap) {
  margin-bottom: 0;
}

:deep(.el-tabs__item) {
  height: 32px;
  line-height: 32px;
  font-size: 13px;
}

:deep(.el-tabs__content) {
  display: none;
}
</style>

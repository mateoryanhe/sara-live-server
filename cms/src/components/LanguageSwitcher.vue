<template>
  <div class="language-switcher" :class="{compact}">
    <span v-if="showLabel" class="language-label">{{ t('common.language') }}</span>
    <el-dropdown trigger="click" @command="switchLocale">
      <el-button size="small" class="language-trigger">
        {{ LOCALE_LABELS[locale as CmsLocale] }}
        <el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
              v-for="item in SUPPORTED_LOCALES"
              :key="item"
              :command="item"
              :class="{active: locale === item}"
          >
            {{ LOCALE_LABELS[item] }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script lang="ts" setup>
import {ArrowDown} from '@element-plus/icons-vue'
import {useI18n} from 'vue-i18n'
import {LOCALE_LABELS, SUPPORTED_LOCALES, setAppLocale, type CmsLocale} from '@/i18n'

withDefaults(defineProps<{
  showLabel?: boolean
  compact?: boolean
}>(), {
  showLabel: false,
  compact: false,
})

const {locale, t} = useI18n()

const switchLocale = (next: CmsLocale) => {
  if (locale.value === next) {
    return
  }
  setAppLocale(next)
}
</script>

<style scoped>
.language-switcher {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.language-trigger {
  min-width: 72px;
}

.language-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

:deep(.el-dropdown-menu__item.active) {
  color: var(--el-color-primary);
  font-weight: 600;
}
</style>

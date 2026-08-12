<script lang="ts" setup>
import {computed} from 'vue'
import {RouterView} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {ElConfigProvider} from 'element-plus'
import {getElementPlusLocale} from '@/i18n'
import type {CmsLocale} from '@/i18n'

const {locale} = useI18n()

const elementLocale = computed(() => getElementPlusLocale(locale.value as CmsLocale))
</script>

<template>
  <ElConfigProvider :locale="elementLocale">
    <RouterView/>
  </ElConfigProvider>
</template>

<style>
#app {
  font-family: 'Segoe UI', system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>

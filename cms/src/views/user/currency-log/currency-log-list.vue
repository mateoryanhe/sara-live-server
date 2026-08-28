<template>
  <CurrencyLogListPanel :currency-type="currencyType" :page-name="pageName"/>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import {useRoute} from 'vue-router'
import CurrencyLogListPanel from './currency-log-list-panel.vue'

const route = useRoute()

const currencyType = computed<1 | 2>(() => (Number(route.meta.currencyType) === 2 ? 2 : 1))
const pageName = computed(() => String(route.name || 'GoldCurrencyLogList'))
</script>

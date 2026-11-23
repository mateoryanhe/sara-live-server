<template>
  <AnchorListPage
      menu-title-key="GuildAnchorListManagement"
      page-permission-name="GuildAnchorListManagement"
      :off-shelf-auto-restore="true"
      :visibility-scoped="true"
  />
</template>

<script lang="ts" setup>
import AnchorListPage from '@/views/user/anchor/anchor-list-page.vue'
</script>

<template>
  <StaticSiteDeployPage
      :deploy-api="officialSiteDeployApi"
      default-deploy-path="/home/ec2-user/cdn/official-site"
      domain-example="web.saralive.net"
      enable-chunked-zip-upload
      page-name="OfficialSiteDeployManagement"
      title-key="menu.OfficialSiteDeployManagement"
  />
</template>

<script lang="ts" setup>
import StaticSiteDeployPage from '@/components/config/static-site-deploy-page.vue'
import {officialSiteDeployApi} from '@/api/modules/official-site-deploy'
</script>

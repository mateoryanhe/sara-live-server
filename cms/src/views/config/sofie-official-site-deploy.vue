<template>
  <StaticSiteDeployPage
      :deploy-api="sofieOfficialSiteDeployApi"
      default-deploy-path="/home/ec2-user/cdn/sofie-official-site"
      domain-example="sofielive.net"
      enable-chunked-zip-upload
      enable-tls-upload
      page-name="SofieOfficialSiteDeployManagement"
      title-key="menu.SofieOfficialSiteDeployManagement"
  />
</template>

<script lang="ts" setup>
import StaticSiteDeployPage from '@/components/config/static-site-deploy-page.vue'
import {sofieOfficialSiteDeployApi} from '@/api/modules/official-site-deploy'
</script>

<template>
  <StaticSiteDeployPage
      :deploy-api="thirdPayOfficialSiteDeployApi"
      default-deploy-path="/home/ec2-user/cdn/third-pay-official-site"
      domain-example="hom.saralive.net"
      enable-chunked-zip-upload
      page-name="ThirdPayOfficialSiteDeployManagement"
      title-key="menu.ThirdPayOfficialSiteDeployManagement"
  />
</template>

<script lang="ts" setup>
import StaticSiteDeployPage from '@/components/config/static-site-deploy-page.vue'
import {thirdPayOfficialSiteDeployApi} from '@/api/modules/official-site-deploy'
</script>

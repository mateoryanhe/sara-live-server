<template>
  <UserAppSubmissionRecordList
      :fetch-list-api="userReportRecordApi.listUserReportRecord"
      menu-title-key="menu.UserReportRecordManagement"
      page-desc-key="pages.userReportRecord.pageDesc"
  />
</template>

<script lang="ts" setup>
import UserAppSubmissionRecordList from './components/user-app-submission-record-list.vue'
import {userReportRecordApi} from '@/api/modules/user-report-record'
</script>

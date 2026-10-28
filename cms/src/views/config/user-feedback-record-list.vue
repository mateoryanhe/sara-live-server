<template>
  <UserAppSubmissionRecordList
      :fetch-list-api="userFeedbackRecordApi.listUserFeedbackRecord"
      menu-title-key="menu.UserFeedbackRecordManagement"
      page-desc-key="pages.userFeedbackRecord.pageDesc"
  />
</template>

<script lang="ts" setup>
import UserAppSubmissionRecordList from './components/user-app-submission-record-list.vue'
import {userFeedbackRecordApi} from '@/api/modules/user-feedback-record'
</script>

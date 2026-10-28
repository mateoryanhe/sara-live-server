import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: 'App POST /userInfo/report（Sofie /sofie/profile/submitComplaint）上报原文，syndb 入库表 user_report_records。',
}

export const userReportRecordMessages = definePageMessagesFromEn(zh, {
  pageDesc: 'Raw App body from POST /userInfo/report (Sofie submitComplaint), table user_report_records.',
})

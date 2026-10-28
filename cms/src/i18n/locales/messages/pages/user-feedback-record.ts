import {definePageMessagesFromEn} from './_define'

const zh = {
  pageDesc: 'App POST /userInfo/feedback（Sofie /sofie/profile/submitOpinion）上报原文，syndb 入库表 user_feedback_records。',
}

export const userFeedbackRecordMessages = definePageMessagesFromEn(zh, {
  pageDesc: 'Raw App body from POST /userInfo/feedback (Sofie submitOpinion), table user_feedback_records.',
})

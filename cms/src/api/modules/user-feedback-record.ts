import request from '../request'
import type {CMSUserAppSubmissionRecordListQuery, CMSUserAppSubmissionRecordItem, PageResponse} from '@/types/api'

export const userFeedbackRecordApi = {
    listUserFeedbackRecord: (data: CMSUserAppSubmissionRecordListQuery) => {
        return request.post<PageResponse<CMSUserAppSubmissionRecordItem>>('/userFeedbackRecord/listUserFeedbackRecord', {
            ...data,
            startTime: data.startTime || 0,
            endTime: data.endTime || 0,
        })
    },
}

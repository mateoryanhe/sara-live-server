import request from '../request'
import type {CMSUserAppSubmissionRecordListQuery, CMSUserAppSubmissionRecordItem, PageResponse} from '@/types/api'

export const userReportRecordApi = {
    listUserReportRecord: (data: CMSUserAppSubmissionRecordListQuery) => {
        return request.post<PageResponse<CMSUserAppSubmissionRecordItem>>('/userReportRecord/listUserReportRecord', {
            ...data,
            startTime: data.startTime || 0,
            endTime: data.endTime || 0,
        })
    },
}

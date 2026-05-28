import {request} from '../request'
import type {LiveRecordItem, LiveRecordQuery, PageResponse} from '@/types/api'

export const liveRecordApi = {
    getLiveRecordList: (params: LiveRecordQuery) => {
        return request.post<PageResponse<LiveRecordItem>>('/liveRecord/cmsLiveRecordList', params)
    },
}

export default liveRecordApi

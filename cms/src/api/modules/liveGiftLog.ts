import {request} from '../request'
import type {LiveGiftLogItem, LiveGiftLogQuery, PageResponse} from '@/types/api'

export const liveGiftLogApi = {
    getLiveGiftLogList: (params: LiveGiftLogQuery) => {
        return request.post<PageResponse<LiveGiftLogItem>>('/liveGiftLog/cmsLiveGiftLogList', params)
    },
}

export default liveGiftLogApi

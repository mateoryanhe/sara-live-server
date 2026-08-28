import {request} from '../request'
import type {DailyEffectiveLiveQuery, GuildAnchorDailyEffectiveLiveItem, LiveRecordItem, LiveRecordQuery, PageResponse, WeeklyUnsettledLiveItem, WeeklyUnsettledLiveQuery} from '@/types/api'

export const liveRecordApi = {
    getLiveRecordList: (params: LiveRecordQuery) => {
        return request.post<PageResponse<LiveRecordItem>>('/liveRecord/cmsLiveRecordList', params)
    },
    getDailyEffectiveLiveList: (params: DailyEffectiveLiveQuery) => {
        return request.post<PageResponse<GuildAnchorDailyEffectiveLiveItem>>('/liveRecord/cmsDailyEffectiveLiveList', params)
    },
    getWeeklyUnsettledLiveList: (params: WeeklyUnsettledLiveQuery) => {
        return request.post<PageResponse<WeeklyUnsettledLiveItem>>('/liveRecord/cmsWeeklyUnsettledLiveList', params)
    },
}

export default liveRecordApi

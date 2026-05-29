import {request} from '../request'
import type {LiveRevenueLogItem, LiveRevenueLogQuery, PageResponse} from '@/types/api'

export const liveRevenueLogApi = {
    getLiveRevenueLogList: (params: LiveRevenueLogQuery) => {
        return request.post<PageResponse<LiveRevenueLogItem>>('/liveRevenueLog/cmsLiveRevenueLogList', params)
    },
}

export default liveRevenueLogApi

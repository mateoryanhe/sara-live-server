import request from '../request'
import type {HotLiveRoomJoinTrendRes} from '@/types/api'

export const trackingEventApi = {
    getHotLiveRoomJoinTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getHotLiveRoomJoinTrend', {})
    },
    getLiveFirstFrameRenderTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getLiveFirstFrameRenderTrend', {})
    },
    getHotLiveRoomLeaveTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getHotLiveRoomLeaveTrend', {})
    },
    getCall1v1InitiateTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getCall1v1InitiateTrend', {})
    },
}

export default trackingEventApi

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
    getGameLiveRoomJoinTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getGameLiveRoomJoinTrend', {})
    },
    getCall1v1InitiateTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getCall1v1InitiateTrend', {})
    },
    getCall1v1RoomCallTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getCall1v1RoomCallTrend', {})
    },
    getCall1v1ConnectSuccessTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getCall1v1ConnectSuccessTrend', {})
    },
    getMiniGameRoundStartTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameRoundStartTrend', {})
    },
    getMiniGameRoundResultTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameRoundResultTrend', {})
    },
    getMiniGameExposureTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameExposureTrend', {})
    },
    getMiniGameWebViewLoadTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameWebViewLoadTrend', {})
    },
}

export default trackingEventApi

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
    getGameLiveWatchStreamTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getGameLiveWatchStreamTrend', {})
    },
    getHotLiveRoomGiftClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getHotLiveRoomGiftClickTrend', {})
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
    getOneToOneRoomListTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getOneToOneRoomListTrend', {})
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
    getMiniGameControllerClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameControllerClickTrend', {})
    },
    getMiniGameSpinClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGameSpinClickTrend', {})
    },
    getMiniGamePlusClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getMiniGamePlusClickTrend', {})
    },
    getOfflineGameRoomPlayTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getOfflineGameRoomPlayTrend', {})
    },
    getShortVideoPlayStartTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getShortVideoPlayStartTrend', {})
    },
    getShortVideoPaywallClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getShortVideoPaywallClickTrend', {})
    },
    getShortVideoPaySuccessTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getShortVideoPaySuccessTrend', {})
    },
    getShortVideoWatchEndTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getShortVideoWatchEndTrend', {})
    },
    getVoiceRoomEnterTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getVoiceRoomEnterTrend', {})
    },
    getVoiceRoomMicRequestTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getVoiceRoomMicRequestTrend', {})
    },
    getVoiceRoomGiftClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getVoiceRoomGiftClickTrend', {})
    },
    getChargePanelEnterTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getChargePanelEnterTrend', {})
    },
    getChargeClickTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getChargeClickTrend', {})
    },
    getChargeResultServerTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getChargeResultServerTrend', {})
    },
    getH5DownloadTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getH5DownloadTrend', {})
    },
    getOfficialSiteDownloadTrend: () => {
        return request.post<HotLiveRoomJoinTrendRes>('/trackingEvent/getOfficialSiteDownloadTrend', {})
    },
}

export default trackingEventApi

import {request} from '../request'
import type {BatchBotAnchorLiveRes, BatchStartBotAnchorLiveReq, BatchStopBotAnchorLiveReq, BotAnchorListItem, CreateBotAnchorReq, PageResponse, QueryBotAnchorListReq, SetBotAnchorStatusReq, StartBotAnchorLiveReq, StopBotAnchorLiveReq, UpdateBotAnchorReq} from '@/types/api'

export const botAnchorApi = {
    getBotAnchorList: (data: QueryBotAnchorListReq) => {
        return request.post<PageResponse<BotAnchorListItem>>('/botAnchor/getBotAnchorList', data)
    },

    createBotAnchor: (data: CreateBotAnchorReq) => {
        return request.post<{ id: string }>('/botAnchor/createBotAnchor', data)
    },

    updateBotAnchor: (data: UpdateBotAnchorReq) => {
        return request.post<{ success: boolean }>('/botAnchor/updateBotAnchor', data)
    },

    setBotAnchorStatus: (data: SetBotAnchorStatusReq) => {
        return request.post<{ success: boolean }>('/botAnchor/setBotAnchorStatus', data)
    },

    startBotAnchorLive: (data: StartBotAnchorLiveReq) => {
        return request.post<{ success: boolean }>('/botAnchor/startBotAnchorLive', data)
    },

    stopBotAnchorLive: (data: StopBotAnchorLiveReq) => {
        return request.post<{ success: boolean }>('/botAnchor/stopBotAnchorLive', data)
    },

    batchStartBotAnchorLive: (data: BatchStartBotAnchorLiveReq) => {
        return request.post<BatchBotAnchorLiveRes>('/botAnchor/batchStartBotAnchorLive', data)
    },

    batchStopBotAnchorLive: (data: BatchStopBotAnchorLiveReq) => {
        return request.post<BatchBotAnchorLiveRes>('/botAnchor/batchStopBotAnchorLive', data)
    },
}

export default botAnchorApi

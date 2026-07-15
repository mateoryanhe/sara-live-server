import {request} from '../request'
import type {BotAnchorListItem, CreateBotAnchorReq, PageResponse, QueryBotAnchorListReq, SetBotAnchorStatusReq, UpdateBotAnchorReq} from '@/types/api'

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
}

export default botAnchorApi

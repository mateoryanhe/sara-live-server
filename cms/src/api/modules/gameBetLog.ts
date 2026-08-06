import {request} from '../request'
import type {GameBetLogItem, GameBetLogQuery, PageResponse} from '@/types/api'

export const gameBetLogApi = {
    getGameBetLogList: (params: GameBetLogQuery) => {
        return request.post<PageResponse<GameBetLogItem>>('/gameBetLog/cmsGameBetLogList', params)
    },
}

export default gameBetLogApi

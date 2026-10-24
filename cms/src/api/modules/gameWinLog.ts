import {request} from '../request'
import type {GameWinLogItem, GameWinLogQuery, PageResponse} from '@/types/api'

export const gameWinLogApi = {
    getGameWinLogList: (params: GameWinLogQuery) => {
        return request.post<PageResponse<GameWinLogItem> & {platforms?: string[]}>('/gameWinLog/cmsGameWinLogList', params)
    },
}

export default gameWinLogApi

import {request} from '../request'
import type {GameCfg, GameCfgQuery} from '@/types/api'

export const gameCfgApi = {
    getGameCfgList: (params: GameCfgQuery) => {
        return request.post<{ total: number; data: GameCfg[] }>('/gameCfg/gameCfgList', params)
    },

    createGameCfg: (data: {
        name: string
        code: string
        liveCover: string
        link?: string
        sort?: number
        status: number
    }) => {
        return request.post<{ id: string }>('/gameCfg/createGameCfg', data)
    },

    updateGameCfg: (data: {
        id: string | number
        name: string
        code: string
        liveCover: string
        link?: string
        sort?: number
        status: number
    }) => {
        return request.post<{ success: boolean }>('/gameCfg/updateGameCfg', data)
    },

    deleteGameCfg: (id: string | number) => {
        return request.post<{ success: boolean }>('/gameCfg/deleteGameCfg', {id})
    },
}

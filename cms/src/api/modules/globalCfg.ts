import {request} from '../request'
import type {DelGlobalCfgReq, GlobalCfg, PageResponse, SaveGlobalCfgReq} from '@/types/api'

const globalCfgApi = {
    // 获取全局配置
    getGlobalCfg: (data?: any) => {
        return request.post<PageResponse<GlobalCfg>>('/globalCfg/getGlobalCfg', data)
    },

    // 保存全局配置
    saveGlobalCfg: (data: SaveGlobalCfgReq) => {
        return request.post<boolean>('/globalCfg/saveGlobalCfg', data)
    },

    // 删除全局配置
    delGlobalCfg: (data: DelGlobalCfgReq) => {
        return request.post<boolean>('/globalCfg/delGlobalCfg', data)
    }
}

export default globalCfgApi

import request from '../request'
import type {
    AddGameShelfReq,
    AddGameShelfRes,
    BatchAddGameShelfReq,
    BatchAddGameShelfRes,
    BatchDeleteGameShelfReq,
    BatchDeleteGameShelfRes,
    DeleteGameShelfReq,
    DeleteGameShelfRes,
    GetGamePlatformCfgRes,
    ReloadVendorGameCacheRes,
    SaveGamePlatformCfgReq,
    SaveGamePlatformCfgRes,
    VendorGame,
    VendorGameQuery,
    GameShelfItem,
    GameShelfQuery,
    UpdateGameShelfReq,
    UpdateGameShelfRes,
    GetMultiplayerConfigUrlReq,
    GetMultiplayerConfigUrlRes,
    CMSGameStartLinkReq,
    CMSGameStartLinkRes,
} from '@/types/api'

export const gamePlatformApi = {
    getGamePlatformCfg: () => {
        return request.post<GetGamePlatformCfgRes>('/gamePlatform/getGamePlatformCfg', {})
    },

    saveGamePlatformCfg: (data: SaveGamePlatformCfgReq) => {
        return request.post<SaveGamePlatformCfgRes>('/gamePlatform/saveGamePlatformCfg', data)
    },

    getVendorGameList: (params: VendorGameQuery) => {
        return request.post<{ total: number; data: VendorGame[] }>('/gamePlatform/vendorGameList', params)
    },

    getGameShelfList: (params: GameShelfQuery) => {
        return request.post<{ total: number; data: GameShelfItem[] }>('/gamePlatform/gameShelfList', params)
    },

    reloadVendorGameCache: () => {
        return request.post<ReloadVendorGameCacheRes>('/gamePlatform/reloadVendorGameCache', {})
    },

    addGameShelf: (data: AddGameShelfReq) => {
        return request.post<AddGameShelfRes>('/gamePlatform/addGameShelf', data)
    },

    deleteGameShelf: (data: DeleteGameShelfReq) => {
        return request.post<DeleteGameShelfRes>('/gamePlatform/deleteGameShelf', data)
    },

    batchAddGameShelf: (data: BatchAddGameShelfReq) => {
        return request.post<BatchAddGameShelfRes>('/gamePlatform/batchAddGameShelf', data)
    },

    batchDeleteGameShelf: (data: BatchDeleteGameShelfReq) => {
        return request.post<BatchDeleteGameShelfRes>('/gamePlatform/batchDeleteGameShelf', data)
    },

    updateGameShelf: (data: UpdateGameShelfReq) => {
        return request.post<UpdateGameShelfRes>('/gamePlatform/updateGameShelf', data)
    },

    getMultiplayerConfigUrl: (data: GetMultiplayerConfigUrlReq) => {
        return request.post<GetMultiplayerConfigUrlRes>('/gamePlatform/getMultiplayerConfigUrl', data)
    },

    getCMSGameStartLink: (data: CMSGameStartLinkReq) => {
        return request.post<CMSGameStartLinkRes>('/gamePlatform/cmsGameStartLink', data)
    },
}

export default gamePlatformApi

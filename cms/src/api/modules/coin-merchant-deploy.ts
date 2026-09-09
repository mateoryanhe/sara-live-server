import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface CoinMerchantDeployInfo {
    id: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    deploySecret: string
    updatedAt: string
}

export interface DeployCoinMerchantZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetCoinMerchantDeployInfoRes {
    info: CoinMerchantDeployInfo | null
}

export interface SaveCoinMerchantDeployCfgReq {
    id: number
    deploySecret: string
}

export interface SaveCoinMerchantDeployCfgRes {
    success: boolean
    id: string
}

export const coinMerchantDeployApi = {
    getCoinMerchantDeployInfo: () => {
        return request.post<GetCoinMerchantDeployInfoRes>('/coinMerchantDeploy/getCoinMerchantDeployInfo', {})
    },

    saveCoinMerchantDeployCfg: (data: SaveCoinMerchantDeployCfgReq) => {
        return request.post<SaveCoinMerchantDeployCfgRes>('/coinMerchantDeploy/saveCoinMerchantDeployCfg', data)
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployCoinMerchantZipRes>('/coinMerchantDeploy/deployZip', formData, {
            timeout: 0,
            onUploadProgress: (event: AxiosProgressEvent) => {
                if (!onUploadProgress || !event.total) {
                    return
                }
                onUploadProgress(Math.min(100, Math.round((event.loaded * 100) / event.total)))
            },
        })
    },
}

export default coinMerchantDeployApi

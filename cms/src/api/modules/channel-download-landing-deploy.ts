import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface ChannelDownloadLandingDeployInfo {
    id: string
    domain: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    updatedAt: string
    lastUploadAt: string
}

export interface DeployChannelDownloadLandingZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetChannelDownloadLandingDeployInfoRes {
    info: ChannelDownloadLandingDeployInfo | null
}

export interface SaveChannelDownloadLandingDeployCfgReq {
    id: number
    domain: string
    deployPath: string
}

export interface SaveChannelDownloadLandingDeployCfgRes {
    success: boolean
    id: string
}

export const channelDownloadLandingDeployApi = {
    getChannelDownloadLandingDeployInfo: () => {
        return request.post<GetChannelDownloadLandingDeployInfoRes>('/channelDownloadLandingDeploy/getChannelDownloadLandingDeployInfo', {})
    },

    saveChannelDownloadLandingDeployCfg: (data: SaveChannelDownloadLandingDeployCfgReq) => {
        return request.post<SaveChannelDownloadLandingDeployCfgRes>('/channelDownloadLandingDeploy/saveChannelDownloadLandingDeployCfg', data)
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployChannelDownloadLandingZipRes>('/channelDownloadLandingDeploy/deployZip', formData, {
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

export default channelDownloadLandingDeployApi

import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface ShareActivityDeployInfo {
    id: string
    domain: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    updatedAt: string
    lastUploadAt: string
}

export interface DeployShareActivityZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetShareActivityDeployInfoRes {
    info: ShareActivityDeployInfo | null
}

export interface SaveShareActivityDeployCfgReq {
    id: number
    domain: string
    deployPath: string
}

export interface SaveShareActivityDeployCfgRes {
    success: boolean
    id: string
}

export const shareActivityDeployApi = {
    getShareActivityDeployInfo: () => {
        return request.post<GetShareActivityDeployInfoRes>('/shareActivityDeploy/getShareActivityDeployInfo', {})
    },

    saveShareActivityDeployCfg: (data: SaveShareActivityDeployCfgReq) => {
        return request.post<SaveShareActivityDeployCfgRes>('/shareActivityDeploy/saveShareActivityDeployCfg', data)
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployShareActivityZipRes>('/shareActivityDeploy/deployZip', formData, {
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

export default shareActivityDeployApi

import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface H5LiveDeployInfo {
    id: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    deploySecret: string
    updatedAt: string
}

export interface DeployH5LiveZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetH5LiveDeployInfoRes {
    info: H5LiveDeployInfo | null
}

export interface SaveH5LiveDeployCfgReq {
    id: number
    deploySecret: string
}

export interface SaveH5LiveDeployCfgRes {
    success: boolean
    id: string
}

export const h5LiveDeployApi = {
    getH5LiveDeployInfo: () => {
        return request.post<GetH5LiveDeployInfoRes>('/h5LiveDeploy/getH5LiveDeployInfo', {})
    },

    saveH5LiveDeployCfg: (data: SaveH5LiveDeployCfgReq) => {
        return request.post<SaveH5LiveDeployCfgRes>('/h5LiveDeploy/saveH5LiveDeployCfg', data)
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployH5LiveZipRes>('/h5LiveDeploy/deployZip', formData, {
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

export default h5LiveDeployApi

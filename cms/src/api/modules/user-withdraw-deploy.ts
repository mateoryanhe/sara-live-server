import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface UserWithdrawDeployInfo {
    domain: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    lastUploadAt: string
}

export interface DeployUserWithdrawZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetUserWithdrawDeployInfoRes {
    info: UserWithdrawDeployInfo | null
}

export interface SaveUserWithdrawDeployCfgReq {
    domain: string
    deployPath: string
}

export interface SaveUserWithdrawDeployCfgRes {
    success: boolean
}

export const userWithdrawDeployApi = {
    getUserWithdrawDeployInfo: () => {
        return request.post<GetUserWithdrawDeployInfoRes>('/userWithdrawDeploy/getUserWithdrawDeployInfo', {})
    },

    saveUserWithdrawDeployCfg: (data: SaveUserWithdrawDeployCfgReq) => {
        return request.post<SaveUserWithdrawDeployCfgRes>('/userWithdrawDeploy/saveUserWithdrawDeployCfg', data)
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployUserWithdrawZipRes>('/userWithdrawDeploy/deployZip', formData, {
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

export default userWithdrawDeployApi

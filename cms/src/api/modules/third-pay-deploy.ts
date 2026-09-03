import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface ThirdPayDeployInfo {
    urlPrefix: string
    deployPath: string
    acceptExt: string
}

export interface DeployThirdPayZipRes {
    fileCount: number
    dirCount: number
    deployPath: string
    urlPrefix: string
}

export interface GetThirdPayDeployInfoRes {
    info: ThirdPayDeployInfo | null
}

export const thirdPayDeployApi = {
    getThirdPayDeployInfo: () => {
        return request.post<GetThirdPayDeployInfoRes>('/thirdPayDeploy/getThirdPayDeployInfo', {})
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployThirdPayZipRes>('/thirdPayDeploy/deployZip', formData, {
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

export default thirdPayDeployApi

import type {AxiosProgressEvent} from 'axios'
import {request} from '../request'

export interface CountryFlagPreviewItem {
    code: string
    nameEn: string
    nameZh: string
    file: string
    icon: string
}

export interface CountryFlagDeployInfo {
    id: string
    version: string
    urlPrefix: string
    deployPath: string
    acceptExt: string
    updatedAt: string
    fileCount: number
    flags: CountryFlagPreviewItem[]
}

export interface DeployCountryFlagZipRes {
    version: string
    fileCount: number
    deployPath: string
    urlPrefix: string
    removed: number
}

export interface GetCountryFlagDeployInfoRes {
    info: CountryFlagDeployInfo | null
}

export const countryFlagDeployApi = {
    getCountryFlagDeployInfo: () => {
        return request.post<GetCountryFlagDeployInfoRes>('/countryFlagDeploy/getCountryFlagDeployInfo', {})
    },

    deployZip: (file: File, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<DeployCountryFlagZipRes>('/countryFlagDeploy/deployZip', formData, {
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

export default countryFlagDeployApi

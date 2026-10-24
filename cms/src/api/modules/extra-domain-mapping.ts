import {request} from '../request'

export interface ExtraDomainMapping {
    id: string
    domain: string
    root: string
    hasTls: boolean
    certFile: string
    keyFile: string
    updatedAt: string
}

export interface ListExtraDomainMappingsRes {
    items: ExtraDomainMapping[]
}

export const extraDomainMappingApi = {
    list: () => request.post<ListExtraDomainMappingsRes>('/domainSite/listExtraDomainMappings', {}),
    save: (data: { domain: string; root: string }) =>
        request.post<{ success: boolean; item: ExtraDomainMapping }>('/domainSite/saveExtraDomainMapping', data),
    saveWithTLS: (data: { domain: string; root: string; cert: File; key: File }, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('domain', data.domain)
        formData.append('root', data.root)
        formData.append('cert', data.cert)
        formData.append('key', data.key)
        return request.post<{ success: boolean; item: ExtraDomainMapping }>('/domainSite/saveExtraDomainMappingUpload', formData, {
            timeout: 120000,
            onUploadProgress: onUploadProgress
                ? (e) => {
                    if (e.total) {
                        onUploadProgress(Math.round((e.loaded / e.total) * 100))
                    }
                }
                : undefined,
        })
    },
    remove: (id: string) =>
        request.post<{ success: boolean }>('/domainSite/deleteExtraDomainMapping', {id}),
}

export default extraDomainMappingApi

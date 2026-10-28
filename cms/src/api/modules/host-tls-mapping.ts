import {request} from '../request'

export interface HostTlsMapping {
    id: string
    domain: string
    hasTls: boolean
    certFile: string
    keyFile: string
    updatedAt: string
}

export interface ListHostTlsMappingsRes {
    items: HostTlsMapping[]
}

export const hostTlsMappingApi = {
    list: () => request.post<ListHostTlsMappingsRes>('/domainSite/listHostTlsMappings', {}),
    saveWithTLS: (data: { domain: string; cert: File; key: File }, onUploadProgress?: (percent: number) => void) => {
        const formData = new FormData()
        formData.append('domain', data.domain)
        formData.append('cert', data.cert)
        formData.append('key', data.key)
        return request.post<{ success: boolean; item: HostTlsMapping }>('/domainSite/saveHostTlsMappingUpload', formData, {
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
        request.post<{ success: boolean }>('/domainSite/deleteHostTlsMapping', {id}),
}

export default hostTlsMappingApi

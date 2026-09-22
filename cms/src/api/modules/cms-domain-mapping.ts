import {request} from '../request'

export interface CMSDomainSiteMapping {
    id: string
    domain: string
    urlPrefix: string
    root: string
    updatedAt: string
}

export interface GetCMSDomainSiteMappingRes {
    mapping: CMSDomainSiteMapping
}

export interface SaveCMSDomainSiteMappingReq {
    domain: string
    root: string
}

export interface SaveCMSDomainSiteMappingRes {
    success: boolean
    mapping: CMSDomainSiteMapping
}

export const cmsDomainMappingApi = {
    get: () => request.post<GetCMSDomainSiteMappingRes>('/domainSite/getCMSDomainSiteMapping', {}),
    save: (data: SaveCMSDomainSiteMappingReq) => {
        return request.post<SaveCMSDomainSiteMappingRes>('/domainSite/saveCMSDomainSiteMapping', data)
    },
}

export default cmsDomainMappingApi

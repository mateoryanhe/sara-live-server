import request from '@/api/request'

export interface CMSDomainAliasItem {
    id: string
    domain: string
    root: string
    updatedAt?: string
}

export interface ListCMSDomainAliasRes {
    primaryDomain: string
    root: string
    items: CMSDomainAliasItem[]
}

export const cmsDomainAliasApi = {
    list: () => request.post<ListCMSDomainAliasRes>('/domainSite/listCMSDomainAlias', {}),
    add: (domain: string) => request.post<{ success: boolean; item: CMSDomainAliasItem }>('/domainSite/addCMSDomainAlias', {domain}),
    remove: (id: string) => request.post<{ success: boolean }>('/domainSite/deleteCMSDomainAlias', {id}),
}

export default cmsDomainAliasApi

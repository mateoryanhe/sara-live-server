import {request} from '../request'
import type {PageQuery, PageResponse} from '@/types/api'

// CMS用户相关类型定义
export interface CMSUser {
    id: string
    name: string
    pwd: string
    status: number
    admin: boolean
    roleId: string
    createdAt: string
    updatedAt: string
}

export interface CMSUserQuery extends PageQuery {
    name?: string
    status?: number
    admin?: boolean
}

// CMS用户管理API
export const cmsUserApi = {
    // 获取CMS用户列表
    getCMSUserList: (params: CMSUserQuery) => {
        return request.post<PageResponse<CMSUser>>('/cmsuser/cmsUserList', params)
    },

    // 创建CMS用户
    createCMSUser: (data: {
        name: string
        pwd: string
        status: number
        admin: boolean
        roleId: string
    }) => {
        return request.post<boolean>('/cmsuser/createCMSUser', data)
    },

    // 更新CMS用户
    updateCMSUser: (data: {
        id: string
        name: string
        pwd?: string
        status: number
        admin: boolean
        roleId: string
    }) => {
        return request.post<boolean>('/cmsuser/updateCMSUser', data)
    },

    // 删除CMS用户
    deleteCMSUser: (id: string) => {
        return request.post<boolean>('/cmsuser/deleteCMSUser', {id})
    },
}
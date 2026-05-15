import {request} from '../request'
import type {PageQuery, PageResponse} from '@/types/api'

// 角色相关类型定义
export interface Role {
    id: string
    name: string
    description: string
    status: number
    createdAt: string
    updatedAt: string
}

export interface RolePermission {
    roleId: string
    permissions: string[]
}

export interface Permission {
    id: number
    module: string
    roleId: number
}

export interface RoleQuery extends PageQuery {
    name?: string
}


// 角色管理API
export const roleApi = {
    // 获取角色列表
    getRoleList: (params: RoleQuery) => {
        return request.post<PageResponse<Role>>('/role/roleList', params)
    },

    // 获取所有角色列表（用于下拉选择）
    getAllRoles: () => {
        return request.post<Role[]>('/role/roleList', {pageIndex: 1, pageSize: 9999})
    },

    // 创建角色
    createRole: (data: {
        name: string
        description: string
        status: number
        permissions: string[]
    }) => {
        return request.post<boolean>('/role/createRole', data)
    },

    // 更新角色
    updateRole: (data: {
        id: string
        name: string
        description: string
        status: number
        permissions?: string[]
    }) => {
        return request.post<boolean>('/role/updateRole', data)
    },

    // 删除角色
    deleteRole: (id: string) => {
        return request.post<boolean>('/role/deleteRole', {id})
    },


    // 获取角色权限列表
    getRolePermissionList: (roleId: number) => {
        return request.post<Permission[]>('/role/permissionList', {roleId})
    },

    // 创建或更新角色权限
    createPermission: (data: Permission[]) => {
        return request.post<boolean>('/role/createPermission', {data})
    },

}


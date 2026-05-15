import type {Permission} from '@/types/api'

// 存储用户权限信息
let userPermissions: Permission[] = []
let isAdmin = false

/**
 * 设置用户权限信息
 * @param modules 用户拥有的模块权限列表
 * @param admin 是否为管理员
 */
export const setUserPermissions = (modules: Permission[], admin: boolean) => {
    userPermissions = modules || []
    isAdmin = admin
}

/**
 * 检查用户是否有访问指定模块的权限
 * @param moduleName 模块名称
 * @returns 是否有权限
 */
export const hasPermission = (moduleName: string): boolean => {
    // 管理员拥有所有权限
    if (isAdmin) {
        return true
    }

    // 检查是否在权限列表中
    return userPermissions.some(module => module.module === moduleName)
}

/**
 * 获取用户权限列表
 * @returns 用户权限列表
 */
export const getUserPermissions = (): Permission[] => {
    return userPermissions
}

/**
 * 检查是否为管理员
 * @returns 是否为管理员
 */
export const getIsAdmin = (): boolean => {
    return isAdmin
}

/**
 * 清除权限信息
 */
export const clearPermissions = () => {
    userPermissions = []
    isAdmin = false
}
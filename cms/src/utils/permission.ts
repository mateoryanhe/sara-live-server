import type {Permission} from '@/types/api'
import {
    getPageFromPermissionKey,
    isButtonPermissionKey,
    buttonPermissionKey,
    pageHasGranularButtons,
    getPageButtons,
} from '@/config/page-buttons'

// 存储用户权限信息
let userPermissions: Permission[] = []
let isAdmin = false
let isSuperAdmin = false

/** 已授权的 module 字符串集合，便于快速查找 */
let permissionModuleSet = new Set<string>()

function rebuildPermissionSet() {
    permissionModuleSet = new Set(
        userPermissions.map(p => p.module).filter(Boolean),
    )
}

/**
 * 设置用户权限信息
 */
export const setUserPermissions = (modules: Permission[], admin: boolean, superAdmin = false) => {
    userPermissions = modules || []
    isAdmin = admin
    isSuperAdmin = superAdmin
    rebuildPermissionSet()
}

/** 是否拥有整页权限（module 恰好等于 pageName） */
export const hasFullPagePermission = (pageName: string): boolean => {
    if (isAdmin) {
        return true
    }
    return permissionModuleSet.has(pageName)
}

/** module 是否为该页面在 page-buttons 中已定义的按钮权限 */
function isKnownPageButton(pageName: string, module: string): boolean {
    if (!isButtonPermissionKey(module)) {
        return false
    }
    const page = getPageFromPermissionKey(module)
    if (page !== pageName) {
        return false
    }
    const action = module.slice(page.length + 1)
    return getPageButtons(pageName).some(btn => btn.key === action)
}

/**
 * 检查用户是否有访问指定页面的权限
 * 拥有整页权限，或拥有该页已定义按钮权限之一，均可进入页面
 */
export const hasPermission = (pageName: string): boolean => {
    if (isAdmin) {
        return true
    }
    if (!pageHasGranularButtons(pageName) && permissionModuleSet.has(pageName)) {
        return true
    }
    for (const module of permissionModuleSet) {
        if (isKnownPageButton(pageName, module)) {
            return true
        }
    }
    return false
}

/**
 * 检查页面内按钮权限
 * @param pageName 页面路由 name
 * @param action 按钮 key（page-buttons 中定义）
 */
export const hasButtonPermission = (pageName: string, action: string): boolean => {
    if (action === 'sync' && isAdmin && !isSuperAdmin) {
        return false
    }
    if (isAdmin) {
        return true
    }
    if (permissionModuleSet.has(buttonPermissionKey(pageName, action))) {
        return true
    }
    if (!pageHasGranularButtons(pageName) && permissionModuleSet.has(pageName)) {
        return true
    }
    return false
}

/**
 * 解析 v-btn-permission 绑定值
 * 支持 'PageName:action' 或 'action'（需配合 binding.arg 传 pageName）
 */
export const resolveButtonPermission = (
    value: string,
    argPage?: string,
): {page: string; action: string} | null => {
    if (!value) {
        return null
    }
    if (value.includes(':')) {
        const page = getPageFromPermissionKey(value)
        const action = value.slice(page.length + 1)
        return {page, action}
    }
    if (argPage) {
        return {page: argPage, action: value}
    }
    return null
}

export const getUserPermissions = (): Permission[] => userPermissions

export const getIsAdmin = (): boolean => isAdmin

export const getIsSuperAdmin = (): boolean => isSuperAdmin

export const clearPermissions = () => {
    userPermissions = []
    isAdmin = false
    isSuperAdmin = false
    permissionModuleSet.clear()
}

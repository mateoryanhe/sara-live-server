import type {RouteRecordRaw} from 'vue-router'
import {layoutRouteGroups} from '@/router/routes'
import {getIsAdmin, hasPermission} from '@/utils/permission'

export interface AccessibleRoute {
    name: string
    path: string
    title: string
}

function buildChildPath(group: RouteRecordRaw, child: RouteRecordRaw): string {
    const groupPath = String(group.path || '').replace(/\/$/, '')
    const childPath = String(child.path || '').replace(/^\//, '')
    if (!childPath) {
        return groupPath || '/'
    }
    return `${groupPath}/${childPath}`.replace(/\/+/g, '/')
}

/** 收集当前用户可访问的业务路由（按菜单顺序） */
export function collectAccessibleRoutes(): AccessibleRoute[] {
    const routes: AccessibleRoute[] = []

    layoutRouteGroups.forEach(group => {
        group.children?.forEach(child => {
            if (!child.name || child.meta?.hidden) {
                return
            }
            const name = String(child.name)
            if (!hasPermission(name)) {
                return
            }
            routes.push({
                name,
                path: buildChildPath(group, child),
                title: String(child.meta?.title || name),
            })
        })
    })

    return routes
}

function findRouteByPath(path: string): AccessibleRoute | undefined {
    const normalized = path.replace(/\/+$/, '') || '/'
    return collectAccessibleRoutes().find(route => route.path === normalized)
}

/** 登录后默认进入的页面：有 Dashboard 权限则进仪表盘，否则进第一个有权限的页面 */
export function getDefaultHomePath(): string {
    if (getIsAdmin() || hasPermission('Dashboard')) {
        return '/dashboard'
    }
    const accessible = collectAccessibleRoutes()
    return accessible[0]?.path || '/login'
}

/** 将目标路径解析为用户可访问的路径 */
export function resolveAccessiblePath(path: string): string {
    const normalized = path.replace(/\/+$/, '') || '/'
    if (normalized === '/' || normalized === '/dashboard') {
        return getDefaultHomePath()
    }

    const matched = findRouteByPath(normalized)
    if (matched) {
        return matched.path
    }

    return getDefaultHomePath()
}

/** 用户是否至少有一个可访问页面 */
export function hasAnyAccessibleRoute(): boolean {
    return getIsAdmin() || collectAccessibleRoutes().length > 0
}

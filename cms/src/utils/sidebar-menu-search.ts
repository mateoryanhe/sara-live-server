import type {RouteRecordRaw} from 'vue-router'
import {layoutRouteGroups} from '@/router/routes'
import {getIsAdmin, hasPermission} from '@/utils/permission'

export interface SidebarMenuSearchItem {
    path: string
    permission: string
    metaTitle: string
    group: string
}

function joinRoutePath(base: string, segment: string): string {
    if (!segment) {
        return base || '/'
    }
    if (segment.startsWith('/')) {
        return segment.replace(/\/+/g, '/')
    }
    const normalizedBase = base.endsWith('/') ? base.slice(0, -1) : base
    if (!normalizedBase) {
        return `/${segment}`.replace(/\/+/g, '/')
    }
    return `${normalizedBase}/${segment}`.replace(/\/+/g, '/')
}

function routeAccessible(route: RouteRecordRaw): boolean {
    if (getIsAdmin()) {
        return true
    }
    const name = route.name
    if (typeof name === 'string' && hasPermission(name)) {
        return true
    }
    const parent = route.meta?.parentPermission
    if (typeof parent === 'string' && parent && hasPermission(parent)) {
        return true
    }
    if (Array.isArray(parent)) {
        return parent.some(p => typeof p === 'string' && p && hasPermission(p))
    }
    return false
}

function collectSidebarMenuSearchItems(
    routes: RouteRecordRaw[],
    basePath: string,
    groupTitles: string[],
    out: SidebarMenuSearchItem[],
): void {
    for (const route of routes) {
        const path = joinRoutePath(basePath, route.path)
        const titlesForChildren = route.meta?.title
            ? [...groupTitles, String(route.meta.title)]
            : groupTitles

        if (route.children?.length) {
            collectSidebarMenuSearchItems(route.children, path, titlesForChildren, out)
            continue
        }
        if (!route.component || route.meta?.hidden) {
            continue
        }
        const permission = route.name
        if (typeof permission !== 'string' || !routeAccessible(route)) {
            continue
        }
        const metaTitle = route.meta?.title ? String(route.meta.title) : permission
        out.push({
            path,
            permission,
            metaTitle,
            group: groupTitles.join(' / '),
        })
    }
}

/** 可搜索的 CMS 菜单页（有权限、非 hidden 叶子路由） */
export function getSidebarMenuSearchItems(): SidebarMenuSearchItem[] {
    const items: SidebarMenuSearchItem[] = []
    collectSidebarMenuSearchItems(layoutRouteGroups, '', [], items)
    const seen = new Set<string>()
    return items.filter(item => {
        if (seen.has(item.path)) {
            return false
        }
        seen.add(item.path)
        return true
    })
}

/** 子序列 + 包含匹配，支持拼音首字母外的模糊定位 */
export function fuzzyMatchMenuQuery(query: string, haystack: string): boolean {
    const q = query.trim().toLowerCase()
    if (!q) {
        return true
    }
    const text = haystack.toLowerCase()
    if (text.includes(q)) {
        return true
    }
    let qi = 0
    for (let i = 0; i < text.length && qi < q.length; i++) {
        if (text[i] === q[qi]) {
            qi++
        }
    }
    return qi === q.length
}

export function filterSidebarMenuSearchItems(
    items: SidebarMenuSearchItem[],
    query: string,
    labelOf: (item: SidebarMenuSearchItem) => string,
    limit = 30,
): SidebarMenuSearchItem[] {
    const q = query.trim()
    if (!q) {
        return items.slice(0, limit)
    }
    const scored: Array<{ item: SidebarMenuSearchItem; score: number }> = []
    for (const item of items) {
        const label = labelOf(item)
        const haystack = [label, item.metaTitle, item.group, item.path, item.permission].join(' ')
        if (!fuzzyMatchMenuQuery(q, haystack)) {
            continue
        }
        const lower = haystack.toLowerCase()
        const ql = q.toLowerCase()
        let score = 0
        if (label.toLowerCase().startsWith(ql)) {
            score += 100
        } else if (label.toLowerCase().includes(ql)) {
            score += 50
        } else if (lower.includes(ql)) {
            score += 20
        } else {
            score += 5
        }
        scored.push({item, score})
    }
    scored.sort((a, b) => b.score - a.score || a.item.path.localeCompare(b.item.path))
    return scored.slice(0, limit).map(row => row.item)
}

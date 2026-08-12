import {ref} from 'vue'
import type {RouteLocationNormalizedLoaded, Router} from 'vue-router'

export interface LayoutTab {
    path: string
    title: string
    name: string
}

const DASHBOARD_PATH = '/dashboard'

const tabs = ref<LayoutTab[]>([
    {
        path: DASHBOARD_PATH,
        title: 'Dashboard',
        name: 'Dashboard',
    },
])

function resolveTabTitle(route: RouteLocationNormalizedLoaded): string {
    const name = route.name?.toString()
    if (name) {
        return name
    }
    const title = route.meta?.title
    if (typeof title === 'string' && title) {
        return title
    }
    return route.path
}

function shouldTrackRoute(route: RouteLocationNormalizedLoaded): boolean {
    if (route.path === '/login' || route.meta?.hidden) {
        return false
    }

    return !!route.name && route.matched.some(record => record.components?.default)
}

export function useLayoutTabs() {
    function addTab(route: RouteLocationNormalizedLoaded) {
        if (!shouldTrackRoute(route)) {
            return
        }

        const exists = tabs.value.some(tab => tab.path === route.path)
        if (exists) {
            return
        }

        tabs.value.push({
            path: route.path,
            title: resolveTabTitle(route),
            name: String(route.name),
        })
    }

    function removeTab(path: string, router: Router) {
        if (path === DASHBOARD_PATH) {
            return
        }

        const index = tabs.value.findIndex(tab => tab.path === path)
        if (index === -1) {
            return
        }

        tabs.value.splice(index, 1)

        if (router.currentRoute.value.path === path) {
            const nextTab = tabs.value[index] || tabs.value[index - 1] || tabs.value[0]
            if (nextTab) {
                void router.push(nextTab.path)
            }
        }
    }

    function closeOtherTabs(keepPath: string, router?: Router) {
        tabs.value = tabs.value.filter(
            tab => tab.path === keepPath || tab.path === DASHBOARD_PATH,
        )

        if (router) {
            const currentPath = router.currentRoute.value.path
            if (currentPath !== keepPath && currentPath !== DASHBOARD_PATH) {
                void router.push(keepPath)
            }
        }
    }

    function closeAllTabs(router: Router) {
        tabs.value = tabs.value.filter(tab => tab.path === DASHBOARD_PATH)
        if (router.currentRoute.value.path !== DASHBOARD_PATH) {
            void router.push(DASHBOARD_PATH)
        }
    }

    return {
        tabs,
        addTab,
        removeTab,
        closeOtherTabs,
        closeAllTabs,
    }
}

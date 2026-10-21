import {nextTick, ref} from 'vue'

/** 侧栏打开新页面时，主内容区 loading（页签切换不启用） */
export const routeNavigating = ref(false)

let inflightDuringSidebarNav = 0
let sidebarNavGeneration = 0

/** 仅侧栏 menu 跳转时调用 */
export function startSidebarRouteLoading() {
    sidebarNavGeneration += 1
    inflightDuringSidebarNav = 0
    routeNavigating.value = true
}

export function currentSidebarNavGeneration(): number {
    return sidebarNavGeneration
}

/** @deprecated 仅保留兼容；请用 startSidebarRouteLoading */
export function startRouteNavigation() {
    startSidebarRouteLoading()
}

export function trackSidebarNavRequestStart(): number {
    if (!routeNavigating.value) {
        return 0
    }
    inflightDuringSidebarNav += 1
    return sidebarNavGeneration
}

export function trackSidebarNavRequestEnd(generation: number) {
    if (!routeNavigating.value || generation !== sidebarNavGeneration) {
        return
    }
    inflightDuringSidebarNav = Math.max(0, inflightDuringSidebarNav - 1)
    scheduleFinishSidebarRouteLoading()
}

function scheduleFinishSidebarRouteLoading() {
    if (!routeNavigating.value) {
        return
    }
    if (inflightDuringSidebarNav > 0) {
        return
    }
    void nextTick(() => {
        if (!routeNavigating.value || inflightDuringSidebarNav > 0) {
            return
        }
        routeNavigating.value = false
    })
}

/** 路由已进入目标页，等待 onMounted 内发起的首批请求 */
export function notifySidebarRouteEntered() {
    if (!routeNavigating.value) {
        return
    }
    void nextTick(() => {
        void nextTick(() => {
            scheduleFinishSidebarRouteLoading()
        })
    })
}

export function finishRouteNavigation() {
    inflightDuringSidebarNav = 0
    routeNavigating.value = false
}

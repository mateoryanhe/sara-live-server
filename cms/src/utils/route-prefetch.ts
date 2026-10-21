const prefetchedPaths = new Set<string>()

const routeChunkLoaders: Record<string, () => Promise<unknown>> = {
    '/log/live/live-record-list': () => import('@/views/log/live/live-record-list.vue'),
    '/log/live/revenue-log-list': () => import('@/views/log/live/revenue-log-list.vue'),
    '/log/live/daily-effective-live-list': () => import('@/views/log/live/daily-effective-live-list.vue'),
    '/log/live/weekly-unsettled-live-list': () => import('@/views/log/live/weekly-unsettled-live-list.vue'),
    '/log/live/video-call-log-list': () => import('@/views/log/live/video-call-log-list.vue'),
}

/** 预拉取路由懒加载 chunk，展开菜单或 hover 时调用 */
export function prefetchRouteChunk(path: string) {
    if (prefetchedPaths.has(path)) {
        return
    }
    const load = routeChunkLoaders[path]
    if (!load) {
        return
    }
    prefetchedPaths.add(path)
    void load()
}

export function prefetchLiveLogRouteChunks() {
    for (const path of Object.keys(routeChunkLoaders)) {
        prefetchRouteChunk(path)
    }
}

import authApi from '@/api/modules/auth'
import {applySessionModules} from '@/utils/auth'
import {
    getLayoutBootstrapPromise,
    layoutBootstrapping,
    layoutShellReady,
    layoutSidebarComponent,
    resetLayoutBootstrap,
    setLayoutBootstrapPromise,
} from '@/utils/layout-bootstrap-state'

export {
    layoutBootstrapping,
    layoutShellReady,
    layoutSidebarComponent,
    resetLayoutBootstrap,
} from '@/utils/layout-bootstrap-state'

export async function loadLayoutNavigation(): Promise<void> {
    const mod = await import('@/components/layout/LayoutSidebar.vue')
    layoutSidebarComponent.value = mod.default
}

export async function refreshCmsSessionPermissions(): Promise<void> {
    const res = await authApi.cmsSessionModules({})
    applySessionModules(res)
}

/** 进入 Layout 时调用；多次调用共用同一 Promise */
export function runLayoutBootstrap(): Promise<void> {
    if (layoutShellReady.value) {
        return Promise.resolve()
    }
    const existing = getLayoutBootstrapPromise()
    if (existing) {
        return existing
    }

    layoutBootstrapping.value = true
    const promise = (async () => {
        try {
            await Promise.all([loadLayoutNavigation(), refreshCmsSessionPermissions()])
            layoutShellReady.value = true
        } catch (err) {
            console.warn('layout bootstrap failed, using cached permissions', err)
            try {
                await loadLayoutNavigation()
            } catch {
                // ignore
            }
            layoutShellReady.value = true
        } finally {
            layoutBootstrapping.value = false
        }
    })()
    setLayoutBootstrapPromise(promise)
    return promise
}

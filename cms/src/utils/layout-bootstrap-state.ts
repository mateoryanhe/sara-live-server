import {ref, shallowRef, type Component} from 'vue'

/** 主框架初始化中：并行加载侧栏 chunk + 刷新权限 */
export const layoutBootstrapping = ref(false)

/** 侧栏 chunk 与权限均已就绪，可挂载导航 */
export const layoutShellReady = ref(false)

export const layoutSidebarComponent = shallowRef<Component | null>(null)

let bootstrapPromise: Promise<void> | null = null

export function getLayoutBootstrapPromise(): Promise<void> | null {
    return bootstrapPromise
}

export function setLayoutBootstrapPromise(p: Promise<void> | null): void {
    bootstrapPromise = p
}

export function resetLayoutBootstrap(): void {
    layoutBootstrapping.value = false
    layoutShellReady.value = false
    layoutSidebarComponent.value = null
    bootstrapPromise = null
}

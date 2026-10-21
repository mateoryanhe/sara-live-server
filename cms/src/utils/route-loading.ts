import {ref} from 'vue'

/** 路由切换中（含懒加载 chunk），用于主内容区 loading */
export const routeNavigating = ref(false)

export function startRouteNavigation() {
    routeNavigating.value = true
}

export function finishRouteNavigation() {
    routeNavigating.value = false
}

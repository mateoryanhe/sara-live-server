import {computed} from 'vue'
import {useRoute} from 'vue-router'
import {hasButtonPermission} from '@/utils/permission'

/**
 * 当前页面按钮权限
 * @param pageName 可选，默认取 route.name
 */
export function usePagePermission(pageName?: string) {
    const route = useRoute()
    const page = computed(() => pageName ?? String(route.name ?? ''))

    const can = (action: string): boolean => hasButtonPermission(page.value, action)

    return {pageName: page, can}
}

import {computed} from 'vue'
import {useRouter} from 'vue-router'
import {usePagePermission} from '@/composables/usePagePermission'

/** 页面内跳转用户详情（需配置 viewUserDetail 按钮权限） */
export function useUserDetailNav(pageName: string) {
    const router = useRouter()
    const {can} = usePagePermission(pageName)
    const canViewUserDetail = computed(() => can('viewUserDetail'))

    const openUserDetail = (userId: string | number | null | undefined) => {
        if (userId == null || userId === '') {
            return
        }
        router.push({
            name: 'UserDetail',
            query: {id: String(userId)},
        })
    }

    return {canViewUserDetail, openUserDetail}
}

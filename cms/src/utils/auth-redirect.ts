import router from '@/router'
import {ElMessage} from 'element-plus'
import {clearAuthSession} from '@/utils/auth'

let redirectingToLogin = false

export function redirectToLogin(message = '登录已失效，请重新登录'): void {
    if (redirectingToLogin) {
        return
    }
    redirectingToLogin = true
    clearAuthSession()

    const current = router.currentRoute.value
    if (current.name === 'Login') {
        redirectingToLogin = false
        return
    }

    ElMessage.warning(message)
    const query = current.path !== '/' && current.path !== '/login'
        ? {redirect: current.fullPath, reason: 'expired'}
        : {reason: 'expired'}

    void router.replace({path: '/login', query}).finally(() => {
        redirectingToLogin = false
    })
}

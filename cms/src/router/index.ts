import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router'
import {hasPermission} from '@/utils/permission'
import {clearAuthSession, isAuthenticated, restoreAuthSession} from '@/utils/auth'
import {getDefaultHomePath, hasAnyAccessibleRoute, resolveAccessiblePath} from '@/utils/accessible-route'
import {ElMessage} from 'element-plus'
import {layoutRouteGroups} from './routes'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/index.vue'),
        meta: {title: '登录', hidden: true},
    },
    {
        path: '/',
        component: () => import('@/views/layout/index.vue'),
        redirect: '/dashboard',
        children: layoutRouteGroups,
    },
]

const router = createRouter({
    history: createWebHistory('/cms/'),
    routes,
})

router.beforeEach((to, from, next) => {
    document.title = to.meta.title
        ? `后台管理系统 - ${to.meta.title}`
        : '后台管理系统'

    restoreAuthSession()

    if (to.name === 'Login') {
        if (isAuthenticated()) {
            const redirect = typeof to.query.redirect === 'string' ? to.query.redirect : '/dashboard'
            next(resolveAccessiblePath(redirect))
            return
        }
        next()
        return
    }

    if (!isAuthenticated()) {
        next({
            path: '/login',
            query: to.fullPath === '/' ? undefined : {redirect: to.fullPath},
        })
        return
    }

    // 根路径重定向到用户有权限的首页
    if (to.path === '/' || (to.path === '/dashboard' && to.name === 'Dashboard' && !hasPermission('Dashboard'))) {
        const home = getDefaultHomePath()
        if (to.path !== home) {
            next(home)
            return
        }
    }

    const moduleName = String(to.name)
    if (hasPermission(moduleName)) {
        next()
        return
    }
    const parentPermission = to.meta.parentPermission
    if (typeof parentPermission === 'string' && parentPermission && hasPermission(parentPermission)) {
        next()
        return
    }
    if (Array.isArray(parentPermission) && parentPermission.some(p => typeof p === 'string' && p && hasPermission(p))) {
        next()
        return
    }

    console.warn(`用户没有访问 ${moduleName} 模块的权限`)

    // 登录后跳转到无权限页面：改去第一个有权限的页面，避免误报错误
    if (from.name === 'Login' || from.name === undefined) {
        const fallback = getDefaultHomePath()
        if (fallback !== '/login' && to.path !== fallback) {
            next(fallback)
            return
        }
    }

    if (!hasAnyAccessibleRoute()) {
        ElMessage.error('您没有任何模块权限，请联系管理员')
        clearAuthSession()
        next({path: '/login'})
        return
    }

    ElMessage.error('您没有权限访问该模块')
    if (from.name && from.name !== 'Login') {
        next(false)
        return
    }
    const fallback = getDefaultHomePath()
    if (fallback !== '/login' && to.path !== fallback) {
        next(fallback)
        return
    }
    clearAuthSession()
    next({
        path: '/login',
        query: {redirect: to.fullPath},
    })
})

export default router

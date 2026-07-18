import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router'
import {hasPermission} from '@/utils/permission'
import {clearAuthSession, isAuthenticated, restoreAuthSession} from '@/utils/auth'
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
            next(redirect)
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

    const moduleName = String(to.name)
    if (hasPermission(moduleName)) {
        next()
        return
    }

    console.warn(`用户没有访问 ${moduleName} 模块的权限`)
    ElMessage.error('您没有权限访问该模块')
    if (from.name && from.name !== 'Login') {
        next(false)
        return
    }
    clearAuthSession()
    next({
        path: '/login',
        query: {redirect: to.fullPath},
    })
})

export default router

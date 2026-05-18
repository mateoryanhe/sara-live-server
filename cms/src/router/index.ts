import {createRouter, createWebHistory, type RouteRecordRaw} from 'vue-router'
import Layout from '@/views/layout/index.vue'
import Login from '@/views/login/index.vue'
import {hasPermission} from '@/utils/permission'
import {ElMessage} from 'element-plus'

// 定义路由配置
const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {title: '登录', hidden: true}
    },
    {
        path: '/',
        redirect: '/login', // 未登录用户重定向到登录页
        children: [
            // 移除仪表盘路由
        ]
    },
    // 账号管理路由
    {
        path: '/account',
        component: Layout,
        redirect: '/account/user-list',
        meta: {title: '账号管理', icon: 'User'},
        children: [
            {
                path: '/account/user-list',
                name: 'UserList',
                component: () => import('@/views/account/user-list.vue'),
                meta: {title: '用户列表'}
            },

            {
                path: '/account/ban-user',
                name: 'BanUser',
                component: () => import('@/views/account/ban-user.vue'),
                meta: {title: '封禁用户', hidden: true} // hidden: true 表示不在菜单中显示
            }
        ]
    },
    // 系统配置路由
    {
        path: '/config',
        component: Layout,
        redirect: '/config/global',
        meta: {title: '系统配置', icon: 'Setting'},
        children: [
            {
                path: '/config/global',
                name: 'GlobalConfig',
                component: () => import('@/views/config/global.vue'),
                meta: {title: '全局配置'}
            }
        ]
    },
    // 直播工会管理路由
    {
        path: '/guild',
        component: Layout,
        redirect: '/guild/guild-list',
        meta: {title: '直播工会', icon: 'OfficeBuilding'},
        children: [
            {
                path: '/guild/guild-list',
                name: 'GuildManagement',
                component: () => import('@/views/guild/guild-list.vue'),
                meta: {title: '工会管理'}
            }
        ]
    },
    // 角色权限管理路由
    {
        path: '/role',
        component: Layout,
        redirect: '/role/role-list',
        meta: {title: '角色权限', icon: 'Lock'},
        children: [
            {
                path: '/role/role-list',
                name: 'RoleManagement',
                component: () => import('@/views/role/role-list.vue'),
                meta: {title: '角色权限管理'}
            },
            {
                path: '/role/module-list',
                name: 'ModuleList',
                component: () => import('@/views/role/module-list.vue'),
                meta: {title: '模块管理'}
            },
            {
                path: '/role/cmsuser-list',
                name: 'CMSUserManagement',
                component: () => import('@/views/role/cmsuser-list.vue'),
                meta: {title: 'CMS用户管理'}
            }
        ]
    }
]

// 创建路由实例
const router = createRouter({
    history: createWebHistory('/cms/'),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    // 设置页面标题
    if (to.meta.title) {
        document.title = `后台管理系统 - ${to.meta.title}`
    } else {
        document.title = '后台管理系统'
    }
    if (to.name == 'Login') {
        console.log('Login-------------')
        next() // 跳转到登录页面
        return;
    }
    // 检查登录状态
    const token = localStorage.getItem('token')
    if (!token) {
        // 未登录且访问非登录页面，跳转到登录页
        next('/login')
        return;
    }
    const moduleName = String(to.name)
    if (hasPermission(moduleName)) {
        next() // 有权限，允许访问
    } else {
        // 没有权限，跳转到登录页面
        console.warn(`用户没有访问 ${moduleName} 模块的权限`)
        ElMessage.error('您没有权限访问该模块')
        next('/login') // 跳转到登录页面
    }
})

export default router
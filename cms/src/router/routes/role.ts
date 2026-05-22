import type {RouteRecordRaw} from 'vue-router'

/** views/role */
export const roleRoutes: RouteRecordRaw = {
    path: '/role',
    meta: {title: '角色权限', icon: 'Lock'},
    redirect: '/role/role-list',
    children: [
        {
            path: 'role-list',
            name: 'RoleManagement',
            component: () => import('@/views/role/role-list.vue'),
            meta: {title: '角色权限管理'},
        },
        {
            path: 'module-list',
            name: 'ModuleList',
            component: () => import('@/views/role/module-list.vue'),
            meta: {title: '模块管理', hidden: true},
        },
        {
            path: 'cmsuser-list',
            name: 'CMSUserManagement',
            component: () => import('@/views/role/cmsuser-list.vue'),
            meta: {title: 'CMS用户管理'},
        },
    ],
}

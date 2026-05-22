import type {RouteRecordRaw} from 'vue-router'
import {userRoutes} from './user'
import {operationRoutes} from './operation'
import {configRoutes} from './config'
import {roleRoutes} from './role'

/** 按 views 目录分类的业务路由分组 */
export const layoutRouteGroups: RouteRecordRaw[] = [
    userRoutes,
    operationRoutes,
    configRoutes,
    roleRoutes,
]

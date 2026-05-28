import type {RouteRecordRaw} from 'vue-router'
import {dashboardRoutes} from './dashboard'
import {userRoutes} from './user'
import {operationRoutes} from './operation'
import {liveRoutes} from './live'
import {logRoutes} from './log'
import {configRoutes} from './config'
import {roleRoutes} from './role'
import {shortVideoRoutes} from './shortvideo'
import {gameRoutes} from './game'

/** 按 views 目录分类的业务路由分组 */
export const layoutRouteGroups: RouteRecordRaw[] = [
    dashboardRoutes,
    userRoutes,
    operationRoutes,
    liveRoutes,
    logRoutes,
    shortVideoRoutes,
    gameRoutes,
    configRoutes,
    roleRoutes,
]

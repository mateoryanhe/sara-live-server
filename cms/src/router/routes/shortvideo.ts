import type {RouteRecordRaw} from 'vue-router'

/** views/shortvideo */
export const shortVideoRoutes: RouteRecordRaw = {
    path: '/shortvideo',
    meta: {title: '短视频', icon: 'VideoCamera'},
    redirect: '/shortvideo/short-video-list',
    children: [
        {
            path: 'short-video-list',
            name: 'ShortVideoManagement',
            component: () => import('@/views/shortvideo/short-video-list.vue'),
            meta: {title: '短视频管理'},
        },
        {
            path: 'short-video-cfg',
            name: 'ShortVideoCfgManagement',
            component: () => import('@/views/shortvideo/short-video-cfg.vue'),
            meta: {title: '短视频配置'},
        },
    ],
}

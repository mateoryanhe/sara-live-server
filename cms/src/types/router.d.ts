import 'vue-router'
import type {PageButtonDef} from '@/config/page-buttons'

declare module 'vue-router' {
    interface RouteMeta {
        title?: string
        icon?: string
        hidden?: boolean
        /** 页面可勾选的按钮权限（未配置则使用 page-buttons.ts 默认/覆盖） */
        buttons?: PageButtonDef[]
        currencyType?: number
    }
}

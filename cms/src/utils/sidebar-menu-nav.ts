/** 根据目标路径展开侧边栏子菜单（index 与 layout 中 el-sub-menu 一致） */
export function resolveSidebarSubmenusToOpen(path: string): string[] {
    const open = new Set<string>()

    const add = (...indices: string[]) => {
        for (const index of indices) {
            open.add(index)
        }
    }

    if (path.startsWith('/tracking/')) {
        add('/tracking')
    }
    if (
        path.startsWith('/user/account/')
        || path === '/user/recharge-order/recharge-order-list'
        || path.startsWith('/operation/recharge/')
        || path.startsWith('/operation/vip/')
        || path === '/user/simulator-device-whitelist'
        || path.startsWith('/user/currency-log/')
    ) {
        add('/user/account')
    }
    if (path.startsWith('/withdraw/')) {
        add('/withdraw')
    }
    if (
        path.startsWith('/operation/guild/')
        || path === '/user/anchor/one-to-one-room-list'
        || path.startsWith('/operation/customer-service/')
        || path === '/user/anchor/live-room-recycle-bin'
    ) {
        add('/anchor-guild')
    }
    if (path.startsWith('/log/live/')) {
        add('/log/live')
    }
    if (
        path.startsWith('/operation/salary/')
        || path === '/operation/guild/guild-transfer-list'
    ) {
        add('/log/settlement')
    }
    if (path.startsWith('/live/')) {
        add('/live')
    }
    if (path.startsWith('/game/') || path.startsWith('/log/game/')) {
        add('/game')
    }
    if (
        path.startsWith('/user/coin-merchant/')
        || path.startsWith('/operation/guild/coin-merchant')
        || path.startsWith('/operation/recharge/coin-merchant')
        || path.startsWith('/log/user/coin-merchant')
        || path.startsWith('/operation/salary/coin-merchant')
    ) {
        add('/coin-merchant-module')
    }
    if (path.startsWith('/shortvideo/')) {
        add('/shortvideo')
    }
    if (
        path.startsWith('/activity/')
        || path.startsWith('/operation/activity-message/')
    ) {
        add('/activity')
    }
    if (path.startsWith('/promotion/')) {
        add('/promotion')
    }
    if (
        path.startsWith('/operation/wallet/')
        || path.startsWith('/operation/recharge/payment')
        || path.startsWith('/operation/salary/anchor-')
    ) {
        add('/settlement')
    }
    if (
        path.startsWith('/user/anchor/')
        || path.startsWith('/user/bot-anchor/')
        || path === '/config/static-cache-cfg'
        || path === '/config/app-pkg-list'
        || path === '/config/app-version-cfg'
        || path === '/config/firebase'
        || path === '/config/privacy-policy'
        || path === '/config/sofie-privacy-policy'
        || path === '/config/audit-code'
        || path === '/config/user-feedback-record'
        || path === '/config/user-report-record'
    ) {
        add('/frontend-module', '/frontend-module/config')
    }
    if (
        path === '/config/h5-live-deploy'
        || path === '/config/coin-merchant-deploy'
        || path === '/config/third-pay-deploy'
        || path === '/config/user-withdraw-deploy'
        || path === '/config/official-site-deploy'
        || path === '/config/third-pay-official-site-deploy'
        || path === '/config/sofie-official-site-deploy'
    ) {
        add('/frontend-module', '/frontend-module/deploy')
    }
    if (path.startsWith('/role/') || path.includes('visibility-list')) {
        add('/role')
    }
    if (path.startsWith('/config/')) {
        add('/config')
        if (
            path.startsWith('/operation/random-nickname/')
            || path === '/config/app-token'
            || path === '/config/account-cfg'
            || path === '/config/server-runtime-cfg'
        ) {
            add('/config/group/basic')
        }
        if (
            path === '/config/simulator-cpu-keyword'
            || path === '/config/device-register-risk-cfg'
            || path === '/config/text-moderation'
        ) {
            add('/config/group/security')
        }
        if (
            path === '/config/google-play'
            || path === '/config/meta-pixel'
            || path === '/config/haipay'
            || path === '/config/cf-email'
            || path === '/config/upload-resource'
            || path === '/config/cms-domain-mapping'
            || path === '/config/extra-domain-mapping'
            || path === '/config/host-tls-mapping'
            || path === '/config/country-flag-deploy'
            || path === '/config/db-backup'
            || path === '/config/data-sync'
        ) {
            add('/config/group/platform')
        }
        if (path === '/config/resource-monitor' || path === '/config/server-log') {
            add('/config/group/ops')
        }
    }
    if (
        path.startsWith('/operation/guild/guild-profile')
        || path.startsWith('/operation/guild/guild-anchor-daily-live')
    ) {
        add('/guild-data')
    }

    return [...open]
}

export function isSidebarRoutableIndex(
    index: string,
    resolve: (to: string) => { matched: Array<{ components?: Record<string, unknown> }> },
): boolean {
    if (!index.startsWith('/')) {
        return false
    }
    const resolved = resolve(index)
    return resolved.matched.some(record => record.components?.default != null)
}

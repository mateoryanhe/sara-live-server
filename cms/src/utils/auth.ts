import type {Permission} from '@/types/api'
import {clearPermissions, setUserPermissions} from '@/utils/permission'

const TOKEN_KEY = 'token'
const AUTH_ID_KEY = 'authId'
const USERNAME_KEY = 'cms_saved_userName'
const PASSWORD_KEY = 'cms_saved_pwd'
const DISPLAY_NAME_KEY = 'username'
const PERMISSIONS_KEY = 'cms_user_permissions'
const ADMIN_KEY = 'cms_user_admin'

export interface SavedCredentials {
    userName: string
    pwd: string
}

export interface AuthSession {
    token: string
    authId: string
    admin: boolean
    modules: Permission[]
}

export function getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
}

export function getAuthId(): string | null {
    return localStorage.getItem(AUTH_ID_KEY)
}

export function isAuthenticated(): boolean {
    const token = getToken()
    if (!token) {
        clearPermissions()
        return false
    }

    const admin = localStorage.getItem(ADMIN_KEY) === 'true'
    const rawModules = localStorage.getItem(PERMISSIONS_KEY)
    if (admin) {
        return true
    }
    if (!rawModules) {
        clearAuthSession()
        return false
    }

    try {
        const modules = JSON.parse(rawModules) as Permission[]
        if (!Array.isArray(modules)) {
            clearAuthSession()
            return false
        }
        return true
    } catch {
        clearAuthSession()
        return false
    }
}

export function restoreAuthSession(): boolean {
    if (!isAuthenticated()) {
        return false
    }

    const admin = localStorage.getItem(ADMIN_KEY) === 'true'
    const rawModules = localStorage.getItem(PERMISSIONS_KEY)
    if (admin) {
        setUserPermissions([], true)
        return true
    }

    try {
        const modules = JSON.parse(rawModules || '[]') as Permission[]
        setUserPermissions(Array.isArray(modules) ? modules : [], false)
    } catch {
        clearPermissions()
        return false
    }
    return true
}

export function setAuthSession(session: AuthSession): void {
    localStorage.setItem(TOKEN_KEY, session.token)
    localStorage.setItem(AUTH_ID_KEY, session.authId)
    localStorage.setItem(ADMIN_KEY, session.admin ? 'true' : 'false')
    localStorage.setItem(PERMISSIONS_KEY, JSON.stringify(session.modules || []))
    setUserPermissions(session.modules || [], session.admin)
}

export function clearAuthSession(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(AUTH_ID_KEY)
    localStorage.removeItem(PERMISSIONS_KEY)
    localStorage.removeItem(ADMIN_KEY)
    clearPermissions()
}

export function saveLoginCredentials(userName: string, pwd: string): void {
    localStorage.setItem(USERNAME_KEY, userName)
    localStorage.setItem(PASSWORD_KEY, pwd)
    localStorage.setItem(DISPLAY_NAME_KEY, userName)
}

export function getSavedCredentials(): SavedCredentials {
    return {
        userName: localStorage.getItem(USERNAME_KEY) || '',
        pwd: localStorage.getItem(PASSWORD_KEY) || '',
    }
}

export function clearSavedCredentials(): void {
    localStorage.removeItem(USERNAME_KEY)
    localStorage.removeItem(PASSWORD_KEY)
    localStorage.removeItem(DISPLAY_NAME_KEY)
}

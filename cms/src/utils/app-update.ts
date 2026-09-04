import {ElMessageBox} from 'element-plus'
import {i18n} from '@/i18n'

const CHUNK_RELOAD_KEY = 'cms_chunk_reload_ts'
const VERSION_PROMPT_KEY = 'cms_version_prompt_ts'
const CHUNK_RELOAD_COOLDOWN_MS = 15_000
const VERSION_PROMPT_COOLDOWN_MS = 60_000
const VERSION_POLL_INTERVAL_MS = 5 * 60_000

type VersionPayload = {
  version?: string
}

function t(key: string): string {
  return String(i18n.global.t(key))
}

function currentAppVersion(): string {
  return String(import.meta.env.VITE_APP_VERSION || 'dev')
}

function versionUrl(): string {
  const base = import.meta.env.BASE_URL || '/'
  const normalized = base.endsWith('/') ? base : `${base}/`
  return `${normalized}version.json`
}

/** Soft reload only; never clears auth/localStorage login session. */
export function softReload(): void {
  window.location.reload()
}

export function reloadForChunkError(): boolean {
  const now = Date.now()
  const last = Number(sessionStorage.getItem(CHUNK_RELOAD_KEY) || 0)
  if (now - last < CHUNK_RELOAD_COOLDOWN_MS) {
    return false
  }
  sessionStorage.setItem(CHUNK_RELOAD_KEY, String(now))
  softReload()
  return true
}

function isChunkLoadError(err: unknown): boolean {
  if (!err) {
    return false
  }
  const message = err instanceof Error ? err.message : String(err)
  const name = err instanceof Error ? err.name : ''
  const text = `${name} ${message}`.toLowerCase()
  return (
      text.includes('failed to fetch dynamically imported module')
      || text.includes('error loading dynamically imported module')
      || text.includes('importing a module script failed')
      || text.includes('loading chunk')
      || text.includes('loading css chunk')
      || text.includes('chunkloaderror')
      || text.includes('unable to preload css')
  )
}

async function fetchRemoteVersion(): Promise<string | null> {
  try {
    const res = await fetch(`${versionUrl()}?t=${Date.now()}`, {
      cache: 'no-store',
      headers: {Accept: 'application/json'},
    })
    if (!res.ok) {
      return null
    }
    const data = (await res.json()) as VersionPayload
    const version = String(data?.version || '').trim()
    return version || null
  } catch {
    return null
  }
}

async function promptReloadForNewVersion(remoteVersion: string): Promise<void> {
  const now = Date.now()
  const last = Number(sessionStorage.getItem(VERSION_PROMPT_KEY) || 0)
  if (now - last < VERSION_PROMPT_COOLDOWN_MS) {
    return
  }
  sessionStorage.setItem(VERSION_PROMPT_KEY, String(now))

  try {
    await ElMessageBox.confirm(
        t('common.appUpdateMessage'),
        t('common.appUpdateTitle'),
        {
          type: 'info',
          confirmButtonText: t('common.appUpdateConfirm'),
          cancelButtonText: t('common.cancel'),
          closeOnClickModal: false,
        },
    )
    sessionStorage.removeItem(CHUNK_RELOAD_KEY)
    softReload()
  } catch {
    // user cancelled; keep session, check again later
    void remoteVersion
  }
}

export async function checkAppVersion(): Promise<void> {
  if (import.meta.env.DEV) {
    return
  }
  const local = currentAppVersion()
  if (!local || local === 'dev') {
    return
  }
  const remote = await fetchRemoteVersion()
  if (!remote || remote === local) {
    return
  }
  await promptReloadForNewVersion(remote)
}

export function setupAppUpdateWatch(): void {
  if (import.meta.env.DEV) {
    return
  }

  const runCheck = () => {
    void checkAppVersion()
  }

  window.setTimeout(runCheck, 8_000)
  window.setInterval(runCheck, VERSION_POLL_INTERVAL_MS)

  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
      runCheck()
    }
  })
  window.addEventListener('focus', runCheck)

  window.addEventListener('error', (event) => {
    const target = event.target as HTMLElement | null
    if (target && (target.tagName === 'SCRIPT' || target.tagName === 'LINK')) {
      reloadForChunkError()
      return
    }
    if (isChunkLoadError(event.error)) {
      reloadForChunkError()
    }
  }, true)

  window.addEventListener('unhandledrejection', (event) => {
    if (isChunkLoadError(event.reason)) {
      reloadForChunkError()
    }
  })
}

export function setupRouterChunkReload(router: {onError: (fn: (err: Error) => void) => void}): void {
  router.onError((err) => {
    if (isChunkLoadError(err)) {
      reloadForChunkError()
    }
  })
}

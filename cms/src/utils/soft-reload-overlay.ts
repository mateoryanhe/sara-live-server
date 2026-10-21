/** sessionStorage：整页 soft reload 期间显示 boot loading，避免白屏 */
export const SOFT_RELOAD_FLAG_KEY = 'cms_soft_reload'
export const SOFT_RELOAD_MSG_KEY = 'cms_soft_reload_msg'

const OVERLAY_ID = 'cms-boot-overlay'

const DEFAULT_OVERLAY_MSG = 'Loading…'

export function softReloadOverlayMessage(): string {
  try {
    return sessionStorage.getItem(SOFT_RELOAD_MSG_KEY) || DEFAULT_OVERLAY_MSG
  } catch {
    return DEFAULT_OVERLAY_MSG
  }
}

export function mountSoftReloadOverlay(message: string): void {
  if (typeof document === 'undefined') {
    return
  }
  let root = document.getElementById(OVERLAY_ID)
  if (!root) {
    root = document.createElement('div')
    root.id = OVERLAY_ID
    root.setAttribute('role', 'status')
    root.setAttribute('aria-live', 'polite')
    root.innerHTML = `
      <div class="cms-boot-overlay__panel">
        <div class="cms-boot-overlay__spinner" aria-hidden="true"></div>
        <p class="cms-boot-overlay__text"></p>
      </div>
    `
    document.body.appendChild(root)
  }
  const textEl = root.querySelector('.cms-boot-overlay__text')
  if (textEl) {
    textEl.textContent = message
  }
}

export function prepareSoftReload(message: string): void {
  try {
    sessionStorage.setItem(SOFT_RELOAD_FLAG_KEY, '1')
    sessionStorage.setItem(SOFT_RELOAD_MSG_KEY, message)
  } catch {
    // ignore private mode
  }
  mountSoftReloadOverlay(message)
}

export function clearSoftReloadOverlay(): void {
  try {
    sessionStorage.removeItem(SOFT_RELOAD_FLAG_KEY)
    sessionStorage.removeItem(SOFT_RELOAD_MSG_KEY)
  } catch {
    // ignore
  }
  document.getElementById(OVERLAY_ID)?.remove()
}

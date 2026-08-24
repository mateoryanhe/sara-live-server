const PASSWORD_CHARS = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789'

export function createRandomPassword(length = 20): string {
  const values = crypto.getRandomValues(new Uint32Array(length))
  return Array.from(values, (value) => PASSWORD_CHARS[value % PASSWORD_CHARS.length]).join('')
}

export async function copyTextToClipboard(value: string): Promise<boolean> {
  const text = value.trim()
  if (!text) {
    return false
  }
  await navigator.clipboard.writeText(text)
  return true
}

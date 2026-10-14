import {Converter} from 'opencc-js'

const toTraditional = Converter({from: 'cn', to: 'tw'})

export function convertTextToZhTw(text: string): string {
  return toTraditional(text)
}

export function convertLocaleDeep<T>(value: T): T {
  if (typeof value === 'string') {
    return convertTextToZhTw(value) as T
  }
  if (Array.isArray(value)) {
    return value.map((item) => convertLocaleDeep(item)) as T
  }
  if (value !== null && typeof value === 'object') {
    const out: Record<string, unknown> = {}
    for (const [key, nested] of Object.entries(value)) {
      out[key] = convertLocaleDeep(nested)
    }
    return out as T
  }
  return value
}

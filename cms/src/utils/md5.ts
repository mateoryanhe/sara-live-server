import md5 from 'js-md5'

/** 浏览器端 MD5(hex 小写),用于币商密码上报 */
export function md5Hex(input: string): string {
  return md5(input)
}

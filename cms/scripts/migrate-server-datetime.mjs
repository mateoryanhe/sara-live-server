import fs from 'fs'
import path from 'path'

const root = path.resolve('src')

const csvFiles = [
  'utils/income-settlement-log-csv.ts',
  'utils/currency-log-csv.ts',
  'utils/live-record-csv.ts',
  'utils/daily-effective-live-csv.ts',
]

for (const rel of csvFiles) {
  const file = path.join(root, rel)
  let s = fs.readFileSync(file, 'utf8')
  if (s.includes('server-datetime')) continue
  s = s.replace(
    /export function format\w+CsvDate\(dateString: string \| null \| undefined\): string \{[\s\S]*?\n\}\n\n?/,
    '',
  )
  if (!s.includes("from '@/utils/server-datetime'")) {
    s = s.replace(/^(import[^\n]+\n)/, "$1import {formatServerDateTimeForExport} from '@/utils/server-datetime'\n")
  }
  s = s.replace(/format\w+CsvDate\(/g, 'formatServerDateTimeForExport(')
  fs.writeFileSync(file, s)
  console.log('csv', rel)
}

const toDayBlock =
  /const toDayStartUnix = \(dateStr: string\): number => \{[\s\S]*?\n\}\n\nconst toDayEndUnix = \(dateStr: string\): number => \{[\s\S]*?\n\}\n\n?/m

const formatDateBlock =
  /const formatDate = \(dateString: string \| null \| undefined\) => \{[\s\S]*?\n\}\n\n/m

function walk(dir) {
  for (const name of fs.readdirSync(dir)) {
    const p = path.join(dir, name)
    if (fs.statSync(p).isDirectory()) walk(p)
    else if (name.endsWith('.vue')) updateVue(p)
  }
}

function insertImport(s, items) {
  const line = `import {${items.join(', ')}} from '@/utils/server-datetime'\n`
  if (s.includes('@/utils/server-datetime')) return s
  const m = s.match(/<script[^>]*>\n((?:import[^\n]+\n)*)/)
  if (!m) return s
  return s.replace(m[0], m[0] + line)
}

function updateVue(file) {
  let s = fs.readFileSync(file, 'utf8')
  const orig = s
  const imports = new Set()

  if (toDayBlock.test(s)) {
    s = s.replace(toDayBlock, '')
    s = s.replace(/\btoDayStartUnix\b/g, 'toServerDayStartUnix')
    s = s.replace(/\btoDayEndUnix\b/g, 'toServerDayEndUnix')
    imports.add('toServerDayStartUnix')
    imports.add('toServerDayEndUnix')
  }

  if (formatDateBlock.test(s)) {
    s = s.replace(formatDateBlock, '')
    imports.add('formatServerDateTime as formatDate')
  }

  if (s.includes('new Date(ts * 1000).toLocaleString()')) {
    s = s.replace(/new Date\(ts \* 1000\)\.toLocaleString\(\)/g, 'formatServerDateTime(ts)')
    imports.add('formatServerDateTime')
  }

  if (imports.size === 0) return

  s = insertImport(s, [...imports])
  if (orig !== s) {
    fs.writeFileSync(file, s)
    console.log('vue', path.relative(root, file))
  }
}

walk(root)

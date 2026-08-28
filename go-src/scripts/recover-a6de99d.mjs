import fs from 'fs'
import path from 'path'

const root = path.resolve(import.meta.dirname, '../..')
const transcript = 'C:/Users/hw/.cursor/projects/d-company-code-sara-live-server/agent-transcripts/d05d4488-6aa4-423c-b540-696d875df8cc/d05d4488-6aa4-423c-b540-696d875df8cc.jsonl'
const lines = fs.readFileSync(transcript, 'utf8').split('\n')

function toRel(filePath) {
  if (!filePath) return null
  const n = filePath.replace(/\\/g, '/')
  const marker = 'sara-live-server/'
  const idx = n.toLowerCase().indexOf(marker)
  if (idx < 0) return null
  return n.slice(idx + marker.length)
}

const start = 5997
const end = 6176
const writes = []
const replaces = []

for (let i = start; i <= end && i < lines.length; i++) {
  let o
  try {
    o = JSON.parse(lines[i])
  } catch {
    continue
  }
  for (const c of o.message?.content || []) {
    if (c.type !== 'tool_use') continue
    const input = c.input || {}
    const rel = toRel(input.path || '')
    if (!rel) continue
    if (c.name === 'Write') writes.push({ line: i, rel, contents: input.contents })
    if (c.name === 'StrReplace') replaces.push({ line: i, rel, old: input.old_string, new: input.new_string })
  }
}

const fileContents = {}

function loadFile(rel) {
  if (rel in fileContents) return fileContents[rel]
  const fp = path.join(root, rel)
  if (fs.existsSync(fp)) {
    fileContents[rel] = fs.readFileSync(fp, 'utf8')
    return fileContents[rel]
  }
  return null
}

for (const w of writes) {
  fileContents[w.rel] = w.contents
}

const touched = new Set(Object.keys(fileContents))
for (const r of replaces) touched.add(r.rel)

for (const rel of touched) loadFile(rel)

let applied = 0
let miss = 0
let nofile = 0

for (const r of replaces) {
  const cur = loadFile(r.rel)
  if (cur == null) {
    console.log('NOFILE', r.line, r.rel)
    nofile++
    continue
  }
  if (!cur.includes(r.old)) {
    console.log('MISS', r.line, r.rel)
    miss++
    continue
  }
  fileContents[r.rel] = cur.replace(r.old, r.new)
  applied++
}

for (const [rel, content] of Object.entries(fileContents)) {
  const fp = path.join(root, rel)
  fs.mkdirSync(path.dirname(fp), { recursive: true })
  fs.writeFileSync(fp, content)
}

console.log(`writes=${writes.length} replaces=${replaces.length} applied=${applied} miss=${miss} nofile=${nofile} saved=${Object.keys(fileContents).length}`)

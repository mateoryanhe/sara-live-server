import fs from 'node:fs'
import path from 'node:path'
import type {Plugin} from 'vite'

export type AppVersionPluginOptions = {
  /** Injected into import.meta.env.VITE_APP_VERSION */
  version: string
}

/** Emit version.json next to index.html and expose VITE_APP_VERSION. */
export function appVersionPlugin(options: AppVersionPluginOptions): Plugin {
  const version = options.version
  let outDir = ''

  return {
    name: 'cms-app-version',
    config() {
      return {
        define: {
          'import.meta.env.VITE_APP_VERSION': JSON.stringify(version),
        },
      }
    },
    configResolved(config) {
      outDir = path.resolve(config.root, config.build.outDir)
    },
    closeBundle() {
      if (!outDir) {
        return
      }
      fs.mkdirSync(outDir, {recursive: true})
      const filePath = path.join(outDir, 'version.json')
      fs.writeFileSync(filePath, `${JSON.stringify({version}, null, 2)}\n`, 'utf8')
    },
  }
}

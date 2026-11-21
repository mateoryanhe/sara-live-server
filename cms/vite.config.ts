import {fileURLToPath, URL} from 'node:url'

import {defineConfig, loadEnv} from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {ElementPlusResolver} from 'unplugin-vue-components/resolvers'
import {appVersionPlugin} from './vite-plugin-app-version'

// https://vite.dev/config/
export default defineConfig(({mode}) => {
    const env = loadEnv(mode, process.cwd(), '')
    const base = env.VITE_BASE_PATH || '/'
    const appVersion = `${mode}-${Date.now().toString(36)}`
    const goApiTarget = env.VITE_GO_API_TARGET || 'http://127.0.0.1:9443'
    const devApiViaProxy =
        mode === 'development' &&
        (env.VITE_API_BASE_URL === '/local-api' ||
            (env.VITE_API_BASE_URL || '').startsWith('/local-api'))

    return {
        base,
        plugins: [
            vue(),
            vueJsx(),
            ...(mode === 'development' ? [vueDevTools()] : []),
            AutoImport({
                resolvers: [ElementPlusResolver()],
            }),
            Components({
                resolvers: [ElementPlusResolver()],
            }),
            appVersionPlugin({version: appVersion}),
        ],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url))
            },
        },
        server: {
            // 0.0.0.0:5174 在部分 Windows 环境会与系统/其它进程冲突；本机开发默认 127.0.0.1。局域网调试可设 VITE_DEV_HOST=0.0.0.0
            host: env.VITE_DEV_HOST || '127.0.0.1',
            port: 5174,
            strictPort: true,
            open: true,
            proxy: devApiViaProxy
                ? {
                      '/local-api': {
                          target: goApiTarget,
                          changeOrigin: true,
                          rewrite: (path) => path.replace(/^\/local-api/, ''),
                      },
                  }
                : undefined,
        },
        build: {
            outDir: 'D:\\root\\cms',
            // outDir 在项目目录之外，Vite 默认不会清空；必须显式清理，避免旧 hash 资源持续累积进发版包。
            emptyOutDir: true,
            assetsDir: 'assets',
            sourcemap: false,
            minify: 'terser',
            rollupOptions: {
                output: {
                    chunkFileNames: 'assets/js/[name]-[hash].js',
                    entryFileNames: 'assets/js/[name]-[hash].js',
                    assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
                    manualChunks(id) {
                        if (!id.includes('node_modules')) {
                            return
                        }
                        // Only split echarts; splitting vue and element-plus causes circular chunk deps.
                        if (id.includes('/echarts/') || id.includes('/zrender/')) {
                            return 'echarts'
                        }
                    },
                }
            },
            chunkSizeWarningLimit: 1000
        }
    }
})

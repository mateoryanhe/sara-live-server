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
            host: '0.0.0.0',
            port: 5173,
            open: true,
        },
        build: {
            outDir: 'D:\\root\\cms',
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

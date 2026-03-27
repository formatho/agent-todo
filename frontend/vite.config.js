import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { visualizer } from 'rollup-plugin-visualizer'
import compression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    vue(),
    // PWA support for offline functionality and improved caching
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.svg'],
      manifest: {
        name: 'Agent Todo Platform',
        short_name: 'Agent Todo',
        description: 'AI-powered task management platform for teams and individuals',
        theme_color: '#3B82F6',
        background_color: '#ffffff',
        display: 'standalone',
        orientation: 'portrait',
        start_url: '/',
        icons: [
          {
            src: 'favicon.svg',
            sizes: 'any',
            type: 'image/svg+xml',
            purpose: 'any'
          }
        ]
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,jpg,jpeg,json}'],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\.formatho\.com\/.*/,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 3600
              }
            }
          },
          {
            urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'font-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 86400
              }
            }
          },
          {
            urlPattern: /^https:\/\/cdnjs\.cloudflare\.com\/.*/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'script-cache',
              expiration: {
                maxEntries: 50,
                maxAgeSeconds: 86400
              }
            }
          }
        ]
      }
    }),
    // Bundle analyzer for optimization insights
    visualizer({
      filename: 'bundle-analysis.html',
      open: false,
      gzipSize: true,
      brotliSize: true
    }),
    // Gzip compression
    compression({
      algorithm: 'gzip',
      ext: '.gz'
    }),
    // Brotli compression (better compression ratio than gzip)
    compression({
      algorithm: 'brotliCompress',
      ext: '.br'
    })
  ],
  build: {
    // Performance optimizations
    target: 'es2020',
    cssCodeSplit: true,
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          // Vendor chunk separation
          vue: ['vue', 'vue-router'],
          pinia: ['pinia'],
          utils: ['axios'],
          ui: ['vuedraggable']
        },
        // Optimize chunk names
        chunkFileNames: (chunkInfo) => {
          const facadeModuleId = chunkInfo.facadeModuleIdByFile
          if (facadeModuleId) {
            const fileName = facadeModuleId.split('/').pop()
            return `js/${fileName.replace(/\.\w+$/, '')}.js`
          }
          return 'js/[name]-[hash].js'
        },
        assetFileNames: (assetInfo) => {
          const info = assetInfo.name.split('.')
          const ext = info[info.length - 1]
          if (/\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/i.test(assetInfo.name)) {
            return `media/${assetInfo.name}`
          }
          if (/\.(png|jpe?g|gif|svg)(\?.*)?$/i.test(assetInfo.name)) {
            return `img/${assetInfo.name}`
          }
          if (ext === 'css') {
            return `css/${info[0]}.css`
          }
          return `assets/${assetInfo.name}`
        }
      }
    },
    // Minimization optimizations
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
        pure_funcs: ['console.log', 'console.info', 'console.warn']
      },
      format: {
        comments: false
      }
    }
  },
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        secure: false
      }
    },
    // Enable HTTP/2 for better performance
    https: false
  },
  // Performance monitoring
  define: {
    __VITE_OPTIMIZE__: true,
    __VITE_PERFORMANCE__: true
  },
  // Optimized preloading
  optimizeDeps: {
    include: [
      'vue',
      'vue-router',
      'pinia',
      'axios',
      'marked'
    ],
    exclude: ['@vue/runtime-core']
  }
})
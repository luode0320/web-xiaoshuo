import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { loadEnv } from 'vite'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  // 根据环境设置API基础URL
  let apiBaseUrl = 'http://localhost:8888'
  if (mode === 'production') {
    apiBaseUrl = 'https://xs.luode.vip'
  } else if (env.VITE_API_BASE_URL) {
    apiBaseUrl = env.VITE_API_BASE_URL
  }

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    server: {
      host: '0.0.0.0',
      port: 3000,
      proxy: {
        '/api': {
          target: apiBaseUrl,
          changeOrigin: true,
          secure: false,
        },
      },
    },
    build: {
      outDir: 'dist',
      assetsDir: 'assets',
      sourcemap: false,
      minify: 'terser',
      rollupOptions: {
        output: {
          manualChunks: {
            vue: ['vue', 'vue-router'],
            element: ['element-plus'],
            utils: ['axios', 'dayjs'],
          },
        },
      },
    },
    // 定义环境变量，构建时会嵌入到代码中
    define: {
      'process.env': {
        VUE_APP_API_BASE_URL: JSON.stringify(apiBaseUrl),
      }
    }
  }
})
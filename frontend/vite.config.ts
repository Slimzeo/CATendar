import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5179,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://localhost:18082',
        changeOrigin: true
      }
    }
  }
})

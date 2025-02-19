import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      'element-plus/lib': 'element-plus/es',
      'element-plus/lib/locale': 'element-plus/es/locale'
    }
  }
})

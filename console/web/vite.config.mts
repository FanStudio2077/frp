import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  base: '',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@shared': fileURLToPath(new URL('../../web/shared', import.meta.url)),
      '@frpc': fileURLToPath(new URL('../../web/frpc/src', import.meta.url)),
    },
    dedupe: ['vue', 'element-plus'],
  },
  build: { assetsDir: '', minify: false },
})

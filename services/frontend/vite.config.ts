import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Proxy API calls to catalog service
      '/api': {
        target: 'http://catalog.kubelab.lan:8081',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})

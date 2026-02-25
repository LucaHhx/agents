import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const host = process.env.TAURI_DEV_HOST

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    host: host || false,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})

import vue from '@vitejs/plugin-vue'
import laravel from 'laravel-vite-plugin'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    laravel({
      input: ['frontend/app.ts', 'frontend/app.scss'],
      ssr: 'frontend/ssr.ts', // Enable SSR
      publicDirectory: 'public',
      buildDirectory: 'bootstrap',
      refresh: true,
    }),
    vue(),
  ],
  build: {
    ssr: true, // Enable SSR
    outDir: 'bootstrap',
    rollupOptions: {
      input: 'frontend/ssr.ts',
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]',
        manualChunks: undefined, // Disable automatic chunk splitting
      },
    },
  },
})
